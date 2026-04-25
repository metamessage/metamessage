import Foundation

public class JSONCTag {
    public var name: String = ""

    public var isNull: Bool = false
    public var example: Bool = false

    public var desc: String = ""
    public var type: ValueType = .unknown
    public var raw: Bool = false
    public var nullable: Bool = false
    public var allowEmpty: Bool = false
    public var unique: Bool = false
    public var defaultValue: String = ""
    public var min: String = ""
    public var max: String = ""
    public var size: Int = 0
    public var enumValues: String = ""
    public var pattern: String = ""
    public var locationOffset: Int = 0
    public var version: Int = 0
    public var mime: String = ""

    public var childDesc: String = ""
    public var childType: ValueType = .unknown
    public var childRaw: Bool = false
    public var childNullable: Bool = false
    public var childAllowEmpty: Bool = false
    public var childUnique: Bool = false
    public var childDefault: String = ""
    public var childMin: String = ""
    public var childMax: String = ""
    public var childSize: Int = 0
    public var childEnum: String = ""
    public var childPattern: String = ""
    public var childLocationOffset: Int = 0
    public var childVersion: Int = 0
    public var childMime: String = ""

    public var isInherit: Bool = false

    public init() {}

    public func inherit(from tag: JSONCTag) {
        self.desc = tag.childDesc
        self.type = tag.childType
        self.raw = tag.childRaw
        self.nullable = tag.childNullable
        self.allowEmpty = tag.childAllowEmpty
        self.unique = tag.childUnique
        self.defaultValue = tag.childDefault
        self.min = tag.childMin
        self.max = tag.childMax
        self.size = tag.childSize
        self.enumValues = tag.childEnum
        self.pattern = tag.childPattern
        self.locationOffset = tag.childLocationOffset
        self.version = tag.childVersion
        self.mime = tag.childMime
    }

    public func stringValue() -> String {
        var parts: [String] = []

        if type != .unknown && !isInherit {
            if type == .string || type == .int || type == .float64 || type == .bool || type == .structType || type == .slice {
            } else {
                if type == .array && size > 0 || type == .enumValue && enumValues != "" {
                } else {
                    parts.append("type=\(type.stringValue)")
                }
            }
        }

        if example {
            parts.append("example")
        }

        if isNull {
            parts.append("is_null")
        }

        if nullable && !isInherit && !isNull {
            parts.append("nullable")
        }

        if desc != "" && !isInherit {
            parts.append("desc=\"\(desc)\"")
        }

        if raw && !isInherit {
            parts.append("raw")
        }

        if allowEmpty && !isInherit {
            parts.append("allow_empty")
        }

        if unique && !isInherit {
            parts.append("unique")
        }

        if defaultValue != "" && !isInherit {
            parts.append("default=\(defaultValue)")
        }

        if min != "" && !isInherit {
            parts.append("min=\(min)")
        }

        if max != "" && !isInherit {
            parts.append("max=\(max)")
        }

        if size != 0 && !isInherit {
            parts.append("size=\(size)")
        }

        if enumValues != "" && !isInherit {
            parts.append("enum=\(enumValues)")
        }

        if pattern != "" && !isInherit {
            parts.append("pattern=\(pattern)")
        }

        if locationOffset != 0 && !isInherit {
            parts.append("location=\(locationOffset)")
        }

        if version != 0 && !isInherit {
            parts.append("version=\(version)")
        }

        if mime != "" && !isInherit {
            parts.append("mime=\(mime)")
        }

        if childDesc != "" {
            parts.append("child_desc=\"\(childDesc)\"")
        }

        if childType != .unknown {
            if childType == .string || childType == .int || childType == .float64 || childType == .bool || childType == .structType || childType == .slice {
            } else {
                if childType == .array && childSize > 0 || childType == .enumValue && childEnum != "" {
                } else {
                    parts.append("child_type=\(childType.stringValue)")
                }
            }
        }

        if childRaw {
            parts.append("child_raw")
        }

        if childNullable {
            parts.append("child_nullable")
        }

        if childAllowEmpty {
            parts.append("child_allow_empty")
        }

        if childUnique {
            parts.append("child_unique")
        }

        if childDefault != "" {
            parts.append("child_default=\(childDefault)")
        }

        if childMin != "" {
            parts.append("child_min=\(childMin)")
        }

        if childMax != "" {
            parts.append("child_max=\(childMax)")
        }

        if childSize != 0 {
            parts.append("child_size=\(childSize)")
        }

        if childEnum != "" {
            parts.append("child_enum=\(childEnum)")
        }

        if childPattern != "" {
            parts.append("child_pattern=\(childPattern)")
        }

        if childLocationOffset != 0 {
            parts.append("child_location=\(childLocationOffset)")
        }

        if childVersion != 0 {
            parts.append("child_version=\(childVersion)")
        }

        if childMime != "" {
            parts.append("child_mime=\(childMime)")
        }

        return parts.joined(separator: "; ")
    }
}

