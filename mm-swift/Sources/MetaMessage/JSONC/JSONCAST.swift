import Foundation

public enum JSONCNodeType: String {
    case unknown
    case object
    case array
    case value
    case doc
}

public protocol JSONCNode {
    func getTag() -> JSONCTag?
    func getType() -> JSONCNodeType
    func getPath() -> String
    func setPath(_ path: String)
}

public struct JSONCField {
    public var key: String
    public var value: JSONCNode

    public init(key: String, value: JSONCNode) {
        self.key = key
        self.value = value
    }
}

public class JSONCObject: JSONCNode {
    public var fields: [JSONCField]
    public var tag: JSONCTag?
    public var path: String

    public init(fields: [JSONCField] = [], tag: JSONCTag? = nil, path: String = "") {
        self.fields = fields
        self.tag = tag
        self.path = path
    }

    public func getTag() -> JSONCTag? {
        return tag
    }

    public func getType() -> JSONCNodeType {
        return .object
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }
}

public class JSONCArray: JSONCNode {
    public var items: [JSONCNode]
    public var tag: JSONCTag?
    public var path: String

    public init(items: [JSONCNode] = [], tag: JSONCTag? = nil, path: String = "") {
        self.items = items
        self.tag = tag
        self.path = path
    }

    public func getTag() -> JSONCTag? {
        return tag
    }

    public func getType() -> JSONCNodeType {
        return .array
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }
}

public class JSONCValue: JSONCNode {
    public var data: Any?
    public var text: String
    public var tag: JSONCTag?
    public var path: String

    public init(data: Any? = nil, text: String = "", tag: JSONCTag? = nil, path: String = "") {
        self.data = data
        self.text = text
        self.tag = tag
        self.path = path
    }

    public func getTag() -> JSONCTag? {
        return tag
    }

    public func getType() -> JSONCNodeType {
        return .value
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }
}

public class JSONCDoc: JSONCNode {
    public var fields: [JSONCField]
    public var tag: JSONCTag?
    public var path: String

    public init(fields: [JSONCField] = [], tag: JSONCTag? = nil, path: String = "") {
        self.fields = fields
        self.tag = tag
        self.path = path
    }

    public func getTag() -> JSONCTag? {
        return tag
    }

    public func getType() -> JSONCNodeType {
        return .doc
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }
}