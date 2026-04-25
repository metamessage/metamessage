<?php

namespace io\metamessage\jsonc;

enum JsoncValueType {
    case Unknown;
    case String;
    case Int;
    case Int8;
    case Int16;
    case Int32;
    case Int64;
    case Uint;
    case Uint8;
    case Uint16;
    case Uint32;
    case Uint64;
    case Float32;
    case Float64;
    case Bool;
    case Bytes;
    case BigInt;
    case DateTime;
    case Date;
    case Time;
    case UUID;
    case Decimal;
    case IP;
    case URL;
    case Email;
    case Enum;
    case Array;
    case Struct;
    case Slice;
    case Map;
    case Null;
    case Raw;
}

class JsoncTag {
    public function __construct(
        public string $name = "",
        public bool $isNull = false,
        public JsoncValueType $type = JsoncValueType::Unknown,
        public string $desc = "",
        public bool $raw = false,
        public bool $nullable = false,
        public bool $allowEmpty = false,
        public bool $unique = false,
        public string $defaultValue = "",
        public string $min = "",
        public string $max = "",
        public int $size = 0,
        public string $enum = "",
        public string $pattern = "",
        public string $location = "",
        public int $version = 0,
        public string $mime = "",
        public string $childDesc = "",
        public JsoncValueType $childType = JsoncValueType::Unknown,
        public bool $childRaw = false,
        public bool $childNullable = false,
        public bool $childAllowEmpty = false,
        public bool $childUnique = false,
        public string $childDefault = "",
        public string $childMin = "",
        public string $childMax = "",
        public int $childSize = 0,
        public string $childEnum = "",
        public string $childPattern = "",
        public string $childLocation = "",
        public int $childVersion = 0,
        public string $childMime = "",
        public bool $isInherit = false
    ) {}
}

abstract class JsoncNode {
    public ?JsoncTag $tag = null;
    public string $path = "";
}

class JsoncValue extends JsoncNode {
    public function __construct(
        public mixed $data = null,
        public string $text = "",
        ?JsoncTag $tag = null,
        string $path = ""
    ) {
        $this->tag = $tag;
        $this->path = $path;
    }
}

class JsoncField {
    public function __construct(
        public string $key,
        public ?JsoncNode $value
    ) {}
}

class JsoncObject extends JsoncNode {
    /** @var JsoncField[] */
    public array $fields = [];

    public function __construct(?JsoncTag $tag = null, string $path = "") {
        $this->tag = $tag;
        $this->path = $path;
    }
}

class JsoncArray extends JsoncNode {
    /** @var JsoncNode[] */
    public array $items = [];

    public function __construct(?JsoncTag $tag = null, string $path = "") {
        $this->tag = $tag;
        $this->path = $path;
    }
}
