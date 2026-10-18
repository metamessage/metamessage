using Xunit;
using MetaMessage.Core;

namespace MetaMessageTests;

/// <summary>
/// v0.5.0 尾定界模式专项测试：golden hex / 尾部槽档位 / 篡改防御。
/// 期望字节与 Go internal/core/footer_test.go、mm-php FooterTest、mm-py footer 一致（逐字节对齐 Go）。
/// </summary>
public class FooterTest
{
    private static string Hex(byte[] b) =>
        string.Join(' ', Array.ConvertAll(b, x => x.ToString("X2")));

    private static byte[] FromHex(string s)
    {
        var parts = s.Split(' ', StringSplitOptions.RemoveEmptyEntries);
        var b = new byte[parts.Length];
        for (int i = 0; i < parts.Length; i++) b[i] = Convert.ToByte(parts[i], 16);
        return b;
    }

    [Fact]
    public void SlotTiersCanonicalWidths()
    {
        var cases = new (int MetaLen, int Width)[] { (1, 1), (253, 1), (254, 2), (255, 2), (256, 3), (65535, 3) };
        foreach (var (metaLen, width) in cases)
        {
            byte[] slot = Footer.FooterSlotBytes(metaLen);
            Assert.Equal(width, slot.Length);
            var data = new List<byte> { 0x00, 0x11, 0x22, 0x33 }.Concat(slot).ToArray();
            var (gotLen, gotW) = Footer.ParseFooterSlot(data, data.Length);
            Assert.Equal(metaLen, gotLen);
            Assert.Equal(width, gotW);
        }
    }

    [Fact]
    public void SlotTooLongRejected()
    {
        Assert.Throws<MmDecodeException>(() => Footer.FooterSlotBytes(65536));
    }

    [Fact]
    public void GoldenSingleFieldDesc()
    {
        // 顶层对象 { "v": "hi" } 带 desc=ab：
        //   key 'v' 内联 [76 81]；value 尾序标量 tag 段 [68 69 82 62 61 1A 03 E7]；object sign [CB]。
        var wire = MetaMessage.Core.MetaMessage.FromJSONC("{\n  // mm: desc=ab\n  \"v\": \"hi\"\n}");
        Assert.Equal("76 81 D2 68 69 82 62 61 1A 03 E7 CB", Hex(wire));
    }

    [Fact]
    public void GoldenMultiFieldOrder()
    {
        var wire = MetaMessage.Core.MetaMessage.FromJSONC("{\n  // mm: example;nullable;desc=ab\n  \"v\": \"x\"\n}");
        Assert.Equal("76 81 D2 78 81 62 61 1A 29 09 05 E8 0C CC", Hex(wire));
        var tag = MetaMessage.Core.MetaMessage.Decode(wire).Tag;
        // 顶层对象 tag 无 example；value 子节点承载该 tag —— 校验解码往返可还原。
        Assert.NotNull(tag);
    }

    [Fact]
    public void ReverseBytesRoundtrip()
    {
        byte[] bs = { 1, 2, 3, 4, 5 };
        Assert.Equal(bs, Footer.ReverseBytes(Footer.ReverseBytes(bs)));
        Assert.Equal(new byte[] { 5, 4, 3, 2, 1 }, Footer.ReverseBytes(bs));
    }

    [Fact]
    public void RoundtripStable()
    {
        string[] srcs =
        {
            "{\n  // mm: desc=数据\n  \"v\": \"1\"\n}",
            "{\n  \"a\": [1, 2]\n}",
        };
        foreach (var src in srcs)
        {
            var wire = MetaMessage.Core.MetaMessage.FromJSONC(src);
            var jsonc = MetaMessage.Core.MetaMessage.DecodeToJsonc(wire);
            var wire2 = MetaMessage.Core.MetaMessage.FromJSONC(jsonc);
            Assert.Equal(Hex(wire), Hex(wire2));
        }
    }

    [Fact]
    public void DecodeDefenseTamperedFooter()
    {
        byte[] golden = FromHex("76 81 D2 68 69 82 62 61 1A 03 E7 CB");

        // 截断尾部 → 定界不再对齐
        Assert.Throws<MmDecodeException>(() => MetaMessage.Core.MetaMessage.Decode(golden[..^1]));

        // meta_len 尾槽篡改为巨大值 → inner 越界
        var a = (byte[])golden.Clone();
        a[10] = 0x7F;
        Assert.Throws<MmDecodeException>(() => MetaMessage.Core.MetaMessage.Decode(a));

        // tag sign 篡改为 bogus 档位 → tag total underflow
        var c = (byte[])golden.Clone();
        c[11] = 0xEB;
        Assert.Throws<MmDecodeException>(() => MetaMessage.Core.MetaMessage.Decode(c));
    }
}
