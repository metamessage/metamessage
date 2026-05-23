<?php

namespace io\metamessage\jsonc;

use PHPUnit\Framework\TestCase;

class JsoncPrinterTest extends TestCase
{

    public function testPrintEmptyObject(): void
    {
        $obj = new JsoncObject();
        $result = JsoncPrinter::ToJSONC($obj);
        $this->assertStringContainsString("{", $result);
        $this->assertStringContainsString("}", $result);
    }

    public function testPrintSimpleObject(): void
    {
        $obj = new JsoncObject();
        $obj->Fields[] = new JsoncField("key", new JsoncValue("value", "\"value\""));
        $result = JsoncPrinter::ToJSONC($obj);
        $this->assertStringContainsString("key", $result);
    }
}