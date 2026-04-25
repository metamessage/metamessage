<?php

namespace io\metamessage\mm;

class WireConstants {
    public const MAX_1 = 0xFF;
    public const MAX_2 = 0xFFFF;
    public const MAX_3 = 0xFFFFFF;
    public const MAX_4 = 0xFFFFFFFF;
    public const MAX_5 = 0xFFFFFFFFFF;
    public const MAX_6 = 0xFFFFFFFFFFFF;
    public const MAX_7 = 0xFFFFFFFFFFFFFF;
    public const MAX_8 = 0xFFFFFFFFFFFFFFFF;

    public const INT_LEN_MASK = 0b11111;
    public const INT_LEN_1 = self::INT_LEN_MASK - 7;
    public const INT_LEN_2 = self::INT_LEN_MASK - 6;
    public const INT_LEN_3 = self::INT_LEN_MASK - 5;
    public const INT_LEN_4 = self::INT_LEN_MASK - 4;
    public const INT_LEN_5 = self::INT_LEN_MASK - 3;
    public const INT_LEN_6 = self::INT_LEN_MASK - 2;
    public const INT_LEN_7 = self::INT_LEN_MASK - 1;
    public const INT_LEN_8 = self::INT_LEN_MASK;

    public const FLOAT_NEG_MASK = 0b10000;
    public const FLOAT_LEN_MASK = 0b01111;
    public const FLOAT_LEN_1 = self::FLOAT_LEN_MASK - 7;
    public const FLOAT_LEN_2 = self::FLOAT_LEN_MASK - 6;
    public const FLOAT_LEN_3 = self::FLOAT_LEN_MASK - 5;
    public const FLOAT_LEN_4 = self::FLOAT_LEN_MASK - 4;
    public const FLOAT_LEN_5 = self::FLOAT_LEN_MASK - 3;
    public const FLOAT_LEN_6 = self::FLOAT_LEN_MASK - 2;
    public const FLOAT_LEN_7 = self::FLOAT_LEN_MASK - 1;
    public const FLOAT_LEN_8 = self::FLOAT_LEN_MASK;

    public const STRING_LEN_MASK = 0b11111;
    public const STRING_LEN_1 = self::STRING_LEN_MASK - 1;
    public const STRING_LEN_2 = self::STRING_LEN_MASK;

    public const BYTES_LEN_MASK = 0b11111;
    public const BYTES_LEN_1 = self::BYTES_LEN_MASK - 1;
    public const BYTES_LEN_2 = self::BYTES_LEN_MASK;

    public const CONTAINER_MASK = 0b10000;
    public const CONTAINER_MAP = 0b00000;
    public const CONTAINER_ARRAY = 0b10000;
    public const CONTAINER_LEN_MASK = 0b01111;
    public const CONTAINER_LEN_1 = self::CONTAINER_LEN_MASK - 1;
    public const CONTAINER_LEN_2 = self::CONTAINER_LEN_MASK;

    public const TAG_LEN_MASK = 0b11111;
    public const TAG_LEN_1 = self::TAG_LEN_MASK - 1;
    public const TAG_LEN_2 = self::TAG_LEN_MASK;
}
