<?php

namespace io\metamessage\core;

use io\metamessage\ir\Object_;
use io\metamessage\ir\Array_;
use io\metamessage\ir\Value;
use io\metamessage\ir\Node;
use io\metamessage\ir\Tag;
use io\metamessage\ir\Field;
use io\metamessage\ir\ValueType;
use io\metamessage\ir\Constants;

class WireDecoder
{
    private array $data;
    private int $offset;

    public function __construct(array $data)
    {
        $this->data = $data;
        $this->offset = 0;
    }

    public function decode(array $encoded): Node
    {
        $this->data = $encoded;
        $this->offset = 0;
        return $this->decodeNode(null, '')[0];
    }

    private function decodeNode(?Tag $tag, string $path): array
    {
        $b = $this->readByte();
        $prefix = $b & Prefix::PREFIX_MASK;

        switch ($prefix) {
            case Prefix::TAG:
                return $this->decodeTag($b, $path);
            case Prefix::SIMPLE:
                return $this->decodeSimple($b, $tag, $path);
            case Prefix::POSITIVE_INT:
                return $this->decodePositiveInt($b, $tag, $path);
            case Prefix::NEGATIVE_INT:
                return $this->decodeNegativeInt($b, $tag, $path);
            case Prefix::FLOAT:
                return $this->decodeFloat($b, $tag, $path);
            case Prefix::STRING:
                return $this->decodeString($b, $tag, $path);
            case Prefix::BYTES:
                return $this->decodeBytes($b, $tag, $path);
            case Prefix::CONTAINER:
                return $this->decodeContainer($b, $tag, $path);
            default:
                throw new MmDecodeException('invalid prefix');
        }
    }

    private function decodeTag(int $prefix, string $path): array
    {
        list($l1, $l2) = $this->tagLen($prefix);

        if ($l1 === 1) {
            $l2 = $this->readByte();
        } elseif ($l1 === 2) {
            $l = $this->readBytes(2);
            $l2 = ($l[0] << 8) | $l[1];
        } elseif ($l1 !== 0) {
            throw new MmDecodeException('invalid data');
        }

        $tag = Tag::newTag();

        $b = $this->readByte();
        $l = $b;
        if ($l < 254) {
        } elseif ($l < 257) {
            $l = $this->readByte();
        } else {
            $l3 = $this->readBytes(2);
            $l = ($l3[0] << 8) | $l3[1];
        }

        while ($l > 0) {
            $n = $this->decodeTagBytes($tag);
            if ($n === 0) {
                throw new MmDecodeException('tag error');
            }
            if ($n > $l) {
                throw new MmDecodeException('tag overflow');
            }
            $l -= $n;
        }

        $length = $l1 + 1 + $l2;

        if ($tag->isNull) {
            switch ($tag->type) {
                case ValueType::DATETIME:
                    $dt = new \DateTime('1970-01-01 00:00:00', new \DateTimeZone('UTC'));
                    $node = new Value($dt, $dt->format('Y-m-d H:i:s'), $tag, $path);
                    break;
                case ValueType::DATE:
                    $dt = new \DateTime('1970-01-01', new \DateTimeZone('UTC'));
                    $node = new Value($dt, $dt->format('Y-m-d'), $tag, $path);
                    break;
                case ValueType::TIME:
                    $dt = new \DateTime('1970-01-01 00:00:00', new \DateTimeZone('UTC'));
                    $node = new Value($dt, $dt->format('H:i:s'), $tag, $path);
                    break;
                case ValueType::I8:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::I16:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::I32:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::I64:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::U:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::U8:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::U16:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::U32:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::U64:
                    $node = new Value(0, '0', $tag, $path);
                    break;
                case ValueType::F32:
                    $node = new Value(0.0, '0.0', $tag, $path);
                    break;
                case ValueType::EMAIL:
                case ValueType::UUID:
                case ValueType::DECIMAL:
                    $node = new Value('', '', $tag, $path);
                    break;
                case ValueType::BIGINT:
                    $node = new Value('0', '0', $tag, $path);
                    break;
                case ValueType::URL:
                    $node = new Value('', '', $tag, $path);
                    break;
                case ValueType::IP:
                    switch ($tag->version) {
                        case 0:
                            $text = '';
                            break;
                        case 4:
                            $text = '0.0.0.0';
                            break;
                        case 6:
                            $text = '::';
                            break;
                        default:
                            throw new MmDecodeException('unsupported IP version: ' . $tag->version);
                    }
                    $node = new Value('', $text, $tag, $path);
                    break;
                default:
                    $result = $this->decodeNode($tag, $path);
                    $node = $result[0];
                    break;
            }
        } else {
            $result = $this->decodeNode($tag, $path);
            $node = $result[0];
        }

        return [$node, $length];
    }

