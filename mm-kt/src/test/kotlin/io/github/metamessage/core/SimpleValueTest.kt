package io.github.metamessage.core

import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class SimpleValueTest {

    @Test
    fun allConstantsAreValid() {
        for (v in 0..31) {
            assertTrue(SimpleValue.isValid(v), "SimpleValue($v) should be valid")
        }
    }

    @Test
    fun outOfRangeIsInvalid() {
        assertFalse(SimpleValue.isValid(32))
        assertFalse(SimpleValue.isValid(100))
    }

    @Test
    fun simpleNullValue() {
        assertEquals(0, SimpleValue.SIMPLE_NULL)
        assertEquals("simple_null", SimpleValue.toString(SimpleValue.SIMPLE_NULL))
    }

    @Test
    fun nullBoolValue() {
        assertEquals(1, SimpleValue.NULL_BOOL)
        assertEquals("null_bool", SimpleValue.toString(SimpleValue.NULL_BOOL))
    }

    @Test
    fun nullIntValue() {
        assertEquals(2, SimpleValue.NULL_INT)
        assertEquals("null_int", SimpleValue.toString(SimpleValue.NULL_INT))
    }

    @Test
    fun nullFloatValue() {
        assertEquals(3, SimpleValue.NULL_FLOAT)
        assertEquals("null_float", SimpleValue.toString(SimpleValue.NULL_FLOAT))
    }

    @Test
    fun nullStringValue() {
        assertEquals(4, SimpleValue.NULL_STRING)
        assertEquals("null_string", SimpleValue.toString(SimpleValue.NULL_STRING))
    }

    @Test
    fun nullBytesValue() {
        assertEquals(5, SimpleValue.NULL_BYTES)
        assertEquals("null_bytes", SimpleValue.toString(SimpleValue.NULL_BYTES))
    }

    @Test
    fun falseValue() {
        assertEquals(6, SimpleValue.FALSE)
        assertEquals("false", SimpleValue.toString(SimpleValue.FALSE))
    }

    @Test
    fun trueValue() {
        assertEquals(7, SimpleValue.TRUE)
        assertEquals("true", SimpleValue.toString(SimpleValue.TRUE))
    }

    @Test
    fun codeAlias() {
        assertEquals(8, SimpleValue.CODE)
        assertEquals("code", SimpleValue.toString(SimpleValue.CODE))
    }

    @Test
    fun messageAlias() {
        assertEquals(9, SimpleValue.MESSAGE)
        assertEquals("message", SimpleValue.toString(SimpleValue.MESSAGE))
    }

    @Test
    fun dataAlias() {
        assertEquals(10, SimpleValue.DATA)
        assertEquals("data", SimpleValue.toString(SimpleValue.DATA))
    }

    @Test
    fun idAlias() {
        assertEquals(18, SimpleValue.ID)
        assertEquals("id", SimpleValue.toString(SimpleValue.ID))
    }

    @Test
    fun nameAlias() {
        assertEquals(19, SimpleValue.NAME)
        assertEquals("name", SimpleValue.toString(SimpleValue.NAME))
    }

    @Test
    fun typeAlias() {
        assertEquals(21, SimpleValue.TYPE)
        assertEquals("type", SimpleValue.toString(SimpleValue.TYPE))
    }

    @Test
    fun keyAlias() {
        assertEquals(31, SimpleValue.KEY)
        assertEquals("key", SimpleValue.toString(SimpleValue.KEY))
    }

    @Test
    fun allFieldNameAliasesMapCorrectly() {
        val expectedNames =
                mapOf(
                        8 to "code",
                        9 to "message",
                        10 to "data",
                        11 to "success",
                        12 to "error",
                        13 to "unknown",
                        14 to "page",
                        15 to "limit",
                        16 to "offset",
                        17 to "total",
                        18 to "id",
                        19 to "name",
                        20 to "description",
                        21 to "type",
                        22 to "version",
                        23 to "status",
                        24 to "url",
                        25 to "create_time",
                        26 to "update_time",
                        27 to "delete_time",
                        28 to "account",
                        29 to "token",
                        30 to "expire_time",
                        31 to "key"
                )
        for ((code, name) in expectedNames) {
            assertEquals(
                    name,
                    SimpleValue.toString(code),
                    "SimpleValue($code) should map to '$name'"
            )
        }
    }

    @Test
    fun outOfRangeToString() {
        assertEquals("SimpleValue(32)", SimpleValue.toString(32))
        assertEquals("SimpleValue(100)", SimpleValue.toString(100))
    }

    @Test
    fun uniqueValues() {
        val seen = mutableSetOf<Int>()
        val values =
                listOf(
                        SimpleValue.SIMPLE_NULL,
                        SimpleValue.NULL_BOOL,
                        SimpleValue.NULL_INT,
                        SimpleValue.NULL_FLOAT,
                        SimpleValue.NULL_STRING,
                        SimpleValue.NULL_BYTES,
                        SimpleValue.FALSE,
                        SimpleValue.TRUE,
                        SimpleValue.CODE,
                        SimpleValue.MESSAGE,
                        SimpleValue.DATA,
                        SimpleValue.SUCCESS,
                        SimpleValue.ERROR,
                        SimpleValue.UNKNOWN,
                        SimpleValue.PAGE,
                        SimpleValue.LIMIT,
                        SimpleValue.OFFSET,
                        SimpleValue.TOTAL,
                        SimpleValue.ID,
                        SimpleValue.NAME,
                        SimpleValue.DESCRIPTION,
                        SimpleValue.TYPE,
                        SimpleValue.VERSION,
                        SimpleValue.STATUS,
                        SimpleValue.URL,
                        SimpleValue.CREATE_TIME,
                        SimpleValue.UPDATE_TIME,
                        SimpleValue.DELETE_TIME,
                        SimpleValue.ACCOUNT,
                        SimpleValue.TOKEN,
                        SimpleValue.EXPIRE_TIME,
                        SimpleValue.KEY
                )
        for (v in values) {
            assertTrue(seen.add(v), "Duplicate SimpleValue: $v")
        }
        assertEquals(32, seen.size)
    }
}
