<?php

namespace io\metamessage\mm;

class MmTag {
    public string $name = '';
    public bool $isNull = false;
    public bool $example = false;
    public string $desc = '';
    public ValueType $type = ValueType::UNKNOWN;
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
    public int $locationHours = 0;
    public int $version = self::DEFAULT_VERSION;
    public string $mime = '';
    public string $childDesc = '';
    public ValueType $childType = ValueType::UNKNOWN;
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
    public int $childLocationHours = 0;
    public int $childVersion = self::DEFAULT_VERSION;
    public string $childMime = '';
    public bool $isInherit = false;

    public const DEFAULT_VERSION = 0;

    public function copy(): self {
        $copy = new self();
        $copy->name = $this->name;
        $copy->isNull = $this->isNull;
        $copy->example = $this->example;
        $copy->desc = $this->desc;
        $copy->type = $this->type;
        $copy->raw = $this->raw;
        $copy->nullable = $this->nullable;
        $copy->allowEmpty = $this->allowEmpty;
        $copy->unique = $this->unique;
        $copy->defaultValue = $this->defaultValue;
        $copy->min = $this->min;
        $copy->max = $this->max;
        $copy->size = $this->size;
        $copy->enumValues = $this->enumValues;
        $copy->pattern = $this->pattern;
        $copy->locationHours = $this->locationHours;
        $copy->version = $this->version;
        $copy->mime = $this->mime;
        $copy->childDesc = $this->childDesc;
        $copy->childType = $this->childType;
        $copy->childRaw = $this->childRaw;
        $copy->childNullable = $this->childNullable;
        $copy->childAllowEmpty = $this->childAllowEmpty;
        $copy->childUnique = $this->childUnique;
        $copy->childDefault = $this->childDefault;
        $copy->childMin = $this->childMin;
        $copy->childMax = $this->childMax;
        $copy->childSize = $this->childSize;
        $copy->childEnum = $this->childEnum;
        $copy->childPattern = $this->childPattern;
        $copy->childLocationHours = $this->childLocationHours;
        $copy->childVersion = $this->childVersion;
        $copy->childMime = $this->childMime;
        $copy->isInherit = $this->isInherit;
        return $copy;
    }

    public function inheritFromArrayParent(?MmTag $parent): void {
        if ($parent === null) {
            return;
        }
        $this->desc = $parent->childDesc;
        $this->type = $parent->childType;
        $this->raw = $parent->childRaw;
        $this->nullable = $parent->childNullable;
        $this->allowEmpty = $parent->childAllowEmpty;
        $this->unique = $parent->childUnique;
        $this->defaultValue = $parent->childDefault;
        $this->min = $parent->childMin;
        $this->max = $parent->childMax;
        $this->size = $parent->childSize;
        $this->enumValues = $parent->childEnum;
        $this->pattern = $parent->childPattern;
        $this->locationHours = $parent->childLocationHours;
        $this->version = $parent->childVersion;
        $this->mime = $parent->childMime;
        $this->isInherit = true;
    }

