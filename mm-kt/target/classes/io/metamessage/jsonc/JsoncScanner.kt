package io.metamessage.jsonc

class JsoncScanner(private val source: String) {
    private val src: List<Char> = source.toList()
    private var pos: Int = 0
    private var line: Int = 1
    private var col: Int = 1
    private var newLine: Boolean = false

    private fun peek(): Char {
        return if (pos >= src.size) '\u0000' else src[pos]
    }

    private fun next(): Char {
        if (pos >= src.size) return '\u0000'
        val ch = src[pos]
        pos++
        if (ch == '\n') {
            newLine = true
            line++
            col = 1
        } else {
            col++
        }
        return ch
    }

    private fun skipWhitespace() {
        while (peek().isWhitespace()) {
            next()
        }
    }

    fun nextToken(): JsoncToken {
        skipWhitespace()

        val ch = peek()
        if (ch == '\u0000') {
            return JsoncToken(JsoncTokenType.EOF, "", line, col)
        }

        val startLine = line
        val startCol = col

        when (ch) {
            '{' -> {
                next()
                return JsoncToken(JsoncTokenType.LBrace, "{", startLine, startCol)
            }
            '}' -> {
                next()
                return JsoncToken(JsoncTokenType.RBrace, "}", startLine, startCol)
            }
            '[' -> {
                next()
                return JsoncToken(JsoncTokenType.LBracket, "[", startLine, startCol)
            }
            ']' -> {
                next()
                return JsoncToken(JsoncTokenType.RBracket, "]", startLine, startCol)
            }
            ':' -> {
                next()
                newLine = false
                return JsoncToken(JsoncTokenType.Colon, ":", startLine, startCol)
            }
            ',' -> {
                next()
                newLine = false
                return JsoncToken(JsoncTokenType.Comma, ",", startLine, startCol)
            }
            '"' -> return scanString()
            '/' -> return scanComment()
            else -> return scanLiteral()
        }
    }

    private fun scanString(): JsoncToken {
        val startLine = line
        val startCol = col
        next()

        val buf = StringBuilder()
        while (true) {
            val ch = next()
            if (ch == '\u0000' || ch == '\n') break
            if (ch == '"') break
            if (ch == '\\') {
                buf.append(ch)
                buf.append(next())
                continue
            }
            buf.append(ch)
        }

        return JsoncToken(JsoncTokenType.String, buf.toString(), startLine, startCol)
    }

    private fun scanComment(): JsoncToken {
        val startLine = line
        val startCol = col
        next()

        if (peek() == '/') {
            val tokenType = if (!newLine) JsoncTokenType.TrailingComment else JsoncTokenType.LeadingComment
            next()
            val buf = StringBuilder()
            while (true) {
                val ch = peek()
                if (ch == '\n' || ch == '\u0000') break
                buf.append(next())
            }
            return JsoncToken(tokenType, buf.toString().trim(), startLine, startCol)
        }

        if (peek() == '*') {
            val tokenType = if (!newLine) JsoncTokenType.TrailingComment else JsoncTokenType.LeadingComment
            next()
            val buf = StringBuilder()
            while (true) {
                if (peek() == '\u0000') break
                if (peek() == '*' && pos + 1 < src.size && src[pos + 1] == '/') {
                    next()
                    next()
                    break
                }
                buf.append(next())
            }
            return JsoncToken(tokenType, buf.toString().trim(), startLine, startCol)
        }

        return JsoncToken(JsoncTokenType.EOF, "", line, col)
    }

    private fun scanLiteral(): JsoncToken {
        val startLine = line
        val startCol = col
        val buf = StringBuilder()

        while (true) {
            val ch = peek()
            if (ch == '\u0000' || " \t\r\n,:{}[]".contains(ch)) break
            buf.append(next())
        }

        val lit = buf.toString()
        return when (lit) {
            "true" -> JsoncToken(JsoncTokenType.True, lit, startLine, startCol)
            "false" -> JsoncToken(JsoncTokenType.False, lit, startLine, startCol)
            "null" -> JsoncToken(JsoncTokenType.Null, lit, startLine, startCol)
            else -> JsoncToken(JsoncTokenType.Number, lit, startLine, startCol)
        }
    }
}