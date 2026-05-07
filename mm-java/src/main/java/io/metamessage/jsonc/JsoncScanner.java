package io.github.metamessage.jsonc;

import java.util.ArrayList;
import java.util.List;

/**
 * JSONC lexer (Go {@code internal/jsonc/scanner}).
 */
public final class JsoncScanner {
    private final char[] src;
    private int pos;
    private int line = 1;
    private int col = 1;
    private boolean newLine = true;

    public JsoncScanner(String input) {
        this.src = input.toCharArray();
    }

    public List<JsoncToken> tokenizeAll() {
        List<JsoncToken> toks = new ArrayList<>();
        for (;;) {
            JsoncToken t = nextToken();
            toks.add(t);
            if (t.type() == JsoncTokenType.EOF) {
                return toks;
            }
        }
    }

    public JsoncToken nextToken() {
        skipWhitespace();
        if (pos >= src.length) {
            return new JsoncToken(JsoncTokenType.EOF, "", line, col);
        }
        int startLine = line;
        int startCol = col;
        char ch = src[pos];
        switch (ch) {
            case '{' -> {
                advance();
                return new JsoncToken(JsoncTokenType.LBRACE, "{", startLine, startCol);
            }
            case '}' -> {
                advance();
                return new JsoncToken(JsoncTokenType.RBRACE, "}", startLine, startCol);
            }
            case '[' -> {
                advance();
                return new JsoncToken(JsoncTokenType.LBRACKET, "[", startLine, startCol);
            }
            case ']' -> {
                advance();
                return new JsoncToken(JsoncTokenType.RBRACKET, "]", startLine, startCol);
            }
            case ':' -> {
                advance();
                newLine = false;
                return new JsoncToken(JsoncTokenType.COLON, ":", startLine, startCol);
            }
            case ',' -> {
                advance();
                newLine = false;
                return new JsoncToken(JsoncTokenType.COMMA, ",", startLine, startCol);
            }
            case '"' -> {
                return scanString(startLine, startCol);
            }
            case '/' -> {
                return scanComment(startLine, startCol);
            }
            default -> {
                return scanLiteral(startLine, startCol);
            }
        }
    }

    private void skipWhitespace() {
        while (pos < src.length) {
            char c = src[pos];
            if (c == ' ' || c == '\t' || c == '\r') {
                advance();
            } else if (c == '\n') {
                advance();
                newLine = true;
            } else {
                break;
            }
        }
    }

    private JsoncToken scanString(int startLine, int startCol) {
        advance();
        StringBuilder buf = new StringBuilder();
        while (pos < src.length) {
            char c = src[pos];
            if (c == '\n' || c == 0) {
                break;
            }
            if (c == '"') {
                advance();
                break;
            }
            if (c == '\\') {
                buf.append(c);
                advance();
                if (pos < src.length) {
                    buf.append(src[pos]);
                    advance();
                }
            } else {
                buf.append(c);
                advance();
            }
        }
        return new JsoncToken(JsoncTokenType.STRING, unescapeString(buf.toString()), startLine, startCol);
    }

    private String unescapeString(String s) {
        if (!s.contains("\\")) {
            return s;
        }
        StringBuilder o = new StringBuilder();
        for (int i = 0; i < s.length(); i++) {
            if (s.charAt(i) == '\\' && i + 1 < s.length()) {
                o.append(s.charAt(++i));
            } else {
                o.append(s.charAt(i));
            }
        }
        return o.toString();
    }

    private JsoncToken scanComment(int startLine, int startCol) {
        advance();
        if (pos < src.length && src[pos] == '/') {
            JsoncTokenType kind = newLine ? JsoncTokenType.LEADING_COMMENT : JsoncTokenType.TRAILING_COMMENT;
            advance();
            StringBuilder buf = new StringBuilder();
            while (pos < src.length && src[pos] != '\n') {
                buf.append(src[pos]);
                advance();
            }
            return new JsoncToken(kind, buf.toString().trim(), startLine, startCol);
        }
        if (pos < src.length && src[pos] == '*') {
            JsoncTokenType kind = newLine ? JsoncTokenType.LEADING_COMMENT : JsoncTokenType.TRAILING_COMMENT;
            advance();
            StringBuilder buf = new StringBuilder();
            while (pos < src.length) {
                if (src[pos] == '*' && pos + 1 < src.length && src[pos + 1] == '/') {
                    advance();
                    advance();
                    break;
                }
                buf.append(src[pos]);
                advance();
            }
            return new JsoncToken(kind, buf.toString().trim(), startLine, startCol);
        }
        return new JsoncToken(JsoncTokenType.EOF, "", startLine, startCol);
    }

    private JsoncToken scanLiteral(int startLine, int startCol) {
        StringBuilder buf = new StringBuilder();
        while (pos < src.length) {
            char c = src[pos];
            if (c == 0 || " \t\r\n,:{}[]/".indexOf(c) >= 0) {
                break;
            }
            buf.append(c);
            advance();
        }
        String lit = buf.toString();
        return switch (lit) {
            case "true" -> new JsoncToken(JsoncTokenType.TRUE, lit, startLine, startCol);
            case "false" -> new JsoncToken(JsoncTokenType.FALSE, lit, startLine, startCol);
            case "null" -> new JsoncToken(JsoncTokenType.NULL, lit, startLine, startCol);
            default -> new JsoncToken(JsoncTokenType.NUMBER, lit, startLine, startCol);
        };
    }

    private void advance() {
        if (pos < src.length) {
            char c = src[pos++];
            if (c == '\n') {
                newLine = true;
                line++;
                col = 1;
            } else {
                col++;
            }
        }
    }
}
