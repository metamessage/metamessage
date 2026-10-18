namespace MetaMessage.Core;

/// <summary>
/// v0.5.0 尾定界辅助（与 Go internal/core/footer.go + tailframe.go 逐字节对齐）：
///
///   1. meta_len 尾部槽编解码（FooterSlotBytes / ParseFooterSlot）+ 字节逆序（ReverseBytes）。
///   2. 解码入口剥帧 Untailify（NodeTailStart / UntailifyNode / UntailifyChildren）：
///      从节点尾 sign 反向定界，把尾定界 wire 还原成「sign 在前」的头部形态，
///      再交给原有正向解码器。
///
/// 编码器各原语已直接按尾定界顺序写出（body→lenBE→sign），本类不含编码侧换位逻辑，
/// 也不设独立 payload_len 槽：各类型 sign 低 5 位的档位/内联语义（v0.3.0 原样）同时承担定界信息。
///
/// meta_len 尾槽三档（canonical 最小宽度，末字节为标记或直接值）：
///   - metaLen <= 253：[metaLen]                     槽宽 1
///   - metaLen <= 255：[metaLen][254]                槽宽 2
///   - metaLen <= 65535：[hi][lo][255]（u16 大端）    槽宽 3
/// </summary>
public class Footer
{
    public const int SlotDirectMax = 253; // 一档可直写的 meta 长度上限（尾字节 < 254）
    public const int SlotMarker1 = 254;   // 二档标记：倒数第 2 字节为 u8 长度
    public const int SlotMarker2 = 255;   // 三档标记：倒数第 3..2 字节为 u16 大端长度

    // FooterSlotBytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
    public static byte[] FooterSlotBytes(int metaLen)
    {
        if (metaLen <= SlotDirectMax)
        {
            return new byte[] { (byte)metaLen };
        }
        else if (metaLen <= WireConstants.MAX_1)
        {
            return new byte[] { (byte)metaLen, (byte)SlotMarker1 };
        }
        else if (metaLen <= WireConstants.MAX_2)
        {
            return new byte[] { (byte)((metaLen >> 8) & 0xFF), (byte)(metaLen & 0xFF), (byte)SlotMarker2 };
        }
        throw new MmDecodeException("tag too long, max length: " + WireConstants.MAX_2 + ", actual: " + metaLen);
    }

    // ParseFooterSlot 从节点尾端反向解析长度槽，返回 [metaLen, slotWidth]。
    public static (int MetaLen, int SlotWidth) ParseFooterSlot(byte[] data, int nodeEnd)
    {
        if (nodeEnd < 1)
        {
            throw new MmDecodeException("footer slot out of range");
        }
        int marker = data[nodeEnd - 1];
        if (marker == SlotMarker1)
        {
            if (nodeEnd < 2)
            {
                throw new MmDecodeException("footer slot truncated");
            }
            return (data[nodeEnd - 2], 2);
        }
        else if (marker == SlotMarker2)
        {
            if (nodeEnd < 3)
            {
                throw new MmDecodeException("footer slot truncated");
            }
            return ((data[nodeEnd - 3] << 8) | data[nodeEnd - 2], 3);
        }
        return (marker, 1);
    }

    // ReverseBytes 返回 bs 的字节逆序拷贝。
    public static byte[] ReverseBytes(byte[] bs)
    {
        var outBytes = new byte[bs.Length];
        for (int i = 0; i < bs.Length; i++)
        {
            outBytes[i] = bs[bs.Length - 1 - i];
        }
        return outBytes;
    }

    // MetaPrefixWidth 正向 meta 前缀槽宽度：与尾部槽同一套 canonical 三档，
    // 宽度恒等于 FooterSlotBytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
    public static int MetaPrefixWidth(int ml)
    {
        if (ml < 254) return 1;
        if (ml <= 255) return 2;
        return 3;
    }

    // MetaPrefixBytes 还原正向 meta 前缀槽。
    public static byte[] MetaPrefixBytes(int ml)
    {
        if (ml < 254)
        {
            return new byte[] { (byte)ml };
        }
        if (ml <= WireConstants.MAX_1)
        {
            return new byte[] { 254, (byte)ml };
        }
        if (ml <= WireConstants.MAX_2)
        {
            return new byte[] { 255, (byte)((ml >> 8) & 0xFF), (byte)(ml & 0xFF) };
        }
        throw new MmDecodeException("meta too long: " + ml);
    }

    // ----- 各类型档位：返回 (w, inline)，w==0 表示内联档 -----

    private static (int W, int Inline) IntLen(int b)
    {
        int l = b & WireConstants.INT_LEN_MASK;
        if (l < WireConstants.INT_LEN_1) return (0, l);
        return (l - WireConstants.INT_LEN_1 + 1, 0);
    }

