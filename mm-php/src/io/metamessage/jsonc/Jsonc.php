<?php

namespace io\metamessage\jsonc;

use io\metamessage\ir\Node;

class Jsonc
{
    public static function ToJSONC(?Node $n): string
    {
        if ($n === null) {
            return '';
        }

        $printer = new JsoncPrinter();
        return $printer::ToJSONC($n);
    }

    public static function ParseFromString(string $s): Node
    {
        $scanner = new JsoncScanner($s);
        $toks = [];
        while (true) {
            $t = $scanner->nextToken();
            $toks[] = $t;
            if ($t->type === JsoncTokenType::EOF) {
                break;
            }
        }

        $parser = new JsoncParser($toks);
        return $parser->parse();
    }

    public static function Print(Node $n): void
    {
        echo self::ToJSONC($n) . "\n";
    }
}
