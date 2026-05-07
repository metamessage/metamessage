package io.github.metamessage.mm

import java.nio.charset.StandardCharsets

object TagFieldParser {
    fun parseOne(c: Cursor, tag: MmTag): Int {
        val start = c.pos
        val b = c.read()
        val prefix = b and 0xF8
        val low = b and 0x07
        when (prefix) {
            MmTag.TagKey.K_IS_NULL -> {
                tag.isNull = (low and 1) == 1
                if (tag.isNull) tag.nullable = true
            }
            MmTag.TagKey.K_EXAMPLE -> tag.example = (low and 1) == 1
            MmTag.TagKey.K_DESC -> readSizedUtf8(c, tag, low, 0)
            MmTag.TagKey.K_TYPE -> tag.type = ValueType.fromCode(c.read())
            MmTag.TagKey.K_RAW -> tag.raw = (low and 1) == 1
            MmTag.TagKey.K_NULLABLE -> tag.nullable = (low and 1) == 1
            MmTag.TagKey.K_ALLOW_EMPTY -> tag.allowEmpty = (low and 1) == 1
            MmTag.TagKey.K_UNIQUE -> tag.unique = (low and 1) == 1
            MmTag.TagKey.K_DEFAULT -> tag.defaultValue = readShortUtf8(c, low)
            MmTag.TagKey.K_MIN -> tag.min = readShortUtf8(c, low)
            MmTag.TagKey.K_MAX -> tag.max = readShortUtf8(c, low)
            MmTag.TagKey.K_SIZE -> tag.size = readUintN(c, low)
            MmTag.TagKey.K_ENUM -> {
                tag.type = ValueType.ENUM
                tag.enumValues = readSizedUtf8Only(c, low)
            }
            MmTag.TagKey.K_PATTERN -> tag.pattern = readShortUtf8(c, low)
            MmTag.TagKey.K_LOCATION -> tag.locationHours = readAscii(c, low).toInt()
            MmTag.TagKey.K_VERSION -> tag.version = readUintN(c, low)
            MmTag.TagKey.K_MIME -> readMime(c, tag, low, true)
            MmTag.TagKey.K_CHILD_DESC -> readSizedUtf8(c, tag, low, 1)
            MmTag.TagKey.K_CHILD_TYPE -> tag.childType = ValueType.fromCode(c.read())
            MmTag.TagKey.K_CHILD_RAW -> tag.childRaw = (low and 1) == 1
            MmTag.TagKey.K_CHILD_NULLABLE -> tag.childNullable = (low and 1) == 1
            MmTag.TagKey.K_CHILD_ALLOW_EMPTY -> tag.childAllowEmpty = (low and 1) == 1
            MmTag.TagKey.K_CHILD_UNIQUE -> tag.childUnique = (low and 1) == 1
            MmTag.TagKey.K_CHILD_DEFAULT -> tag.childDefault = readShortUtf8(c, low)
            MmTag.TagKey.K_CHILD_MIN -> tag.childMin = readShortUtf8(c, low)
            MmTag.TagKey.K_CHILD_MAX -> tag.childMax = readShortUtf8(c, low)
            MmTag.TagKey.K_CHILD_SIZE -> tag.childSize = readUintN(c, low)
            MmTag.TagKey.K_CHILD_ENUM -> {
                tag.childType = ValueType.ENUM
                tag.childEnum = readSizedUtf8Only(c, low)
            }
            MmTag.TagKey.K_CHILD_PATTERN -> tag.childPattern = readShortUtf8(c, low)
            MmTag.TagKey.K_CHILD_LOCATION -> tag.childLocationHours = readAscii(c, low).toInt()
            MmTag.TagKey.K_CHILD_VERSION -> tag.childVersion = readUintN(c, low)
            MmTag.TagKey.K_CHILD_MIME -> readMime(c, tag, low, false)
            else -> throw MmDecodeException("invalid tag field prefix: 0x${Integer.toHexString(prefix)}")
        }
        return c.pos - start
    }

    private fun readSizedUtf8(c: Cursor, tag: MmTag, low: Int, mode: Int) {
        val s = readSizedUtf8Only(c, low)
        if (mode == 0) tag.desc = s else tag.childDesc = s
    }

    private fun readSizedUtf8Only(c: Cursor, low: Int): String {
        return if (low <= 5) {
            String(c.readBytes(low), StandardCharsets.UTF_8)
        } else if (low == 6) {
            val l = c.read()
            String(c.readBytes(l), StandardCharsets.UTF_8)
        } else {
            val hi = c.read()
            val lo = c.read()
            val l = (hi shl 8) or lo
            String(c.readBytes(l), StandardCharsets.UTF_8)
        }
    }

    private fun readShortUtf8(c: Cursor, low: Int): String {
        return if (low < 7) {
            String(c.readBytes(low), StandardCharsets.UTF_8)
        } else {
            val l = c.read()
            String(c.readBytes(l), StandardCharsets.UTF_8)
        }
    }

    private fun readAscii(c: Cursor, low: Int): String {
        return String(c.readBytes(low), StandardCharsets.US_ASCII)
    }

    private fun readUintN(c: Cursor, low: Int): Int {
        require(low < 8) { "uint field length" }
        var v = 0
        for (i in 0 until low) {
            v = (v shl 8) or c.read()
        }
        return v
    }

    private fun readMime(c: Cursor, tag: MmTag, low: Int, self: Boolean) {
        if (low < 7) {
            if (self) tag.mime = MimeWire.toString(low)
            else tag.childMime = MimeWire.toString(low)
        } else {
            val l2 = c.read()
            if (self) tag.mime = MimeWire.toString(l2)
            else tag.childMime = MimeWire.toString(l2)
        }
    }

    class Cursor(var d: ByteArray, var pos: Int) {
        fun read(): Int {
            if (pos >= d.size) throw MmDecodeException("eof")
            return d[pos++].toInt() and 0xFF
        }

        fun readBytes(n: Int): ByteArray {
            if (pos + n > d.size) throw MmDecodeException("eof")
            val r = ByteArray(n)
            System.arraycopy(d, pos, r, 0, n)
            pos += n
            return r
        }
    }
}
