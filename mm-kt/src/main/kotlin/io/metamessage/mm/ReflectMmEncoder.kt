package io.metamessage.mm

import java.math.BigInteger
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.ZoneOffset
import java.util.UUID

object ReflectMmEncoder {
    fun encode(root: Any): ByteArray {
        val enc = WireEncoder()
        encodeValue(enc, root, rootTagForClass(root.javaClass))
        return enc.toByteArray()
    }

    private fun rootTagForClass(c: Class<*>): MmTag {
        val ann = c.getAnnotation(MM::class.java)
        val t = if (ann != null) MmTag.fromAnnotation(ann) else MmTag.empty()
        if (t.type == ValueType.UNKNOWN) t.type = ValueType.STRUCT
        if (t.name.isEmpty()) t.name = CamelToSnake.convert(c.simpleName)
        return t
    }

    private fun encodeValue(enc: WireEncoder, v: Any?, tag: MmTag) {
        val work = tag.copy()
        if (v == null) {
            if (!work.nullable && !work.isNull) {
                throw IllegalArgumentException("null for non-nullable")
            }
            work.isNull = true
        }
        // Handle wrapper types directly
        if (v is Integer || v is Long || v is Float || v is Double || v is Boolean || v is String) {
            val payload = WireEncoder()
            encodeScalarPayload(payload, v, work)
            enc.encodeTaggedPayload(payload.toByteArray(), work.toBytes())
            return
        }
        if (work.type == ValueType.STRUCT) {
            if (v == null) throw IllegalArgumentException("null struct")
            encodeStruct(enc, v, work)
            return
        }
        if (work.type == ValueType.SLICE || work.type == ValueType.ARRAY) {
            @Suppress("UNCHECKED_CAST")
            encodeList(enc, v as List<*>, work)
            return
        }
        if (work.type == ValueType.MAP) {
            @Suppress("UNCHECKED_CAST")
            encodeMap(enc, v as Map<*, *>, work)
            return
        }
        val payload = WireEncoder()
        encodeScalarPayload(payload, v, work)
        enc.encodeTaggedPayload(payload.toByteArray(), work.toBytes())
    }

    private fun encodeStruct(enc: WireEncoder, o: Any, objTag: MmTag) {
        val keysPacked = GrowableByteBuf()
        val valsPacked = GrowableByteBuf()
        val tmp = WireEncoder()
        for (f in o.javaClass.declaredFields) {
            if (java.lang.reflect.Modifier.isStatic(f.modifiers)) continue
            f.isAccessible = true
            val mm = f.getAnnotation(MM::class.java)
            if (mm != null && mm.name == "-") continue
            var ft = TypeInference.forField(f)
            if (mm != null) {
                val ann = MmTag.fromAnnotation(mm)
                if (ann.type != ValueType.UNKNOWN) ft.type = ann.type
                if (ann.childType != ValueType.UNKNOWN) ft.childType = ann.childType
                if (ann.name.isNotEmpty()) ft.name = ann.name
                mergeAnnotations(ft, ann)
            }
            val key = fieldKey(f, ft, mm)
            val fv = f.get(o)
            tmp.reset()
            encodeValue(tmp, fv, ft)
            valsPacked.writeAll(tmp.toByteArray())
            tmp.reset()
            tmp.encodeString(key)
            keysPacked.writeAll(tmp.toByteArray())
        }
        tmp.reset()
        tmp.encodeArrayPayload(keysPacked.copyRange(0, keysPacked.length()))
        val mapBody = GrowableByteBuf()
        mapBody.writeAll(tmp.toByteArray())
        mapBody.writeAll(valsPacked.copyRange(0, valsPacked.length()))
        tmp.reset()
        tmp.encodeMapPayload(mapBody.copyRange(0, mapBody.length()))
        enc.encodeTaggedPayload(tmp.toByteArray(), objTag.toBytes())
    }

    private fun mergeAnnotations(dst: MmTag, src: MmTag) {
        if (src.desc.isNotEmpty()) dst.desc = src.desc
        dst.nullable = dst.nullable || src.nullable
        dst.raw = dst.raw || src.raw
        dst.allowEmpty = dst.allowEmpty || src.allowEmpty
        dst.unique = dst.unique || src.unique
        if (src.defaultValue.isNotEmpty()) dst.defaultValue = src.defaultValue
        if (src.enumValues.isNotEmpty()) {
            dst.enumValues = src.enumValues
            dst.type = ValueType.ENUM
        }
        dst.locationHours = src.locationHours
        dst.version = src.version
        if (src.mime.isNotEmpty()) dst.mime = src.mime
        dst.childDesc = src.childDesc
        if (src.childType != ValueType.UNKNOWN) dst.childType = src.childType
        dst.childNullable = dst.childNullable || src.childNullable
        if (src.childEnum.isNotEmpty()) {
            dst.childEnum = src.childEnum
            dst.childType = ValueType.ENUM
        }
    }

    private fun fieldKey(f: java.lang.reflect.Field, ft: MmTag, mm: MM?): String {
        if (mm != null && mm.name.isNotEmpty() && mm.name != "-") return mm.name
        if (ft.name.isNotEmpty() && ft.name != "-") return ft.name
        return CamelToSnake.convert(f.name)
    }

