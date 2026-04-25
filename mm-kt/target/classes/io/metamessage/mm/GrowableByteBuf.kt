package io.metamessage.mm

import java.util.Arrays

class GrowableByteBuf {
    private var buf = ByteArray(1024)
    private var len = 0

    fun write(vararg bs: Byte) {
        ensure(bs.size)
        System.arraycopy(bs, 0, buf, len, bs.size)
        len += bs.size
    }

    fun write(bs: ByteArray, off: Int, length: Int) {
        ensure(length)
        System.arraycopy(bs, off, buf, len, length)
        len += length
    }

    fun writeAll(bs: ByteArray) {
        if (bs.isEmpty()) return
        write(bs, 0, bs.size)
    }

    fun size(): Int = len

    fun copyRange(start: Int, endExclusive: Int): ByteArray {
        return Arrays.copyOfRange(buf, start, endExclusive)
    }

    fun reset() {
        len = 0
    }

    fun entireArrayView(): ByteArray = buf

    fun length(): Int = len

    private fun ensure(n: Int) {
        if (len + n > MAX_CAP) {
            throw IllegalStateException("maximum size exceeded")
        }
        if (len + n > buf.size) {
            var newCap = buf.size * 2
            if (newCap < len + n) {
                newCap = len + n
            }
            buf = Arrays.copyOf(buf, newCap)
        }
    }

    companion object {
        private const val MAX_CAP = 1024 * 1024 * 1024
    }
}
