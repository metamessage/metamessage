<?php

namespace io\metamessage\mm;

class TagFieldParser {

    public static function parse(array &$data, int &$offset): ?MmTag {
        $t = MmTag::empty();
        $maxOffset = count($data);
        $maxFieldCount = 255;
        $fieldCount = 0;

        while ($offset < $maxOffset) {
            if ($fieldCount++ > $maxFieldCount) {
                throw new MmDecodeException('Tag fields overflow');
            }

            $b = $data[$offset++];
            $k = $b >> 3;

            switch ($k) {
                case MmTag::K_IS_NULL:
                    $t->isNull = true;
                    break;
                case MmTag::K_EXAMPLE:
                    $t->example = true;
                    break;
                case MmTag::K_DESC:
                    $t->desc = self::readSizedString($data, $offset, $b);
                    break;
                case MmTag::K_TYPE:
                    $t->type = ValueType::fromCode($data[$offset++]);
                    break;
                case MmTag::K_RAW:
                    $t->raw = true;
                    break;
                case MmTag::K_NULLABLE:
                    $t->nullable = true;
                    break;
                case MmTag::K_ALLOW_EMPTY:
                    $t->allowEmpty = true;
                    break;
                case MmTag::K_UNIQUE:
                    $t->unique = true;
                    break;
                case MmTag::K_DEFAULT:
                    $t->defaultValue = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_MIN:
                    $t->min = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_MAX:
                    $t->max = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_SIZE:
                    $t->size = self::readUint64($data, $offset, $b);
                    break;
                case MmTag::K_ENUM:
                    $t->enumValues = self::readSizedString($data, $offset, $b);
                    $t->type = ValueType::ENUM;
                    break;
                case MmTag::K_PATTERN:
                    $t->pattern = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_LOCATION:
                    $n = $b & 7;
                    $t->locationHours = (int)self::readString($data, $offset, $n);
                    break;
                case MmTag::K_VERSION:
                    $t->version = self::readUint64($data, $offset, $b);
                    break;
                case MmTag::K_MIME:
                    $n = $b & 7;
                    if ($n < 7) {
                        $t->mime = MimeWire::toString($n);
                    } else {
                        $t->mime = MimeWire::toString($data[$offset++]);
                    }
                    break;
                case MmTag::K_CHILD_DESC:
                    $t->childDesc = self::readSizedString($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_TYPE:
                    $t->childType = ValueType::fromCode($data[$offset++]);
                    break;
                case MmTag::K_CHILD_RAW:
                    $t->childRaw = true;
                    break;
                case MmTag::K_CHILD_NULLABLE:
                    $t->childNullable = true;
                    break;
                case MmTag::K_CHILD_ALLOW_EMPTY:
                    $t->childAllowEmpty = true;
                    break;
                case MmTag::K_CHILD_UNIQUE:
                    $t->childUnique = true;
                    break;
                case MmTag::K_CHILD_DEFAULT:
                    $t->childDefault = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_MIN:
                    $t->childMin = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_MAX:
                    $t->childMax = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_SIZE:
                    $t->childSize = self::readUint64($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_ENUM:
                    $t->childEnum = self::readSizedString($data, $offset, $b);
                    $t->childType = ValueType::ENUM;
                    break;
                case MmTag::K_CHILD_PATTERN:
                    $t->childPattern = self::readShortString($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_LOCATION:
                    $n = $b & 7;
                    $t->childLocationHours = (int)self::readString($data, $offset, $n);
                    break;
                case MmTag::K_CHILD_VERSION:
                    $t->childVersion = self::readUint64($data, $offset, $b);
                    break;
                case MmTag::K_CHILD_MIME:
                    $n = $b & 7;
                    if ($n < 7) {
                        $t->childMime = MimeWire::toString($n);
                    } else {
                        $t->childMime = MimeWire::toString($data[$offset++]);
                    }
                    break;
                default:
                    // 对于未知的标签字段，我们已经在循环开始时增加了 $offset
                    // 所以这里不需要再增加，直接跳过
                    break;
            }
        }

        return $t;
    }

    private static function readString(array &$data, int &$offset, int $n): string {
        $s = '';
        for ($i = 0; $i < $n; $i++) {
            $s .= chr($data[$offset++]);
        }
        return $s;
    }

    private static function readSizedString(array &$data, int &$offset, int $b): string {
        $n = $b & 7;
        if ($n <= 5) {
            return self::readString($data, $offset, $n);
        } elseif ($n == 6) {
            $l = $data[$offset++];
            return self::readString($data, $offset, $l);
        } else {
            $l = ($data[$offset++] << 8) | $data[$offset++];
            return self::readString($data, $offset, $l);
        }
    }

    private static function readShortString(array &$data, int &$offset, int $b): string {
        $n = $b & 7;
        if ($n < 7) {
            return self::readString($data, $offset, $n);
        } else {
            $l = $data[$offset++];
            return self::readString($data, $offset, $l);
        }
    }

    private static function readUint64(array &$data, int &$offset, int $b): int {
        $n = $b & 7;
        switch ($n) {
            case 0:
                return $data[$offset++];
            case 1:
                return ($data[$offset++] << 8) | $data[$offset++];
            case 2:
                return ($data[$offset++] << 16) | ($data[$offset++] << 8) | $data[$offset++];
            case 3:
                return ($data[$offset++] << 24) | ($data[$offset++] << 16) | ($data[$offset++] << 8) | $data[$offset++];
            default:
                throw new MmDecodeException('Invalid uint64 encoding');
        }
    }
}
