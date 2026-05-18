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
    fun nullBoolValue() {
        assertEquals(0, SimpleValue.NULL_BOOL)
        assertEquals("null_bool", SimpleValue.toString(SimpleValue.NULL_BOOL))
    }

    @Test
    fun nullIntValue() {
        assertEquals(1, SimpleValue.NULL_INT)
        assertEquals("null_int", SimpleValue.toString(SimpleValue.NULL_INT))
    }

    @Test
    fun nullFloatValue() {
        assertEquals(2, SimpleValue.NULL_FLOAT)
        assertEquals("null_float", SimpleValue.toString(SimpleValue.NULL_FLOAT))
    }

    @Test
    fun nullStringValue() {
        assertEquals(3, SimpleValue.NULL_STRING)
        assertEquals("null_string", SimpleValue.toString(SimpleValue.NULL_STRING))
    }

    @Test
    fun nullBytesValue() {
        assertEquals(4, SimpleValue.NULL_BYTES)
        assertEquals("null_bytes", SimpleValue.toString(SimpleValue.NULL_BYTES))
    }

    @Test
    fun falseValue() {
        assertEquals(5, SimpleValue.FALSE)
        assertEquals("false", SimpleValue.toString(SimpleValue.FALSE))
    }

    @Test
    fun trueValue() {
        assertEquals(6, SimpleValue.TRUE)
        assertEquals("true", SimpleValue.toString(SimpleValue.TRUE))
    }

    @Test
    fun codeAlias() {
        assertEquals(7, SimpleValue.CODE)
        assertEquals("code", SimpleValue.toString(SimpleValue.CODE))
    }

    @Test
    fun messageAlias() {
        assertEquals(8, SimpleValue.MESSAGE)
        assertEquals("message", SimpleValue.toString(SimpleValue.MESSAGE))
    }

    @Test
    fun dataAlias() {
        assertEquals(9, SimpleValue.DATA)
        assertEquals("data", SimpleValue.toString(SimpleValue.DATA))
    }

    @Test
    fun idAlias() {
        assertEquals(17, SimpleValue.ID)
        assertEquals("id", SimpleValue.toString(SimpleValue.ID))
    }

    @Test
    fun nameAlias() {
        assertEquals(18, SimpleValue.NAME)
        assertEquals("name", SimpleValue.toString(SimpleValue.NAME))
    }

    @Test
    fun typeAlias() {
        assertEquals(20, SimpleValue.TYPE)
        assertEquals("type", SimpleValue.toString(SimpleValue.TYPE))
    }

    @Test
    fun keyAlias() {
        assertEquals(30, SimpleValue.KEY)
        assertEquals("key", SimpleValue.toString(SimpleValue.KEY))
    }

    @Test
    fun valueAlias() {
        assertEquals(31, SimpleValue.VAL)
        assertEquals("value", SimpleValue.toString(SimpleValue.VAL))
    }

    @Test
    fun allFieldNameAliasesMapCorrectly() {
        val expectedNames =
                mapOf(
                        7 to "code",
                        8 to "message",
                        9 to "data",
                        10 to "success",
                        11 to "error",
                        12 to "unknown",
                        13 to "page",
                        14 to "limit",
                        15 to "offset",
                        16 to "total",
                        17 to "id",
                        18 to "name",
                        19 to "description",
                        20 to "type",
                        21 to "version",
                        22 to "status",
                        23 to "url",
                        24 to "create_time",
                        25 to "update_time",
                        26 to "delete_time",
                        27 to "account",
                        28 to "token",
                        29 to "expire_time",
                        30 to "key",
                        31 to "value"
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
                        SimpleValue.KEY,
                        SimpleValue.VAL
                )
        for (v in values) {
            assertTrue(seen.add(v), "Duplicate SimpleValue: $v")
        }
        assertEquals(32, seen.size)
    }
}
