<?php

namespace io\metamessage\mm;

class ReflectMmEncoder {
    public static function encode($root): array {
        $enc = new WireEncoder();
        self::encodeValue($enc, $root, self::rootTagForClass($root::class));
        return $enc->toByteArray();
    }

    private static function rootTagForClass(string $className): MmTag {
        $reflection = new \ReflectionClass($className);
        $ann = $reflection->getAttributes(MM::class)[0] ?? null;
        $t = $ann ? MmTag::fromAnnotation($ann->newInstance()) : MmTag::empty();
        if ($t->type === ValueType::UNKNOWN) {
            $t->type = ValueType::STRUCT;
        }
        if (empty($t->name)) {
            $shortName = $reflection->getShortName();
            $t->name = CamelToSnake::convert($shortName);
        }
        return $t;
    }

    private static function encodeValue(WireEncoder $enc, $v, MmTag $tag): void {
        $work = $tag->copy();
        if ($v === null) {
            if (!$work->nullable && !$work->isNull) {
                throw new \Exception('Null for non-nullable');
            }
            $work->isNull = true;
        }
        if ($work->type === ValueType::STRUCT) {
            if ($v === null) {
                throw new \Exception('Null struct');
            }
            self::encodeStruct($enc, $v, $work);
            return;
        }
        if ($work->type === ValueType::SLICE || $work->type === ValueType::ARRAY) {
            self::encodeList($enc, $v, $work);
            return;
        }
        if ($work->type === ValueType::UNKNOWN) {
            $work->type = self::inferTypeFromValue($v);
        }
        $payload = new WireEncoder();
        self::encodeScalarPayload($payload, $v, $work);
        $enc->encodeTaggedPayload($payload->toByteArray(), $work->toBytes());
    }

    private static function inferTypeFromValue($v): ValueType {
        if ($v === null) {
            return ValueType::STRING;
        }
        switch (gettype($v)) {
            case 'string':
                return ValueType::STRING;
            case 'integer':
                return ValueType::INT;
            case 'double':
                return ValueType::FLOAT64;
            case 'boolean':
                return ValueType::BOOL;
            case 'array':
                return ValueType::SLICE;
            case 'object':
                if ($v instanceof \DateTime || $v instanceof \DateTimeImmutable) {
                    return ValueType::DATETIME;
                }
                return ValueType::STRUCT;
            default:
                return ValueType::STRING;
        }
    }

    private static function encodeStruct(WireEncoder $enc, $o, MmTag $objTag): void {
        $keysPacked = new GrowableByteBuf();
        $valsPacked = new GrowableByteBuf();
        $tmp = new WireEncoder();
        $reflection = new \ReflectionObject($o);
        foreach ($reflection->getProperties() as $f) {
            if ($f->isStatic()) {
                continue;
            }
            $f->setAccessible(true);
            $ann = $f->getAttributes(MM::class)[0] ?? null;
            if ($ann) {
                $mm = $ann->newInstance();
                if ($mm->name === '-') {
                    continue;
                }
            }
            $ft = TypeInference::forProperty($f);
            if ($ann) {
                $mm = $ann->newInstance();
                $annTag = MmTag::fromAnnotation($mm);
                if ($annTag->type !== ValueType::UNKNOWN) {
                    $ft->type = $annTag->type;
                }
                if ($annTag->childType !== ValueType::UNKNOWN) {
                    $ft->childType = $annTag->childType;
                }
                if (!empty($annTag->name)) {
                    $ft->name = $annTag->name;
                }
                self::mergeAnnotations($ft, $annTag);
            }
            $key = self::fieldKey($f, $ft, $ann);
            $fv = $f->getValue($o);
            $tmp->reset();
            self::encodeValue($tmp, $fv, $ft);
            $valsPacked->writeAll($tmp->toByteArray());
            $tmp->reset();
            $tmp->encodeString($key);
            $keysPacked->writeAll($tmp->toByteArray());
        }
        $tmp->reset();
        $tmp->encodeArrayPayload($keysPacked->copyRange(0, $keysPacked->length()));
        $mapBody = new GrowableByteBuf();
        $mapBody->writeAll($tmp->toByteArray());
        $mapBody->writeAll($valsPacked->copyRange(0, $valsPacked->length()));
        $tmp->reset();
        $tmp->encodeObjectPayload($mapBody->copyRange(0, $mapBody->length()));
        $enc->encodeTaggedPayload($tmp->toByteArray(), $objTag->toBytes());
    }

