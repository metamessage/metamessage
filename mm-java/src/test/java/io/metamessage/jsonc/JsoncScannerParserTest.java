package io.metamessage.jsonc;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import io.metamessage.mm.ValueType;
import java.util.List;
import org.junit.jupiter.api.Test;

class JsoncScannerParserTest {

    @Test
    void scannerDistinguishesLeadingAndTrailingComments() {
        String src = """
                // mm: type=struct
                {
                    // mm: type=i
                    "Code": 1, // mm: desc=ok
                    /* mm: type=str */
                    "Name": "A\nB"
                }
                """;

        List<JsoncToken> toks = new JsoncScanner(src).tokenizeAll();
        assertFalse(toks.isEmpty());
        assertTrue(toks.stream().anyMatch(t -> t.type() == JsoncTokenType.LEADING_COMMENT));
        assertTrue(toks.stream().anyMatch(t -> t.type() == JsoncTokenType.TRAILING_COMMENT));
        assertEquals(JsoncTokenType.EOF, toks.get(toks.size() - 1).type());
    }

    @Test
    void parserParsesObjectAndSnakeCaseKeys() {
        String src = """
                {
                  "UserName": "liz",
                  "IsAdmin": true,
                  "Score": 12.5
                }
                """;

        JcNode n = Jsonc.parseFromString(src);
        JcNode.JcObject o = assertInstanceOf(JcNode.JcObject.class, n);
        assertEquals(3, o.fields().size());
        assertEquals("user_name", o.fields().get(0).key());
        assertEquals("is_admin", o.fields().get(1).key());
        assertEquals("score", o.fields().get(2).key());
    }

    @Test
    void parserParsesArrayWithLeadingTag() {
        String src = """
                // mm: type=slice
                [1, 2, 3] // mm: desc=nums
                """;
        JcNode n = Jsonc.parseFromString(src);
        JcNode.JcArray arr = assertInstanceOf(JcNode.JcArray.class, n);
        assertNotNull(arr.tag());
        assertEquals(ValueType.SLICE, arr.tag().type);
        assertEquals(3, arr.items().size());
    }

    @Test
    void parserRejectsNullLiteral() {
        JsoncException ex = assertThrows(JsoncException.class, () -> Jsonc.parseFromString("{\"n\":null}"));
        assertTrue(ex.getMessage().contains("null literal"));
    }

    @Test
    void parserHonorsMaxDepth() {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < 40; i++) {
            sb.append('[');
        }
        for (int i = 0; i < 40; i++) {
            sb.append(']');
        }
        JsoncException ex = assertThrows(JsoncException.class, () -> Jsonc.parseFromString(sb.toString()));
        assertTrue(ex.getMessage().contains("max depth"));
    }

    @Test
    void printerOutputsCommentsAndEscapedStrings() {
        String raw = "a\nb\"c";
        JcNode.JcValue val = new JcNode.JcValue(raw, raw, tag(ValueType.STRING), "p");
        JcNode.JcObject obj = new JcNode.JcObject(tag(ValueType.STRUCT), "", List.of(new JcField("name", val)));

        String out = Jsonc.toString(obj);
        assertTrue(out.contains("// mm:"), out);
        assertTrue(out.contains("\\n"), out);
        assertTrue(out.contains("\\\""), out);
        assertTrue(out.contains("\"name\""), out);
    }

    private static io.metamessage.mm.MmTag tag(ValueType type) {
        io.metamessage.mm.MmTag t = io.metamessage.mm.MmTag.empty();
        t.type = type;
        t.desc = "d";
        return t;
    }
}
