package io.github.metamessage.core

import io.github.metamessage.MM
import io.github.metamessage.ir.Array as AstArray
import io.github.metamessage.ir.Node
import io.github.metamessage.ir.Object as AstObject
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.Value as AstValue
import io.github.metamessage.ir.ValueType
import java.time.LocalDate
import java.time.LocalTime
import java.util.UUID

object Encoder {

    fun encode(value: Any): ByteArray {
        val node = valueToNode(value, rootTagForClass(value.javaClass), "")
        return encodeNode(node)
    }

    fun encodeNode(node: Node): ByteArray {
        val enc = WireEncoder()
        encodeNodeInternal(enc, node)
        return enc.toByteArray()
    }

    private fun encodeNodeInternal(enc: WireEncoder, node: Node) {
        when (node) {
            is AstObject -> encodeObjectNode(enc, node)
            is AstArray -> encodeArrayNode(enc, node)
            is AstValue -> encodeValueNode(enc, node)
            is io.github.metamessage.ir.Doc -> encodeDocNode(enc, node)
        }
    }

    private fun encodeObjectNode(enc: WireEncoder, obj: AstObject) {
        val keysPacked = GrowableByteBuf()
        val valsPacked = GrowableByteBuf()
        val tmp = WireEncoder()
        for (field in obj.fields) {
            tmp.reset()
            encodeNodeInternal(tmp, field.value)
            valsPacked.writeAll(tmp.toByteArray())

            tmp.reset()
            tmp.encodeString(field.key)
            keysPacked.writeAll(tmp.toByteArray())
        }
        tmp.reset()
        tmp.encodeArrayPayload(keysPacked.copyRange(0, keysPacked.length()))
        val mapBody = GrowableByteBuf()
        mapBody.writeAll(tmp.toByteArray())
        mapBody.writeAll(valsPacked.copyRange(0, valsPacked.length()))
        tmp.reset()
        tmp.encodeObjectPayload(mapBody.copyRange(0, mapBody.length()))
        enc.encodeTaggedPayload(tmp.toByteArray(), obj.tag?.toBytes() ?: ByteArray(0))
    }

    private fun encodeDocNode(enc: WireEncoder, doc: io.github.metamessage.ir.Doc) {
        val keysPacked = GrowableByteBuf()
        val valsPacked = GrowableByteBuf()
        val tmp = WireEncoder()
        for (field in doc.fields) {
            tmp.reset()
            encodeNodeInternal(tmp, field.value)
            valsPacked.writeAll(tmp.toByteArray())

            tmp.reset()
            tmp.encodeString(field.key)
            keysPacked.writeAll(tmp.toByteArray())
        }
        tmp.reset()
        tmp.encodeArrayPayload(keysPacked.copyRange(0, keysPacked.length()))
        val mapBody = GrowableByteBuf()
        mapBody.writeAll(tmp.toByteArray())
        mapBody.writeAll(valsPacked.copyRange(0, valsPacked.length()))
        tmp.reset()
        tmp.encodeObjectPayload(mapBody.copyRange(0, mapBody.length()))
        enc.encodeTaggedPayload(tmp.toByteArray(), doc.tag?.toBytes() ?: ByteArray(0))
    }

    private fun encodeArrayNode(enc: WireEncoder, arr: AstArray) {
        val body = GrowableByteBuf()
        val el = WireEncoder()
        for (item in arr.items) {
            el.reset()
            encodeNodeInternal(el, item)
            body.writeAll(el.toByteArray())
        }
        el.reset()
        el.encodeArrayPayload(body.copyRange(0, body.length()))
        enc.encodeTaggedPayload(el.toByteArray(), arr.tag?.toBytes() ?: ByteArray(0))
    }

