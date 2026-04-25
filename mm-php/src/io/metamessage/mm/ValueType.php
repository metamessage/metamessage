<?php

namespace io\metamessage\mm;

enum ValueType {
    case UNKNOWN;
    case DOC;
    case SLICE;
    case ARRAY;
    case STRUCT;
    case MAP;
    case STRING;
    case BYTES;
    case BOOL;
    case INT;
    case INT8;
    case INT16;
    case INT32;
    case INT64;
    case UINT;
    case UINT8;
    case UINT16;
    case UINT32;
    case UINT64;
    case FLOAT32;
    case FLOAT64;
    case BIGINT;
    case DATETIME;
    case DATE;
    case TIME;
    case UUID;
    case DECIMAL;
    case IP;
    case URL;
    case EMAIL;
    case ENUM;
    case IMAGE;
    case VIDEO;

    public function code(): int {
        return array_search($this, self::cases());
    }

    public static function fromCode(int $code): self {
        $cases = self::cases();
        return $cases[$code] ?? self::UNKNOWN;
    }

    public function wireName(): string {
        return match($this) {
            self::UNKNOWN => 'unknown',
            self::DOC => 'doc',
            self::SLICE => 'slice',
            self::ARRAY => 'arr',
            self::STRUCT => 'struct',
            self::MAP => 'map',
            self::STRING => 'str',
            self::BYTES => 'bytes',
            self::BOOL => 'bool',
            self::INT => 'i',
            self::INT8 => 'i8',
            self::INT16 => 'i16',
            self::INT32 => 'i32',
            self::INT64 => 'i64',
            self::UINT => 'u',
            self::UINT8 => 'u8',
            self::UINT16 => 'u16',
            self::UINT32 => 'u32',
            self::UINT64 => 'u64',
            self::FLOAT32 => 'f32',
            self::FLOAT64 => 'f64',
            self::BIGINT => 'bi',
            self::DATETIME => 'datetime',
            self::DATE => 'date',
            self::TIME => 'time',
            self::UUID => 'uuid',
            self::DECIMAL => 'decimal',
            self::IP => 'ip',
            self::URL => 'url',
            self::EMAIL => 'email',
            self::ENUM => 'enum',
            self::IMAGE => 'image',
            self::VIDEO => 'video',
        };
    }

    public static function parseWireName(string $s): self {
        if (empty($s)) {
            return self::UNKNOWN;
        }
        $wireNameMap = array_reduce(self::cases(), function($map, $case) {
            $map[strtolower($case->wireName())] = $case;
            return $map;
        }, []);
        return $wireNameMap[strtolower(trim($s))] ?? self::UNKNOWN;
    }
}
