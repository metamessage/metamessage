package io.github.metamessage.core

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.*
import java.math.BigInteger
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.ValueType
import java.time.ZoneId


class ValidatorTest {

    // ==================== Integer Tests ====================

    @Test
    fun validateInt8_ValidValue() {
        val tag = Tag().apply { type = ValueType.INT8 }
        val result = tag.validateInt8(42.toByte())
        assertTrue(result.valid)
        assertEquals("42", result.text)
    }

    @Test
    fun validateInt8_MinValue() {
        val tag = Tag().apply { type = ValueType.INT8 }
        val result = tag.validateInt8((-128).toByte())
        assertTrue(result.valid)
    }

    @Test
    fun validateInt8_MaxValue() {
        val tag = Tag().apply { type = ValueType.INT8 }
        val result = tag.validateInt8(127.toByte())
        assertTrue(result.valid)
    }

    @Test
    fun validateInt8_WithMinConstraint() {
        val tag = Tag().apply {
            type = ValueType.INT8
            min = "10"
        }
        val result = tag.validateInt8(15.toByte())
        assertTrue(result.valid)
    }

    @Test
    fun validateInt8_BelowMinConstraint() {
        val tag = Tag().apply {
            type = ValueType.INT8
            min = "10"
        }
        val result = tag.validateInt8(5.toByte())
        assertFalse(result.valid)
        assertTrue(result.error?.contains("less than the minimum") ?: false)
    }

    @Test
    fun validateInt8_WithMaxConstraint() {
        val tag = Tag().apply {
            type = ValueType.INT8
            max = "50"
        }
        val result = tag.validateInt8(42.toByte())
        assertTrue(result.valid)
    }

    @Test
    fun validateInt8_ExceedsMaxConstraint() {
        val tag = Tag().apply {
            type = ValueType.INT8
            max = "50"
        }
        val result = tag.validateInt8(60.toByte())
        assertFalse(result.valid)
        assertTrue(result.error?.contains("exceeds the maximum") ?: false)
    }

    @Test
    fun validateInt16_ValidValue() {
        val tag = Tag().apply { type = ValueType.INT16 }
        val result = tag.validateInt16(1234.toShort())
        assertTrue(result.valid)
        assertEquals("1234", result.text)
    }

    @Test
    fun validateInt32_ValidValue() {
        val tag = Tag().apply { type = ValueType.INT32 }
        val result = tag.validateInt32(123456789)
        assertTrue(result.valid)
        assertEquals("123456789", result.text)
    }

    @Test
    fun validateInt64_ValidValue() {
        val tag = Tag().apply { type = ValueType.INT64 }
        val result = tag.validateInt64(1234567890123456789L)
        assertTrue(result.valid)
        assertEquals("1234567890123456789", result.text)
    }

    @Test
    fun validateUint_ValidValue() {
        val tag = Tag().apply { type = ValueType.UINT }
        val result = tag.validateUint(4294967295L)
        assertTrue(result.valid)
    }

    @Test
    fun validateUint8_ValidValue() {
        val tag = Tag().apply { type = ValueType.UINT8 }
        val result = tag.validateUint8(255.toShort())
        assertTrue(result.valid)
    }

    @Test
    fun validateUint16_ValidValue() {
        val tag = Tag().apply { type = ValueType.UINT16 }
        val result = tag.validateUint16(65535)
        assertTrue(result.valid)
    }

    @Test
    fun validateUint32_ValidValue() {
        val tag = Tag().apply { type = ValueType.UINT32 }
        val result = tag.validateUint32(4294967295L)
        assertTrue(result.valid)
    }

    @Test
    fun validateUint64_ValidValue() {
        val tag = Tag().apply { type = ValueType.UINT64 }
        val result = tag.validateUint64(BigInteger("18446744073709551615"))
        assertTrue(result.valid)
    }

    // ==================== Float Tests ====================

    @Test
    fun validateFloat32_ValidValue() {
        val tag = Tag().apply { type = ValueType.FLOAT32 }
        val result = tag.validateFloat32(3.14f)
        assertTrue(result.valid)
    }

