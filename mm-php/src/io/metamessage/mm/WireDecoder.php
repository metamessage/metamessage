<?php

namespace io\metamessage\mm;

class WireDecoder {
    private array $data;
    private int $offset;

    public function __construct(array $data) {
        $this->data = $data;
        $this->offset = 0;
    }

    public function decode(): MmTree {
        $d = $this->decodeNode(null);
        if ($this->offset !== count($this->data)) {
            throw new MmDecodeException('Trailing bytes at ' . $this->offset . ' len ' . count($this->data));
        }
        return $d['node'];
    }

    private function decodeNode(?MmTag $inherited): array {
        $start = $this->offset;
        if ($this->offset >= count($this->data)) {
            throw new MmDecodeException('EOF');
        }
        $b = $this->data[$this->offset++];
        $prefix = $b & Prefix::PREFIX_MASK;
        switch ($prefix) {
            case Prefix::TAG:
                return $this->decodeTagged($b, $inherited, $start);
            case Prefix::SIMPLE:
                return ['node' => $this->decodeSimple($b, $inherited), 'consumed' => $this->offset - $start];
            case Prefix::POSITIVE_INT:
                return $this->decodePositiveInt($b, $inherited, $start);
            case Prefix::NEGATIVE_INT:
                return $this->decodeNegativeInt($b, $inherited, $start);
            case Prefix::FLOAT:
                return $this->decodeFloat($b, $inherited, $start);
            case Prefix::STRING:
                return $this->decodeString($b, $inherited, $start);
            case Prefix::BYTES:
                return $this->decodeBytes($b, $inherited, $start);
            case Prefix::CONTAINER:
                return $this->decodeContainer($b, $inherited, $start);
            default:
                throw new MmDecodeException('Invalid prefix');
        }
    }

    private function decodeTagged(int $firstByte, ?MmTag $inherited, int $start): array {
        $tl = $this->tagOuterLen($firstByte);
        $l1 = $tl[0];
        $l2 = $tl[1];
        if ($l1 === 1) {
            $l2 = $this->data[$this->offset++];
        } elseif ($l1 === 2) {
            $l2 = ($this->data[$this->offset++] << 8) | $this->data[$this->offset++];
        }
        $innerStart = $this->offset;
        $innerEnd = $innerStart + $l2;
        if ($innerEnd > count($this->data)) {
            throw new MmDecodeException('Tag frame past EOF');
        }
        // 读取标签字段长度
        $innerFieldLen = 0;
        $tagDataStart = $this->offset;
        if ($tagDataStart < $innerEnd) {
            $tb = $this->data[$tagDataStart];
            if ($tb < 254) {
                $innerFieldLen = $tb;
                $tagDataStart += 1;
            } elseif ($tb === 254) {
                if ($tagDataStart + 1 < $innerEnd) {
                    $innerFieldLen = $this->data[$tagDataStart + 1];
                    $tagDataStart += 2;
                } else {
                    // 没有足够的字节来读取标签字段长度，将 innerFieldLen 设置为 0
                    $innerFieldLen = 0;
                }
            } else {
                if ($tagDataStart + 2 < $innerEnd) {
                    $innerFieldLen = ($this->data[$tagDataStart + 1] << 8) | $this->data[$tagDataStart + 2];
                    $tagDataStart += 3;
                } else {
                    // 没有足够的字节来读取标签字段长度，将 innerFieldLen 设置为 0
                    $innerFieldLen = 0;
                }
            }
        }
        // 确保 tagDataStart 不超过 innerEnd
        if ($tagDataStart > $innerEnd) {
            $tagDataStart = $innerEnd;
            $innerFieldLen = 0;
        }
        // 确保 innerFieldLen 不会导致 tagDataEnd 超过 innerEnd
        if ($tagDataStart + $innerFieldLen > $innerEnd) {
            $innerFieldLen = 0;
        }
        // 计算标签数据的结束位置
        $tagDataEnd = $tagDataStart + $innerFieldLen;
        if ($tagDataEnd > $innerEnd) {
            throw new MmDecodeException('Tag fields overflow: tagDataStart=' . $tagDataStart . ', innerFieldLen=' . $innerFieldLen . ', tagDataEnd=' . $tagDataEnd . ', innerEnd=' . $innerEnd);
        }
        // 更新偏移量到标签数据的起始位置
        $this->offset = $tagDataStart;
        $tag = MmTag::empty();
        $tagData = array_slice($this->data, $tagDataStart, $innerFieldLen);
        $tagOffset = 0;
        $parsedTag = TagFieldParser::parse($tagData, $tagOffset);
        if ($parsedTag) {
            $tag = $parsedTag;
        }
        // 将偏移量设置为标签数据的结束位置，准备解码 payload
        $this->offset = $tagDataEnd;
        // 解码 payload
        $node;
        if ($this->offset < $innerEnd) {
            if ($tag->isNull) {
                $synthetic = $this->nullScalarForTag($tag);
                $node = $synthetic ?? $this->decodeNode($tag)['node'];
            } else {
                $node = $this->decodeNode($tag)['node'];
            }
        } else {
            // 没有 payload 数据，使用 nullScalarForTag 方法获取默认值
            $synthetic = $this->nullScalarForTag($tag);
            $node = $synthetic ?? new MmTree\MmScalar('', '', $tag);
        }
        if ($this->offset !== $innerEnd) {
            throw new MmDecodeException('Tag payload size mismatch: at ' . $this->offset . ' expected ' . $innerEnd);
        }
        return ['node' => $node, 'consumed' => $this->offset - $start];
    }