    private fun encodeValueNode(enc: WireEncoder, value: AstValue) {
        val tag = value.tag ?: Tag.empty()
        val tmp = WireEncoder()
        when (tag.type) {
            ValueType.DATETIME -> {
                if (!tag.isNull) {
                    val timeValue = value.data ?: throw IllegalArgumentException("null datetime")
                    tmp.encodeInt64(TimeUtil.epochSeconds(timeValue))
                }
            }
            ValueType.DATE -> {
                if (!tag.isNull) {
                    val dateValue =
                            value.data as? LocalDate
                                    ?: throw IllegalArgumentException("invalid date")
                    tmp.encodeInt64(TimeUtil.daysSinceEpochUtc(dateValue))
                }
            }
            ValueType.TIME -> {
                if (!tag.isNull) {
                    val timeValue =
                            value.data as? LocalTime
                                    ?: throw IllegalArgumentException("invalid time")
                    tmp.encodeInt64(TimeUtil.secondsOfDay(timeValue).toLong())
                }
            }
            ValueType.I -> {
                if (tag.isNull) {
                    tmp.encodeSimple(SimpleValue.NULL_INT)
                } else {
                    tmp.encodeInt64((value.data as Number).toLong())
                }
            }
            ValueType.I8 -> {
                if (!tag.isNull) {
                    tmp.encodeInt64((value.data as Number).toLong())
                }
            }
            ValueType.I16 -> {
                if (!tag.isNull) {
                    tmp.encodeInt64((value.data as Number).toLong())
                }
            }
            ValueType.I32 -> {
                if (!tag.isNull) {
                    tmp.encodeInt64((value.data as Number).toLong())
                }
            }
            ValueType.I64 -> {
                if (!tag.isNull) {
                    tmp.encodeInt64((value.data as Number).toLong())
                }
            }
            ValueType.U -> {
                if (!tag.isNull) {
                    tmp.encodeU64((value.data as Number).toLong())
                }
            }
            ValueType.U8 -> {
                if (!tag.isNull) {
                    tmp.encodeU64((value.data as Number).toLong())
                }
            }
            ValueType.U16 -> {
                if (!tag.isNull) {
                    tmp.encodeU64((value.data as Number).toLong())
                }
            }
            ValueType.U32 -> {
                if (!tag.isNull) {
                    tmp.encodeU64((value.data as Number).toLong())
                }
            }
            ValueType.U64 -> {
                if (!tag.isNull) {
                    tmp.encodeU64(value.data as Long)
                }
            }
            ValueType.F32 -> {
                if (!tag.isNull) {
                    tmp.encodeFloatString(value.text)
                }
            }
            ValueType.F64 -> {
                if (tag.isNull) {
                    tmp.encodeSimple(SimpleValue.NULL_FLOAT)
                } else {
                    tmp.encodeFloatString(value.text)
                }
            }
            ValueType.STR -> {
                if (tag.isNull) {
                    tmp.encodeSimple(SimpleValue.NULL_STRING)
                } else {
                    tmp.encodeString(value.text)
                }
            }
            ValueType.EMAIL -> {
                if (!tag.isNull) {
                    tmp.encodeString(value.text)
                }
            }
            ValueType.UUID -> {
                if (!tag.isNull) {
                    val uuid = value.data as? UUID ?: throw IllegalArgumentException("invalid uuid")
                    tmp.encodeBytes(uuidBytes(uuid))
                }
            }
            ValueType.DECIMAL -> {
                if (!tag.isNull) {
                    tmp.encodeFloatString(value.text)
                }
            }
            ValueType.URL -> {
                if (!tag.isNull) {
                    tmp.encodeString(value.text)
                }
            }
            ValueType.IP -> {
                if (!tag.isNull) {
                    val ip =
                            value.data as? java.net.InetAddress
                                    ?: throw IllegalArgumentException("invalid ip")
                    when (tag.version) {
                        0 -> tmp.encodeString(value.text)
                        4 -> tmp.encodeBytes(ip.address)
                        6 ->
                                if (value.text.length < 16) {
                                    tmp.encodeString(value.text)
                                } else {
                                    tmp.encodeBytes(ip.address)
                                }
                        else ->
                                throw IllegalArgumentException(
                                        "unsupported IP version: ${tag.version}"
                                )
                    }
                }
            }
            ValueType.BYTES -> {
                if (tag.isNull) {
                    tmp.encodeSimple(SimpleValue.NULL_BYTES)
                } else {
                    val bytes = value.data as? ByteArray ?: ByteArray(0)
                    tmp.encodeBytes(bytes)
                }
            }
            ValueType.BIGINT -> {
                if (!tag.isNull) {
                    tmp.encodeBigIntDecimal(value.text)
                }
            }
            ValueType.BOOL -> {
                if (tag.isNull) {
                    tmp.encodeSimple(SimpleValue.NULL_BOOL)
                } else {
                    val boolValue = value.data as? Boolean ?: false
                    tmp.encodeBool(boolValue)
                }
            }
            ValueType.ENUM -> {
                if (!tag.isNull) {
                    tmp.encodeInt64((value.data as Number).toLong())
                }
            }
            else -> {
                if (tag.isNull) {
                    tmp.encodeSimple(SimpleValue.NULL_STRING)
                } else {
                    tmp.encodeString(value.text)
                }
            }
        }
        enc.encodeTaggedPayload(tmp.toByteArray(), tag.toBytes())
    }

    fun rootTagForClass(c: Class<*>): Tag {
        val ann = c.getAnnotation(MM::class.java)
        val t = if (ann != null) Tag.fromAnnotation(ann) else Tag.empty()
        if (t.type == ValueType.UNKNOWN) t.type = ValueType.OBJ
        if (t.name.isEmpty()) t.name = CamelToSnake.convert(c.simpleName)
        return t
    }

    private fun uuidBytes(u: UUID): ByteArray {
        val bb = java.nio.ByteBuffer.allocate(16)
        bb.putLong(u.mostSignificantBits)
        bb.putLong(u.leastSignificantBits)
        return bb.array()
    }
}
