// MetaMessage Swift test harness - parse JSONC file and re-print to JSONC.
import Foundation
import MetaMessage

guard CommandLine.arguments.count > 1 else {
    fputs("usage: harness [--encode|--decode] <file.jsonc>\n", stderr)
    exit(1)
}

let arg1 = CommandLine.arguments[1]

if arg1 == "--encode" {
    guard CommandLine.arguments.count > 2 else {
        fputs("usage: harness --encode <file.jsonc>\n", stderr)
        exit(1)
    }
    let path = CommandLine.arguments[2]
    guard let input = try? String(contentsOfFile: path, encoding: .utf8) else {
        fputs("read error\n", stderr)
        exit(1)
    }
    do {
        let data = try MetaMessage.fromJSONC(input)
        let hex = data.map { String(format: "%02x", $0) }.joined()
        print(hex, terminator: "")
    } catch {
        fputs("encode error: \(error)\n", stderr)
        exit(1)
    }
    exit(0)
}

if arg1 == "--decode" {
    let inputData = FileHandle.standardInput.availableData
    guard let hexStr = String(data: inputData, encoding: .utf8) else {
        fputs("error reading hex from stdin\n", stderr)
        exit(1)
    }
    let trimmedHex = hexStr.trimmingCharacters(in: CharacterSet.whitespacesAndNewlines)
    var data = Data()
    var index = trimmedHex.startIndex
    while index < trimmedHex.endIndex {
        let nextIndex = trimmedHex.index(index, offsetBy: 2, limitedBy: trimmedHex.endIndex) ?? trimmedHex.endIndex
        if let byte = UInt8(trimmedHex[index..<nextIndex], radix: 16) {
            data.append(byte)
        }
        index = nextIndex
    }
    do {
        let output = try MetaMessage.toJSONC(data)
        print(output, terminator: "")
    } catch {
        fputs("decode error: \(error)\n", stderr)
        exit(1)
    }
    exit(0)
}

// Existing behavior
let path = CommandLine.arguments[1]
guard let input = try? String(contentsOfFile: path, encoding: .utf8) else {
    fputs("read error: cannot open \(path)\n", stderr)
    exit(1)
}

guard let node = try? parseJSONC(input) else {
    fputs("parse error\n", stderr)
    exit(1)
}

let output = JSONCPrinter().print(node)
print(output, terminator: "")