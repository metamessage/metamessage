using Xunit;
using MetaMessage.Core;

namespace MetaMessageTests;

/// <summary>
/// Boundary tests for the 2-byte length bucket (MAX_2 = 65535).
///
/// MAX_2 itself must be accepted: it is representable by the 2-byte length
/// field (0xFF 0xFF). A `length &lt; MAX_2` guard would fall through to the
/// "too long" error for exactly that value, while the wire format and every
/// other language port accept it.
/// </summary>
public class WireEncoderBoundaryTest
{
    private const int Max2 = 0xFFFF; // 65535
    private const int WantWire = Max2 + 3; // 1 prefix + 2 length + payload

    private static byte[] AsciiBytes(int n) => Enumerable.Repeat((byte)'a', n).ToArray();

    private static string AsciiString(int n) => new('a', n);

    private static void AssertLenPrefix(byte[] wire)
    {
        // v0.5.0 尾定界：[payload][len hi][len lo][sign]，2 字节长度档落在尾部。
        Assert.Equal(WantWire, wire.Length);
        Assert.Equal(0xFF, wire[wire.Length - 3]);
        Assert.Equal(0xFF, wire[wire.Length - 2]);
    }

    [Fact]
    public void EncodeStringAtMax2Byte()
    {
        var enc = new WireEncoder();
        int n = enc.EncodeString(AsciiString(Max2));
        Assert.Equal(WantWire, n);
        AssertLenPrefix(enc.ToByteArray());
    }

    [Fact]
    public void EncodeBytesAtMax2Byte()
    {
        var enc = new WireEncoder();
        int n = enc.EncodeBytes(AsciiBytes(Max2));
        Assert.Equal(WantWire, n);
        AssertLenPrefix(enc.ToByteArray());
    }

    [Fact]
    public void EncodeArrayPayloadAtMax2Byte()
    {
        var enc = new WireEncoder();
        int n = enc.EncodeArrayPayload(AsciiBytes(Max2));
        Assert.Equal(WantWire, n);
        AssertLenPrefix(enc.ToByteArray());
    }

    [Fact]
    public void EncodeObjectPayloadAtMax2Byte()
    {
        var enc = new WireEncoder();
        int n = enc.EncodeObjectPayload(AsciiBytes(Max2));
        Assert.Equal(WantWire, n);
        AssertLenPrefix(enc.ToByteArray());
    }

    [Fact]
    public void EncodeAtMax2BytePlusOneStillRejected()
    {
        // String is capped at MAX_2 by the wire format (STRING_LEN_2 is the top tier),
        // so MAX_2+1 must reject.
        var enc = new WireEncoder();
        Assert.Throws<Exception>(() => enc.EncodeString(AsciiString(Max2 + 1)));
    }

    [Fact]
    public void EncodeBytesArrayObjectAtMax2BytePlusOneUsesNextTier()
    {
        // Bytes / ArrayPayload / ObjectPayload have 4 length tiers (v0.3.0 alignment
        // with the Go reference), so MAX_2+1 legitimately falls into the 3-byte tier
        // instead of being rejected.
        var enc = new WireEncoder();
        int payloadLen = Max2 + 1;
        Assert.Null(Record.Exception(() => enc.EncodeBytes(AsciiBytes(payloadLen))));
        Assert.Null(Record.Exception(() => enc.EncodeArrayPayload(AsciiBytes(payloadLen))));
        Assert.Null(Record.Exception(() => enc.EncodeObjectPayload(AsciiBytes(payloadLen))));
    }
}
