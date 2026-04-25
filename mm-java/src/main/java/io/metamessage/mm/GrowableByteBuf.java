package io.metamessage.mm;

import java.util.Arrays;

final class GrowableByteBuf {
    private static final int MAX_CAP = 1024 * 1024 * 1024;

    private byte[] buf = new byte[1024];
    private int len;

    void write(byte... bs) {
        ensure(bs.length);
        System.arraycopy(bs, 0, buf, len, bs.length);
        len += bs.length;
    }

    void write(byte[] bs, int off, int length) {
        ensure(length);
        System.arraycopy(bs, off, buf, len, length);
        len += length;
    }

    void writeAll(byte[] bs) {
        if (bs.length == 0) {
            return;
        }
        write(bs, 0, bs.length);
    }

    int size() {
        return len;
    }

    byte[] copyRange(int start, int endExclusive) {
        return Arrays.copyOfRange(buf, start, endExclusive);
    }

    void reset() {
        len = 0;
    }

    byte[] entireArrayView() {
        return buf;
    }

    int length() {
        return len;
    }

    private void ensure(int n) {
        if (len + n > MAX_CAP) {
            throw new IllegalStateException("maximum size exceeded");
        }
        if (len + n > buf.length) {
            int newCap = buf.length * 2;
            if (newCap < len + n) {
                newCap = len + n;
            }
            buf = Arrays.copyOf(buf, newCap);
        }
    }
}
