<?php

namespace io\metamessage\mm;

class WireEncoder {
    private GrowableByteBuf $buf;

    public function __construct() {
        $this->buf = new GrowableByteBuf();
    }

    public function toByteArray(): array {
        return $this->buf->copyRange(0, $this->buf->length());
    }

    public function finishTakeLast(int $writtenFromEnd): array {
        $end = $this->buf->length();
        return $this->buf->copyRange($end - $writtenFromEnd, $end);
    }

    public function reset(): void {
        $this->buf->reset();
    }

    public function size(): int {
        return $this->buf->length();
    }

    public function encodeSimple(int $simpleValue): int {
        $start = $this->buf->length();
        $this->buf->write(Prefix::SIMPLE | $simpleValue);
        return $this->buf->length() - $start;
    }

    public function encodeBool(bool $v): int {
        return $this->encodeSimple($v ? SimpleValue::TRUE : SimpleValue::FALSE);
    }

    public function encodeInt64(int $v): int {
        if ($v >= 0) {
            return $this->encodeUintWithPrefix(Prefix::POSITIVE_INT, $v);
        } else {
            $uv = $v === PHP_INT_MIN ? PHP_INT_MAX + 1 : -$v;
            return $this->encodeUintWithPrefix(Prefix::NEGATIVE_INT, $uv);
        }
    }

    public function encodeUint64(int $uv): int {
        if ($uv < 0) {
            throw new \Exception('Expected unsigned');
        }
        return $this->encodeUintWithPrefix(Prefix::POSITIVE_INT, $uv);
    }

