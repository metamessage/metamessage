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

                    var parsedData: Any = text

                    switch tag.type {
                    case .str:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid string: \(text), valid: \"\"")
                            }
                            parsedData = ""
                        } else {
                            let strResult = validator.validate(text, tag: tag)
                            if !strResult.isValid {
                                throw JSONCParserError.invalidData(strResult.errors.joined(separator: ", "))
                            }
                        }

                    case .bytes:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid bytes: \(text), valid: \"\"")
                            }
                            parsedData = Data()
                        } else {
                            guard let data = Data(base64Encoded: text) else {
                                throw JSONCParserError.invalidData("invalid base64 bytes: \(text)")
                            }
                            let bytesResult = validator.validate(data, tag: tag)
                            if !bytesResult.isValid {
                                throw JSONCParserError.invalidData(bytesResult.errors.joined(separator: ", "))
                            }
                            parsedData = data
                        }

                    case .datetime:
                        if tag.isNull {
                            let formatter = DateFormatter()
                            formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                            formatter.timeZone = TimeZone(abbreviation: "UTC")
                            let defaultTime = formatter.string(from: Date(timeIntervalSince1970: 0))
                            if text != defaultTime {
                                throw JSONCParserError.invalidData("invalid datetime: \(text), valid: \(defaultTime)")
                            }
                            parsedData = Date(timeIntervalSince1970: 0)
                        } else {
                            let formatter = DateFormatter()
                            formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                            formatter.timeZone = TimeZone(abbreviation: "UTC")
                            if let date = formatter.date(from: text) {
                                let dtResult = validator.validate(date, tag: tag)
                                if !dtResult.isValid {
                                    throw JSONCParserError.invalidData(dtResult.errors.joined(separator: ", "))
                                }
                                parsedData = date
                            } else {
                                throw JSONCParserError.invalidData("invalid datetime: \(text)")
                            }
                        }

                    case .date:
                        if tag.isNull {
                            let formatter = DateFormatter()
                            formatter.dateFormat = "yyyy-MM-dd"
                            formatter.timeZone = TimeZone(abbreviation: "UTC")
                            let defaultDate = formatter.string(from: Date(timeIntervalSince1970: 0))
                            if text != defaultDate {
                                throw JSONCParserError.invalidData("invalid date: \(text), valid: \(defaultDate)")
                            }
                            parsedData = Date(timeIntervalSince1970: 0)
                        } else {
                            let formatter = DateFormatter()
                            formatter.dateFormat = "yyyy-MM-dd"
                            formatter.timeZone = TimeZone(abbreviation: "UTC")
                            if let date = formatter.date(from: text) {
                                let dResult = validator.validate(date, tag: tag)
                                if !dResult.isValid {
                                    throw JSONCParserError.invalidData(dResult.errors.joined(separator: ", "))
                                }
                                parsedData = date
                            } else {
                                throw JSONCParserError.invalidData("invalid date: \(text)")
                            }
                        }

                    case .time:
                        if tag.isNull {
                            let formatter = DateFormatter()
                            formatter.dateFormat = "HH:mm:ss"
                            formatter.timeZone = TimeZone(abbreviation: "UTC")
                            let defaultTime = formatter.string(from: Date(timeIntervalSince1970: 0))
                            if text != defaultTime {
                                throw JSONCParserError.invalidData("invalid time: \(text), valid: \(defaultTime)")
                            }
                            parsedData = Date(timeIntervalSince1970: 0)
                        } else {
                            let formatter = DateFormatter()
                            formatter.dateFormat = "HH:mm:ss"
                            formatter.timeZone = TimeZone(abbreviation: "UTC")
                            if let date = formatter.date(from: text) {
                                let tResult = validator.validate(date, tag: tag)
                                if !tResult.isValid {
                                    throw JSONCParserError.invalidData(tResult.errors.joined(separator: ", "))
                                }
                                parsedData = date
                            } else {
                                throw JSONCParserError.invalidData("invalid time: \(text)")
                            }
                        }

                    case .uuid:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid uuid: \(text), valid: \"\"")
                            }
                            parsedData = [UInt8](repeating: 0, count: 16)
                        } else {
                            let hexStr = text.replacingOccurrences(of: "-", with: "")
                            var uuidBytes = [UInt8]()
                            var index = hexStr.startIndex
                            while index < hexStr.endIndex {
                                let next = hexStr.index(index, offsetBy: 2, limitedBy: hexStr.endIndex) ?? hexStr.endIndex
                                if let byte = UInt8(hexStr[index..<next], radix: 16) {
                                    uuidBytes.append(byte)
                                }
                                index = next
                            }
                            if uuidBytes.count == 16 {
                                parsedData = uuidBytes
                            }
                            let uuidResult = validator.validate(parsedData, tag: tag)
                            if !uuidResult.isValid {
                                throw JSONCParserError.invalidData(uuidResult.errors.joined(separator: ", "))
                            }
                        }

                    case .decimal:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid decimal: \(text), valid: \"\"")
                            }
                            parsedData = ""
                        } else {
                            let decResult = validator.validate(text, tag: tag)
                            if !decResult.isValid {
                                throw JSONCParserError.invalidData(decResult.errors.joined(separator: ", "))
                            }
                        }

                    case .ip:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid ip: \(text), valid: \"\"")
                            }
                            parsedData = ""
                        } else {
                            let ipResult = validator.validate(text, tag: tag)
                            if !ipResult.isValid {
                                throw JSONCParserError.invalidData(ipResult.errors.joined(separator: ", "))
                            }
                        }

                    case .url:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid url: \(text), valid: \"\"")
                            }
                            parsedData = ""
                        } else {
                            let urlResult = validator.validate(text, tag: tag)
                            if !urlResult.isValid {
                                throw JSONCParserError.invalidData(urlResult.errors.joined(separator: ", "))
                            }
                        }

                    case .email:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid email: \(text), valid: \"\"")
                            }
                            parsedData = ""
                        } else {
                            let emailResult = validator.validate(text, tag: tag)
                            if !emailResult.isValid {
                                throw JSONCParserError.invalidData(emailResult.errors.joined(separator: ", "))
                            }
                        }

                    case .enums:
                        if tag.enums.isEmpty {
                            throw JSONCParserError.invalidData("enum empty")
                        }
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid enums: \(text), valid: \"\"")
                            }
                            parsedData = -1
                        } else {
                            let enumResult = validator.validate(text, tag: tag)
                            if !enumResult.isValid {
                                throw JSONCParserError.invalidData(enumResult.errors.joined(separator: ", "))
                            }
                        }

                    case .media:
                        if tag.isNull {
                            if text != "" {
                                throw JSONCParserError.invalidData("invalid media: \(text), valid: \"\"")
                            }
                            parsedData = Data()
                        } else {
                            guard let data = Data(base64Encoded: text) else {
                                throw JSONCParserError.invalidData("invalid base64 media: \(text)")
                            }
                            parsedData = data
                        }

                    default:
                        let strResult = validator.validate(text, tag: tag)
                        if !strResult.isValid {
                            throw JSONCParserError.invalidData(strResult.errors.joined(separator: ", "))
                        }
                    }

                    let value = NodeScalar(data: parsedData, text: text, tag: tag, path: path)
                    return value
                }

                let value = NodeScalar(data: text, text: text, tag: tag, path: path)
                return value

        case .number:
            let tag = preTag ?? consumeCommentsFor(tok.line)
            let text = tok.literal

            if let tag = tag {
                if tag.type == .unknown {
                    if text.contains(".") {
                        tag.type = .f64
                    } else if text.hasPrefix("-") {
                        tag.type = .i
                    } else {
                        tag.type = .i
                    }
                }

                if tag.isNull {
                    if text.contains(".") {
                        if text != "0.0" {
                            throw JSONCParserError.invalidData("invalid float: \(text), valid: 0.0")
                        }
                        switch tag.type {
                        case .f32:
                            return NodeScalar(data: Float(0.0), text: text, tag: tag, path: path)
                        case .decimal:
                            return NodeScalar(data: "", text: text, tag: tag, path: path)
                        default:
                            return NodeScalar(data: Double(0.0), text: text, tag: tag, path: path)
                        }
                    } else {
                        if text != "0" {
                            throw JSONCParserError.invalidData("invalid int: \(text), valid: 0")
                        }
                        if tag.type == .bigint {
                            return NodeScalar(data: "0", text: text, tag: tag, path: path)
                        }
                        return NodeScalar(data: Int(0), text: text, tag: tag, path: path)
                    }
                }
            }

            var data: Any?

            if text.contains(".") {
                if let tag = tag {
                    switch tag.type {
                    case .f32:
                        if let f = Float(text) {
                            data = f
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .decimal:
                        data = text
                        let numResult = validator.validate(data!, tag: tag)
                        if !numResult.isValid {
                            throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                        }
                    default:
                        if let d = Double(text) {
                            data = d
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    }
                } else {
                    data = Double(text)
                }
            } else if text.hasPrefix("-") {
                if let tag = tag {
                    switch tag.type {
                    case .i8:
                        if let v = Int8(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i16:
                        if let v = Int16(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i32:
                        if let v = Int32(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i64:
                        if let v = Int64(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .bigint:
                        data = text
                        let numResult = validator.validate(data!, tag: tag)
                        if !numResult.isValid {
                            throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                        }
                    default:
                        if let v = Int(text) {
                            data = v
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        } else if let v = Int64(text) {
                            data = v
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    }
                } else {
                    if let ival = Int(text) {
                        data = ival
                    } else {
                        data = Int64(text)
                    }
                }
            } else {
                if let tag = tag {
                    switch tag.type {
                    case .u:
                        if let v = UInt(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .u8:
                        if let v = UInt8(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .u16:
                        if let v = UInt16(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .u32:
                        if let v = UInt32(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .u64:
                        if let v = UInt64(text) {
                            if v <= UInt64(Int.max) {
                                data = Int(v)
                            } else {
                                data = v
                            }
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i8:
                        if let v = Int8(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i16:
                        if let v = Int16(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i32:
                        if let v = Int32(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .i64:
                        if let v = Int64(text) {
                            data = Int(v)
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    case .bigint:
                        data = text
                        let numResult = validator.validate(data!, tag: tag)
                        if !numResult.isValid {
                            throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                        }
                    default:
                        if let uval = UInt64(text) {
                            if uval > UInt64(Int.max) {
                                data = uval
                            } else {
                                data = Int(uval)
                            }
                            let numResult = validator.validate(data!, tag: tag)
                            if !numResult.isValid {
                                throw JSONCParserError.invalidData(numResult.errors.joined(separator: ", "))
                            }
                        }
                    }
                } else {
                    if let uval = UInt64(text) {
                        if uval > UInt64(Int.max) {
                            data = uval
                        } else {
                            data = Int(uval)
                        }
                    }
                }
            }

            let value = NodeScalar(data: data, text: text, tag: tag, path: path)
            return value

        case .trueValue:
            let tag = preTag ?? consumeCommentsFor(tok.line)
            if let tag = tag {
                if tag.type == .unknown {
                    tag.type = .bool
                }
                switch tag.type {
                case .bool:
                    if tag.isNull {
                        throw JSONCParserError.invalidData("bool must false when bool is null")
                    }
                    let trueResult = validator.validate(true, tag: tag)
                    if !trueResult.isValid {
                        throw JSONCParserError.invalidData(trueResult.errors.joined(separator: ", "))
                    }
                default:
                    throw JSONCParserError.invalidData("unsupported type \(tag.type.stringValue) for boolean literal")
                }
            }
            let value = NodeScalar(data: true, text: "true", tag: tag, path: path)
            return value

        case .falseValue:
            let tag = preTag ?? consumeCommentsFor(tok.line)
            if let tag = tag {
                if tag.type == .unknown {
                    tag.type = .bool
                }
                switch tag.type {
                case .bool:
                    if tag.isNull {
                        return NodeScalar(data: false, text: "false", tag: tag, path: path)
                    }
                    let falseResult = validator.validate(false, tag: tag)
                    if !falseResult.isValid {
                        throw JSONCParserError.invalidData(falseResult.errors.joined(separator: ", "))
                    }
                default:
                    throw JSONCParserError.invalidData("unsupported type \(tag.type.stringValue) for boolean literal")
                }
            }
            let value = NodeScalar(data: false, text: "false", tag: tag, path: path)
            return value

        case .nullValue:
            throw JSONCParserError.invalidData("null is not supported")

        default:
            throw JSONCParserError.unexpectedToken("Unexpected token: \(tok.type)")
        }
    }

    private func parseObject(_ openLine: Int, _ path: String, _ preTag: Tag? = nil) throws -> NodeObject {
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

        var localPath = path
        if let tag = tag, !tag.name.isEmpty {
            if localPath.isEmpty {
                localPath = tag.name
            } else {
                localPath = "\(localPath).\(tag.name)"
            }
        }

        let obj = NodeObject(tag: tag, path: localPath)

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

            let key = camelToSnake(keyTok.literal)

            let tok2 = peek()
            _ = next()

            if let parentTag = tag {
                let ownTag = consumeCommentsFor(tok2.line)
                let childTag = ownTag ?? Tag()

                if parentTag.type == .map {
                    childTag.inherit(from: parentTag)
                    if childTag.example {
                        parentTag.isEmpty = true
                    }
                }

                let childPath: String
                if parentTag.type == .map {
                    childPath = "\(localPath)[\(key)]"
                } else {
                    childPath = "\(localPath).\(key)"
                }
                if let val = try parseNode(childPath, childTag) {
                    let field = Field(key: key, value: val)
                    obj.fields.append(field)
                }
            } else {
                let childPath = "\(localPath).\(key)"
                if let val = try parseNode(childPath) {
                    let field = Field(key: key, value: val)
                    obj.fields.append(field)
                }
            }

            if peek().type == .comma {
                _ = next()
            }
        }

        if let tag = tag, !tag.example {
            if tag.type == .obj || tag.type == .map {
                let structResult = validator.validate(obj, tag: tag)
                if !structResult.isValid {
                    throw JSONCParserError.invalidData(structResult.errors.joined(separator: ", "))
                }
            }
        }

        return obj
    }

    private func parseArray(_ openLine: Int, _ path: String, _ preTag: Tag? = nil) throws -> NodeArray {
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
        // Reset the container's type so it gets typed as vec/arr instead of
        // inheriting childType from the field's tag (which is meant for items).
        tag?.type = .unknown
        if let tag = tag {
            if tag.type == .unknown {
                tag.type = .vec
            }
        }

        var localPath = path
        if let tag = tag, !tag.name.isEmpty {
            if localPath.isEmpty {
                localPath = tag.name
            } else {
                localPath = "\(localPath).\(tag.name)"
            }
        }

        let arr = NodeArray(tag: tag, path: localPath)

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

            let itemPath = "\(localPath)[\(index)]"
            if let parentTag = tag {
                let elemOwnTag: Tag?
                if openLine != tok.line {
                    elemOwnTag = consumeCommentsFor(tok.line)
                } else {
                    elemOwnTag = nil
                }
                let childTag = elemOwnTag ?? Tag()
                // Items inherit from the field's original tag (preTag) which
                // carries child_* attributes, falling back to the array tag.
                childTag.inherit(from: preTag ?? parentTag)
                if childTag.example {
                    parentTag.isEmpty = true
                }

                if let item = try parseNode(itemPath, childTag) {
                    if let value = item as? NodeScalar, childTag.type == .bigint, value.data == nil {
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

        if node is NodeScalar {
            (node as? NodeScalar)?.tag = mergeTag(existing, tag)
        } else if node is NodeObject {
            (node as? NodeObject)?.tag = mergeTag(existing, tag)
        } else if node is NodeArray {
            (node as? NodeArray)?.tag = mergeTag(existing, tag)
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