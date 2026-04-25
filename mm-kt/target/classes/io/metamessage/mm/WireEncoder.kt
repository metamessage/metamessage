package io.metamessage.mm

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
        val start = buf.length()
        buf.write((Prefix.SIMPLE or simpleValue).toByte())
        return buf.length() - start
    }

    fun encodeBool(v: Boolean): Int = encodeSimple(if (v) SimpleValue.TRUE else SimpleValue.FALSE)

    fun encodeInt64(v: Long): Int {
        return if (v >= 0) {
            encodeUintWithPrefix(Prefix.POSITIVE_INT, v)
        } else {
            val uv = if (v == Long.MIN_VALUE) 1L shl 63 else -v
            encodeUintWithPrefix(Prefix.NEGATIVE_INT, uv)
        }
    }

    fun encodeUint64(uv: Long): Int {
        require(uv >= 0) { "expected unsigned" }
        return encodeUintWithPrefix(Prefix.POSITIVE_INT, uv)
    }

    private fun encodeUintWithPrefix(prefix: Int, uv: Long): Int {
        val start = buf.length()
        when {
            uv < WireConstants.INT_LEN_1 -> buf.write((prefix or uv.toInt()).toByte())
            uv <= WireConstants.MAX_1 -> {
                buf.write((prefix or WireConstants.INT_LEN_1).toByte())
                buf.write(uv.toByte())
            }
            uv <= WireConstants.MAX_2 -> {
                buf.write((prefix or WireConstants.INT_LEN_2).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            uv <= WireConstants.MAX_3 -> {
                buf.write((prefix or WireConstants.INT_LEN_3).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            uv <= WireConstants.MAX_4 -> {
                buf.write((prefix or WireConstants.INT_LEN_4).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            uv <= WireConstants.MAX_5 -> {
                buf.write((prefix or WireConstants.INT_LEN_5).toByte())
                buf.write((uv shr 32).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            uv <= WireConstants.MAX_6 -> {
                buf.write((prefix or WireConstants.INT_LEN_6).toByte())
                buf.write((uv shr 40).toByte())
                buf.write((uv shr 32).toByte())
                buf.write((uv shr 24).toByte())
                buf.write((uv shr 16).toByte())
                buf.write((uv shr 8).toByte())
                buf.write(uv.toByte())
            }
            uv <= WireConstants.MAX_7 -> {
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
            buf.write((sign or p.mantissa.toInt()).toByte())
        } else {
            var mantissa = p.mantissa
            when {
                mantissa <= WireConstants.MAX_1 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_1).toByte())
                    buf.write(p.exponent)
                    buf.write(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_2 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_2).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_3 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_3).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 16).toByte())
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_4 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_4).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 24).toByte())
                    buf.write((mantissa shr 16).toByte())
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_5 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_5).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 32).toByte())
                    buf.write((mantissa shr 24).toByte())
                    buf.write((mantissa shr 16).toByte())
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_6 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_6).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 40).toByte())
                    buf.write((mantissa shr 32).toByte())
                    buf.write((mantissa shr 24).toByte())
                    buf.write((mantissa shr 16).toByte())
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
                mantissa <= WireConstants.MAX_7 -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_7).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 48).toByte())
                    buf.write((mantissa shr 40).toByte())
                    buf.write((mantissa shr 32).toByte())
                    buf.write((mantissa shr 24).toByte())
                    buf.write((mantissa shr 16).toByte())
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
                else -> {
                    buf.write((sign or WireConstants.FLOAT_LEN_8).toByte())
                    buf.write(p.exponent)
                    buf.write((mantissa shr 56).toByte())
                    buf.write((mantissa shr 48).toByte())
                    buf.write((mantissa shr 40).toByte())
                    buf.write((mantissa shr 32).toByte())
                    buf.write((mantissa shr 24).toByte())
                    buf.write((mantissa shr 16).toByte())
                    buf.write((mantissa shr 8).toByte())
                    buf.write(mantissa.toByte())
                }
            }
        }
        return buf.length() - start
    }

    fun encodeString(s: String): Int {
        val utf = s.toByteArray(StandardCharsets.UTF_8)
        val length = utf.size
        val start = buf.length()
        val sign = Prefix.STRING
        when {
            length < WireConstants.STRING_LEN_1 -> {
                buf.write((sign or length).toByte())
                buf.writeAll(utf)
            }
            length < WireConstants.MAX_1 -> {
                buf.write((sign or WireConstants.STRING_LEN_1).toByte())
                buf.write(length.toByte())
                buf.writeAll(utf)
            }
            length < WireConstants.MAX_2 -> {
                buf.write((sign or WireConstants.STRING_LEN_2).toByte())
                buf.write((length shr 8).toByte())
                buf.write(length.toByte())
                buf.writeAll(utf)
            }
            else -> throw IllegalArgumentException("string too long")
        }
        return buf.length() - start
    }

    fun encodeBytes(data: ByteArray): Int {
        val length = data.size
        val start = buf.length()
        val sign = Prefix.BYTES
        when {
            length < WireConstants.BYTES_LEN_1 -> {
                buf.write((sign or length).toByte())
                buf.writeAll(data)
            }
            length < WireConstants.MAX_1 -> {
                buf.write((sign or WireConstants.BYTES_LEN_1).toByte())
                buf.write(length.toByte())
                buf.writeAll(data)
            }
            length < WireConstants.MAX_2 -> {
                buf.write((sign or WireConstants.BYTES_LEN_2).toByte())
                buf.write((length shr 8).toByte())
                buf.write(length.toByte())
                buf.writeAll(data)
            }
            else -> throw IllegalArgumentException("bytes too long")
        }
        return buf.length() - start
    }

    fun encodeArrayPayload(payload: ByteArray): Int = encodeContainer(payload, Prefix.CONTAINER or WireConstants.CONTAINER_ARRAY)

    fun encodeMapPayload(payload: ByteArray): Int = encodeContainer(payload, Prefix.CONTAINER or WireConstants.CONTAINER_MAP)

    private fun encodeContainer(payload: ByteArray, baseSign: Int): Int {
        val length = payload.size
        val start = buf.length()
        when {
            length < WireConstants.CONTAINER_LEN_1 -> {
                buf.write((baseSign or length).toByte())
                buf.writeAll(payload)
            }
            length < WireConstants.MAX_1 -> {
                buf.write((baseSign or WireConstants.CONTAINER_LEN_1).toByte())
                buf.write(length.toByte())
                buf.writeAll(payload)
            }
            length < WireConstants.MAX_2 -> {
                buf.write((baseSign or WireConstants.CONTAINER_LEN_2).toByte())
                buf.write((length shr 8).toByte())
                buf.write(length.toByte())
                buf.writeAll(payload)
            }
            else -> throw IllegalArgumentException("container payload too long")
        }
        return buf.length() - start
    }

    fun encodeTagInner(tagBytes: ByteArray): Int {
        if (tagBytes.isEmpty()) return 0
        require(tagBytes.size <= WireConstants.MAX_2) { "tag too long" }
        val start = buf.length()
        val length = tagBytes.size
        when {
            length < 254 -> {
                buf.write(length.toByte())
                buf.writeAll(tagBytes)
            }
            length < 257 -> {
                buf.write(254.toByte())
                buf.write(length.toByte())
                buf.writeAll(tagBytes)
            }
            else -> {
                buf.write(255.toByte())
                buf.write((length shr 8).toByte())
                buf.write(length.toByte())
                buf.writeAll(tagBytes)
            }
        }
        return buf.length() - start
    }

    fun encodeTaggedPayload(payload: ByteArray, rawTagFields: ByteArray): Int {
        if (rawTagFields.isEmpty()) {
            buf.writeAll(payload)
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
}