public func parseMMTag(_ tagStr: String) -> JSONCTag? {
    var tag = tagStr.trimmingCharacters(in: .whitespaces)

    if tag.hasPrefix("//") {
        tag = String(tag.dropFirst(2))
    }
    tag = tag.trimmingCharacters(in: .whitespaces)

    if tag.hasPrefix("mm:") {
        tag = String(tag.dropFirst(3))
    }
    tag = tag.trimmingCharacters(in: .whitespaces)

    if tag.isEmpty {
        return JSONCTag()
    }

    let result = JSONCTag()
    let parts = tag.split(separator: ";").map { String($0).trimmingCharacters(in: .whitespaces) }

    for part in parts {
        guard !part.isEmpty else { continue }

        var key: String
        var value: String

        if let equalIndex = part.firstIndex(of: "=") {
            key = String(part[..<equalIndex]).trimmingCharacters(in: .whitespaces)
            value = String(part[part.index(after: equalIndex)...]).trimmingCharacters(in: .whitespaces)
        } else {
            key = part.trimmingCharacters(in: .whitespaces)
            value = ""
        }

        let lowerKey = key.lowercased()

        switch lowerKey {
        case "is_null":
            result.isNull = true
            result.nullable = true

        case "example":
            result.example = true

        case "desc":
            result.desc = value

        case "type":
            if let t = ValueType.parse(value) {
                result.type = t
            }

        case "raw":
            result.raw = true

        case "nullable":
            result.nullable = true

        case "allow_empty":
            result.allowEmpty = true

        case "unique":
            result.unique = true

        case "default":
            result.defaultValue = value

        case "pattern":
            result.pattern = value

        case "min":
            result.min = value

        case "max":
            result.max = value

        case "size":
            if let size = Int(value) {
                result.size = size
            }

        case "enum":
            result.type = .enumValue
            result.enumValues = value

        case "location":
            if let offset = Int(value), offset >= -12, offset <= 14 {
                result.locationOffset = offset
            }

        case "version":
            if let ver = Int(value), ver >= 1, ver <= 10 {
                result.version = ver
            }

        case "mime":
            result.mime = value

        case "child_desc":
            result.childDesc = value

        case "child_type":
            if let t = ValueType.parse(value) {
                result.childType = t
            }

        case "child_raw":
            result.childRaw = true

        case "child_nullable":
            result.childNullable = true

        case "child_allow_empty":
            result.childAllowEmpty = true

        case "child_unique":
            result.childUnique = true

        case "child_default":
            result.childDefault = value

        case "child_pattern":
            result.childPattern = value

        case "child_min":
            result.childMin = value

        case "child_max":
            result.childMax = value

        case "child_size":
            if let size = Int(value) {
                result.childSize = size
            }

        case "child_enum":
            result.childEnum = value
            result.childType = .enumValue

        case "child_location":
            if let offset = Int(value), offset >= -12, offset <= 14 {
                result.childLocationOffset = offset
            }

        case "child_version":
            if let ver = Int(value), ver >= 1, ver <= 10 {
                result.childVersion = ver
            }

        case "child_mime":
            result.childMime = value

        default:
            break
        }
    }

    return result
}