    @Test
    fun validateFloat32_WithMinMax() {
        val tag = Tag().apply {
            type = ValueType.FLOAT32
            min = "1.0"
            max = "10.0"
        }
        val result = tag.validateFloat32(3.14f)
        assertTrue(result.valid)
    }

    @Test
    fun validateFloat64_ValidValue() {
        val tag = Tag().apply { type = ValueType.FLOAT64 }
        val result = tag.validateFloat64(3.141592653589793)
        assertTrue(result.valid)
    }

    // ==================== String Tests ====================

    @Test
    fun validateString_ValidValue() {
        val tag = Tag().apply { type = ValueType.STRING }
        val result = tag.validateString("hello")
        assertTrue(result.valid)
        assertEquals("hello", result.text)
    }

    @Test
    fun validateString_WithPattern() {
        val tag = Tag().apply {
            type = ValueType.STRING
            pattern = "^[a-z]+$"
        }
        val result = tag.validateString("hello")
        assertTrue(result.valid)
    }

    @Test
    fun validateString_PatternMismatch() {
        val tag = Tag().apply {
            type = ValueType.STRING
            pattern = "^[a-z]+$"
        }
        val result = tag.validateString("Hello123")
        assertFalse(result.valid)
    }

    @Test
    fun validateString_WithMinMaxLength() {
        val tag = Tag().apply {
            type = ValueType.STRING
            min = "3"
            max = "10"
        }
        val result = tag.validateString("hello")
        assertTrue(result.valid)
    }

    // ==================== Bytes Tests ====================

    @Test
    fun validateBytes_ValidValue() {
        val tag = Tag().apply { type = ValueType.BYTES }
        val result = tag.validateBytes(byteArrayOf(0x12, 0x34, 0x56))
        assertTrue(result.valid)
        assertNotNull(result.text)
    }

    // ==================== Bool Tests ====================

    @Test
    fun validateBool_True() {
        val tag = Tag().apply { type = ValueType.BOOL }
        val result = tag.validateBool(true)
        assertTrue(result.valid)
        assertEquals("true", result.text)
    }

    @Test
    fun validateBool_False() {
        val tag = Tag().apply { type = ValueType.BOOL }
        val result = tag.validateBool(false)
        assertTrue(result.valid)
        assertEquals("false", result.text)
    }

    // ==================== UUID Tests ====================

    @Test
    fun validateUUID_ValidValue() {
        val tag = Tag().apply { type = ValueType.UUID }
        val result = tag.validateUUID("550e8400-e29b-41d4-a716-446655440000")
        assertTrue(result.valid)
    }

    @Test
    fun validateUUID_InvalidFormat() {
        val tag = Tag().apply { type = ValueType.UUID }
        val result = tag.validateUUID("invalid-uuid")
        assertFalse(result.valid)
    }

    // ==================== Email Tests ====================

    @Test
    fun validateEmail_ValidValue() {
        val tag = Tag().apply { type = ValueType.EMAIL }
        val result = tag.validateEmail("test@example.com")
        assertTrue(result.valid)
    }

    @Test
    fun validateEmail_InvalidFormat() {
        val tag = Tag().apply { type = ValueType.EMAIL }
        val result = tag.validateEmail("invalid-email")
        assertFalse(result.valid)
    }

    // ==================== Enum Tests ====================

    @Test
    fun validateEnum_ValidValue() {
        val tag = Tag().apply {
            type = ValueType.ENUM
            enum = "apple|banana|cherry"
        }
        val result = tag.validateEnum("banana")
        assertTrue(result.valid)
    }

    @Test
    fun validateEnum_InvalidValue() {
        val tag = Tag().apply {
            type = ValueType.ENUM
            enum = "apple|banana|cherry"
        }
        val result = tag.validateEnum("orange")
        assertFalse(result.valid)
    }

    // ==================== DateTime Tests ====================

