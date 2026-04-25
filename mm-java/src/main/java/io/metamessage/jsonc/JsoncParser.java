package io.metamessage.jsonc;

import io.metamessage.mm.CamelToSnake;
import io.metamessage.mm.MmLineParser;
import io.metamessage.mm.MmTag;
import io.metamessage.mm.ValueType;
import java.util.ArrayList;
import java.util.List;

/**
 * JSONC parser (Go {@code internal/jsonc/parser}).
 */
public final class JsoncParser {
    private static final int MAX_DEPTH = 32;

    private final List<JsoncToken> toks;
    private int pos;
    private final List<JsoncToken> pending = new ArrayList<>();
    private int depth;

    public JsoncParser(List<JsoncToken> toks) {
        this.toks = toks;
    }

    public JcNode parseDocument() {
        JcNode val = null;
        while (true) {
            JsoncToken tok = peek();
            if (tok.type() == JsoncTokenType.EOF) {
                return val;
            }
            if (tok.type() == JsoncTokenType.LEADING_COMMENT) {
                if (!pending.isEmpty()) {
                    JsoncToken last = pending.get(pending.size() - 1);
                    if (tok.line() - last.line() > 1) {
                        pending.clear();
                    }
                }
                pending.add(tok);
                next();
                continue;
            }
            if (tok.type() == JsoncTokenType.TRAILING_COMMENT) {
                if (val != null) {
                    MmTag p = parseCommentToTag(tok.literal());
                    if (p != null) {
                        val = mergeNodeTag(val, p);
                    }
                }
                next();
                continue;
            }
            val = parse("");
            return val;
        }
    }

    private JcNode mergeNodeTag(JcNode n, MmTag parsed) {
        if (n == null || parsed == null) {
            return n;
        }
        return setTag(n, MmTag.merge(childTagOf(n), parsed));
    }

    private JcNode parse(String path) {
        JsoncToken tok = next();
        return switch (tok.type()) {
            case EOF -> null;
            case LBRACE -> parseObject(tok.line(), path);
            case LBRACKET -> parseArray(tok.line(), path);
            case STRING -> parseStringToken(path, tok);
            case NUMBER -> JsoncValueParser.numberLiteral(
                    tok.literal(), coalesce(consumeTagForLine(tok.line())), path);
            case TRUE -> JsoncValueParser.boolLiteral(
                    true, coalesce(consumeTagForLine(tok.line())), path);
            case FALSE -> JsoncValueParser.boolLiteral(
                    false, coalesce(consumeTagForLine(tok.line())), path);
            case NULL -> throw new JsoncException("null literal not supported in this parser");
            default -> throw new JsoncException("unexpected " + tok.type());
        };
    }

    private static MmTag coalesce(MmTag t) {
        return t == null ? MmTag.empty() : t;
    }

    private JcNode parseStringToken(String path, JsoncToken tok) {
        MmTag tag = coalesce(consumeTagForLine(tok.line()));
        if (tag.type == ValueType.UNKNOWN) {
            tag = tag.copy();
            tag.type = ValueType.STRING;
        }
        if (isSimpleFieldKeyword(tok.literal())) {
            MmTag t = tag.copy();
            t.type = ValueType.STRING;
            return new JcNode.JcValue(tok.literal(), tok.literal(), t, path);
        }
        return JsoncValueParser.stringLiteral(tok.literal(), tag, path);
    }

    private static boolean isSimpleFieldKeyword(String text) {
        return switch (text) {
            case "code", "message", "data", "success", "error", "unknown", "page", "limit", "offset", "total", "id", "name",
                    "description", "type", "version", "status", "url", "create_time", "update_time", "delete_time", "account",
                    "token", "expire_time", "key", "value" -> true;
            default -> false;
        };
    }

    private MmTag consumeTagForLine(int line) {
        if (pending.isEmpty()) {
            return null;
        }
        JsoncToken last = pending.get(pending.size() - 1);
        if (line - last.line() > 1) {
            pending.clear();
            return null;
        }
        MmTag out = MmTag.empty();
        for (JsoncToken ct : pending) {
            MmTag p = parseCommentToTag(ct.literal());
            if (p != null) {
                out = MmTag.merge(out, p);
            }
        }
        pending.clear();
        return out;
    }

    private MmTag parseCommentToTag(String lit) {
        if (lit == null) {
            return null;
        }
        if (lit.startsWith("mm:") || lit.contains("mm:")) {
            return MmLineParser.parse(lit);
        }
        return null;
    }

