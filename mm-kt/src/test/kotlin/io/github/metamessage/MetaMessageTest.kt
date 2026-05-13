package io.github.metamessage

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.*
import java.time.LocalDateTime
import java.time.LocalDate
import java.time.LocalTime
import java.math.BigInteger
import java.util.UUID

import io.github.metamessage.MM
import io.github.metamessage.ast.ValueType

class MetaMessageTest {

    @MM
    class Person(var name: String = "Ada", var age: Int = 40)

    @MM
    class Team(var teamName: String = "core", @MM(childType = ValueType.STRING) var members: List<String> = listOf("a", "b"))

    @MM
    class Clock(@MM(type = ValueType.DATETIME) var `when`: LocalDateTime = LocalDateTime.of(2024, 6, 1, 12, 0, 0))

    @Test
    fun roundtripSimpleStruct() {
        val p = Person()
        val wire = MetaMessage.encodeFromValue(p)
        val out = MetaMessage.decodeToValue(wire, Person::class.java)
        assertEquals(p.name, out.name)
        assertEquals(p.age, out.age)
    }

    @Test
    fun roundtripListField() {
        val t = Team()
        val wire = MetaMessage.encodeFromValue(t)
        val out = MetaMessage.decodeToValue(wire, Team::class.java)
        assertEquals(t.teamName, out.teamName)
        assertEquals(t.members, out.members)
    }

    @Test
    fun roundtripDateTime() {
        val c = Clock()
        val wire = MetaMessage.encodeFromValue(c)
        val out = MetaMessage.decodeToValue(wire, Clock::class.java)
        assertEquals(c.`when`, out.`when`)
    }

    @MM
    class BasicTypes(
        var bool: Boolean = true,
        var int: Int = 42,
        var long: Long = 1234567890L,
        var float: Float = 3.14f,
        var double: Double = 3.141592653589793,
        var string: String = "Hello, World!"
    )

    @Test
    fun roundtripBasicTypes() {
        val bt = BasicTypes()
        val wire = MetaMessage.encodeFromValue(bt)
        val out = MetaMessage.decodeToValue(wire, BasicTypes::class.java)
        assertEquals(bt.bool, out.bool)
        assertEquals(bt.int, out.int)
        assertEquals(bt.long, out.long)
        assertEquals(bt.float, out.float)
        assertEquals(bt.double, out.double)
        assertEquals(bt.string, out.string)
    }

    @MM(desc = "A map of strings to integers")
    class MapTypes(
        var map: Map<String, Int> = mapOf("a" to 1, "b" to 2, "c" to 3)
    )

    @Test
    fun roundtripMapTypes() {
        val mt = MapTypes()
        val wire = MetaMessage.encodeFromValue(mt)
        val out = MetaMessage.decodeToValue(wire, MapTypes::class.java)
        assertEquals(mt.map.size, out.map.size)
        for ((key, value) in mt.map) {
            var found = false
            try {
                for ((outKey, outValue) in out.map) {
                    if (outKey.toString() == key.toString() && outValue == value) {
                        found = true
                        break
                    }
                }
            } catch (e: ClassCastException) {
                println("捕获到异常: ${e.message}")
                e.printStackTrace()
            }
            assertTrue(found, "Value not found for key: $key")
        }
    }

    @MM
    class AllIntegerTypes(
        @MM(type = ValueType.INT8) var int8: Byte = 10,
        @MM(type = ValueType.INT16) var int16: Short = 100,
        @MM(type = ValueType.INT32) var int32: Int = 1000,
        @MM(type = ValueType.INT64) var int64: Long = 10000L,
        @MM(type = ValueType.UINT8) var uint8: Byte = 10,
        @MM(type = ValueType.UINT16) var uint16: Short = 100,
        @MM(type = ValueType.UINT32) var uint32: Int = 1000,
        @MM(type = ValueType.UINT64) var uint64: Long = 10000L
    )

