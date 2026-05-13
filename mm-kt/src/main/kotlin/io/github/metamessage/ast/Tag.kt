package io.github.metamessage.ast


import io.github.metamessage.mm.WireConstants
import io.github.metamessage.MM

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
    fun copy(): Tag = Tag(
        name, isNull, example, desc, type, raw, nullable, allowEmpty, unique,
        default, min, max, size, enum, pattern, location, version, mime,
        childDesc, childType, childRaw, childNullable, childAllowEmpty, childUnique,
        childDefault, childMin, childMax, childSize, childEnum, childPattern,
        childLocation, childVersion, childMime, isInherit
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

    override fun toString(): String {
        val parts = mutableListOf<String>()

        if (type != ValueType.UNKNOWN && !isInherit) {
            if (!(type == ValueType.STRING ||
                    type == ValueType.INT ||
                    type == ValueType.FLOAT64 ||
                    type == ValueType.BOOL ||
                    type == ValueType.STRUCT ||
                    type == ValueType.SLICE)) {
                if (!((type == ValueType.ARRAY && size > 0) ||
                        (type == ValueType.ENUM && enum.isNotEmpty()))) {
                    parts.add("type=${type.toString()}")
                }
            }
        }

        if (example) {
            parts.add("example")
        }

        if (isNull) {
            parts.add("is_null")
        }

        if (nullable && !isInherit) {
            if (!isNull) {
                parts.add("nullable")
            }
        }

        if (desc.isNotEmpty() && !isInherit) {
            parts.add("desc=${desc}")
        }

        if (raw && !isInherit) {
            parts.add("raw")
        }

        if (allowEmpty && !isInherit) {
            parts.add("allow_empty")
        }

        if (unique && !isInherit) {
            parts.add("unique")
        }

        if (default.isNotEmpty() && !isInherit) {
            parts.add("default=${default}")
        }

        if (min.isNotEmpty() && !isInherit) {
            parts.add("min=${min}")
        }

        if (max.isNotEmpty() && !isInherit) {
            parts.add("max=${max}")
        }

        if (size != 0 && !isInherit) {
            parts.add("size=${size}")
        }

        if (enum.isNotEmpty() && !isInherit) {
            parts.add("enum=${enum}")
        }

        if (pattern.isNotEmpty() && !isInherit) {
            parts.add("pattern=${pattern}")
        }

        if (location != 0 && !isInherit) {
            parts.add("location=${location}")
        }

        if (version != 0 && !isInherit) {
            parts.add("version=${version}")
        }

        if (mime.isNotEmpty() && !isInherit) {
            parts.add("mime=${mime}")
        }

        if (childDesc.isNotEmpty()) {
            parts.add("child_desc=${childDesc}")
        }

        if (childType != ValueType.UNKNOWN) {
            if (!(childType == ValueType.STRING ||
                    childType == ValueType.INT ||
                    childType == ValueType.FLOAT64 ||
                    childType == ValueType.BOOL ||
                    childType == ValueType.STRUCT ||
                    childType == ValueType.SLICE)) {
                if (!((childType == ValueType.ARRAY && childSize > 0) ||
                        (childType == ValueType.ENUM && childEnum.isNotEmpty()))) {
                    parts.add("child_type=${childType.toString()}")
                }
            }
        }

        if (childRaw) {
            parts.add("child_raw")
        }

        if (childNullable) {
            parts.add("child_nullable")
        }

        if (childAllowEmpty) {
            parts.add("child_allow_empty")
        }

        if (childUnique) {
            parts.add("child_unique")
        }

        if (childDefault.isNotEmpty()) {
            parts.add("child_default=${childDefault}")
        }

        if (childMin.isNotEmpty()) {
            parts.add("child_min=${childMin}")
        }

        if (childMax.isNotEmpty()) {
            parts.add("child_max=${childMax}")
        }

        if (childSize != 0) {
            parts.add("child_size=${childSize}")
        }

        if (childEnum.isNotEmpty()) {
            parts.add("child_enum=${childEnum}")
        }

        if (childPattern.isNotEmpty()) {
            parts.add("child_pattern=${childPattern}")
        }

        if (childLocation != 0) {
            parts.add("child_location=${childLocation}")
        }

        if (childVersion != 0) {
            parts.add("child_version=${childVersion}")
        }

        if (childMime.isNotEmpty()) {
            parts.add("child_mime=${childMime}")
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
        if (version != DEFAULT_VERSION && !isInherit) encodeUint64(w, TagKey.K_VERSION, version.toLong())
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
        if (childVersion != DEFAULT_VERSION) encodeUint64(w, TagKey.K_CHILD_VERSION, childVersion.toLong())
        if (childMime.isNotEmpty()) {
            val m = Mime.parse(childMime)
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
            ValueType.ENUM -> if (enum.isNotEmpty()) false else true
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
            name, isNull, example, desc, type, raw, nullable, allowEmpty, unique,
            default, min, max, size, enum, pattern, location, version, mime,
            childDesc, childType, childRaw, childNullable, childAllowEmpty, childUnique,
            childDefault, childMin, childMax, childSize, childEnum, childPattern,
            childLocation, childVersion, childMime, isInherit
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

        private val emailRegex = Pattern.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
        private val decimalRegex = Pattern.compile("^-?\\d+\\.\\d+$")
        private val uuidRegex = Pattern.compile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")
        private val ipv4Regex = Pattern.compile("^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$")
        private val ipv6Regex = Pattern.compile("^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$")
        private val urlRegex = Pattern.compile("^https?://[\\w.-]+(?:/[\\w./?%&=-]*)?$")
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
            val seen = mutableSetOf<Any>()
            for (i in value.indices) {
                val data = value[i]
                val key = data ?: "null"
                if (seen.contains(key)) {
                    return ValidationResult(false, "array duplicate value found: $data, index: $i")
                }
                seen.add(key)
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

    fun validateString(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (!allowEmpty) {
                return ValidationResult(false, "type string not allow empty value \"$value\"")
            }
            return ValidationResult(true, data = value, text = value)
        }

        if (pattern.isNotEmpty()) {
            try {
                val regex = Pattern.compile(pattern)
                if (!regex.matcher(value).matches()) {
                    return ValidationResult(false, "value \"$value\" does not match pattern $pattern")
                }
            } catch (e: Exception) {
                return ValidationResult(false, "pattern \"$pattern\" compile err: ${e.message}")
            }
        }

        val length = value.length

        if (min.isNotEmpty()) {
            val mini = min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int: $min")
            }
            if (length < mini) {
                return ValidationResult(false, "string length $length is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int: $max")
            }
            if (length > maxi) {
                return ValidationResult(false, "string length $length exceeds the maximum limit $maxi")
            }
        }

        if (size > 0 && length != size) {
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
            val seen = mutableSetOf<Any>()
            for (i in value.indices) {
                val data = value[i]
                val key = data ?: "null"
                if (seen.contains(key)) {
                    return ValidationResult(false, "slice duplicate value found: $data, index: $i")
                }
                seen.add(key)
            }
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

    fun validateBytes(value: ByteArray): ValidationResult {
        val length = value.size

        if (length == 0) {
            if (!allowEmpty) {
                return ValidationResult(false, "type []byte not allow empty value []byte{}")
            }
            return ValidationResult(true, data = value, text = "")
        }

        if (min.isNotEmpty()) {
            val mini = min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int: $min")
            }
            if (length < mini) {
                return ValidationResult(false, "[]byte length $length is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int: $max")
            }
            if (length > maxi) {
                return ValidationResult(false, "[]byte length $length exceeds the maximum limit $maxi")
            }
        }

        if (size > 0 && length != size) {
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

    fun validateInt8(value: Byte): ValidationResult {
        if (value == 0.toByte()) {
            if (!allowEmpty) {
                return ValidationResult(false, "type int8 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = -128L
        val maxVal = 127L

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int8: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of int8 range [-128, 127]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int8: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of int8 range [-128, 127]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type int16 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = -32768L
        val maxVal = 32767L

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int16: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of int16 range [-32768, 32767]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int16: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of int16 range [-32768, 32767]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type int32 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = -2147483648L
        val maxVal = 2147483647L

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int32: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of int32 range [-2147483648, 2147483647]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int32: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of int32 range [-2147483648, 2147483647]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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

    fun validateInt(value: Int): ValidationResult {
        return validateInt32(value)
    }

    fun validateInt64(value: Long): ValidationResult {
        if (value == 0L) {
            if (!allowEmpty) {
                return ValidationResult(false, "type int64 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = Long.MIN_VALUE
        val maxVal = Long.MAX_VALUE

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int64: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of int64 range [$minVal, $maxVal]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int64: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of int64 range [$minVal, $maxVal]")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type uint not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = 0L
        val maxVal = 4294967295L

        if (value < minVal || value > maxVal) {
            return ValidationResult(false, "value $value is out of uint range [0, 4294967295]")
        }

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of uint range [0, 4294967295]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of uint range [0, 4294967295]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type uint8 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = 0L
        val maxVal = 255L

        if (value < minVal || value > maxVal) {
            return ValidationResult(false, "value $value is out of uint8 range [0, 255]")
        }

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint8: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of uint8 range [0, 255]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint8: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of uint8 range [0, 255]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type uint16 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = 0L
        val maxVal = 65535L

        if (value < minVal || value > maxVal) {
            return ValidationResult(false, "value $value is out of uint16 range [0, 65535]")
        }

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint16: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of uint16 range [0, 65535]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint16: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of uint16 range [0, 65535]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type uint32 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = 0L
        val maxVal = 4294967295L

        if (value < minVal || value > maxVal) {
            return ValidationResult(false, "value $value is out of uint32 range [0, 4294967295]")
        }

        if (min.isNotEmpty()) {
            val mini = min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint32: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of uint32 range [0, 4294967295]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint32: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of uint32 range [0, 4294967295]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type uint64 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val minVal = BigInteger.ZERO
        val maxVal = BigInteger("18446744073709551615")

        if (value < minVal || value > maxVal) {
            return ValidationResult(false, "value $value is out of uint64 range [0, 18446744073709551615]")
        }

        if (min.isNotEmpty()) {
            val mini: BigInteger
            try {
                mini = BigInteger(min)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse tag.min as uint64: $min")
            }
            if (mini < minVal || mini > maxVal) {
                return ValidationResult(false, "tag.min $mini is out of uint64 range [0, 18446744073709551615]")
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
                return ValidationResult(false, "failed to parse tag.max as uint64: $max")
            }
            if (maxi < minVal || maxi > maxVal) {
                return ValidationResult(false, "tag.max $maxi is out of uint64 range [0, 18446744073709551615]")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type float32 not allow empty value 0.0")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (min.isNotEmpty()) {
            val mini = min.toFloatOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as float32: $min")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toFloatOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as float32: $max")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type float64 not allow empty value 0.0")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (min.isNotEmpty()) {
            val mini = min.toDoubleOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as float64: $min")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toDoubleOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as float64: $max")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type bigint not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (min.isNotEmpty()) {
            val mini: BigInteger
            try {
                mini = BigInteger(min)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse tag.min as bigint: $min")
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
                return ValidationResult(false, "failed to parse tag.max as bigint: $max")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type bigint not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    fun validateDateTime(value: LocalDateTime): ValidationResult {
        if (value == LocalDateTime.MIN) {
            if (!allowEmpty) {
                return ValidationResult(false, "type datetime not allow empty ${value.toString()}")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        val text = value.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"))

        return ValidationResult(true, data = value, text = text)
    }

    fun validateDate(value: LocalDate): ValidationResult {
        if (value == LocalDate.MIN) {
            if (!allowEmpty) {
                return ValidationResult(false, "type date not allow empty ${value.toString()}")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        val text = value.toString()

        return ValidationResult(true, data = value, text = text)
    }

    fun validateTime(value: LocalTime): ValidationResult {
        if (value == LocalTime.MIN) {
            if (!allowEmpty) {
                return ValidationResult(false, "type time not allow empty ${value.toString()}")
            }
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        val text = value.toString()

        return ValidationResult(true, data = value, text = text)
    }

    fun validateUUID(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (!allowEmpty) {
                return ValidationResult(false, "type uuid not allow empty value \"\"")
            }
            return ValidationResult(true, data = ByteArray(16), text = value)
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
            val uuidVersion = cleanValue.substring(12, 13).toInt(16)
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

    fun validateEmail(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (!allowEmpty) {
                return ValidationResult(false, "type email not allow empty value \"\"")
            }
            return ValidationResult(true, data = value, text = value)
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
            if (!allowEmpty) {
                return ValidationResult(false, "type enum not allow empty value \"\"")
            }
            return ValidationResult(true, data = -1, text = value)
        }

        if (enum.isEmpty()) {
            return ValidationResult(false, "enum not defined")
        }

        val enums = enum.split("|")
        var idx = -1
        for (i in enums.indices) {
            val enumValue = enums[i]
            if (enumValue.trim() == value) {
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
            if (!allowEmpty) {
                return ValidationResult(false, "type image not allow empty value []byte{}")
            }
            return ValidationResult(true, data = value, text = "")
        }

        if (min.isNotEmpty()) {
            val mini = min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int: $min")
            }
            if (length < mini) {
                return ValidationResult(false, "[]byte length $length < min $mini")
            }
        }

        if (max.isNotEmpty()) {
            val maxi = max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int: $max")
            }
            if (length > maxi) {
                return ValidationResult(false, "[]byte length $length > max $maxi")
            }
        }

        if (size > 0 && length != size) {
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

    fun validateDecimal(value: String): ValidationResult {
        if (value.isEmpty()) {
            if (!allowEmpty) {
                return ValidationResult(false, "type decimal not allow empty value \"\"")
            }
            return ValidationResult(true, data = value, text = value)
        }

        if (!decimalRegex.matcher(value).matches()) {
            return ValidationResult(false, "value '$value' does not match decimal pattern")
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
            if (!allowEmpty) {
                return ValidationResult(false, "type ip not allow empty value \"\"")
            }
            return ValidationResult(true, data = ByteArray(0), text = value)
        }

        val isIpv4 = ipv4Regex.matcher(value).matches()
        val isIpv6 = ipv6Regex.matcher(value).matches()

        if (!isIpv4 && !isIpv6) {
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
            if (!allowEmpty) {
                return ValidationResult(false, "type url not allow empty value \"\"")
            }
            return ValidationResult(true, data = value, text = value)
        }

        if (!urlRegex.matcher(value).matches()) {
            return ValidationResult(false, "value '$value' does not match URL pattern")
        }

        if (desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (location != 0) {
            return ValidationResult(false, "type url not support location UTC$location")
        }

        return ValidationResult(true, data = value, text = value)
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