    @Test
    fun validateDateTime_ValidValue() {
        val tag = Tag().apply { type = ValueType.DATETIME }
        val date = java.util.Date()
        val localDateTime = date.toInstant()
            .atZone(ZoneId.systemDefault())
            .toLocalDateTime()
        val result = tag.validateDateTime(localDateTime)
        assertTrue(result.valid)
    }

    @Test
    fun validateDate_ValidValue() {
        val tag = Tag().apply { type = ValueType.DATE }
        val date = java.util.Date()
        val localDate = date.toInstant()
            .atZone(ZoneId.systemDefault())
            .toLocalDate()
        val result = tag.validateDate(localDate)
        assertTrue(result.valid)
    }

    @Test
    fun validateTime_ValidValue() {
        val tag = Tag().apply { type = ValueType.TIME }
        val date = java.util.Date()
        val localTime = date.toInstant()
            .atZone(ZoneId.systemDefault())
            .toLocalTime()  
        val result = tag.validateTime(localTime)
        assertTrue(result.valid)
    }

    // ==================== Array Tests ====================

    @Test
    fun validateArray_ValidValue() {
        val tag = Tag().apply { type = ValueType.ARRAY }
        val result = tag.validateArray(listOf(1, 2, 3))
        assertTrue(result.valid)
    }

    @Test
    fun validateArray_ChildUnique() {
        val tag = Tag().apply {
            type = ValueType.ARRAY
            childUnique = true
        }
        val result = tag.validateArray(listOf(1, 2, 3))
        assertTrue(result.valid)
    }

    @Test
    fun validateArray_DuplicateValues() {
        val tag = Tag().apply {
            type = ValueType.ARRAY
            childUnique = true
        }
        val result = tag.validateArray(listOf(1, 2, 2, 3))
        assertFalse(result.valid)
        assertTrue(result.error?.contains("duplicate") ?: false)
    }

    // ==================== Struct Tests ====================

    @Test
    fun validateStruct_Valid() {
        val tag = Tag().apply { type = ValueType.STRUCT }
        val result = tag.validateStruct()
        assertTrue(result.valid)
    }

    // ==================== BigInt Tests ====================

    @Test
    fun validateBigInt_ValidValue() {
        val tag = Tag().apply { type = ValueType.BIGINT }
        val result = tag.validateBigInt(BigInteger("123456789012345678901234567890"))
        assertTrue(result.valid)
    }

    // ==================== Boundary Tests ====================

    @Test
    fun validateInt8_ZeroWithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.INT8
            allowEmpty = true
        }
        val result = tag.validateInt8(0.toByte())
        assertTrue(result.valid)
    }

    @Test
    fun validateInt8_ZeroWithoutAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.INT8
            allowEmpty = false
        }
        val result = tag.validateInt8(0.toByte())
        assertFalse(result.valid)
    }

    @Test
    fun validateString_EmptyWithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.STRING
            allowEmpty = true
        }
        val result = tag.validateString("")
        assertTrue(result.valid)
    }

    @Test
    fun validateString_EmptyWithoutAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.STRING
            allowEmpty = false
        }
        val result = tag.validateString("")
        assertFalse(result.valid)
    }

    // ==================== Error Handling Tests ====================

    @Test
    fun validateInt8_InvalidMinValue() {
        val tag = Tag().apply {
            type = ValueType.INT8
            min = "invalid"
        }
        val result = tag.validateInt8(42.toByte())
        assertFalse(result.valid)
        assertTrue(result.error?.contains("failed to parse") ?: false)
    }

    @Test
    fun validateInt8_MinOutOfRange() {
        val tag = Tag().apply {
            type = ValueType.INT8
            min = "200"
        }
        val result = tag.validateInt8(42.toByte())
        assertFalse(result.valid)
        assertTrue(result.error?.contains("out of int8 range") ?: false)
    }

    @Test
    fun validateUint64_InvalidMinValue() {
        val tag = Tag().apply {
            type = ValueType.UINT64
            min = "invalid"
        }
        val result = tag.validateUint64(BigInteger("100"))
        assertFalse(result.valid)
    }
}
