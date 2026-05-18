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

    fun encodeUint64(uv: Long): Int {
        require(uv >= 0) { "expected unsigned" }
        return encodeUintWithPrefix(Prefix.POSITIVE_INT, uv)
    }

    private fun encodeUintWithPrefix(prefix: Int, uv: Long): Int {
        val start = buf.length()
        when {
            java.lang.Long.compareUnsigned(uv, WireConstants.INT_LEN_1.toLong()) < 0 ->
                    buf.write((prefix or uv.toInt()).toByte())
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_1) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_1).toByte())
                buf.write(uv.toByte())
            }
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_2) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_2).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_3) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_3).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_4) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_4).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_5) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_5).toByte())
                buf.write((uv shr 32).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_6) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_6).toByte())
                buf.write((uv shr 40).toByte())
                buf.write((uv shr 32).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            java.lang.Long.compareUnsigned(uv, WireConstants.MAX_7) <= 0 -> {
                buf.write((prefix or WireConstants.INT_LEN_7).toByte())
                buf.write((uv shr 48).toByte())
                buf.write((uv shr 40).toByte())
                buf.write((uv shr 32).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            else -> {
                buf.write((prefix or WireConstants.INT_LEN_8).toByte())
                buf.write((uv shr 56).toByte())
                buf.write((uv shr 48).toByte())
                buf.write((uv shr 40).toByte())
                buf.write((uv shr 32).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
        }
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
            when {
                mantissa <= WireConstants.MAX_1 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_1).toByte())
                    writeByte(p.exponent)
                    writeByte(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_2 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_2).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_3 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_3).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 16).toByte())
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_4 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_4).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 24).toByte())
                    writeByte((mantissa shr 16).toByte())
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_5 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_5).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 32).toByte())
                    writeByte((mantissa shr 24).toByte())
                    writeByte((mantissa shr 16).toByte())
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_6 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_6).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 40).toByte())
                    writeByte((mantissa shr 32).toByte())
                    writeByte((mantissa shr 24).toByte())
                    writeByte((mantissa shr 16).toByte())
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_7 -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_7).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 48).toByte())
                    writeByte((mantissa shr 40).toByte())
                    writeByte((mantissa shr 32).toByte())
                    writeByte((mantissa shr 24).toByte())
                    writeByte((mantissa shr 16).toByte())
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
                else -> {
                    writeByte((sign or WireConstants.FLOAT_LEN_8).toByte())
                    writeByte(p.exponent)
                    writeByte((mantissa shr 56).toByte())
                    writeByte((mantissa shr 48).toByte())
                    writeByte((mantissa shr 40).toByte())
                    writeByte((mantissa shr 32).toByte())
                    writeByte((mantissa shr 24).toByte())
                    writeByte((mantissa shr 16).toByte())
                    writeByte((mantissa shr 8).toByte())
                    writeByte(mantissa.toByte())
                }
            }
        }
        return buf.length() - start
    }

    fun encodeString(s: String): Int {
        val utf = s.toByteArray(StandardCharsets.UTF_8)
        val length = utf.size
        require(length <= WireConstants.MAX_2) { "string too long" }
        val sign = Prefix.STRING
        return when {
            length < WireConstants.STRING_LEN_1 ->
                    writeBytesWithPrefix(utf, (sign or length).toByte())
            length < WireConstants.MAX_1 ->
                    writeBytesWithPrefix(
                            utf,
                            (sign or WireConstants.STRING_LEN_1).toByte(),
                            length.toByte()
                    )
            else ->
                    writeBytesWithPrefix(
                            utf,
                            (sign or WireConstants.STRING_LEN_2).toByte(),
                            (length shr 8).toByte(),
                            length.toByte()
                    )
        }
    }

    fun encodeBytes(data: ByteArray): Int {
        val length = data.size
        require(length <= WireConstants.MAX_2) { "bytes too long" }
        val sign = Prefix.BYTES
        return when {
            length < WireConstants.BYTES_LEN_1 ->
                    writeBytesWithPrefix(data, (sign or length).toByte())
            length < WireConstants.MAX_1 ->
                    writeBytesWithPrefix(
                            data,
                            (sign or WireConstants.BYTES_LEN_1).toByte(),
                            length.toByte()
                    )
            else ->
                    writeBytesWithPrefix(
                            data,
                            (sign or WireConstants.BYTES_LEN_2).toByte(),
                            (length shr 8).toByte(),
                            length.toByte()
                    )
        }
    }

    fun encodeArrayPayload(payload: ByteArray): Int =
            encodeContainer(payload, Prefix.CONTAINER or WireConstants.CONTAINER_ARRAY)

    fun encodeObjectPayload(payload: ByteArray): Int =
            encodeContainer(payload, Prefix.CONTAINER or WireConstants.CONTAINER_MAP)

    private fun encodeContainer(payload: ByteArray, baseSign: Int): Int {
        val length = payload.size
        require(length <= WireConstants.MAX_2) { "container payload too long" }
        return when {
            length < WireConstants.CONTAINER_LEN_1 ->
                    writeBytesWithPrefix(payload, (baseSign or length).toByte())
            length < WireConstants.MAX_1 ->
                    writeBytesWithPrefix(
                            payload,
                            (baseSign or WireConstants.CONTAINER_LEN_1).toByte(),
                            length.toByte()
                    )
            else ->
                    writeBytesWithPrefix(
                            payload,
                            (baseSign or WireConstants.CONTAINER_LEN_2).toByte(),
                            (length shr 8).toByte(),
                            length.toByte()
                    )
        }
    }

    fun encodeTagInner(tagBytes: ByteArray): Int {
        if (tagBytes.isEmpty()) return 0
        require(tagBytes.size <= WireConstants.MAX_2) { "tag too long" }
        val length = tagBytes.size
        return when {
            length < 254 -> writeBytesWithPrefix(tagBytes, length.toByte())
            length < 257 -> writeBytesWithPrefix(tagBytes, 254.toByte(), length.toByte())
            else ->
                    writeBytesWithPrefix(
                            tagBytes,
                            255.toByte(),
                            (length shr 8).toByte(),
                            length.toByte()
                    )
        }
    }

    fun encodeTaggedPayload(payload: ByteArray, rawTagFields: ByteArray): Int {
        if (rawTagFields.isEmpty()) {
            writeBytes(payload)
            return payload.size
        }
        val tEnc = WireEncoder()
        tEnc.encodeTagInner(rawTagFields)
        val tagEncoded = tEnc.toByteArray()
        val length = tagEncoded.size + payload.size
        require(length <= WireConstants.MAX_2) { "tag+payload too long" }
        val start = buf.length()
        val sign = Prefix.TAG
        when {
            length < WireConstants.TAG_LEN_1 -> {
                buf.write((sign or length).toByte())
                buf.writeAll(tagEncoded)
                buf.writeAll(payload)
            }
            length < WireConstants.MAX_1 -> {
                buf.write((sign or WireConstants.TAG_LEN_1).toByte())
                buf.write(length.toByte())
                buf.writeAll(tagEncoded)
                buf.writeAll(payload)
            }
            else -> {
                buf.write((sign or WireConstants.TAG_LEN_2).toByte())
                buf.write((length shr 8).toByte())
                buf.write(length.toByte())
                buf.writeAll(tagEncoded)
                buf.writeAll(payload)
            }
        }
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

    private fun writeBytesWithPrefix(bs: ByteArray, vararg prefix: Byte): Int {
        for (b in prefix) {
            buf.write(b)
        }
        buf.writeAll(bs)
        return prefix.size + bs.size
    }
}
