import Foundation

public class Tag {
    public var name: String = ""

    public var isNull: Bool = false
    public var example: Bool = false

    public var desc: String = ""
    public var type: ValueType = .unknown
    public var deprecated: Bool = false
    public var nullable: Bool = false
    public var allowEmpty: Bool = false
    public var unique: Bool = false
    public var defaultVal: String = ""
    public var min: String = ""
    public var max: String = ""
    public var size: Int = 0
    public var enums: String = ""
    public var pattern: String = ""
    public var location: Int = 0
    public var version: Int = 0
    public var mime: String = ""
    public var more: Int = 0

    public var childDesc: String = ""
    public var childType: ValueType = .unknown
    public var childRaw: Bool = false
    public var childNullable: Bool = false
    public var childAllowEmpty: Bool = false
    public var childUnique: Bool = false
    public var childDefaultVal: String = ""
    public var childMin: String = ""
    public var childMax: String = ""
    public var childSize: Int = 0
    public var childEnums: String = ""
    public var childPattern: String = ""
    public var childLocation: Int = 0
    public var childVersion: Int = 0
    public var childMime: String = ""

    public var isInherit: Bool = false

    public init() {}

    public func inherit(from tag: Tag) {
        if !tag.childDesc.isEmpty { self.desc = tag.childDesc }
        if tag.childType != .unknown { self.type = tag.childType }
        if tag.childNullable { self.nullable = true }
        if tag.childAllowEmpty { self.allowEmpty = true }
        if tag.childUnique { self.unique = true }
        if !tag.childDefaultVal.isEmpty { self.defaultVal = tag.childDefaultVal }
        if !tag.childMin.isEmpty { self.min = tag.childMin }
        if !tag.childMax.isEmpty { self.max = tag.childMax }
        if tag.childSize != 0 { self.size = tag.childSize }
        if !tag.childEnums.isEmpty { self.enums = tag.childEnums }
        if !tag.childPattern.isEmpty { self.pattern = tag.childPattern }
        if tag.childLocation != 0 { self.location = tag.childLocation }
        if tag.childVersion != 0 { self.version = tag.childVersion }
        if !tag.childMime.isEmpty { self.mime = tag.childMime }
    }

    public func stringValue() -> String {
        var parts: [String] = []

        if type != .unknown && !isInherit {
            if type == .str || type == .i || type == .f64 || type == .bool || type == .obj || type == .vec {
            } else {
                if type == .arr && size > 0 || type == .enums && enums != "" {
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

        if deprecated && !isInherit {
            parts.append("deprecated")
        }

        if allowEmpty && !isInherit {
            parts.append("allow_empty")
        }

        if unique && !isInherit {
            parts.append("unique")
        }

        if defaultVal != "" && !isInherit {
            parts.append("default_val=\(defaultVal)")
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

        if enums != "" && !isInherit {
            parts.append("enums=\(enums)")
        }

        if pattern != "" && !isInherit {
            parts.append("pattern=\(pattern)")
        }

        if location != 0 && !isInherit {
            parts.append("location=\(location)")
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
            if childType == .str || childType == .i || childType == .f64 || childType == .bool || childType == .obj || childType == .vec {
            } else {
                if childType == .arr && childSize > 0 || childType == .enums && childEnums != "" {
                } else {
                    parts.append("child_type=\(childType.stringValue)")
                }
            }
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

        if childDefaultVal != "" {
            parts.append("child_default_val=\(childDefaultVal)")
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

        if childEnums != "" {
            parts.append("child_enums=\(childEnums)")
        }

        if childPattern != "" {
            parts.append("child_pattern=\(childPattern)")
        }

        if childLocation != 0 {
            parts.append("child_location=\(childLocation)")
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

public func parseMMTag(_ tagStr: String) -> Tag? {
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
        return Tag()
    }

    let result = Tag()
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
            result.desc = value.trimmingCharacters(in: CharacterSet(charactersIn: "\""))

        case "type":
            if let t = ValueType.parse(value) {
                result.type = t
            }

        case "raw":
            result.deprecated = true

        case "deprecated":
            result.deprecated = true

        case "nullable":
            result.nullable = true

        case "allow_empty":
            result.allowEmpty = true

        case "unique":
            result.unique = true

        case "default", "default_val":
            result.defaultVal = value

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

        case "enum", "enums":
            result.type = .enums
            result.enums = value

        case "location":
            if let offset = Int(value), offset >= -12, offset <= 14 {
                result.location = offset
            }

        case "version":
            if let ver = Int(value), ver >= 1, ver <= 10 {
                result.version = ver
            }

        case "mime":
            result.mime = value

        case "child_desc":
            result.childDesc = value.trimmingCharacters(in: CharacterSet(charactersIn: "\""))

        case "child_type":
            if let t = ValueType.parse(value) {
                result.childType = t
            }

        case "child_raw":
            break

        case "child_nullable":
            result.childNullable = true

        case "child_allow_empty":
            result.childAllowEmpty = true

        case "child_unique":
            result.childUnique = true

        case "child_default_val":
            result.childDefaultVal = value

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

        case "child_enums":
            result.childEnums = value
            result.childType = .enums

        case "child_location":
            if let offset = Int(value), offset >= -12, offset <= 14 {
                result.childLocation = offset
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