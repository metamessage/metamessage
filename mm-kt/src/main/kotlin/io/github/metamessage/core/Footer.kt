package io.github.metamessage.core

/**
 * v0.5.0 尾定界辅助（与 Go internal/core/footer.go + tailframe.go 逐字节对齐）：
 *
 *   1. meta_len 尾部槽编解码（footerSlotBytes / parseFooterSlot）+ 字节逆序（reverseBytes）。
 *   2. 解码入口剥帧 untailify（nodeTailStart / untailifyNode / untailifyChildren）：
 *      从节点尾 sign 反向定界，把尾定界 wire 还原成「sign 在前」的头部形态，再交给原有正向解码器。
 *
 * 编码器各原语已直接按尾定界顺序写出（body→lenBE→sign），本文件不含编码侧换位逻辑，
 * 也不设独立 payload_len 槽：各类型 sign 低 5 位的档位/内联语义（v0.3.0 原样）同时承担定界信息。
 */
object Footer {
    const val SLOT_DIRECT_MAX = 253 // 一档可直写的 meta 长度上限（尾字节 < 254）
    const val SLOT_MARKER_1 = 254 // 二档标记：倒数第 2 字节为 u8 长度
    const val SLOT_MARKER_2 = 255 // 三档标记：倒数第 3..2 字节为 u16 大端长度

    /** footerSlotBytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。 */
    fun footerSlotBytes(metaLen: Int): ByteArray =
            when {
                metaLen <= SLOT_DIRECT_MAX -> byteArrayOf(metaLen.toByte())
                metaLen <= 0xFF -> byteArrayOf(metaLen.toByte(), SLOT_MARKER_1.toByte())
                metaLen <= 0xFFFF ->
                        byteArrayOf(
                                ((metaLen shr 8) and 0xFF).toByte(),
                                (metaLen and 0xFF).toByte(),
                                SLOT_MARKER_2.toByte()
                        )
                else -> throw MmDecodeException("tag too long, max length: 65535, actual: $metaLen")
            }

    /** parseFooterSlot 从节点尾端反向解析长度槽，返回 [metaLen, slotWidth]。 */
    fun parseFooterSlot(data: ByteArray, nodeEnd: Int): IntArray {
        if (nodeEnd < 1) throw MmDecodeException("footer slot out of range")
        return when (val marker = data[nodeEnd - 1].toInt() and 0xFF) {
            SLOT_MARKER_1 -> {
                if (nodeEnd < 2) throw MmDecodeException("footer slot truncated")
                intArrayOf(data[nodeEnd - 2].toInt() and 0xFF, 2)
            }
            SLOT_MARKER_2 -> {
                if (nodeEnd < 3) throw MmDecodeException("footer slot truncated")
                intArrayOf(((data[nodeEnd - 3].toInt() and 0xFF) shl 8) or (data[nodeEnd - 2].toInt() and 0xFF), 3)
            }
            else -> intArrayOf(marker, 1)
        }
    }

    /** reverseBytes 返回 bs 的字节逆序拷贝。 */
    fun reverseBytes(bs: ByteArray): ByteArray {
        val out = ByteArray(bs.size)
        for (i in bs.indices) out[bs.size - 1 - i] = bs[i]
        return out
    }

    /** metaPrefixWidth 正向 meta 前缀槽宽度：与尾部槽同一套 canonical 三档，宽度恒等于 footerSlotBytes 槽宽。 */
    fun metaPrefixWidth(ml: Int): Int = when {
        ml < 254 -> 1
        ml <= 255 -> 2
        else -> 3
    }

    /** metaPrefixBytes 还原正向 meta 前缀槽。 */
    fun metaPrefixBytes(ml: Int): ByteArray =
            when {
                ml < 254 -> byteArrayOf(ml.toByte())
                ml <= 0xFF -> byteArrayOf(254.toByte(), ml.toByte())
                ml <= 0xFFFF -> byteArrayOf(255.toByte(), ((ml shr 8) and 0xFF).toByte(), (ml and 0xFF).toByte())
                else -> throw MmDecodeException("meta too long: $ml")
            }

    // ----- 各类型档位：返回 [w, inline]，w==0 表示内联档 -----

