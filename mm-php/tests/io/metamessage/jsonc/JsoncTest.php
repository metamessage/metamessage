<?php

namespace io\metamessage\jsonc;

use PHPUnit\Framework\TestCase;

class JsoncParserTest extends TestCase {

    public function testParseEmptyObject(): void {
        $source = "{}";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertCount(0, $result->fields);
    }

    public function testParseSimpleObject(): void {
        $source = '{"key": "value"}';
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertCount(1, $result->fields);
        $this->assertEquals("key", $result->fields[0]->key);
    }

    public function testParseNumber(): void {
        $source = "123";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncValue::class, $result);
        $this->assertEquals(123, $result->data);
    }

    public function testParseFloat(): void {
        $source = "3.14";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncValue::class, $result);
        $this->assertEqualsWithDelta(3.14, $result->data, 0.001);
    }

    public function testParseBoolean(): void {
        $sourceTrue = "true";
        $resultTrue = parseJsonc($sourceTrue);
        $this->assertInstanceOf(JsoncValue::class, $resultTrue);
        $this->assertEquals(true, $resultTrue->data);

        $sourceFalse = "false";
        $resultFalse = parseJsonc($sourceFalse);
        $this->assertInstanceOf(JsoncValue::class, $resultFalse);
        $this->assertEquals(false, $resultFalse->data);
    }

    public function testParseNull(): void {
        $source = "null";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncValue::class, $result);
        $this->assertNull($result->data);
        $this->assertTrue($result->tag->isNull);
    }

    public function testParseArray(): void {
        $source = "[1, 2, 3]";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncArray::class, $result);
        $this->assertCount(3, $result->items);
    }

    public function testParseNestedObject(): void {
        $source = '{"outer": {"inner": "value"}}';
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertEquals("outer", $result->fields[0]->key);
        $this->assertInstanceOf(JsoncObject::class, $result->fields[0]->value);
        $inner = $result->fields[0]->value;
        $this->assertEquals("inner", $inner->fields[0]->key);
    }

    public function testParseWithLineComment(): void {
        $source = "{\n// this is a comment\n\"key\": \"value\"\n}";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
    }

    public function testParseWithBlockComment(): void {
        $source = "{/* this is a block comment */\"key\": \"value\"}";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
    }

    public function testParseWithTrailingComma(): void {
        $source = '{"key": "value",}';
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertCount(1, $result->fields);
    }
}

class JsoncPrinterTest extends TestCase {

    public function testPrintEmptyObject(): void {
        $obj = new JsoncObject();
        $result = JsoncPrinter::toString($obj);
        $this->assertStringContainsString("{", $result);
        $this->assertStringContainsString("}", $result);
    }

    public function testPrintSimpleObject(): void {
        $obj = new JsoncObject();
        $obj->fields[] = new JsoncField("key", new JsoncValue("value", "\"value\""));
        $result = JsoncPrinter::toString($obj);
        $this->assertStringContainsString("key", $result);
    }

    public function testPrintCompact(): void {
        $obj = new JsoncObject();
        $obj->fields[] = new JsoncField("key", new JsoncValue("value", "\"value\""));
        $result = JsoncPrinter::toCompactString($obj);
        $this->assertStringNotContainsString("\n", $result);
    }
}

class JsoncScannerTest extends TestCase {

    public function testScanEmptyInput(): void {
        $scanner = new JsoncScanner("");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::EOF, $token->type);
    }

    public function testScanLBrace(): void {
        $scanner = new JsoncScanner("{");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::LBrace, $token->type);
    }

    public function testScanRBrace(): void {
        $scanner = new JsoncScanner("}");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::RBrace, $token->type);
    }

    public function testScanString(): void {
        $scanner = new JsoncScanner("\"hello\"");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::String, $token->type);
        $this->assertEquals("hello", $token->literal);
    }

    public function testScanNumber(): void {
        $scanner = new JsoncScanner("123");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::Number, $token->type);
        $this->assertEquals("123", $token->literal);
    }

    public function testScanTrue(): void {
        $scanner = new JsoncScanner("true");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::True, $token->type);
    }

    public function testScanFalse(): void {
        $scanner = new JsoncScanner("false");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::False, $token->type);
    }

    public function testScanLineCommentTrailing(): void {
        $scanner = new JsoncScanner("// this is a comment\n");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::TrailingComment, $token->type);
    }

    public function testScanBlockCommentTrailing(): void {
        $scanner = new JsoncScanner("/* block comment */");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::TrailingComment, $token->type);
    }
}
