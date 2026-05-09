package io.github.metamessage.mm

import java.math.BigInteger
import java.nio.ByteBuffer
import java.nio.charset.StandardCharsets
import java.time.Instant
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.ZoneId
import java.time.ZoneOffset
import java.time.format.DateTimeFormatter
import java.util.UUID
import io.github.metamessage.ast.Tag
import io.github.metamessage.ast.ValueType
import io.github.metamessage.ast.Field
import io.github.metamessage.ast.Node
import io.github.metamessage.ast.Value
import io.github.metamessage.ast.Array
import io.github.metamessage.ast.Object

import io.github.metamessage.jsonc.toJsonc

class Decoder() {
    private var offset = 0
    private lateinit var data: ByteArray


    fun decode(data: ByteArray): Node {
        this.data = data
        this.offset = 0

        val d = decodeNode(null)
        if (offset != data.size) {
            throw MmDecodeException("trailing bytes at $offset len ${data.size}")
        }
        // println("22222222${toJsonc(d.node)}")
        return d.node
    }

    private data class Decoded(val node: Node, val consumed: Int)

    private fun decodeNode(inherited: Tag?): Decoded {
        val start = offset
        if (offset >= data.size) {
            throw MmDecodeException("eof")
        }
        val b = data[offset++].toInt() and 0xFF
        val prefix = b and Prefix.PREFIX_MASK
        val d = when (prefix) {
            Prefix.TAG -> decodeTagged(b, start)
            Prefix.SIMPLE -> Decoded(decodeSimple(b, inherited), offset - start)
            Prefix.POSITIVE_INT -> decodePositiveInt(b, inherited, start)
            Prefix.NEGATIVE_INT -> decodeNegativeInt(b, inherited, start)
            Prefix.FLOAT -> decodeFloat(b, inherited, start)
            Prefix.STRING -> decodeString(b, inherited, start)
            Prefix.BYTES -> decodeBytes(b, inherited, start)
            Prefix.CONTAINER -> decodeContainer(b, inherited, start)
            else -> throw MmDecodeException("invalid prefix")
        }
        return d
    }

    private fun decodeTagged(firstByte: Int, start: Int): Decoded {
        val tl = tagOuterLen(firstByte)
        val l1 = tl[0]
        var l2 = tl[1]
        if (l1 == 1) {
            l2 = data[offset++].toInt() and 0xFF
        } else if (l1 == 2) {
            l2 = ((data[offset++].toInt() and 0xFF) shl 8) or (data[offset++].toInt() and 0xFF)
        }
        val innerStart = offset
        val innerEnd = innerStart + l2
        if (innerEnd > data.size) {
            throw MmDecodeException("tag frame past eof")
        }
        val tb = data[offset++].toInt() and 0xFF
        val innerFieldLen = readInnerLen(tb)
        val fieldsEnd = offset + innerFieldLen
        if (fieldsEnd > innerEnd) {
            throw MmDecodeException("tag fields overflow")
        }
        val tag = Tag.empty()
        while (offset < fieldsEnd) {
            val n = TagFieldParser.parseOne(TagFieldParser.Cursor(data, offset), tag)
            if (n <= 0) {
                throw MmDecodeException("tag error")
            }
            offset += n
        }
        val node: Node
        if (tag.isNull) {
            val synthetic = nullScalarForTag(tag)
            if (synthetic != null) {
                node = synthetic
            } else {
                val inner = decodeNode(tag)
                node = inner.node
            }
        } else {
            val inner = decodeNode(tag)
            node = inner.node
        }
        // Ensure offset reaches innerEnd
        offset = innerEnd
        return Decoded(node, offset - start)
    }