    private JcNode.JcObject parseObject(int openLine, String path) {
        depth++;
        if (depth > MAX_DEPTH) {
            throw new JsoncException("max depth " + MAX_DEPTH);
        }
        try {
            MmTag tag = consumeTagForLine(openLine);
            if (tag == null) {
                tag = MmTag.empty();
            }
            if (tag.type == ValueType.UNKNOWN) {
                tag = tag.copy();
                tag.type = ValueType.STRUCT;
            }
            if (tag.name != null && !tag.name.isEmpty()) {
                if (path.isEmpty()) {
                    path = tag.name;
                } else {
                    path = path + "." + tag.name;
                }
            }
            List<JcField> fields = new ArrayList<>();
            JcNode lastVal = null;
            String p0 = path;
            while (true) {
                JsoncToken t = peek();
                if (t.type() == JsoncTokenType.EOF) {
                    break;
                }
                if (t.type() == JsoncTokenType.RBRACE) {
                    next();
                    break;
                }
                if (t.type() == JsoncTokenType.LEADING_COMMENT) {
                    if (!pending.isEmpty()) {
                        JsoncToken last = pending.get(pending.size() - 1);
                        if (t.line() - last.line() > 1) {
                            pending.clear();
                        }
                    }
                    pending.add(t);
                    next();
                    continue;
                }
                if (t.type() == JsoncTokenType.TRAILING_COMMENT) {
                    if (lastVal != null) {
                        MmTag p = parseCommentToTag(t.literal());
                        if (p != null) {
                            lastVal = mergeNodeTag(lastVal, p);
                            if (!fields.isEmpty()) {
                                JcField old = fields.get(fields.size() - 1);
                                fields.set(fields.size() - 1, new JcField(old.key(), lastVal));
                            }
                        }
                    }
                    next();
                    continue;
                }
                JsoncToken keyT = next();
                if (keyT.type() != JsoncTokenType.STRING) {
                    throw new JsoncException("expected string key");
                }
                if (next().type() != JsoncTokenType.COLON) {
                    throw new JsoncException("expected :");
                }
                String keyStr = CamelToSnake.convert(keyT.literal());
                String pa = p0 + "." + keyStr;
                if (tag.type == ValueType.MAP) {
                    pa = p0 + "[" + keyStr + "]";
                }
                JcNode val = parse(pa);
                if (val == null) {
                    continue;
                }
                if (tag != null) {
                    MmTag c = childTagOf(val).copy();
                    c.inheritFromArrayParent(tag);
                    val = setTag(val, c);
                }
                fields.add(new JcField(keyStr, val));
                lastVal = val;
                if (peek().type() == JsoncTokenType.COMMA) {
                    next();
                }
            }
            return new JcNode.JcObject(tag, p0, fields);
        } finally {
            depth--;
        }
    }

    private MmTag childTagOf(JcNode n) {
        if (n instanceof JcNode.JcValue v) {
            return v.tag();
        }
        if (n instanceof JcNode.JcObject o) {
            return o.tag();
        }
        if (n instanceof JcNode.JcArray a) {
            return a.tag();
        }
        throw new IllegalStateException();
    }

    private JcNode setTag(JcNode n, MmTag t) {
        if (n instanceof JcNode.JcValue v) {
            return new JcNode.JcValue(v.data(), v.text(), t, v.path());
        }
        if (n instanceof JcNode.JcObject o) {
            return new JcNode.JcObject(t, o.path(), o.fields());
        }
        if (n instanceof JcNode.JcArray a) {
            return new JcNode.JcArray(t, a.path(), a.items());
        }
        throw new IllegalStateException();
    }

    private JcNode.JcArray parseArray(int openLine, String path) {
        depth++;
        if (depth > MAX_DEPTH) {
            throw new JsoncException("max depth");
        }
        try {
            MmTag tag = consumeTagForLine(openLine);
            if (tag == null) {
                tag = MmTag.empty();
            }
            if (tag.type == ValueType.UNKNOWN) {
                tag = tag.copy();
                tag.type = tag.size > 0 ? ValueType.ARRAY : ValueType.SLICE;
            }
            if (tag.name != null && !tag.name.isEmpty()) {
                path = path + "." + tag.name;
            }
            List<JcNode> items = new ArrayList<>();
            JcNode lastItem = null;
            int i = 0;
            while (true) {
                JsoncToken t = peek();
                if (t.type() == JsoncTokenType.EOF) {
                    break;
                }
                if (t.type() == JsoncTokenType.RBRACKET) {
                    next();
                    break;
                }
                if (t.type() == JsoncTokenType.LEADING_COMMENT) {
                    if (!pending.isEmpty()) {
                        JsoncToken last = pending.get(pending.size() - 1);
                        if (t.line() - last.line() > 1) {
                            pending.clear();
                        }
                    }
                    pending.add(t);
                    next();
                    continue;
                }
                if (t.type() == JsoncTokenType.TRAILING_COMMENT) {
                    if (lastItem != null) {
                        MmTag p = parseCommentToTag(t.literal());
                        if (p != null) {
                            lastItem = mergeNodeTag(lastItem, p);
                            if (!items.isEmpty()) {
                                items.set(items.size() - 1, lastItem);
                            }
                        }
                    }
                    next();
                    continue;
                }
                String pa = path + "[" + i + "]";
                JcNode item = parse(pa);
                if (item == null) {
                    continue;
                }
                MmTag ch = childTagOf(item);
                if (ch != null) {
                    MmTag c = ch.copy();
                    c.inheritFromArrayParent(tag);
                    item = setTag(item, c);
                }
                items.add(item);
                lastItem = item;
                i++;
                if (peek().type() == JsoncTokenType.COMMA) {
                    next();
                }
            }
            return new JcNode.JcArray(tag, path, items);
        } finally {
            depth--;
        }
    }

    private JsoncToken peek() {
        if (pos >= toks.size()) {
            return new JsoncToken(JsoncTokenType.EOF, "", 0, 0);
        }
        return toks.get(pos);
    }

    private JsoncToken next() {
        if (pos >= toks.size()) {
            return new JsoncToken(JsoncTokenType.EOF, "", 0, 0);
        }
        return toks.get(pos++);
    }
}
