<?php

namespace io\metamessage\ir;

class ValidationResult
{
    public bool $isValid;
    public array $errors;

    public function __construct()
    {
        $this->isValid = true;
        $this->errors = [];
    }

    public function addError(string $error): void
    {
        $this->isValid = false;
        $this->errors[] = $error;
    }
}

class MmValidator
{
    private const MAX_DESC_LENGTH = 65535;

    public function validate(mixed $value, Tag $tag): ValidationResult
    {
        $result = new ValidationResult();

        switch ($tag->type) {
            case ValueType::STR:
                $this->validateStr($value, $tag, $result);
                break;
            case ValueType::BYTES:
                $this->validateBytes($value, $tag, $result);
                break;
            case ValueType::BOOL:
                $this->validateBool($value, $tag, $result);
                break;
            case ValueType::I:
                $this->validateI($value, $tag, $result);
                break;
            case ValueType::I8:
                $this->validateI8($value, $tag, $result);
                break;
            case ValueType::I16:
                $this->validateI16($value, $tag, $result);
                break;
            case ValueType::I32:
                $this->validateI32($value, $tag, $result);
                break;
            case ValueType::I64:
                $this->validateI64($value, $tag, $result);
                break;
            case ValueType::U:
                $this->validateU($value, $tag, $result);
                break;
            case ValueType::U8:
                $this->validateU8($value, $tag, $result);
                break;
            case ValueType::U16:
                $this->validateU16($value, $tag, $result);
                break;
            case ValueType::U32:
                $this->validateU32($value, $tag, $result);
                break;
            case ValueType::U64:
                $this->validateU64($value, $tag, $result);
                break;
            case ValueType::F32:
                $this->validateF32($value, $tag, $result);
                break;
            case ValueType::F64:
                $this->validateF64($value, $tag, $result);
                break;
            case ValueType::BIGINT:
                $this->validateBigint($value, $tag, $result);
                break;
            case ValueType::DATETIME:
                $this->validateDatetime($value, $tag, $result);
                break;
            case ValueType::DATE:
                $this->validateDate($value, $tag, $result);
                break;
            case ValueType::TIME:
                $this->validateTime($value, $tag, $result);
                break;
            case ValueType::UUID:
                $this->validateUUID($value, $tag, $result);
                break;
            case ValueType::DECIMAL:
                $this->validateDecimal($value, $tag, $result);
                break;
            case ValueType::IP:
                $this->validateIP($value, $tag, $result);
                break;
            case ValueType::URL:
                $this->validateURL($value, $tag, $result);
                break;
            case ValueType::EMAIL:
                $this->validateEmail($value, $tag, $result);
                break;
            case ValueType::ENUM:
                $this->validateEnum($value, $tag, $result);
                break;
            case ValueType::IMAGE:
                $this->validateImage($value, $tag, $result);
                break;
            case ValueType::ARR:
                $this->validateArr($value, $tag, $result);
                break;
            case ValueType::VEC:
                $this->validateVec($value, $tag, $result);
                break;
            case ValueType::OBJ:
                $this->validateObj($tag, $result);
                break;
            case ValueType::MAP:
                $this->validateMap($tag, $result);
                break;
            default:
                break;
        }

        return $result;
    }

    private function checkDescAndLocation(Tag $tag, ValidationResult $result, string $typeName): bool
    {
        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
            return false;
        }

        if ($tag->locationHours !== 0) {
            $result->addError("type $typeName not support location UTC{$tag->locationHours}");
            return false;
        }

