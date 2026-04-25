import Foundation

public enum JSONCParserError: Error {
    case unexpectedToken(String)
    case unexpectedEndOfData
    case invalidData(String)
    case maxDepthExceeded
}

public class JSONCParser {
    private var tokens: [JSONCToken]
    private var pos: Int
    private var pendingComments: [JSONCToken]
    private var depth: Int
    private let maxDepth: Int = 32

    public init(tokens: [JSONCToken]) {
        self.tokens = tokens
        self.pos = 0
        self.pendingComments = []
        self.depth = 0
    }

    private func peek() -> JSONCToken {
        guard pos < tokens.count else {
            return JSONCToken(type: .eof, line: 0, column: 0)
        }
        return tokens[pos]
    }

    private func next() -> JSONCToken {
        let token = peek()
        pos += 1
        return token
    }

    private func consumeCommentsFor(_ anchorLine: Int) -> JSONCTag? {
        guard !pendingComments.isEmpty else { return nil }

        let last = pendingComments[pendingComments.count - 1]
        if anchorLine - last.line > 1 {
            pendingComments = []
            return nil
        }

        var result: JSONCTag?
        for comment in pendingComments {
            if let parsed = parseCommentToTag(comment.literal) {
                result = mergeTag(result, parsed)
            }
        }

        pendingComments = []
        return result
    }

    private func mergeTag(_ dst: JSONCTag?, _ src: JSONCTag) -> JSONCTag {
        if dst == nil {
            return src
        }

        let merged = JSONCTag()
        merged.name = src.name.isEmpty ? (dst?.name ?? "") : src.name

        if src.isNull { merged.isNull = src.isNull }
        if src.example { merged.example = src.example }
        if !src.desc.isEmpty { merged.desc = src.desc }
        if src.type != .unknown { merged.type = src.type }
        if src.raw { merged.raw = src.raw }
        if src.nullable { merged.nullable = src.nullable }
        if src.allowEmpty { merged.allowEmpty = src.allowEmpty }
        if src.unique { merged.unique = src.unique }
        if !src.defaultValue.isEmpty { merged.defaultValue = src.defaultValue }
        if !src.min.isEmpty { merged.min = src.min }
        if !src.max.isEmpty { merged.max = src.max }
        if src.size != 0 { merged.size = src.size }
        if !src.enumValues.isEmpty { merged.enumValues = src.enumValues }
        if !src.pattern.isEmpty { merged.pattern = src.pattern }
        if src.locationOffset != 0 { merged.locationOffset = src.locationOffset }
        if src.version != 0 { merged.version = src.version }
        if !src.mime.isEmpty { merged.mime = src.mime }

        if !src.childDesc.isEmpty { merged.childDesc = src.childDesc }
        if src.childType != .unknown { merged.childType = src.childType }
        if src.childRaw { merged.childRaw = src.childRaw }
        if src.childNullable { merged.childNullable = src.childNullable }
        if src.childAllowEmpty { merged.childAllowEmpty = src.childAllowEmpty }
        if src.childUnique { merged.childUnique = src.childUnique }
        if !src.childDefault.isEmpty { merged.childDefault = src.childDefault }
        if !src.childMin.isEmpty { merged.childMin = src.childMin }
        if !src.childMax.isEmpty { merged.childMax = src.childMax }
        if src.childSize != 0 { merged.childSize = src.childSize }
        if !src.childEnum.isEmpty { merged.childEnum = src.childEnum }
        if !src.childPattern.isEmpty { merged.childPattern = src.childPattern }
        if src.childLocationOffset != 0 { merged.childLocationOffset = src.childLocationOffset }
        if src.childVersion != 0 { merged.childVersion = src.childVersion }
        if !src.childMime.isEmpty { merged.childMime = src.childMime }

        return merged
    }

    private func parseCommentToTag(_ literal: String) -> JSONCTag? {
        if literal.hasPrefix("mm:") {
            return parseMMTag(String(literal.dropFirst(3)))
        }
        return nil
    }

    public func parse() throws -> JSONCNode? {
        var result: JSONCNode?

        while true {
            let tok = peek()
            if tok.type == .eof {
                return result
            }

            if tok.type == .leadingComment {
                if !pendingComments.isEmpty {
                    let last = pendingComments[pendingComments.count - 1]
                    if tok.line - last.line > 1 {
                        pendingComments = []
                    }
                }
                pendingComments.append(tok)
                _ = next()
                continue
            }

            if tok.type == .trailingComment {
                if let val = result {
                    if let parsed = parseCommentToTag(tok.literal) {
                        mergeNodeTag(val, parsed)
                    }
                }
                _ = next()
                continue
            }

            result = try parseNode("")
        }
    }

