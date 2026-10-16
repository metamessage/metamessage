package io.github.metamessage.core

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
        val required = len.toLong() + n.toLong()
        if (required > MAX_CAP) {
            throw IllegalStateException("exceeds format limit (max $MAX_CAP bytes)")
        }
        if (required > RUNTIME_MAX_ARRAY) {
            throw IllegalStateException("exceeds runtime allocatable limit (ByteArray < 2GB)")
        }
        if (required > buf.size) {
            var newCap = buf.size.toLong() * 2
            if (newCap < required) {
                newCap = required.toLong()
            }
            if (newCap > RUNTIME_MAX_ARRAY) {
                newCap = RUNTIME_MAX_ARRAY
            }
            buf = Arrays.copyOf(buf, newCap.toInt())
        }
    }

    companion object {
        // v0.3.0：格式上限对齐 Max4Byte(0xFFFFFFFF) ≈ 4GB
        private const val MAX_CAP = 0xFFFFFFFFL
        // JVM 实际可分配数组上限（Int 索引，留 8 字节头余量）
        private const val RUNTIME_MAX_ARRAY = 2147483647L - 8
    }
}
