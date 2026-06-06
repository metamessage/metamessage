<?php

namespace io\metamessage\jsonc;

use PHPUnit\Framework\TestCase;

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

    public function testScanLineComment(): void {
        $scanner = new JsoncScanner("// this is a comment\n");
        $token = $scanner->nextToken();
        $this->assertEquals(JsoncTokenType::Comment, $token->type);
    }
}