    private function nullScalarForTag(MmTag $tag): ?MmTree {
        switch ($tag->type) {
            case ValueType::DATETIME:
                $z = $this->zoneForHours($tag->locationHours);
                $dt = new \DateTime(TimeUtil::EPOCH, new \DateTimeZone($z ?? 'UTC'));
                return new MmTree\MmScalar($dt, $dt->format('Y-m-d H:i:s'), $tag);
            case ValueType::DATE:
                $d = new \DateTime(TimeUtil::EPOCH, new \DateTimeZone('UTC'));
                return new MmTree\MmScalar($d, $d->format('Y-m-d'), $tag);
            case ValueType::TIME:
                $t = new \DateTime('1970-01-01 00:00:00', new \DateTimeZone('UTC'));
                return new MmTree\MmScalar($t, $t->format('H:i:s'), $tag);
            case ValueType::INT8:
            case ValueType::INT16:
            case ValueType::INT32:
            case ValueType::INT:
                return new MmTree\MmScalar(0, '0', $tag);
            case ValueType::INT64:
                return new MmTree\MmScalar(0, '0', $tag);
            case ValueType::UINT:
            case ValueType::UINT8:
            case ValueType::UINT16:
            case ValueType::UINT32:
                return new MmTree\MmScalar(0, '0', $tag);
            case ValueType::UINT64:
                return new MmTree\MmScalar(0, '0', $tag);
            case ValueType::FLOAT32:
                return new MmTree\MmScalar(0.0, '0.0', $tag);
            case ValueType::FLOAT64:
                return new MmTree\MmScalar(0.0, '0.0', $tag);
            case ValueType::EMAIL:
            case ValueType::UUID:
            case ValueType::DECIMAL:
            case ValueType::URL:
                return new MmTree\MmScalar('', '', $tag);
            case ValueType::BIGINT:
                return new MmTree\MmScalar('0', '0', $tag);
            case ValueType::IP:
                return new MmTree\MmScalar('', $this->ipNullText($tag->version), $tag);
            default:
                return null;
        }
    }

    private function ipNullText(int $version): string {
        return match($version) {
            4 => '0.0.0.0',
            6 => '::',
            default => ''
        };
    }

    private function zoneForHours(int $hours): ?string {
        return $hours === 0 ? null : TimeUtil::zoneFromHourOffset($hours);
    }

    private function tagOuterLen(int $firstByte): array {
        $l = $firstByte & WireConstants::TAG_LEN_MASK;
        if ($l < WireConstants::TAG_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::TAG_LEN_1) {
            return [1, 0];
        } else {
            return [2, 0];
        }
    }