    private fun intLen(b: Int): IntArray {
        val l = b and WireConstants.INT_LEN_MASK
        return if (l < WireConstants.INT_LEN_1) intArrayOf(0, l) else intArrayOf(l - WireConstants.INT_LEN_1 + 1, 0)
    }

    private fun floatLen(b: Int): IntArray {
        val l = b and WireConstants.FLOAT_LEN_MASK
        return if (l < WireConstants.FLOAT_LEN_1) intArrayOf(0, l) else intArrayOf(l - WireConstants.FLOAT_LEN_1 + 1, 0)
    }

    private fun stringLen(b: Int): IntArray {
        val l = b and WireConstants.STRING_LEN_MASK
        return when {
            l < WireConstants.STRING_LEN_1 -> intArrayOf(0, l)
            l == WireConstants.STRING_LEN_1 -> intArrayOf(1, 0)
            else -> intArrayOf(2, 0)
        }
    }

    private fun bytesLen(b: Int): IntArray {
        val l = b and WireConstants.BYTES_LEN_MASK
        return if (l < WireConstants.BYTES_LEN_1) intArrayOf(0, l) else intArrayOf(l - WireConstants.BYTES_LEN_1 + 1, 0)
    }

    private fun containerLen(b: Int): IntArray {
        val l = b and WireConstants.CONTAINER_LEN_MASK
        return if (l < WireConstants.CONTAINER_LEN_1) intArrayOf(0, l) else intArrayOf(l - WireConstants.CONTAINER_LEN_1 + 1, 0)
    }

    private fun tagLen(b: Int): IntArray {
        val l = b and WireConstants.TAG_LEN_MASK
        return if (l < WireConstants.TAG_LEN_1) intArrayOf(0, l) else intArrayOf(l - WireConstants.TAG_LEN_1 + 1, 0)
    }

    /** tierOf 依 sign 的 3 位类型前缀分派，返回 [prefix, w, inline]。 */
    private fun tierOf(sign: Int): IntArray {
        val p = sign and Prefix.PREFIX_MASK
        return when (p) {
            Prefix.SIMPLE -> intArrayOf(p, 0, 0)
            Prefix.POSITIVE_INT, Prefix.NEGATIVE_INT -> intArrayOf(p) + intLen(sign)
            Prefix.FLOAT -> intArrayOf(p) + floatLen(sign)
            Prefix.STRING -> intArrayOf(p) + stringLen(sign)
            Prefix.BYTES -> intArrayOf(p) + bytesLen(sign)
            Prefix.CONTAINER -> intArrayOf(p) + containerLen(sign)
            Prefix.TAG -> intArrayOf(p) + tagLen(sign)
            else -> throw MmDecodeException("invalid prefix $p")
        }
    }

    private fun tierBytesBE(length: Long, w: Int): ByteArray {
        if (w == 0) return ByteArray(0)
        val out = ByteArray(w)
        for (i in w - 1 downTo 0) out[w - 1 - i] = ((length shr (8 * i)) and 0xFF).toByte()
        return out
    }

    // readTierBE 从 [end-w, end) 读大端档位长度（用 Long 防 0xFFFFFFFF 溢出）。
    private fun readTierBE(data: ByteArray, end: Int, w: Int): Long {
        if (end - w < 0) throw MmDecodeException("truncated: tier bytes out of range")
        var v = 0L
        for (i in end - w until end) v = v * 256 + (data[i].toInt() and 0xFF)
        return v
    }

    // nodeTailStart 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
    private fun nodeTailStart(wire: ByteArray, end: Int): Int {
        if (end < 1 || end > wire.size) throw MmDecodeException("truncated: node end $end out of range")
        val sign = wire[end - 1].toInt() and 0xFF
        val t = tierOf(sign)
        val p = t[0]
        val w = t[1]
        val inline = t[2]
        return when (p) {
            Prefix.SIMPLE -> end - 1
            Prefix.POSITIVE_INT, Prefix.NEGATIVE_INT -> if (w == 0) end - 1 else end - 1 - w
            // float body 实宽 = w + 1（含 exp）→ wire [exp][mant][sign]
            Prefix.FLOAT -> if (w == 0) end - 1 else end - 2 - w
            Prefix.STRING, Prefix.BYTES, Prefix.CONTAINER ->
                    if (w == 0) end - 1 - inline else (end - 1 - w - readTierBE(wire, end - 1, w)).toInt()
            Prefix.TAG -> {
                val total = if (w == 0) inline.toLong() else readTierBE(wire, end - 1, w)
                val slotEnd = end - 1 - w
                val slot = parseFooterSlot(wire, slotEnd)
                val ml = slot[0]
                val slotW = slot[1]
                // 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
                // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
                val inner = total - ml - slotW
                if (inner < 0) throw MmDecodeException("tag total underflow")
                (slotEnd - slotW - ml - inner).toInt()
            }
            else -> throw MmDecodeException("invalid prefix")
        }
    }

