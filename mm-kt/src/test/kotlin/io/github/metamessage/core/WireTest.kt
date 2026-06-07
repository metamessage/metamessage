package io.github.metamessage.core

import io.github.metamessage.ir.*
import io.github.metamessage.ir.NodeArray as AstArray
import java.math.BigInteger
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class WireTest {

    @Test
    fun encodeDecodeBoolTrue() {
        val enc = WireEncoder()
        enc.encodeBool(true)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals(true, (node as NodeScalar).data)
    }

    @Test
    fun encodeDecodeBoolFalse() {
        val enc = WireEncoder()
        enc.encodeBool(false)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals(false, (node as NodeScalar).data)
    }

    @Test
    fun encodeDecodeIntPositive() {
        val enc = WireEncoder()
        enc.encodeInt64(42L)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals(42L, ((node as NodeScalar).data as Number).toLong())
    }

    @Test
    fun encodeDecodeIntNegative() {
        val enc = WireEncoder()
        enc.encodeInt64(-42L)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals(-42L, ((node as NodeScalar).data as Number).toLong())
    }

    @Test
    fun encodeDecodeIntZero() {
        val enc = WireEncoder()
        enc.encodeInt64(0L)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals(0L, ((node as NodeScalar).data as Number).toLong())
    }

    @Test
    fun encodeDecodeIntMax() {
        val tag = Tag.empty().apply { type = ValueType.I64 }
        val node = NodeScalar(data = Long.MAX_VALUE, text = Long.MAX_VALUE.toString(), tag = tag)
        val encoded = Encoder.encodeNode(node)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is NodeScalar)
        assertEquals(Long.MAX_VALUE, (decoded as NodeScalar).data)
    }

    @Test
    fun encodeDecodeIntMin() {
        val tag = Tag.empty().apply { type = ValueType.I64 }
        val node = NodeScalar(data = Long.MIN_VALUE, text = Long.MIN_VALUE.toString(), tag = tag)
        val encoded = Encoder.encodeNode(node)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is NodeScalar)
        assertEquals(Long.MIN_VALUE, (decoded as NodeScalar).data)
    }

    @Test
    fun encodeDecodeUintLarge() {
        val tag = Tag.empty().apply { type = ValueType.U64 }
        val node = NodeScalar(data = 1L shl 50, text = (1L shl 50).toString(), tag = tag)
        val encoded = Encoder.encodeNode(node)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is NodeScalar)
        assertEquals(1L shl 50, (decoded as NodeScalar).data)
    }

    @Test
    fun encodeDecodeStringSimple() {
        val enc = WireEncoder()
        enc.encodeString("hello")
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals("hello", (node as NodeScalar).data)
    }

    @Test
    fun encodeDecodeStringEmpty() {
        val enc = WireEncoder()
        enc.encodeString("")
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals("", (node as NodeScalar).data)
    }

    @Test
    fun encodeDecodeStringUnicode() {
        val enc = WireEncoder()
        enc.encodeString("你好世界")
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertEquals("你好世界", (node as NodeScalar).data)
    }

    @Test
    fun encodeDecodeBytes() {
        val data = byteArrayOf(0x12, 0x34, 0x56, 0x78)
        val enc = WireEncoder()
        enc.encodeBytes(data)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertArrayEquals(data, (node as NodeScalar).data as ByteArray)
    }

    @Test
    fun encodeDecodeBytesEmpty() {
        val data = ByteArray(0)
        val enc = WireEncoder()
        enc.encodeBytes(data)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertArrayEquals(data, (node as NodeScalar).data as ByteArray)
    }

    @Test
    fun encodeDecodeBinaryData() {
        val data = ByteArray(256) { (it % 256).toByte() }
        val enc = WireEncoder()
        enc.encodeBytes(data)
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertArrayEquals(data, (node as NodeScalar).data as ByteArray)
    }

    @Test
    fun encodeDecodeFloatString() {
        val enc = WireEncoder()
        enc.encodeFloatString("3.14")
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        assertTrue((node as NodeScalar).data is Double)
    }

    @Test
    fun encodeDecodeFloatStringZero() {
        val enc = WireEncoder()
        enc.encodeFloatString("0.0")
        val node = Decoder().decode(enc.toByteArray())
        assertTrue(node is NodeScalar)
        val d = (node as NodeScalar).data as Double
        assertEquals(0.0, d, 0.0)
    }

    @Test
    fun encodeDecodeObjectEmpty() {
        val obj = NodeObject()
        val encoded = Encoder.encodeNode(obj)
        val node = Decoder().decode(encoded)
        assertTrue(node is NodeObject)
        assertEquals(0, (node as NodeObject).fields.size)
    }

    @Test
    fun encodeDecodeObjectWithFields() {
        val obj =
                NodeObject().apply {
                    fields.add(
                            Field(
                                    "name",
                                    NodeScalar(
                                            data = "Alice",
                                            text = "\"Alice\"",
                                            tag = Tag.empty().apply { type = ValueType.STR }
                                    )
                            )
                    )
                    fields.add(
                            Field(
                                    "age",
                                    NodeScalar(
                                            data = 30L,
                                            text = "30",
                                            tag = Tag.empty().apply { type = ValueType.I }
                                    )
                            )
                    )
                }
        val encoded = Encoder.encodeNode(obj)
        val node = Decoder().decode(encoded)
        assertTrue(node is NodeObject)
        assertEquals(2, (node as NodeObject).fields.size)
    }

    @Test
    fun encodeDecodeArrayEmpty() {
        val arr = AstArray()
        val encoded = Encoder.encodeNode(arr)
        val node = Decoder().decode(encoded)
        assertTrue(node is AstArray)
        assertEquals(0, (node as AstArray).items.size)
    }

    @Test
    fun encodeDecodeArrayWithItems() {
        val arr =
                AstArray().apply {
                    items.add(
                            NodeScalar(
                                    data = 1L,
                                    text = "1",
                                    tag = Tag.empty().apply { type = ValueType.I }
                            )
                    )
                    items.add(
                            NodeScalar(
                                    data = 2L,
                                    text = "2",
                                    tag = Tag.empty().apply { type = ValueType.I }
                            )
                    )
                    items.add(
                            NodeScalar(
                                    data = 3L,
                                    text = "3",
                                    tag = Tag.empty().apply { type = ValueType.I }
                            )
                    )
                }
        val encoded = Encoder.encodeNode(arr)
        val node = Decoder().decode(encoded)
        assertTrue(node is AstArray)
        assertEquals(3, (node as AstArray).items.size)
    }

    @Test
    fun encodeDecodeSimpleCodeAliases() {
        for (code in SimpleValue.CODE..SimpleValue.VAL) {
            val enc = WireEncoder()
            enc.encodeSimple(code)
            val node = Decoder().decode(enc.toByteArray())
            assertTrue(node is NodeScalar, "Failed for code=$code")
            val name = SimpleValue.toString(code)
            assertEquals(name, (node as NodeScalar).text)
        }
    }

    @Test
    fun encodeDecodeBigIntRoundTrip() {
        val tag = Tag.empty().apply { type = ValueType.BIGINT }
        val bi = BigInteger("12345678901234567890")
        val node = NodeScalar(data = bi, text = bi.toString(), tag = tag)
        val encoded = Encoder.encodeNode(node)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is NodeScalar)
        assertEquals(bi, (decoded as NodeScalar).data)
    }

    @Test
    fun encodeDecodeSmallInts() {
        val testValues = listOf(-128L, -1L, 0L, 1L, 127L, 255L, 256L, 32767L, 65535L)
        for (v in testValues) {
            val enc = WireEncoder()
            enc.encodeInt64(v)
            val node = Decoder().decode(enc.toByteArray())
            assertTrue(node is NodeScalar)
            assertEquals(v, ((node as NodeScalar).data as Number).toLong())
        }
    }

    @Test
    fun encodeDecodeNestedNodeObject() {
        val inner = NodeObject().apply { fields.add(Field("x", NodeScalar(data = 1L, text = "1"))) }
        val outer = NodeObject().apply { fields.add(Field("inner", inner)) }
        val encoded = Encoder.encodeNode(outer)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is NodeObject)
        val outerObj = decoded as NodeObject
        assertEquals(1, outerObj.fields.size)
        assertTrue(outerObj.fields[0].value is NodeObject)
    }
}
