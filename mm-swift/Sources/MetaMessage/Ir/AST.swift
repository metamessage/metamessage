import Foundation

public enum NodeType: String {
    case unknown
    case object
    case array
    case value
    case doc
}

public protocol Node {
    func getTag() -> Tag?
    func getType() -> NodeType
    func getPath() -> String
    func setPath(_ path: String)
    func setTag(_ tag: Tag)
}

public struct Field {
    public var key: String
    public var value: Node

    public init(key: String, value: Node) {
        self.key = key
        self.value = value
    }
}

public class MMObject: Node {
    public var fields: [Field]
    public var tag: Tag?
    public var path: String

    public init(fields: [Field] = [], tag: Tag? = nil, path: String = "") {
        self.fields = fields
        self.tag = tag
        self.path = path
    }

    public func getTag() -> Tag? {
        return tag
    }

    public func getType() -> NodeType {
        return .object
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }

    public func setTag(_ tag: Tag) {
        self.tag = tag
    }
}

public class MMArray: Node {
    public var items: [Node]
    public var tag: Tag?
    public var path: String

    public init(items: [Node] = [], tag: Tag? = nil, path: String = "") {
        self.items = items
        self.tag = tag
        self.path = path
    }

    public func getTag() -> Tag? {
        return tag
    }

    public func getType() -> NodeType {
        return .array
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }

    public func setTag(_ tag: Tag) {
        self.tag = tag
    }
}

public class Value: Node {
    public var data: Any?
    public var text: String
    public var tag: Tag?
    public var path: String

    public init(data: Any? = nil, text: String = "", tag: Tag? = nil, path: String = "") {
        self.data = data
        self.text = text
        self.tag = tag
        self.path = path
    }

    public func getTag() -> Tag? {
        return tag
    }

    public func getType() -> NodeType {
        return .value
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }

    public func setTag(_ tag: Tag) {
        self.tag = tag
    }
}

public class MMDoc: Node {
    public var fields: [Field]
    public var tag: Tag?
    public var path: String

    public init(fields: [Field] = [], tag: Tag? = nil, path: String = "") {
        self.fields = fields
        self.tag = tag
        self.path = path
    }

    public func getTag() -> Tag? {
        return tag
    }

    public func getType() -> NodeType {
        return .doc
    }

    public func getPath() -> String {
        return path
    }

    public func setPath(_ path: String) {
        self.path = path
    }

    public func setTag(_ tag: Tag) {
        self.tag = tag
    }
}