    private function encodeUintWithPrefix(int $prefix, int $uv): int {
        $start = $this->buf->length();
        if ($uv < WireConstants::INT_LEN_1) {
            $this->buf->write($prefix | $uv);
        } elseif ($uv <= WireConstants::MAX_1) {
            $this->buf->write($prefix | WireConstants::INT_LEN_1, $uv);
        } elseif ($uv <= WireConstants::MAX_2) {
            $this->buf->write($prefix | WireConstants::INT_LEN_2, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_3) {
            $this->buf->write($prefix | WireConstants::INT_LEN_3, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_4) {
            $this->buf->write($prefix | WireConstants::INT_LEN_4, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } else {
            throw new \Exception('Integer too large');
        }
        return $this->buf->length() - $start;
    }

    public function encodeFloatString(string $s): int {
        $p = FloatCodec::parseDecimalString($s);
        $start = $this->buf->length();
        $sign = Prefix::FLOAT;
        if ($p['negative']) {
            $sign |= WireConstants::FLOAT_NEG_MASK;
        }
        if ($p['exponent'] === -1 && $p['mantissa'] <= 7) {
            $this->buf->write($sign | $p['mantissa']);
        } else {
            $mantissa = $p['mantissa'];
            if ($mantissa <= WireConstants::MAX_1) {
                $this->buf->write($sign | WireConstants::FLOAT_LEN_1, $p['exponent'], $mantissa);
            } elseif ($mantissa <= WireConstants::MAX_2) {
                $this->buf->write($sign | WireConstants::FLOAT_LEN_2, $p['exponent'], ($mantissa >> 8) & 0xFF, $mantissa & 0xFF);
            } else {
                throw new \Exception('Float mantissa too large');
            }
        }
        return $this->buf->length() - $start;
    }

    public function encodeString(string $s): int {
        $utf = unpack('C*', $s);
        $length = count($utf);
        $start = $this->buf->length();
        $sign = Prefix::STRING;
        if ($length < WireConstants::STRING_LEN_1) {
            $this->buf->write($sign | $length, ...$utf);
        } elseif ($length < WireConstants::MAX_1) {
            $this->buf->write($sign | WireConstants::STRING_LEN_1, $length, ...$utf);
        } elseif ($length < WireConstants::MAX_2) {
            $this->buf->write($sign | WireConstants::STRING_LEN_2, ($length >> 8) & 0xFF, $length & 0xFF, ...$utf);
        } else {
            throw new \Exception('String too long');
        }
        return $this->buf->length() - $start;
    }

    public function encodeBytes(array $data): int {
        $length = count($data);
        $start = $this->buf->length();
        $sign = Prefix::BYTES;
        if ($length < WireConstants::BYTES_LEN_1) {
            $this->buf->write($sign | $length, ...$data);
        } elseif ($length < WireConstants::MAX_1) {
            $this->buf->write($sign | WireConstants::BYTES_LEN_1, $length, ...$data);
        } elseif ($length < WireConstants::MAX_2) {
            $this->buf->write($sign | WireConstants::BYTES_LEN_2, ($length >> 8) & 0xFF, $length & 0xFF, ...$data);
        } else {
            throw new \Exception('Bytes too long');
        }
        return $this->buf->length() - $start;
    }

    public function encodeArrayPayload(array $payload): int {
        return $this->encodeContainer($payload, Prefix::CONTAINER | WireConstants::CONTAINER_ARRAY);
    }

    public function encodeObjectPayload(array $payload): int {
        return $this->encodeContainer($payload, Prefix::CONTAINER | WireConstants::CONTAINER_MAP);
    }

    private function encodeContainer(array $payload, int $baseSign): int {
        $length = count($payload);
        $start = $this->buf->length();
        if ($length < WireConstants::CONTAINER_LEN_1) {
            $this->buf->write($baseSign | $length, ...$payload);
        } elseif ($length < WireConstants::MAX_1) {
            $this->buf->write($baseSign | WireConstants::CONTAINER_LEN_1, $length, ...$payload);
        } elseif ($length < WireConstants::MAX_2) {
            $this->buf->write($baseSign | WireConstants::CONTAINER_LEN_2, ($length >> 8) & 0xFF, $length & 0xFF, ...$payload);
        } else {
            throw new \Exception('Container payload too long');
        }
        return $this->buf->length() - $start;
    }

    public function encodeTagInner(array $tagBytes): int {
        if (empty($tagBytes)) {
            return 0;
        }
        if (count($tagBytes) > WireConstants::MAX_2) {
            throw new \Exception('Tag too long');
        }
        $start = $this->buf->length();
        $length = count($tagBytes);
        if ($length < 254) {
            $this->buf->write($length, ...$tagBytes);
        } elseif ($length < 257) {
            $this->buf->write(254, $length, ...$tagBytes);
        } else {
            $this->buf->write(255, ($length >> 8) & 0xFF, $length & 0xFF, ...$tagBytes);
        }
        return $this->buf->length() - $start;
    }

    public function encodeTaggedPayload(array $payload, array $rawTagFields): int {
        if (empty($rawTagFields)) {
            $this->buf->writeAll($payload);
            return count($payload);
        }
        $tEnc = new WireEncoder();
        $tEnc->encodeTagInner($rawTagFields);
        $tagEncoded = $tEnc->toByteArray();
        $length = count($tagEncoded) + count($payload);
        if ($length > WireConstants::MAX_2) {
            throw new \Exception('Tag+payload too long');
        }
        $start = $this->buf->length();
        $sign = Prefix::TAG;
        if ($length < WireConstants::TAG_LEN_1) {
            $this->buf->write($sign | $length, ...$tagEncoded, ...$payload);
        } elseif ($length < WireConstants::MAX_1) {
            $this->buf->write($sign | WireConstants::TAG_LEN_1, $length, ...$tagEncoded, ...$payload);
        } else {
            $this->buf->write($sign | WireConstants::TAG_LEN_2, ($length >> 8) & 0xFF, $length & 0xFF, ...$tagEncoded, ...$payload);
        }
        return $this->buf->length() - $start;
    }

    public function encodeBigIntDecimal(string $s): int {
        $bits = BigIntWireCodec::encodeSignedDecimal($s);
        $inner = array_merge([strlen($s)], $bits);
        return $this->encodeBytes($inner);
    }

    public function sliceFrom(int $start): int {
        return $this->buf->length() - $start;
    }

    public function copyLast(int $n): array {
        $end = $this->buf->length();
        return $this->buf->copyRange($end - $n, $end);
    }
}