        return true;
    }

    private function validateStr(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type string not allow empty value \"$val\"");
            return;
        }

        if ($tag->pattern !== '') {
            $re = @preg_match($tag->pattern, '');
            if ($re === false) {
                $result->addError("pattern \"{$tag->pattern}\" compile error");
                return;
            }

            if (!preg_match($tag->pattern, $val)) {
                $result->addError("value \"$val\" does not match pattern {$tag->pattern}");
                return;
            }
        }

        $l = mb_strlen($val, 'UTF-8');

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                $result->addError("failed to parse t.Min as int");
                return;
            }
            if ($l < $mini) {
                $result->addError("string length $l is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                $result->addError("failed to parse t.Max as int");
                return;
            }
            if ($l > $maxi) {
                $result->addError("string length $l exceeds the maximum limit $maxi");
                return;
            }
        }

        if ($tag->size !== 0) {
            if ($l !== $tag->size) {
                $result->addError("string length $l != size {$tag->size}");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'string');
    }

    private function validateBytes(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a byte string');
            return;
        }

        $val = $value;
        $l = strlen($val);

        if ($l === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type []byte not allow empty value []byte{}');
            return;
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                $result->addError("failed to parse t.Min as int");
                return;
            }
            if ($l < $mini) {
                $result->addError("[]byte length $l is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                $result->addError("failed to parse t.Max as int");
                return;
            }
            if ($l > $maxi) {
                $result->addError("[]byte length $l exceeds the maximum limit $maxi");
                return;
            }
        }

        if ($tag->size !== 0) {
            if ($l !== $tag->size) {
                $result->addError("[]byte length $l != size {$tag->size}");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, '[]byte');
    }

    private function validateBool(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_bool($value)) {
            $result->addError('value must be a boolean');
            return;
        }

        if ($tag->allowEmpty) {
            $result->addError('type bool not support allow empty');
            return;
        }

        $this->checkDescAndLocation($tag, $result, 'bool');
    }

    private function validateI(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type int not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseInt($tag->min, 64);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as int");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseInt($tag->max, 64);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as int");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'int');
    }

    private function validateI8(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type int8 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseInt($tag->min, 8);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as int8");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseInt($tag->max, 8);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as int8");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'int8');
    }

    private function validateI16(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type int16 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseInt($tag->min, 16);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as int16");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseInt($tag->max, 16);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as int16");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'int16');
    }

    private function validateI32(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type int32 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseInt($tag->min, 32);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as int32");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseInt($tag->max, 32);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as int32");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'int32');
    }

    private function validateI64(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type int64 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseInt($tag->min, 64);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as int64");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseInt($tag->max, 64);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as int64");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'int64');
    }

    private function validateU(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type uint not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseUint($tag->min, 64);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as uint");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseUint($tag->max, 64);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as uint");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'uint');
    }

    private function validateU8(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type uint8 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseUint($tag->min, 8);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as uint8");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseUint($tag->max, 8);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as uint8");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'uint8');
    }

    private function validateU16(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type uint16 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseUint($tag->min, 16);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as uint16");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseUint($tag->max, 16);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as uint16");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'uint16');
    }

    private function validateU32(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type uint32 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseUint($tag->min, 32);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as uint32");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseUint($tag->max, 32);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as uint32");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'uint32');
    }

    private function validateU64(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_int($value)) {
            $result->addError('value must be an integer');
            return;
        }

        $val = $value;

        if ($val === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("type uint64 not allow empty value $val");
            return;
        }

        if ($tag->min !== '') {
            $mini = $this->parseUint($tag->min, 64);
            if ($mini === null) {
                $result->addError("failed to parse t.Min as uint64");
                return;
            }
            if ($val < $mini) {
                $result->addError("value $val is less than the minimum limit $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = $this->parseUint($tag->max, 64);
            if ($maxi === null) {
                $result->addError("failed to parse t.Max as uint64");
                return;
            }
            if ($val > $maxi) {
                $result->addError("value $val exceeds the maximum limit $maxi");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'uint64');
    }

    private function validateF32(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_float($value) && !is_int($value)) {
            $result->addError('value must be a float');
            return;
        }

        $val = (float) $value;

        if ($val === 0.0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type float32 not allow empty value 0.0');
            return;
        }

        if ($tag->min !== '') {
            $mini = (float) $tag->min;
            if ((string) $mini !== $tag->min && $tag->min !== '0' && $tag->min !== '0.0') {
                $result->addError("failed to parse t.Min as float32");
                return;
            }
            if ($val < $mini) {
                $result->addError(sprintf('%f < min %f', $val, $mini));
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = (float) $tag->max;
            if ((string) $maxi !== $tag->max && $tag->max !== '0' && $tag->max !== '0.0') {
                $result->addError("failed to parse t.Max as float32");
                return;
            }
            if ($val > $maxi) {
                $result->addError(sprintf('%f > max %f', $val, $maxi));
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'float32');
    }

    private function validateF64(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_float($value) && !is_int($value)) {
            $result->addError('value must be a float');
            return;
        }

        $val = (float) $value;

        if ($val === 0.0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type float64 not allow empty value 0.0');
            return;
        }

        if ($tag->min !== '') {
            $mini = (float) $tag->min;
            if ((string) $mini !== $tag->min && $tag->min !== '0' && $tag->min !== '0.0') {
                $result->addError("failed to parse t.Min as float64");
                return;
            }
            if ($val < $mini) {
                $result->addError(sprintf('%f < min %f', $val, $mini));
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = (float) $tag->max;
            if ((string) $maxi !== $tag->max && $tag->max !== '0' && $tag->max !== '0.0') {
                $result->addError("failed to parse t.Max as float64");
                return;
            }
            if ($val > $maxi) {
                $result->addError(sprintf('%f > max %f', $val, $maxi));
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'float64');
    }

    private function validateBigint(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        $val = $value;
        if (is_int($val)) {
            $val = (string) $val;
        }
        if (!is_string($val)) {
            $result->addError('value must be a string representing a big integer');
            return;
        }

        if (!preg_match('/^-?\d+$/', $val)) {
            $result->addError('value must be a valid big integer string');
            return;
        }

        if ($this->bigIntSign($val) === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type big.Int not allow empty value 0');
            return;
        }

        if ($tag->min !== '') {
            if (!preg_match('/^-?\d+$/', $tag->min)) {
                $result->addError("invalid min \"{$tag->min}\" for big.Int");
                return;
            }
            if ($this->bigIntCmp($val, $tag->min) === -1) {
                $result->addError("big.Int length $val < min {$tag->min}");
                return;
            }
        }

        if ($tag->max !== '') {
            if (!preg_match('/^-?\d+$/', $tag->max)) {
                $result->addError("invalid max \"{$tag->max}\" for big.Int");
                return;
            }
            if ($this->bigIntCmp($val, $tag->max) === 1) {
                $result->addError("big.Int length $val > max {$tag->max}");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'big.Int');
    }

    private function validateDatetime(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!$value instanceof \DateTime && !is_int($value) && !is_string($value)) {
            $result->addError('value must be a DateTime, timestamp, or string');
            return;
        }

        try {
            if (is_int($value)) {
                $dt = new \DateTime('@' . $value, new \DateTimeZone('UTC'));
            } elseif (is_string($value)) {
                $dt = new \DateTime($value, new \DateTimeZone('UTC'));
            } else {
                $dt = clone $value;
            }
        } catch (\Exception $e) {
            $result->addError('value is not a valid datetime');
            return;
        }

        $tz = new \DateTimeZone('UTC');
        if ($tag->locationHours !== 0) {
            $offsetSeconds = $tag->locationHours * 3600;
            $sign = $offsetSeconds >= 0 ? '+' : '-';
            $offsetSeconds = abs($offsetSeconds);
            $hours = intdiv($offsetSeconds, 3600);
            $minutes = intdiv($offsetSeconds % 3600, 60);
            $tzName = sprintf('%s%02d:%02d', $sign, $hours, $minutes);
            $tz = new \DateTimeZone($tzName);
        }

        $dt->setTimezone($tz);
        $dt->setTime((int) $dt->format('H'), (int) $dt->format('i'), (int) $dt->format('s'));
        $format = $dt->format('Y-m-d H:i:s');

        if ($dt->getTimestamp() === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("datetime type does not allow empty \"$format\". you can set allow_empty or child_allow_empty to allow it.");
            return;
        }

        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
        }
    }

    private function validateDate(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!$value instanceof \DateTime && !is_int($value) && !is_string($value)) {
            $result->addError('value must be a DateTime, timestamp, or string');
            return;
        }

        try {
            if (is_int($value)) {
                $dt = new \DateTime('@' . $value, new \DateTimeZone('UTC'));
            } elseif (is_string($value)) {
                $dt = new \DateTime($value, new \DateTimeZone('UTC'));
            } else {
                $dt = clone $value;
            }
        } catch (\Exception $e) {
            $result->addError('value is not a valid date');
            return;
        }

        $tz = new \DateTimeZone('UTC');
        if ($tag->locationHours !== 0) {
            $offsetSeconds = $tag->locationHours * 3600;
            $sign = $offsetSeconds >= 0 ? '+' : '-';
            $offsetSeconds = abs($offsetSeconds);
            $hours = intdiv($offsetSeconds, 3600);
            $minutes = intdiv($offsetSeconds % 3600, 60);
            $tzName = sprintf('%s%02d:%02d', $sign, $hours, $minutes);
            $tz = new \DateTimeZone($tzName);
        }

        $dt->setTimezone($tz);
        $dt->setTime((int) $dt->format('H'), (int) $dt->format('i'), (int) $dt->format('s'));
        $format = $dt->format('Y-m-d');

        if ($dt->getTimestamp() === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("date type does not allow empty \"$format\". you can set allow_empty or child_allow_empty to allow it.");
            return;
        }

        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
        }
    }

    private function validateTime(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!$value instanceof \DateTime && !is_int($value) && !is_string($value)) {
            $result->addError('value must be a DateTime, timestamp, or string');
            return;
        }

        try {
            if (is_int($value)) {
                $dt = new \DateTime('@' . $value, new \DateTimeZone('UTC'));
            } elseif (is_string($value)) {
                $dt = new \DateTime($value, new \DateTimeZone('UTC'));
            } else {
                $dt = clone $value;
            }
        } catch (\Exception $e) {
            $result->addError('value is not a valid time');
            return;
        }

        $tz = new \DateTimeZone('UTC');
        if ($tag->locationHours !== 0) {
            $offsetSeconds = $tag->locationHours * 3600;
            $sign = $offsetSeconds >= 0 ? '+' : '-';
            $offsetSeconds = abs($offsetSeconds);
            $hours = intdiv($offsetSeconds, 3600);
            $minutes = intdiv($offsetSeconds % 3600, 60);
            $tzName = sprintf('%s%02d:%02d', $sign, $hours, $minutes);
            $tz = new \DateTimeZone($tzName);
        }

        $dt->setTimezone($tz);
        $dt->setTime((int) $dt->format('H'), (int) $dt->format('i'), (int) $dt->format('s'));
        $format = $dt->format('H:i:s');

        if ($dt->getTimestamp() === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError("time type does not allow empty \"$format\". you can set allow_empty or child_allow_empty to allow it.");
            return;
        }

        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
        }
    }

    private function validateUUID(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type uuid not allow empty value ""');
            return;
        }

        $uuidPattern = '/^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$/';
        if (!preg_match($uuidPattern, $val)) {
            $result->addError("value '$val' does not match UUID pattern");
            return;
        }

        $hex = str_replace('-', '', $val);
        if (strlen($hex) !== 32) {
            $result->addError("invalid uuid: $val");
            return;
        }

        $uuidBytes = [];
        for ($i = 0; $i < 32; $i += 2) {
            $uuidBytes[] = hexdec(substr($hex, $i, 2));
        }

        if ($tag->version !== Tag::DEFAULT_VERSION && $tag->version !== (int) (($uuidBytes[6] >> 4) & 0x0F)) {
            $result->addError('invalid uuid version');
            return;
        }

        $this->checkDescAndLocation($tag, $result, 'uuid');
    }

    private function validateDecimal(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type decimal not allow empty value ""');
            return;
        }

        $decimalPattern = '/^-?\d+\.\d+$/';
        if (!preg_match($decimalPattern, $val)) {
            $result->addError("invalid decimal \"$val\", must be like \"0.0\"");
            return;
        }

        $this->checkDescAndLocation($tag, $result, 'decimal');
    }

    private function validateIP(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '' || $val === '<nil>') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type ip not allow empty value ""');
            return;
        }

        if ($tag->version === 4) {
            if (!filter_var($val, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4)) {
                $result->addError("invalid ipv4: $val");
                return;
            }
        }

        if ($tag->version === 6) {
            if (!filter_var($val, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6) || filter_var($val, FILTER_VALIDATE_IP, FILTER_FLAG_IPV4)) {
                $result->addError("invalid ipv6: $val");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'ip');
    }

    private function validateURL(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type url not allow empty value ""');
            return;
        }

        $parsed = parse_url($val);

        if ($parsed === false) {
            $result->addError("invalid url: $val");
            return;
        }

        $scheme = $parsed['scheme'] ?? '';
        if ($scheme !== 'http' && $scheme !== 'https') {
            $result->addError("invalid url: $val");
            return;
        }

        $host = $parsed['host'] ?? '';
        if ($host === '') {
            $result->addError("invalid url: $val");
            return;
        }

        $this->checkDescAndLocation($tag, $result, 'url');
    }

    private function validateEmail(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type email not allow empty value ""');
            return;
        }

        $emailPattern = '/^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/';
        if (!preg_match($emailPattern, $val)) {
            $result->addError("value '$val' does not match email pattern");
            return;
        }

        $this->checkDescAndLocation($tag, $result, 'email');
    }

    private function validateEnum(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }

        $val = $value;

        if ($val === '') {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type enum not allow empty value ""');
            return;
        }

        $enums = explode('|', $tag->enumValues);
        $idx = -1;
        foreach ($enums as $i => $s) {
            if (trim($s) === $val) {
                $idx = $i;
                break;
            }
        }

        if ($idx === -1) {
            $enumList = implode(', ', $enums);
            $result->addError("value '$val' not found in enum: [$enumList]");
            return;
        }

        $this->checkDescAndLocation($tag, $result, 'enum');
    }

    private function validateImage(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (!is_string($value)) {
            $result->addError('value must be a byte string');
            return;
        }

        $val = $value;
        $l = strlen($val);

        if ($l === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type image not allow empty value []byte{}');
            return;
        }

        if ($tag->min !== '') {
            $mini = (int) $tag->min;
            if ((string) $mini !== $tag->min) {
                $result->addError("failed to parse t.Min as int");
                return;
            }
            if ($l < $mini) {
                $result->addError("[]byte length $l < min $mini");
                return;
            }
        }

        if ($tag->max !== '') {
            $maxi = (int) $tag->max;
            if ((string) $maxi !== $tag->max) {
                $result->addError("failed to parse t.Max as int");
                return;
            }
            if ($l > $maxi) {
                $result->addError("[]byte length $l > max $maxi");
                return;
            }
        }

        if ($tag->size !== 0) {
            if ($l !== $tag->size) {
                $result->addError("[]byte length $l != size {$tag->size}");
                return;
            }
        }

        $this->checkDescAndLocation($tag, $result, 'image');
    }

    private function validateArr(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
            return;
        }

        if ($tag->locationHours !== 0) {
            $result->addError("type array not support location UTC{$tag->locationHours}");
            return;
        }

        if (!is_array($value)) {
            $result->addError('value must be an array');
            return;
        }

        $l = count($value);

        if ($l === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type array not allow empty');
            return;
        }

        if ($tag->size > 0) {
            if ($l > $tag->size) {
                $result->addError('type array over size');
                return;
            }
        }

        if ($tag->childUnique && $l > 0) {
            $seen = [];
            foreach ($value as $i => $item) {
                $key = is_object($item) ? spl_object_hash($item) : (is_array($item) ? serialize($item) : $item);
                if (isset($seen[$key])) {
                    $result->addError("array duplicate value found at index: $i");
                    return;
                }
                $seen[$key] = true;
            }
        }
    }

    private function validateVec(mixed $value, Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if ($tag->locationHours !== 0) {
            $result->addError("type vec not support location UTC{$tag->locationHours}");
            return;
        }

        if (!is_array($value)) {
            $result->addError('value must be an vec');
            return;
        }

        $l = count($value);

        if ($l === 0) {
            if ($tag->allowEmpty) {
                return;
            }
            $result->addError('type vec not allow empty');
            return;
        }

        if ($tag->childUnique && $l > 0) {
            $seen = [];
            foreach ($value as $i => $item) {
                $key = is_object($item) ? spl_object_hash($item) : (is_array($item) ? serialize($item) : $item);
                if (isset($seen[$key])) {
                    $result->addError("vec duplicate value found at index: $i");
                    return;
                }
                $seen[$key] = true;
            }
        }
    }

    private function validateObj(Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
            return;
        }

        if ($tag->locationHours !== 0) {
            $result->addError("type struct not support location UTC{$tag->locationHours}");
            return;
        }
    }

    private function validateMap(Tag $tag, ValidationResult $result): void
    {
        if ($tag->isNull) {
            return;
        }

        if (strlen($tag->desc) > self::MAX_DESC_LENGTH) {
            $result->addError('desc length exceeds 65535 bytes');
            return;
        }

        if ($tag->locationHours !== 0) {
            $result->addError("type map not support location UTC{$tag->locationHours}");
            return;
        }
    }

    private function parseInt(string $s, int $bitSize): ?int
    {
        if (!is_numeric($s)) {
            return null;
        }

        $val = (int) $s;

        if ((string) $val !== $s && $s !== '0') {
            $val2 = floatval($s);
            if ($val2 !== (float) $val) {
                return null;
            }
            $val = (int) $val2;
        }

        switch ($bitSize) {
            case 8:
                if ($val < -128 || $val > 127) {
                    return null;
                }
                break;
            case 16:
                if ($val < -32768 || $val > 32767) {
                    return null;
                }
                break;
            case 32:
                if ($val < -2147483648 || $val > 2147483647) {
                    return null;
                }
                break;
        }

        return $val;
    }

    private function parseUint(string $s, int $bitSize): ?int
    {
        if (!is_numeric($s)) {
            return null;
        }

        $val = (int) $s;
        if ($val < 0) {
            return null;
        }

        if ((string) $val !== $s && $s !== '0') {
            return null;
        }

        switch ($bitSize) {
            case 8:
                if ($val > 255) {
                    return null;
                }
                break;
            case 16:
                if ($val > 65535) {
                    return null;
                }
                break;
            case 32:
                if ($val > 4294967295) {
                    return null;
                }
                break;
        }

        return $val;
    }

    private function bigIntSign(string $s): int
    {
        if ($s === '0' || $s === '-0') {
            return 0;
        }
        return $s[0] === '-' ? -1 : 1;
    }

    private function bigIntCmp(string $a, string $b): int
    {
        $negA = ($a[0] === '-');
        $negB = ($b[0] === '-');

        if ($negA && !$negB) {
            return -1;
        }
        if (!$negA && $negB) {
            return 1;
        }

        $aAbs = ltrim($a, '-');
        $bAbs = ltrim($b, '-');

        $lenA = strlen($aAbs);
        $lenB = strlen($bAbs);

        if ($lenA !== $lenB) {
            if ($negA) {
                return $lenA > $lenB ? -1 : 1;
            }
            return $lenA > $lenB ? 1 : -1;
        }

        $cmp = strcmp($aAbs, $bAbs);
        if ($negA) {
            return -$cmp;
        }
        return $cmp <=> 0;
    }
}
