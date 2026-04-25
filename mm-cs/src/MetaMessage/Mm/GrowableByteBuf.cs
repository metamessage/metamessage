namespace MetaMessage.Mm;

public class GrowableByteBuf
{
    private byte[] _buffer;
    private int _length;
    private const int DefaultCapacity = 1024;

    public GrowableByteBuf(int capacity = DefaultCapacity)
    {
        _buffer = new byte[capacity];
        _length = 0;
    }

    public int Length => _length;

    public void Write(int value)
    {
        EnsureCapacity(_length + 1);
        _buffer[_length++] = (byte)value;
    }

    public void Write(int value, byte b)
    {
        EnsureCapacity(_length + 2);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b;
    }

    public void Write(int value, byte b1, byte b2)
    {
        EnsureCapacity(_length + 3);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
    }

    public void Write(int value, byte b1, byte b2, byte b3)
    {
        EnsureCapacity(_length + 4);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
    }

    public void Write(int value, byte b1, byte b2, byte b3, byte b4)
    {
        EnsureCapacity(_length + 5);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
        _buffer[_length++] = b4;
    }

    public void Write(int value, byte b1, byte b2, byte b3, byte b4, byte b5)
    {
        EnsureCapacity(_length + 6);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
        _buffer[_length++] = b4;
        _buffer[_length++] = b5;
    }

    public void Write(int value, byte b1, byte b2, byte b3, byte b4, byte b5, byte b6)
    {
        EnsureCapacity(_length + 7);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
        _buffer[_length++] = b4;
        _buffer[_length++] = b5;
        _buffer[_length++] = b6;
    }

    public void Write(int value, byte b1, byte b2, byte b3, byte b4, byte b5, byte b6, byte b7)
    {
        EnsureCapacity(_length + 8);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
        _buffer[_length++] = b4;
        _buffer[_length++] = b5;
        _buffer[_length++] = b6;
        _buffer[_length++] = b7;
    }

    public void Write(int value, byte b1, byte b2, byte b3, byte b4, byte b5, byte b6, byte b7, byte b8)
    {
        EnsureCapacity(_length + 9);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
        _buffer[_length++] = b4;
        _buffer[_length++] = b5;
        _buffer[_length++] = b6;
        _buffer[_length++] = b7;
        _buffer[_length++] = b8;
    }

    public void Write(int value, byte b1, byte b2, byte b3, byte b4, byte b5, byte b6, byte b7, byte b8, byte b9)
    {
        EnsureCapacity(_length + 10);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        _buffer[_length++] = b3;
        _buffer[_length++] = b4;
        _buffer[_length++] = b5;
        _buffer[_length++] = b6;
        _buffer[_length++] = b7;
        _buffer[_length++] = b8;
        _buffer[_length++] = b9;
    }

    public void WriteWithBytes(int value, byte[] bytes)
    {
        EnsureCapacity(_length + 1 + bytes.Length);
        _buffer[_length++] = (byte)value;
        if (bytes.Length > 0)
        {
            Array.Copy(bytes, 0, _buffer, _length, bytes.Length);
            _length += bytes.Length;
        }
    }

    public void WriteWithBytes(int value, byte b, byte[] bytes)
    {
        EnsureCapacity(_length + 2 + bytes.Length);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b;
        if (bytes.Length > 0)
        {
            Array.Copy(bytes, 0, _buffer, _length, bytes.Length);
            _length += bytes.Length;
        }
    }

    public void WriteWithBytes(int value, byte b1, byte b2, byte[] bytes)
    {
        EnsureCapacity(_length + 3 + bytes.Length);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        if (bytes.Length > 0)
        {
            Array.Copy(bytes, 0, _buffer, _length, bytes.Length);
            _length += bytes.Length;
        }
    }

    public void WriteWithMultipleBytes(int value, params byte[][] byteArrays)
    {
        int totalLength = 1;
        foreach (var arr in byteArrays)
        {
            totalLength += arr.Length;
        }
        EnsureCapacity(_length + totalLength);
        _buffer[_length++] = (byte)value;
        foreach (var arr in byteArrays)
        {
            if (arr.Length > 0)
            {
                Array.Copy(arr, 0, _buffer, _length, arr.Length);
                _length += arr.Length;
            }
        }
    }

    public void WriteWithMultipleBytes(int value, byte b, params byte[][] byteArrays)
    {
        int totalLength = 2;
        foreach (var arr in byteArrays)
        {
            totalLength += arr.Length;
        }
        EnsureCapacity(_length + totalLength);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b;
        foreach (var arr in byteArrays)
        {
            if (arr.Length > 0)
            {
                Array.Copy(arr, 0, _buffer, _length, arr.Length);
                _length += arr.Length;
            }
        }
    }

    public void WriteWithMultipleBytes(int value, byte b1, byte b2, params byte[][] byteArrays)
    {
        int totalLength = 3;
        foreach (var arr in byteArrays)
        {
            totalLength += arr.Length;
        }
        EnsureCapacity(_length + totalLength);
        _buffer[_length++] = (byte)value;
        _buffer[_length++] = b1;
        _buffer[_length++] = b2;
        foreach (var arr in byteArrays)
        {
            if (arr.Length > 0)
            {
                Array.Copy(arr, 0, _buffer, _length, arr.Length);
                _length += arr.Length;
            }
        }
    }

    public void WriteAll(byte[] bytes)
    {
        if (bytes == null || bytes.Length == 0)
            return;

        EnsureCapacity(_length + bytes.Length);
        Array.Copy(bytes, 0, _buffer, _length, bytes.Length);
        _length += bytes.Length;
    }

    public byte[] CopyRange(int start, int end)
    {
        if (start < 0 || end > _length || start >= end)
            return Array.Empty<byte>();

        int length = end - start;
        byte[] result = new byte[length];
        Array.Copy(_buffer, start, result, 0, length);
        return result;
    }

    public void Reset()
    {
        _length = 0;
    }

    public byte[] ToArray()
    {
        byte[] result = new byte[_length];
        Array.Copy(_buffer, 0, result, 0, _length);
        return result;
    }

    private void EnsureCapacity(int required)
    {
        if (required <= _buffer.Length)
            return;

        int newCapacity = Math.Max(_buffer.Length * 2, required);
        byte[] newBuffer = new byte[newCapacity];
        Array.Copy(_buffer, 0, newBuffer, 0, _length);
        _buffer = newBuffer;
    }
}