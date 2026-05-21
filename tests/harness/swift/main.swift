// MetaMessage Swift test harness - parse JSONC file and re-print to JSONC.
import Foundation
import MetaMessage

guard CommandLine.arguments.count > 1 else {
    fputs(stderr, "usage: harness <file.jsonc>\n")
    exit(1)
}

let path = CommandLine.arguments[1]
guard let input = try? String(contentsOfFile: path, encoding: .utf8) else {
    fputs(stderr, "read error: cannot open \(path)\n")
    exit(1)
}

guard let node = try? parseJSONC(input) else {
    fputs(stderr, "parse error\n")
    exit(1)
}

let output = JSONCPrinter().print(node)
print(output, terminator: "")