    private func parseNode(_ path: String) throws -> JSONCNode? {
        let tok = next()

        switch tok.type {
        case .eof:
            return nil

        case .lBrace:
            return try parseObject(tok.line, path)

        case .lBracket:
            return try parseArray(tok.line, path)

        case .string:
            var tag = try consumeCommentsFor(tok.line) ?? JSONCTag()
            if tag.type == .unknown {
                tag.type = .string
            }
            let text = tok.literal

            return JSONCValue(data: text, text: text, tag: tag, path: path)

        case .number:
            var tag = try consumeCommentsFor(tok.line) ?? JSONCTag()
            if tag.type == .unknown {
                if tok.literal.contains(".") {
                    tag.type = .float64
                } else if tok.literal.hasPrefix("-") {
                    tag.type = .int
                } else {
                    tag.type = .int
                }
            }

            var data: Any?
            if tok.literal.contains(".") {
                data = Double(tok.literal)
            } else if tok.literal.hasPrefix("-") {
                data = Int64(tok.literal)
            } else {
                if let uval = UInt64(tok.literal) {
                    if uval > UInt64(Int.max) {
                        data = uval
                    } else {
                        data = Int(uval)
                    }
                }
            }

            return JSONCValue(data: data, text: tok.literal, tag: tag, path: path)

        case .trueValue:
            var tag = try consumeCommentsFor(tok.line) ?? JSONCTag()
            if tag.type == .unknown {
                tag.type = .bool
            }
            return JSONCValue(data: true, text: "true", tag: tag, path: path)

        case .falseValue:
            var tag = try consumeCommentsFor(tok.line) ?? JSONCTag()
            if tag.type == .unknown {
                tag.type = .bool
            }
            return JSONCValue(data: false, text: "false", tag: tag, path: path)

        case .nullValue:
            var tag = try consumeCommentsFor(tok.line) ?? JSONCTag()
            if tag.type == .unknown {
                tag.type = .unknown
            }
            tag.isNull = true
            return JSONCValue(data: nil, text: "null", tag: tag, path: path)

        default:
            throw JSONCParserError.unexpectedToken("Unexpected token: \(tok.type)")
        }
    }

    private func parseObject(_ openLine: Int, _ path: String) throws -> JSONCObject {
        depth += 1
        if depth > maxDepth {
            throw JSONCParserError.maxDepthExceeded
        }

        defer { depth -= 1 }

        var tag = try consumeCommentsFor(openLine) ?? JSONCTag()
        if tag.type == .unknown {
            tag.type = .structType
        }

        let obj = JSONCObject(tag: tag, path: path)

        while true {
            let tok = peek()
            if tok.type == .eof {
                break
            }
            if tok.type == .rBrace {
                _ = next()
                break
            }

            if tok.type == .leadingComment {
                if !pendingComments.isEmpty {
                    let last = pendingComments[pendingComments.count - 1]
                    if tok.line - last.line > 1 {
                        pendingComments = []
                    }
                }
                pendingComments.append(tok)
                _ = next()
                continue
            }

            if tok.type == .trailingComment {
                if let lastField = obj.fields.last, let val = lastField.value as? JSONCNode {
                    if let parsed = parseCommentToTag(tok.literal) {
                        mergeNodeTag(val, parsed)
                    }
                }
                _ = next()
                continue
            }

            let keyTok = next()
            guard keyTok.type == .string else {
                throw JSONCParserError.unexpectedToken("Expected string key")
            }

            let key = keyTok.literal

            _ = next()

            let childPath = "\(path).\(key)"
            if let val = try parseNode(childPath) {
                let childTag = val.getTag()
                if let ct = childTag, let t = tag as JSONCTag? {
                    ct.inherit(from: t)
                }
                let field = JSONCField(key: key, value: val)
                obj.fields.append(field)
            }

            if peek().type == .comma {
                _ = next()
            }
        }

        return obj
    }

    private func parseArray(_ openLine: Int, _ path: String) throws -> JSONCArray {
        depth += 1
        if depth > maxDepth {
            throw JSONCParserError.maxDepthExceeded
        }

        defer { depth -= 1 }

        var tag = try consumeCommentsFor(openLine) ?? JSONCTag()
        if tag.type == .unknown {
            if tag.size > 0 {
                tag.type = .array
            } else {
                tag.type = .slice
            }
        }

        let arr = JSONCArray(tag: tag, path: path)

        var index = 0
        while true {
            let tok = peek()
            if tok.type == .eof {
                break
            }
            if tok.type == .rBracket {
                _ = next()
                break
            }

            if tok.type == .leadingComment {
                if !pendingComments.isEmpty {
                    let last = pendingComments[pendingComments.count - 1]
                    if tok.line - last.line > 1 {
                        pendingComments = []
                    }
                }
                pendingComments.append(tok)
                _ = next()
                continue
            }

            if tok.type == .trailingComment {
                if let lastItem = arr.items.last, let val = lastItem as? JSONCNode {
                    if let parsed = parseCommentToTag(tok.literal) {
                        mergeNodeTag(val, parsed)
                    }
                }
                _ = next()
                continue
            }

            let itemPath = "\(path)[\(index)]"
            if let item = try parseNode(itemPath) {
                let childTag = item.getTag()
                if let ct = childTag, let t = tag as JSONCTag? {
                    ct.inherit(from: t)
                }
                arr.items.append(item)
                index += 1
            }

            if peek().type == .comma {
                _ = next()
            }
        }

        return arr
    }

    private func mergeNodeTag(_ node: JSONCNode, _ tag: JSONCTag) {
        guard let existing = node.getTag() else { return }

        if node is JSONCValue {
            (node as? JSONCValue)?.tag = mergeTag(existing, tag)
        } else if node is JSONCObject {
            (node as? JSONCObject)?.tag = mergeTag(existing, tag)
        } else if node is JSONCArray {
            (node as? JSONCArray)?.tag = mergeTag(existing, tag)
        }
    }
}

public func parseJSONC(_ input: String) throws -> JSONCNode? {
    let scanner = JSONCScanner(input: input)
    var tokens: [JSONCToken] = []

    while true {
        let token = scanner.nextToken()
        tokens.append(token)
        if token.type == .eof {
            break
        }
    }

    let parser = JSONCParser(tokens: tokens)
    return try parser.parse()
}