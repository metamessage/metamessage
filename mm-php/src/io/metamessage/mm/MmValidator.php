<?php

namespace io\metamessage\mm;

class ValidationResult {
    public bool $isValid;
    public array $errors;
    
    public function __construct() {
        $this->isValid = true;
        $this->errors = [];
    }
    
    public function addError(string $error): void {
        $this->isValid = false;
        $this->errors[] = $error;
    }
}

class MmValidator {
    public function validate(mixed $value, MmTag $tag): ValidationResult {
        $result = new ValidationResult();

        if ($value === null) {
            if (!$tag->nullable) {
                $result->addError('value is required');
            }
            return $result;
        }

        switch ($tag->type) {
            case ValueType::BOOL:
                $this->validateBool($value, $tag, $result);
                break;
            case ValueType::INT:
                $this->validateInt($value, $tag, $result);
                break;
            case ValueType::INT8:
                $this->validateInt8($value, $tag, $result);
                break;
            case ValueType::INT16:
                $this->validateInt16($value, $tag, $result);
                break;
            case ValueType::INT32:
                $this->validateInt32($value, $tag, $result);
                break;
            case ValueType::INT64:
                $this->validateInt64($value, $tag, $result);
                break;
            case ValueType::UINT:
                $this->validateUint($value, $tag, $result);
                break;
            case ValueType::UINT8:
                $this->validateUint8($value, $tag, $result);
                break;
            case ValueType::UINT16:
                $this->validateUint16($value, $tag, $result);
                break;
            case ValueType::UINT32:
                $this->validateUint32($value, $tag, $result);
                break;
            case ValueType::UINT64:
                $this->validateUint64($value, $tag, $result);
                break;
            case ValueType::BIGINT:
                $this->validateBigInt($value, $tag, $result);
                break;
            case ValueType::FLOAT32:
                $this->validateFloat32($value, $tag, $result);
                break;
            case ValueType::FLOAT64:
                $this->validateFloat64($value, $tag, $result);
                break;
            case ValueType::DECIMAL:
                $this->validateFloat($value, $tag, $result);
                break;
            case ValueType::STRING:
                $this->validateString($value, $tag, $result);
                break;
            case ValueType::EMAIL:
                $this->validateEmail($value, $tag, $result);
                break;
            case ValueType::URL:
                $this->validateUrl($value, $tag, $result);
                break;
            case ValueType::BYTES:
                $this->validateBytes($value, $tag, $result);
                break;
            case ValueType::UUID:
                $this->validateUuid($value, $tag, $result);
                break;
            case ValueType::DATETIME:
            case ValueType::DATE:
            case ValueType::TIME:
                $this->validateDateTime($value, $tag, $result);
                break;
            case ValueType::ENUM:
                $this->validateEnum($value, $tag, $result);
                break;
            case ValueType::ARRAY:
            case ValueType::SLICE:
                $this->validateArray($value, $tag, $result);
                break;
            case ValueType::STRUCT:
                $this->validateStruct($value, $tag, $result);
                break;
            default:
                break;
        }

        return $result;
    }

    private function validateBool(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_bool($value)) {
            $result->addError('value must be a boolean');
        }
    }

    private function validateInt(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
        }
    }

    private function validateInt8(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < -128 || $value > 127) {
            $result->addError('value out of range for int8');
        }
    }

    private function validateInt16(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < -32768 || $value > 32767) {
            $result->addError('value out of range for int16');
        }
    }

    private function validateInt32(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < -2147483648 || $value > 2147483647) {
            $result->addError('value out of range for int32');
        }
    }

    private function validateInt64(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
    }

    private function validateUint(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < 0 || $value > 4294967295) {
            $result->addError('value out of range for uint');
        }
    }

    private function validateUint8(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < 0 || $value > 255) {
            $result->addError('value out of range for uint8');
        }
    }

    private function validateUint16(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < 0 || $value > 65535) {
            $result->addError('value out of range for uint16');
        }
    }

    private function validateUint32(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
        if ($value < 0 || $value > 4294967295) {
            $result->addError('value out of range for uint32');
        }
    }

    private function validateUint64(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_int($value)) {
            $result->addError('value must be a number');
            return;
        }
    }

    private function validateBigInt(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (is_int($value)) {
            return;
        }
        if (!is_string($value)) {
            $result->addError('value must be a bigint string');
            return;
        }
        if (!extension_loaded('gmp')) {
            $result->addError('GMP extension is not loaded');
            return;
        }
        try {
            gmp_init($value);
        } catch (\Exception $e) {
            $result->addError('value must be a valid bigint string');
        }
    }

    private function validateFloat32(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_float($value) && !is_int($value)) {
            $result->addError('value must be a float');
            return;
        }
        $floatValue = (float)$value;
        if ($floatValue < -3.4028234663852886e+38 || $floatValue > 3.4028234663852886e+38) {
            $result->addError('value out of range for float32');
        }
    }

    private function validateFloat64(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_float($value) && !is_int($value)) {
            $result->addError('value must be a float');
            return;
        }
    }

    private function validateFloat(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_float($value) && !is_int($value)) {
            $result->addError('value must be a float');
        }
    }
    
    private function validateString(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
        }
    }
    
    private function validateEmail(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
            return;
        }
        
        if (!empty($value) && !filter_var($value, FILTER_VALIDATE_EMAIL)) {
            $result->addError('value is not a valid email');
        }
    }
    
    private function validateUrl(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
            return;
        }
        
        if (!empty($value) && !filter_var($value, FILTER_VALIDATE_URL)) {
            $result->addError('value is not a valid url');
        }
    }
    
    private function validateBytes(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value) && !is_array($value)) {
            $result->addError('value must be a byte array');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
        }
    }
    
    private function validateUuid(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
            return;
        }
        
        if (!empty($value) && !preg_match('/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i', $value)) {
            $result->addError('value is not a valid uuid');
        }
    }
    
    private function validateDateTime(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value) && !$value instanceof \DateTime) {
            $result->addError('value must be a datetime');
        }
    }
    
    private function validateEnum(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_string($value)) {
            $result->addError('value must be a string');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
            return;
        }
        
        if (!empty($value) && !empty($tag->enumValues)) {
            $enumValues = explode('|', $tag->enumValues);
            $enumValues = array_map('trim', $enumValues);
            if (!in_array($value, $enumValues)) {
                $result->addError('value is not in enum');
            }
        }
    }
    
    private function validateArray(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_array($value)) {
            $result->addError('value must be an array');
            return;
        }
        
        if (empty($value) && !$tag->allowEmpty) {
            $result->addError('value is empty');
        }
        
        // 验证数组元素
        foreach ($value as $index => $item) {
            $childTag = MmTag::empty();
            $childTag->inheritFromArrayParent($tag);
            $itemResult = $this->validate($item, $childTag);
            if (!$itemResult->isValid) {
                foreach ($itemResult->errors as $error) {
                    $result->addError("item $index: $error");
                }
            }
        }
    }
    
    private function validateStruct(mixed $value, MmTag $tag, ValidationResult $result): void {
        if (!is_array($value) && !is_object($value)) {
            $result->addError('value must be an object');
        }
    }
}

// 全局验证器实例
$validator = new MmValidator();
