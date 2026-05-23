<?php

namespace io\metamessage\core;

use io\metamessage\ir\Object_;
use io\metamessage\ir\Array_;
use io\metamessage\ir\Value;
use io\metamessage\ir\Node;
use io\metamessage\ir\Tag;
use io\metamessage\ir\ValueType;

class WireEncoder
{
    private array $buf;
    private int $offset;
    private int $maxCap;

    private const DEFAULT_BUF_SIZE = 1024;
    private const MAX_CAP = 1024 * 1024 * 1024;

    public function __construct()
    {
        $this->buf = array_fill(0, self::DEFAULT_BUF_SIZE, 0);
        $this->offset = 0;
        $this->maxCap = self::MAX_CAP;
    }

    public function reset(): void
    {
        $this->offset = 0;
    }

    public function encode(Node $node): array
    {
        if ($node instanceof Object_) {
            $n = $this->encodeNodeObject($node);
        } elseif ($node instanceof Array_) {
            $n = $this->encodeNodeArray($node);
        } elseif ($node instanceof Value) {
            $n = $this->encodeNodeValue($node);
        } else {
            throw new \Exception('encode error: unsupported type ' . get_class($node));
        }

        $out = array_slice($this->buf, $this->offset - $n, $n);
        $this->offset = 0;
        return $out;
    }

    private function encodeNodeObject(Object_ $obj): int
    {
        $bufKey = [];
        $buf = [];
        $tag = $obj->getTag();

        foreach ($obj->Fields as $field) {
            $val = $field->Value;
            if ($val instanceof Object_) {
                $n = $this->encodeNodeObject($val);
            } elseif ($val instanceof Array_) {
                $n = $this->encodeNodeArray($val);
            } elseif ($val instanceof Value) {
                $n = $this->encodeNodeValue($val);
            } else {
                throw new \Exception('unsupported type ' . get_class($val));
            }

            $encodedSub = array_slice($this->buf, $this->offset - $n, $n);
            array_splice($buf, count($buf), 0, $encodedSub);

            $ns = $this->encodeString($field->Key);
            $encodedKey = array_slice($this->buf, $this->offset - $ns, $ns);
            array_splice($bufKey, count($bufKey), 0, $encodedKey);
        }

        $nk = $this->encodeArray($bufKey);
        $encodedKeyArray = array_slice($this->buf, $this->offset - $nk, $nk);
        $bufAll = array_merge($encodedKeyArray, $buf);

        $n = $this->encodeObject($bufAll);

        $n1 = $this->encodeComment(array_slice($this->buf, $this->offset - $n, $n), $tag);
        if ($n1 === 0) {
            return $n;
        }
        return $n1;
    }

    private function encodeComment(array $payload, ?Tag $tag): int
    {
        if ($tag === null) {
            return 0;
        }

        $ns = $this->encodeT($tag->bytes());
        if ($ns === 0) {
            return 0;
        }

        return $this->encodeTag($payload, array_slice($this->buf, $this->offset - $ns, $ns));
    }

    private function encodeNodeArray(Array_ $arr): int
    {
        $buf = [];
        $tag = $arr->getTag();

        foreach ($arr->Items as $item) {
            if ($item instanceof Object_) {
                $n = $this->encodeNodeObject($item);
            } elseif ($item instanceof Array_) {
                $n = $this->encodeNodeArray($item);
            } elseif ($item instanceof Value) {
                $n = $this->encodeNodeValue($item);
            } else {
                throw new \Exception('unsupported type ' . get_class($item));
            }

            $encodedSub = array_slice($this->buf, $this->offset - $n, $n);
            array_splice($buf, count($buf), 0, $encodedSub);
        }

        $n = $this->encodeArray($buf);

        $n1 = $this->encodeComment(array_slice($this->buf, $this->offset - $n, $n), $tag);
        if ($n1 === 0) {
            return $n;
        }
        return $n1;
    }

    private function encodeNodeValue(Value $val): int
    {
        $tag = $val->getTag();

        switch ($tag->type) {
            case ValueType::DATETIME:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeDateTime($val->Data);
                }
                break;

            case ValueType::DATE:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeDate($val->Data);
                }
                break;