    private static (int W, int Inline) FloatLen(int b)
    {
        int l = b & WireConstants.FLOAT_LEN_MASK;
        if (l < WireConstants.FLOAT_LEN_1) return (0, l);
        return (l - WireConstants.FLOAT_LEN_1 + 1, 0);
    }

    private static (int W, int Inline) StringLen(int b)
    {
        int l = b & WireConstants.STRING_LEN_MASK;
        if (l < WireConstants.STRING_LEN_1) return (0, l);
        if (l == WireConstants.STRING_LEN_1) return (1, 0);
        return (2, 0);
    }

    private static (int W, int Inline) BytesLen(int b)
    {
        int l = b & WireConstants.BYTES_LEN_MASK;
        if (l < WireConstants.BYTES_LEN_1) return (0, l);
        return (l - WireConstants.BYTES_LEN_1 + 1, 0);
    }

    private static (int W, int Inline) ContainerLen(int b)
    {
        int l = b & WireConstants.CONTAINER_LEN_MASK;
        if (l < WireConstants.CONTAINER_LEN_1) return (0, l);
        return (l - WireConstants.CONTAINER_LEN_1 + 1, 0);
    }

    private static (int W, int Inline) TagLen(int b)
    {
        int l = b & WireConstants.TAG_LEN_MASK;
        if (l < WireConstants.TAG_LEN_1) return (0, l);
        return (l - WireConstants.TAG_LEN_1 + 1, 0);
    }

    // TierOf 依 sign 的 3 位类型前缀分派，返回 (prefix, w, inline)。
    private static (int P, int W, int Inline) TierOf(int sign)
    {
        int p = sign & Prefix.PREFIX_MASK;
        switch (p)
        {
            case Prefix.SIMPLE:
                return (p, 0, 0);
            case Prefix.POSITIVE_INT:
            case Prefix.NEGATIVE_INT:
                { var (w, inline) = IntLen(sign); return (p, w, inline); }
            case Prefix.FLOAT:
                { var (w, inline) = FloatLen(sign); return (p, w, inline); }
            case Prefix.STRING:
                { var (w, inline) = StringLen(sign); return (p, w, inline); }
            case Prefix.BYTES:
                { var (w, inline) = BytesLen(sign); return (p, w, inline); }
            case Prefix.CONTAINER:
                { var (w, inline) = ContainerLen(sign); return (p, w, inline); }
            case Prefix.TAG:
                { var (w, inline) = TagLen(sign); return (p, w, inline); }
        }
        throw new MmDecodeException("invalid prefix " + p);
    }

    private static byte[] TierBytesBE(long length, int w)
    {
        if (w == 0) return Array.Empty<byte>();
        var outBytes = new byte[w];
        for (int i = w - 1; i >= 0; i--)
        {
            outBytes[w - 1 - i] = (byte)((length >> (8 * i)) & 0xFF);
        }
        return outBytes;
    }

    // ReadTierBE 从 [end-w, end) 读大端档位长度（用 long 防 0xFFFFFFFF 溢出）。
    private static long ReadTierBE(byte[] data, int end, int w)
    {
        if (end - w < 0)
        {
            throw new MmDecodeException("truncated: tier bytes out of range");
        }
        long v = 0;
        for (int i = end - w; i < end; i++)
        {
            v = v * 256 + data[i];
        }
        return v;
    }

    // NodeTailStart 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
    private static int NodeTailStart(byte[] wire, int end)
    {
        if (end < 1 || end > wire.Length)
        {
            throw new MmDecodeException("truncated: node end " + end + " out of range");
        }
        int sign = wire[end - 1];
        var (p, w, inline) = TierOf(sign);
        switch (p)
        {
            case Prefix.SIMPLE:
                return end - 1;
            case Prefix.POSITIVE_INT:
            case Prefix.NEGATIVE_INT:
                return w == 0 ? end - 1 : end - 1 - w;
            case Prefix.FLOAT:
                // float body 实宽 = w + 1（含 exp）→ wire [exp][mant][sign]
                return w == 0 ? end - 1 : end - 2 - w;
            case Prefix.STRING:
            case Prefix.BYTES:
            case Prefix.CONTAINER:
                if (w == 0)
                {
                    return end - 1 - inline;
                }
                return (int)(end - 1 - w - ReadTierBE(wire, end - 1, w));
            case Prefix.TAG:
                long total = w == 0 ? inline : ReadTierBE(wire, end - 1, w);
                int slotEnd = end - 1 - w;
                var (ml, slotW) = ParseFooterSlot(wire, slotEnd);
                // 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
                // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
                long inner = total - ml - slotW;
                if (inner < 0)
                {
                    throw new MmDecodeException("tag total underflow");
                }
                return (int)(slotEnd - slotW - ml - inner);
        }
        throw new MmDecodeException("invalid prefix");
    }

