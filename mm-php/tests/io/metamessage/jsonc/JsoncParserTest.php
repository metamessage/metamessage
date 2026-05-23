<?php

namespace io\metamessage\jsonc;

use ParseError;
use PHPUnit\Framework\TestCase;

class JsoncParserTest extends TestCase
{

    public function testParseEmptyObject(): void
    {
        $source = "{}";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertCount(0, $result->Fields);
    }

    public function testParseSimpleObject(): void
    {
        $source = '{"key": "value"}';
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertCount(1, $result->Fields);
        $this->assertEquals("key", $result->Fields[0]->Key);
    }

    public function testParseNumber(): void
    {
        $source = "123";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncValue::class, $result);
        $this->assertEquals(123, $result->Data);
    }

    public function testParseFloat(): void
    {
        $source = "3.14";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncValue::class, $result);
        $this->assertEqualsWithDelta(3.14, $result->Data, 0.001);
    }

    public function testParseBoolean(): void
    {
        $sourceTrue = "true";
        $resultTrue = parseJsonc($sourceTrue);
        $this->assertInstanceOf(JsoncValue::class, $resultTrue);
        $this->assertEquals(true, $resultTrue->Data);

        $sourceFalse = "false";
        $resultFalse = parseJsonc($sourceFalse);
        $this->assertInstanceOf(JsoncValue::class, $resultFalse);
        $this->assertEquals(false, $resultFalse->Data);
    }

    public function testParseNullThrowsException(): void
    {
        $this->expectException(\Exception::class);
        $this->expectExceptionMessage('null literal is not supported');
        parseJsonc("null");
    }

    public function testParseArray(): void
    {
        $source = "[1, 2, 3]";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncArray::class, $result);
        $this->assertCount(3, $result->Items);
    }

    public function testParseNestedObject(): void
    {
        $source = '{"outer": {"inner": "value"}}';
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertEquals("outer", $result->Fields[0]->Key);
        $this->assertInstanceOf(JsoncObject::class, $result->Fields[0]->Value);
        $inner = $result->Fields[0]->Value;
        $this->assertEquals("inner", $inner->Fields[0]->Key);
    }

    public function testParseWithLineComment(): void
    {
        $source = "{\n// this is a comment\n\"key\": \"value\"\n}";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
    }

    public function testParseWithBlockComment(): void
    {
        $source = "{/* this is a block comment */\"key\": \"value\"}";
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
    }

    public function testParseWithTrailingComma(): void
    {
        $source = '{"key": "value",}';
        $result = parseJsonc($source);
        $this->assertInstanceOf(JsoncObject::class, $result);
        $this->assertCount(1, $result->Fields);
    }
}
