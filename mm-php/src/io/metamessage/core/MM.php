<?php

namespace io\metamessage\core;

use Attribute;
use io\metamessage\ir\ValueType;

#[Attribute(Attribute::TARGET_CLASS | Attribute::TARGET_PROPERTY | Attribute::IS_REPEATABLE)]
class MM
{
    public string $name = '';
    public bool $isNull = false;
    public bool $example = false;
    public string $desc = '';
    public ValueType $type;
    public bool $raw = false;
    public bool $nullable = false;
    public bool $allowEmpty = false;
    public bool $unique = false;
    public string $defaultValue = '';
    public string $min = '';
    public string $max = '';
    public int $size = 0;
    public string $enumValues = '';
    public string $pattern = '';
    public int $location = 0;
    public int $version = 0;
    public string $mime = '';

    public string $childDesc = '';
    public ValueType $childType;
    public bool $childRaw = false;
    public bool $childNullable = false;
    public bool $childAllowEmpty = false;
    public bool $childUnique = false;
    public string $childDefault = '';
    public string $childMin = '';
    public string $childMax = '';
    public int $childSize = 0;
    public string $childEnum = '';
    public string $childPattern = '';
    public int $childLocation = 0;
    public int $childVersion = 0;
    public string $childMime = '';

    public function __construct(
        string $name = '',
        bool $isNull = false,
        bool $example = false,
        string $desc = '',
        ValueType $type = ValueType::UNKNOWN,
        bool $raw = false,
        bool $nullable = false,
        bool $allowEmpty = false,
        bool $unique = false,
        string $defaultValue = '',
        string $min = '',
        string $max = '',
        int $size = 0,
        string $enumValues = '',
        string $pattern = '',
        int $location = 0,
        int $version = 0,
        string $mime = '',
        string $childDesc = '',
        ValueType $childType = ValueType::UNKNOWN,
        bool $childRaw = false,
        bool $childNullable = false,
        bool $childAllowEmpty = false,
        bool $childUnique = false,
        string $childDefault = '',
        string $childMin = '',
        string $childMax = '',
        int $childSize = 0,
        string $childEnum = '',
        string $childPattern = '',
        int $childLocation = 0,
        int $childVersion = 0,
        string $childMime = '',
    ) {
        $this->name = $name;
        $this->isNull = $isNull;
        $this->example = $example;
        $this->desc = $desc;
        $this->type = $type;
        $this->raw = $raw;
        $this->nullable = $nullable;
        $this->allowEmpty = $allowEmpty;
        $this->unique = $unique;
        $this->defaultValue = $defaultValue;
        $this->min = $min;
        $this->max = $max;
        $this->size = $size;
        $this->enumValues = $enumValues;
        $this->pattern = $pattern;
        $this->location = $location;
        $this->version = $version;
        $this->mime = $mime;
        $this->childDesc = $childDesc;
        $this->childType = $childType;
        $this->childRaw = $childRaw;
        $this->childNullable = $childNullable;
        $this->childAllowEmpty = $childAllowEmpty;
        $this->childUnique = $childUnique;
        $this->childDefault = $childDefault;
        $this->childMin = $childMin;
        $this->childMax = $childMax;
        $this->childSize = $childSize;
        $this->childEnum = $childEnum;
        $this->childPattern = $childPattern;
        $this->childLocation = $childLocation;
        $this->childVersion = $childVersion;
        $this->childMime = $childMime;
    }
}