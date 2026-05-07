package io.github.metamessage.mm;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import org.junit.jupiter.api.Test;

class MmLineParserTest {

    @Test
    void parseAndFormatComplexDirectives() {
        String line = "// mm: type=u16; example; desc=\"user id\"; nullable; allow_empty; unique; default=1; min=1; max=9; size=3; enum=a|b; pattern=^[a-z]+$; location=8; version=2; mime=json; child_desc=\"child d\"; child_type=i32; child_raw; child_nullable; child_allow_empty; child_unique; child_default=7; child_min=1; child_max=8; child_size=9; child_enum=x|y; child_pattern=^[0-9]+$; child_location=9; child_version=3; child_mime=png";

        MmTag t = MmLineParser.parse(line);
        assertEquals(ValueType.ENUM, t.type);
        assertTrue(t.example);
        assertEquals("user id", t.desc);
        assertTrue(t.nullable);
        assertTrue(t.allowEmpty);
        assertTrue(t.unique);
        assertEquals("1", t.defaultValue);
        assertEquals("1", t.min);
        assertEquals("9", t.max);
        assertEquals(3, t.size);
        assertEquals("a|b", t.enumValues);
        assertEquals("^[a-z]+$", t.pattern);
        assertEquals(8, t.locationHours);
        assertEquals(2, t.version);
        assertEquals("json", t.mime);

        assertEquals("child d", t.childDesc);
        assertEquals(ValueType.ENUM, t.childType);
        assertTrue(t.childRaw);
        assertTrue(t.childNullable);
        assertTrue(t.childAllowEmpty);
        assertTrue(t.childUnique);
        assertEquals("7", t.childDefault);
        assertEquals("1", t.childMin);
        assertEquals("8", t.childMax);
        assertEquals(9, t.childSize);
        assertEquals("x|y", t.childEnum);
        assertEquals("^[0-9]+$", t.childPattern);
        assertEquals(9, t.childLocationHours);
        assertEquals(3, t.childVersion);
        assertEquals("png", t.childMime);

        String out = MmTagDirectives.format(t);
        assertTrue(out.contains("enum=a|b"), out);
        assertTrue(out.contains("child_enum=x|y"), out);
        assertTrue(out.contains("desc=\"user id\""), out);
    }

    @Test
    void parseIgnoresUnknownAndInvalidVersion() {
        MmTag t = MmLineParser.parse("mm: foo=bar; version=100; type=zzz; child_version=100; child_type=bbb; is_null");
        assertTrue(t.isNull);
        assertEquals(ValueType.UNKNOWN, t.type);
        assertEquals(MmTag.DEFAULT_VERSION, t.version);
        assertEquals(ValueType.UNKNOWN, t.childType);
        assertEquals(MmTag.DEFAULT_VERSION, t.childVersion);
    }
}
