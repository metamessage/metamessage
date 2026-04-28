<?php

namespace io\metamessage\jsonc;

class JsoncPrinter {
    private const INDENT = "    ";

    public static function toString(JsoncNode $node, int $indentLevel = 0): string {
        if ($node instanceof JsoncObject) {
            return self::objectToString($node, $indentLevel);
        } elseif ($node instanceof JsoncArray) {
            return self::arrayToString($node, $indentLevel);
        } elseif ($node instanceof JsoncValue) {
            return self::valueToString($node, $indentLevel);
        }
        return "";
    }

    private static function objectToString(JsoncObject $obj, int $indentLevel): string {
        $sb = "{\n";

        foreach ($obj->fields as $field) {
            for ($i = 0; $i < $indentLevel + 1; $i++) {
                $sb .= self::INDENT;
            }

            $fieldValue = $field->value;
            if ($fieldValue instanceof JsoncValue && $fieldValue->tag !== null && $fieldValue->tag->desc !== "") {
                $sb .= "// mm: " . self::tagToString($fieldValue->tag) . "\n";
                for ($i = 0; $i < $indentLevel + 1; $i++) {
                    $sb .= self::INDENT;
                }
            }

            $sb .= "\"" . $field->key . "\": ";
            $sb .= self::toString($fieldValue ?? new JsoncValue(), $indentLevel + 1);
            $sb .= ",\n";
        }

        for ($i = 0; $i < $indentLevel; $i++) {
            $sb .= self::INDENT;
        }
        $sb .= "}";
        return $sb;
    }

    private static function arrayToString(JsoncArray $arr, int $indentLevel): string {
        $sb = "[\n";

        foreach ($arr->items as $item) {
            for ($i = 0; $i < $indentLevel + 1; $i++) {
                $sb .= self::INDENT;
            }
            $sb .= self::toString($item, $indentLevel + 1);
            $sb .= ",\n";
        }

        for ($i = 0; $i < $indentLevel; $i++) {
            $sb .= self::INDENT;
        }
        $sb .= "]";
        return $sb;
    }

    private static function valueToString(JsoncValue $value, int $indentLevel): string {
        $sb = "";
        if ($value->tag !== null && $value->tag->desc !== "") {
            $sb .= "// mm: " . self::tagToString($value->tag) . "\n";
            for ($i = 0; $i < $indentLevel; $i++) {
                $sb .= self::INDENT;
            }
        }
        $sb .= self::valueToStringOnly($value);
        return $sb;
    }

    private static function valueToStringOnly(JsoncValue $value): string {
        $tag = $value->tag;
        $type = $tag?->type ?? JsoncValueType::Unknown;

        $needsQuotes = match ($type) {
            JsoncValueType::String,
            JsoncValueType::Bytes,
            JsoncValueType::DateTime,
            JsoncValueType::Date,
            JsoncValueType::Time,
            JsoncValueType::UUID,
            JsoncValueType::IP,
            JsoncValueType::URL,
            JsoncValueType::Email,
            JsoncValueType::Enum => true,
            default => false,
        };

        if ($needsQuotes) {
            return "\"" . $value->text . "\"";
        } else {
            return $value->text;
        }
    }

    private static function tagToString(JsoncTag $tag): string {
        $parts = [];

        if ($tag->type !== JsoncValueType::Unknown) {
            $parts[] = "type=" . self::typeToString($tag->type);
        }
        if ($tag->desc !== "") {
            $parts[] = "desc=" . $tag->desc;
        }
        if ($tag->nullable) {
            $parts[] = "nullable";
        }
        if ($tag->isNull) {
            $parts[] = "is_null";
        }
        if ($tag->raw) {
            $parts[] = "raw";
        }
        if ($tag->defaultValue !== "") {
            $parts[] = "default=" . $tag->defaultValue;
        }
        if ($tag->enum !== "") {
            $parts[] = "enum=" . $tag->enum;
        }

        return implode("; ", $parts);
    }

    private static function typeToString(JsoncValueType $type): string {
        return match ($type) {
            JsoncValueType::String => "str",
            JsoncValueType::Int => "i",
            JsoncValueType::Int8 => "i8",
            JsoncValueType::Int16 => "i16",
            JsoncValueType::Int32 => "i32",
            JsoncValueType::Int64 => "i64",
            JsoncValueType::Uint => "u",
            JsoncValueType::Uint8 => "u8",
            JsoncValueType::Uint16 => "u16",
            JsoncValueType::Uint32 => "u32",
            JsoncValueType::Uint64 => "u64",
            JsoncValueType::Float32 => "f32",
            JsoncValueType::Float64 => "f64",
            JsoncValueType::Bool => "bool",
            JsoncValueType::Bytes => "bytes",
            JsoncValueType::BigInt => "bi",
            JsoncValueType::DateTime => "datetime",
            JsoncValueType::Date => "date",
            JsoncValueType::Time => "time",
            JsoncValueType::UUID => "uuid",
            JsoncValueType::Decimal => "decimal",
            JsoncValueType::IP => "ip",
            JsoncValueType::URL => "url",
            JsoncValueType::Email => "email",
            JsoncValueType::Enum => "enum",
            JsoncValueType::Array => "arr",
            JsoncValueType::Struct => "obj",
            JsoncValueType::Null => "null",
            JsoncValueType::Unknown => "unknown",
            default => "unknown",
        };
    }

    public static function toCompactString(JsoncNode $node): string {
        if ($node instanceof JsoncObject) {
            return self::compactObject($node);
        } elseif ($node instanceof JsoncArray) {
            return self::compactArray($node);
        } elseif ($node instanceof JsoncValue) {
            return self::compactValue($node);
        }
        return "";
    }

    private static function compactObject(JsoncObject $obj): string {
        $fields = array_map(
            fn(JsoncField $f) => "\"" . $f->key . "\": " . self::toCompactString($f->value ?? new JsoncValue()),
            $obj->fields
        );
        return "{" . implode(",", $fields) . "}";
    }

    private static function compactArray(JsoncArray $arr): string {
        $items = array_map(
            fn(JsoncNode $item) => self::toCompactString($item),
            $arr->items
        );
        return "[" . implode(",", $items) . "]";
    }

    private static function compactValue(JsoncValue $value): string {
        return self::valueToStringOnly($value);
    }
}

function printJsonc(JsoncNode $node): string {
    return JsoncPrinter::toString($node);
}

function printJsoncCompact(JsoncNode $node): string {
    return JsoncPrinter::toCompactString($node);
}
