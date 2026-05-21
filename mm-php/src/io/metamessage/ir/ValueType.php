<?php

namespace io\metamessage\ir;

enum ValueType
{
    case UNKNOWN;
    case DOC;
    case VEC;
    case ARR;
    case OBJ;
    case MAP;
    case STR;
    case BYTES;
    case BOOL;
    case I;
    case I8;
    case I16;
    case I32;
    case I64;
    case U;
    case U8;
    case U16;
    case U32;
    case U64;
    case F32;
    case F64;
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

    public function code(): int
    {
        return array_search($this, self::cases());
    }

    public static function fromCode(int $code): self
    {
        $cases = self::cases();
        return $cases[$code] ?? self::UNKNOWN;
    }

    public function wireName(): string
    {
        return match ($this) {
            self::UNKNOWN => 'unknown',
            self::DOC => 'doc',
            self::VEC => 'vec',
            self::ARR => 'arr',
            self::OBJ => 'obj',
            self::MAP => 'map',
            self::STR => 'str',
            self::BYTES => 'bytes',
            self::BOOL => 'bool',
            self::I => 'i',
            self::I8 => 'i8',
            self::I16 => 'i16',
            self::I32 => 'i32',
            self::I64 => 'i64',
            self::U => 'u',
            self::U8 => 'u8',
            self::U16 => 'u16',
            self::U32 => 'u32',
            self::U64 => 'u64',
            self::F32 => 'f32',
            self::F64 => 'f64',
            self::BIGINT => 'bigint',
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

    public static function parseWireName(string $s): self
    {
        if (empty($s)) {
            return self::UNKNOWN;
        }
        $key = strtolower(trim($s));
        if ($key === 'struct') {
            return self::OBJ;
        }
        $wireNameMap = array_reduce(self::cases(), function ($map, $case) {
            $map[strtolower($case->wireName())] = $case;
            return $map;
        }, []);
        return $wireNameMap[$key] ?? self::UNKNOWN;
    }
}
