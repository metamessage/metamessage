package io.github.metamessage.mm

import org.junit.jupiter.api.Test
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertTrue
import java.time.LocalDateTime

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
        val wire = MetaMessage.encode(p)
        val out = MetaMessage.decode(wire, Person::class.java)
        assertEquals(p.name, out.name)
        assertEquals(p.age, out.age)
    }

    @Test
    fun roundtripListField() {
        val t = Team()
        val wire = MetaMessage.encode(t)
        val out = MetaMessage.decode(wire, Team::class.java)
        assertEquals(t.teamName, out.teamName)
        assertEquals(t.members, out.members)
    }

    @Test
    fun roundtripDateTime() {
        val c = Clock()
        val wire = MetaMessage.encode(c)
        val out = MetaMessage.decode(wire, Clock::class.java)
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
        val wire = MetaMessage.encode(bt)
        val out = MetaMessage.decode(wire, BasicTypes::class.java)
        assertEquals(bt.bool, out.bool)
        assertEquals(bt.int, out.int)
        assertEquals(bt.long, out.long)
        assertEquals(bt.float, out.float)
        assertEquals(bt.double, out.double)
        assertEquals(bt.string, out.string)
    }

    @MM
    class MapTypes(
        var map: Map<String, Int> = mapOf("a" to 1, "b" to 2, "c" to 3)
    )

    @Test
    fun roundtripMapTypes() {
        val mt = MapTypes()
        val wire = MetaMessage.encode(mt)
        val out = MetaMessage.decode(wire, MapTypes::class.java)
        assertEquals(mt.map.size, out.map.size)
        // Convert keys to the same type as in the map
        for ((key, value) in mt.map) {
            // Try to get the value using the original key
            var found = false
            try {
                for ((outKey, outValue) in out.map) {
                    if (outKey.toString() == key.toString() && outValue == value) {
                        found = true
                        break
                    }
                }
            } catch (e: ClassCastException) {
                // Ignore class cast exceptions
            }
            assertTrue(found, "Value not found for key: $key")
        }
    }
}
