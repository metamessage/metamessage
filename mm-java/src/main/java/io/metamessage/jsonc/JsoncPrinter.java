package io.github.metamessage.jsonc;

import io.github.metamessage.mm.MmTag;
import io.github.metamessage.mm.MmTagDirectives;
import io.github.metamessage.mm.ValueType;

/**
 * Renders a {@link JcNode} as JSONC text (Go {@code internal/jsonc/to_string.go}).
 */
public final class JsoncPrinter {
    private static final String INDENT = "\t";

    private JsoncPrinter() {}

    public static String toString(JcNode n) {
        if (n == null) {
            return "";
        }
        StringBuilder b = new StringBuilder();
        writeLeadingComment(b, n.tag(), 0);
        writeNode(b, n, 0);
        return b.toString();
    }

    private static void writeIndent(StringBuilder b, int indent) {
        b.append(INDENT.repeat(Math.max(0, indent)));
    }

    private static void writeLeadingComment(StringBuilder b, MmTag tag, int indent) {
        if (tag == null) {
            return;
        }
        String s = MmTagDirectives.format(tag);
        if (s == null || s.isEmpty()) {
            return;
        }
        b.append('\n');
        writeIndent(b, indent);
        b.append("// mm: ").append(s).append('\n');
    }

    private static void writeNode(StringBuilder b, JcNode n, int indent) {
        if (n instanceof JcNode.JcValue v) {
            writeValue(b, v);
        } else if (n instanceof JcNode.JcObject o) {
            writeObject(b, o, indent);
        } else if (n instanceof JcNode.JcArray a) {
            writeArray(b, a, indent);
        }
    }

    private static void writeValue(StringBuilder b, JcNode.JcValue v) {
        MmTag tag = v.tag();
        if (tag == null) {
            return;
        }
        ValueType t = tag.type;
        if (t == null) {
            t = ValueType.STRING;
        }
        switch (t) {
            case STRING, BYTES, DATETIME, DATE, TIME, UUID, IP, URL, EMAIL, ENUM -> b.append(jsonQuote(v.text()));
            case INT, INT8, INT16, INT32, INT64, UINT, UINT8, UINT16, UINT32, UINT64, BIGINT, DECIMAL, BOOL -> b.append(v.text());
            case FLOAT32, FLOAT64 -> b.append(v.text());
            default -> b.append(v.text());
        }
    }

    private static void writeObject(StringBuilder b, JcNode.JcObject o, int indent) {
        b.append("{\n");
        for (JcField f : o.fields()) {
            JcNode child = f.value();
            writeLeadingComment(b, child.tag(), indent + 1);
            writeIndent(b, indent + 1);
            b.append(jsonQuote(f.key())).append(": ");
            writeNode(b, child, indent + 1);
            b.append(",\n");
        }
        writeIndent(b, indent);
        b.append('}');
    }

    private static void writeArray(StringBuilder b, JcNode.JcArray a, int indent) {
        b.append("[\n");
        for (JcNode item : a.items()) {
            writeLeadingComment(b, item.tag(), indent + 1);
            writeIndent(b, indent + 1);
            writeNode(b, item, indent + 1);
            b.append(",\n");
        }
        writeIndent(b, indent);
        b.append(']');
    }

    static String jsonQuote(String s) {
        if (s == null) {
            return "null";
        }
        StringBuilder sb = new StringBuilder(s.length() + 8);
        sb.append('"');
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            switch (c) {
                case '"' -> sb.append("\\\"");
                case '\\' -> sb.append("\\\\");
                case '\b' -> sb.append("\\b");
                case '\f' -> sb.append("\\f");
                case '\n' -> sb.append("\\n");
                case '\r' -> sb.append("\\r");
                case '\t' -> sb.append("\\t");
                default -> {
                    if (c < 0x20) {
                        sb.append(String.format(java.util.Locale.ROOT, "\\u%04x", (int) c));
                    } else {
                        sb.append(c);
                    }
                }
            }
        }
        sb.append('"');
        return sb.toString();
    }
}
