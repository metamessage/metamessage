package io.github.metamessage

import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.ValueType
import java.math.BigInteger
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class ZeroValueTest {

    @MM
    class AllZeroWithAllowEmpty(
            @MM(type = ValueType.I8, allowEmpty = true) var i8: Byte = 0,
            @MM(type = ValueType.I16, allowEmpty = true) var i16: Short = 0,
            @MM(type = ValueType.I32, allowEmpty = true) var i32: Int = 0,
            @MM(type = ValueType.I64, allowEmpty = true) var i64: Long = 0L,
            @MM(type = ValueType.U8, allowEmpty = true) var u8: Byte = 0,
            @MM(type = ValueType.U16, allowEmpty = true) var u16: Short = 0,
            @MM(type = ValueType.U32, allowEmpty = true) var u32: Int = 0,
            @MM(type = ValueType.U64, allowEmpty = true) var u64: Long = 0L,
            @MM(type = ValueType.F32, allowEmpty = true) var f32: Float = 0.0f,
            @MM(type = ValueType.F64, allowEmpty = true) var f64: Double = 0.0,
            @MM(type = ValueType.BIGINT, allowEmpty = true) var bigint: BigInteger = BigInteger.ZERO
    )

    @MM
    class ZeroWithoutAllowEmpty(
            @MM(type = ValueType.I32) var age: Int = 0
    )

    @Test
    fun encodeAllZeroValuesWithAllowEmpty() {
        val obj = AllZeroWithAllowEmpty()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, AllZeroWithAllowEmpty::class.java)
        assertEquals(0.toByte(), out.i8)
        assertEquals(0.toShort(), out.i16)
        assertEquals(0, out.i32)
        assertEquals(0L, out.i64)
        assertEquals(0.toByte(), out.u8)
        assertEquals(0.toShort(), out.u16)
        assertEquals(0, out.u32)
        assertEquals(0L, out.u64)
        assertEquals(0.0f, out.f32)
        assertEquals(0.0, out.f64)
        assertEquals(BigInteger.ZERO, out.bigint)
    }

    @Test
    fun encodeZeroWithoutAllowEmptyFails() {
        val obj = ZeroWithoutAllowEmpty()
        try {
            MetaMessage.encodeFromValue(obj)
            fail("Should have thrown exception for zero value without allowEmpty")
        } catch (e: IllegalArgumentException) {
            assertTrue(e.message?.contains("not allow empty") ?: false, "Error: ${e.message}")
        }
    }

    @Test
    fun directValidateI32WithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.I32
            allowEmpty = true
        }
        val result = tag.validateI32(0)
        assertTrue(result.valid, "allowEmpty=true should allow 0, but got: ${result.error}")
    }

    @Test
    fun directValidateI32WithoutAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.I32
            allowEmpty = false
        }
        val result = tag.validateI32(0)
        assertFalse(result.valid)
    }

    @Test
    fun directValidateU8WithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.U8
            allowEmpty = true
        }
        val result = tag.validateU8(0.toShort())
        assertTrue(result.valid, "allowEmpty=true should allow 0, but got: ${result.error}")
    }

    @Test
    fun directValidateI64WithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.I64
            allowEmpty = true
        }
        val result = tag.validateI64(0L)
        assertTrue(result.valid, "allowEmpty=true should allow 0, but got: ${result.error}")
    }

    @Test
    fun directValidateF32WithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.F32
            allowEmpty = true
        }
        val result = tag.validateF32(0.0f)
        assertTrue(result.valid, "allowEmpty=true should allow 0.0, but got: ${result.error}")
    }

    @Test
    fun directValidateBigintWithAllowEmpty() {
        val tag = Tag().apply {
            type = ValueType.BIGINT
            allowEmpty = true
        }
        val result = tag.validateBigint(BigInteger.ZERO)
        assertTrue(result.valid, "allowEmpty=true should allow 0, but got: ${result.error}")
    }
}
