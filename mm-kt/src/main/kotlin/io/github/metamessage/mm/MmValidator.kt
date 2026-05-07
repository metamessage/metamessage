package io.github.metamessage.mm

import java.math.BigInteger
import java.util.Base64
import java.util.regex.Pattern

// Validation result class
class ValidationResult(
    val valid: Boolean,
    val error: String? = null,
    val data: Any? = null,
    val text: String? = null
)

// Validator class
object MmValidator {
    // Regex patterns
    private val emailRegex = Pattern.compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
    private val decimalRegex = Pattern.compile("^-?\\d+\\.\\d+$")
    private val uuidRegex = Pattern.compile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

    // Validate array
    fun validateArray(value: List<*>, tag: MmTag): ValidationResult {
        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type array not support location UTC${tag.locationHours}")
        }

        val length = value.size

        if (length == 0) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type array not allow empty")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (tag.size > 0 && length > tag.size) {
            return ValidationResult(false, "type array over size")
        }

        if (tag.childUnique) {
            val seen = mutableSetOf<Any>()
            for (i in value.indices) {
                val data = value[i]
                val key = if (data is Any) data else "null"
                if (seen.contains(key)) {
                    return ValidationResult(false, "array duplicate value found: $data, index: $i")
                }
                seen.add(key)
            }
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate struct
    fun validateStruct(tag: MmTag): ValidationResult {
        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type struct not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true)
    }

    // Validate string
    fun validateString(value: String, tag: MmTag): ValidationResult {
        if (value.isEmpty()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type string not allow empty value \"$value\"")
            }
            return ValidationResult(true, data = value, text = value)
        }

        if (tag.pattern.isNotEmpty()) {
            try {
                val regex = Pattern.compile(tag.pattern)
                if (!regex.matcher(value).matches()) {
                    return ValidationResult(false, "value \"$value\" does not match pattern ${tag.pattern}")
                }
            } catch (e: Exception) {
                return ValidationResult(false, "pattern \"${tag.pattern}\" compile err: ${e.message}")
            }
        }

        val length = value.length

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int: ${tag.min}")
            }
            if (length < mini) {
                return ValidationResult(false, "string length $length is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int: ${tag.max}")
            }
            if (length > maxi) {
                return ValidationResult(false, "string length $length exceeds the maximum limit $maxi")
            }
        }

        if (tag.size > 0 && length != tag.size) {
            return ValidationResult(false, "string length $length != size ${tag.size}")
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type string not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value)
    }

    // Validate bytes
    fun validateBytes(value: ByteArray, tag: MmTag): ValidationResult {
        val length = value.size

        if (length == 0) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type []byte not allow empty value []byte{}")
            }
            return ValidationResult(true, data = value, text = "")
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int: ${tag.min}")
            }
            if (length < mini) {
                return ValidationResult(false, "[]byte length $length is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int: ${tag.max}")
            }
            if (length > maxi) {
                return ValidationResult(false, "[]byte length $length exceeds the maximum limit $maxi")
            }
        }

        if (tag.size > 0 && length != tag.size) {
            return ValidationResult(false, "[]byte length $length != size ${tag.size}")
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type []byte not support location UTC${tag.locationHours}")
        }

        // Convert to base64 string
        val text = Base64.getEncoder().encodeToString(value)

        return ValidationResult(true, data = value, text = text)
    }

    // Validate bool
    fun validateBool(value: Boolean, tag: MmTag): ValidationResult {
        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.allowEmpty) {
            return ValidationResult(false, "type bool not support allow empty")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type bool not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate int8
    fun validateInt8(value: Byte, tag: MmTag): ValidationResult {
        if (value == 0.toByte()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type int8 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = -128L
        val max = 127L

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int8: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of int8 range [-128, 127]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int8: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of int8 range [-128, 127]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type int8 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate int16
    fun validateInt16(value: Short, tag: MmTag): ValidationResult {
        if (value == 0.toShort()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type int16 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = -32768L
        val max = 32767L

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int16: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of int16 range [-32768, 32767]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int16: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of int16 range [-32768, 32767]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type int16 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate int32
    fun validateInt32(value: Int, tag: MmTag): ValidationResult {
        if (value == 0) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type int32 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = -2147483648L
        val max = 2147483647L

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int32: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of int32 range [-2147483648, 2147483647]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int32: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of int32 range [-2147483648, 2147483647]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type int32 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate int64
    fun validateInt64(value: Long, tag: MmTag): ValidationResult {
        if (value == 0L) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type int64 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = Long.MIN_VALUE
        val max = Long.MAX_VALUE

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int64: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of int64 range [" + Long.MIN_VALUE + ", " + Long.MAX_VALUE + "]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int64: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of int64 range [" + Long.MIN_VALUE + ", " + Long.MAX_VALUE + "]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type int64 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate uint
    fun validateUint(value: Long, tag: MmTag): ValidationResult {
        if (value == 0L) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type uint not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = 0L
        val max = 4294967295L

        if (value < min || value > max) {
            return ValidationResult(false, "value $value is out of uint range [0, 4294967295]")
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of uint range [0, 4294967295]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of uint range [0, 4294967295]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type uint not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate uint8
    fun validateUint8(value: Short, tag: MmTag): ValidationResult {
        if (value == 0.toShort()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type uint8 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = 0L
        val max = 255L

        if (value < min || value > max) {
            return ValidationResult(false, "value $value is out of uint8 range [0, 255]")
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint8: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of uint8 range [0, 255]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint8: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of uint8 range [0, 255]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type uint8 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate uint16
    fun validateUint16(value: Int, tag: MmTag): ValidationResult {
        if (value == 0) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type uint16 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = 0L
        val max = 65535L

        if (value < min || value > max) {
            return ValidationResult(false, "value $value is out of uint16 range [0, 65535]")
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint16: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of uint16 range [0, 65535]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint16: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of uint16 range [0, 65535]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type uint16 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate uint32
    fun validateUint32(value: Long, tag: MmTag): ValidationResult {
        if (value == 0L) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type uint32 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = 0L
        val max = 4294967295L

        if (value < min || value > max) {
            return ValidationResult(false, "value $value is out of uint32 range [0, 4294967295]")
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toLongOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as uint32: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of uint32 range [0, 4294967295]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toLongOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as uint32: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of uint32 range [0, 4294967295]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type uint32 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate uint64
    fun validateUint64(value: BigInteger, tag: MmTag): ValidationResult {
        val zero = BigInteger.ZERO

        if (value == zero) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type uint64 not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        val min = BigInteger.ZERO
        val max = BigInteger("18446744073709551615")

        if (value < min || value > max) {
            return ValidationResult(false, "value $value is out of uint64 range [0, 18446744073709551615]")
        }

        if (tag.min.isNotEmpty()) {
            val mini: BigInteger
            try {
                mini = BigInteger(tag.min)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse tag.min as uint64: ${tag.min}")
            }
            if (mini < min || mini > max) {
                return ValidationResult(false, "tag.min $mini is out of uint64 range [0, 18446744073709551615]")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi: BigInteger
            try {
                maxi = BigInteger(tag.max)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse tag.max as uint64: ${tag.max}")
            }
            if (maxi < min || maxi > max) {
                return ValidationResult(false, "tag.max $maxi is out of uint64 range [0, 18446744073709551615]")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type uint64 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate float32
    fun validateFloat32(value: Float, tag: MmTag): ValidationResult {
        if (value == 0.0f) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type float32 not allow empty value 0.0")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toFloatOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as float32: ${tag.min}")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toFloatOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as float32: ${tag.max}")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type float32 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate float64
    fun validateFloat64(value: Double, tag: MmTag): ValidationResult {
        if (value == 0.0) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type float64 not allow empty value 0.0")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toDoubleOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as float64: ${tag.min}")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toDoubleOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as float64: ${tag.max}")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type float64 not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate BigInteger
    fun validateBigInteger(value: BigInteger, tag: MmTag): ValidationResult {
        val zero = BigInteger.ZERO

        if (value == zero) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type bigint not allow empty value $value")
            }
            return ValidationResult(true, data = value, text = value.toString())
        }

        if (tag.min.isNotEmpty()) {
            val mini: BigInteger
            try {
                mini = BigInteger(tag.min)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse tag.min as bigint: ${tag.min}")
            }
            if (value < mini) {
                return ValidationResult(false, "value $value is less than the minimum limit $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi: BigInteger
            try {
                maxi = BigInteger(tag.max)
            } catch (e: NumberFormatException) {
                return ValidationResult(false, "failed to parse tag.max as bigint: ${tag.max}")
            }
            if (value > maxi) {
                return ValidationResult(false, "value $value exceeds the maximum limit $maxi")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type bigint not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value.toString())
    }

    // Validate datetime
    fun validateDateTime(value: java.util.Date, tag: MmTag): ValidationResult {
        if (value.time == 0L) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type datetime not allow empty ${value.toString()}")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        // Format as datetime string
        val format = value.toString()

        return ValidationResult(true, data = value, text = format)
    }

    // Validate UUID
    fun validateUUID(value: String, tag: MmTag): ValidationResult {
        if (value.isEmpty()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type uuid not allow empty value \"\"")
            }
            return ValidationResult(true, data = ByteArray(16), text = value)
        }

        if (!uuidRegex.matcher(value).matches()) {
            return ValidationResult(false, "value '$value' does not match UUID pattern")
        }

        // Parse UUID to bytes (simplified)
        val uuidBytes = ByteArray(16)
        val parts = value.replace("-", "").split("(?<=\\G.{2})").filter { it.isNotEmpty() }
        parts.forEachIndexed { index, part ->
            if (index < 16) {
                uuidBytes[index] = part.toInt(16).toByte()
            }
        }

        if (tag.version != MmTag.DEFAULT_VERSION) {
            // Extract version from UUID
            val version = value.substring(14, 15).toInt(16)
            if (tag.version != version) {
                return ValidationResult(false, "invalid uuid version")
            }
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type uuid not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = uuidBytes, text = value)
    }

    // Validate email
    fun validateEmail(value: String, tag: MmTag): ValidationResult {
        if (value.isEmpty()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type email not allow empty value \"\"")
            }
            return ValidationResult(true, data = value, text = value)
        }

        if (!emailRegex.matcher(value).matches()) {
            return ValidationResult(false, "value '$value' does not match email pattern")
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type email not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = value, text = value)
    }

    // Validate enum
    fun validateEnum(value: String, tag: MmTag): ValidationResult {
        if (value.isEmpty()) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type enum not allow empty value \"\"")
            }
            return ValidationResult(true, data = -1, text = value)
        }

        if (tag.enumValues.isEmpty()) {
            return ValidationResult(false, "enum not defined")
        }

        val enums = tag.enumValues.split("|")
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

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type enum not support location UTC${tag.locationHours}")
        }

        return ValidationResult(true, data = idx, text = value)
    }

    // Validate image
    fun validateImage(value: ByteArray, tag: MmTag): ValidationResult {
        val length = value.size

        if (length == 0) {
            if (!tag.allowEmpty) {
                return ValidationResult(false, "type image not allow empty value []byte{}")
            }
            return ValidationResult(true, data = value, text = "")
        }

        if (tag.min.isNotEmpty()) {
            val mini = tag.min.toIntOrNull()
            if (mini == null) {
                return ValidationResult(false, "failed to parse tag.min as int: ${tag.min}")
            }
            if (length < mini) {
                return ValidationResult(false, "[]byte length $length < min $mini")
            }
        }

        if (tag.max.isNotEmpty()) {
            val maxi = tag.max.toIntOrNull()
            if (maxi == null) {
                return ValidationResult(false, "failed to parse tag.max as int: ${tag.max}")
            }
            if (length > maxi) {
                return ValidationResult(false, "[]byte length $length > max $maxi")
            }
        }

        if (tag.size > 0 && length != tag.size) {
            return ValidationResult(false, "[]byte length $length != size ${tag.size}")
        }

        if (tag.desc.length > 65535) {
            return ValidationResult(false, "desc length exceeds 65535 bytes")
        }

        if (tag.locationHours != 0) {
            return ValidationResult(false, "type image not support location UTC${tag.locationHours}")
        }

        // Convert to base64 string
        val text = Base64.getEncoder().encodeToString(value)

        return ValidationResult(true, data = value, text = text)
    }

    // Validate any value based on type
    fun validate(value: Any?, tag: MmTag): ValidationResult {
        return when (tag.type) {
            ValueType.ARRAY -> {
                if (value is List<*>) {
                    validateArray(value, tag)
                } else {
                    ValidationResult(false, "expected array, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.STRUCT -> {
                validateStruct(tag)
            }
            ValueType.STRING -> {
                if (value is String) {
                    validateString(value, tag)
                } else {
                    ValidationResult(false, "expected string, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.BYTES -> {
                if (value is ByteArray) {
                    validateBytes(value, tag)
                } else {
                    ValidationResult(false, "expected bytes, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.BOOL -> {
                if (value is Boolean) {
                    validateBool(value, tag)
                } else {
                    ValidationResult(false, "expected bool, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.INT8 -> {
                when (value) {
                    is Byte -> validateInt8(value, tag)
                    is Long -> validateInt8(value.toByte(), tag)
                    is Int -> validateInt8(value.toByte(), tag)
                    else -> ValidationResult(false, "expected int8, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.INT16 -> {
                when (value) {
                    is Short -> validateInt16(value, tag)
                    is Long -> validateInt16(value.toShort(), tag)
                    is Int -> validateInt16(value.toShort(), tag)
                    else -> ValidationResult(false, "expected int16, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.INT32 -> {
                when (value) {
                    is Int -> validateInt32(value, tag)
                    is Long -> validateInt32(value.toInt(), tag)
                    else -> ValidationResult(false, "expected int32, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.INT64, ValueType.INT -> {
                when (value) {
                    is Long -> validateInt64(value, tag)
                    else -> ValidationResult(false, "expected int64, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.UINT -> {
                when (value) {
                    is Long -> validateUint(value, tag)
                    is Int -> validateUint(value.toLong(), tag)
                    is BigInteger -> validateUint(value.toLong(), tag)
                    else -> ValidationResult(false, "expected uint, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.UINT8 -> {
                when (value) {
                    is Short -> validateUint8(value, tag)
                    is Long -> validateUint8(value.toShort(), tag)
                    is Int -> validateUint8(value.toShort(), tag)
                    else -> ValidationResult(false, "expected uint8, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.UINT16 -> {
                when (value) {
                    is Int -> validateUint16(value, tag)
                    is Long -> validateUint16(value.toInt(), tag)
                    else -> ValidationResult(false, "expected uint16, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.UINT32 -> {
                when (value) {
                    is Long -> validateUint32(value, tag)
                    is Int -> validateUint32(value.toLong(), tag)
                    else -> ValidationResult(false, "expected uint32, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.UINT64 -> {
                when (value) {
                    is BigInteger -> validateUint64(value, tag)
                    is Long -> validateUint64(BigInteger.valueOf(value), tag)
                    is Short -> validateUint64(BigInteger.valueOf(value.toLong()), tag)
                    is Int -> validateUint64(BigInteger.valueOf(value.toLong()), tag)
                    else -> ValidationResult(false, "expected uint64, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.FLOAT32 -> {
                when (value) {
                    is Float -> validateFloat32(value, tag)
                    is Double -> validateFloat32(value.toFloat(), tag)
                    else -> ValidationResult(false, "expected float32, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.FLOAT64 -> {
                when (value) {
                    is Double -> validateFloat64(value, tag)
                    is Float -> validateFloat64(value.toDouble(), tag)
                    else -> ValidationResult(false, "expected float64, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.BIGINT -> {
                when (value) {
                    is BigInteger -> validateBigInteger(value, tag)
                    is Long -> validateBigInteger(BigInteger.valueOf(value), tag)
                    is Int -> validateBigInteger(BigInteger.valueOf(value.toLong()), tag)
                    is Short -> validateBigInteger(BigInteger.valueOf(value.toLong()), tag)
                    else -> ValidationResult(false, "expected bigint, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.DATETIME, ValueType.DATE, ValueType.TIME -> {
                if (value is java.util.Date) {
                    validateDateTime(value, tag)
                } else {
                    ValidationResult(false, "expected Date, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.UUID -> {
                if (value is String) {
                    validateUUID(value, tag)
                } else {
                    ValidationResult(false, "expected string, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.EMAIL -> {
                if (value is String) {
                    validateEmail(value, tag)
                } else {
                    ValidationResult(false, "expected string, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.ENUM -> {
                if (value is String) {
                    validateEnum(value, tag)
                } else {
                    ValidationResult(false, "expected string, got ${value?.javaClass?.simpleName}")
                }
            }
            ValueType.IMAGE -> {
                if (value is ByteArray) {
                    validateImage(value, tag)
                } else {
                    ValidationResult(false, "expected bytes, got ${value?.javaClass?.simpleName}")
                }
            }
            else -> {
                ValidationResult(true, data = value, text = value?.toString() ?: "null")
            }
        }
    }
}
