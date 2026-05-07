package io.github.metamessage.mm

import java.nio.charset.StandardCharsets
import java.util.Objects

class MmTag(
    var name: String = "",
    var isNull: Boolean = false,
    var example: Boolean = false,
    var desc: String = "",
    var type: ValueType = ValueType.UNKNOWN,
    var raw: Boolean = false,
    var nullable: Boolean = false,
    var allowEmpty: Boolean = false,
    var unique: Boolean = false,
    var defaultValue: String = "",
    var min: String = "",
    var max: String = "",
    var size: Int = 0,
    var enumValues: String = "",
    var pattern: String = "",
    var locationHours: Int = 0,
    var version: Int = DEFAULT_VERSION,
    var mime: String = "",
    var childDesc: String = "",
    var childType: ValueType = ValueType.UNKNOWN,
    var childRaw: Boolean = false,
    var childNullable: Boolean = false,
    var childAllowEmpty: Boolean = false,
    var childUnique: Boolean = false,
    var childDefault: String = "",
    var childMin: String = "",
    var childMax: String = "",
    var childSize: Int = 0,
    var childEnum: String = "",
    var childPattern: String = "",
    var childLocationHours: Int = 0,
    var childVersion: Int = DEFAULT_VERSION,
    var childMime: String = "",
    var isInherit: Boolean = false
) {
    fun copy(): MmTag = MmTag(
        name, isNull, example, desc, type, raw, nullable, allowEmpty, unique,
        defaultValue, min, max, size, enumValues, pattern, locationHours, version, mime,
        childDesc, childType, childRaw, childNullable, childAllowEmpty, childUnique,
        childDefault, childMin, childMax, childSize, childEnum, childPattern,
        childLocationHours, childVersion, childMime, isInherit
    )

    fun inheritFromArrayParent(parent: MmTag?) {
        if (parent == null) return
        desc = parent.childDesc
        type = parent.childType
        raw = parent.childRaw
        nullable = parent.childNullable
        allowEmpty = parent.childAllowEmpty
        unique = parent.childUnique
        defaultValue = parent.childDefault
        min = parent.childMin
        max = parent.childMax
        size = parent.childSize
        enumValues = parent.childEnum
        pattern = parent.childPattern
        locationHours = parent.childLocationHours
        version = parent.childVersion
        mime = parent.childMime
        isInherit = true
    }

    fun toBytes(): ByteArray {
        val w = TagByteWriter()
        if (example) w.writeByte((TagKey.K_EXAMPLE or 1).toByte())
        if (isNull) w.writeByte((TagKey.K_IS_NULL or 1).toByte())
        if (nullable && !isInherit) {
            if (!isNull) {
                w.writeByte((TagKey.K_NULLABLE or 1).toByte())
            }
        }
        if (desc.isNotEmpty() && !isInherit) {
            writeSizedString(w, TagKey.K_DESC, desc)
        }
        if (type != ValueType.UNKNOWN && !isInherit) {
            if (shouldEmitExplicitType()) {
                w.writeByte(TagKey.K_TYPE.toByte())
                w.writeByte(type.code())
            }
        }
        if (raw && !isInherit) w.writeByte((TagKey.K_RAW or 1).toByte())
        if (allowEmpty && !isInherit) w.writeByte((TagKey.K_ALLOW_EMPTY or 1).toByte())
        if (unique && !isInherit) w.writeByte((TagKey.K_UNIQUE or 1).toByte())
        if (defaultValue.isNotEmpty() && !isInherit) {
            writeShortString(w, TagKey.K_DEFAULT, defaultValue)
        }
        if (min.isNotEmpty() && !isInherit) writeShortString(w, TagKey.K_MIN, min)
        if (max.isNotEmpty() && !isInherit) writeShortString(w, TagKey.K_MAX, max)
        if (size != 0 && !isInherit) encodeUint64(w, TagKey.K_SIZE, size.toLong())
        if (enumValues.isNotEmpty() && !isInherit) writeSizedString(w, TagKey.K_ENUM, enumValues)
        if (pattern.isNotEmpty() && !isInherit) writeShortString(w, TagKey.K_PATTERN, pattern)
        if (locationHours != 0 && !isInherit) {
            val v = locationHours.toString()
            w.writeByte((TagKey.K_LOCATION or v.length).toByte())
            w.writeAscii(v)
        }
        if (version != DEFAULT_VERSION && !isInherit) encodeUint64(w, TagKey.K_VERSION, version.toLong())
        if (mime.isNotEmpty() && !isInherit) {
            val m = MimeWire.parse(mime)
            if (m < 7) {
                w.writeByte((TagKey.K_MIME or m).toByte())
            } else {
                w.writeByte((TagKey.K_MIME or 7).toByte())
                w.writeByte(m.toByte())
            }
        }
        if (childDesc.isNotEmpty()) writeSizedString(w, TagKey.K_CHILD_DESC, childDesc)
        if (childType != ValueType.UNKNOWN) {
            if (shouldEmitChildType()) {
                w.writeByte(TagKey.K_CHILD_TYPE.toByte())
                w.writeByte(childType.code())
            }
        }
        if (childRaw) w.writeByte((TagKey.K_CHILD_RAW or 1).toByte())
        if (childNullable) w.writeByte((TagKey.K_CHILD_NULLABLE or 1).toByte())
        if (childAllowEmpty) w.writeByte((TagKey.K_CHILD_ALLOW_EMPTY or 1).toByte())
        if (childUnique) w.writeByte((TagKey.K_CHILD_UNIQUE or 1).toByte())
        if (childDefault.isNotEmpty()) writeShortString(w, TagKey.K_CHILD_DEFAULT, childDefault)
        if (childMin.isNotEmpty()) writeShortString(w, TagKey.K_CHILD_MIN, childMin)
        if (childMax.isNotEmpty()) writeShortString(w, TagKey.K_CHILD_MAX, childMax)
        if (childSize != 0) encodeUint64(w, TagKey.K_CHILD_SIZE, childSize.toLong())
        if (childEnum.isNotEmpty()) writeSizedString(w, TagKey.K_CHILD_ENUM, childEnum)
        if (childPattern.isNotEmpty()) writeShortString(w, TagKey.K_CHILD_PATTERN, childPattern)
        if (childLocationHours != 0) {
            val v = childLocationHours.toString()
            w.writeByte((TagKey.K_CHILD_LOCATION or v.length).toByte())
            w.writeAscii(v)
        }
        if (childVersion != DEFAULT_VERSION) encodeUint64(w, TagKey.K_CHILD_VERSION, childVersion.toLong())
        if (childMime.isNotEmpty()) {
            val m = MimeWire.parse(childMime)
            if (m < 7) {
                w.writeByte((TagKey.K_CHILD_MIME or m).toByte())
            } else {
                w.writeByte((TagKey.K_CHILD_MIME or 7).toByte())
                w.writeByte(m.toByte())
            }
        }
        return w.toByteArray()
    }

    private fun shouldEmitExplicitType(): Boolean {
        return when (type) {
            ValueType.STRING, ValueType.BYTES, ValueType.INT, ValueType.FLOAT64,
            ValueType.BOOL, ValueType.STRUCT, ValueType.SLICE -> false
            ValueType.ARRAY -> if (size > 0) false else true
            ValueType.ENUM -> if (enumValues.isNotEmpty()) false else true
            else -> true
        }
    }

    private fun shouldEmitChildType(): Boolean {
        return when (childType) {
            ValueType.STRING, ValueType.INT, ValueType.FLOAT64,
            ValueType.BOOL, ValueType.STRUCT, ValueType.SLICE -> false
            ValueType.ARRAY -> if (childSize > 0) false else true
            ValueType.ENUM -> if (childEnum.isNotEmpty()) false else true
            else -> true
        }
    }

    private fun writeSizedString(w: TagByteWriter, key: Int, s: String) {
        val b = s.toByteArray(StandardCharsets.UTF_8)
        val l = b.size
        when {
            l <= 5 -> {
                w.writeByte((key or l).toByte())
                w.writeBytes(b)
            }
            l <= 0xFF -> {
                w.writeByte((key or 6).toByte())
                w.writeByte(l.toByte())
                w.writeBytes(b)
            }
            else -> {
                w.writeByte((key or 7).toByte())
                w.writeByte((l shr 8).toByte())
                w.writeByte(l.toByte())
                w.writeBytes(b)
            }
        }
    }

    private fun writeShortString(w: TagByteWriter, key: Int, s: String) {
        val b = s.toByteArray(StandardCharsets.UTF_8)
        val l = b.size
        if (l < 7) {
            w.writeByte((key or l).toByte())
            w.writeBytes(b)
        } else {
            w.writeByte((key or 7).toByte())
            w.writeByte(l.toByte())
            w.writeBytes(b)
        }
    }

    private fun encodeUint64(buf: TagByteWriter, sign: Int, uv: Long) {
        require(uv >= 0) { "unsigned expected" }
        when {
            uv <= WireConstants.MAX_1 -> {
                buf.writeByte(sign.toByte())
                buf.writeByte(uv.toByte())
            }
            uv <= WireConstants.MAX_2 -> {
                buf.writeByte((sign or 1).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
            uv <= WireConstants.MAX_3 -> {
                buf.writeByte((sign or 2).toByte())
                buf.writeByte((uv shr 16).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
            uv <= WireConstants.MAX_4 -> {
                buf.writeByte((sign or 3).toByte())
                buf.writeByte((uv shr 24).toByte())
                buf.writeByte((uv shr 16).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
            uv <= WireConstants.MAX_5 -> {
                buf.writeByte((sign or 4).toByte())
                buf.writeByte((uv shr 32).toByte())
                buf.writeByte((uv shr 24).toByte())
                buf.writeByte((uv shr 16).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
            uv <= WireConstants.MAX_6 -> {
                buf.writeByte((sign or 5).toByte())
                buf.writeByte((uv shr 40).toByte())
                buf.writeByte((uv shr 32).toByte())
                buf.writeByte((uv shr 24).toByte())
                buf.writeByte((uv shr 16).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
            uv <= WireConstants.MAX_7 -> {
                buf.writeByte((sign or 6).toByte())
                buf.writeByte((uv shr 48).toByte())
                buf.writeByte((uv shr 40).toByte())
                buf.writeByte((uv shr 32).toByte())
                buf.writeByte((uv shr 24).toByte())
                buf.writeByte((uv shr 16).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
            else -> {
                buf.writeByte((sign or 7).toByte())
                buf.writeByte((uv shr 56).toByte())
                buf.writeByte((uv shr 48).toByte())
                buf.writeByte((uv shr 40).toByte())
                buf.writeByte((uv shr 32).toByte())
                buf.writeByte((uv shr 24).toByte())
                buf.writeByte((uv shr 16).toByte())
                buf.writeByte((uv shr 8).toByte())
                buf.writeByte(uv.toByte())
            }
        }
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is MmTag) return false
        return isNull == other.isNull &&
                example == other.example &&
                raw == other.raw &&
                nullable == other.nullable &&
                allowEmpty == other.allowEmpty &&
                unique == other.unique &&
                size == other.size &&
                locationHours == other.locationHours &&
                version == other.version &&
                childRaw == other.childRaw &&
                childNullable == other.childNullable &&
                childAllowEmpty == other.childAllowEmpty &&
                childUnique == other.childUnique &&
                childSize == other.childSize &&
                childLocationHours == other.childLocationHours &&
                childVersion == other.childVersion &&
                type == other.type &&
                desc == other.desc &&
                defaultValue == other.defaultValue &&
                min == other.min &&
                max == other.max &&
                enumValues == other.enumValues &&
                pattern == other.pattern &&
                mime == other.mime &&
                childType == other.childType &&
                childDesc == other.childDesc &&
                childDefault == other.childDefault &&
                childMin == other.childMin &&
                childMax == other.childMax &&
                childEnum == other.childEnum &&
                childPattern == other.childPattern &&
                childMime == other.childMime
    }

    override fun hashCode(): Int {
        return Objects.hash(
            name, isNull, example, desc, type, raw, nullable, allowEmpty, unique,
            defaultValue, min, max, size, enumValues, pattern, locationHours, version, mime,
            childDesc, childType, childRaw, childNullable, childAllowEmpty, childUnique,
            childDefault, childMin, childMax, childSize, childEnum, childPattern,
            childLocationHours, childVersion, childMime, isInherit
        )
    }

    private class TagByteWriter {
        private var buf = ByteArray(64)
        private var len = 0

        fun writeByte(b: Byte) {
            ensure(1)
            buf[len++] = b
        }

        fun writeBytes(b: ByteArray) {
            ensure(b.size)
            System.arraycopy(b, 0, buf, len, b.size)
            len += b.size
        }

        fun writeAscii(s: String) {
            val b = s.toByteArray(StandardCharsets.US_ASCII)
            writeBytes(b)
        }

        fun toByteArray(): ByteArray = buf.copyOf(len)

        private fun ensure(n: Int) {
            if (len + n > buf.size) {
                buf = buf.copyOf(Math.max(buf.size * 2, len + n))
            }
        }
    }

    companion object {
        const val DEFAULT_VERSION = 0

        fun empty(): MmTag = MmTag()

        fun fromAnnotation(ann: MM?): MmTag {
            val t = MmTag()
            if (ann == null) return t
            t.name = ann.name
            t.isNull = ann.isNull
            t.example = ann.example
            t.desc = ann.desc
            t.type = ann.type
            t.raw = ann.raw
            t.nullable = ann.nullable
            t.allowEmpty = ann.allowEmpty
            t.unique = ann.unique
            t.defaultValue = ann.defaultValue
            t.min = ann.min
            t.max = ann.max
            t.size = ann.size
            t.enumValues = ann.enumValues
            if (t.enumValues.isNotEmpty()) {
                t.type = ValueType.ENUM
            }
            t.pattern = ann.pattern
            t.locationHours = ann.location
            t.version = ann.version
            t.mime = ann.mime
            t.childDesc = ann.childDesc
            t.childType = ann.childType
            t.childRaw = ann.childRaw
            t.childNullable = ann.childNullable
            t.childAllowEmpty = ann.childAllowEmpty
            t.childUnique = ann.childUnique
            t.childDefault = ann.childDefault
            t.childMin = ann.childMin
            t.childMax = ann.childMax
            t.childSize = ann.childSize
            t.childEnum = ann.childEnum
            if (t.childEnum.isNotEmpty()) {
                t.childType = ValueType.ENUM
            }
            t.childPattern = ann.childPattern
            t.childLocationHours = ann.childLocation
            t.childVersion = ann.childVersion
            t.childMime = ann.childMime
            return t
        }
    }

    object TagKey {
        const val K_IS_NULL = 0 shl 3
        const val K_EXAMPLE = 1 shl 3
        const val K_DESC = 2 shl 3
        const val K_TYPE = 3 shl 3
        const val K_RAW = 4 shl 3
        const val K_NULLABLE = 5 shl 3
        const val K_ALLOW_EMPTY = 6 shl 3
        const val K_UNIQUE = 7 shl 3
        const val K_DEFAULT = 8 shl 3
        const val K_MIN = 9 shl 3
        const val K_MAX = 10 shl 3
        const val K_SIZE = 11 shl 3
        const val K_ENUM = 12 shl 3
        const val K_PATTERN = 13 shl 3
        const val K_LOCATION = 14 shl 3
        const val K_VERSION = 15 shl 3
        const val K_MIME = 16 shl 3
        const val K_CHILD_DESC = 17 shl 3
        const val K_CHILD_TYPE = 18 shl 3
        const val K_CHILD_RAW = 19 shl 3
        const val K_CHILD_NULLABLE = 20 shl 3
        const val K_CHILD_ALLOW_EMPTY = 21 shl 3
        const val K_CHILD_UNIQUE = 22 shl 3
        const val K_CHILD_DEFAULT = 23 shl 3
        const val K_CHILD_MIN = 24 shl 3
        const val K_CHILD_MAX = 25 shl 3
        const val K_CHILD_SIZE = 26 shl 3
        const val K_CHILD_ENUM = 27 shl 3
        const val K_CHILD_PATTERN = 28 shl 3
        const val K_CHILD_LOCATION = 29 shl 3
        const val K_CHILD_VERSION = 30 shl 3
        const val K_CHILD_MIME = 31 shl 3
    }
}
