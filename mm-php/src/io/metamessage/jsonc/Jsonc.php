<?php

namespace io\metamessage\jsonc;

use io\metamessage\ir\Node;
use io\metamessage\ir\Object_;
use io\metamessage\ir\Array_;
use io\metamessage\ir\Value;
use io\metamessage\ir\Field;

class_alias(Object_::class, 'io\\metamessage\\jsonc\\JsoncObject');
class_alias(Array_::class, 'io\\metamessage\\jsonc\\JsoncArray');
class_alias(Value::class, 'io\\metamessage\\jsonc\\JsoncValue');
class_alias(Field::class, 'io\\metamessage\\jsonc\\JsoncField');

function parseJsonc(string $s): Node
{
    return Jsonc::ParseFromString($s);
}

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