    private function readInnerLen(int $b): int {
        $l = $b;
        if ($l < 254) {
            return $l;
        } elseif ($l === 254) {
            $len = $this->data[$this->offset];
            $this->offset++;
            return $len;
        } else {
            $len = ($this->data[$this->offset] << 8) | $this->data[$this->offset + 1];
            $this->offset += 2;
            return $len;
        }
    }

    private function readInnerLenWithoutOffset(int $b): int {
        $l = $b;
        if ($l < 254) {
            return $l;
        } elseif ($l === 254) {
            return $this->data[$this->offset];
        } else {
            return ($this->data[$this->offset] << 8) | $this->data[$this->offset + 1];
        }
    }

    private function decodeSimple(int $first, ?MmTag $inherited): MmTree {
        $tag = $inherited ? clone $inherited : MmTag::empty();
        $sv = $first & Prefix::SUFFIX_MASK;
        switch ($sv) {
            case SimpleValue::FALSE:
                $tag->type = ValueType::BOOL;
                return new MmTree\MmScalar(false, 'false', $tag);
            case SimpleValue::TRUE:
                $tag->type = ValueType::BOOL;
                return new MmTree\MmScalar(true, 'true', $tag);
            case SimpleValue::NULL_BOOL:
                return $this->nullBool($tag);
            case SimpleValue::NULL_INT:
                return $this->nullInt($tag);
            case SimpleValue::NULL_FLOAT:
                return $this->nullFloat($tag);
            case SimpleValue::NULL_STRING:
                return $this->nullString($tag);
            case SimpleValue::NULL_BYTES:
                return $this->nullBytes($tag);
            default:
                throw new MmDecodeException('Unsupported simple: ' . $sv);
        }
    }