    private function decodeTagBytes(Tag $tag): int
    {
        $b = $this->readByte();
        $prefix = $b & 0xF8;
        $remain = $b & 0x07;

        switch ($prefix) {
            case Tag::K_IS_NULL:
                $tag->isNull = ($remain & 0x01) === 1;
                if ($tag->isNull) {
                    $tag->nullable = true;
                }
                return 1;

            case Tag::K_EXAMPLE:
                $tag->example = ($remain & 0x01) === 1;
                return 1;

            case Tag::K_DESC:
                $n = 1;
                if ($remain <= 5) {
                    $desc = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $desc .= chr($this->readByte());
                        $n++;
                    }
                    $tag->desc = $desc;
                } elseif ($remain <= 6) {
                    $l = $this->readByte();
                    $n++;
                    $desc = '';
                    for ($i = 0; $i < $l; $i++) {
                        $desc .= chr($this->readByte());
                        $n++;
                    }
                    $tag->desc = $desc;
                } elseif ($remain <= 7) {
                    $lh = $this->readByte();
                    $ll = $this->readByte();
                    $l = ($lh << 8) | $ll;
                    $n += 2;
                    $desc = '';
                    for ($i = 0; $i < $l; $i++) {
                        $desc .= chr($this->readByte());
                        $n++;
                    }
                    $tag->desc = $desc;
                }
                return $n;

            case Tag::K_TYPE:
                $b2 = $this->readByte();
                $tag->type = ValueType::fromCode($b2);
                return 2;

            case Tag::K_RAW:
                $tag->raw = ($remain & 0x01) === 1;
                return 1;

            case Tag::K_NULLABLE:
                $tag->nullable = ($remain & 0x01) === 1;
                return 1;

            case Tag::K_ALLOW_EMPTY:
                $tag->allowEmpty = true;
                return 1;

            case Tag::K_UNIQUE:
                $tag->unique = true;
                return 1;

            case Tag::K_DEFAULT:
                $n = 1;
                if ($remain < 7) {
                    $default = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $default .= chr($this->readByte());
                        $n++;
                    }
                    $tag->defaultValue = $default;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $default = '';
                    for ($i = 0; $i < $l; $i++) {
                        $default .= chr($this->readByte());
                        $n++;
                    }
                    $tag->defaultValue = $default;
                }
                return $n;

