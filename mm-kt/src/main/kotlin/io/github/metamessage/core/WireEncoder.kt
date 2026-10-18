package io.github.metamessage.core

import java.nio.charset.StandardCharsets

class WireEncoder {

    private val buf = GrowableByteBuf()

    fun toByteArray(): ByteArray = buf.copyRange(0, buf.length())

    fun finishTakeLast(writtenFromEnd: Int): ByteArray {
        val end = buf.length()
        return buf.copyRange(end - writtenFromEnd, end)
    }

    fun reset() {
        buf.reset()
    }

    fun size(): Int = buf.length()

    fun encodeSimple(simpleValue: Int): Int {
        require(simpleValue in 0..31) { "SimpleValue out of range: $simpleValue" }
        return writeByte((Prefix.SIMPLE or simpleValue).toByte())
    }

    fun encodeBool(v: Boolean): Int = encodeSimple(if (v) SimpleValue.TRUE else SimpleValue.FALSE)

    fun encodeInt64(v: Long): Int {
        if (v >= 0) {
            return encodeUintWithPrefix(Prefix.POSITIVE_INT, v)
        }
        val uv = if (v == Long.MIN_VALUE) 1L shl 63 else -v
        return encodeUintWithPrefix(Prefix.NEGATIVE_INT, uv)
    }

    fun encodeU64(uv: Long): Int {
        require(uv >= 0) { "expected unsigned" }
        return encodeUintWithPrefix(Prefix.POSITIVE_INT, uv)
    }

