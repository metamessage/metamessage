import Foundation

public class NodeDecoder {
    private let buffer: MMBuffer

    public init(data: Data) {
        self.buffer = MMBuffer(data: data)
    }

    public func decode() throws -> Node {
        return try decodeNode(tag: nil, path: "$")
    }

    private func decodeNode(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.peek() else {
            throw MMError.unexpectedEndOfData
        }

        guard let prefix = getPrefix(byte) else {
            throw MMError.invalidPrefix
        }

        switch prefix {
        case .simple:
            return try decodeSimple(tag: tag, path: path)
        case .positiveInt:
            return try decodePositiveInt(tag: tag, path: path)
        case .negativeInt:
            return try decodeNegativeInt(tag: tag, path: path)
        case .prefixFloat:
            return try decodeFloat(tag: tag, path: path)
        case .prefixString:
            return try decodeString(tag: tag, path: path)
        case .prefixBytes:
            return try decodeBytes(tag: tag, path: path)
        case .container:
            return try decodeContainer(tag: tag, path: path)
        case .prefixTag:
            return try decodeTag(path: path)
        }
    }

    // MARK: - Tag decode

    private func decodeTag(path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = tagLen(byte)
        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let data = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        var pos = 0
        guard pos < data.count else {
            throw MMError.unexpectedEndOfData
        }

        let tagBodyLenByte = data[pos]
        pos += 1
        var tagBodyLen = Int(tagBodyLenByte)
        if tagBodyLenByte >= 254 {
            if tagBodyLenByte == 254 {
                guard pos < data.count else { throw MMError.unexpectedEndOfData }
                tagBodyLen = Int(data[pos])
                pos += 1
            } else {
                guard pos + 1 < data.count else { throw MMError.unexpectedEndOfData }
                tagBodyLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                pos += 2
            }
        }

        let tag = Tag()
        let tagFieldEnd = pos + tagBodyLen
        while pos < tagFieldEnd {
            guard pos < data.count else { break }
            let fieldByte = data[pos]
            let key = fieldByte & 0xF8
            let fieldLen = Int(fieldByte) & 0x07
            pos += 1

            switch key {
            case TagKey.type:
                guard pos < data.count else { break }
                tag.type = ValueType(rawValue: data[pos]) ?? .unknown
                pos += 1

            case TagKey.isNull:
                tag.isNull = (fieldLen & 0x01) == 1
                if tag.isNull {
                    tag.nullable = true
                }

            case TagKey.deprecated:
                tag.deprecated = (fieldLen & 0x01) == 1

            case TagKey.nullable:
                tag.nullable = (fieldLen & 0x01) == 1

            case TagKey.allowEmpty:
                tag.allowEmpty = (fieldLen & 0x01) == 1

            case TagKey.unique:
                tag.unique = (fieldLen & 0x01) == 1

            case TagKey.example:
                tag.example = (fieldLen & 0x01) == 1

            case TagKey.location:
                if fieldLen <= 5 {
                    guard pos + fieldLen <= data.count else { break }
                    let valueBytes = data[pos..<pos+fieldLen]
                    if let locStr = String(bytes: valueBytes, encoding: .utf8), let loc = Int(locStr) {
                        tag.location = loc
                    }
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1
                    guard pos + strLen <= data.count else { break }
                    let valueBytes = data[pos..<pos+strLen]
                    if let locStr = String(bytes: valueBytes, encoding: .utf8), let loc = Int(locStr) {
                        tag.location = loc
                    }
                    pos += strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2
                    guard pos + strLen <= data.count else { break }
                    let valueBytes = data[pos..<pos+strLen]
                    if let locStr = String(bytes: valueBytes, encoding: .utf8), let loc = Int(locStr) {
                        tag.location = loc
                    }
                    pos += strLen
                }

            case TagKey.enums:
                if fieldLen <= 5 {
                    guard pos + fieldLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+fieldLen], encoding: .utf8) {
                        tag.enums = enumStr
                        tag.type = .enums
                    }
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        tag.enums = enumStr
                        tag.type = .enums
                    }
                    pos += strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        tag.enums = enumStr
                        tag.type = .enums
                    }
                    pos += strLen
                }

            case TagKey.mime:
                let (mimeVal, mimeAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if mimeAdv > 0 {
                    tag.mime = Mime.toString(mimeVal)
                    tag.type = .media
                    pos += mimeAdv
                } else {
                    break
                }

            case TagKey.childMime:
                let (childMimeVal, childMimeAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if childMimeAdv > 0 {
                    tag.childMime = Mime.toString(childMimeVal)
                    tag.childType = .media
                    pos += childMimeAdv
                } else {
                    break
                }

            case TagKey.childType:
                guard pos < data.count else { break }
                tag.childType = ValueType(rawValue: data[pos]) ?? .unknown
                pos += 1

            case TagKey.childEnums:
                if fieldLen <= 5 {
                    guard pos + fieldLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+fieldLen], encoding: .utf8) {
                        tag.childEnums = enumStr
                        tag.childType = .enums
                    }
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        tag.childEnums = enumStr
                        tag.childType = .enums
                    }
                    pos += strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        tag.childEnums = enumStr
                        tag.childType = .enums
                    }
                    pos += strLen
                }

            case TagKey.size:
                let (sizeVal, sizeAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if sizeAdv > 0 {
                    tag.size = Int(sizeVal)
                    pos += sizeAdv
                } else {
                    break
                }

            case TagKey.version:
                let (verVal, verAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if verAdv > 0 {
                    tag.version = Int(verVal)
                    pos += verAdv
                } else {
                    break
                }

            case TagKey.more:
                let (moreVal, moreAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if moreAdv > 0 {
                    tag.more = Int(moreVal)
                    pos += moreAdv
                } else {
                    break
                }

            case TagKey.childSize:
                let (csVal, csAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if csAdv > 0 {
                    tag.childSize = Int(csVal)
                    pos += csAdv
                } else {
                    break
                }

            case TagKey.childVersion:
                let (cvVal, cvAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if cvAdv > 0 {
                    tag.childVersion = Int(cvVal)
                    pos += cvAdv
                } else {
                    break
                }

            case TagKey.desc:
                let (descStr, descAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if descAdv > 0 {
                    tag.desc = descStr
                    pos += descAdv
                }

            case TagKey.defaultVal:
                let (dvStr, dvAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if dvAdv > 0 {
                    tag.defaultVal = dvStr
                    pos += dvAdv
                }

            case TagKey.min:
                let (minStr, minAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if minAdv > 0 {
                    tag.min = minStr
                    pos += minAdv
                }

            case TagKey.max:
                let (maxStr, maxAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if maxAdv > 0 {
                    tag.max = maxStr
                    pos += maxAdv
                }

            case TagKey.childDesc:
                let (cdStr, cdAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if cdAdv > 0 {
                    tag.childDesc = cdStr
                    pos += cdAdv
                }

            case TagKey.childNullable:
                tag.childNullable = (fieldLen & 0x01) == 1

            case TagKey.childAllowEmpty:
                tag.childAllowEmpty = (fieldLen & 0x01) == 1

            case TagKey.childUnique:
                tag.childUnique = (fieldLen & 0x01) == 1

            case TagKey.childDefaultVal:
                let (cdvStr, cdvAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if cdvAdv > 0 {
                    tag.childDefaultVal = cdvStr
                    pos += cdvAdv
                }

            case TagKey.childMin:
                let (cminStr, cminAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if cminAdv > 0 {
                    tag.childMin = cminStr
                    pos += cminAdv
                }

            case TagKey.childMax:
                let (cmaxStr, cmaxAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if cmaxAdv > 0 {
                    tag.childMax = cmaxStr
                    pos += cmaxAdv
                }

            case TagKey.childPattern:
                let (cpStr, cpAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if cpAdv > 0 {
                    tag.childPattern = cpStr
                    pos += cpAdv
                }

            case TagKey.childLocation:
                let (clStr, clAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if clAdv > 0, let loc = Int(clStr) {
                    tag.childLocation = loc
                    pos += clAdv
                }

            case TagKey.pattern:
                let (pStr, pAdv) = decodeTagString(data, pos: pos, fieldLen: fieldLen)
                if pAdv > 0 {
                    tag.pattern = pStr
                    pos += pAdv
                }

            default:
                let isU64Field = key == TagKey.size || key == TagKey.version ||
                                 key == TagKey.more || key == TagKey.childSize ||
                                 key == TagKey.childVersion
                if isU64Field {
                    pos += fieldLen + 1
                } else if fieldLen <= 5 {
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1 + strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2 + strLen
                }
            }
        }

        pos = tagFieldEnd

        guard pos < data.count else {
            // Tag with no payload — return a null value or the tag itself
            if tag.isNull {
                return tagNullNode(tag, path: path)
            }
            throw MMError.unexpectedEndOfData
        }

        let payloadData = Data(data[pos...])
        let innerDecoder = NodeDecoder(data: payloadData)
        let value = try innerDecoder.decodeNode(tag: tag, path: path)

        // Attach the tag to the decoded node
        if let val = value as? NodeScalar {
            val.tag = tag
            val.path = path
            return val
        } else if let arr = value as? NodeArray {
            arr.tag = tag
            arr.path = path
            return arr
        } else if let obj = value as? NodeObject {
            obj.tag = tag
            obj.path = path
            return obj
        }

        return value
    }

    private func tagNullNode(_ tag: Tag, path: String) -> Node {
        switch tag.type {
        case .i, .i8, .i16, .i32, .i64:
            return NodeScalar(data: 0, text: "0", tag: tag, path: path)
        case .u, .u8, .u16, .u32, .u64:
            return NodeScalar(data: UInt(0), text: "0", tag: tag, path: path)
        case .f32, .f64:
            return NodeScalar(data: 0.0, text: "0.0", tag: tag, path: path)
        case .bool:
            return NodeScalar(data: false, text: "false", tag: tag, path: path)
        case .str, .email, .url, .datetime, .date, .time, .enums, .media:
            return NodeScalar(data: "", text: "", tag: tag, path: path)
        case .bytes, .uuid, .ip:
            return NodeScalar(data: Data(), text: "", tag: tag, path: path)
        default:
            return NodeScalar(data: nil, text: "null", tag: tag, path: path)
        }
    }

    // MARK: - Simple

    private func decodeSimple(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let resolvedTag = tag ?? Tag()

        guard let simpleValue = MMSimpleValue(rawValue: byte) else {
            throw MMError.invalidData
        }

        // Handle null before modifying tag type — null must keep .unknown type
        if simpleValue == .null {
            if resolvedTag.type != .unknown {
                throw MMError.invalidData
            }
            return NodeNull(tag: resolvedTag, path: path)
        }

        if resolvedTag.type == .unknown {
            resolvedTag.type = .bool
        }

        switch simpleValue {
        case .null:
            // Unreachable — handled above
            throw MMError.invalidData
        case .trueValue:
            return NodeScalar(data: true, text: "true", tag: resolvedTag, path: path)
        case .falseValue:
            return NodeScalar(data: false, text: "false", tag: resolvedTag, path: path)
        case .nullBool:
            return NodeScalar(data: false, text: "false", tag: resolvedTag, path: path)
        case .nullInt:
            return NodeScalar(data: 0, text: "0", tag: resolvedTag, path: path)
        case .nullFloat:
            return NodeScalar(data: 0.0, text: "0.0", tag: resolvedTag, path: path)
        case .nullString:
            return NodeScalar(data: "", text: "", tag: resolvedTag, path: path)
        case .nullBytes:
            return NodeScalar(data: Data(), text: "", tag: resolvedTag, path: path)
        default:
            // Handle null-like values via the tag
            if let tag = tag, tag.isNull {
                return NodeScalar(data: nil, text: "null", tag: resolvedTag, path: path)
            }
            // Fall through to simple string codes
            // Simple string codes (code, message, data, success, error, etc.)
            let text: String
            switch simpleValue {
            case .code: text = "code"
            case .message: text = "message"
            case .data: text = "data"
            case .success: text = "success"
            case .error: text = "error"
            case .unknown: text = "unknown"
            case .page: text = "page"
            case .limit: text = "limit"
            case .offset: text = "offset"
            case .total: text = "total"
            case .id: text = "id"
            case .name: text = "name"
            case .description: text = "description"
            case .typeValue: text = "type"
            case .version: text = "version"
            case .status: text = "status"
            case .url: text = "url"
            case .createTime: text = "create_time"
            case .updateTime: text = "update_time"
            case .deleteTime: text = "delete_time"
            case .account: text = "account"
            case .token: text = "token"
            case .expireTime: text = "expire_time"
            case .key: text = "key"
            default: text = "null"
            }
            resolvedTag.type = .str
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)
        }
    }

    // MARK: - Positive Int

    private func decodePositiveInt(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = intLen(byte)

        var value: UInt64 = 0
        for _ in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            value = (value << 8) | UInt64(b)
        }

        if extraBytes == 0 {
            let suffix = Int(byte & 0x1F)
            value = UInt64(suffix)
        }

        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            resolvedTag.type = .i
        }

        return makeIntNode(value: value, negative: false, tag: resolvedTag, path: path)
    }

    // MARK: - Negative Int

    private func decodeNegativeInt(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = intLen(byte)

        var value: UInt64 = 0
        for _ in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            value = (value << 8) | UInt64(b)
        }

        if extraBytes == 0 {
            let suffix = Int(byte & 0x1F)
            value = UInt64(suffix)
        }

        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            resolvedTag.type = .i
        }

        return makeIntNode(value: value, negative: true, tag: resolvedTag, path: path)
    }

    private func makeIntNode(value: UInt64, negative: Bool, tag: Tag, path: String) -> Node {
        switch tag.type {
        case .i:
            let v: Int
            if negative {
                if value == 9223372036854775808 {
                    v = Int(Int64.min)
                } else {
                    v = -Int(value)
                }
            } else {
                v = Int(value)
            }
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .i8:
            let v: Int8 = negative ? -Int8(truncatingIfNeeded: value) : Int8(truncatingIfNeeded: value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .i16:
            let v: Int16 = negative ? -Int16(truncatingIfNeeded: value) : Int16(truncatingIfNeeded: value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .i32:
            let v: Int32 = negative ? -Int32(truncatingIfNeeded: value) : Int32(truncatingIfNeeded: value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .i64:
            let v: Int64
            if negative {
                if value == 9223372036854775808 {
                    v = Int64.min
                } else {
                    v = -Int64(value)
                }
            } else {
                v = Int64(value)
            }
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .u:
            let v = UInt(value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .u8:
            let v = UInt8(truncatingIfNeeded: value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .u16:
            let v = UInt16(truncatingIfNeeded: value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .u32:
            let v = UInt32(truncatingIfNeeded: value)
            return NodeScalar(data: v, text: "\(v)", tag: tag, path: path)

        case .u64:
            return NodeScalar(data: value, text: "\(value)", tag: tag, path: path)

        case .bool:
            let v = value != 0
            return NodeScalar(data: v, text: v ? "true" : "false", tag: tag, path: path)

        case .datetime:
            let ts: Int64
            if negative {
                if value == 9223372036854775808 {
                    ts = Int64.min
                } else {
                    ts = -Int64(value)
                }
            } else {
                ts = Int64(value)
            }
            let adjustedTs = TimeInterval(ts) + TimeInterval(tag.location * 3600)
            let date = Date(timeIntervalSince1970: adjustedTs)
            let formatter = DateFormatter()
            formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
            formatter.timeZone = TimeZone(abbreviation: "UTC")
            let text = formatter.string(from: date)
            return NodeScalar(data: date, text: text, tag: tag, path: path)

        case .date:
            var days = Int(value)
            if negative { days = -days }
            let ref = Date(timeIntervalSince1970: 0)
            let date = Calendar(identifier: .gregorian).date(byAdding: .day, value: days, to: ref) ?? ref
            let formatter = DateFormatter()
            formatter.dateFormat = "yyyy-MM-dd"
            formatter.timeZone = TimeZone(abbreviation: "UTC")
            let text = formatter.string(from: date)
            return NodeScalar(data: date, text: text, tag: tag, path: path)

        case .time:
            var secs = Int(value)
            if negative { secs = -secs }
            secs = max(0, min(secs, 86399))
            let h = secs / 3600
            let m = (secs % 3600) / 60
            let s = secs % 60
            let text = String(format: "%02d:%02d:%02d", h, m, s)
            return NodeScalar(data: text, text: text, tag: tag, path: path)

        case .enums:
            if !tag.enums.isEmpty {
                let enumValues = tag.enums.split(separator: "|").map { $0.trimmingCharacters(in: .whitespaces) }
                let idx = Int(value)
                if idx >= 0 && idx < enumValues.count {
                    let text = String(enumValues[idx])
                    return NodeScalar(data: text, text: text, tag: tag, path: path)
                }
            }
            return NodeScalar(data: String(value), text: String(value), tag: tag, path: path)

        default:
            if negative {
                return NodeScalar(data: -Int64(value), text: "-\(value)", tag: tag, path: path)
            }
            return NodeScalar(data: Int64(value), text: "\(value)", tag: tag, path: path)
        }
    }

    // MARK: - Float

    private func decodeFloat(tag: Tag?, path: String) throws -> Node {
        guard let prefix = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let prefixLower = prefix & 0x0F
        var v: Double = 0

        if prefixLower <= 6 {
            v = Double(prefixLower) / 10.0
            if (prefix & MMConstants.floatPositiveNegativeMask) != 0 {
                v = -v
            }
        } else {
            guard let expByte = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            let exponent = Int8(bitPattern: expByte)

            let (extraBytes, _) = floatLen(prefix)
            var mantissa: UInt64 = 0
            for _ in 0..<extraBytes {
                guard let b = buffer.read() else {
                    throw MMError.unexpectedEndOfData
                }
                mantissa = (mantissa << 8) | UInt64(b)
            }

            let decimalStr = mantissaToDecimal(mantissa, exponent)
            guard let parsed = Double(decimalStr) else {
                throw MMError.invalidData
            }
            v = parsed

            if (prefix & MMConstants.floatPositiveNegativeMask) != 0 {
                v = -v
            }
        }

        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            resolvedTag.type = .f64
        }

        switch resolvedTag.type {
        case .f32:
            let f = Float(v)
            return NodeScalar(data: f, text: formatFloat(f), tag: resolvedTag, path: path)
        case .f64:
            return NodeScalar(data: v, text: formatDouble(v), tag: resolvedTag, path: path)
        case .decimal:
            return NodeScalar(data: v, text: formatDouble(v), tag: resolvedTag, path: path)
        default:
            return NodeScalar(data: v, text: formatDouble(v), tag: resolvedTag, path: path)
        }
    }

    private func formatDouble(_ value: Double) -> String {
        if value == value.rounded() && value != 0 {
            return String(format: "%.1f", value)
        }
        // Use scientific if very small/large
        let s = String(value).lowercased()
        if s.contains("e") {
            return scientificToDecimal(s)
        }
        return s
    }

    private func formatFloat(_ value: Float) -> String {
        return formatDouble(Double(value))
    }

    private func scientificToDecimal(_ s: String) -> String {
        let parts = s.split(separator: "e", maxSplits: 1)
        guard parts.count == 2, let exponent = Int(parts[1]) else { return s }
        var mantissaStr = String(parts[0])
        if exponent >= 0 {
            let dotIdx = mantissaStr.firstIndex(of: ".") ?? mantissaStr.endIndex
            mantissaStr.remove(at: dotIdx == mantissaStr.endIndex ? mantissaStr.endIndex : dotIdx)
            if exponent + 1 >= mantissaStr.count {
                mantissaStr += String(repeating: "0", count: exponent + 1 - mantissaStr.count)
            } else {
                let insertIdx = mantissaStr.index(mantissaStr.startIndex, offsetBy: exponent + 1)
                mantissaStr.insert(".", at: insertIdx)
            }
        } else {
            let absExp = -exponent
            if let dotIdx = mantissaStr.firstIndex(of: ".") {
                mantissaStr.remove(at: dotIdx)
            }
            if absExp >= mantissaStr.count {
                mantissaStr = "0." + String(repeating: "0", count: absExp - mantissaStr.count) + mantissaStr
            } else {
                let insertIdx = mantissaStr.index(mantissaStr.startIndex, offsetBy: mantissaStr.count - absExp)
                mantissaStr.insert(".", at: insertIdx)
                if mantissaStr.hasPrefix(".") {
                    mantissaStr = "0" + mantissaStr
                }
            }
        }
        return mantissaStr
    }

    // MARK: - String

    private func decodeString(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = stringLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        var text = ""
        if totalLen > 0 {
            guard let bytes = buffer.read(totalLen) else {
                throw MMError.unexpectedEndOfData
            }
            guard let str = String(bytes: bytes, encoding: .utf8) else {
                throw MMError.invalidData
            }
            text = str
        }

        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            resolvedTag.type = .str
        }

        switch resolvedTag.type {
        case .str:
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)
        case .email:
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)
        case .url:
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)
        case .ip:
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)
        default:
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)
        }
    }

    // MARK: - Bytes

    private func decodeBytes(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = bytesLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        var dataVal = Data()
        if totalLen > 0 {
            guard let result = buffer.read(totalLen) else {
                throw MMError.unexpectedEndOfData
            }
            dataVal = Data(result)
        }

        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            resolvedTag.type = .bytes
        }

        switch resolvedTag.type {
        case .bytes:
            let text = dataVal.base64EncodedString()
            return NodeScalar(data: dataVal, text: text, tag: resolvedTag, path: path)

        case .media:
            let text = dataVal.base64EncodedString()
            return NodeScalar(data: dataVal, text: text, tag: resolvedTag, path: path)

        case .bigint:
            let bytes = [UInt8](dataVal)
            let text = NodeDecoder.decodeBigIntFromBytes(bytes)
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)

        case .uuid:
            guard dataVal.count == 16 else {
                throw MMError.invalidData
            }
            let bytes = [UInt8](dataVal)
            let text = String(format: "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
                bytes[0], bytes[1], bytes[2], bytes[3],
                bytes[4], bytes[5], bytes[6], bytes[7],
                bytes[8], bytes[9], bytes[10], bytes[11],
                bytes[12], bytes[13], bytes[14], bytes[15])
            return NodeScalar(data: text, text: text, tag: resolvedTag, path: path)

        case .ip:
            let ipStr = dataVal.map { String($0) }.joined(separator: ".")
            return NodeScalar(data: dataVal, text: ipStr, tag: resolvedTag, path: path)

        default:
            let text = dataVal.base64EncodedString()
            return NodeScalar(data: dataVal, text: text, tag: resolvedTag, path: path)
        }
    }

    // MARK: - Container

    private func decodeContainer(tag: Tag?, path: String) throws -> Node {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let isArrayContainer = isArray(byte)
        let (extraBytes, len) = containerLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        if isArrayContainer {
            return try decodeArrayElements(totalLen: totalLen, tag: tag, path: path)
        } else {
            return try decodeObjectElements(totalLen: totalLen, tag: tag, path: path)
        }
    }

    private func decodeArrayElements(totalLen: Int, tag: Tag?, path: String) throws -> Node {
        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            if resolvedTag.size != 0 {
                resolvedTag.type = .arr
            } else {
                resolvedTag.type = .vec
            }
        }

        let arr = NodeArray(items: [], tag: resolvedTag, path: path)
        let startPos = buffer.position()
        let endPos = startPos + totalLen

        while buffer.position() < endPos {
            let elementTag = Tag()
            elementTag.inherit(from: resolvedTag)
            let elementPath = "\(path)[\(arr.items.count)]"
            let element = try decodeNode(tag: elementTag, path: elementPath)
            arr.items.append(element)
        }

        return arr
    }

    private func decodeObjectElements(totalLen: Int, tag: Tag?, path: String) throws -> Node {
        let resolvedTag = tag ?? Tag()
        if resolvedTag.type == .unknown {
            resolvedTag.type = .obj
        }

        let keyArrayValue = try decodeNode(tag: resolvedTag, path: path)
        guard let keyArray = keyArrayValue as? NodeArray else {
            throw MMError.invalidData
        }

        let obj = NodeObject(fields: [], tag: resolvedTag, path: path)

        for item in keyArray.items {
            guard let keyVal = item as? NodeScalar, let key = keyVal.data as? String else {
                continue
            }

            let fieldTag = Tag()
            fieldTag.inherit(from: resolvedTag)
            let fieldPath = "\(path).\(key)"
            let fieldValue = try decodeNode(tag: fieldTag, path: fieldPath)
            obj.fields.append(Field(key: key, value: fieldValue))
        }

        return obj
    }

    /// Decode a bigint from bit-packed bytes.
    /// Format: [digit_count (1 byte)] + [bit-packed data]
    /// Bit-packed: sign(1 bit) + groups of 10/7/4 bits for 3/2/1 decimal digits
    static func decodeBigIntFromBytes(_ bytes: [UInt8]) -> String {
        guard bytes.count > 1 else { return "0" }
        let digitLen = Int(bytes[0])
        guard digitLen > 0 else { return "0" }

        let bitData = Array(bytes[1...])
        let totalBits = bitData.count * 8
        var bitPos = 0

        func readBit() -> Int {
            guard bitPos < totalBits else { return 0 }
            let byteIdx = bitPos / 8
            let bitInByte = 7 - (bitPos % 8)
            let bit = (Int(bitData[byteIdx]) >> bitInByte) & 1
            bitPos += 1
            return bit
        }

        func readBits(_ n: Int) -> Int {
            var val = 0
            for _ in 0..<n {
                val = (val << 1) | readBit()
            }
            return val
        }

        let sign = readBit()
        let neg = sign == 1

        var result = ""
        var remaining = digitLen
        while remaining > 0 {
            if remaining >= 3 {
                let val = readBits(10)
                result += String(format: "%03d", val)
                remaining -= 3
            } else if remaining == 2 {
                let val = readBits(7)
                result += String(format: "%02d", val)
                remaining -= 2
            } else {
                let val = readBits(4)
                result += String(val)
                remaining -= 1
            }
        }

        // Trim leading zeros
        var trimmed = ""
        var started = false
        for c in result {
            if !started {
                if c != "0" {
                    started = true
                    trimmed.append(c)
                }
            } else {
                trimmed.append(c)
            }
        }
        let finalStr = trimmed.isEmpty ? "0" : trimmed

        if neg && finalStr != "0" {
            return "-" + finalStr
        }
        return finalStr
    }
}

// MARK: - Backward Compatibility Decoder

public class Decoder {
    private let buffer: MMBuffer

    public init(data: Data) {
        self.buffer = MMBuffer(data: data)
    }

    public enum DecodedValue {
        case bool(Bool)
        case int(Int64)
        case uint(UInt64)
        case float(Double)
        case string(String)
        case bigint(String)
        case data(Data)
        case array([DecodedValue])
        case object([String: DecodedValue])
        case null
    }

    public func decode() throws -> DecodedValue {
        return try decodeValue(tag: nil)
    }

    private func decodeValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.peek() else {
            throw MMError.unexpectedEndOfData
        }

        guard let prefix = getPrefix(byte) else {
            throw MMError.invalidPrefix
        }

        switch prefix {
        case .simple:
            return try decodeSimpleValue(tag: tag)
        case .positiveInt:
            return try decodePositiveIntValue(tag: tag)
        case .negativeInt:
            return try decodeNegativeIntValue(tag: tag)
        case .prefixFloat:
            return try decodeFloatValue(tag: tag)
        case .prefixString:
            return try decodeStringValue(tag: tag)
        case .prefixBytes:
            return try decodeBytesValue(tag: tag)
        case .container:
            return try decodeContainerValue(tag: tag)
        case .prefixTag:
            return try decodeTagValue(tag: tag)
        }
    }

    private func decodeSimpleValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        guard let simpleValue = MMSimpleValue(rawValue: byte) else {
            throw MMError.invalidData
        }

        switch simpleValue {
        case .trueValue: return .bool(true)
        case .falseValue: return .bool(false)
        case .nullBool: return .bool(false)
        case .nullInt: return .int(0)
        case .nullFloat: return .float(0.0)
        case .nullString: return .string("")
        case .nullBytes: return .data(Data())
        default: return .null
        }
    }

    private func decodePositiveIntValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = intLen(byte)

        var value: UInt64 = 0
        for _ in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            value = (value << 8) | UInt64(b)
        }

        if extraBytes == 0 {
            let suffix = Int(byte & 0x1F)
            value = UInt64(suffix)
        }

        if value > UInt64(Int64.max) {
            return .uint(value)
        }
        return .int(Int64(value))
    }

    private func decodeNegativeIntValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, _) = intLen(byte)

        var value: UInt64 = 0
        for _ in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            value = (value << 8) | UInt64(b)
        }

        if extraBytes == 0 {
            let suffix = Int(byte & 0x1F)
            value = UInt64(suffix)
        }

        if value > UInt64(Int64.max) {
            return .int(Int64.min)
        }
        return .int(-Int64(value))
    }

    private func decodeFloatValue(tag: Tag?) throws -> DecodedValue {
        guard let prefix = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let prefixLower = prefix & 0x0F

        if prefixLower <= 6 {
            let val = Double(prefixLower) / 10.0
            if (prefix & MMConstants.floatPositiveNegativeMask) != 0 {
                return .float(-val)
            }
            return .float(val)
        }

        guard let expByte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }
        let exponent = Int8(bitPattern: expByte)

        let (extraBytes, _) = floatLen(prefix)
        var mantissa: UInt64 = 0
        for _ in 0..<extraBytes {
            guard let b = buffer.read() else {
                throw MMError.unexpectedEndOfData
            }
            mantissa = (mantissa << 8) | UInt64(b)
        }

        let decimalStr = mantissaToDecimal(mantissa, exponent)
        guard let parsed = Double(decimalStr) else {
            throw MMError.invalidData
        }

        if (prefix & MMConstants.floatPositiveNegativeMask) != 0 {
            return .float(-parsed)
        }
        return .float(parsed)
    }

    private func decodeStringValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = stringLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let bytes = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        guard let string = String(bytes: bytes, encoding: .utf8) else {
            throw MMError.invalidData
        }

        return .string(string)
    }

    private func decodeBytesValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = bytesLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let result = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        return .data(Data(result))
    }

    private func decodeContainerValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let isArrayContainer = isArray(byte)
        let (extraBytes, len) = containerLen(byte)

        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        if isArrayContainer {
            var elements: [DecodedValue] = []
            let startPos = buffer.position()
            let endPos = startPos + totalLen
            while buffer.position() < endPos {
                let element = try decodeValue(tag: tag)
                elements.append(element)
            }
            return .array(elements)
        } else {
            var dict: [String: DecodedValue] = [:]

            let keyArrayValue = try decodeValue(tag: nil)
            guard case .array(let keyItems) = keyArrayValue else {
                throw MMError.invalidData
            }

            for item in keyItems {
                guard case .string(let key) = item else {
                    continue
                }
                let value = try decodeValue(tag: tag)
                dict[key] = value
            }

            return .object(dict)
        }
    }

    private func decodeTagValue(tag: Tag?) throws -> DecodedValue {
        guard let byte = buffer.read() else {
            throw MMError.unexpectedEndOfData
        }

        let (extraBytes, len) = tagLen(byte)
        var totalLen: Int
        if extraBytes > 0 {
            guard let bytes = buffer.read(extraBytes) else {
                throw MMError.unexpectedEndOfData
            }
            totalLen = 0
            for b in bytes {
                totalLen = (totalLen << 8) | Int(b)
            }
        } else {
            totalLen = len
        }

        guard let data = buffer.read(totalLen) else {
            throw MMError.unexpectedEndOfData
        }

        var pos = 0
        guard pos < data.count else {
            throw MMError.unexpectedEndOfData
        }

        let tagBodyLenByte = data[pos]
        pos += 1
        var tagBodyLen = Int(tagBodyLenByte)
        if tagBodyLenByte >= 254 {
            if tagBodyLenByte == 254 {
                guard pos < data.count else { throw MMError.unexpectedEndOfData }
                tagBodyLen = Int(data[pos])
                pos += 1
            } else {
                guard pos + 1 < data.count else { throw MMError.unexpectedEndOfData }
                tagBodyLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                pos += 2
            }
        }

        let decodedTag = Tag()
        let tagFieldEnd = pos + tagBodyLen
        while pos < tagFieldEnd {
            guard pos < data.count else { break }
            let fieldByte = data[pos]
            let key = fieldByte & 0xF8
            let fieldLen = Int(fieldByte) & 0x07
            pos += 1

            switch key {
            case TagKey.type:
                guard pos < data.count else { break }
                decodedTag.type = ValueType(rawValue: data[pos]) ?? .unknown
                pos += 1

            case TagKey.isNull:
                decodedTag.isNull = (fieldLen & 0x01) == 1
                if decodedTag.isNull {
                    decodedTag.nullable = true
                }

            case TagKey.nullable:
                decodedTag.nullable = (fieldLen & 0x01) == 1

            case TagKey.deprecated:
                decodedTag.deprecated = (fieldLen & 0x01) == 1

            case TagKey.allowEmpty:
                decodedTag.allowEmpty = (fieldLen & 0x01) == 1

            case TagKey.unique:
                decodedTag.unique = (fieldLen & 0x01) == 1

            case TagKey.example:
                decodedTag.example = (fieldLen & 0x01) == 1

            case TagKey.location:
                if fieldLen <= 5 {
                    guard pos + fieldLen <= data.count else { break }
                    let valueBytes = data[pos..<pos+fieldLen]
                    if let locStr = String(bytes: valueBytes, encoding: .utf8), let loc = Int(locStr) {
                        decodedTag.location = loc
                    }
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1
                    guard pos + strLen <= data.count else { break }
                    let valueBytes = data[pos..<pos+strLen]
                    if let locStr = String(bytes: valueBytes, encoding: .utf8), let loc = Int(locStr) {
                        decodedTag.location = loc
                    }
                    pos += strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2
                    guard pos + strLen <= data.count else { break }
                    let valueBytes = data[pos..<pos+strLen]
                    if let locStr = String(bytes: valueBytes, encoding: .utf8), let loc = Int(locStr) {
                        decodedTag.location = loc
                    }
                    pos += strLen
                }

            case TagKey.enums:
                if fieldLen <= 5 {
                    guard pos + fieldLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+fieldLen], encoding: .utf8) {
                        decodedTag.enums = enumStr
                        decodedTag.type = .enums
                    }
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        decodedTag.enums = enumStr
                        decodedTag.type = .enums
                    }
                    pos += strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        decodedTag.enums = enumStr
                        decodedTag.type = .enums
                    }
                    pos += strLen
                }

            case TagKey.mime:
                let (mimeVal, mimeAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if mimeAdv > 0 {
                    decodedTag.mime = Mime.toString(mimeVal)
                    decodedTag.type = .media
                    pos += mimeAdv
                }

            case TagKey.childType:
                guard pos < data.count else { break }
                decodedTag.childType = ValueType(rawValue: data[pos]) ?? .unknown
                pos += 1

            case TagKey.childEnums:
                if fieldLen <= 5 {
                    guard pos + fieldLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+fieldLen], encoding: .utf8) {
                        decodedTag.childEnums = enumStr
                        decodedTag.childType = .enums
                    }
                    pos += fieldLen
                } else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        decodedTag.childEnums = enumStr
                        decodedTag.childType = .enums
                    }
                    pos += strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2
                    guard pos + strLen <= data.count else { break }
                    if let enumStr = String(bytes: data[pos..<pos+strLen], encoding: .utf8) {
                        decodedTag.childEnums = enumStr
                        decodedTag.childType = .enums
                    }
                    pos += strLen
                }

            case TagKey.childMime:
                let (childMimeVal, childMimeAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if childMimeAdv > 0 {
                    decodedTag.childMime = Mime.toString(childMimeVal)
                    decodedTag.childType = .media
                    pos += childMimeAdv
                }

            case TagKey.size:
                let (sizeVal, sizeAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if sizeAdv > 0 {
                    decodedTag.size = Int(sizeVal)
                    pos += sizeAdv
                }

            case TagKey.childSize:
                let (csVal, csAdv) = decodeTagFieldU64(data, pos: pos, fieldLen: fieldLen)
                if csAdv > 0 {
                    decodedTag.childSize = Int(csVal)
                    pos += csAdv
                }

            default:
                if fieldLen <= 5 { pos += fieldLen }
                else if fieldLen == 6 {
                    guard pos < data.count else { break }
                    let strLen = Int(data[pos])
                    pos += 1 + strLen
                } else if fieldLen == 7 {
                    guard pos + 1 < data.count else { break }
                    let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
                    pos += 2 + strLen
                }
            }
        }

        pos = tagFieldEnd

        if decodedTag.isNull {
            if decodedTag.type != .unknown {
                return tagNullValue(decodedTag.type)
            }
        }

        guard pos < data.count else {
            if decodedTag.isNull {
                return .null
            }
            throw MMError.unexpectedEndOfData
        }

        let payloadData = Data(data[pos...])
        let innerDecoder = Decoder(data: payloadData)
        let value = try innerDecoder.decodeValue(tag: decodedTag)

        // Apply child type transformations
        if decodedTag.childType != .unknown {
            switch value {
            case .array(let elements):
                var result: [DecodedValue] = []
                for element in elements {
                    result.append(applyChildType(element, childType: decodedTag.childType, childEnums: decodedTag.childEnums))
                }
                return .array(result)
            case .object(let dict):
                var result: [String: DecodedValue] = [:]
                for (key, val) in dict {
                    result[key] = applyChildType(val, childType: decodedTag.childType, childEnums: decodedTag.childEnums)
                }
                return .object(result)
            default:
                return applyChildType(value, childType: decodedTag.childType, childEnums: decodedTag.childEnums)
            }
        }

        if decodedTag.type != .unknown {
            return applyTagType(value, tag: decodedTag)
        }
        return value
    }

    private func applyChildType(_ value: DecodedValue, childType: ValueType, childEnums: String) -> DecodedValue {
        switch childType {
        case .bytes, .media:
            if case .data(let d) = value {
                return .string(d.base64EncodedString())
            }
        case .bigint:
            if case .data(let d) = value {
                return .bigint(NodeDecoder.decodeBigIntFromBytes([UInt8](d)))
            }
        case .time:
            if case .int(let seconds) = value {
                let secs = max(0, min(seconds, 86399))
                let h = secs / 3600
                let m = (secs % 3600) / 60
                let s = secs % 60
                return .string(String(format: "%02d:%02d:%02d", h, m, s))
            }
        case .date:
            if case .int(let days) = value {
                let ref = Date(timeIntervalSince1970: 0)
                let date = Calendar(identifier: .gregorian).date(byAdding: .day, value: Int(days), to: ref) ?? ref
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }
        case .datetime:
            if case .int(let ts) = value {
                let date = Date(timeIntervalSince1970: TimeInterval(ts))
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }
        case .uuid:
            if case .data(let d) = value, d.count == 16 {
                let bytes = [UInt8](d)
                let hexStr = String(format: "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
                    bytes[0], bytes[1], bytes[2], bytes[3],
                    bytes[4], bytes[5], bytes[6], bytes[7],
                    bytes[8], bytes[9], bytes[10], bytes[11],
                    bytes[12], bytes[13], bytes[14], bytes[15])
                return .string(hexStr)
            }
        case .enums:
            if !childEnums.isEmpty {
                let enumValues = childEnums.split(separator: "|").map { $0.trimmingCharacters(in: .whitespaces) }
                if case .int(let idx) = value, idx >= 0, idx < enumValues.count {
                    return .string(String(enumValues[Int(idx)]))
                }
                if case .uint(let idx) = value, idx < enumValues.count {
                    return .string(String(enumValues[Int(idx)]))
                }
            }
        default:
            break
        }
        return value
    }

    private func applyTagType(_ value: DecodedValue, tag: Tag) -> DecodedValue {
        switch tag.type {
        case .datetime:
            if case .int(let ts) = value {
                let adjustedTs = TimeInterval(ts) + TimeInterval(tag.location * 3600)
                let date = Date(timeIntervalSince1970: adjustedTs)
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }
            if case .uint(let ts) = value {
                let adjustedTs = TimeInterval(ts) + TimeInterval(tag.location * 3600)
                let date = Date(timeIntervalSince1970: adjustedTs)
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd HH:mm:ss"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }

        case .date:
            if case .int(let days) = value {
                let ref = Date(timeIntervalSince1970: 0)
                let date = Calendar(identifier: .gregorian).date(byAdding: .day, value: Int(days), to: ref) ?? ref
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }
            if case .uint(let days) = value {
                let ref = Date(timeIntervalSince1970: 0)
                let date = Calendar(identifier: .gregorian).date(byAdding: .day, value: Int(days), to: ref) ?? ref
                let formatter = DateFormatter()
                formatter.dateFormat = "yyyy-MM-dd"
                formatter.timeZone = TimeZone(abbreviation: "UTC")
                return .string(formatter.string(from: date))
            }

        case .time:
            if case .int(let secs) = value {
                let hours = secs / 3600
                let mins = (secs % 3600) / 60
                let sec = secs % 60
                return .string(String(format: "%02d:%02d:%02d", hours, mins, sec))
            }
            if case .uint(let secs) = value {
                let hours = secs / 3600
                let mins = (secs % 3600) / 60
                let sec = secs % 60
                return .string(String(format: "%02d:%02d:%02d", hours, mins, sec))
            }

        case .uuid:
            if case .data(let d) = value, d.count == 16 {
                let bytes = [UInt8](d)
                let hexStr = String(format: "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
                    bytes[0], bytes[1], bytes[2], bytes[3],
                    bytes[4], bytes[5], bytes[6], bytes[7],
                    bytes[8], bytes[9], bytes[10], bytes[11],
                    bytes[12], bytes[13], bytes[14], bytes[15])
                return .string(hexStr)
            }

        case .enums:
            if !tag.enums.isEmpty {
                let enumValues = tag.enums.split(separator: "|").map { $0.trimmingCharacters(in: .whitespaces) }
                if case .int(let idx) = value, idx >= 0, idx < enumValues.count {
                    return .string(String(enumValues[Int(idx)]))
                }
                if case .uint(let idx) = value, idx < enumValues.count {
                    return .string(String(enumValues[Int(idx)]))
                }
            }

        case .media:
            if case .data(let d) = value {
                return .string(d.base64EncodedString())
            }

        case .bytes:
            if case .data(let d) = value {
                return .string(d.base64EncodedString())
            }

        default:
            break
        }

        return value
    }

    private func tagNullValue(_ type: ValueType) -> DecodedValue {
        switch type {
        case .i, .i8, .i16, .i32, .i64:
            return .int(0)
        case .u, .u8, .u16, .u32, .u64:
            return .uint(0)
        case .f32, .f64:
            return .float(0.0)
        case .bool:
            return .bool(false)
        case .str, .email, .url, .datetime, .date, .time, .enums, .media:
            return .string("")
        case .bytes, .uuid, .ip:
            return .data(Data())
        default:
            return .null
        }
    }
}

