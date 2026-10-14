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
        Assert.Equal(WantWire, wire.Length);
        Assert.Equal(0xFF, wire[1]);
        Assert.Equal(0xFF, wire[2]);
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
        var enc = new WireEncoder();
        Assert.Throws<Exception>(() => enc.EncodeString(AsciiString(Max2 + 1)));
        Assert.Throws<Exception>(() => enc.EncodeBytes(AsciiBytes(Max2 + 1)));
        Assert.Throws<Exception>(() => enc.EncodeArrayPayload(AsciiBytes(Max2 + 1)));
        Assert.Throws<Exception>(() => enc.EncodeObjectPayload(AsciiBytes(Max2 + 1)));
    }
}