    public function toBytes(): array {
        $w = new TagByteWriter();
        if ($this->example) {
            $w->writeByte(self::K_IS_NULL | 1);
        }
        if ($this->isNull) {
            $w->writeByte(self::K_IS_NULL | 1);
        }
        if ($this->nullable && !$this->isInherit) {
            if (!$this->isNull) {
                $w->writeByte(self::K_NULLABLE | 1);
            }
        }
        if (!empty($this->desc) && !$this->isInherit) {
            $this->writeSizedString($w, self::K_DESC, $this->desc);
        }
        if ($this->type !== ValueType::UNKNOWN && !$this->isInherit) {
            if ($this->shouldEmitExplicitType()) {
                $w->writeByte(self::K_TYPE);
                $w->writeByte($this->type->code());
            }
        }
        if ($this->raw && !$this->isInherit) {
            $w->writeByte(self::K_RAW | 1);
        }
        if ($this->allowEmpty && !$this->isInherit) {
            $w->writeByte(self::K_ALLOW_EMPTY | 1);
        }
        if ($this->unique && !$this->isInherit) {
            $w->writeByte(self::K_UNIQUE | 1);
        }
        if (!empty($this->defaultValue) && !$this->isInherit) {
            $this->writeShortString($w, self::K_DEFAULT, $this->defaultValue);
        }
        if (!empty($this->min) && !$this->isInherit) {
            $this->writeShortString($w, self::K_MIN, $this->min);
        }
        if (!empty($this->max) && !$this->isInherit) {
            $this->writeShortString($w, self::K_MAX, $this->max);
        }
        if ($this->size !== 0 && !$this->isInherit) {
            $this->encodeUint64($w, self::K_SIZE, $this->size);
        }
        if (!empty($this->enumValues) && !$this->isInherit) {
            $this->type = ValueType::ENUM;
            $this->writeSizedString($w, self::K_ENUM, $this->enumValues);
        }
        if (!empty($this->pattern) && !$this->isInherit) {
            $this->writeShortString($w, self::K_PATTERN, $this->pattern);
        }
        if ($this->locationHours !== 0 && !$this->isInherit) {
            $v = (string)$this->locationHours;
            $w->writeByte(self::K_LOCATION | strlen($v));
            $w->writeAscii($v);
        }
        if ($this->version !== self::DEFAULT_VERSION && !$this->isInherit) {
            $this->encodeUint64($w, self::K_VERSION, $this->version);
        }
        if (!empty($this->mime) && !$this->isInherit) {
            $m = MimeWire::parse($this->mime);
            if ($m < 7) {
                $w->writeByte(self::K_MIME | $m);
            } else {
                $w->writeByte(self::K_MIME | 7);
                $w->writeByte($m);
            }
        }
        if (!empty($this->childDesc)) {
            $this->writeSizedString($w, self::K_CHILD_DESC, $this->childDesc);
        }
        if ($this->childType !== ValueType::UNKNOWN) {
            if ($this->shouldEmitChildType()) {
                $w->writeByte(self::K_CHILD_TYPE);
                $w->writeByte($this->childType->code());
            }
        }
        if ($this->childRaw) {
            $w->writeByte(self::K_CHILD_RAW | 1);
        }
        if ($this->childNullable) {
            $w->writeByte(self::K_CHILD_NULLABLE | 1);
        }
        if ($this->childAllowEmpty) {
            $w->writeByte(self::K_CHILD_ALLOW_EMPTY | 1);
        }
        if ($this->childUnique) {
            $w->writeByte(self::K_CHILD_UNIQUE | 1);
        }
        if (!empty($this->childDefault)) {
            $this->writeShortString($w, self::K_CHILD_DEFAULT, $this->childDefault);
        }
        if (!empty($this->childMin)) {
            $this->writeShortString($w, self::K_CHILD_MIN, $this->childMin);
        }
        if (!empty($this->childMax)) {
            $this->writeShortString($w, self::K_CHILD_MAX, $this->childMax);
        }
        if ($this->childSize !== 0) {
            $this->encodeUint64($w, self::K_CHILD_SIZE, $this->childSize);
        }
        if (!empty($this->childEnum)) {
            $this->childType = ValueType::ENUM;
            $this->writeSizedString($w, self::K_CHILD_ENUM, $this->childEnum);
        }
        if (!empty($this->childPattern)) {
            $this->writeShortString($w, self::K_CHILD_PATTERN, $this->childPattern);
        }
        if ($this->childLocationHours !== 0) {
            $v = (string)$this->childLocationHours;
            $w->writeByte(self::K_CHILD_LOCATION | strlen($v));
            $w->writeAscii($v);
        }
        if ($this->childVersion !== self::DEFAULT_VERSION) {
            $this->encodeUint64($w, self::K_CHILD_VERSION, $this->childVersion);
        }
        if (!empty($this->childMime)) {
            $m = MimeWire::parse($this->childMime);
            if ($m < 7) {
                $w->writeByte(self::K_CHILD_MIME | $m);
            } else {
                $w->writeByte(self::K_CHILD_MIME | 7);
                $w->writeByte($m);
            }
        }
        return $w->toByteArray();
    }

    private function shouldEmitExplicitType(): bool {
        return match($this->type) {
            ValueType::STRING, ValueType::BYTES, ValueType::INT, ValueType::FLOAT64,
            ValueType::BOOL, ValueType::STRUCT, ValueType::SLICE => false,
            ValueType::ARRAY => $this->size > 0 ? false : true,
            ValueType::ENUM => !empty($this->enumValues) ? false : true,
            default => true
        };
    }

    private function shouldEmitChildType(): bool {
        return match($this->childType) {
            ValueType::STRING, ValueType::INT, ValueType::FLOAT64,
            ValueType::BOOL, ValueType::STRUCT, ValueType::SLICE => false,
            ValueType::ARRAY => $this->childSize > 0 ? false : true,
            ValueType::ENUM => !empty($this->childEnum) ? false : true,
            default => true
        };
    }

    private function writeSizedString(TagByteWriter $w, int $key, string $s): void {
        $b = unpack('C*', $s);
        $l = count($b);
        if ($l <= 5) {
            $w->writeByte($key | $l, ...$b);
        } elseif ($l <= 0xFF) {
            $w->writeByte($key | 6, $l, ...$b);
        } else {
            $w->writeByte($key | 7, ($l >> 8) & 0xFF, $l & 0xFF, ...$b);
        }
    }

    private function writeShortString(TagByteWriter $w, int $key, string $s): void {
        $b = unpack('C*', $s);
        $l = count($b);
        if ($l < 7) {
            $w->writeByte($key | $l, ...$b);
        } else {
            $w->writeByte($key | 7, $l, ...$b);
        }
    }