            case ValueType::TIME:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeTime($val->Data);
                }
                break;

            case ValueType::I:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_INT);
                } else {
                    $n = $this->encodeInt64($val->Data);
                }
                break;

            case ValueType::I8:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeInt64($val->Data);
                }
                break;

            case ValueType::I16:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeInt64($val->Data);
                }
                break;

            case ValueType::I32:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeInt64($val->Data);
                }
                break;

            case ValueType::I64:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeInt64($val->Data);
                }
                break;

            case ValueType::U:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeU64($val->Data);
                }
                break;

            case ValueType::U8:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeU64($val->Data);
                }
                break;

            case ValueType::U16:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeU64($val->Data);
                }
                break;

            case ValueType::U32:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeU64($val->Data);
                }
                break;

            case ValueType::U64:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeU64($val->Data);
                }
                break;

            case ValueType::F32:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeFloat($val->Text);
                }
                break;

            case ValueType::F64:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_FLOAT);
                } else {
                    $n = $this->encodeFloat($val->Text);
                }
                break;

            case ValueType::STR:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeString($val->Text);
                }
                break;

            case ValueType::EMAIL:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeString($val->Text);
                }
                break;

            case ValueType::UUID:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeBytes($val->Data);
                }
                break;

            case ValueType::DECIMAL:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeFloat($val->Text);
                }
                break;

            case ValueType::URL:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeString($val->Text);
                }
                break;

            case ValueType::IP:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $text = $val->Text;
                    switch ($tag->version) {
                        case 0:
                            $n = $this->encodeString($text);
                            break;
                        case 4:
                            $n = $this->encodeBytes($val->Data);
                            break;
                        case 6:
                            if (strlen($text) < 16) {
                                $n = $this->encodeString($text);
                            } else {
                                $n = $this->encodeBytes($val->Data);
                            }
                            break;
                        default:
                            throw new \Exception('unsupported IP version: ' . $tag->version);
                    }
                }
                break;

            case ValueType::BYTES:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_BYTES);
                } else {
                    $n = $this->encodeBytes($val->Data);
                }
                break;

            case ValueType::BIGINT:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeBigInt($val->Text);
                }
                break;

            case ValueType::BOOL:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_BOOL);
                } else {
                    $n = $this->encodeBool($val->Data);
                }
                break;

            case ValueType::ENUM:
                if ($tag->isNull) {
                    $n = $this->encodeSimple(SimpleValue::NULL_STRING);
                } else {
                    $n = $this->encodeInt64($val->Data);
                }
                break;

            default:
                throw new \Exception('type error: unsupported type: ' . $val->getTag()->type->name . ', value: ' . var_export($val->Data, true));
        }

        $n1 = $this->encodeComment(array_slice($this->buf, $this->offset - $n, $n), $tag);
        if ($n1 === 0) {
            return $n;
        }
        return $n1;
    }

    public function encodeSimple(int $value): int
    {
        $sign = Prefix::SIMPLE | $value;
        return $this->writeByte($sign);
    }

    public function encodeBool(bool $v): int
    {
        $value = $v ? SimpleValue::TRUE : SimpleValue::FALSE;
        return $this->encodeSimple($value);
    }

    public function encodeU64(int $uv): int
    {
        if ($uv < 0) {
            throw new \Exception('encodeU64: negative value');
        }
        return $this->encodeInt(Prefix::POSITIVE_INT, $uv);
    }

    public function encodeInt64(int $v): int
    {
        if ($v >= 0) {
            return $this->encodeInt(Prefix::POSITIVE_INT, $v);
        }

        $sign = Prefix::NEGATIVE_INT;
        if ($v === PHP_INT_MIN) {
            // -9223372036854775808 → 9223372036854775808 unsigned
            $uv = 9223372036854775808;
        } else {
            $uv = -$v;
        }
        return $this->encodeInt($sign, $uv);
    }

    private function encodeInt(int $sign, int $uv): int
    {
        if ($uv < WireConstants::INT_LEN_1) {
            $sign |= $uv;
            return $this->writeByte($sign);
        } elseif ($uv <= WireConstants::MAX_1) {
            $sign |= WireConstants::INT_LEN_1;
            return $this->writeByte($sign, $uv);
        } elseif ($uv <= WireConstants::MAX_2) {
            $sign |= WireConstants::INT_LEN_2;
            return $this->writeByte($sign, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_3) {
            $sign |= WireConstants::INT_LEN_3;
            return $this->writeByte($sign, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_4) {
            $sign |= WireConstants::INT_LEN_4;
            return $this->writeByte($sign, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_5) {
            $sign |= WireConstants::INT_LEN_5;
            return $this->writeByte($sign, ($uv >> 32) & 0xFF, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_6) {
            $sign |= WireConstants::INT_LEN_6;
            return $this->writeByte($sign, ($uv >> 40) & 0xFF, ($uv >> 32) & 0xFF, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_7) {
            $sign |= WireConstants::INT_LEN_7;
            return $this->writeByte($sign, ($uv >> 48) & 0xFF, ($uv >> 40) & 0xFF, ($uv >> 32) & 0xFF, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_8) {
            $sign |= WireConstants::INT_LEN_8;
            return $this->writeByte($sign, ($uv >> 56) & 0xFF, ($uv >> 48) & 0xFF, ($uv >> 40) & 0xFF, ($uv >> 32) & 0xFF, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } else {
            throw new \Exception('invalid byte size');
        }
    }

    public function encodeBigInt(string $s): int
    {
        $bits = BigIntWireCodec::encodeSignedDecimal($s);
        $inner = array_merge([strlen($s)], $bits);
        return $this->encodeBytes($inner);
    }

    public function encodeDateTime(\DateTime $t): int
    {
        $v = $t->getTimestamp();
        return $this->encodeInt64($v);
    }

    public function encodeDate(\DateTime $t): int
    {
        $default = new \DateTime('1970-01-01 00:00:00', new \DateTimeZone('UTC'));
        $v1 = clone $t;
        $v1->setTimezone(new \DateTimeZone('UTC'));
        $v1->setTime(0, 0, 0);
        $diff = $v1->diff($default);
        $v = (int)$diff->format('%r%a');
        return $this->encodeInt64($v);
    }

    public function encodeTime(\DateTime $t): int
    {
        $v1 = clone $t;
        $v1->setTimezone(new \DateTimeZone('UTC'));
        $v = $v1->format('H') * 3600 + $v1->format('i') * 60 + $v1->format('s');
        return $this->encodeInt(Prefix::POSITIVE_INT, (int)$v);
    }

    public function encodeFloat(string $s): int
    {
        $p = FloatCodec::parseDecimalString($s);
        $sign = Prefix::FLOAT;
        if ($p['negative']) {
            $sign |= WireConstants::FLOAT_NEG_MASK;
        }

        $exponent = $p['exponent'];
        $mantissa = $p['mantissa'];

        if ($exponent === -1 && $mantissa <= 7) {
            $sign |= $mantissa;
            return $this->writeByte($sign);
        } elseif ($mantissa <= WireConstants::MAX_1) {
            $sign |= WireConstants::FLOAT_LEN_1;
            return $this->writeByte($sign, $exponent & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_2) {
            $sign |= WireConstants::FLOAT_LEN_2;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_3) {
            $sign |= WireConstants::FLOAT_LEN_3;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 16) & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_4) {
            $sign |= WireConstants::FLOAT_LEN_4;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 24) & 0xFF, ($mantissa >> 16) & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_5) {
            $sign |= WireConstants::FLOAT_LEN_5;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 32) & 0xFF, ($mantissa >> 24) & 0xFF, ($mantissa >> 16) & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_6) {
            $sign |= WireConstants::FLOAT_LEN_6;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 40) & 0xFF, ($mantissa >> 32) & 0xFF, ($mantissa >> 24) & 0xFF, ($mantissa >> 16) & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_7) {
            $sign |= WireConstants::FLOAT_LEN_7;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 48) & 0xFF, ($mantissa >> 40) & 0xFF, ($mantissa >> 32) & 0xFF, ($mantissa >> 24) & 0xFF, ($mantissa >> 16) & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } elseif ($mantissa <= WireConstants::MAX_8) {
            $sign |= WireConstants::FLOAT_LEN_8;
            return $this->writeByte($sign, $exponent & 0xFF, ($mantissa >> 56) & 0xFF, ($mantissa >> 48) & 0xFF, ($mantissa >> 40) & 0xFF, ($mantissa >> 32) & 0xFF, ($mantissa >> 24) & 0xFF, ($mantissa >> 16) & 0xFF, ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
        } else {
            throw new \Exception('unsupported mantissa length');
        }
    }

    public function encodeString(string $s): int
    {
        $utf = unpack('C*', $s);
        if ($utf === false) {
            $utf = [];
        }
        $length = count($utf);
        if ($length > WireConstants::MAX_2) {
            throw new \Exception('string too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $length);
        }

        $sign = Prefix::STRING;
        if ($length < WireConstants::STRING_LEN_1) {
            $sign |= $length;
            return $this->writeStringWithPrefix($s, $sign);
        } elseif ($length < WireConstants::MAX_1) {
            $sign |= WireConstants::STRING_LEN_1;
            return $this->writeStringWithPrefix($s, $sign, $length);
        } else {
            $sign |= WireConstants::STRING_LEN_2;
            return $this->writeStringWithPrefix($s, $sign, ($length >> 8) & 0xFF, $length & 0xFF);
        }
    }

    public function encodeBytes(array $bs): int
    {
        $length = count($bs);
        if ($length > WireConstants::MAX_2) {
            throw new \Exception('bytes too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $length);
        }

        $sign = Prefix::BYTES;
        if ($length < WireConstants::BYTES_LEN_1) {
            $sign |= $length;
            return $this->writeBytesWithPrefix($bs, $sign);
        } elseif ($length < WireConstants::MAX_1) {
            $sign |= WireConstants::BYTES_LEN_1;
            return $this->writeBytesWithPrefix($bs, $sign, $length);
        } else {
            $sign |= WireConstants::BYTES_LEN_2;
            return $this->writeBytesWithPrefix($bs, $sign, ($length >> 8) & 0xFF, $length & 0xFF);
        }
    }

    public function encodeTag(array $payload, array $tag): int
    {
        if (empty($tag)) {
            return 0;
        }

        $length = count($payload) + count($tag);
        if ($length > WireConstants::MAX_2) {
            throw new \Exception('tag+payload too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $length);
        }

        $sign = Prefix::TAG;

        if ($length < WireConstants::TAG_LEN_1) {
            $sign |= $length;
            $ns = $this->writeBytesWithPrefix($tag, $sign);
            $nb = $this->writeBytes($payload);
        } elseif ($length < WireConstants::MAX_1) {
            $sign |= WireConstants::TAG_LEN_1;
            $ns = $this->writeBytesWithPrefix($tag, $sign, $length);
            $nb = $this->writeBytes($payload);
        } else {
            $sign |= WireConstants::TAG_LEN_2;
            $ns = $this->writeBytesWithPrefix($tag, $sign, ($length >> 8) & 0xFF, $length & 0xFF);
            $nb = $this->writeBytes($payload);
        }

        return $nb + $ns;
    }

    public function encodeT(array $bs): int
    {
        $length = count($bs);
        if ($length === 0) {
            return 0;
        }

        if ($length > WireConstants::MAX_2) {
            throw new \Exception('tag too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $length);
        }

        if ($length < 254) {
            return $this->writeBytesWithPrefix($bs, $length);
        }

        if ($length < 257) {
            return $this->writeBytesWithPrefix($bs, 254, $length);
        }

        return $this->writeBytesWithPrefix($bs, 255, ($length >> 8) & 0xFF, $length & 0xFF);
    }

    public function encodeObject(array $bs): int
    {
        $length = count($bs);
        if ($length > WireConstants::MAX_2) {
            throw new \Exception('object too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $length);
        }

        $sign = Prefix::CONTAINER | WireConstants::CONTAINER_MAP;
        if ($length < WireConstants::CONTAINER_LEN_1) {
            $sign |= $length;
            return $this->writeBytesWithPrefix($bs, $sign);
        } elseif ($length < WireConstants::MAX_1) {
            $sign |= WireConstants::CONTAINER_LEN_1;
            return $this->writeBytesWithPrefix($bs, $sign, $length);
        } else {
            $sign |= WireConstants::CONTAINER_LEN_2;
            return $this->writeBytesWithPrefix($bs, $sign, ($length >> 8) & 0xFF, $length & 0xFF);
        }
    }

    public function encodeArray(array $bs): int
    {
        $length = count($bs);
        if ($length > WireConstants::MAX_2) {
            throw new \Exception('array too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $length);
        }

        $sign = Prefix::CONTAINER | WireConstants::CONTAINER_ARRAY;
        if ($length < WireConstants::CONTAINER_LEN_1) {
            $sign |= $length;
            return $this->writeBytesWithPrefix($bs, $sign);
        } elseif ($length < WireConstants::MAX_1) {
            $sign |= WireConstants::CONTAINER_LEN_1;
            return $this->writeBytesWithPrefix($bs, $sign, $length);
        } else {
            $sign |= WireConstants::CONTAINER_LEN_2;
            return $this->writeBytesWithPrefix($bs, $sign, ($length >> 8) & 0xFF, $length & 0xFF);
        }
    }

    private function ensureCapacity(int $required): void
    {
        if ($required > $this->maxCap) {
            throw new \Exception('maximum size exceeded');
        }

        $currentCap = count($this->buf);
        if ($required > $currentCap) {
            $newCap = $currentCap * 2;
            if ($newCap > $this->maxCap || $newCap < $required) {
                $newCap = $required;
            }
            $newBuf = array_fill(0, $newCap, 0);
            for ($i = 0; $i < $this->offset; $i++) {
                $newBuf[$i] = $this->buf[$i];
            }
            $this->buf = $newBuf;
        }
    }

    public function writeByte(int ...$bs): int
    {
        $l = count($bs);
        $required = $this->offset + $l;
        $this->ensureCapacity($required);

        for ($i = 0; $i < $l; $i++) {
            $this->buf[$this->offset + $i] = $bs[$i] & 0xFF;
        }
        $this->offset += $l;
        return $l;
    }

    public function writeBytes(array $bs): int
    {
        $l = count($bs);
        if ($l === 0) {
            return 0;
        }

        $required = $this->offset + $l;
        $this->ensureCapacity($required);

        for ($i = 0; $i < $l; $i++) {
            $this->buf[$this->offset + $i] = $bs[$i] & 0xFF;
        }
        $this->offset += $l;
        return $l;
    }

    public function writeString(string $s): int
    {
        $l = strlen($s);
        $required = $this->offset + $l;
        $this->ensureCapacity($required);

        for ($i = 0; $i < $l; $i++) {
            $this->buf[$this->offset + $i] = ord($s[$i]);
        }
        $this->offset += $l;
        return $l;
    }

    public function writeBytesWithPrefix(array $bs, int ...$prefix): int
    {
        $lp = count($prefix);
        $l = $lp + count($bs);
        $required = $this->offset + $l;
        $this->ensureCapacity($required);

        for ($i = 0; $i < $lp; $i++) {
            $this->buf[$this->offset + $i] = $prefix[$i] & 0xFF;
        }
        for ($i = 0; $i < count($bs); $i++) {
            $this->buf[$this->offset + $lp + $i] = $bs[$i] & 0xFF;
        }
        $this->offset += $l;
        return $l;
    }

    public function writeStringWithPrefix(string $s, int ...$prefix): int
    {
        $lp = count($prefix);
        $l = $lp + strlen($s);
        $required = $this->offset + $l;
        $this->ensureCapacity($required);

        for ($i = 0; $i < $lp; $i++) {
            $this->buf[$this->offset + $i] = $prefix[$i] & 0xFF;
        }
        for ($i = 0; $i < strlen($s); $i++) {
            $this->buf[$this->offset + $lp + $i] = ord($s[$i]);
        }
        $this->offset += $l;
        return $l;
    }

    public function getEncodedBytes(int $written): array
    {
        return array_slice($this->buf, $this->offset - $written, $written);
    }
}