    private static function mergeAnnotations(MmTag $dst, MmTag $src): void {
        if (!empty($src->desc)) {
            $dst->desc = $src->desc;
        }
        $dst->nullable |= $src->nullable;
        $dst->raw |= $src->raw;
        $dst->allowEmpty |= $src->allowEmpty;
        $dst->unique |= $src->unique;
        if (!empty($src->defaultValue)) {
            $dst->defaultValue = $src->defaultValue;
        }
        if (!empty($src->enumValues)) {
            $dst->enumValues = $src->enumValues;
            $dst->type = ValueType::ENUM;
        }
        $dst->locationHours = $src->locationHours;
        $dst->version = $src->version;
        if (!empty($src->mime)) {
            $dst->mime = $src->mime;
        }
        $dst->childDesc = $src->childDesc;
        if ($src->childType !== ValueType::UNKNOWN) {
            $dst->childType = $src->childType;
        }
        $dst->childNullable |= $src->childNullable;
        if (!empty($src->childEnum)) {
            $dst->childEnum = $src->childEnum;
            $dst->childType = ValueType::ENUM;
        }
    }

    private static function fieldKey(\ReflectionProperty $f, MmTag $ft, ?\ReflectionAttribute $ann): string {
        if ($ann) {
            $mm = $ann->newInstance();
            if (!empty($mm->name) && $mm->name !== '-') {
                return $mm->name;
            }
        }
        if (!empty($ft->name) && $ft->name !== '-') {
            return $ft->name;
        }
        return CamelToSnake::convert($f->getName());
    }

    private static function encodeScalarPayload(WireEncoder $w, $v, MmTag $tag): void {
        if ($tag->isNull) {
            switch ($tag->type) {
                case ValueType::INT:
                    $w->encodeSimple(SimpleValue::NULL_INT);
                    return;
                case ValueType::STRING:
                    $w->encodeSimple(SimpleValue::NULL_STRING);
                    return;
                case ValueType::BYTES:
                    $w->encodeSimple(SimpleValue::NULL_BYTES);
                    return;
                case ValueType::FLOAT32:
                case ValueType::FLOAT64:
                    $w->encodeSimple(SimpleValue::NULL_FLOAT);
                    return;
                case ValueType::BOOL:
                    $w->encodeSimple(SimpleValue::NULL_BOOL);
                    return;
                default:
                    $w->encodeSimple(SimpleValue::NULL_STRING);
                    return;
            }
        }
        switch ($tag->type) {
            case ValueType::BOOL:
                $w->encodeBool((bool)$v);
                break;
            case ValueType::INT:
            case ValueType::INT8:
            case ValueType::INT16:
            case ValueType::INT32:
            case ValueType::UINT:
            case ValueType::UINT16:
                $w->encodeInt64((int)$v);
                break;
            case ValueType::INT64:
            case ValueType::UINT32:
            case ValueType::UINT64:
                $w->encodeInt64((int)$v);
                break;
            case ValueType::FLOAT32:
                $w->encodeFloatString((string)$v);
                break;
            case ValueType::FLOAT64:
            case ValueType::DECIMAL:
                $w->encodeFloatString((string)$v);
                break;
            case ValueType::STRING:
            case ValueType::EMAIL:
            case ValueType::URL:
                $w->encodeString($v ?? '');
                break;
            case ValueType::BYTES:
                $w->encodeBytes($v ?? []);
                break;
            case ValueType::BIGINT:
                $w->encodeBigIntDecimal($v ?? '0');
                break;
            case ValueType::UUID:
                $w->encodeBytes(self::uuidBytes($v));
                break;
            case ValueType::DATETIME:
                $sec = TimeUtil::epochSeconds($v);
                $w->encodeInt64($sec);
                break;
            case ValueType::DATE:
                $days = TimeUtil::daysSinceEpochUtc($v);
                $w->encodeInt64($days);
                break;
            case ValueType::TIME:
                $sec = TimeUtil::secondsOfDay($v);
                $w->encodeInt64($sec);
                break;
            case ValueType::ENUM:
                $w->encodeInt64((int)$v);
                break;
            default:
                throw new \Exception('Unsupported scalar ' . $tag->type->name);
        }
    }

    private static function encodeList(WireEncoder $enc, $list, MmTag $tag): void {
        if ($list === null) {
            $nt = $tag->copy();
            $nt->isNull = true;
            $enc->encodeTaggedPayload([], $nt->toBytes());
            return;
        }
        $body = new GrowableByteBuf();
        $el = new WireEncoder();
        foreach ($list as $item) {
            $el->reset();
            $et = MmTag::empty();
            $et->inheritFromArrayParent($tag);
            if ($et->type === ValueType::UNKNOWN && $item !== null) {
                $et->type = TypeInference::valueTypeForComponent(gettype($item));
                if ($et->type === ValueType::UNKNOWN) {
                    $et->type = ValueType::STRUCT;
                }
            }
            self::encodeValue($el, $item, $et);
            $body->writeAll($el->toByteArray());
        }
        $el->reset();
        $el->encodeArrayPayload($body->copyRange(0, $body->length()));
        $enc->encodeTaggedPayload($el->toByteArray(), $tag->toBytes());
    }

    private static function uuidBytes(string $uuid): array {
        $uuid = str_replace('-', '', $uuid);
        $bytes = [];
        for ($i = 0; $i < 32; $i += 2) {
            $bytes[] = hexdec(substr($uuid, $i, 2));
        }
        return $bytes;
    }
}