    // UntailifyNode：wire 形态 [start,end) → head 形态字节。
    private static byte[] UntailifyNode(byte[] wire, int start, int end)
    {
        if (start < 0 || end > wire.Length || start >= end)
        {
            throw new MmDecodeException("untailify: bad range [" + start + "," + end + ")");
        }
        int sign = wire[end - 1];
        var (p, w, inline) = TierOf(sign);

        if (p == Prefix.SIMPLE)
        {
            return new byte[] { (byte)sign };
        }
        if (p == Prefix.POSITIVE_INT || p == Prefix.NEGATIVE_INT || p == Prefix.FLOAT)
        {
            if (w == 0)
            {
                return new byte[] { (byte)sign };
            }
            byte[] body = Sub(wire, start, end - 1 - start);
            return Concat(new byte[] { (byte)sign }, body);
        }
        if (p == Prefix.STRING || p == Prefix.BYTES || p == Prefix.CONTAINER)
        {
            long l = w == 0 ? inline : ReadTierBE(wire, end - 1, w);
            byte[] body = Sub(wire, start, end - 1 - w - start);
            if (body.Length != l)
            {
                throw new MmDecodeException("untailify: payload " + body.Length + ", tier says " + l);
            }
            byte[] head = Concat(new byte[] { (byte)sign }, TierBytesBE(l, w));
            if (p == Prefix.CONTAINER)
            {
                return Concat(head, UntailifyChildren(body));
            }
            return Concat(head, body);
        }
        if (p == Prefix.TAG)
        {
            long total = w == 0 ? inline : ReadTierBE(wire, end - 1, w);
            int slotEnd = end - 1 - w;
            var (ml, slotW) = ParseFooterSlot(wire, slotEnd);
            int innerLen = (int)(total - ml - slotW);
            if (innerLen < 0)
            {
                throw new MmDecodeException("untailify tag: total underflow");
            }
            int metaEnd = slotEnd - slotW;
            byte[] meta = ReverseBytes(Sub(wire, metaEnd - ml, ml));
            byte[] innerWire = Sub(wire, start, metaEnd - ml - start);
            if (innerWire.Length != innerLen)
            {
                throw new MmDecodeException("untailify tag: inner " + innerWire.Length + ", want " + innerLen);
            }
            byte[] inner = innerLen > 0 ? UntailifyNode(innerWire, 0, innerWire.Length) : Array.Empty<byte>();
            byte[] outBytes = Concat(new byte[] { (byte)sign }, TierBytesBE(total, w));
            outBytes = Concat(outBytes, MetaPrefixBytes(ml));
            outBytes = Concat(outBytes, meta);
            return Concat(outBytes, inner);
        }
        throw new MmDecodeException("untailify: unreachable");
    }

    // UntailifyChildren：wire 容器 payload 内顺序 wire 子节点逐个还原（从末端反向扫描，正序拼接）。
    private static byte[] UntailifyChildren(byte[] payload)
    {
        var nodes = new List<byte[]>();
        int end = payload.Length;
        while (end > 0)
        {
            int start = NodeTailStart(payload, end);
            if (start < 0 || start >= end)
            {
                throw new MmDecodeException("untailify children: bad boundary " + start + ".." + end);
            }
            nodes.Add(UntailifyNode(payload, start, end));
            end = start;
        }
        byte[] outBytes = Array.Empty<byte>();
        for (int i = nodes.Count - 1; i >= 0; i--)
        {
            outBytes = Concat(outBytes, nodes[i]);
        }
        return outBytes;
    }

    // Untailify 解码入口：整帧 wire 形态 → head 形态。
    public static byte[] Untailify(byte[] wire)
    {
        if (wire.Length == 0)
        {
            throw new MmDecodeException("empty wire");
        }
        int start = NodeTailStart(wire, wire.Length);
        if (start != 0)
        {
            throw new MmDecodeException("footer boundary mismatch: node start " + start + ", want 0");
        }
        return UntailifyNode(wire, start, wire.Length);
    }

    // ----- 字节数组辅助 -----

    private static byte[] Sub(byte[] arr, int start, int len)
    {
        if (len <= 0) return Array.Empty<byte>();
        if (start < 0 || start + len > arr.Length)
        {
            throw new MmDecodeException("untailify: slice out of range");
        }
        var outBytes = new byte[len];
        Array.Copy(arr, start, outBytes, 0, len);
        return outBytes;
    }

    private static byte[] Concat(byte[] a, byte[] b)
    {
        if (a.Length == 0) return b;
        if (b.Length == 0) return a;
        var outBytes = new byte[a.Length + b.Length];
        Array.Copy(a, 0, outBytes, 0, a.Length);
        Array.Copy(b, 0, outBytes, a.Length, b.Length);
        return outBytes;
    }
}
