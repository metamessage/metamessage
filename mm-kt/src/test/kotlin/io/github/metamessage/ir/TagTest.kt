package io.github.metamessage.ir

import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class TagTest {

    @Test
    fun emptyTag() {
        val t = Tag.empty()
        assertEquals("", t.name)
        assertEquals(ValueType.UNKNOWN, t.type)
        assertEquals("", t.desc)
        assertFalse(t.isNull)
        assertFalse(t.nullable)
        assertEquals(0, t.version)
    }

    @Test
    fun copyCreatesIndependentTag() {
        val t = Tag(name = "test", type = ValueType.I, desc = "a number")
        val copy = t.copy()
        assertEquals(t.name, copy.name)
        assertEquals(t.type, copy.type)
        assertEquals(t.desc, copy.desc)

        copy.desc = "changed"
        assertEquals("a number", t.desc)
        assertEquals("changed", copy.desc)
    }

    @Test
    fun parseMMTagBasicType() {
        val t = Tag.parseMMTag("mm: type=str")
        assertEquals(ValueType.STR, t.type)
    }

    @Test
    fun parseMMTagWithDesc() {
        val t = Tag.parseMMTag("mm: desc=用户名")
        assertEquals("用户名", t.desc)
    }

    @Test
    fun parseMMTagWithName() {
        val t = Tag.parseMMTag("mm: name=userId")
        assertEquals("userId", t.name)
    }

    @Test
    fun parseMMTagWithNullable() {
        val t = Tag.parseMMTag("mm: nullable")
        assertTrue(t.nullable)
    }

    @Test
    fun parseMMTagWithIsNull() {
        val t = Tag.parseMMTag("mm: is_null")
        assertTrue(t.isNull)
        assertTrue(t.nullable)
    }

    @Test
    fun parseMMTagWithExample() {
        val t = Tag.parseMMTag("mm: example")
        assertTrue(t.example)
    }

    @Test
    fun parseMMTagWithAllowEmpty() {
        val t = Tag.parseMMTag("mm: allow_empty")
        assertTrue(t.allowEmpty)
    }

    @Test
    fun parseMMTagWithUnique() {
        val t = Tag.parseMMTag("mm: unique")
        assertTrue(t.unique)
    }

    @Test
    fun parseMMTagWithDefault() {
        val t = Tag.parseMMTag("mm: default=42")
        assertEquals("42", t.default)
    }

    @Test
    fun parseMMTagWithMin() {
        val t = Tag.parseMMTag("mm: min=10")
        assertEquals("10", t.min)
    }

    @Test
    fun parseMMTagWithMax() {
        val t = Tag.parseMMTag("mm: max=100")
        assertEquals("100", t.max)
    }

    @Test
    fun parseMMTagWithSize() {
        val t = Tag.parseMMTag("mm: size=5")
        assertEquals(5, t.size)
    }

    @Test
    fun parseMMTagWithEnum() {
        val t = Tag.parseMMTag("mm: enum=a|b|c")
        assertEquals(ValueType.ENUM, t.type)
        assertEquals("a|b|c", t.enum)
    }

    @Test
    fun parseMMTagWithPattern() {
        val t = Tag.parseMMTag("mm: pattern=^[0-9]+$")
        assertEquals("^[0-9]+$", t.pattern)
    }

    @Test
    fun parseMMTagWithLocation() {
        val t = Tag.parseMMTag("mm: location=8")
        assertEquals(8, t.location)
    }

    @Test
    fun parseMMTagWithVersion() {
        val t = Tag.parseMMTag("mm: version=5")
        assertEquals(5, t.version)
    }

    @Test
    fun parseMMTagWithMime() {
        val t = Tag.parseMMTag("mm: mime=application/json")
        assertEquals("application/json", t.mime)
    }

    @Test
    fun parseMMTagWithChildDesc() {
        val t = Tag.parseMMTag("mm: child_desc=元素描述")
        assertEquals("元素描述", t.childDesc)
    }

    @Test
    fun parseMMTagWithChildType() {
        val t = Tag.parseMMTag("mm: child_type=i32")
        assertEquals(ValueType.I32, t.childType)
    }

    @Test
    fun parseMMTagWithChildAllowEmpty() {
        val t = Tag.parseMMTag("mm: child_allow_empty")
        assertTrue(t.childAllowEmpty)
    }

    @Test
    fun parseMMTagWithChildUnique() {
        val t = Tag.parseMMTag("mm: child_unique")
        assertTrue(t.childUnique)
    }

    @Test
    fun parseMMTagWithQuotedValue() {
        val t = Tag.parseMMTag("mm: name=\"myField\"")
        assertEquals("myField", t.name)
    }

    @Test
    fun parseMMTagMultipleFields() {
        val t = Tag.parseMMTag("mm: type=str; desc=用户名; nullable; min=1; max=100")
        assertEquals(ValueType.STR, t.type)
        assertEquals("用户名", t.desc)
        assertTrue(t.nullable)
        assertEquals("1", t.min)
        assertEquals("100", t.max)
    }

    @Test
    fun parseMMTagFromComment() {
        val t = Tag.parseMMTag("// mm: type=i; desc=年龄")
        assertEquals(ValueType.I, t.type)
        assertEquals("年龄", t.desc)
    }

    @Test
    fun parseMMTagEmpty() {
        val t = Tag.parseMMTag("")
        assertEquals(ValueType.UNKNOWN, t.type)
        assertEquals("", t.desc)
    }

    @Test
    fun mergeTagCopiesType() {
        val dst = Tag()
        val src = Tag().apply { type = ValueType.I32 }
        Tag.mergeTag(dst, src)
        assertEquals(ValueType.I32, dst.type)
    }

    @Test
    fun mergeTagCopiesDesc() {
        val dst = Tag()
        val src = Tag().apply { desc = "my field" }
        Tag.mergeTag(dst, src)
        assertEquals("my field", dst.desc)
    }

    @Test
    fun mergeTagCopiesNullable() {
        val dst = Tag()
        val src = Tag().apply { nullable = true }
        Tag.mergeTag(dst, src)
        assertTrue(dst.nullable)
    }

    @Test
    fun mergeTagNullSourceDoesNothing() {
        val dst = Tag().apply { type = ValueType.I }
        Tag.mergeTag(dst, null)
        assertEquals(ValueType.I, dst.type)
    }

    @Test
    fun mergeTagDoesNotOverwriteExisting() {
        val dst = Tag().apply { desc = "original" }
        val src = Tag().apply { type = ValueType.I32 }
        Tag.mergeTag(dst, src)
        assertEquals("original", dst.desc)
        assertEquals(ValueType.I32, dst.type)
    }

    @Test
    fun mergeTagCopiesChildFields() {
        val dst = Tag()
        val src =
                Tag().apply {
                    childDesc = "child desc"
                    childType = ValueType.I32
                    childAllowEmpty = true
                }
        Tag.mergeTag(dst, src)
        assertEquals("child desc", dst.childDesc)
        assertEquals(ValueType.I32, dst.childType)
        assertTrue(dst.childAllowEmpty)
    }

    @Test
    fun inheritFromArrayParentCopiesChildFields() {
        val parent =
                Tag().apply {
                    childDesc = "element description"
                    childType = ValueType.I32
                    childMin = "1"
                    childMax = "100"
                }
        val child = Tag()
        child.inheritFromArrayParent(parent)

        assertEquals("element description", child.desc)
        assertEquals(ValueType.I32, child.type)
        assertEquals("1", child.min)
        assertEquals("100", child.max)
        assertTrue(child.isInherit)
    }

    @Test
    fun inheritFromArrayParentNullDoesNothing() {
        val child = Tag().apply { type = ValueType.STR }
        child.inheritFromArrayParent(null)
        assertEquals(ValueType.STR, child.type)
        assertFalse(child.isInherit)
    }

    @Test
    fun toStringBasicType() {
        val t = Tag().apply { type = ValueType.UUID }
        val s = t.toString()
        assertTrue(s.contains("type=uuid"))
    }

    @Test
    fun toStringWithDesc() {
        val t = Tag().apply { desc = "hello world" }
        val s = t.toString()
        assertTrue(s.contains("desc=\"hello world\""))
    }

    @Test
    fun toStringWithDescNeedsQuoting() {
        val t = Tag().apply { desc = "a \"quoted\" string" }
        val s = t.toString()
        assertTrue(s.contains("desc="))
        assertTrue(s.contains("quoted"))
    }

    @Test
    fun toStringWithMultipleFields() {
        val t =
                Tag().apply {
                    type = ValueType.I
                    desc = "a number"
                    min = "1"
                    max = "100"
                    nullable = true
                }
        val s = t.toString()
        assertTrue(s.contains("desc="))
        assertTrue(s.contains("min=1"))
        assertTrue(s.contains("max=100"))
        assertTrue(s.contains("nullable"))
    }

    @Test
    fun toStringEmptyTag() {
        val t = Tag()
        assertEquals("", t.toString())
    }

    @Test
    fun toStringWithChildFields() {
        val t =
                Tag().apply {
                    type = ValueType.VEC
                    childType = ValueType.I32
                    childDesc = "elements"
                }
        val s = t.toString()
        assertTrue(s.contains("child_type=i32"))
        assertTrue(s.contains("child_desc=\"elements\""))
    }

    @Test
    fun toBytesEmptyTag() {
        val t = Tag()
        assertEquals(0, t.toBytes().size)
    }

    @Test
    fun toBytesWithType() {
        val t = Tag().apply { type = ValueType.UUID }
        val bytes = t.toBytes()
        assertTrue(bytes.isNotEmpty())
    }

    @Test
    fun toBytesWithDesc() {
        val t = Tag().apply { desc = "test" }
        val bytes = t.toBytes()
        assertTrue(bytes.isNotEmpty())
    }

    @Test
    fun toBytesWithChildFields() {
        val t =
                Tag().apply {
                    childType = ValueType.I32
                    childDesc = "child"
                }
        val bytes = t.toBytes()
        assertTrue(bytes.isNotEmpty())
    }

    @Test
    fun getPatternReturnsNullForEmpty() {
        val t = Tag()
        assertNull(t.getPattern())
    }

    @Test
    fun getPatternReturnsCompiledPattern() {
        val t = Tag().apply { pattern = "^[a-z]+$" }
        val p = t.getPattern()
        assertNotNull(p)
        assertTrue(p!!.matcher("hello").matches())
        assertFalse(p.matcher("Hello123").matches())
    }
}
