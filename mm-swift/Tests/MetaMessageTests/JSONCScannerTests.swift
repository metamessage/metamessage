import XCTest
@testable import MetaMessage

final class JSONCScannerTests: XCTestCase {
    func testScanEmpty() {
        let scanner = JSONCScanner(input: "")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .eof)
    }

    func testScanBraces() {
        let scanner = JSONCScanner(input: "{}")
        XCTAssertEqual(scanner.nextToken().type, .lBrace)
        XCTAssertEqual(scanner.nextToken().type, .rBrace)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanBrackets() {
        let scanner = JSONCScanner(input: "[]")
        XCTAssertEqual(scanner.nextToken().type, .lBracket)
        XCTAssertEqual(scanner.nextToken().type, .rBracket)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanColon() {
        let scanner = JSONCScanner(input: ":")
        XCTAssertEqual(scanner.nextToken().type, .colon)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanComma() {
        let scanner = JSONCScanner(input: ",")
        XCTAssertEqual(scanner.nextToken().type, .comma)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanString() {
        let scanner = JSONCScanner(input: "\"hello\"")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .string)
        XCTAssertEqual(token.literal, "hello")
    }

    func testScanEmptyString() {
        let scanner = JSONCScanner(input: "\"\"")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .string)
        XCTAssertEqual(token.literal, "")
    }

    func testScanStringWithEscape() {
        let scanner = JSONCScanner(input: "\"hello\\\"world\"")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .string)
        XCTAssertEqual(token.literal, "hello\\\"world")
    }

    func testScanNumber() {
        let scanner = JSONCScanner(input: "123")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .number)
        XCTAssertEqual(token.literal, "123")
    }

    func testScanNegativeNumber() {
        let scanner = JSONCScanner(input: "-456")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .number)
        XCTAssertEqual(token.literal, "-456")
    }

    func testScanFloatNumber() {
        let scanner = JSONCScanner(input: "3.14")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .number)
        XCTAssertEqual(token.literal, "3.14")
    }

    func testScanTrue() {
        let scanner = JSONCScanner(input: "true")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .trueValue)
        XCTAssertEqual(token.literal, "true")
    }

    func testScanFalse() {
        let scanner = JSONCScanner(input: "false")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .falseValue)
        XCTAssertEqual(token.literal, "false")
    }

    func testScanNull() {
        let scanner = JSONCScanner(input: "null")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .nullValue)
        XCTAssertEqual(token.literal, "null")
    }

    func testScanLineComment() {
        let scanner = JSONCScanner(input: "// this is a comment")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .leadingComment)
        XCTAssertEqual(token.literal, "this is a comment")
    }

    func testScanBlockComment() {
        let scanner = JSONCScanner(input: "/* this is a block comment */")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .leadingComment)
        XCTAssertEqual(token.literal, "this is a block comment")
    }

    func testScanMultiLineBlockComment() {
        let scanner = JSONCScanner(input: "/* this is\na block\ncomment */")
        let token = scanner.nextToken()
        XCTAssertEqual(token.type, .leadingComment)
        XCTAssertEqual(token.literal, "this is\na block\ncomment")
    }

    func testScanWhitespace() {
        let scanner = JSONCScanner(input: "   \t\n\r  {}")
        XCTAssertEqual(scanner.nextToken().type, .lBrace)
        XCTAssertEqual(scanner.nextToken().type, .rBrace)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanMixed() {
        let scanner = JSONCScanner(input: """
        {
            "name": "test",
            "value": 123,
            "active": true
        }
        """)

        XCTAssertEqual(scanner.nextToken().type, .lBrace)
        XCTAssertEqual(scanner.nextToken().type, .string)
        XCTAssertEqual(scanner.nextToken().type, .colon)
        XCTAssertEqual(scanner.nextToken().type, .string)
        XCTAssertEqual(scanner.nextToken().type, .comma)
        XCTAssertEqual(scanner.nextToken().type, .string)
        XCTAssertEqual(scanner.nextToken().type, .colon)
        XCTAssertEqual(scanner.nextToken().type, .number)
        XCTAssertEqual(scanner.nextToken().type, .comma)
        XCTAssertEqual(scanner.nextToken().type, .string)
        XCTAssertEqual(scanner.nextToken().type, .colon)
        XCTAssertEqual(scanner.nextToken().type, .trueValue)
        XCTAssertEqual(scanner.nextToken().type, .rBrace)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanLineCommentWithJSON() {
        let scanner = JSONCScanner(input: """
        // mm:type=str
        "value"
        """)
        XCTAssertEqual(scanner.nextToken().type, .leadingComment)
        XCTAssertEqual(scanner.nextToken().type, .string)
        XCTAssertEqual(scanner.nextToken().type, .eof)
    }

    func testScanLineNumbers() {
        let scanner = JSONCScanner(input: "\"a\"\n\"b\"\n\"c\"")
        let t1 = scanner.nextToken()
        XCTAssertEqual(t1.line, 1)
        XCTAssertEqual(t1.column, 1)

        let t2 = scanner.nextToken()
        XCTAssertEqual(t2.line, 2)

        let t3 = scanner.nextToken()
        XCTAssertEqual(t3.line, 3)
    }
}