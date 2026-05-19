<?php

namespace io\metamessage\jsonc;

use io\metamessage\ir\Object_;
use io\metamessage\ir\Array_;
use io\metamessage\ir\Value;
use io\metamessage\ir\Node;
use io\metamessage\ir\Tag;
use io\metamessage\ir\ValueType;

class JsoncPrinter
{
    private const INDENT_UNIT = "\t";

    public static function ToJSONC(Node $n): string
    {
        if ($n === null) {
            return '';
        }
        return self::writeNodeJSONC($n, 0);
    }

    private static function writeNodeJSONC(Node $n, int $indent): string
    {
        if ($n instanceof Value) {
            return self::writeValueJSONC($n);
        } elseif ($n instanceof Object_) {
            return self::writeObjectJSONC($n, $indent);
        } elseif ($n instanceof Array_) {
            return self::writeArrayJSONC($n, $indent);
        }
        return '';
    }

    private static function writeValueJSONC(Value $v): string
    {
        if ($v->getTag() === null) {
            return '';
        }

        switch ($v->getTag()->type) {
            case ValueType::STRING:
            case ValueType::BYTES:
            case ValueType::DATETIME:
            case ValueType::DATE:
            case ValueType::TIME:
            case ValueType::UUID:
            case ValueType::IP:
            case ValueType::URL:
            case ValueType::EMAIL:
            case ValueType::ENUM:
                return '"' . $v->Text . '"';

            case ValueType::INT:
            case ValueType::INT8:
            case ValueType::INT16:
            case ValueType::INT32:
            case ValueType::INT64:
            case ValueType::UINT:
            case ValueType::UINT8:
            case ValueType::UINT16:
            case ValueType::UINT32:
            case ValueType::UINT64:
            case ValueType::BIGINT:
            case ValueType::DECIMAL:
            case ValueType::BOOL:
            case ValueType::FLOAT32:
            case ValueType::FLOAT64:
                return $v->Text;

            default:
                return $v->Text;
        }
    }

    private static function writeObjectJSONC(Object_ $o, int $indent): string
    {
        $sb = "{\n";

        foreach ($o->Fields as $f) {
            $val = $f->Value;
            if ($val !== null) {
                $sb .= self::writeLeadingComments($val->getTag(), $indent + 1);
            }

            self::writeIndent($sb, $indent + 1);

            $sb .= '"' . $f->Key . '": ';

            if ($val !== null) {
                $sb .= self::writeNodeJSONC($val, $indent + 1);
            } else {
                $sb .= 'null';
            }

            $sb .= ",\n";
        }

        self::writeIndent($sb, $indent);
        $sb .= '}';
        return $sb;
    }

    private static function writeArrayJSONC(Array_ $a, int $indent): string
    {
        $sb = "[\n";

        foreach ($a->Items as $item) {
            $sb .= self::writeLeadingComments($item->getTag(), $indent + 1);

            self::writeIndent($sb, $indent + 1);

            $sb .= self::writeNodeJSONC($item, $indent + 1);

            $sb .= ",\n";
        }

        self::writeIndent($sb, $indent);
        $sb .= ']';
        return $sb;
    }

    private static function writeLeadingComments(?Tag $tag, int $indent): string
    {
        $sb = '';
        if ($tag === null) {
            return $sb;
        }
        $tagStr = $tag->toString();
        if ($tagStr !== '') {
            $sb .= "\n";
            self::writeIndent($sb, $indent);
            $sb .= '// mm: ' . $tagStr . "\n";
        }
        return $sb;
    }

    private static function writeIndent(string &$sb, int $indent): void
    {
        $sb .= str_repeat(self::INDENT_UNIT, $indent);
    }
}
