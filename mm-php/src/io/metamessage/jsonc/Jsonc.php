<?php

namespace io\metamessage\jsonc;

use io\metamessage\ir\Node;
use io\metamessage\ir\NodeObject;
use io\metamessage\ir\NodeArray;
use io\metamessage\ir\NodeScalar;
use io\metamessage\ir\Field;

class_alias(NodeObject::class, 'io\\metamessage\\jsonc\\JsoncObject');
class_alias(NodeArray::class, 'io\\metamessage\\jsonc\\JsoncArray');
class_alias(NodeScalar::class, 'io\\metamessage\\jsonc\\JsoncValue');
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
