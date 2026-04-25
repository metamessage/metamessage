<?php

namespace io\metamessage\mm;

class Prefix {
    public const SIMPLE = 0b000 << 5;
    public const POSITIVE_INT = 0b001 << 5;
    public const NEGATIVE_INT = 0b010 << 5;
    public const FLOAT = 0b011 << 5;
    public const STRING = 0b100 << 5;
    public const BYTES = 0b101 << 5;
    public const CONTAINER = 0b110 << 5;
    public const TAG = 0b111 << 5;

    public const PREFIX_MASK = 0b11100000;
    public const SUFFIX_MASK = 0b00011111;

    public static function of(int $b): int {
        return $b & self::PREFIX_MASK;
    }

    public static function suffix(int $b): int {
        return $b & self::SUFFIX_MASK;
    }
}