    private function nullBool(MmTag $tag): MmTree {
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::BOOL;
        }
        if ($tag->type !== ValueType::BOOL) {
            throw new MmDecodeException('null_bool type mismatch');
        }
        $tag->isNull = true;
        return new MmTree\MmScalar(false, 'false', $tag);
    }

    private function nullInt(MmTag $tag): MmTree {
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::INT;
        }
        if ($tag->type !== ValueType::INT) {
            throw new MmDecodeException('null_int type mismatch');
        }
        $tag->isNull = true;
        return new MmTree\MmScalar(0, '0', $tag);
    }

    private function nullFloat(MmTag $tag): MmTree {
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::FLOAT64;
        }
        if ($tag->type !== ValueType::FLOAT32 && $tag->type !== ValueType::FLOAT64) {
            throw new MmDecodeException('null_float type mismatch');
        }
        $tag->isNull = true;
        return new MmTree\MmScalar(0.0, '0.0', $tag);
    }

    private function nullString(MmTag $tag): MmTree {
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::STRING;
        }
        if ($tag->type !== ValueType::STRING) {
            throw new MmDecodeException('null_string type mismatch');
        }
        $tag->isNull = true;
        return new MmTree\MmScalar('', '', $tag);
    }

    private function nullBytes(MmTag $tag): MmTree {
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::BYTES;
        }
        if ($tag->type !== ValueType::BYTES) {
            throw new MmDecodeException('null_bytes type mismatch');
        }
        $tag->isNull = true;
        return new MmTree\MmScalar([], '', $tag);
    }

    private function decodePositiveInt(int $first, ?MmTag $inherited, int $start): array {
        $v = $this->readUintBody($first);
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::INT;
        }
        return ['node' => $this->mapUintToTree($tag, $v), 'consumed' => $this->offset - $start];
    }

    private function decodeNegativeInt(int $first, ?MmTag $inherited, int $start): array {
        $v = $this->readUintBody($first);
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::INT;
        }
        return ['node' => $this->mapNegativeInt($tag, $v), 'consumed' => $this->offset - $start];
    }

    private function readUintBody(int $first): int {
        $l1 = $this->intLenExtraBytes($first);
        $low = $first & WireConstants::INT_LEN_MASK;
        if ($l1 === 0) {
            return $low;
        } else {
            $v = 0;
            for ($i = 0; $i < $l1; $i++) {
                $v = ($v << 8) | $this->data[$this->offset++];
            }
            return $v;
        }
    }

    private function intLenExtraBytes(int $first): int {
        $l = $first & WireConstants::INT_LEN_MASK;
        return $l < WireConstants::INT_LEN_1 ? 0 : $l - WireConstants::INT_LEN_1 + 1;
    }

    private function mapUintToTree(MmTag $tag, int $v): MmTree {
        switch ($tag->type) {
            case ValueType::INT:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::INT8:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::INT16:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::INT32:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::INT64:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::UINT:
            case ValueType::UINT8:
            case ValueType::UINT16:
            case ValueType::UINT32:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::UINT64:
                return new MmTree\MmScalar($v, (string)$v, $tag);
            case ValueType::DATETIME:
                return $this->decodeDateTime($tag, $v);
            case ValueType::DATE:
                return $this->decodeDate($tag, $v);
            case ValueType::TIME:
                return $this->decodeTime($tag, $v);
            case ValueType::ENUM:
                return $this->decodeEnum($tag, $v);
            default:
                throw new MmDecodeException('Unsupported int type: ' . $tag->type->name);
        }
    }

    private function mapNegativeInt(MmTag $tag, int $v): MmTree {
        $negV = -$v;
        switch ($tag->type) {
            case ValueType::INT:
            case ValueType::INT8:
            case ValueType::INT16:
            case ValueType::INT32:
            case ValueType::INT64:
                return new MmTree\MmScalar($negV, (string)$negV, $tag);
            default:
                throw new MmDecodeException('Unsupported neg int type: ' . $tag->type->name);
        }
    }

    private function decodeDateTime(MmTag $tag, int $v): MmTree {
        if ($tag->isNull) {
            return new MmTree\MmScalar(null, '', $tag);
        }
        $ins = new \DateTime('@' . $v, new \DateTimeZone('UTC'));
        return new MmTree\MmScalar($ins, $ins->format('Y-m-d H:i:s'), $tag);
    }

    private function decodeDate(MmTag $tag, int $v): MmTree {
        if ($tag->isNull) {
            return new MmTree\MmScalar(null, '', $tag);
        }
        $d = TimeUtil::dateFromDays($v);
        return new MmTree\MmScalar($d, $d->format('Y-m-d'), $tag);
    }

    private function decodeTime(MmTag $tag, int $v): MmTree {
        if ($tag->isNull) {
            return new MmTree\MmScalar(null, '', $tag);
        }
        if ($v > 86399) {
            throw new MmDecodeException('Time out of range');
        }
        $t = TimeUtil::timeFromSeconds($v);
        return new MmTree\MmScalar($t, $t->format('H:i:s'), $tag);
    }

    private function decodeEnum(MmTag $tag, int $v): MmTree {
        if ($tag->isNull) {
            return new MmTree\MmScalar(-1, '', $tag);
        }
        if (empty($tag->enumValues)) {
            throw new MmDecodeException('Enum without labels');
        }
        $parts = explode('|', $tag->enumValues);
        if ($v >= count($parts)) {
            throw new MmDecodeException('Enum index out of range');
        }
        $label = trim($parts[$v]);
        return new MmTree\MmScalar($v, $label, $tag);
    }

    private function decodeFloat(int $first, ?MmTag $inherited, int $start): array {
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::FLOAT64;
        }
        $p = $first & Prefix::PREFIX_MASK;
        $val;
        if ($p === Prefix::FLOAT && ($first & ~WireConstants::FLOAT_NEG_MASK) <= (Prefix::FLOAT | 7)) {
            $mantissa = $first & WireConstants::FLOAT_LEN_MASK;
            $val = $mantissa / 10.0;
            if (($first & WireConstants::FLOAT_NEG_MASK) !== 0) {
                $val = -$val;
            }
        } else {
            $exp = $this->data[$this->offset++];
            $l1 = $this->floatLenExtraBytes($first);
            $low = $first & WireConstants::FLOAT_LEN_MASK;
            $mantissa;
            if ($l1 === 0) {
                $mantissa = $low;
            } else {
                $mantissa = 0;
                for ($i = 0; $i < $l1; $i++) {
                    $mantissa = ($mantissa << 8) | $this->data[$this->offset++];
                }
            }
            $dec = FloatCodec::mantissaToDecimal($mantissa, $exp);
            $val = (float)$dec;
            if (($first & WireConstants::FLOAT_NEG_MASK) !== 0) {
                $val = -$val;
            }
        }
        $node;
        switch ($tag->type) {
            case ValueType::FLOAT32:
            case ValueType::FLOAT64:
            case ValueType::DECIMAL:
                $node = new MmTree\MmScalar($val, (string)$val, $tag);
                break;
            default:
                throw new MmDecodeException('Bad float tag ' . $tag->type->name);
        }
        return ['node' => $node, 'consumed' => $this->offset - $start];
    }

    private function floatLenExtraBytes(int $first): int {
        $l = $first & WireConstants::FLOAT_LEN_MASK;
        return $l < WireConstants::FLOAT_LEN_1 ? 0 : $l - WireConstants::FLOAT_LEN_1 + 1;
    }

    private function decodeString(int $first, ?MmTag $inherited, int $start): array {
        $sl = $this->stringLen($first);
        $l1 = $sl[0];
        $l2 = $sl[1];
        if ($l1 === 1) {
            $l2 = $this->data[$this->offset++];
        } elseif ($l1 === 2) {
            $l2 = ($this->data[$this->offset++] << 8) | $this->data[$this->offset++];
        }
        $bs = $l2 > 0 ? $this->readBytes($l2) : [];
        $text = pack('C*', ...$bs);
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::STRING;
        }
        $node;
        switch ($tag->type) {
            case ValueType::STRING:
            case ValueType::EMAIL:
            case ValueType::URL:
            case ValueType::IP:
                $node = new MmTree\MmScalar($text, $text, $tag);
                break;
            default:
                throw new MmDecodeException('Unsupported string type: ' . $tag->type->name);
        }
        return ['node' => $node, 'consumed' => $this->offset - $start];
    }

    private function stringLen(int $first): array {
        $l = $first & WireConstants::STRING_LEN_MASK;
        if ($l < WireConstants::STRING_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::STRING_LEN_1) {
            return [1, 0];
        } else {
            return [2, 0];
        }
    }

    private function decodeBytes(int $first, ?MmTag $inherited, int $start): array {
        $bl = $this->bytesLen($first);
        $l1 = $bl[0];
        $l2 = $bl[1];
        if ($l1 === 1) {
            $l2 = $this->data[$this->offset++];
        } elseif ($l1 === 2) {
            $l2 = ($this->data[$this->offset++] << 8) | $this->data[$this->offset++];
        }
        $bs = $l2 > 0 ? $this->readBytes($l2) : [];
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::BYTES;
        }
        $node;
        switch ($tag->type) {
            case ValueType::BYTES:
                $node = new MmTree\MmScalar($bs, '', $tag);
                break;
            case ValueType::BIGINT:
                $node = $this->bigintFromBytes($bs, $tag);
                break;
            case ValueType::UUID:
                if (count($bs) !== 16) {
                    throw new MmDecodeException('UUID length');
                }
                $uuid = $this->uuidFromBytes($bs);
                $node = new MmTree\MmScalar($uuid, $uuid, $tag);
                break;
            case ValueType::IP:
                $node = new MmTree\MmScalar($bs, '', $tag);
                break;
            default:
                throw new MmDecodeException('Unsupported bytes type: ' . $tag->type->name);
        }
        return ['node' => $node, 'consumed' => $this->offset - $start];
    }

    private function bytesLen(int $first): array {
        $l = $first & WireConstants::BYTES_LEN_MASK;
        if ($l < WireConstants::BYTES_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::BYTES_LEN_1) {
            return [1, 0];
        } else {
            return [2, 0];
        }
    }

    private function uuidFromBytes(array $bs): string {
        $hex = implode('', array_map(function($b) { return str_pad(dechex($b), 2, '0', STR_PAD_LEFT); }, $bs));
        return substr($hex, 0, 8) . '-' . substr($hex, 8, 4) . '-' . substr($hex, 12, 4) . '-' . substr($hex, 16, 4) . '-' . substr($hex, 20);
    }

    private function bigintFromBytes(array $bs, MmTag $tag): MmTree {
        if (empty($bs)) {
            return new MmTree\MmScalar('0', '0', $tag);
        }
        $n = $bs[0];
        $body = array_slice($bs, 1);
        $bits = $this->bigintBits($body);
        $neg = !empty($bits) && $bits[0] === 1;
        $digits = BigIntWireCodec::decodePositive($body, $n);
        $bi = $neg ? '-' . $digits : $digits;
        return new MmTree\MmScalar($bi, $bi, $tag);
    }

    private function bigintBits(array $data): array {
        $bits = [];
        foreach ($data as $bt) {
            for ($i = 7; $i >= 0; $i--) {
                $bits[] = ($bt >> $i) & 1;
            }
        }
        return $bits;
    }

    private function decodeContainer(int $first, ?MmTag $inherited, int $start): array {
        $isArray = ($first & WireConstants::CONTAINER_MASK) === WireConstants::CONTAINER_ARRAY;
        return $isArray ? $this->decodeArray($first, $inherited, $start) : $this->decodeObject($first, $inherited, $start);
    }

    private function decodeArray(int $first, ?MmTag $inherited, int $start): array {
        $cl = $this->containerLen($first);
        $l1 = $cl[0];
        $l2 = $cl[1];
        if ($l1 === 1) {
            $l2 = $this->data[$this->offset++];
        } elseif ($l1 === 2) {
            $l2 = ($this->data[$this->offset++] << 8) | $this->data[$this->offset++];
        }
        $bodyStart = $this->offset;
        $bodyEnd = $bodyStart + $l2;
        if ($bodyEnd > count($this->data)) {
            throw new MmDecodeException('Array past EOF');
        }
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = $tag->size > 0 ? ValueType::ARRAY : ValueType::SLICE;
        }
        $items = [];
        while ($this->offset < $bodyEnd) {
            $elemTag = MmTag::empty();
            $elemTag->inheritFromArrayParent($tag);
            $el = $this->decodeNode($elemTag);
            $items[] = $el['node'];
        }
        if ($this->offset !== $bodyEnd) {
            throw new MmDecodeException('Array body misaligned');
        }
        return ['node' => new MmTree\MmArray($tag, $items), 'consumed' => $this->offset - $start];
    }

    private function decodeObject(int $first, ?MmTag $inherited, int $start): array {
        $cl = $this->containerLen($first);
        $l1 = $cl[0];
        $l2 = $cl[1];
        if ($l1 === 1) {
            $l2 = $this->data[$this->offset++];
        } elseif ($l1 === 2) {
            $l2 = ($this->data[$this->offset++] << 8) | $this->data[$this->offset++];
        }
        $innerStart = $this->offset;
        $innerEnd = $innerStart + $l2;
        if ($innerEnd > count($this->data)) {
            throw new MmDecodeException('Object past EOF');
        }
        $tag = $inherited ? clone $inherited : MmTag::empty();
        if ($tag->type === ValueType::UNKNOWN) {
            $tag->type = ValueType::STRUCT;
        }
        $keyPrefixPos = $this->offset;
        $keyPrefix = $this->data[$this->offset++];
        $keysDec = $this->decodeArray($keyPrefix, null, $keyPrefixPos);
        $keys = $keysDec['node'];
        if (!($keys instanceof MmTree\MmArray)) {
            throw new MmDecodeException('Expected array for keys');
        }
        $fields = [];
        $i = 0;
        while ($this->offset < $innerEnd) {
            $elemTag = MmTag::empty();
            $elemTag->inheritFromArrayParent($tag);
            $val = $this->decodeNode($elemTag);
            $key = $keys->items[$i]->text;
            $fields[] = [$key, $val['node']];
            $i++;
        }
        if ($this->offset !== $innerEnd) {
            throw new MmDecodeException('Object body misaligned');
        }
        return ['node' => new MmTree\MmObject($tag, $fields), 'consumed' => $this->offset - $start];
    }

    private function containerLen(int $first): array {
        $l = $first & WireConstants::CONTAINER_LEN_MASK;
        if ($l < WireConstants::CONTAINER_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::CONTAINER_LEN_1) {
            return [1, 0];
        } else {
            return [2, 0];
        }
    }

    private function readBytes(int $n): array {
        if ($this->offset + $n > count($this->data)) {
            throw new MmDecodeException('EOF');
        }
        $r = array_slice($this->data, $this->offset, $n);
        $this->offset += $n;
        return $r;
    }
}
