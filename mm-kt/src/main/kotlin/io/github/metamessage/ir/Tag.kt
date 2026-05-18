package io.github.metamessage.ir

import io.github.metamessage.MM
import io.github.metamessage.core.WireConstants
import java.math.BigInteger
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.time.format.DateTimeFormatter
import java.util.Base64
import java.util.Objects
import java.util.regex.Pattern

class Tag(
        var name: String = "",
        var isNull: Boolean = false,
        var example: Boolean = false,
        var desc: String = "",
        var type: ValueType = ValueType.UNKNOWN,
        var raw: Boolean = false,
        var nullable: Boolean = false,
        var allowEmpty: Boolean = false,
        var unique: Boolean = false,
        var default: String = "",
        var min: String = "",
        var max: String = "",
        var size: Int = 0,
        var enum: String = "",
        var pattern: String = "",
        var location: Int = 0,
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
        var childLocation: Int = 0,
        var childVersion: Int = DEFAULT_VERSION,
        var childMime: String = "",
        var isInherit: Boolean = false
) {
    fun copy(): Tag =
            Tag(
                    name,
                    isNull,
                    example,
                    desc,
                    type,
                    raw,
                    nullable,
                    allowEmpty,
                    unique,
                    default,
                    min,
                    max,
                    size,
                    enum,
                    pattern,
                    location,
                    version,
                    mime,
                    childDesc,
                    childType,
                    childRaw,
                    childNullable,
                    childAllowEmpty,
                    childUnique,
                    childDefault,
                    childMin,
                    childMax,
                    childSize,
                    childEnum,
                    childPattern,
                    childLocation,
                    childVersion,
                    childMime,
                    isInherit
            )

    fun inheritFromArrayParent(parent: Tag?) {
        if (parent == null) return

        isInherit = true

        if (parent.childDesc.isNotEmpty()) {
            desc = parent.childDesc
        }

        if (parent.childType != ValueType.UNKNOWN) {
            type = parent.childType
        }

        if (parent.childRaw) {
            raw = parent.childRaw
        }

        if (parent.childNullable) {
            nullable = parent.childNullable
        }

        if (parent.childAllowEmpty) {
            allowEmpty = parent.childAllowEmpty
        }

        if (parent.childUnique) {
            unique = parent.childUnique
        }

        if (parent.childDefault.isNotEmpty()) {
            default = parent.childDefault
        }

        if (parent.childMin.isNotEmpty()) {
            min = parent.childMin
        }

        if (parent.childMax.isNotEmpty()) {
            max = parent.childMax
        }

        if (parent.childSize != 0) {
            size = parent.childSize
        }

        if (parent.childEnum.isNotEmpty()) {
            enum = parent.childEnum
        }

        if (parent.childPattern.isNotEmpty()) {
            pattern = parent.childPattern
        }

        if (parent.childLocation != 0) {
            location = parent.childLocation
        }

        if (parent.childVersion != DEFAULT_VERSION) {
            version = parent.childVersion
        }

        if (parent.childMime.isNotEmpty()) {
            mime = parent.childMime
        }
    }

    fun getPattern(): Pattern? {
        if (pattern.isEmpty()) return null
        return Pattern.compile(pattern)
    }

    override fun toString(): String {
        val parts = mutableListOf<String>()

        if (type != ValueType.UNKNOWN && !isInherit) {
            if (!(type == ValueType.STRING ||
                            type == ValueType.INT ||
                            type == ValueType.FLOAT64 ||
                            type == ValueType.BOOL ||
                            type == ValueType.STRUCT ||
                            type == ValueType.SLICE)
            ) {
                if (!((type == ValueType.ARRAY && size > 0) ||
                                (type == ValueType.ENUM && enum.isNotEmpty()))
                ) {
                    parts.add("${T_TYPE}=${type.toString()}")
                }
            }
        }

        if (example) {
            parts.add(T_EXAMPLE)
        }

        if (isNull) {
            parts.add(T_IS_NULL)
        }

        if (nullable && !isInherit) {
            if (!isNull) {
                parts.add(T_NULLABLE)
            }
        }

        if (desc.isNotEmpty() && !isInherit) {
            parts.add("${T_DESC}=${quote(desc)}")
        }

        if (raw && !isInherit) {
            parts.add(T_RAW)
        }

        if (allowEmpty && !isInherit) {
            parts.add(T_ALLOW_EMPTY)
        }

        if (unique && !isInherit) {
            parts.add(T_UNIQUE)
        }

        if (default.isNotEmpty() && !isInherit) {
            parts.add("${T_DEFAULT}=${default}")
        }

        if (min.isNotEmpty() && !isInherit) {
            parts.add("${T_MIN}=${min}")
        }

        if (max.isNotEmpty() && !isInherit) {
            parts.add("${T_MAX}=${max}")
        }

        if (size != 0 && !isInherit) {
            parts.add("${T_SIZE}=${size}")
        }

        if (enum.isNotEmpty() && !isInherit) {
            parts.add("${T_ENUM}=${enum}")
        }

        if (pattern.isNotEmpty() && !isInherit) {
            parts.add("${T_PATTERN}=${pattern}")
        }

        if (location != 0 && !isInherit) {
            parts.add("${T_LOCATION}=${location}")
        }

        if (version != DEFAULT_VERSION && !isInherit) {
            parts.add("${T_VERSION}=${version}")
        }

        if (mime.isNotEmpty() && !isInherit) {
            parts.add("${T_MIME}=${mime}")
        }

        if (childDesc.isNotEmpty()) {
            parts.add("${T_CHILD_DESC}=${quote(childDesc)}")
        }

        if (childType != ValueType.UNKNOWN) {
            if (!(childType == ValueType.STRING ||
                            childType == ValueType.INT ||
                            childType == ValueType.FLOAT64 ||
                            childType == ValueType.BOOL ||
                            childType == ValueType.STRUCT ||
                            childType == ValueType.SLICE)
            ) {
                if (!((childType == ValueType.ARRAY && childSize > 0) ||
                                (childType == ValueType.ENUM && childEnum.isNotEmpty()))
                ) {
                    parts.add("${T_CHILD_TYPE}=${childType.toString()}")
                }
            }
        }

        if (childRaw) {
            parts.add(T_CHILD_RAW)
        }

        if (childNullable) {
            parts.add(T_CHILD_NULLABLE)
        }

        if (childAllowEmpty) {
            parts.add(T_CHILD_ALLOW_EMPTY)
        }

        if (childUnique) {
            parts.add(T_CHILD_UNIQUE)
        }

        if (childDefault.isNotEmpty()) {
            parts.add("${T_CHILD_DEFAULT}=${childDefault}")
        }

        if (childMin.isNotEmpty()) {
            parts.add("${T_CHILD_MIN}=${childMin}")
        }

        if (childMax.isNotEmpty()) {
            parts.add("${T_CHILD_MAX}=${childMax}")
        }

        if (childSize != 0) {
            parts.add("${T_CHILD_SIZE}=${childSize}")
        }

        if (childEnum.isNotEmpty()) {
            parts.add("${T_CHILD_ENUM}=${childEnum}")
        }

        if (childPattern.isNotEmpty()) {
            parts.add("${T_CHILD_PATTERN}=${childPattern}")
        }

        if (childLocation != 0) {
            parts.add("${T_CHILD_LOCATION}=${childLocation}")
        }

        if (childVersion != DEFAULT_VERSION) {
            parts.add("${T_CHILD_VERSION}=${childVersion}")
        }

        if (childMime.isNotEmpty()) {
            parts.add("${T_CHILD_MIME}=${childMime}")
        }

        return parts.joinToString("; ")
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
        if (default.isNotEmpty() && !isInherit) {
            writeShortString(w, TagKey.K_DEFAULT, default)
        }
        if (min.isNotEmpty() && !isInherit) writeShortString(w, TagKey.K_MIN, min)
        if (max.isNotEmpty() && !isInherit) writeShortString(w, TagKey.K_MAX, max)
        if (size != 0 && !isInherit) encodeUint64(w, TagKey.K_SIZE, size.toLong())
        if (enum.isNotEmpty() && !isInherit) writeSizedString(w, TagKey.K_ENUM, enum)
        if (pattern.isNotEmpty() && !isInherit) writeShortString(w, TagKey.K_PATTERN, pattern)
        if (location != 0 && !isInherit) {
            val v = location.toString()
            w.writeByte((TagKey.K_LOCATION or v.length).toByte())
            w.writeAscii(v)
        }
        if (version != DEFAULT_VERSION && !isInherit)
                encodeUint64(w, TagKey.K_VERSION, version.toLong())
        if (mime.isNotEmpty() && !isInherit) {
            val m = Mime.parse(mime)
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
        if (childLocation != 0) {
            val v = childLocation.toString()
            w.writeByte((TagKey.K_CHILD_LOCATION or v.length).toByte())
            w.writeAscii(v)
        }
        if (childVersion != DEFAULT_VERSION)
                encodeUint64(w, TagKey.K_CHILD_VERSION, childVersion.toLong())
        if (childMime.isNotEmpty()) {
            val l = childMime.length
            if (l < 7) {
                w.writeByte((TagKey.K_CHILD_MIME or l).toByte())
                w.writeBytes(childMime.toByteArray(charset("UTF-8")))
            } else {
                w.writeByte((TagKey.K_CHILD_MIME or 7).toByte())
                w.writeByte(l.toByte())
                w.writeBytes(childMime.toByteArray(charset("UTF-8")))
            }
        }
        return w.toByteArray()
    }

    private fun shouldEmitExplicitType(): Boolean {
        return when (type) {
            ValueType.STRING,
            ValueType.BYTES,
            ValueType.INT,
            ValueType.FLOAT64,
            ValueType.BOOL,
            ValueType.STRUCT,
            ValueType.SLICE -> false
            ValueType.ARRAY -> if (size > 0) false else true
            ValueType.ENUM -> if (enum.isNotEmpty()) false else true
            else -> true
        }
    }

    private fun shouldEmitChildType(): Boolean {
        return when (childType) {
            ValueType.STRING,
            ValueType.INT,
            ValueType.FLOAT64,
            ValueType.BOOL,
            ValueType.STRUCT,
            ValueType.SLICE -> false
            ValueType.ARRAY -> if (childSize > 0) false else true
            ValueType.ENUM -> if (childEnum.isNotEmpty()) false else true
            else -> true
        }
    }

    private fun writeSizedString(w: TagByteWriter, key: Int, s: String) {
        val b = s.toByteArray(charset("UTF-8"))
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
        val b = s.toByteArray(charset("UTF-8"))
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
        if (other !is Tag) return false
        return isNull == other.isNull &&
                example == other.example &&
                raw == other.raw &&
                nullable == other.nullable &&
                allowEmpty == other.allowEmpty &&
                unique == other.unique &&
                size == other.size &&
                location == other.location &&
                version == other.version &&
                childRaw == other.childRaw &&
                childNullable == other.childNullable &&
                childAllowEmpty == other.childAllowEmpty &&
                childUnique == other.childUnique &&
                childSize == other.childSize &&
                childLocation == other.childLocation &&
                childVersion == other.childVersion &&
                type == other.type &&
                desc == other.desc &&
                default == other.default &&
                min == other.min &&
                max == other.max &&
                enum == other.enum &&
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
                name,
                isNull,
                example,
                desc,
                type,
                raw,
                nullable,
                allowEmpty,
                unique,
                default,
                min,
                max,
                size,
                enum,
                pattern,
                location,
                version,
                mime,
                childDesc,
                childType,
                childRaw,
                childNullable,
                childAllowEmpty,
                childUnique,
                childDefault,
                childMin,
                childMax,
                childSize,
                childEnum,
                childPattern,
                childLocation,
                childVersion,
                childMime,
                isInherit
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
            val b = s.toByteArray(charset("US-ASCII"))
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

        const val T_IS_NULL = "is_null"
        const val T_EXAMPLE = "example"

        const val T_NAME = "name"
        const val T_DESC = "desc"
        const val T_TYPE = "type"
        const val T_RAW = "raw"
        const val T_NULLABLE = "nullable"
        const val T_ALLOW_EMPTY = "allow_empty"
        const val T_UNIQUE = "unique"
        const val T_DEFAULT = "default"
        const val T_MIN = "min"
        const val T_MAX = "max"
        const val T_SIZE = "size"
        const val T_ENUM = "enum"
        const val T_PATTERN = "pattern"
        const val T_LOCATION = "location"
        const val T_VERSION = "version"
        const val T_MIME = "mime"

        const val T_CHILD_DESC = "child_desc"
        const val T_CHILD_TYPE = "child_type"
        const val T_CHILD_RAW = "child_raw"
        const val T_CHILD_NULLABLE = "child_nullable"
        const val T_CHILD_ALLOW_EMPTY = "child_allow_empty"
        const val T_CHILD_UNIQUE = "child_unique"
        const val T_CHILD_DEFAULT = "child_default"
        const val T_CHILD_MIN = "child_min"
        const val T_CHILD_MAX = "child_max"
        const val T_CHILD_SIZE = "child_size"
        const val T_CHILD_ENUM = "child_enum"
        const val T_CHILD_PATTERN = "child_pattern"
        const val T_CHILD_LOCATION = "child_location"
        const val T_CHILD_VERSION = "child_version"
        const val T_CHILD_MIME = "child_mime"

        private val emailRegex =
                Pattern.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
        private val decimalRegex = Pattern.compile("^-?\\d+\\.\\d+$")
        private val uuidRegex =
                Pattern.compile(
                        "^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$"
                )

        fun empty(): Tag = Tag()

        fun fromAnnotation(ann: MM?): Tag {
            val t = Tag()
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
            t.default = ann.default
            t.min = ann.min
            t.max = ann.max
            t.size = ann.size
            t.enum = ann.enum
            if (t.enum.isNotEmpty()) {
                t.type = ValueType.ENUM
            }
            t.pattern = ann.pattern
            t.location = ann.location
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
            t.childLocation = ann.childLocation
            t.childVersion = ann.childVersion
            t.childMime = ann.childMime
            return t
        }

        fun mergeTag(dst: Tag, src: Tag?): Tag {
            if (src == null) return dst

            if (src.isNull) dst.isNull = src.isNull
            if (src.example) dst.example = src.example
            if (src.desc.isNotEmpty()) dst.desc = src.desc
            if (src.type != ValueType.UNKNOWN) dst.type = src.type
            if (src.raw) dst.raw = true
            if (src.nullable) dst.nullable = true
            if (src.allowEmpty) dst.allowEmpty = true
            if (src.unique) dst.unique = true
            if (src.default.isNotEmpty()) dst.default = src.default
            if (src.min.isNotEmpty()) dst.min = src.min
            if (src.max.isNotEmpty()) dst.max = src.max
            if (src.size != 0) dst.size = src.size
            if (src.enum.isNotEmpty()) dst.enum = src.enum
            if (src.pattern.isNotEmpty()) dst.pattern = src.pattern
            if (src.location != 0) dst.location = src.location
            if (src.version != DEFAULT_VERSION) dst.version = src.version
            if (src.mime.isNotEmpty()) dst.mime = src.mime

            if (src.childDesc.isNotEmpty()) dst.childDesc = src.childDesc
            if (src.childType != ValueType.UNKNOWN) dst.childType = src.childType
            if (src.childRaw) dst.childRaw = true
            if (src.childNullable) dst.childNullable = true
            if (src.childAllowEmpty) dst.childAllowEmpty = true
            if (src.childUnique) dst.childUnique = true
            if (src.childDefault.isNotEmpty()) dst.childDefault = src.childDefault
            if (src.childMin.isNotEmpty()) dst.childMin = src.childMin
            if (src.childMax.isNotEmpty()) dst.childMax = src.childMax
            if (src.childSize != 0) dst.childSize = src.childSize
            if (src.childEnum.isNotEmpty()) dst.childEnum = src.childEnum
            if (src.childPattern.isNotEmpty()) dst.childPattern = src.childPattern
            if (src.childLocation != 0) dst.childLocation = src.childLocation
            if (src.childVersion != DEFAULT_VERSION) dst.childVersion = src.childVersion
            if (src.childMime.isNotEmpty()) dst.childMime = src.childMime

            return dst
        }

        fun parseMMTag(tag: String): Tag {
            val r = Tag()
            var t = tag.trim()
            t = t.removePrefix("//").trim()
            t = t.removePrefix("mm:").trim()
            if (t.isEmpty()) return r

            val parts = splitTag(t)
            for (p in parts) {
                if (p.isEmpty()) continue

                val k: String
                val v: String
                if (p.contains("=")) {
                    val kv = p.split("=", limit = 2)
                    k = kv[0].trim()
                    v = kv[1].trim()
                } else {
                    k = p.trim()
                    v = ""
                }

                var value = v
                if (value.length >= 2 && value.startsWith("\"") && value.endsWith("\"")) {
                    value = unquote(value)
                }

                val lower = k.lowercase()
                when (lower) {
                    T_NAME -> r.name = value
                    T_IS_NULL -> {
                        r.isNull = true
                        r.nullable = true
                    }
                    T_EXAMPLE -> r.example = true
                    T_DESC -> r.desc = value
                    T_TYPE -> r.type = ValueType.parseWireName(value)
                    T_RAW -> r.raw = true
                    T_NULLABLE -> r.nullable = true
                    T_ALLOW_EMPTY -> r.allowEmpty = true
                    T_UNIQUE -> r.unique = true
                    T_DEFAULT -> r.default = value
                    T_MIN -> r.min = value
                    T_MAX -> r.max = value
                    T_SIZE -> {
                        val u = value.toLongOrNull()
                        if (u != null && u <= Int.MAX_VALUE.toLong()) {
                            r.size = u.toInt()
                        }
                    }
                    T_ENUM -> {
                        r.type = ValueType.ENUM
                        r.enum = value
                    }
                    T_PATTERN -> r.pattern = value
                    T_LOCATION -> {
                        val d = value.toIntOrNull()
                        if (d != null && d >= -12 && d <= 14) {
                            r.location = d
                        }
                    }
                    T_VERSION -> {
                        val d = value.toIntOrNull()
                        if (d != null && d in 1..10) {
                            r.version = d
                        }
                    }
                    T_MIME -> r.mime = value
                    T_CHILD_DESC -> r.childDesc = value
                    T_CHILD_TYPE -> r.childType = ValueType.parseWireName(value)
                    T_CHILD_RAW -> r.childRaw = true
                    T_CHILD_NULLABLE -> r.childNullable = true
                    T_CHILD_ALLOW_EMPTY -> r.childAllowEmpty = true
                    T_CHILD_UNIQUE -> r.childUnique = true
                    T_CHILD_DEFAULT -> r.childDefault = value
                    T_CHILD_MIN -> r.childMin = value
                    T_CHILD_MAX -> r.childMax = value
                    T_CHILD_SIZE -> {
                        val u = value.toLongOrNull()
                        if (u != null && u <= Int.MAX_VALUE.toLong()) {
                            r.childSize = u.toInt()
                        }
                    }
                    T_CHILD_ENUM -> {
                        r.childType = ValueType.ENUM
                        r.childEnum = value
                    }
                    T_CHILD_PATTERN -> r.childPattern = value
                    T_CHILD_LOCATION -> {
                        val d = value.toIntOrNull()
                        if (d != null && d >= -12 && d <= 14) {
                            r.childLocation = d
                        }
                    }
                    T_CHILD_VERSION -> {
                        val d = value.toIntOrNull()
                        if (d != null && d in 1..10) {
                            r.childVersion = d
                        }
                    }
                    T_CHILD_MIME -> r.childMime = value
                }
            }
            return r
        }

        private fun splitTag(tag: String): List<String> {
            if (tag.isEmpty()) return emptyList()
            return tag.split(";").map { it.trim() }
        }

        internal fun quote(s: String): String {
            val sb = StringBuilder()
            sb.append('"')
            for (c in s) {
                when (c) {
                    '\\' -> sb.append("\\\\")
                    '"' -> sb.append("\\\"")
                    '\n' -> sb.append("\\n")
                    '\r' -> sb.append("\\r")
                    '\t' -> sb.append("\\t")
                    '\b' -> sb.append("\\b")
                    in '\u0000'..'\u0007', '\u000B', '\u000C', in '\u000E'..'\u001F' -> {
                        sb.append(String.format("\\u%04x", c.code))
                    }
                    else -> sb.append(c)
                }
            }
            sb.append('"')
            return sb.toString()
        }

        private fun unquote(s: String): String {
            val sb = StringBuilder()
            var i = 1
            while (i < s.length - 1) {
                val c = s[i]
                if (c == '\\' && i + 1 < s.length - 1) {
                    when (s[i + 1]) {
                        '\\' -> {
                            sb.append('\\')
                            i++
                        }
                        '"' -> {
                            sb.append('"')
                            i++
                        }
                        'n' -> {
                            sb.append('\n')
                            i++
                        }
                        'r' -> {
                            sb.append('\r')
                            i++
                        }
                        't' -> {
                            sb.append('\t')
                            i++
                        }
                        'u' -> {
                            if (i + 5 < s.length) {
                                val hex = s.substring(i + 2, i + 6)
                                try {
                                    sb.append(hex.toInt(16).toChar())
                                    i += 5
                                } catch (_: NumberFormatException) {
                                    sb.append(c)
                                }
                            } else {
                                sb.append(c)
                            }
                        }
                        else -> sb.append(c)
                    }
                } else {
                    sb.append(c)
                }
                i++
            }
            return sb.toString()
        }
    }

    data class ValidationResult(
            val valid: Boolean,
            val error: String? = null,
            val data: Any? = null,
            val text: String? = null
    )

    fun validateArray(value: List<*>): ValidationResult {
        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type array not support location UTC$location")
        }

        val length = value.size

        if (length == 0) {
            if (!allowEmpty) {
                return ValidationResult(false, "type array not allow empty")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (size > 0 && length > size) {
            return ValidationResult(false, "type array over size")
        }

        if (childUnique) {
            val seen = mutableSetOf<Any?>()
            for (i in value.indices) {
                val data = value[i]
                if (seen.contains(data)) {
                    return ValidationResult(false, "array duplicate value found: $data, index: $i")
                }
                seen.add(data)
            }
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateStruct(): ValidationResult {
        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type struct not support location UTC$location")
        }

        return ValidationResult(true)
    }

    fun validateMap(): ValidationResult {
        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type map not support location UTC$location")
        }

        return ValidationResult(true)
    }

    fun validateString(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = value)
            }
            return ValidationResult(false, "type string not allow empty value \"$value\"")
        }

        if (pattern.isNotEmpty()) {
            try {
                val regex = Pattern.compile(pattern)
                if (!regex.matcher(value).matches()) {
                    return ValidationResult(
                            false,
                            "value \"$value\" does not match pattern $pattern"
                    )
                }
            } catch (e: Exception) {
                return ValidationResult(false, "pattern \"$pattern\" compile err: ${e.message}")
            }
        }

        val length = value.codePointCount(0, value.length)

        if (min.isNotEmpty()) {
            val mini = min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int: $min")
            }
            if (length < mini) {
                return ValidationResult(
                        false,
                        "string length $length is less than the minimum limit $mini"
                )
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int: $max")
            }
            if (length > maxi) {
                return ValidationResult(
                        false,
                        "string length $length exceeds the maximum limit $maxi"
                )
            }
        }

        if (size != 0 && length != size) {
            return ValidationResult(false, "string length $length != size $size")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type string not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value)
    }

    fun validateBytes(value: ByteArray): ValidationResult {
        val length = value.size

        if (length == 0) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "")
            }
            return ValidationResult(false, "type []byte not allow empty value []byte{}")
        }

        if (min.isNotEmpty()) {
            val mini = min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int: $min")
            }
            if (length < mini) {
                return ValidationResult(
                        false,
                        "[]byte length $length is less than the minimum limit $mini"
                )
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int: $max")
            }
            if (length > maxi) {
                return ValidationResult(
                        false,
                        "[]byte length $length exceeds the maximum limit $maxi"
                )
            }
        }

        if (size != 0 && length != size) {
            return ValidationResult(false, "[]byte length $length != size $size")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type []byte not support location UTC$location")
        }

        val text = Base64.getEncoder().encodeToString(value)

        return ValidationResult(true, data = value, text = text)
    }

    fun validateBool(value: Boolean): ValidationResult {
        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (allowEmpty) {
            return ValidationResult(false, "type bool not support allow empty")
        }

        if (location != 0) {
            return ValidationResult(false, "type bool not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateSlice(value: List<*>): ValidationResult {
        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type slice not support location UTC$location")
        }

        val length = value.size

        if (length == 0) {
            if (!allowEmpty) {
                return ValidationResult(false, "type slice not allow empty")
            }
            return ValidationResult(true)
        }

        if (childUnique) {
            val seen = mutableSetOf<Any?>()
            for (i in value.indices) {
                val data = value[i]
                if (seen.contains(data)) {
                    return ValidationResult(false, "slice duplicate value found: $data, index: $i")
                }
                seen.add(data)
            }
        }

        return ValidationResult(true)
    }

    fun validateInt(value: Int): ValidationResult {
        if (value == 0) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type int not allow empty value $value")
        }

        val val64 = value.toLong()

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type int not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateInt8(value: Byte): ValidationResult {
        if (value == 0.toByte()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type int8 not allow empty value $value")
        }

        val val64 = value.toLong()

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int8: $min")
            }
            if (mini < Byte.MIN_VALUE.toLong() || mini > Byte.MAX_VALUE.toLong()) {
                return ValidationResult(
                        false,
                        "failed to parse t.Min as int8: value $mini out of int8 range"
                )
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int8: $max")
            }
            if (maxi < Byte.MIN_VALUE.toLong() || maxi > Byte.MAX_VALUE.toLong()) {
                return ValidationResult(
                        false,
                        "failed to parse t.Max as int8: value $maxi out of int8 range"
                )
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type int8 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateInt16(value: Short): ValidationResult {
        if (value == 0.toShort()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type int16 not allow empty value $value")
        }

        val val64 = value.toLong()

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int16: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int16: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type int16 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateInt32(value: Int): ValidationResult {
        if (value == 0) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type int32 not allow empty value $value")
        }

        val val64 = value.toLong()

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int32: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int32: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type int32 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateInt64(value: Long): ValidationResult {
        if (value == 0L) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type int64 not allow empty value $value")
        }

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int64: $min")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int64: $max")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type int64 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateUint(value: Long): ValidationResult {
        if (value == 0L) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type uint not allow empty value $value")
        }

        val val64 = value

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as uint: $min")
            }
            if (mini < 0) {
                return ValidationResult(false, "failed to parse t.Min as uint: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as uint: $max")
            }
            if (maxi < 0) {
                return ValidationResult(false, "failed to parse t.Max as uint: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type uint not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateUint8(value: Short): ValidationResult {
        if (value == 0.toShort()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type uint8 not allow empty value $value")
        }

        val val64 = value.toLong()

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as uint8: $min")
            }
            if (mini < 0) {
                return ValidationResult(false, "failed to parse t.Min as uint8: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as uint8: $max")
            }
            if (maxi < 0) {
                return ValidationResult(false, "failed to parse t.Max as uint8: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type uint8 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateUint16(value: Int): ValidationResult {
        if (value == 0) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type uint16 not allow empty value $value")
        }

        val val64 = value.toLong()

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as uint16: $min")
            }
            if (mini < 0) {
                return ValidationResult(false, "failed to parse t.Min as uint16: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as uint16: $max")
            }
            if (maxi < 0) {
                return ValidationResult(false, "failed to parse t.Max as uint16: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type uint16 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateUint32(value: Long): ValidationResult {
        if (value == 0L) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type uint32 not allow empty value $value")
        }

        val val64 = value

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as uint32: $min")
            }
            if (mini < 0) {
                return ValidationResult(false, "failed to parse t.Min as uint32: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "value $val64 is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as uint32: $max")
            }
            if (maxi < 0) {
                return ValidationResult(false, "failed to parse t.Max as uint32: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "value $val64 exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type uint32 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateUint64(value: BigInteger): ValidationResult {
        val zero = BigInteger.ZERO

        if (value == zero) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type uint64 not allow empty value $value")
        }

        if (min.isNotEmpty()) {
            val mini: BigInteger
            try {
                mini = BigInteger(min)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse t.Min as uint64: $min")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi: BigInteger
            try {
                maxi = BigInteger(max)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse t.Max as uint64: $max")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type uint64 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateFloat32(value: Float): ValidationResult {
        if (value == 0.0f) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0.0")
            }
            return ValidationResult(false, "type float32 not allow empty value 0.0")
        }

        val val64 = value.toDouble()

        if (min.isNotEmpty()) {
            val mini = min.toDoubleOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as float32: $min")
            }
            if (val64 < mini) {
                return ValidationResult(false, "$val64 < min $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toDoubleOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as float32: $max")
            }
            if (val64 > maxi) {
                return ValidationResult(false, "$val64 > max $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type float32 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateFloat64(value: Double): ValidationResult {
        if (value == 0.0) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0.0")
            }
            return ValidationResult(false, "type float64 not allow empty value 0.0")
        }

        if (min.isNotEmpty()) {
            val mini = min.toDoubleOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as float64: $min")
            }
            if (value < mini) {
                return ValidationResult(false, "$value < min $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toDoubleOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as float64: $max")
            }
            if (value > maxi) {
                return ValidationResult(false, "$value > max $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type float64 not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateBigInt(value: BigInteger): ValidationResult {
        val zero = BigInteger.ZERO

        if (value == zero) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "0")
            }
            return ValidationResult(false, "type big.Int not allow empty value 0")
        }

        if (min.isNotEmpty()) {
            val mini: BigInteger
            try {
                mini = BigInteger(min)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "invalid min \"$min\" for big.Int")
            }
            if (value < mini) {
                return ValidationResult(false, "big.Int length ${value} < min ${mini}")
            }
        }

        if (max.isNotEmpty()) {
            val maxi: BigInteger
            try {
                maxi = BigInteger(max)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "invalid max \"$max\" for big.Int")
            }
            if (value > maxi) {
                return ValidationResult(false, "big.Int length ${value} > max ${maxi}")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type big.Int not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateDateTime(value: LocalDateTime): ValidationResult {
        val epoch = LocalDateTime.of(1970, 1, 1, 0, 0, 0)

        val format = value.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"))
        if (value == epoch) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = format)
            }
            return ValidationResult(
                    false,
                    "datetime type does not allow empty \"$format\". you can set allow_empty or child_allow_empty to allow it."
            )
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        return ValidationResult(true, data = value, text = format)
    }

    fun validateDate(value: LocalDate): ValidationResult {
        val epoch = LocalDate.of(1970, 1, 1)

        val format = value.toString()
        if (value == epoch) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = format)
            }
            return ValidationResult(
                    false,
                    "date type does not allow empty \"$format\". you can set allow_empty or child_allow_empty to allow it."
            )
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        return ValidationResult(true, data = value, text = format)
    }

    fun validateTime(value: LocalTime): ValidationResult {
        val zero = LocalTime.of(0, 0, 0)

        val format = value.toString()
        if (value == zero) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = format)
            }
            return ValidationResult(
                    false,
                    "time type does not allow empty \"$format\". you can set allow_empty or child_allow_empty to allow it."
            )
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        return ValidationResult(true, data = value, text = format)
    }

    fun validateUUID(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = ByteArray(16), text = value)
            }
            return ValidationResult(false, "type uuid not allow empty value \"\"")
        }

        if (!uuidRegex.matcher(value).matches()) {
            return ValidationResult(false, "value '$value' does not match UUID pattern")
        }

        val uuidBytes = ByteArray(16)
        val cleanValue = value.replace("-", "")
        for (i in 0 until 16) {
            val start = i * 2
            val end = start + 2
            if (end <= cleanValue.length) {
                val part = cleanValue.substring(start, end)
                uuidBytes[i] = part.toInt(16).toByte()
            }
        }

        if (version != DEFAULT_VERSION) {
            val uuidVersion = ((uuidBytes[6].toInt() shr 4) and 0x0F)
            if (version != uuidVersion) {
                return ValidationResult(false, "invalid uuid version")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type uuid not support location UTC$location")
        }

        return ValidationResult(true, data = uuidBytes, text = value)
    }

    fun validateDecimal(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = value)
            }
            return ValidationResult(false, "type decimal not allow empty value \"\"")
        }

        if (!decimalRegex.matcher(value).matches()) {
            return ValidationResult(false, "invalid decimal \"$value\", must be like \"0.0\"")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type decimal not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value)
    }

    fun validateIP(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = ByteArray(0), text = "")
            }
            return ValidationResult(false, "type ip not allow empty value \"\"")
        }

        try {
            val addr = java.net.InetAddress.getByName(value)
            val ipBytes = addr.address

            if (version == 4 && ipBytes.size != 4) {
                return ValidationResult(false, "invalid ipv4: $value")
            }

            if (version == 6 && ipBytes.size != 16) {
                return ValidationResult(false, "invalid ipv6: $value")
            }
        } catch (e: Exception) {
            return ValidationResult(false, "value '$value' does not match IP pattern")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type ip not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value)
    }

    fun validateURL(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "")
            }
            return ValidationResult(false, "type url not allow empty value \"\"")
        }

        try {
            val url = java.net.URI(value).toURL()
            if (url.protocol != "http" && url.protocol != "https") {
                return ValidationResult(false, "invalid url: $value")
            }
            if (url.host.isEmpty()) {
                return ValidationResult(false, "invalid url: $value")
            }
        } catch (e: Exception) {
            return ValidationResult(false, "invalid url: $value")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type url not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value)
    }

    fun validateEmail(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = value)
            }
            return ValidationResult(false, "type email not allow empty value \"\"")
        }

        if (!emailRegex.matcher(value).matches()) {
            return ValidationResult(false, "value '$value' does not match email pattern")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type email not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value)
    }

    fun validateEnum(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (allowEmpty) {
                return ValidationResult(true, data = -1, text = value)
            }
            return ValidationResult(false, "type enum not allow empty value \"\"")
        }

        val enums = enum.split("|")
        var idx = -1
        for (i in enums.indices) {
            if (enums[i].trim() == value) {
                idx = i
                break
            }
        }

        if (idx == -1) {
            return ValidationResult(false, "value '$value' not found in enum: $enums")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type enum not support location UTC$location")
        }

        return ValidationResult(true, data = idx, text = value)
    }

    fun validateImage(value: ByteArray): ValidationResult {
        val length = value.size

        if (length == 0) {
            if (allowEmpty) {
                return ValidationResult(true, data = value, text = "")
            }
            return ValidationResult(false, "type image not allow empty value []byte{}")
        }

        if (min.isNotEmpty()) {
            val mini = min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse t.Min as int: $min")
            }
            if (length < mini) {
                return ValidationResult(false, "[]byte length $length < min $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse t.Max as int: $max")
            }
            if (length > maxi) {
                return ValidationResult(false, "[]byte length $length > max $maxi")
            }
        }

        if (size != 0 && length != size) {
            return ValidationResult(false, "[]byte length $length != size $size")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type image not support location UTC$location")
        }

        val text = Base64.getEncoder().encodeToString(value)

        return ValidationResult(true, data = value, text = text)
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
