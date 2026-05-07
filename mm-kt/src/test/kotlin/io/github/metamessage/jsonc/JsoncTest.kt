package io.github.metamessage.jsonc

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.*

class JsoncParserTest {

    @Test
    fun parseEmptyObject() {
        val source = "{}"
        val result = parseJsonc(source)
        assertTrue(result is JsoncObject)
        assertEquals(0, (result as JsoncObject).fields.size)
    }

    @Test
    fun parseSimpleObject() {
        val source = """{"key": "value"}"""
        val result = parseJsonc(source)
        assertTrue(result is JsoncObject)
        val obj = result as JsoncObject
        assertEquals(1, obj.fields.size)
        assertEquals("key", obj.fields[0].key)
    }

    @Test
    fun parseNumber() {
        val source = "123"
        val result = parseJsonc(source)
        assertTrue(result is JsoncValue)
        val value = result as JsoncValue
        assertEquals(123L, value.data)
    }

    @Test
    fun parseFloat() {
        val source = "3.14"
        val result = parseJsonc(source)
        assertTrue(result is JsoncValue)
        val value = result as JsoncValue
        assertEquals(3.14, value.data)
    }

    @Test
    fun parseBoolean() {
        val sourceTrue = "true"
        val resultTrue = parseJsonc(sourceTrue)
        assertTrue(resultTrue is JsoncValue)
        assertEquals(true, (resultTrue as JsoncValue).data)

        val sourceFalse = "false"
        val resultFalse = parseJsonc(sourceFalse)
        assertTrue(resultFalse is JsoncValue)
        assertEquals(false, (resultFalse as JsoncValue).data)
    }

    @Test
    fun parseNull() {
        val source = "null"
        val result = parseJsonc(source)
        assertTrue(result is JsoncValue)
        val value = result as JsoncValue
        assertNull(value.data)
        assertTrue(value.tag?.isNull == true)
    }

    @Test
    fun parseArray() {
        val source = "[1, 2, 3]"
        val result = parseJsonc(source)
        assertTrue(result is JsoncArray)
        val arr = result as JsoncArray
        assertEquals(3, arr.items.size)
    }

    @Test
    fun parseNestedObject() {
        val source = """{"outer": {"inner": "value"}}"""
        val result = parseJsonc(source)
        assertTrue(result is JsoncObject)
        val outer = result as JsoncObject
        assertEquals("outer", outer.fields[0].key)
        assertTrue(outer.fields[0].value is JsoncObject)
        val inner = outer.fields[0].value as JsoncObject
        assertEquals("inner", inner.fields[0].key)
    }

    @Test
    fun parseWithLineComment() {
        val source = """
            {
                // this is a comment
                "key": "value"
            }
        """.trimIndent()
        val result = parseJsonc(source)
        assertTrue(result is JsoncObject)
    }

    @Test
    fun parseWithBlockComment() {
        val source = """
            {
                /* this is a
                   block comment */
                "key": "value"
            }
        """.trimIndent()
        val result = parseJsonc(source)
        assertTrue(result is JsoncObject)
    }

    @Test
    fun parseWithTrailingComma() {
        val source = """{"key": "value",}"""
        val result = parseJsonc(source)
        assertTrue(result is JsoncObject)
        val obj = result as JsoncObject
        assertEquals(1, obj.fields.size)
    }
}

class JsoncPrinterTest {

    @Test
    fun printEmptyObject() {
        val obj = JsoncObject()
        val result = JsoncPrinter.toString(obj)
        assertTrue(result.contains("{"))
        assertTrue(result.contains("}"))
    }

    @Test
    fun printSimpleObject() {
        val obj = JsoncObject()
        obj.fields.add(JsoncField("key", JsoncValue(data = "value", text = "\"value\"")))
        val result = JsoncPrinter.toString(obj)
        assertTrue(result.contains("key"))
    }

    @Test
    fun printCompact() {
        val obj = JsoncObject()
        obj.fields.add(JsoncField("key", JsoncValue(data = "value", text = "\"value\"")))
        val result = JsoncPrinter.toCompactString(obj)
        assertFalse(result.contains("\n"))
    }
}

class JsoncScannerTest {

    @Test
    fun scanEmptyInput() {
        val scanner = JsoncScanner("")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.EOF, token.type)
    }

    @Test
    fun scanLBrace() {
        val scanner = JsoncScanner("{")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.LBrace, token.type)
    }

    @Test
    fun scanRBrace() {
        val scanner = JsoncScanner("}")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.RBrace, token.type)
    }

    @Test
    fun scanString() {
        val scanner = JsoncScanner("\"hello\"")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.String, token.type)
        assertEquals("hello", token.literal)
    }

    @Test
    fun scanNumber() {
        val scanner = JsoncScanner("123")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.Number, token.type)
        assertEquals("123", token.literal)
    }

    @Test
    fun scanTrue() {
        val scanner = JsoncScanner("true")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.True, token.type)
    }

    @Test
    fun scanFalse() {
        val scanner = JsoncScanner("false")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.False, token.type)
    }

    @Test
    fun scanLineCommentTrailing() {
        val scanner = JsoncScanner("// this is a comment\n")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.TrailingComment, token.type)
    }

    @Test
    fun scanBlockCommentTrailing() {
        val scanner = JsoncScanner("/* block comment */")
        val token = scanner.nextToken()
        assertEquals(JsoncTokenType.TrailingComment, token.type)
    }
}