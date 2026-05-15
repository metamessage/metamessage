package io.github.metamessage.jsonc

import io.github.metamessage.ir.Object
import io.github.metamessage.ir.Value
import io.github.metamessage.ir.ValueType
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.Field

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.*

class JsoncMMTest {

    @Test
    fun parseMmTagInComment() {
        val source = """
            {
                // mm: type=str; desc=用戶名
                "name": "張三"
            }
        """.trimIndent()
        val result = parseJsonc(source)
        assertTrue(result is Object)
        val obj = result as Object
        assertEquals("name", obj.fields[0].key)
        assertEquals("張三", obj.fields[0].value.let { it as? Value }?.data)
        assertEquals("用戶名", obj.fields[0].value.let { it as? Value }?.tag?.desc)
        assertEquals(ValueType.STRING, obj.fields[0].value.let { it as? Value }?.tag?.type)
    }

    @Test
    fun parseMmTagBlockComment() {
        val source = """
            {
                /* mm: type=i; desc=年齡 */
                "age": 25
            }
        """.trimIndent()
        val result = parseJsonc(source)
        assertTrue(result is Object)
        val obj = result as Object
        assertEquals("age", obj.fields[0].key)
        assertEquals(25L, obj.fields[0].value.let { it as? Value }?.data)
        assertEquals("年齡", obj.fields[0].value.let { it as? Value }?.tag?.desc)
        assertEquals(ValueType.INT, obj.fields[0].value.let { it as? Value }?.tag?.type)
    }

    @Test
    fun parseNonMmCommentIgnored() {
        val source = """
            {
                // 這是普通註釋，不是 tag
                "name": "李四"
            }
        """.trimIndent()
        val result = parseJsonc(source)
        assertTrue(result is Object)
        val obj = result as Object
        assertEquals("name", obj.fields[0].key)
        val valueTag = obj.fields[0].value.let { it as? Value }?.tag
        assertNotNull(valueTag)
        assertEquals("", valueTag!!.desc)
        assertEquals(ValueType.STRING, valueTag.type)
    }

    @Test
    fun printWithMmTag() {
        val obj = Object()
        val tag = Tag()
        tag.type = ValueType.UUID
        tag.desc = "user id"
        val value = Value(data = "550e8400-e29b-41d4-a716-446655440000", text = "\"550e8400-e29b-41d4-a716-446655440000\"", tag = tag)
        obj.fields.add(Field("userId", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("// mm:"))
        assertTrue(output.contains("type=uuid"))
        assertTrue(output.contains("desc="))
        assertTrue(output.contains("\"550e8400-e29b-41d4-a716-446655440000\""))
    }

    @Test
    fun printNumberWithoutQuotes() {
        val obj = Object()
        val tag = Tag()
        tag.type = ValueType.INT
        tag.desc = "年齡"
        val value = Value(data = 25L, text = "25", tag = tag)
        obj.fields.add(Field("age", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("25"))
        assertFalse(output.contains("\"25\""))
    }

    @Test
    fun printFloatWithoutQuotes() {
        val obj = Object()
        val tag = Tag()
        tag.type = ValueType.FLOAT64
        tag.desc = "價格"
        val value = Value(data = 3.14, text = "3.14", tag = tag)
        obj.fields.add(Field("price", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("3.14"))
        assertFalse(output.contains("\"3.14\""))
    }

    @Test
    fun printBoolWithoutQuotes() {
        val obj = Object()
        val tag = Tag()
        tag.type = ValueType.BOOL
        val value = Value(data = true, text = "true", tag = tag)
        obj.fields.add(Field("active", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("true"))
        assertFalse(output.contains("\"true\""))
    }

    @Test
    fun printCompactNoQuotesForNumbers() {
        val value = Value(data = 123L, text = "123")
        val output = JsoncPrinter.toCompactString(value)
        assertEquals("123", output)
    }

    @Test
    fun printCompactNoQuotesForBool() {
        val value = Value(data = false, text = "false")
        val output = JsoncPrinter.toCompactString(value)
        assertEquals("false", output)
    }
}