    private function encodeUint64(TagByteWriter $buf, int $sign, int $uv): void {
        if ($uv < 0) {
            throw new \Exception('unsigned expected');
        }
        if ($uv <= WireConstants::MAX_1) {
            $buf->writeByte($sign, $uv);
        } elseif ($uv <= WireConstants::MAX_2) {
            $buf->writeByte($sign | 1, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_3) {
            $buf->writeByte($sign | 2, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } elseif ($uv <= WireConstants::MAX_4) {
            $buf->writeByte($sign | 3, ($uv >> 24) & 0xFF, ($uv >> 16) & 0xFF, ($uv >> 8) & 0xFF, $uv & 0xFF);
        } else {
            throw new \Exception('uint64 too large');
        }
    }

    public static function empty(): self {
        return new self();
    }

    public static function fromAnnotation(MM $ann): self {
        $t = new self();
        $t->name = $ann->name;
        $t->isNull = $ann->isNull;
        $t->example = $ann->example;
        $t->desc = $ann->desc;
        $t->type = $ann->type;
        $t->raw = $ann->raw;
        $t->nullable = $ann->nullable;
        $t->allowEmpty = $ann->allowEmpty;
        $t->unique = $ann->unique;
        $t->defaultValue = $ann->defaultValue;
        $t->min = $ann->min;
        $t->max = $ann->max;
        $t->size = $ann->size;
        $t->enumValues = $ann->enumValues;
        if (!empty($t->enumValues)) {
            $t->type = ValueType::ENUM;
        }
        $t->pattern = $ann->pattern;
        $t->locationHours = $ann->location;
        $t->version = $ann->version;
        $t->mime = $ann->mime;
        $t->childDesc = $ann->childDesc;
        $t->childType = $ann->childType;
        $t->childRaw = $ann->childRaw;
        $t->childNullable = $ann->childNullable;
        $t->childAllowEmpty = $ann->childAllowEmpty;
        $t->childUnique = $ann->childUnique;
        $t->childDefault = $ann->childDefault;
        $t->childMin = $ann->childMin;
        $t->childMax = $ann->childMax;
        $t->childSize = $ann->childSize;
        $t->childEnum = $ann->childEnum;
        if (!empty($t->childEnum)) {
            $t->childType = ValueType::ENUM;
        }
        $t->childPattern = $ann->childPattern;
        $t->childLocationHours = $ann->childLocation;
        $t->childVersion = $ann->childVersion;
        $t->childMime = $ann->childMime;
        return $t;
    }

    public const K_IS_NULL = 0 << 3;
    public const K_EXAMPLE = 1 << 3;
    public const K_DESC = 2 << 3;
    public const K_TYPE = 3 << 3;
    public const K_RAW = 4 << 3;
    public const K_NULLABLE = 5 << 3;
    public const K_ALLOW_EMPTY = 6 << 3;
    public const K_UNIQUE = 7 << 3;
    public const K_DEFAULT = 8 << 3;
    public const K_MIN = 9 << 3;
    public const K_MAX = 10 << 3;
    public const K_SIZE = 11 << 3;
    public const K_ENUM = 12 << 3;
    public const K_PATTERN = 13 << 3;
    public const K_LOCATION = 14 << 3;
    public const K_VERSION = 15 << 3;
    public const K_MIME = 16 << 3;
    public const K_CHILD_DESC = 17 << 3;
    public const K_CHILD_TYPE = 18 << 3;
    public const K_CHILD_RAW = 19 << 3;
    public const K_CHILD_NULLABLE = 20 << 3;
    public const K_CHILD_ALLOW_EMPTY = 21 << 3;
    public const K_CHILD_UNIQUE = 22 << 3;
    public const K_CHILD_DEFAULT = 23 << 3;
    public const K_CHILD_MIN = 24 << 3;
    public const K_CHILD_MAX = 25 << 3;
    public const K_CHILD_SIZE = 26 << 3;
    public const K_CHILD_ENUM = 27 << 3;
    public const K_CHILD_PATTERN = 28 << 3;
    public const K_CHILD_LOCATION = 29 << 3;
    public const K_CHILD_VERSION = 30 << 3;
    public const K_CHILD_MIME = 31 << 3;
}

class TagByteWriter {
    private array $buf;
    private int $len;

    public function __construct() {
        $this->buf = array_fill(0, 64, 0);
        $this->len = 0;
    }

    public function writeByte(int ...$bs): void {
        foreach ($bs as $b) {
            if ($this->len >= count($this->buf)) {
                $this->buf = array_pad($this->buf, count($this->buf) * 2, 0);
            }
            $this->buf[$this->len++] = $b;
        }
    }

    public function writeAscii(string $s): void {
        $b = unpack('C*', $s);
        $this->writeByte(...$b);
    }

    public function toByteArray(): array {
        return array_slice($this->buf, 0, $this->len);
    }
}
