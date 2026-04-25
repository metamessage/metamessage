<?php

namespace io\metamessage\mm;

use Attribute;

#[Attribute(Attribute::TARGET_PROPERTY | Attribute::TARGET_CLASS)]
class MM {
    public function __construct(
        public string $name = '',
        public ValueType $type = ValueType::UNKNOWN,
        public bool $isNull = false,
        public bool $example = false,
        public string $desc = '',
        public bool $raw = false,
        public bool $nullable = false,
        public bool $allowEmpty = false,
        public bool $unique = false,
        public string $defaultValue = '',
        public string $min = '',
        public string $max = '',
        public int $size = 0,
        public string $enumValues = '',
        public string $pattern = '',
        public int $location = 0,
        public int $version = 0,
        public string $mime = '',
        public string $childDesc = '',
        public ValueType $childType = ValueType::UNKNOWN,
        public bool $childRaw = false,
        public bool $childNullable = false,
        public bool $childAllowEmpty = false,
        public bool $childUnique = false,
        public string $childDefault = '',
        public string $childMin = '',
        public string $childMax = '',
        public int $childSize = 0,
        public string $childEnum = '',
        public string $childPattern = '',
        public int $childLocation = 0,
        public int $childVersion = 0,
        public string $childMime = ''
    ) {}
}