            case Tag::K_MIN:
                $n = 1;
                if ($remain < 7) {
                    $min = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $min .= chr($this->readByte());
                        $n++;
                    }
                    $tag->min = $min;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $min = '';
                    for ($i = 0; $i < $l; $i++) {
                        $min .= chr($this->readByte());
                        $n++;
                    }
                    $tag->min = $min;
                }
                return $n;

            case Tag::K_MAX:
                $n = 1;
                if ($remain < 7) {
                    $max = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $max .= chr($this->readByte());
                        $n++;
                    }
                    $tag->max = $max;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $max = '';
                    for ($i = 0; $i < $l; $i++) {
                        $max .= chr($this->readByte());
                        $n++;
                    }
                    $tag->max = $max;
                }
                return $n;

            case Tag::K_SIZE:
                $n = 1;
                $size = 0;
                for ($i = 0; $i <= $remain; $i++) {
                    $size = ($size << 8) | $this->readByte();
                    $n++;
                }
                $tag->size = $size;
                return $n;

            case Tag::K_ENUM:
                $n = 1;
                $tag->type = ValueType::ENUM;
                if ($remain <= 5) {
                    $enum = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $enum .= chr($this->readByte());
                        $n++;
                    }
                    $tag->enumValues = $enum;
                } elseif ($remain === 6) {
                    $l = $this->readByte();
                    $n++;
                    $enum = '';
                    for ($i = 0; $i < $l; $i++) {
                        $enum .= chr($this->readByte());
                        $n++;
                    }
                    $tag->enumValues = $enum;
                } elseif ($remain === 7) {
                    $lh = $this->readByte();
                    $ll = $this->readByte();
                    $l = ($lh << 8) | $ll;
                    $n += 2;
                    $enum = '';
                    for ($i = 0; $i < $l; $i++) {
                        $enum .= chr($this->readByte());
                        $n++;
                    }
                    $tag->enumValues = $enum;
                }
                return $n;

            case Tag::K_PATTERN:
                $n = 1;
                if ($remain < 7) {
                    $pattern = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $pattern .= chr($this->readByte());
                        $n++;
                    }
                    $tag->pattern = $pattern;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $pattern = '';
                    for ($i = 0; $i < $l; $i++) {
                        $pattern .= chr($this->readByte());
                        $n++;
                    }
                    $tag->pattern = $pattern;
                }
                return $n;

            case Tag::K_LOCATION:
                $n = 1;
                $loc = '';
                for ($i = 0; $i < $remain; $i++) {
                    $loc .= chr($this->readByte());
                    $n++;
                }
                $tag->locationHours = (int)$loc;
                return $n;

            case Tag::K_VERSION:
                $n = 1;
                $ver = 0;
                for ($i = 0; $i <= $remain; $i++) {
                    $ver = ($ver << 8) | $this->readByte();
                    $n++;
                }
                $tag->version = $ver;
                return $n;

            case Tag::K_MIME:
                if ($remain < 7) {
                    $tag->mime = $remain;
                    return 1;
                } else {
                    $l2 = $this->readByte();
                    $tag->mime = $l2;
                    return 2;
                }

            case Tag::K_CHILD_DESC:
                $n = 1;
                if ($remain <= 5) {
                    $childDesc = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $childDesc .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childDesc = $childDesc;
                } elseif ($remain === 6) {
                    $l = $this->readByte();
                    $n++;
                    $childDesc = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childDesc .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childDesc = $childDesc;
                } elseif ($remain === 7) {
                    $lh = $this->readByte();
                    $ll = $this->readByte();
                    $l = ($lh << 8) | $ll;
                    $n += 2;
                    $childDesc = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childDesc .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childDesc = $childDesc;
                }
                return $n;

            case Tag::K_CHILD_TYPE:
                $b2 = $this->readByte();
                $tag->childType = ValueType::fromCode($b2);
                return 2;

            case Tag::K_CHILD_RAW:
                $tag->childRaw = ($remain & 0x01) === 1;
                return 1;

            case Tag::K_CHILD_NULLABLE:
                $tag->childNullable = ($remain & 0x01) === 1;
                return 1;

            case Tag::K_CHILD_ALLOW_EMPTY:
                $tag->childAllowEmpty = true;
                return 1;

            case Tag::K_CHILD_UNIQUE:
                $tag->childUnique = true;
                return 1;

            case Tag::K_CHILD_DEFAULT:
                $n = 1;
                if ($remain < 7) {
                    $childDefault = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $childDefault .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childDefault = $childDefault;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $childDefault = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childDefault .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childDefault = $childDefault;
                }
                return $n;

            case Tag::K_CHILD_MIN:
                $n = 1;
                if ($remain < 7) {
                    $childMin = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $childMin .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childMin = $childMin;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $childMin = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childMin .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childMin = $childMin;
                }
                return $n;

            case Tag::K_CHILD_MAX:
                $n = 1;
                if ($remain < 7) {
                    $childMax = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $childMax .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childMax = $childMax;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $childMax = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childMax .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childMax = $childMax;
                }
                return $n;

            case Tag::K_CHILD_SIZE:
                $n = 1;
                $childSize = 0;
                for ($i = 0; $i <= $remain; $i++) {
                    $childSize = ($childSize << 8) | $this->readByte();
                    $n++;
                }
                $tag->childSize = $childSize;
                return $n;

            case Tag::K_CHILD_ENUM:
                $n = 1;
                $tag->childType = ValueType::ENUM;
                if ($remain <= 5) {
                    $childEnum = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $childEnum .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childEnum = $childEnum;
                } elseif ($remain === 6) {
                    $l = $this->readByte();
                    $n++;
                    $childEnum = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childEnum .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childEnum = $childEnum;
                } elseif ($remain === 7) {
                    $lh = $this->readByte();
                    $ll = $this->readByte();
                    $l = ($lh << 8) | $ll;
                    $n += 2;
                    $childEnum = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childEnum .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childEnum = $childEnum;
                }
                return $n;

            case Tag::K_CHILD_PATTERN:
                $n = 1;
                if ($remain < 7) {
                    $childPattern = '';
                    for ($i = 0; $i < $remain; $i++) {
                        $childPattern .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childPattern = $childPattern;
                } else {
                    $l = $this->readByte();
                    $n++;
                    $childPattern = '';
                    for ($i = 0; $i < $l; $i++) {
                        $childPattern .= chr($this->readByte());
                        $n++;
                    }
                    $tag->childPattern = $childPattern;
                }
                return $n;

            case Tag::K_CHILD_LOCATION:
                $n = 1;
                $loc = '';
                for ($i = 0; $i < $remain; $i++) {
                    $loc .= chr($this->readByte());
                    $n++;
                }
                $tag->childLocationHours = (int)$loc;
                return $n;

            case Tag::K_CHILD_VERSION:
                $n = 1;
                $ver = 0;
                for ($i = 0; $i <= $remain; $i++) {
                    $ver = ($ver << 8) | $this->readByte();
                    $n++;
                }
                $tag->childVersion = $ver;
                return $n;

            case Tag::K_CHILD_MIME:
                if ($remain < 7) {
                    $tag->childMime = $remain;
                    return 1;
                } else {
                    $l2 = $this->readByte();
                    $tag->childMime = $l2;
                    return 2;
                }

            default:
                throw new MmDecodeException('unsupported tag key: ' . $prefix);
        }
    }

    private function decodeString(int $prefix, ?Tag $tag, string $path): array
    {
        list($l1, $l2) = $this->stringLen($prefix);

        if ($l1 === 1) {
            $l2 = $this->readByte();
        } elseif ($l1 === 2) {
            $l = $this->readBytes(2);
            $l2 = ($l[0] << 8) | $l[1];
        }

        $bs = [];
        if ($l2 > 0) {
            $bs = $this->readBytes($l2);
        }

        $text = pack('C*', ...$bs);

        if ($tag === null) {
            $tag = Tag::newTag();
        }

        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::STR;
        }

        switch ($tag->type) {
            case ValueType::EMAIL:
                $data = $text;
                break;
            case ValueType::URL:
                $data = $text;
                break;
            case ValueType::IP:
                $data = $text;
                break;
            case ValueType::STR:
                $data = $text;
                break;
            default:
                throw new MmDecodeException('unsupported string type: ' . $tag->type->name);
        }

        $node = new Value($data, $text, $tag, $path);
        $length = $l1 + 1 + $l2;
        return [$node, $length];
    }

    private function decodeBytes(int $prefix, ?Tag $tag, string $path): array
    {
        list($l1, $l2) = $this->bytesLen($prefix);

        if ($l1 === 1) {
            $l2 = $this->readByte();
        } elseif ($l1 === 2) {
            $l = $this->readBytes(2);
            $l2 = ($l[0] << 8) | $l[1];
        }

        $bs = [];
        if ($l2 > 0) {
            $bs = $this->readBytes($l2);
        }

        if ($tag === null) {
            $tag = Tag::newTag();
        }

        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::BYTES;
        }

        $text = '';
        switch ($tag->type) {
            case ValueType::BIGINT:
                $digits = BigIntWireCodec::decodePositive(array_slice($bs, 1), $bs[0]);
                $bits = $this->bigintBits(array_slice($bs, 1));
                $neg = !empty($bits) && $bits[0] === 1;
                $text = $neg ? '-' . $digits : $digits;
                $data = $text;
                break;

            case ValueType::BYTES:
                $data = $bs;
                $text = base64_encode(pack('C*', ...$bs));
                break;

            case ValueType::UUID:
                $data = $bs;
                $text = $this->uuidToString($bs);
                break;

            case ValueType::IP:
                $data = $bs;
                $text = $this->bytesToIPString($bs);
                break;

            default:
                throw new MmDecodeException('unsupported bytes type: ' . $tag->type->name);
        }

        $node = new Value($data, $text, $tag, $path);
        $length = $l1 + 1 + $l2;
        return [$node, $length];
    }

    private function decodePositiveInt(int $prefix, ?Tag $tag, string $path): array
    {
        list($l1, $l2) = $this->intLen($prefix);
        $v = $this->readUintBody($prefix, $l1, $l2);

        if ($tag === null) {
            $tag = Tag::newTag();
        }

        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::I;
        }

        $text = (string)$v;
        switch ($tag->type) {
            case ValueType::I:
                $data = (int)$v;
                break;
            case ValueType::I8:
                $data = (int)$v;
                break;
            case ValueType::I16:
                $data = (int)$v;
                break;
            case ValueType::I32:
                $data = (int)$v;
                break;
            case ValueType::I64:
                $data = (int)$v;
                break;
            case ValueType::U:
                $data = (int)$v;
                break;
            case ValueType::U8:
                $data = (int)$v;
                break;
            case ValueType::U16:
                $data = (int)$v;
                break;
            case ValueType::U32:
                $data = (int)$v;
                break;
            case ValueType::U64:
                $data = (int)$v;
                break;
            case ValueType::DATETIME:
                if ($tag->isNull) {
                    $data = null;
                    $text = '';
                } else {
                    if ($v > PHP_INT_MAX) {
                        throw new MmDecodeException('decodeDateTime: time value out of range');
                    }
                    $dt = new \DateTime('@' . $v, new \DateTimeZone('UTC'));
                    $data = $dt;
                    $text = $dt->format('Y-m-d H:i:s');
                }
                break;
            case ValueType::DATE:
                if ($tag->isNull) {
                    $data = null;
                    $text = '';
                } else {
                    $default = new \DateTime('1970-01-01 00:00:00', new \DateTimeZone('UTC'));
                    $dt = clone $default;
                    $dt->modify('+' . $v . ' days');
                    $dt->setTime(0, 0, 0);
                    $data = $dt;
                    $text = $dt->format('Y-m-d');
                }
                break;
            case ValueType::TIME:
                if ($tag->isNull) {
                    $data = null;
                    $text = '';
                } else {
                    if ($v > 86399) {
                        throw new MmDecodeException('decodeTime: time value out of range (0-86399)');
                    }
                    $hour = intdiv($v, 3600);
                    $minute = intdiv($v % 3600, 60);
                    $second = $v % 60;
                    $dt = new \DateTime(sprintf('1970-01-01 %02d:%02d:%02d', $hour, $minute, $second), new \DateTimeZone('UTC'));
                    $data = $dt;
                    $text = $dt->format('H:i:s');
                }
                break;
            case ValueType::ENUM:
                if ($tag->isNull) {
                    $data = -1;
                    $text = '';
                } else {
                    if ($tag->enumValues !== '') {
                        $enums = explode('|', $tag->enumValues);
                        $d = (int)$v;
                        if ($d >= count($enums)) {
                            throw new MmDecodeException('enum index out of range');
                        }
                        $data = $d;
                        $text = trim($enums[$d]);
                    } else {
                        throw new MmDecodeException('only enum are supported');
                    }
                }
                break;
            default:
                throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
        }

        $node = new Value($data, $text, $tag, $path);
        $length = $l1 + 1;
        return [$node, $length];
    }

    private function decodeNegativeInt(int $prefix, ?Tag $tag, string $path): array
    {
        list($l1, $l2) = $this->intLen($prefix);
        $v = $this->readUintBody($prefix, $l1, $l2);

        if ($tag === null) {
            $tag = Tag::newTag();
        }

        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::I;
        }

        $text = '-' . $v;

        switch ($tag->type) {
            case ValueType::I:
                $data = -(int)$v;
                break;
            case ValueType::I8:
                $data = -(int)$v;
                break;
            case ValueType::I16:
                $data = -(int)$v;
                break;
            case ValueType::I32:
                $data = -(int)$v;
                break;
            case ValueType::I64:
                $data = -(int)$v;
                break;
            case ValueType::DATETIME:
                if ($tag->isNull) {
                    $data = null;
                    $text = '';
                } else {
                    $dt = new \DateTime('@-' . $v, new \DateTimeZone('UTC'));
                    $data = $dt;
                    $text = $dt->format('Y-m-d H:i:s');
                }
                break;
            case ValueType::DATE:
                if ($tag->isNull) {
                    $data = null;
                    $text = '';
                } else {
                    $default = new \DateTime('1970-01-01 00:00:00', new \DateTimeZone('UTC'));
                    $dt = clone $default;
                    $dt->modify('-' . $v . ' days');
                    $dt->setTime(0, 0, 0);
                    $data = $dt;
                    $text = $dt->format('Y-m-d');
                }
                break;
            default:
                throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
        }

        $node = new Value($data, $text, $tag, $path);
        $length = $l1 + 1;
        return [$node, $length];
    }

    private function decodeFloat(int $prefix, ?Tag $tag, string $path): array
    {
        list($l1, $l2) = $this->floatLen($prefix);
        $p = $prefix;

        $v = 0.0;

        if ($p >= Prefix::FLOAT && $p <= Prefix::FLOAT + 7) {
            $v = ($p & 0xF) / 10;
            $length = 1;
        } elseif ($p >= (Prefix::FLOAT | WireConstants::FLOAT_NEG_MASK) && $p <= (Prefix::FLOAT | WireConstants::FLOAT_NEG_MASK) + 7) {
            $v = - (($p & 0xF) / 10);
            $length = 1;
        } else {
            $exp = $this->readByte();
            if ($exp > 127) {
                $exp -= 256;
            }
            $mantissa = $this->readUintBodyFromLen($l1, $l2);

            $dec = FloatCodec::mantissaToDecimal($mantissa, $exp);
            $v = (float)$dec;
            if (($p & WireConstants::FLOAT_NEG_MASK) !== 0) {
                $v = -$v;
            }
            $length = $l1 + 2;
        }

        if ($tag === null) {
            $tag = Tag::newTag();
        }

        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::F64;
        }

        switch ($tag->type) {
            case ValueType::F32:
                $data = (float)$v;
                $text = (string)$v;
                break;
            case ValueType::F64:
                $data = $v;
                $text = (string)$v;
                break;
            case ValueType::DECIMAL:
                $data = (string)$v;
                $text = (string)$v;
                break;
            default:
                throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
        }

        $node = new Value($data, $text, $tag, $path);
        return [$node, $length];
    }

    private function decodeSimple(int $prefix, ?Tag $tag, string $path): array
    {
        if ($tag === null) {
            $tag = Tag::newTag();
        }

        $sv = $prefix & Prefix::SUFFIX_MASK;

        switch ($sv) {
            case SimpleValue::FALSE:
                $tag->type = ValueType::BOOL;
                $node = new Value(false, Constants::FALSE, $tag, $path);
                break;
            case SimpleValue::TRUE:
                $tag->type = ValueType::BOOL;
                $node = new Value(true, Constants::TRUE, $tag, $path);
                break;
            case SimpleValue::NULL_BOOL:
                if ($tag->type === ValueType::UNKNOWN) {
                    $tag->type = ValueType::BOOL;
                } elseif ($tag->type !== ValueType::BOOL) {
                    throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
                }
                $node = new Value(false, Constants::FALSE, $tag, $path);
                break;
            case SimpleValue::NULL_INT:
                if ($tag->type === ValueType::UNKNOWN) {
                    $tag->type = ValueType::I;
                } elseif ($tag->type !== ValueType::I) {
                    throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
                }
                $node = new Value(0, '0', $tag, $path);
                break;
            case SimpleValue::NULL_FLOAT:
                if ($tag->type === ValueType::UNKNOWN) {
                    $tag->type = ValueType::F64;
                } elseif ($tag->type !== ValueType::F64 && $tag->type !== ValueType::F32) {
                    throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
                }
                $node = new Value(0.0, '0.0', $tag, $path);
                break;
            case SimpleValue::NULL_STRING:
                if ($tag->type === ValueType::UNKNOWN) {
                    $tag->type = ValueType::STR;
                } elseif ($tag->type !== ValueType::STR) {
                    throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
                }
                $node = new Value('', '', $tag, $path);
                break;
            case SimpleValue::NULL_BYTES:
                if ($tag->type === ValueType::UNKNOWN) {
                    $tag->type = ValueType::BYTES;
                } elseif ($tag->type !== ValueType::BYTES) {
                    throw new MmDecodeException('unsupported value types: ' . $tag->type->name);
                }
                $node = new Value([], '', $tag, $path);
                break;
            case SimpleValue::CODE:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_CODE_STR, $tag, $path);
                break;
            case SimpleValue::MESSAGE:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_MESSAGE_STR, $tag, $path);
                break;
            case SimpleValue::DATA:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_DATA_STR, $tag, $path);
                break;
            case SimpleValue::SUCCESS:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_SUCCESS_STR, $tag, $path);
                break;
            case SimpleValue::ERROR:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_ERROR_STR, $tag, $path);
                break;
            case SimpleValue::UNKNOWN:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_UNKNOWN_STR, $tag, $path);
                break;
            case SimpleValue::PAGE:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_PAGE_STR, $tag, $path);
                break;
            case SimpleValue::LIMIT:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_LIMIT_STR, $tag, $path);
                break;
            case SimpleValue::OFFSET:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_OFFSET_STR, $tag, $path);
                break;
            case SimpleValue::TOTAL:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_TOTAL_STR, $tag, $path);
                break;
            case SimpleValue::ID:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_ID_STR, $tag, $path);
                break;
            case SimpleValue::NAME:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_NAME_STR, $tag, $path);
                break;
            case SimpleValue::DESCRIPTION:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_DESCRIPTION_STR, $tag, $path);
                break;
            case SimpleValue::TYPE:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_TYPE_STR, $tag, $path);
                break;
            case SimpleValue::VERSION:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_VERSION_STR, $tag, $path);
                break;
            case SimpleValue::STATUS:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_STATUS_STR, $tag, $path);
                break;
            case SimpleValue::URL:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_URL_STR, $tag, $path);
                break;
            case SimpleValue::CREATE_TIME:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_CREATE_TIME_STR, $tag, $path);
                break;
            case SimpleValue::UPDATE_TIME:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_UPDATE_TIME_STR, $tag, $path);
                break;
            case SimpleValue::DELETE_TIME:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_DELETE_TIME_STR, $tag, $path);
                break;
            case SimpleValue::ACCOUNT:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_ACCOUNT_STR, $tag, $path);
                break;
            case SimpleValue::TOKEN:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_TOKEN_STR, $tag, $path);
                break;
            case SimpleValue::EXPIRE_TIME:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_EXPIRE_TIME_STR, $tag, $path);
                break;
            case SimpleValue::KEY:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_KEY_STR, $tag, $path);
                break;
            case SimpleValue::VAL:
                $tag->type = ValueType::STR;
                $node = new Value(null, Constants::SIMPLE_VAL_STR, $tag, $path);
                break;
            default:
                throw new MmDecodeException('unsupported value: ' . $prefix);
        }

        return [$node, 1];
    }

    private function decodeContainer(int $prefix, ?Tag $tag, string $path): array
    {
        if (($prefix & WireConstants::CONTAINER_MASK) === WireConstants::CONTAINER_ARRAY) {
            return $this->decodeArr($prefix, $tag, $path);
        }
        return $this->decodeObj($prefix, $tag, $path);
    }

    private function decodeArr(int $prefix, ?Tag $tag, string $path): array
    {
        if ($tag === null) {
            $tag = Tag::newTag();
            $tag->type = ValueType::VEC;
        }
        if ($tag->type === ValueType::UNKNOWN) {
            if ($tag->size > 0) {
                $tag->type = ValueType::ARR;
            } else {
                $tag->type = ValueType::VEC;
            }
        }

        list($l1, $l2) = $this->containerLen($prefix);

        if ($l1 === 1) {
            if (count($this->data) < $this->offset + 1) {
                throw new MmDecodeException($path . ': invalid data');
            }
            $l2 = $this->readByte();
        } elseif ($l1 === 2) {
            if (count($this->data) < $this->offset + 2) {
                throw new MmDecodeException($path . ': invalid data');
            }
            $l = $this->readBytes(2);
            $l2 = ($l[0] << 8) | $l[1];
        }

        $node = new Array_();
        $node->Tag = $tag;
        $node->Path = $path;

        $index = 0;
        while ($index < $l2) {
            $tagValue = Tag::newTag();
            $tagValue->inherit($tag);

            $p = $path . '[' . $index . ']';
            list($n, $l) = $this->decodeNode($tagValue, $p);
            if ($l <= 0) {
                throw new MmDecodeException($p . ': decode error');
            }

            $node->Items[] = $n;
            $index += $l;
        }

        $length = $l1 + 1 + $l2;
        return [$node, $length];
    }

    private function decodeObj(int $prefix, ?Tag $tag, string $path): array
    {
        if ($tag === null) {
            $tag = Tag::newTag();
            $tag->type = ValueType::OBJ;
        }
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::OBJ;
        }

        list($l1, $l2) = $this->containerLen($prefix);

        if ($l1 === 1) {
            if (count($this->data) < $this->offset + 1) {
                throw new MmDecodeException($path . ': invalid data');
            }
            $l2 = $this->readByte();
        } elseif ($l1 === 2) {
            if (count($this->data) < $this->offset + 2) {
                throw new MmDecodeException($path . ': invalid data');
            }
            $l = $this->readBytes(2);
            $l2 = ($l[0] << 8) | $l[1];
        }

        $node = new Object_();
        $node->Tag = $tag;
        $node->Path = $path;

        $lArray = $this->readByte();
        list($nKeys, $lKeys) = $this->decodeArr($lArray, $tag, $path);

        $index = $lKeys;
        $i = 0;
        while ($index < $l2) {
            $tagValue = Tag::newTag();
            $tagValue->inherit($tag);

            $key = $nKeys->Items[$i]->Text;
            $p = $path . '.' . $key;
            list($n, $l) = $this->decodeNode($tagValue, $p);
            if ($l <= 0) {
                throw new MmDecodeException($p . ': decode error');
            }

            $node->Fields[] = new Field($key, $n);
            $index += $l;
            $i++;
        }

        $length = $l1 + 1 + $l2;
        return [$node, $length];
    }

    private function readByte(): int
    {
        if ($this->offset >= count($this->data)) {
            throw new MmDecodeException('EOF');
        }
        $b = $this->data[$this->offset];
        $this->offset++;
        return $b;
    }

    private function readBytes(int $n): array
    {
        if ($n < 0) {
            throw new MmDecodeException('invalid length');
        }
        if ($this->offset + $n > count($this->data)) {
            throw new MmDecodeException('EOF');
        }
        $bs = array_slice($this->data, $this->offset, $n);
        $this->offset += $n;
        return $bs;
    }

    private function readUintBody(int $prefix, int $l1, int $l2): int
    {
        if ($l1 === 0) {
            return $l2;
        }

        $v = 0;
        for ($i = 0; $i < $l1; $i++) {
            $v = ($v << 8) | $this->readByte();
        }
        return $v;
    }

    private function readUintBodyFromLen(int $l1, int $l2): int
    {
        if ($l1 === 0) {
            return $l2;
        }

        $v = 0;
        for ($i = 0; $i < $l1; $i++) {
            $v = ($v << 8) | $this->readByte();
        }
        return $v;
    }

    private function tagLen(int $first): array
    {
        $l = $first & WireConstants::TAG_LEN_MASK;
        if ($l === WireConstants::TAG_LEN_1) {
            return [1, 0];
        } elseif ($l === WireConstants::TAG_LEN_2) {
            return [2, 0];
        } else {
            return [0, $l];
        }
    }

    private function containerLen(int $first): array
    {
        $l = $first & WireConstants::CONTAINER_LEN_MASK;
        if ($l === WireConstants::CONTAINER_LEN_1) {
            return [1, 0];
        } elseif ($l === WireConstants::CONTAINER_LEN_2) {
            return [2, 0];
        } else {
            return [0, $l];
        }
    }

    private function stringLen(int $first): array
    {
        $l = $first & WireConstants::STRING_LEN_MASK;
        if ($l < WireConstants::STRING_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::STRING_LEN_1) {
            return [1, $l];
        } else {
            return [2, $l];
        }
    }

    private function bytesLen(int $first): array
    {
        $l = $first & WireConstants::BYTES_LEN_MASK;
        if ($l < WireConstants::BYTES_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::BYTES_LEN_1) {
            return [1, $l];
        } else {
            return [2, $l];
        }
    }

    private function intLen(int $first): array
    {
        $l = $first & WireConstants::INT_LEN_MASK;
        if ($l < WireConstants::INT_LEN_1) {
            return [0, $l];
        } else {
            return [$l - WireConstants::INT_LEN_1 + 1, 0];
        }
    }

    private function floatLen(int $first): array
    {
        $l = $first & WireConstants::FLOAT_LEN_MASK;
        if ($l < WireConstants::FLOAT_LEN_1) {
            return [0, $l];
        } else {
            return [$l - WireConstants::FLOAT_LEN_1 + 1, 0];
        }
    }

    private function uuidToString(array $bs): string
    {
        if (count($bs) !== 16) {
            throw new MmDecodeException('UUID length');
        }
        $hex = '';
        for ($i = 0; $i < 16; $i++) {
            $hex .= sprintf('%02x', $bs[$i]);
        }
        return substr($hex, 0, 8) . '-' . substr($hex, 8, 4) . '-' . substr($hex, 12, 4) . '-' . substr($hex, 16, 4) . '-' . substr($hex, 20);
    }

    private function bytesToIPString(array $bs): string
    {
        if (count($bs) === 4) {
            return implode('.', $bs);
        } elseif (count($bs) === 16) {
            $parts = [];
            for ($i = 0; $i < 16; $i += 2) {
                $parts[] = sprintf('%x', ($bs[$i] << 8) | $bs[$i + 1]);
            }
            return implode(':', $parts);
        }
        return '';
    }

    private function bigintBits(array $data): array
    {
        $bits = [];
        foreach ($data as $bt) {
            for ($i = 7; $i >= 0; $i--) {
                $bits[] = ($bt >> $i) & 1;
            }
        }
        return $bits;
    }
}