    private fun nullScalarForTag(tag: Tag): Node? {
        return when (tag.type) {
            ValueType.DATETIME -> {
                val z = zoneForHours(tag.location)
                val dt = LocalDateTime.ofInstant(Instant.EPOCH, z ?: ZoneOffset.UTC)
                Value(dt, dt.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")), tag)
            }
            ValueType.DATE -> {
                val d = LocalDate.of(1970, 1, 1)
                Value(d, d.toString(), tag)
            }
            ValueType.TIME -> Value(LocalTime.MIDNIGHT, LocalTime.MIDNIGHT.toString(), tag)
            ValueType.INT8 -> Value(0.toByte(), "0", tag)
            ValueType.INT16 -> Value(0.toShort(), "0", tag)
            ValueType.INT32 -> Value(0, "0", tag)
            ValueType.INT64 -> Value(0L, "0", tag)
            ValueType.UINT, ValueType.UINT8, ValueType.UINT16, ValueType.UINT32 -> Value(0, "0", tag)
            ValueType.UINT64 -> Value(0L, "0", tag)
            ValueType.FLOAT32 -> Value(0f, "0.0", tag)
            ValueType.FLOAT64 -> Value(0.0, "0.0", tag)
            ValueType.EMAIL, ValueType.UUID, ValueType.DECIMAL -> Value("", "", tag)
            ValueType.BIGINT -> Value(BigInteger.ZERO, "0", tag)
            ValueType.URL -> Value("", "", tag)
            ValueType.IP -> Value(null, ipNullText(tag.version), tag)
            else -> null
        }
    }

    private fun ipNullText(version: Int): String {
        return when (version) {
            4 -> "0.0.0.0"
            6 -> "::"
            else -> ""
        }
    }

    private fun zoneForHours(hours: Int): ZoneId? {
        return if (hours == 0) null else ZoneOffset.ofHours(hours)
    }

    private fun tagOuterLen(firstByte: Int): IntArray {
        val l = firstByte and WireConstants.TAG_LEN_MASK
        return when {
            l < WireConstants.TAG_LEN_1 -> intArrayOf(0, l)
            l == WireConstants.TAG_LEN_1 -> intArrayOf(1, 0)
            l == WireConstants.TAG_LEN_2 -> intArrayOf(2, 0)
            else -> intArrayOf(0, l)
        }
    }

    private fun readInnerLen(b: Int): Int {
        val l = b
        return when {
            l < 254 -> l
            l < 257 -> data[offset++].toInt() and 0xFF
            else -> ((data[offset++].toInt() and 0xFF) shl 8) or (data[offset++].toInt() and 0xFF)
        }
    }

    private fun decodeSimple(first: Int, inherited: Tag?): Node {
        val tag = inherited?.copy() ?: Tag.empty()
        val sv = first and Prefix.SUFFIX_MASK
        return when (sv) {
            SimpleValue.FALSE -> {
                tag.type = ValueType.BOOL
                Value(false, "false", tag)
            }
            SimpleValue.TRUE -> {
                tag.type = ValueType.BOOL
                Value(true, "true", tag)
            }
            SimpleValue.NULL_BOOL -> nullBool(tag)
            SimpleValue.NULL_INT -> nullInt(tag)
            SimpleValue.NULL_FLOAT -> nullFloat(tag)
            SimpleValue.NULL_STRING -> nullString(tag)
            SimpleValue.NULL_BYTES -> nullBytes(tag)
            else -> throw MmDecodeException("unsupported simple: $sv")
        }
    }

    private fun nullBool(tag: Tag): Node {
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.BOOL
        if (tag.type != ValueType.BOOL) throw MmDecodeException("null_bool type mismatch")
        return Value(false, "false", tag)
    }

    private fun nullInt(tag: Tag): Node {
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.INT
        if (tag.type != ValueType.INT) throw MmDecodeException("null_int type mismatch")
        return Value(0, "0", tag)
    }

    private fun nullFloat(tag: Tag): Node {
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.FLOAT64
        if (tag.type != ValueType.FLOAT32 && tag.type != ValueType.FLOAT64) {
            throw MmDecodeException("null_float type mismatch")
        }
        return if (tag.type == ValueType.FLOAT32) {
            Value(0f, "0.0", tag)
        } else {
            Value(0.0, "0.0", tag)
        }
    }

    private fun nullString(tag: Tag): Node {
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.STRING
        if (tag.type != ValueType.STRING) throw MmDecodeException("null_string type mismatch")
        return Value("", "", tag)
    }

    private fun nullBytes(tag: Tag): Node {
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.BYTES
        if (tag.type != ValueType.BYTES) throw MmDecodeException("null_bytes type mismatch")
        return Value(ByteArray(0), "", tag)
    }

    private fun decodePositiveInt(first: Int, inherited: Tag?, start: Int): Decoded {
        val v = readUintBody(first)
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.INT
        return Decoded(mapUintToTree(tag, v), offset - start)
    }

    private fun decodeNegativeInt(first: Int, inherited: Tag?, start: Int): Decoded {
        val v = readUintBody(first)
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.INT
        return Decoded(mapNegativeInt(tag, v), offset - start)
    }

    private fun readUintBody(first: Int): Long {
        val l1 = intLenExtraBytes(first)
        val low = first and WireConstants.INT_LEN_MASK
        val v = if (l1 == 0) {
            low.toLong()
        } else {
            var v = 0L
            for (i in 0 until l1) {
                v = (v shl 8) or (data[offset++].toInt() and 0xFF).toLong()
            }
            v
        }
        return v
    }

    private fun intLenExtraBytes(first: Int): Int {
        val l = first and WireConstants.INT_LEN_MASK
        return if (l < WireConstants.INT_LEN_1) 0 else l - WireConstants.INT_LEN_1 + 1
    }

    private fun mapUintToTree(tag: Tag, v: Long): Node {
        return when (tag.type) {
            ValueType.INT -> Value(v.toInt(), v.toString(), tag)
            ValueType.INT8 -> Value(v.toByte(), v.toString(), tag)
            ValueType.INT16 -> Value(v.toShort(), v.toString(), tag)
            ValueType.INT32 -> Value(v.toInt(), v.toString(), tag)
            ValueType.INT64 -> Value(v, v.toString(), tag)
            ValueType.UINT -> Value(v.toInt(), v.toString(), tag)
            ValueType.UINT8 -> Value(v.toShort(), v.toString(), tag)
            ValueType.UINT16 -> Value(v.toInt(), v.toString(), tag)
            ValueType.UINT32 -> Value(v.toInt(), v.toString(), tag)
            ValueType.UINT64 -> Value(v, v.toString(), tag)
            ValueType.DATETIME -> decodeDateTime(tag, v)
            ValueType.DATE -> decodeDate(tag, v)
            ValueType.TIME -> decodeTime(tag, v)
            ValueType.ENUM -> decodeEnum(tag, v)
            else -> throw MmDecodeException("unsupported int type: ${tag.type}")
        }
    }

    private fun mapNegativeInt(tag: Tag, v: Long): Node {
        return when (tag.type) {
            ValueType.INT -> Value((-v).toInt(), "-$v", tag)
            ValueType.INT8 -> Value((-v).toByte(), "-$v", tag)
            ValueType.INT16 -> Value((-v).toShort(), "-$v", tag)
            ValueType.INT32 -> Value((-v).toInt(), "-$v", tag)
            ValueType.INT64 -> Value(-v, "-$v", tag)
            else -> throw MmDecodeException("unsupported neg int type: ${tag.type}")
        }
    }

    private fun decodeDateTime(tag: Tag, v: Long): Node {
        if (tag.isNull) return Value(null, "", tag)
        val ins = Instant.ofEpochSecond(v)
        val z = zoneForHours(tag.location)
        val ldt = LocalDateTime.ofInstant(ins, z ?: ZoneOffset.UTC)
        return Value(ldt, ldt.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")), tag)
    }

    private fun decodeDate(tag: Tag, v: Long): Node {
        if (tag.isNull) return Value(null, "", tag)
        if (v > Integer.MAX_VALUE) throw MmDecodeException("date overflow")
        val d = TimeUtil.dateFromDays(v)
        return Value(d, d.toString(), tag)
    }

    private fun decodeTime(tag: Tag, v: Long): Node {
        if (tag.isNull) return Value(null, "", tag)
        if (v > 86399) throw MmDecodeException("time out of range")
        val t = TimeUtil.timeFromSeconds(v.toInt())
        return Value(t, t.toString(), tag)
    }

    private fun decodeEnum(tag: Tag, v: Long): Node {
        if (tag.isNull) return Value(-1, "", tag)
        if (tag.enum.isEmpty()) throw MmDecodeException("enum without labels")
        val parts = tag.enum.split("\\|")
        if (v >= parts.size) throw MmDecodeException("enum index out of range")
        val label = parts[v.toInt()].trim()
        return Value(v.toInt(), label, tag)
    }

    private fun decodeFloat(first: Int, inherited: Tag?, start: Int): Decoded {
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.FLOAT64
        val l = first and WireConstants.FLOAT_LEN_MASK
        val `val`: Double = if (l < WireConstants.FLOAT_LEN_1) {
            var `val` = (first and 0xF).toDouble() / 10.0
            if ((first and WireConstants.FLOAT_NEG_MASK) != 0) `val` = -`val`
            `val`
        } else {
            val exp = data[offset++].toByte()
            val l1 = floatLenExtraBytes(first)
            val mantissa = if (l1 == 0) {
                0
            } else {
                var m = 0L
                for (i in 0 until l1) {
                    m = (m shl 8) or (data[offset++].toInt() and 0xFF).toLong()
                }
                m
            }
            val dec = FloatCodec.mantissaToDecimal(mantissa, exp)
            var `val` = dec.toDouble()
            if ((first and WireConstants.FLOAT_NEG_MASK) != 0) `val` = -`val`
            `val`
        }
        val node = when (tag.type) {
            ValueType.FLOAT32 -> Value(`val`.toFloat(), `val`.toString(), tag)
            ValueType.FLOAT64, ValueType.DECIMAL -> Value(`val`, `val`.toString(), tag)
            else -> throw MmDecodeException("bad float tag ${tag.type}")
        }
        return Decoded(node, offset - start)
    }

    private fun floatLenExtraBytes(first: Int): Int {
        val l = first and WireConstants.FLOAT_LEN_MASK
        return if (l < WireConstants.FLOAT_LEN_1) {
            0
        } else {
            l - WireConstants.FLOAT_LEN_1 + 1
        }
    }

    private fun decodeString(first: Int, inherited: Tag?, start: Int): Decoded {
        val sl = stringLen(first)
        var l2 = sl[1]
        if (sl[0] == 1) {
            l2 = data[offset++].toInt() and 0xFF
        } else if (sl[0] == 2) {
            l2 = ((data[offset++].toInt() and 0xFF) shl 8) or (data[offset++].toInt() and 0xFF)
        }
        val bs = if (l2 > 0) readBytes(l2) else ByteArray(0)
        val text = String(bs, StandardCharsets.UTF_8)
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.STRING
        val node = when (tag.type) {
            ValueType.STRING, ValueType.EMAIL -> Value(text, text, tag)
            ValueType.URL -> Value(text, text, tag)
            ValueType.IP -> Value(text, text, tag)
            else -> throw MmDecodeException("unsupported string type: ${tag.type}")
        }
        return Decoded(node, offset - start)
    }

    private fun stringLen(first: Int): IntArray {
        val l = first and WireConstants.STRING_LEN_MASK
        return when {
            l < WireConstants.STRING_LEN_1 -> intArrayOf(0, l)
            l == WireConstants.STRING_LEN_1 -> intArrayOf(1, 0)
            else -> intArrayOf(2, 0)
        }
    }

    private fun decodeBytes(first: Int, inherited: Tag?, start: Int): Decoded {
        val bl = bytesLen(first)
        var l2 = bl[1]
        if (bl[0] == 1) {
            l2 = data[offset++].toInt() and 0xFF
        } else if (bl[0] == 2) {
            l2 = ((data[offset++].toInt() and 0xFF) shl 8) or (data[offset++].toInt() and 0xFF)
        }
        val bs = if (l2 > 0) readBytes(l2) else ByteArray(0)
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.BYTES
        val node = when (tag.type) {
            ValueType.BYTES -> Value(bs, "", tag)
            ValueType.BIGINT -> bigintFromBytes(bs, tag)
            ValueType.UUID -> {
                if (bs.size != 16) throw MmDecodeException("uuid length")
                val u = uuidFromBytes(bs)
                Value(u, u.toString(), tag)
            }
            ValueType.IP -> Value(bs, "", tag)
            else -> throw MmDecodeException("unsupported bytes type: ${tag.type}")
        }
        return Decoded(node, offset - start)
    }

    private fun bytesLen(first: Int): IntArray {
        val l = first and WireConstants.BYTES_LEN_MASK
        return when {
            l < WireConstants.BYTES_LEN_1 -> intArrayOf(0, l)
            l == WireConstants.BYTES_LEN_1 -> intArrayOf(1, 0)
            else -> intArrayOf(2, 0)
        }
    }

    private fun uuidFromBytes(bs: ByteArray): UUID {
        val bb = ByteBuffer.wrap(bs)
        val hi = bb.long
        val lo = bb.long
        return UUID(hi, lo)
    }

    private fun bigintFromBytes(bs: ByteArray, tag: Tag): Node {
        if (bs.isEmpty()) return Value(BigInteger.ZERO, "0", tag)
        val n = bs[0].toInt() and 0xFF
        val body = bs.copyOfRange(1, bs.size)
        val bits = bigintBits(body)
        val neg = bits.isNotEmpty() && bits[0] == 1
        val digits = BigIntWireCodec.decodePositive(body, n)
        val bi = BigInteger(if (neg) "-$digits" else digits)
        return Value(bi, bi.toString(), tag)
    }

    private fun bigintBits(data: ByteArray): List<Int> {
        val bits = mutableListOf<Int>()
        for (bt in data) {
            for (i in 7 downTo 0) {
                bits.add((bt.toInt() shr i) and 1)
            }
        }
        return bits
    }

    private fun decodeContainer(first: Int, inherited: Tag?, start: Int): Decoded {
        val isArray = (first and WireConstants.CONTAINER_MASK) == WireConstants.CONTAINER_ARRAY
        val isMap = (first and WireConstants.CONTAINER_MASK) == WireConstants.CONTAINER_MAP
        return if (isArray) {
            decodeArray(first, inherited, start)
        } else if (isMap) {
            decodeObject(first, inherited, start)
        } else {
            decodeObject(first, inherited, start)
        }
    }

    private fun decodeArray(first: Int, inherited: Tag?, start: Int): Decoded {
        val cl = containerLen(first)
        var l2 = cl[1]
        if (cl[0] == 1) {
            l2 = data[offset++].toInt() and 0xFF
        } else if (cl[0] == 2) {
            l2 = ((data[offset++].toInt() and 0xFF) shl 8) or (data[offset++].toInt() and 0xFF)
        }
        val bodyStart = offset
        val bodyEnd = bodyStart + l2
        if (bodyEnd > data.size) {
            throw MmDecodeException("array past eof")
        }
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) {
            tag.type = if (tag.size > 0) ValueType.ARRAY else ValueType.SLICE
        }
        val items = mutableListOf<Node>()
        while (offset < bodyEnd) {
            val elemTag = Tag.empty()
            elemTag.inheritFromArrayParent(tag)
            val el = decodeNode(elemTag)
            items.add(el.node)
        }
        if (offset != bodyEnd) {
            throw MmDecodeException("array body misaligned")
        }
        return Decoded(Array(items, tag = tag), offset - start)
    }

    private fun decodeObject(first: Int, inherited: Tag?, start: Int): Decoded {
        val cl = containerLen(first)
        var l2 = cl[1]
        if (cl[0] == 1) {
            l2 = data[offset++].toInt() and 0xFF
        } else if (cl[0] == 2) {
            l2 = ((data[offset++].toInt() and 0xFF) shl 8) or (data[offset++].toInt() and 0xFF)
        }
        val innerStart = offset
        val innerEnd = innerStart + l2
        if (innerEnd > data.size) {
            throw MmDecodeException("object past eof")
        }
        val tag = inherited?.copy() ?: Tag.empty()
        if (tag.type == ValueType.UNKNOWN) tag.type = ValueType.STRUCT
        val keyPrefixPos = offset
        val keyPrefix = data[offset++].toInt() and 0xFF
        val keysDec = decodeArray(keyPrefix, null, keyPrefixPos)
        val keys = keysDec.node as Array
        val fields = mutableListOf<Field>()
        var i = 0
        while (offset < innerEnd && i < keys.items.size) {
            val elemTag = Tag.empty()
            elemTag.inheritFromArrayParent(tag)
            val `val` = decodeNode(elemTag)
            val key = (keys.items[i] as Value).text
            fields.add(Field(key, `val`.node))
            i++
        }
        // Ensure offset reaches innerEnd
        offset = innerEnd
        return Decoded(Object(fields, tag = tag), offset - start)
    }

    private fun containerLen(first: Int): IntArray {
        val l = first and WireConstants.CONTAINER_LEN_MASK
        return when {
            l < WireConstants.CONTAINER_LEN_1 -> intArrayOf(0, l)
            l == WireConstants.CONTAINER_LEN_1 -> intArrayOf(1, 0)
            else -> intArrayOf(2, 0)
        }
    }

    private fun readBytes(n: Int): ByteArray {
        if (offset + n > data.size) throw MmDecodeException("eof")
        val r = ByteArray(n)
        System.arraycopy(data, offset, r, 0, n)
        offset += n
        return r
    }
}