// MARK: - Tag field U64 decoder
// Reads a uint64 value from tag field bytes.
// fieldLen (lower 3 bits of the tag key byte) encodes the number of extra bytes.
// Value = (fieldLen + 1) bytes, big-endian.
func decodeTagFieldU64(_ data: [UInt8], pos: Int, fieldLen: Int) -> (UInt64, Int) {
    let byteCount = fieldLen + 1
    guard pos + byteCount <= data.count else {
        return (0, 0)
    }
    var value: UInt64 = 0
    for i in 0..<byteCount {
        value = (value << 8) | UInt64(data[pos + i])
    }
    return (value, byteCount)
}

// MARK: - Tag field string decoder
// Reads a string value from tag field bytes.
// fieldLen encodes the length: 0-5 = direct length, 6 = 1-byte length, 7 = 2-byte length
func decodeTagString(_ data: [UInt8], pos: Int, fieldLen: Int) -> (String, Int) {
    if fieldLen <= 5 {
        guard pos + fieldLen <= data.count else { return ("", 0) }
        if fieldLen == 0 { return ("", 0) }
        if let str = String(bytes: data[pos..<pos+fieldLen], encoding: .utf8) {
            return (str, fieldLen)
        }
        return ("", fieldLen)
    } else if fieldLen == 6 {
        guard pos < data.count else { return ("", 0) }
        let strLen = Int(data[pos])
        guard pos + 1 + strLen <= data.count else { return ("", 0) }
        if let str = String(bytes: data[(pos+1)..<(pos+1+strLen)], encoding: .utf8) {
            return (str, 1 + strLen)
        }
        return ("", 1 + strLen)
    } else if fieldLen == 7 {
        guard pos + 1 < data.count else { return ("", 0) }
        let strLen = (Int(data[pos]) << 8) | Int(data[pos + 1])
        guard pos + 2 + strLen <= data.count else { return ("", 0) }
        if let str = String(bytes: data[(pos+2)..<(pos+2+strLen)], encoding: .utf8) {
            return (str, 2 + strLen)
        }
        return ("", 2 + strLen)
    }
    return ("", 0)
}

// MARK: - Mantissa to decimal string converter
// Converts a mantissa and exponent into a decimal string representation.
func mantissaToDecimal(_ mantissa: UInt64, _ exponent: Int8) -> String {
    let numStr = String(mantissa)
    let decimalPos = numStr.count + Int(exponent)

    if decimalPos <= 0 {
        return "0." + String(repeating: "0", count: -decimalPos) + numStr
    } else if decimalPos > 0 && decimalPos < numStr.count {
        let prefix = String(numStr[numStr.startIndex..<numStr.index(numStr.startIndex, offsetBy: decimalPos)])
        let suffix = String(numStr[numStr.index(numStr.startIndex, offsetBy: decimalPos)...])
        return prefix + "." + suffix
    } else {
        let trailingZeros = decimalPos - numStr.count
        return numStr + String(repeating: "0", count: trailingZeros)
    }
}