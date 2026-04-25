import MetaMessage

// 创建 Person 对象
let person = Person(name: "Ed", age: 30)
print("Original: Name=\(person.name), Age=\(person.age)")

// 编码到 Wire 格式
let wire = MetaMessage.encode(Person(name: "Ed", age: 30))
print("Encoded: \(bytesToHex(wire))")

// 从 Wire 解码
if let decoded = try? MetaMessage.decode(wire) {
    print("Decoded: \(decoded)")
}

func bytesToHex(_ data: Data) -> String {
    return data.map { String(format: "%02x", $0) }.joined()
}
