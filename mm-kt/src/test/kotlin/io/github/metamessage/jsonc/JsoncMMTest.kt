package io.github.metamessage.jsonc

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
        assertTrue(result is JsoncObject)
        val obj = result as JsoncObject
        assertEquals("name", obj.fields[0].key)
        assertEquals("張三", obj.fields[0].value?.let { it as? JsoncValue }?.data)
        assertEquals("用戶名", obj.fields[0].value?.let { it as? JsoncValue }?.tag?.desc)
        assertEquals(JsoncValueType.String, obj.fields[0].value?.let { it as? JsoncValue }?.tag?.type)
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
        assertTrue(result is JsoncObject)
        val obj = result as JsoncObject
        assertEquals("age", obj.fields[0].key)
        assertEquals(25L, obj.fields[0].value?.let { it as? JsoncValue }?.data)
        assertEquals("年齡", obj.fields[0].value?.let { it as? JsoncValue }?.tag?.desc)
        assertEquals(JsoncValueType.Int, obj.fields[0].value?.let { it as? JsoncValue }?.tag?.type)
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
        assertTrue(result is JsoncObject)
        val obj = result as JsoncObject
        assertEquals("name", obj.fields[0].key)
        val valueTag = obj.fields[0].value?.let { it as? JsoncValue }?.tag
        assertNotNull(valueTag)
        assertEquals("", valueTag!!.desc)
        assertEquals(JsoncValueType.String, valueTag.type)
    }

    @Test
    fun printWithMmTag() {
        val obj = JsoncObject()
        val tag = JsoncTag()
        tag.type = JsoncValueType.String
        tag.desc = "用戶名"
        val value = JsoncValue(data = "張三", text = "張三", tag = tag)
        obj.fields.add(JsoncField("name", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("// mm:"))
        assertTrue(output.contains("type=str"))
        assertTrue(output.contains("desc=用戶名"))
        assertTrue(output.contains("\"張三\""))
    }

    @Test
    fun printNumberWithoutQuotes() {
        val obj = JsoncObject()
        val tag = JsoncTag()
        tag.type = JsoncValueType.Int
        tag.desc = "年齡"
        val value = JsoncValue(data = 25L, text = "25", tag = tag)
        obj.fields.add(JsoncField("age", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("25"))
        assertFalse(output.contains("\"25\""))
    }

    @Test
    fun printFloatWithoutQuotes() {
        val obj = JsoncObject()
        val tag = JsoncTag()
        tag.type = JsoncValueType.Float64
        tag.desc = "價格"
        val value = JsoncValue(data = 3.14, text = "3.14", tag = tag)
        obj.fields.add(JsoncField("price", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("3.14"))
        assertFalse(output.contains("\"3.14\""))
    }

    @Test
    fun printBoolWithoutQuotes() {
        val obj = JsoncObject()
        val tag = JsoncTag()
        tag.type = JsoncValueType.Bool
        val value = JsoncValue(data = true, text = "true", tag = tag)
        obj.fields.add(JsoncField("active", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("true"))
        assertFalse(output.contains("\"true\""))
    }

    @Test
    fun printNullWithoutQuotes() {
        val obj = JsoncObject()
        val tag = JsoncTag()
        tag.isNull = true
        tag.type = JsoncValueType.Null
        val value = JsoncValue(data = null, text = "null", tag = tag)
        obj.fields.add(JsoncField("optional", value))

        val output = JsoncPrinter.toString(obj)
        assertTrue(output.contains("null"))
        assertFalse(output.contains("\"null\""))
    }

    @Test
    fun printCompactNoQuotesForNumbers() {
        val value = JsoncValue(data = 123L, text = "123")
        val output = JsoncPrinter.toCompactString(value)
        assertEquals("123", output)
    }

    @Test
    fun printCompactNoQuotesForBool() {
        val value = JsoncValue(data = false, text = "false")
        val output = JsoncPrinter.toCompactString(value)
        assertEquals("false", output)
    }
}