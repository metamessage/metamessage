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

    private func consumeCommentsFor(_ anchorLine: Int) -> Tag? {
        guard !pendingComments.isEmpty else { return nil }

        let last = pendingComments[pendingComments.count - 1]
        if anchorLine - last.line > 1 {
            pendingComments = []
            return nil
        }
        if last.line > anchorLine {
            return nil
        }

        var result: Tag?
        for comment in pendingComments {
            if let parsed = parseCommentToTag(comment.literal) {
                result = mergeTag(result, parsed)
            }
        }

        pendingComments = []
        return result
    }

    private func mergeTag(_ dst: Tag?, _ src: Tag) -> Tag {
        if dst == nil {
            return src
        }

        let merged = Tag()
        merged.name = src.name.isEmpty ? (dst?.name ?? "") : src.name

        if src.isNull { merged.isNull = src.isNull }
        if src.example { merged.example = src.example }
        if !src.desc.isEmpty { merged.desc = src.desc }
        if src.type != .unknown { merged.type = src.type }
        if src.deprecated { merged.deprecated = src.deprecated }
        if src.nullable { merged.nullable = src.nullable }
        if src.allowEmpty { merged.allowEmpty = src.allowEmpty }
        if src.unique { merged.unique = src.unique }
        if !src.defaultVal.isEmpty { merged.defaultVal = src.defaultVal }
        if !src.min.isEmpty { merged.min = src.min }
        if !src.max.isEmpty { merged.max = src.max }
        if src.size != 0 { merged.size = src.size }
        if !src.enums.isEmpty { merged.enums = src.enums }
        if !src.pattern.isEmpty { merged.pattern = src.pattern }
        if src.location != 0 { merged.location = src.location }
        if src.version != 0 { merged.version = src.version }
        if !src.mime.isEmpty { merged.mime = src.mime }

        if !src.childDesc.isEmpty { merged.childDesc = src.childDesc }
        if src.childType != .unknown { merged.childType = src.childType }
        if src.childNullable { merged.childNullable = src.childNullable }
        if src.childAllowEmpty { merged.childAllowEmpty = src.childAllowEmpty }
        if src.childUnique { merged.childUnique = src.childUnique }
        if !src.childDefaultVal.isEmpty { merged.childDefaultVal = src.childDefaultVal }
        if !src.childMin.isEmpty { merged.childMin = src.childMin }
        if !src.childMax.isEmpty { merged.childMax = src.childMax }
        if src.childSize != 0 { merged.childSize = src.childSize }
        if !src.childEnums.isEmpty { merged.childEnums = src.childEnums }
        if !src.childPattern.isEmpty { merged.childPattern = src.childPattern }
        if src.childLocation != 0 { merged.childLocation = src.childLocation }
        if src.childVersion != 0 { merged.childVersion = src.childVersion }
        if !src.childMime.isEmpty { merged.childMime = src.childMime }

        return merged
    }

    private func parseCommentToTag(_ literal: String) -> Tag? {
        if literal.hasPrefix("mm:") {
            return parseMMTag(String(literal.dropFirst(3)))
        }
        return nil
    }

    public func parse() throws -> Node? {
        var result: Node?

        while true {
            let tok = peek()
            if tok.type == .eof {
                return result
            }

            if tok.type == .comment {
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


            result = try parseNode("")
        }
    }

    private func parseNode(_ path: String, _ preTag: Tag? = nil) throws -> Node? {
        let tok = next()

        switch tok.type {
        case .eof:
            return nil

        case .lBrace:
            return try parseObject(tok.line, path, preTag)

        case .lBracket:
            return try parseArray(tok.line, path, preTag)

        case .string:
                let tag = preTag ?? consumeCommentsFor(tok.line)
                let text = tok.literal

                if let tag = tag {
                    if tag.type == .unknown {
                        tag.type = .str
                    }
                    if tag.isNull {
                        if text != "" {
                            throw JSONCParserError.invalidData("invalid string: \(text), valid: \"\"")
                        }
                        return Value(data: "", text: "", tag: tag, path: path)
                    }
                }

                var parsedData: Any = text
                if let tag = tag {
                    switch tag.type {
                    case .datetime:
                        let formatter = DateFormatter()
                        formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                        formatter.timeZone = TimeZone(abbreviation: "UTC")
                        if let date = formatter.date(from: text) {
                            parsedData = date
                        }
                    case .date:
                        let formatter = DateFormatter()
                        formatter.dateFormat = "yyyy-MM-dd"
                        formatter.timeZone = TimeZone(abbreviation: "UTC")
                        if let date = formatter.date(from: text) {
                            parsedData = date
                        }
                    case .time:
                        let formatter = DateFormatter()
                        formatter.dateFormat = "HH:mm:ss"
                        formatter.timeZone = TimeZone(abbreviation: "UTC")
                        if let date = formatter.date(from: text) {
                            parsedData = date
                        }
                    case .uuid:
                        let hexStr = text.replacingOccurrences(of: "-", with: "")
                        var uuidBytes = [UInt8]()
                        var index = hexStr.startIndex
                        while index < hexStr.endIndex {
                            let next = hexStr.index(index, offsetBy: 2)
                            if let byte = UInt8(hexStr[index..<next], radix: 16) {
                                uuidBytes.append(byte)
                            }
                            index = next
                        }
                        if uuidBytes.count == 16 {
                            parsedData = uuidBytes
                        }
                    case .bytes:
                        break
                    default:
                        break
                    }

                    let stringResult = validator.validate(parsedData, tag: tag)
                    if !stringResult.isValid {
                        throw JSONCParserError.invalidData(stringResult.errors.joined(separator: ", "))
                    }
                }

                let value = Value(data: parsedData, text: text, tag: tag, path: path)
                return value

        case .number:
            let tag = preTag ?? consumeCommentsFor(tok.line)
            if let tag = tag {
                if tag.type == .unknown {
                    if tok.literal.contains(".") {
                        tag.type = .f64
                    } else if tok.literal.hasPrefix("-") {
                        tag.type = .i
                    } else {
                        tag.type = .i
                    }
                }

                if tag.isNull {
                    if tok.literal.contains(".") {
                        if tok.literal != "0.0" {
                            throw JSONCParserError.invalidData("invalid float: \(tok.literal), valid: 0.0")
                        }
                        return Value(data: tag.type == .f32 ? Float(0.0) : Double(0.0), text: tok.literal, tag: tag, path: path)
                    } else {
                        if tok.literal != "0" {
                            throw JSONCParserError.invalidData("invalid int: \(tok.literal), valid: 0")
                        }
                        if tag.type == .bigint {
                            return Value(data: "0", text: tok.literal, tag: tag, path: path)
                        }
                        return Value(data: Int(0), text: tok.literal, tag: tag, path: path)
                    }
                }
            }

            var data: Any?

            if let tag = tag, tag.type == .bigint {
                data = tok.literal
                let numberResult = validator.validate(data!, tag: tag)
                if !numberResult.isValid {
                    throw JSONCParserError.invalidData(numberResult.errors.joined(separator: ", "))
                }
            } else {
                if tok.literal.contains(".") {
                    if let tag = tag, tag.type == .f32 {
                        data = Float(tok.literal)
                    } else {
                        data = Double(tok.literal)
                    }
                } else if tok.literal.hasPrefix("-") {
                    if let ival = Int(tok.literal) {
                        data = ival
                    } else {
                        data = Int64(tok.literal)
                    }
                } else {
                    if let uval = UInt64(tok.literal) {
                        if uval > UInt64(Int.max) {
                            data = uval
                        } else {
                            data = Int(uval)
                        }
                    }
                }

                if let tag = tag {
                    let numberResult = validator.validate(data, tag: tag)
                    if !numberResult.isValid {
                        throw JSONCParserError.invalidData(numberResult.errors.joined(separator: ", "))
                    }
                }
            }

            let value = Value(data: data, text: tok.literal, tag: tag, path: path)
            return value

        case .trueValue:
            let tag = preTag ?? consumeCommentsFor(tok.line)
            if let tag = tag {
                if tag.type == .unknown {
                    tag.type = .bool
                }
                if tag.isNull {
                    throw JSONCParserError.invalidData("bool must false when bool is null")
                }
                let trueResult = validator.validate(true, tag: tag)
                if !trueResult.isValid {
                    throw JSONCParserError.invalidData(trueResult.errors.joined(separator: ", "))
                }
            }
            let value = Value(data: true, text: "true", tag: tag, path: path)
            return value

        case .falseValue:
            let tag = preTag ?? consumeCommentsFor(tok.line)
            if let tag = tag {
                if tag.type == .unknown {
                    tag.type = .bool
                }
                if tag.isNull {
                    return Value(data: false, text: "false", tag: tag, path: path)
                }
                let falseResult = validator.validate(false, tag: tag)
                if !falseResult.isValid {
                    throw JSONCParserError.invalidData(falseResult.errors.joined(separator: ", "))
                }
            }
            let value = Value(data: false, text: "false", tag: tag, path: path)
            return value

        case .nullValue:
            throw JSONCParserError.invalidData("null is not supported")

        default:
            throw JSONCParserError.unexpectedToken("Unexpected token: \(tok.type)")
        }
    }

    private func parseObject(_ openLine: Int, _ path: String, _ preTag: Tag? = nil) throws -> MMObject {
        depth += 1
        if depth > maxDepth {
            throw JSONCParserError.maxDepthExceeded
        }

        defer { depth -= 1 }

        var tag = consumeCommentsFor(openLine)
        if let pt = preTag {
            if tag == nil {
                tag = Tag()
            }
            tag!.inherit(from: pt)
        }
        if let tag = tag {
            if tag.type == .unknown {
                tag.type = .obj
            }
        }

        let obj = MMObject(tag: tag, path: path)

        if let tag = tag, !tag.example {
            let structResult = validator.validate(obj, tag: tag)
            if !structResult.isValid {
                throw JSONCParserError.invalidData(structResult.errors.joined(separator: ", "))
            }
        }

        while true {
            let tok = peek()
            if tok.type == .eof {
                break
            }
            if tok.type == .rBrace {
                _ = next()
                break
            }

            if tok.type == .comment {
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


            let keyTok = next()
            guard keyTok.type == .string else {
                throw JSONCParserError.unexpectedToken("Expected string key")
            }

            let key = keyTok.literal

            _ = next()

            let childPath = "\(path).\(key)"
            let valueTok = peek()
            if let parentTag = tag {
                let ownTag = consumeCommentsFor(valueTok.line)
                let childTag = ownTag ?? Tag()
                childTag.inherit(from: parentTag)

                if let val = try parseNode(childPath, childTag) {
                    let field = Field(key: key, value: val)
                    obj.fields.append(field)
                }
            } else {
                if let val = try parseNode(childPath) {
                    let field = Field(key: key, value: val)
                    obj.fields.append(field)
                }
            }

            if peek().type == .comma {
                _ = next()
            }
        }

        return obj
    }

    private func parseArray(_ openLine: Int, _ path: String, _ preTag: Tag? = nil) throws -> MMArray {
        depth += 1
        if depth > maxDepth {
            throw JSONCParserError.maxDepthExceeded
        }

        defer { depth -= 1 }

        var tag = consumeCommentsFor(openLine)
        if let pt = preTag {
            if tag == nil {
                tag = Tag()
            }
            tag!.inherit(from: pt)
        }
        if let tag = tag {
            if tag.type == .unknown {
                if tag.size > 0 {
                    tag.type = .arr
                } else {
                    tag.type = .vec
                }
            }
        }

        let arr = MMArray(tag: tag, path: path)

        if let tag = tag, !tag.example {
            let arrayResult = validator.validate(arr, tag: tag)
            if !arrayResult.isValid {
                throw JSONCParserError.invalidData(arrayResult.errors.joined(separator: ", "))
            }
        }

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

            if tok.type == .comment {
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

            let itemPath = "\(path)[\(index)]"
            if let parentTag = tag {
                let elemOwnTag = consumeCommentsFor(tok.line)
                let childTag = elemOwnTag ?? Tag()
                childTag.inherit(from: parentTag)

                if let item = try parseNode(itemPath, childTag) {
                    if let value = item as? Value, childTag.type == .bigint, value.data == nil {
                        value.data = value.text
                        let bigIntResult = validator.validate(value.data!, tag: childTag)
                        if !bigIntResult.isValid {
                            throw JSONCParserError.invalidData(bigIntResult.errors.joined(separator: ", "))
                        }
                    }
                    arr.items.append(item)
                    index += 1
                }
            } else {
                if let item = try parseNode(itemPath) {
                    arr.items.append(item)
                    index += 1
                }
            }

            if peek().type == .comma {
                _ = next()
            }
        }

        return arr
    }

    private func mergeNodeTag(_ node: Node, _ tag: Tag) {
        guard let existing = node.getTag() else { return }

        if node is Value {
            (node as? Value)?.tag = mergeTag(existing, tag)
        } else if node is MMObject {
            (node as? MMObject)?.tag = mergeTag(existing, tag)
        } else if node is MMArray {
            (node as? MMArray)?.tag = mergeTag(existing, tag)
        }
    }
}

public func parseJSONC(_ input: String) throws -> Node? {
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