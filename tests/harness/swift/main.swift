// MetaMessage Swift test harness - parse JSONC file and re-print to JSONC.
import Foundation
import MetaMessage

guard CommandLine.arguments.count > 1 else {
    fputs("usage: harness <file.jsonc>\n", stderr)
    exit(1)
}

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