    // v0.5.0 尾定界：一趟正向直接按 [body BE][sign] 写出，档位/内联语义与 v0.3.x 一致，
    // 仅物理落盘顺序改为尾部形态（sign 落尾，长度天然已知，零回写）。
    private fun encodeUintWithPrefix(prefix: Int, uv: Long): Int {
        val start = buf.length()
        if (java.lang.Long.compareUnsigned(uv, WireConstants.INT_LEN_1.toLong()) < 0) {
            buf.write((prefix or uv.toInt()).toByte())
            return buf.length() - start
        }
        val width = when {
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_1) <= 0 -> 1
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_2) <= 0 -> 2
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_3) <= 0 -> 3
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_4) <= 0 -> 4
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_5) <= 0 -> 5
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_6) <= 0 -> 6
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_7) <= 0 -> 7
            else -> 8
        }
        for (i in width - 1 downTo 0) buf.write((uv shr (8 * i)).toByte())
        buf.write((prefix or (WireConstants.INT_LEN_1 - 1 + width)).toByte())
        return buf.length() - start
    }

    fun encodeFloatString(s: String): Int {
        val p = FloatCodec.parseDecimalString(s)
        val start = buf.length()
        var sign = Prefix.FLOAT
        if (p.negative) sign = sign or WireConstants.FLOAT_NEG_MASK
        if (p.exponent.toInt() == -1 && p.mantissa <= 7) {
            writeByte((sign or p.mantissa.toInt()).toByte())
        } else {
            val mantissa = p.mantissa
            // 尾定界：body [exp][mantissa BE w字节]，sign(档位) 落尾
            val width = when {
                mantissa <= WireConstants.MAX_1 -> 1
                mantissa <= WireConstants.MAX_2 -> 2
                mantissa <= WireConstants.MAX_3 -> 3
                mantissa <= WireConstants.MAX_4 -> 4
                mantissa <= WireConstants.MAX_5 -> 5
                mantissa <= WireConstants.MAX_6 -> 6
                mantissa <= WireConstants.MAX_7 -> 7
                else -> 8
            }
            writeByte(p.exponent)
            for (i in width - 1 downTo 0) writeByte((mantissa shr (8 * i)).toByte())
            writeByte((sign or (WireConstants.FLOAT_LEN_1 - 1 + width)).toByte())
        }
        return buf.length() - start
    }

    fun encodeString(s: String): Int {
        val utf = s.toByteArray(StandardCharsets.UTF_8)
        val length = utf.size
        require(length <= WireConstants.MAX_2) { "string too long" }
        val sign = Prefix.STRING
        // 尾定界：[body][lenBE 档位][sign]
        return when {
            length < WireConstants.STRING_LEN_1 ->
                    writeBytesThenSuffix(utf, (sign or length).toByte())
            length < WireConstants.MAX_1 ->
                    writeBytesThenSuffix(utf, length.toByte(), (sign or WireConstants.STRING_LEN_1).toByte())
            else ->
                    writeBytesThenSuffix(
                            utf,
                            (length shr 8).toByte(),
                            length.toByte(),
                            (sign or WireConstants.STRING_LEN_2).toByte()
                    )
        }
    }

    fun encodeBytes(data: ByteArray): Int {
        val length = data.size
        require(length <= WireConstants.MAX_4) { "bytes too long" }
        val sign = Prefix.BYTES
        // 尾定界：[body][lenBE 档位][sign]
        return when {
            length < WireConstants.BYTES_LEN_1 ->
                    writeBytesThenSuffix(data, (sign or length).toByte())
            // v0.3.0：1 字节档边界对齐 Go 参考（length == 255 落 1 字节档）
            length <= WireConstants.MAX_1 ->
                    writeBytesThenSuffix(data, length.toByte(), (sign or WireConstants.BYTES_LEN_1).toByte())
            length <= WireConstants.MAX_2 ->
                    writeBytesThenSuffix(
                            data,
                            (length shr 8).toByte(),
                            length.toByte(),
                            (sign or WireConstants.BYTES_LEN_2).toByte()
                    )
            length <= WireConstants.MAX_3 ->
                    writeBytesThenSuffix(
                            data,
                            (length shr 16).toByte(),
                            (length shr 8).toByte(),
                            length.toByte(),
                            (sign or WireConstants.BYTES_LEN_3).toByte()
                    )
            // 最后一档必须 <=：length == MAX_4 本身可表达，否则静默丢数据
            else ->
                    writeBytesThenSuffix(
                            data,
                            (length shr 24).toByte(),
                            (length shr 16).toByte(),
                            (length shr 8).toByte(),
                            length.toByte(),
                            (sign or WireConstants.BYTES_LEN_4).toByte()
                    )
        }
    }

    fun encodeArrayPayload(payload: ByteArray): Int =
            encodeContainer(payload, Prefix.CONTAINER or WireConstants.CONTAINER_ARRAY)

    fun encodeObjectPayload(payload: ByteArray): Int =
            encodeContainer(payload, Prefix.CONTAINER or WireConstants.CONTAINER_MAP)

    private fun encodeContainer(payload: ByteArray, baseSign: Int): Int {
        val length = payload.size
        require(length <= WireConstants.MAX_4) { "container payload too long" }
        // 尾定界：[payload][lenBE 档位][sign]
        return when {
            length < WireConstants.CONTAINER_LEN_1 ->
                    writeBytesThenSuffix(payload, (baseSign or length).toByte())
            length <= WireConstants.MAX_1 ->
                    writeBytesThenSuffix(payload, length.toByte(), (baseSign or WireConstants.CONTAINER_LEN_1).toByte())
            length <= WireConstants.MAX_2 ->
                    writeBytesThenSuffix(
                            payload,
                            (length shr 8).toByte(),
                            length.toByte(),
                            (baseSign or WireConstants.CONTAINER_LEN_2).toByte()
                    )
            length <= WireConstants.MAX_3 ->
                    writeBytesThenSuffix(
                            payload,
                            (length shr 16).toByte(),
                            (length shr 8).toByte(),
                            length.toByte(),
                            (baseSign or WireConstants.CONTAINER_LEN_3).toByte()
                    )
            else ->
                    writeBytesThenSuffix(
                            payload,
                            (length shr 24).toByte(),
                            (length shr 16).toByte(),
                            (length shr 8).toByte(),
                            length.toByte(),
                            (baseSign or WireConstants.CONTAINER_LEN_4).toByte()
                    )
        }
    }

    // v0.5.0 尾定界 tag 包装段（一趟正向，零回写），与 Go encode_tag.go 逐字节对齐：
    //   [inner 已编码 payload wire][reverse(meta 正序)][meta_len 尾部槽][total_len 档位大端][tag sign]
    // meta 为 Tag.toBytes() 正序字节；total = len(payload)+len(meta)+len(尾槽)。
    fun encodeTaggedPayload(payload: ByteArray, meta: ByteArray): Int {
        if (meta.isEmpty()) {
            writeBytes(payload)
            return payload.size
        }
        val slot = Footer.footerSlotBytes(meta.size)
        val length = payload.size + meta.size + slot.size
        require(length.toLong() <= WireConstants.MAX_4) { "tag+payload too long" }
        val start = buf.length()
        var sign = Prefix.TAG
        val tier: ByteArray = when {
            length < WireConstants.TAG_LEN_1 -> {
                sign = sign or length
                ByteArray(0)
            }
            length <= WireConstants.MAX_1 -> {
                sign = sign or WireConstants.TAG_LEN_1
                byteArrayOf(length.toByte())
            }
            length <= WireConstants.MAX_2 -> {
                sign = sign or WireConstants.TAG_LEN_2
                byteArrayOf((length shr 8).toByte(), length.toByte())
            }
            length <= WireConstants.MAX_3 -> {
                sign = sign or WireConstants.TAG_LEN_3
                byteArrayOf((length shr 16).toByte(), (length shr 8).toByte(), length.toByte())
            }
            // 最后一档必须 <=：total == MAX_4 本身可表达。
            else -> {
                sign = sign or WireConstants.TAG_LEN_4
                byteArrayOf((length shr 24).toByte(), (length shr 16).toByte(), (length shr 8).toByte(), length.toByte())
            }
        }
        buf.writeAll(payload)
        buf.writeAll(Footer.reverseBytes(meta))
        buf.writeAll(slot)
        buf.writeAll(tier)
        buf.write(sign.toByte())
        return buf.length() - start
    }

    fun encodeBigIntDecimal(s: String): Int {
        val bits = BigIntWireCodec.encodeSignedDecimal(s)
        val inner = ByteArray(1 + bits.size)
        inner[0] = s.length.toByte()
        System.arraycopy(bits, 0, inner, 1, bits.size)
        return encodeBytes(inner)
    }

    fun sliceFrom(start: Int): Int = buf.length() - start

    fun copyLast(n: Int): ByteArray {
        val end = buf.length()
        return buf.copyRange(end - n, end)
    }

    private fun writeByte(vararg bs: Byte): Int {
        buf.write(*bs)
        return bs.size
    }

    private fun writeBytes(bs: ByteArray): Int {
        buf.writeAll(bs)
        return bs.size
    }

    private fun writeBytesThenSuffix(bs: ByteArray, vararg suffix: Byte): Int {
        buf.writeAll(bs)
        for (b in suffix) {
            buf.write(b)
        }
        return bs.size + suffix.size
    }
}