    @Test
    fun roundtripAllIntegerTypes() {
        val obj = AllIntegerTypes()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, AllIntegerTypes::class.java)
        assertEquals(obj.int8, out.int8)
        assertEquals(obj.int16, out.int16)
        assertEquals(obj.int32, out.int32)
        assertEquals(obj.int64, out.int64)
        assertEquals(obj.uint8, out.uint8)
        assertEquals(obj.uint16, out.uint16)
        assertEquals(obj.uint32, out.uint32)
        assertEquals(obj.uint64, out.uint64)
    }

    @MM
    class FloatTypes(
        @MM(type = ValueType.FLOAT32) var float32: Float = 3.14f,
        @MM(type = ValueType.FLOAT64) var float64: Double = 3.141592653589793
    )

    @Test
    fun roundtripFloatTypes() {
        val obj = FloatTypes()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, FloatTypes::class.java)
        assertEquals(obj.float32, out.float32, 0.0001f)
        assertEquals(obj.float64, out.float64, 0.0000001)
    }

    @MM
    class StringTypes(
        @MM(type = ValueType.STRING) var str: String = "test",
        @MM(type = ValueType.URL) var url: String = "https://example.com",
        @MM(type = ValueType.EMAIL) var email: String = "test@example.com",
        @MM(type = ValueType.IP) var ip: String = "192.168.1.1",
        @MM(type = ValueType.UUID) var uuid: String = "550e8400-e29b-41d4-a716-446655440000"
    )

    @Test
    fun roundtripStringTypes() {
        val obj = StringTypes()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, StringTypes::class.java)
        assertEquals(obj.str, out.str)
        assertEquals(obj.url, out.url)
        assertEquals(obj.email, out.email)
        assertEquals(obj.ip, out.ip)
        assertEquals(obj.uuid, out.uuid)
    }

    @MM
    class DateTypes(
        @MM(type = ValueType.DATETIME) var datetime: LocalDateTime = LocalDateTime.of(2024, 12, 31, 23, 59, 59),
        @MM(type = ValueType.DATE) var date: LocalDate = LocalDate.of(2024, 12, 31),
        @MM(type = ValueType.TIME) var time: LocalTime = LocalTime.of(23, 59, 59)
    )

    @Test
    fun roundtripDateTypes() {
        val obj = DateTypes()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, DateTypes::class.java)
        assertEquals(obj.datetime, out.datetime)
        assertEquals(obj.date, out.date)
        assertEquals(obj.time, out.time)
    }

    @MM
    class BytesType(
        @MM(type = ValueType.BYTES) var bytes: ByteArray = "Hello".toByteArray()
    )

    @Test
    fun roundtripBytes() {
        val obj = BytesType()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, BytesType::class.java)
        assertArrayEquals(obj.bytes, out.bytes)
    }

    @MM
    class BigIntType(
        @MM(type = ValueType.BIGINT) var bigint: BigInteger = BigInteger("123456789012345678901234567890")
    )

    @Test
    fun roundtripBigInt() {
        val obj = BigIntType()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, BigIntType::class.java)
        assertEquals(obj.bigint, out.bigint)
    }

    @MM
    class BoolEdgeCases(
        var boolTrue: Boolean = true,
        var boolFalse: Boolean = false
    )

    @Test
    fun roundtripBoolEdgeCases() {
        val obj = BoolEdgeCases()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, BoolEdgeCases::class.java)
        assertTrue(out.boolTrue)
        assertFalse(out.boolFalse)
    }

    @MM
    class SliceTypes(
        @MM(type = ValueType.SLICE, childType = ValueType.INT) var intSlice: List<Int> = listOf(1, 2, 3),
        @MM(type = ValueType.SLICE, childType = ValueType.STRING) var strSlice: List<String> = listOf("a", "b", "c")
    )

    @Test
    fun roundtripSlices() {
        val obj = SliceTypes()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, SliceTypes::class.java)
        assertEquals(obj.intSlice, out.intSlice)
        assertEquals(obj.strSlice, out.strSlice)
    }

    @MM
    class NestedStruct(
        var name: String = "parent",
        var child: ChildStruct = ChildStruct()
    )

    @MM
    class ChildStruct(
        var value: Int = 42
    )

    @Test
    fun roundtripNestedStruct() {
        val obj = NestedStruct()
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, NestedStruct::class.java)
        assertEquals(obj.name, out.name)
        assertEquals(obj.child.value, out.child.value)
    }

    @MM
    class NullableTypes(
        var nullableString: String? = "test",
        var nullableInt: Int? = 42,
        var nullableBool: Boolean? = true
    )

    @Test
    fun roundtripNullableWithValues() {
        val obj = NullableTypes("test", 42, true)
        val wire = MetaMessage.encodeFromValue(obj)
        val out = MetaMessage.decodeToValue(wire, NullableTypes::class.java)
        assertEquals("test", out.nullableString)
        assertEquals(42, out.nullableInt)
        assertTrue(out.nullableBool ?: false)
    }
}