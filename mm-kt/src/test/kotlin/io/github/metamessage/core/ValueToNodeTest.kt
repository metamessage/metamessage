package io.github.metamessage.core

import io.github.metamessage.MM
import io.github.metamessage.ir.*
import io.github.metamessage.ir.Array as AstArray
import java.math.BigInteger
import java.time.LocalDate
import java.time.LocalDateTime
import java.time.LocalTime
import java.util.UUID
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class ValueToNodeTest {

    @MM class SimpleUser(var name: String = "Alice", var age: Int = 30)

    @Test
    fun simpleStructEncodeDecode() {
        val user = SimpleUser()
        val encoded = Encoder.encode(user)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is Object)
        val obj = decoded as Object
        assertEquals(2, obj.fields.size)
        val nameField = obj.fields.find { it.key == "name" }
        assertNotNull(nameField)
        assertEquals("Alice", (nameField!!.value as Value).data)
        val ageField = obj.fields.find { it.key == "age" }
        assertNotNull(ageField)
        assertEquals(30L, ((ageField!!.value as Value).data as Number).toLong())
    }

    @Test
    fun nullToNode() {
        val node = valueToNode(null, Tag.empty().apply { type = ValueType.I }, "")
        assertTrue(node is Value)
        assertNull((node as Value).data)
    }

    @Test
    fun stringToNode() {
        val node = valueToNode("hello", Tag.empty().apply { type = ValueType.STR }, "")
        assertTrue(node is Value)
        assertEquals("hello", (node as Value).data)
    }

    @Test
    fun intToNode() {
        val node = valueToNode(42, Tag.empty().apply { type = ValueType.I }, "")
        assertTrue(node is Value)
        assertEquals(42L, ((node as Value).data as Number).toLong())
    }

    @Test
    fun longToNode() {
        val node = valueToNode(1234567890L, Tag.empty().apply { type = ValueType.I64 }, "")
        assertTrue(node is Value)
        assertEquals(1234567890L, (node as Value).data)
    }

    @Test
    fun booleanTrueToNode() {
        val node = valueToNode(true, Tag.empty().apply { type = ValueType.BOOL }, "")
        assertTrue(node is Value)
        assertEquals(true, (node as Value).data)
    }

    @Test
    fun booleanFalseToNode() {
        val node = valueToNode(false, Tag.empty().apply { type = ValueType.BOOL }, "")
        assertTrue(node is Value)
        assertEquals(false, (node as Value).data)
    }

    @Test
    fun floatToNode() {
        val node = valueToNode(3.14f, Tag.empty().apply { type = ValueType.F32 }, "")
        assertTrue(node is Value)
        assertTrue((node as Value).data is Float)
    }

    @Test
    fun doubleToNode() {
        val node = valueToNode(3.141592653589793, Tag.empty().apply { type = ValueType.F64 }, "")
        assertTrue(node is Value)
        assertTrue((node as Value).data is Double)
        assertEquals(3.141592653589793, (node as Value).data as Double, 0.0000001)
    }

    @Test
    fun byteArrayToNode() {
        val data = byteArrayOf(0x12, 0x34, 0x56)
        val node = valueToNode(data, Tag.empty().apply { type = ValueType.BYTES }, "")
        assertTrue(node is Value)
        assertArrayEquals(data, (node as Value).data as ByteArray)
    }

    @Test
    fun localDateTimeToNode() {
        val dt = LocalDateTime.of(2024, 6, 1, 12, 0, 0)
        val node = valueToNode(dt, Tag.empty().apply { type = ValueType.DATETIME }, "")
        assertTrue(node is Value)
        assertEquals(dt, (node as Value).data)
    }

    @Test
    fun localDateToNode() {
        val d = LocalDate.of(2024, 12, 31)
        val node = valueToNode(d, Tag.empty().apply { type = ValueType.DATE }, "")
        assertTrue(node is Value)
        assertEquals(d, (node as Value).data)
    }

    @Test
    fun localTimeToNode() {
        val t = LocalTime.of(23, 59, 59)
        val node = valueToNode(t, Tag.empty().apply { type = ValueType.TIME }, "")
        assertTrue(node is Value)
        assertEquals(t, (node as Value).data)
    }

    @Test
    fun bigIntegerToNode() {
        val bi = BigInteger("123456789012345678901234567890")
        val node = valueToNode(bi, Tag.empty().apply { type = ValueType.BIGINT }, "")
        assertTrue(node is Value)
        assertEquals(bi, (node as Value).data)
    }

    @Test
    fun uuidToNode() {
        val uuid = UUID.randomUUID()
        val node = valueToNode(uuid, Tag.empty().apply { type = ValueType.UUID }, "")
        assertTrue(node is Value)
        assertTrue((node as Value).data is ByteArray)
    }

    @Test
    fun byteToNode() {
        val node = valueToNode(42.toByte(), Tag.empty().apply { type = ValueType.I8 }, "")
        assertTrue(node is Value)
        assertEquals(42L, ((node as Value).data as Number).toLong())
    }

    @Test
    fun shortToNode() {
        val node = valueToNode(1234.toShort(), Tag.empty().apply { type = ValueType.I16 }, "")
        assertTrue(node is Value)
        assertEquals(1234L, ((node as Value).data as Number).toLong())
    }

    @Test
    fun listToNode() {
        val list = listOf(1, 2, 3)
        val tag =
                Tag.empty().apply {
                    type = ValueType.VEC
                    childType = ValueType.I
                }
        val node = valueToNode(list, tag, "")
        assertTrue(node is AstArray)
        assertEquals(3, (node as AstArray).items.size)
    }

    @Test
    fun emptyListToNode() {
        val list = emptyList<Int>()
        val tag =
                Tag.empty().apply {
                    type = ValueType.VEC
                    childType = ValueType.I
                }
        val node = valueToNode(list, tag, "")
        assertTrue(node is AstArray)
        assertEquals(1, (node as AstArray).items.size)
    }

    @Test
    fun mapToNode() {
        val map = mapOf("a" to 1, "b" to 2)
        val tag = Tag.empty().apply { type = ValueType.MAP }
        val node = valueToNode(map, tag, "")
        assertTrue(node is Object)
    }

    @Test
    fun structWithNullableFieldsToNode() {
        @MM class Nullables(var name: String? = "test", var age: Int? = 42)
        val obj = Nullables()
        val encoded = Encoder.encode(obj)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is Object)
        val o = decoded as Object
        assertTrue(o.fields.size >= 2)
    }

    @Test
    fun structWithDateTimeToNode() {
        @MM
        class Event(
                @MM(type = ValueType.DATETIME)
                var when_: LocalDateTime = LocalDateTime.of(2024, 6, 1, 12, 0, 0)
        )
        val event = Event()
        val encoded = Encoder.encode(event)
        val decoded = Decoder().decode(encoded)
        assertTrue(decoded is Object)
        val fields = (decoded as Object).fields
        assertTrue(fields.isNotEmpty())
    }

    @Test
    fun nilToNodeReturnsCorrectValue() {
        val values =
                listOf(
                        ValueType.I to null,
                        ValueType.I8 to null,
                        ValueType.STR to null,
                        ValueType.BOOL to null,
                        ValueType.F64 to null,
                        ValueType.BYTES to null
                )
        for ((vt, _) in values) {
            val node = nilToNode(vt)
            assertTrue(node is Value)
            assertNull((node as Value).data)
            assertEquals(vt, node.tag?.type)
        }
    }
}