    // untailifyNode：wire 形态 [start,end) → head 形态字节。
    private fun untailifyNode(wire: ByteArray, start: Int, end: Int): ByteArray {
        if (start < 0 || end > wire.size || start >= end) {
            throw MmDecodeException("untailify: bad range [$start,$end)")
        }
        val sign = wire[end - 1].toInt() and 0xFF
        val t = tierOf(sign)
        val p = t[0]
        val w = t[1]
        val inline = t[2]

        if (p == Prefix.SIMPLE) return byteArrayOf(sign.toByte())
        if (p == Prefix.POSITIVE_INT || p == Prefix.NEGATIVE_INT || p == Prefix.FLOAT) {
            if (w == 0) return byteArrayOf(sign.toByte())
            val body = wire.copyOfRange(start, end - 1)
            return byteArrayOf(sign.toByte()) + body
        }
        if (p == Prefix.STRING || p == Prefix.BYTES || p == Prefix.CONTAINER) {
            val l = if (w == 0) inline.toLong() else readTierBE(wire, end - 1, w)
            val body = wire.copyOfRange(start, end - 1 - w)
            if (body.size.toLong() != l) {
                throw MmDecodeException("untailify: payload ${body.size}, tier says $l")
            }
            val head = byteArrayOf(sign.toByte()) + tierBytesBE(l, w)
            return if (p == Prefix.CONTAINER) head + untailifyChildren(body) else head + body
        }
        if (p == Prefix.TAG) {
            val total = if (w == 0) inline.toLong() else readTierBE(wire, end - 1, w)
            val slotEnd = end - 1 - w
            val slot = parseFooterSlot(wire, slotEnd)
            val ml = slot[0]
            val slotW = slot[1]
            val innerLen = (total - ml - slotW).toInt()
            if (innerLen < 0) throw MmDecodeException("untailify tag: total underflow")
            val metaEnd = slotEnd - slotW
            val meta = reverseBytes(wire.copyOfRange(metaEnd - ml, metaEnd))
            val innerWire = wire.copyOfRange(start, metaEnd - ml)
            if (innerWire.size != innerLen) {
                throw MmDecodeException("untailify tag: inner ${innerWire.size}, want $innerLen")
            }
            val inner = if (innerLen > 0) untailifyNode(innerWire, 0, innerWire.size) else ByteArray(0)
            var out = byteArrayOf(sign.toByte()) + tierBytesBE(total, w)
            out += metaPrefixBytes(ml)
            out += meta
            out += inner
            return out
        }
        throw MmDecodeException("untailify: unreachable")
    }

    // untailifyChildren：wire 容器 payload 内顺序子节点逐个还原（从末端反向扫描，正序拼接）。
    private fun untailifyChildren(payload: ByteArray): ByteArray {
        val nodes = ArrayList<ByteArray>()
        var end = payload.size
        while (end > 0) {
            val start = nodeTailStart(payload, end)
            if (start < 0 || start >= end) {
                throw MmDecodeException("untailify children: bad boundary $start..$end")
            }
            nodes.add(untailifyNode(payload, start, end))
            end = start
        }
        var out = ByteArray(0)
        for (i in nodes.indices.reversed()) out += nodes[i]
        return out
    }

    // untailify 解码入口：整帧 wire 形态 → head 形态。
    fun untailify(wire: ByteArray): ByteArray {
        if (wire.isEmpty()) throw MmDecodeException("empty wire")
        val start = nodeTailStart(wire, wire.size)
        if (start != 0) {
            throw MmDecodeException("footer boundary mismatch: node start $start, want 0")
        }
        return untailifyNode(wire, start, wire.size)
    }
}
