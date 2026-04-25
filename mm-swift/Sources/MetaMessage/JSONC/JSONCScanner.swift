import Foundation

public enum JSONCTokenType: Equatable {
    case eof
    case lBrace
    case rBrace
    case lBracket
    case rBracket
    case colon
    case comma
    case string
    case number
    case trueValue
    case falseValue
    case nullValue
    case leadingComment
    case trailingComment
}

public struct JSONCToken {
    public let type: JSONCTokenType
    public let literal: String
    public let line: Int
    public let column: Int

    public init(type: JSONCTokenType, literal: String = "", line: Int, column: Int) {
        self.type = type
        self.literal = literal
        self.line = line
        self.column = column
    }
}

public class JSONCScanner {
    private var src: [Character]
    private var pos: Int
    private var line: Int
    private var col: Int
    private var newLine: Bool

    public init(input: String) {
        self.src = Array(input)
        self.pos = 0
        self.line = 1
        self.col = 1
        self.newLine = false
    }

    private func peek() -> Character {
        guard pos < src.count else { return "\0" }
        return src[pos]
    }

    private func next() -> Character {
        guard pos < src.count else { return "\0" }
        let ch = src[pos]
        pos += 1

        if ch == "\n" {
            newLine = true
            line += 1
            col = 1
        } else {
            col += 1
        }
        return ch
    }

    private func skipWhitespace() {
        while let ch = peekIfNotEmpty(), ch.isWhitespace {
            _ = next()
        }
    }

    private func peekIfNotEmpty() -> Character? {
        guard pos < src.count else { return nil }
        return src[pos]
    }

    public func nextToken() -> JSONCToken {
        skipWhitespace()

        let ch = peekIfNotEmpty()
        if ch == nil || ch == "\0" {
            return JSONCToken(type: .eof, line: line, column: col)
        }

        let startLine = line
        let startCol = col

        switch ch {
        case "{":
            _ = next()
            return JSONCToken(type: .lBrace, line: startLine, column: startCol)
        case "}":
            _ = next()
            return JSONCToken(type: .rBrace, line: startLine, column: startCol)
        case "[":
            _ = next()
            return JSONCToken(type: .lBracket, line: startLine, column: startCol)
        case "]":
            _ = next()
            return JSONCToken(type: .rBracket, line: startLine, column: startCol)
        case ":":
            _ = next()
            newLine = false
            return JSONCToken(type: .colon, line: startLine, column: startCol)
        case ",":
            _ = next()
            newLine = false
            return JSONCToken(type: .comma, line: startLine, column: startCol)
        case "\"":
            return scanString(startLine: startLine, startCol: startCol)
        case "/":
            return scanComment(startLine: startLine, startCol: startCol)
        default:
            return scanLiteral(startLine: startLine, startCol: startCol)
        }
    }

    private func scanString(startLine: Int, startCol: Int) -> JSONCToken {
        _ = next()

        var buf = ""
        while let ch = peekIfNotEmpty(), ch != "\0" {
            let c = next()
            if c == "\n" || c == "\0" {
                break
            }
            if c == "\"" {
                break
            }
            if c == "\\" {
                buf.append(c)
                if let escaped = peekIfNotEmpty(), escaped != "\0" {
                    buf.append(next())
                }
                continue
            }
            buf.append(c)
        }

        return JSONCToken(type: .string, literal: buf, line: startLine, column: startCol)
    }

    private func scanComment(startLine: Int, startCol: Int) -> JSONCToken {
        _ = next()

        if peekIfNotEmpty() == "/" {
            let commentType: JSONCTokenType = newLine ? .leadingComment : .trailingComment
            _ = next()

            var buf = ""
            while let ch = peekIfNotEmpty(), ch != "\0", ch != "\n" {
                buf.append(next())
            }
            let trimmed = buf.trimmingCharacters(in: .whitespaces)

            return JSONCToken(type: commentType, literal: trimmed, line: startLine, column: startCol)
        }

        if peekIfNotEmpty() == "*" {
            let commentType: JSONCTokenType = newLine ? .leadingComment : .trailingComment
            _ = next()

            var buf = ""
            while let ch = peekIfNotEmpty(), ch != "\0" {
                if ch == "*" && pos + 1 < src.count && src[pos + 1] == "/" {
                    _ = next()
                    _ = next()
                    break
                }
                buf.append(next())
            }
            let trimmed = buf.trimmingCharacters(in: .whitespaces)

            return JSONCToken(type: commentType, literal: trimmed, line: startLine, column: startCol)
        }

        return JSONCToken(type: .eof, line: startLine, column: startCol)
    }

    private func scanLiteral(startLine: Int, startCol: Int) -> JSONCToken {
        var buf = ""

        while let ch = peekIfNotEmpty(), ch != "\0", !ch.isWhitespace, !isDelimiter(ch) {
            buf.append(next())
        }

        switch buf {
        case "true":
            return JSONCToken(type: .trueValue, literal: buf, line: startLine, column: startCol)
        case "false":
            return JSONCToken(type: .falseValue, literal: buf, line: startLine, column: startCol)
        case "null":
            return JSONCToken(type: .nullValue, literal: buf, line: startLine, column: startCol)
        default:
            return JSONCToken(type: .number, literal: buf, line: startLine, column: startCol)
        }
    }

    private func isDelimiter(_ ch: Character) -> Bool {
        return ch == "{" || ch == "}" || ch == "[" || ch == "]" || ch == ":" || ch == ","
    }
}