<?php

namespace io\metamessage\jsonc;

class Jsonc {
    public static function parseFromString(string $s): JsoncNode {
        return parseJsonc($s);
    }

    public static function parseFromBytes(string $b): JsoncNode {
        return parseJsonc($b);
    }

    public static function toString(JsoncNode $n): string {
        return JsoncPrinter::toString($n);
    }

    public static function toCompactString(JsoncNode $n): string {
        return JsoncPrinter::toCompactString($n);
    }

    public static function print(JsoncNode $n): void {
        echo self::toString($n) . "\n";
    }
}
