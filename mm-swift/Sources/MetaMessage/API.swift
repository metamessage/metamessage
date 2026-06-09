import Foundation

// MARK: - Simple encode/decode

/// Encode a boolean value to MetaMessage wire format.
public func encode(_ value: Bool) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode an integer value to MetaMessage wire format.
public func encode(_ value: Int) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode an Int8 value to MetaMessage wire format.
public func encode(_ value: Int8) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode an Int16 value to MetaMessage wire format.
public func encode(_ value: Int16) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode an Int32 value to MetaMessage wire format.
public func encode(_ value: Int32) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode an Int64 value to MetaMessage wire format.
public func encode(_ value: Int64) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a UInt value to MetaMessage wire format.
public func encode(_ value: UInt) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a UInt8 value to MetaMessage wire format.
public func encode(_ value: UInt8) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a UInt16 value to MetaMessage wire format.
public func encode(_ value: UInt16) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a UInt32 value to MetaMessage wire format.
public func encode(_ value: UInt32) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a UInt64 value to MetaMessage wire format.
public func encode(_ value: UInt64) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a Float value to MetaMessage wire format.
public func encode(_ value: Float) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a Double value to MetaMessage wire format.
public func encode(_ value: Double) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a String value to MetaMessage wire format.
public func encode(_ value: String) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Encode a Data value to MetaMessage wire format.
public func encode(_ value: Data) -> Data {
    return MetaMessageEncoder().encode(value)
}

/// Decode MetaMessage wire format to a decoded value.
public func decode(_ data: Data) throws -> Decoder.DecodedValue {
    return try MetaMessageDecoder(data: data).decode()
}

/// Convert wire format data directly to JSONC string.
public func toJSONC(_ data: Data) throws -> String {
    let decoder = NodeDecoder(data: data)
    let node = try decoder.decode()
    let printer = JSONCPrinter()
    return printer.print(node)
}

/// Validate a value against a tag specification.
public func validate(_ value: Any?, tag: Tag) -> ValidationResult {
    return validator.validate(value, tag: tag)
}

/// Convert a Swift value to MetaMessage wire format with a tag string.
public func fromValue(_ value: Any?, tag: String) throws -> Data {
    let node = try valueToNode(value, tag: tag)
    let encoder = Encoder()
    return try encodeNode(node, with: encoder)
}

/// Parse a JSONC string and encode to MetaMessage wire format.
public func fromJSONC(_ s: String) throws -> Data {
    guard let node = try parseJSONC(s) else {
        throw MMError.invalidData
    }
    let encoder = Encoder()
    return try encodeNode(node, with: encoder)
}

// MARK: - High-Level API

/// Convert a Swift value directly to JSONC string.
public func valueToJSONC(_ value: Any?, name: String) throws -> String {
    let node = try valueToNode(value, tag: name)
    return JSONCPrinter().print(node)
}

/// Parse JSONC, re-encode to wire format, then decode to DecodedValue.
public func jsoncToValue(_ jsonc: String) throws -> Decoder.DecodedValue {
    guard let node = try parseJSONC(jsonc) else {
        throw MMError.invalidData
    }
    let encoder = Encoder()
    let data = try encodeNode(node, with: encoder)
    return try MetaMessageDecoder(data: data).decode()
}

/// Parse JSONC and bind the result to an existing object.
public func bindFromJSONC(_ inString: String, to out: AnyObject) throws {
    guard let node = try parseJSONC(inString) else {
        throw MMError.invalidData
    }
    try bindNode(node, to: out)
}

// MARK: - TS-compatible API

/// Encode a value to MetaMessage wire format (TS-compatible).
///
/// - Parameters:
///   - value: The value to encode.
///   - tag: Optional tag for metadata.
/// - Returns: Wire format data.
public func encodeFromValue(_ value: Any?, tag: Tag? = nil) throws -> Data {
    return try fromValue(value, tag: tag?.name ?? "")
}

/// Encode a JSONC string to MetaMessage wire format (TS-compatible).
///
/// - Parameter jsonc: JSONC string.
/// - Returns: Wire format data.
public func encodeFromJsonc(_ jsonc: String) throws -> Data {
    return try fromJSONC(jsonc)
}

/// Decode MetaMessage wire format to a decoded value (TS-compatible).
///
/// - Parameter wire: Wire format data.
/// - Returns: A DecodedValue representing the decoded data.
public func decodeToValue(_ wire: Data) throws -> Decoder.DecodedValue {
    return try decode(wire)
}

/// Decode MetaMessage wire format to JSONC string (TS-compatible).
///
/// - Parameter wire: Wire format data.
/// - Returns: JSONC string.
public func decodeToJsonc(_ wire: Data) throws -> String {
    return try toJSONC(wire)
}

/// Convert a value to JSONC string (TS-compatible).
///
/// - Parameters:
///   - value: The value to convert.
///   - tag: Optional tag for metadata.
/// - Returns: JSONC string.
public func valueToJsonc(_ value: Any?, tag: Tag? = nil) throws -> String {
    return try valueToJSONC(value, name: tag?.name ?? "")
}