package io.metamessage.mm;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class MetaMessageTest {

    @MM
    static class Person {
        public String name = "Ada";
        public int age = 40;
    }

    @MM
    static class Team {
        public String teamName = "core";

        @MM(childType = ValueType.STRING)
        public List<String> members = List.of("a", "b");
    }

    @Test
    void roundtripSimpleStruct() throws Exception {
        Person p = new Person();
        byte[] wire = MetaMessage.encode(p);
        Person out = MetaMessage.decode(wire, Person.class);
        assertEquals(p.name, out.name);
        assertEquals(p.age, out.age);
    }

    @Test
    void roundtripListField() throws Exception {
        Team t = new Team();
        byte[] wire = MetaMessage.encode(t);
        Team out = MetaMessage.decode(wire, Team.class);
        assertEquals(t.teamName, out.teamName);
        assertEquals(t.members, out.members);
    }

    @MM
    static class Clock {
        @MM(type = ValueType.DATETIME)
        public LocalDateTime when = LocalDateTime.of(2024, 6, 1, 12, 0, 0);
    }

    @Test
    void roundtripDateTime() throws Exception {
        Clock c = new Clock();
        byte[] wire = MetaMessage.encode(c);
        Clock out = MetaMessage.decode(wire, Clock.class);
        assertEquals(c.when, out.when);
    }
    
    @MM
    static class BasicTypes {
        public boolean boolValue = true;
        public byte byteValue = 127;
        public short shortValue = 32767;
        public int intValue = 2147483647;
        public long longValue = 9223372036854775807L;
        public float floatValue = 3.14f;
        public double doubleValue = 3.141592653589793;
        public String stringValue = "Hello, World!";
    }
    
    @Test
    void roundtripBasicTypes() throws Exception {
        BasicTypes bt = new BasicTypes();
        byte[] wire = MetaMessage.encode(bt);
        BasicTypes out = MetaMessage.decode(wire, BasicTypes.class);
        assertEquals(bt.boolValue, out.boolValue);
        assertEquals(bt.byteValue, out.byteValue);
        assertEquals(bt.shortValue, out.shortValue);
        assertEquals(bt.intValue, out.intValue);
        assertEquals(bt.longValue, out.longValue);
        assertEquals(bt.floatValue, out.floatValue, 0.0001f);
        assertEquals(bt.doubleValue, out.doubleValue, 0.0001);
        assertEquals(bt.stringValue, out.stringValue);
    }
    
    @MM
    static class MapTypes {
        public Map<String, String> stringMap = Map.of("key1", "value1", "key2", "value2");
        public Map<String, Integer> intMap = Map.of("key1", 1, "key2", 2);
    }
    
    @Test
    void roundtripMapTypes() throws Exception {
        MapTypes mt = new MapTypes();
        byte[] wire = MetaMessage.encode(mt);
        MapTypes out = MetaMessage.decode(wire, MapTypes.class);
        assertEquals(mt.stringMap, out.stringMap);
        assertEquals(mt.intMap, out.intMap);
    }
}