    private fun encodeScalarPayload(w: WireEncoder, v: Any?, tag: MmTag) {
        if (tag.isNull && tag.type == ValueType.INT) {
            w.encodeSimple(SimpleValue.NULL_INT)
            return
        }
        if (tag.isNull && tag.type == ValueType.STRING) {
            w.encodeSimple(SimpleValue.NULL_STRING)
            return
        }
        if (tag.isNull && tag.type == ValueType.BYTES) {
            w.encodeSimple(SimpleValue.NULL_BYTES)
            return
        }
        if (tag.isNull && (tag.type == ValueType.FLOAT32 || tag.type == ValueType.FLOAT64)) {
            w.encodeSimple(SimpleValue.NULL_FLOAT)
            return
        }
        when (tag.type) {
            ValueType.BOOL -> w.encodeBool(v == true)
            ValueType.INT -> w.encodeInt64((v as Number).toInt().toLong())
            ValueType.INT8 -> w.encodeInt64((v as Number).toByte().toLong())
            ValueType.INT16 -> w.encodeInt64((v as Number).toShort().toLong())
            ValueType.INT32, ValueType.UINT, ValueType.UINT16 -> w.encodeInt64((v as Number).toInt().toLong())
            ValueType.INT64, ValueType.UINT32, ValueType.UINT64 -> w.encodeInt64((v as Number).toLong())
            ValueType.FLOAT32 -> w.encodeFloatString((v as Number).toDouble().toString())
            ValueType.FLOAT64, ValueType.DECIMAL -> w.encodeFloatString((v as Number).toDouble().toString())
            ValueType.STRING, ValueType.EMAIL, ValueType.URL -> w.encodeString(v?.toString() ?: "")
            ValueType.BYTES -> w.encodeBytes(v as? ByteArray ?: ByteArray(0))
            ValueType.BIGINT -> w.encodeBigIntDecimal(v?.toString() ?: "0")
            ValueType.UUID -> w.encodeBytes(uuidBytes(v as UUID))
            ValueType.DATETIME -> {
                val sec = when (v) {
                    is LocalDateTime -> v.toEpochSecond(ZoneOffset.UTC)
                    is Instant -> v.epochSecond
                    else -> throw IllegalArgumentException("datetime")
                }
                w.encodeInt64(sec)
            }
            ValueType.DATE -> w.encodeInt64(TimeUtil.daysSinceEpochUtc(v as LocalDate))
            ValueType.TIME -> w.encodeInt64((v as LocalTime).toSecondOfDay().toLong())
            ValueType.ENUM -> w.encodeInt64((v as Number).toInt().toLong())
            else -> throw UnsupportedOperationException("scalar ${tag.type}")
        }
    }

    private fun encodeList(enc: WireEncoder, list: List<*>?, tag: MmTag) {
        if (list == null) {
            val nt = tag.copy()
            nt.isNull = true
            enc.encodeTaggedPayload(ByteArray(0), nt.toBytes())
            return
        }
        val body = GrowableByteBuf()
        val el = WireEncoder()
        for (item in list) {
            el.reset()
            val et = MmTag.empty()
            et.inheritFromArrayParent(tag)
            if (et.type == ValueType.UNKNOWN && item != null) {
                et.type = TypeInference.valueTypeForComponent(item.javaClass)
                if (et.type == ValueType.UNKNOWN) et.type = ValueType.STRUCT
            }
            encodeValue(el, item, et)
            body.writeAll(el.toByteArray())
        }
        el.reset()
        el.encodeArrayPayload(body.copyRange(0, body.length()))
        enc.encodeTaggedPayload(el.toByteArray(), tag.toBytes())
    }

    private fun encodeMap(enc: WireEncoder, map: Map<*, *>?, tag: MmTag) {
        if (map == null) {
            val nt = tag.copy()
            nt.isNull = true
            enc.encodeTaggedPayload(ByteArray(0), nt.toBytes())
            return
        }
        val keysPacked = GrowableByteBuf()
        val valsPacked = GrowableByteBuf()
        val tmp = WireEncoder()
        for ((key, value) in map) {
            // Encode key
            tmp.reset()
            tmp.encodeString(key?.toString() ?: "")
            keysPacked.writeAll(tmp.toByteArray())
            
            // Encode value
            tmp.reset()
            val valueTag = MmTag.empty()
            // For map values, we'll use the childType if available
            if (tag.childType != ValueType.UNKNOWN) {
                valueTag.type = tag.childType
            } else if (value != null) {
                valueTag.type = TypeInference.valueTypeForComponent(value.javaClass)
                if (valueTag.type == ValueType.UNKNOWN) {
                    valueTag.type = ValueType.STRUCT
                }
            }
            // Handle wrapper types directly
            if (value is Int || value is Long || value is Float || value is Double || value is Boolean || value is String) {
                when (value) {
                    is Int -> valueTag.type = ValueType.INT
                    is Long -> valueTag.type = ValueType.INT64
                    is Float -> valueTag.type = ValueType.FLOAT32
                    is Double -> valueTag.type = ValueType.FLOAT64
                    is Boolean -> valueTag.type = ValueType.BOOL
                    is String -> valueTag.type = ValueType.STRING
                }
            }
            encodeValue(tmp, value, valueTag)
            valsPacked.writeAll(tmp.toByteArray())
        }
        tmp.reset()
        tmp.encodeArrayPayload(keysPacked.copyRange(0, keysPacked.length()))
        val mapBody = GrowableByteBuf()
        mapBody.writeAll(tmp.toByteArray())
        mapBody.writeAll(valsPacked.copyRange(0, valsPacked.length()))
        tmp.reset()
        tmp.encodeMapPayload(mapBody.copyRange(0, mapBody.length()))
        enc.encodeTaggedPayload(tmp.toByteArray(), tag.toBytes())
    }

    private fun uuidBytes(u: UUID): ByteArray {
        val bb = java.nio.ByteBuffer.allocate(16)
        bb.putLong(u.mostSignificantBits)
        bb.putLong(u.leastSignificantBits)
        return bb.array()
    }
}
