package io.github.metamessage.jsonc;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.github.metamessage.mm.MmTag;
import io.github.metamessage.mm.ValueType;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;

class JsoncBinderTest {

    public static class Child {
        public String Name;
    }

    public static class Person {
        public int Code;
        public String Name;
        public Child Child;
        public List<String> Tags;
        public Map<String, Object> Meta;
    }

    @Test
    void bindStructWithNestedListAndMap() {
        JcNode.JcObject node = new JcNode.JcObject(
                tag(ValueType.STRUCT),
                "",
                List.of(
                        new JcField("code", new JcNode.JcValue(7, "7", tag(ValueType.INT), "p.code")),
                        new JcField("name", new JcNode.JcValue("liz", "liz", tag(ValueType.STRING), "p.name")),
                        new JcField("child", new JcNode.JcObject(tag(ValueType.STRUCT), "p.child", List.of(
                                new JcField("name", new JcNode.JcValue("c", "c", tag(ValueType.STRING), "p.child.name"))))),
                        new JcField("tags", new JcNode.JcArray(tag(ValueType.SLICE), "p.tags", List.of(
                                new JcNode.JcValue("a", "a", tag(ValueType.STRING), ""),
                                new JcNode.JcValue("b", "b", tag(ValueType.STRING), "")))),
                        new JcField("meta", new JcNode.JcObject(tag(ValueType.MAP), "p.meta", List.of(
                                new JcField("x", new JcNode.JcValue(1, "1", tag(ValueType.INT), "")))))));

        Person out = new Person();
        Jsonc.bind(node, out);

        assertEquals(7, out.Code);
        assertEquals("liz", out.Name);
        assertNotNull(out.Child);
        assertEquals("c", out.Child.Name);
        assertEquals(List.of("a", "b"), out.Tags);
        assertEquals(1, ((Number) out.Meta.get("x")).intValue());
    }

    @Test
    void bindArrayIntoJavaArray() {
        JcNode.JcArray arr = new JcNode.JcArray(
                tag(ValueType.ARRAY, 3),
                "arr",
                List.of(
                        new JcNode.JcValue(1, "1", tag(ValueType.INT), "arr[0]"),
                        new JcNode.JcValue(2, "2", tag(ValueType.INT), "arr[1]"),
                        new JcNode.JcValue(3, "3", tag(ValueType.INT), "arr[2]")));

        int[] out = new int[3];
        Jsonc.bind(arr, out);
        assertEquals(1, out[0]);
        assertEquals(2, out[1]);
        assertEquals(3, out[2]);
    }

    @Test
    void bindFailureCases() {
        assertThrows(JsoncException.class, () -> Jsonc.bind(null, new Object()));

        JcNode.JcObject mapNode = new JcNode.JcObject(tag(ValueType.MAP), "", List.of());
        assertThrows(JsoncException.class, () -> Jsonc.bind(mapNode, new Object()));

        JcNode.JcArray arr = new JcNode.JcArray(tag(ValueType.ARRAY, 2), "", List.of(
                new JcNode.JcValue(1, "1", tag(ValueType.INT), ""),
                new JcNode.JcValue(2, "2", tag(ValueType.INT), "")));
        int[] wrong = new int[1];
        JsoncException ex = assertThrows(JsoncException.class, () -> Jsonc.bind(arr, wrong));
        assertTrue(ex.getMessage().contains("array length mismatch") || ex.getMessage().contains("item count"));
    }

    @Test
    void bindFromStringFactoryPath() {
        String src = """
                {
                  "Code": 9,
                  "Name": "neo"
                }
                """;
        Person p = Jsonc.bindFromString(src, Person.class);
        assertEquals(9, p.Code);
        assertEquals("neo", p.Name);
    }

    @Test
    void bindArrayToListTarget() {
        JcNode.JcArray arr = new JcNode.JcArray(tag(ValueType.SLICE), "", List.of(
                new JcNode.JcValue("a", "a", tag(ValueType.STRING), ""),
                new JcNode.JcValue("b", "b", tag(ValueType.STRING), "")));
        List<Object> out = new ArrayList<>();
        Jsonc.bind(arr, out);
        assertEquals(List.of("a", "b"), out);
    }

    private static MmTag tag(ValueType type) {
        MmTag t = MmTag.empty();
        t.type = type;
        return t;
    }

    private static MmTag tag(ValueType type, int size) {
        MmTag t = tag(type);
        t.size = size;
        return t;
    }
}
