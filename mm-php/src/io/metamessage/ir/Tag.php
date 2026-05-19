<?php

namespace io\metamessage\ir;

use io\metamessage\core\MM;
use io\metamessage\core\WireConstants;

class Tag
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
    public int $locationHours = 0;
    public int $version = self::DEFAULT_VERSION;
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
    public int $childLocationHours = 0;
    public int $childVersion = self::DEFAULT_VERSION;
    public string $childMime = '';

    public bool $isInherit = false;

    public const DEFAULT_VERSION = 0;

    const T_IS_NULL = 'is_null';
    const T_EXAMPLE = 'example';

    const T_NAME = 'name';
    const T_DESC = 'desc';
    const T_TYPE = 'type';
    const T_RAW = 'raw';
    const T_NULLABLE = 'nullable';
    const T_ALLOW_EMPTY = 'allow_empty';
    const T_UNIQUE = 'unique';
    const T_DEFAULT = 'default';
    const T_MIN = 'min';
    const T_MAX = 'max';
    const T_SIZE = 'size';
    const T_ENUM = 'enum';
    const T_PATTERN = 'pattern';
    const T_LOCATION = 'location';
    const T_VERSION = 'version';
    const T_MIME = 'mime';

    const T_CHILD_DESC = 'child_desc';
    const T_CHILD_TYPE = 'child_type';
    const T_CHILD_RAW = 'child_raw';
    const T_CHILD_NULLABLE = 'child_nullable';
    const T_CHILD_ALLOW_EMPTY = 'child_allow_empty';
    const T_CHILD_UNIQUE = 'child_unique';
    const T_CHILD_DEFAULT = 'child_default';
    const T_CHILD_MIN = 'child_min';
    const T_CHILD_MAX = 'child_max';
    const T_CHILD_SIZE = 'child_size';
    const T_CHILD_ENUM = 'child_enum';
    const T_CHILD_PATTERN = 'child_pattern';
    const T_CHILD_LOCATION = 'child_location';
    const T_CHILD_VERSION = 'child_version';
    const T_CHILD_MIME = 'child_mime';

    const K_IS_NULL = 0 << 3;
    const K_EXAMPLE = 1 << 3;
    const K_DESC = 2 << 3;
    const K_TYPE = 3 << 3;
    const K_RAW = 4 << 3;
    const K_NULLABLE = 5 << 3;
    const K_ALLOW_EMPTY = 6 << 3;
    const K_UNIQUE = 7 << 3;
    const K_DEFAULT = 8 << 3;
    const K_MIN = 9 << 3;
    const K_MAX = 10 << 3;
    const K_SIZE = 11 << 3;
    const K_ENUM = 12 << 3;
    const K_PATTERN = 13 << 3;
    const K_LOCATION = 14 << 3;
    const K_VERSION = 15 << 3;
    const K_MIME = 16 << 3;
    const K_CHILD_DESC = 17 << 3;
    const K_CHILD_TYPE = 18 << 3;
    const K_CHILD_RAW = 19 << 3;
    const K_CHILD_NULLABLE = 20 << 3;
    const K_CHILD_ALLOW_EMPTY = 21 << 3;
    const K_CHILD_UNIQUE = 22 << 3;
    const K_CHILD_DEFAULT = 23 << 3;
    const K_CHILD_MIN = 24 << 3;
    const K_CHILD_MAX = 25 << 3;
    const K_CHILD_SIZE = 26 << 3;
    const K_CHILD_ENUM = 27 << 3;
    const K_CHILD_PATTERN = 28 << 3;
    const K_CHILD_LOCATION = 29 << 3;
    const K_CHILD_VERSION = 30 << 3;
    const K_CHILD_MIME = 31 << 3;

    const MAX_1_BYTE = 0xFF;
    const MAX_2_BYTE = 0xFFFF;
    const MAX_3_BYTE = 0xFFFFFF;
    const MAX_4_BYTE = 0xFFFFFFFF;
    const MAX_5_BYTE = 0xFFFFFFFFFF;
    const MAX_6_BYTE = 0xFFFFFFFFFFFF;
    const MAX_7_BYTE = 0xFFFFFFFFFFFFFF;
    const MAX_8_BYTE = 0xFFFFFFFFFFFFFFFF;

    const INT_LEN_MASK = 0b11111;
    const INT_LEN_1_BYTE = self::INT_LEN_MASK - 7;
    const INT_LEN_2_BYTE = self::INT_LEN_MASK - 6;
    const INT_LEN_3_BYTE = self::INT_LEN_MASK - 5;
    const INT_LEN_4_BYTE = self::INT_LEN_MASK - 4;
    const INT_LEN_5_BYTE = self::INT_LEN_MASK - 3;
    const INT_LEN_6_BYTE = self::INT_LEN_MASK - 2;
    const INT_LEN_7_BYTE = self::INT_LEN_MASK - 1;
    const INT_LEN_8_BYTE = self::INT_LEN_MASK;

    public function __construct()
    {
        $this->type = ValueType::UNKNOWN;
        $this->childType = ValueType::UNKNOWN;
    }

    public static function newTag(): self
    {
        return new self();
    }

    public static function empty(): self
    {
        return new self();
    }

    public static function fromAnnotation(MM $ann): self
    {
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

    public function copy(): self
    {
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

    public function inheritFromArrayParent(?Tag $parent): void
    {
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

    public function toBytes(): array
    {
        return $this->bytes();
    }

    public function inherit(Tag $tag): void
    {
        $this->isInherit = true;

        if ($tag->childDesc !== '') {
            $this->desc = $tag->childDesc;
        }

        if ($tag->childType !== ValueType::UNKNOWN) {
            $this->type = $tag->childType;
        }

        if ($tag->childRaw) {
            $this->raw = $tag->childRaw;
        }

        if ($tag->childNullable) {
            $this->nullable = $tag->childNullable;
        }

        if ($tag->childAllowEmpty) {
            $this->allowEmpty = $tag->childAllowEmpty;
        }

        if ($tag->childUnique) {
            $this->unique = $tag->childUnique;
        }

        if ($tag->childDefault !== '') {
            $this->defaultValue = $tag->childDefault;
        }

        if ($tag->childMin !== '') {
            $this->min = $tag->childMin;
        }

        if ($tag->childMax !== '') {
            $this->max = $tag->childMax;
        }

        if ($tag->childSize !== 0) {
            $this->size = $tag->childSize;
        }

        if ($tag->childEnum !== '') {
            $this->enumValues = $tag->childEnum;
        }

        if ($tag->childPattern !== '') {
            $this->pattern = $tag->childPattern;
        }

        if ($tag->childLocationHours !== 0) {
            $this->locationHours = $tag->childLocationHours;
        }

        if ($tag->childVersion !== self::DEFAULT_VERSION) {
            $this->version = $tag->childVersion;
        }

        if ($tag->childMime !== '') {
            $this->mime = $tag->childMime;
        }
    }

    public function getPattern(): ?string
    {
        if ($this->pattern === '') {
            return null;
        }
        return $this->pattern;
    }

    public function json(): string
    {
        return json_encode($this, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE);
    }

    public function toString(): string
    {
        $b = '';
        $first = true;
        $add = function (string $s) use (&$b, &$first): void {
            if (!$first) {
                $b .= '; ';
            }
            $b .= $s;
            $first = false;
        };

        if ($this->type !== ValueType::UNKNOWN && !$this->isInherit) {
            if ($this->type === ValueType::STRING ||
                $this->type === ValueType::INT ||
                $this->type === ValueType::FLOAT64 ||
                $this->type === ValueType::BOOL ||
                $this->type === ValueType::STRUCT ||
                $this->type === ValueType::SLICE) {
            } else {
                if ($this->type === ValueType::ARRAY && $this->size > 0 ||
                    $this->type === ValueType::ENUM && $this->enumValues !== '') {
                } else {
                    $add(self::T_TYPE . '=' . $this->type->wireName());
                }
            }
        }

        if ($this->example) {
            $add(self::T_EXAMPLE);
        }

        if ($this->isNull) {
            $add(self::T_IS_NULL);
        }

        if ($this->nullable && !$this->isInherit) {
            if (!$this->isNull) {
                $add(self::T_NULLABLE);
            }
        }

        if ($this->desc !== '' && !$this->isInherit) {
            $add(self::T_DESC . '=' . $this->quote($this->desc));
        }

        if ($this->raw && !$this->isInherit) {
            $add(self::T_RAW);
        }

        if ($this->allowEmpty && !$this->isInherit) {
            $add(self::T_ALLOW_EMPTY);
        }

        if ($this->unique && !$this->isInherit) {
            $add(self::T_UNIQUE);
        }

        if ($this->defaultValue !== '' && !$this->isInherit) {
            $add(self::T_DEFAULT . '=' . $this->defaultValue);
        }

        if ($this->min !== '' && !$this->isInherit) {
            $add(self::T_MIN . '=' . $this->min);
        }

        if ($this->max !== '' && !$this->isInherit) {
            $add(self::T_MAX . '=' . $this->max);
        }

        if ($this->size !== 0 && !$this->isInherit) {
            $add(self::T_SIZE . '=' . $this->size);
        }

        if ($this->enumValues !== '' && !$this->isInherit) {
            $add(self::T_ENUM . '=' . $this->enumValues);
        }

        if ($this->pattern !== '' && !$this->isInherit) {
            $add(self::T_PATTERN . '=' . $this->pattern);
        }

        if ($this->locationHours !== 0 && !$this->isInherit) {
            $add(self::T_LOCATION . '=' . $this->locationHours);
        }

        if ($this->version !== self::DEFAULT_VERSION && !$this->isInherit) {
            $add(self::T_VERSION . '=' . $this->version);
        }

        if ($this->mime !== '' && !$this->isInherit) {
            $add(self::T_MIME . '=' . $this->mime);
        }

        if ($this->childDesc !== '') {
            $add(self::T_CHILD_DESC . '=' . $this->quote($this->childDesc));
        }

        if ($this->childType !== ValueType::UNKNOWN) {
            if ($this->childType === ValueType::STRING ||
                $this->childType === ValueType::INT ||
                $this->childType === ValueType::FLOAT64 ||
                $this->childType === ValueType::BOOL ||
                $this->childType === ValueType::STRUCT ||
                $this->childType === ValueType::SLICE) {
            } else {
                if ($this->childType === ValueType::ARRAY && $this->childSize > 0 ||
                    $this->childType === ValueType::ENUM && $this->childEnum !== '') {
                } else {
                    $add(self::T_CHILD_TYPE . '=' . $this->childType->wireName());
                }
            }
        }

        if ($this->childRaw) {
            $add(self::T_CHILD_RAW);
        }

        if ($this->childNullable) {
            $add(self::T_CHILD_NULLABLE);
        }

        if ($this->childAllowEmpty) {
            $add(self::T_CHILD_ALLOW_EMPTY);
        }

        if ($this->childUnique) {
            $add(self::T_CHILD_UNIQUE);
        }

        if ($this->childDefault !== '') {
            $add(self::T_CHILD_DEFAULT . '=' . $this->childDefault);
        }

        if ($this->childMin !== '') {
            $add(self::T_CHILD_MIN . '=' . $this->childMin);
        }

        if ($this->childMax !== '') {
            $add(self::T_CHILD_MAX . '=' . $this->childMax);
        }

        if ($this->childSize !== 0) {
            $add(self::T_CHILD_SIZE . '=' . $this->childSize);
        }

        if ($this->childEnum !== '') {
            $add(self::T_CHILD_ENUM . '=' . $this->childEnum);
        }

        if ($this->childPattern !== '') {
            $add(self::T_CHILD_PATTERN . '=' . $this->childPattern);
        }

        if ($this->childLocationHours !== 0) {
            $add(self::T_CHILD_LOCATION . '=' . $this->childLocationHours);
        }

        if ($this->childVersion !== self::DEFAULT_VERSION) {
            $add(self::T_CHILD_VERSION . '=' . $this->childVersion);
        }

        if ($this->childMime !== '') {
            $add(self::T_CHILD_MIME . '=' . $this->childMime);
        }

        return $b;
    }

    public function bytes(): array
    {
        $w = new TagByteWriter();

        if ($this->example) {
            $w->writeByte(self::K_EXAMPLE | 1);
        }

        if ($this->isNull) {
            $w->writeByte(self::K_IS_NULL | 1);
        }

        if ($this->nullable && !$this->isInherit) {
            if (!$this->isNull) {
                $w->writeByte(self::K_NULLABLE | 1);
            }
        }

        if ($this->desc !== '' && !$this->isInherit) {
            $l = strlen($this->desc);
            if ($l <= 5) {
                $w->writeByte(self::K_DESC | $l);
                $w->writeAscii($this->desc);
            } elseif ($l <= 0xFF) {
                $w->writeByte(self::K_DESC | 6);
                $w->writeByte($l);
                $w->writeAscii($this->desc);
            } elseif ($l <= 0xFFFF) {
                $w->writeByte(self::K_DESC | 7);
                $w->writeByte(($l >> 8) & 0xFF);
                $w->writeByte($l & 0xFF);
                $w->writeAscii($this->desc);
            }
        }

        if ($this->type !== ValueType::UNKNOWN && !$this->isInherit) {
            if ($this->type === ValueType::STRING ||
                $this->type === ValueType::BYTES ||
                $this->type === ValueType::INT ||
                $this->type === ValueType::FLOAT64 ||
                $this->type === ValueType::BOOL ||
                $this->type === ValueType::STRUCT ||
                $this->type === ValueType::SLICE) {
            } else {
                if ($this->type === ValueType::ARRAY && $this->size > 0 ||
                    $this->type === ValueType::ENUM && $this->enumValues !== '') {
                } else {
                    $w->writeByte(self::K_TYPE);
                    $w->writeByte($this->type->code());
                }
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

        if ($this->defaultValue !== '' && !$this->isInherit) {
            $l = strlen($this->defaultValue);
            if ($l < 7) {
                $w->writeByte(self::K_DEFAULT | $l);
                $w->writeAscii($this->defaultValue);
            } else {
                $w->writeByte(self::K_DEFAULT | 7);
                $w->writeByte($l);
                $w->writeAscii($this->defaultValue);
            }
        }

        if ($this->min !== '' && !$this->isInherit) {
            $l = strlen($this->min);
            if ($l < 7) {
                $w->writeByte(self::K_MIN | $l);
                $w->writeAscii($this->min);
            } else {
                $w->writeByte(self::K_MIN | 7);
                $w->writeByte($l);
                $w->writeAscii($this->min);
            }
        }

        if ($this->max !== '' && !$this->isInherit) {
            $l = strlen($this->max);
            if ($l < 7) {
                $w->writeByte(self::K_MAX | $l);
                $w->writeAscii($this->max);
            } else {
                $w->writeByte(self::K_MAX | 7);
                $w->writeByte($l);
                $w->writeAscii($this->max);
            }
        }

        if ($this->size !== 0 && !$this->isInherit) {
            self::encodeU64Static($w, self::K_SIZE, $this->size);
        }

        if ($this->enumValues !== '' && !$this->isInherit) {
            $l = strlen($this->enumValues);
            if ($l <= 5) {
                $w->writeByte(self::K_ENUM | $l);
                $w->writeAscii($this->enumValues);
            } elseif ($l <= 0xFF) {
                $w->writeByte(self::K_ENUM | 6);
                $w->writeByte($l);
                $w->writeAscii($this->enumValues);
            } elseif ($l <= 0xFFFF) {
                $w->writeByte(self::K_ENUM | 7);
                $w->writeByte(($l >> 8) & 0xFF);
                $w->writeByte($l & 0xFF);
                $w->writeAscii($this->enumValues);
            }
        }

        if ($this->pattern !== '' && !$this->isInherit) {
            $l = strlen($this->pattern);
            if ($l < 7) {
                $w->writeByte(self::K_PATTERN | $l);
                $w->writeAscii($this->pattern);
            } else {
                $w->writeByte(self::K_PATTERN | 7);
                $w->writeByte($l);
                $w->writeAscii($this->pattern);
            }
        }

        if ($this->locationHours !== 0 && !$this->isInherit) {
            $v = (string)$this->locationHours;
            $w->writeByte(self::K_LOCATION | strlen($v));
            $w->writeAscii($v);
        }

        if ($this->version !== self::DEFAULT_VERSION && !$this->isInherit) {
            self::encodeU64Static($w, self::K_VERSION, $this->version);
        }

        if ($this->mime !== '' && !$this->isInherit) {
            $l = Mime::parse($this->mime);
            if ($l < 7) {
                $w->writeByte(self::K_MIME | $l);
            } else {
                $w->writeByte(self::K_MIME | 7);
                $w->writeByte($l);
            }
        }

        if ($this->childDesc !== '') {
            $l = strlen($this->childDesc);
            if ($l <= 5) {
                $w->writeByte(self::K_CHILD_DESC | $l);
                $w->writeAscii($this->childDesc);
            } elseif ($l <= 0xFF) {
                $w->writeByte(self::K_CHILD_DESC | 6);
                $w->writeByte($l);
                $w->writeAscii($this->childDesc);
            } elseif ($l <= 0xFFFF) {
                $w->writeByte(self::K_CHILD_DESC | 7);
                $w->writeByte(($l >> 8) & 0xFF);
                $w->writeByte($l & 0xFF);
                $w->writeAscii($this->childDesc);
            }
        }

        if ($this->childType !== ValueType::UNKNOWN) {
            if ($this->childType === ValueType::STRING ||
                $this->childType === ValueType::INT ||
                $this->childType === ValueType::FLOAT64 ||
                $this->childType === ValueType::BOOL ||
                $this->childType === ValueType::STRUCT ||
                $this->childType === ValueType::SLICE) {
            } else {
                if ($this->childType === ValueType::ARRAY && $this->childSize > 0 ||
                    $this->childType === ValueType::ENUM && $this->childEnum !== '') {
                } else {
                    $w->writeByte(self::K_CHILD_TYPE);
                    $w->writeByte($this->childType->code());
                }
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

        if ($this->childDefault !== '') {
            $l = strlen($this->childDefault);
            if ($l < 7) {
                $w->writeByte(self::K_CHILD_DEFAULT | $l);
                $w->writeAscii($this->childDefault);
            } else {
                $w->writeByte(self::K_CHILD_DEFAULT | 7);
                $w->writeByte($l);
                $w->writeAscii($this->childDefault);
            }
        }

        if ($this->childMin !== '') {
            $l = strlen($this->childMin);
            if ($l < 7) {
                $w->writeByte(self::K_CHILD_MIN | $l);
                $w->writeAscii($this->childMin);
            } else {
                $w->writeByte(self::K_CHILD_MIN | 7);
                $w->writeByte($l);
                $w->writeAscii($this->childMin);
            }
        }

        if ($this->childMax !== '') {
            $l = strlen($this->childMax);
            if ($l < 7) {
                $w->writeByte(self::K_CHILD_MAX | $l);
                $w->writeAscii($this->childMax);
            } else {
                $w->writeByte(self::K_CHILD_MAX | 7);
                $w->writeByte($l);
                $w->writeAscii($this->childMax);
            }
        }

        if ($this->childSize !== 0) {
            self::encodeU64Static($w, self::K_CHILD_SIZE, $this->childSize);
        }

        if ($this->childEnum !== '') {
            $l = strlen($this->childEnum);
            if ($l <= 5) {
                $w->writeByte(self::K_CHILD_ENUM | $l);
                $w->writeAscii($this->childEnum);
            } elseif ($l <= 0xFF) {
                $w->writeByte(self::K_CHILD_ENUM | 6);
                $w->writeByte($l);
                $w->writeAscii($this->childEnum);
            } elseif ($l <= 0xFFFF) {
                $w->writeByte(self::K_CHILD_ENUM | 7);
                $w->writeByte(($l >> 8) & 0xFF);
                $w->writeByte($l & 0xFF);
                $w->writeAscii($this->childEnum);
            }
        }

        if ($this->childPattern !== '') {
            $l = strlen($this->childPattern);
            if ($l < 7) {
                $w->writeByte(self::K_CHILD_PATTERN | $l);
                $w->writeAscii($this->childPattern);
            } else {
                $w->writeByte(self::K_CHILD_PATTERN | 7);
                $w->writeByte($l);
                $w->writeAscii($this->childPattern);
            }
        }

        if ($this->childLocationHours !== 0) {
            $v = (string)$this->childLocationHours;
            $w->writeByte(self::K_CHILD_LOCATION | strlen($v));
            $w->writeAscii($v);
        }

        if ($this->childVersion !== self::DEFAULT_VERSION) {
            self::encodeU64Static($w, self::K_CHILD_VERSION, $this->childVersion);
        }

        if ($this->childMime !== '') {
            $l = strlen($this->childMime);
            if ($l < 7) {
                $w->writeByte(self::K_CHILD_MIME | $l);
                $w->writeAscii($this->childMime);
            } else {
                $w->writeByte(self::K_CHILD_MIME | 7);
                $w->writeByte($l);
                $w->writeAscii($this->childMime);
            }
        }

        return $w->toByteArray();
    }

    public static function encodeU64Static(TagByteWriter $buf, int $sign, int $uv): void
    {
        if ($uv < 0) {
            throw new \Exception('unsigned expected');
        }
        if ($uv <= self::MAX_1_BYTE) {
            $sign |= 0;
            $buf->writeByte($sign);
            $buf->writeByte($uv);
        } elseif ($uv <= self::MAX_2_BYTE) {
            $sign |= 1;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        } elseif ($uv <= self::MAX_3_BYTE) {
            $sign |= 2;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 16) & 0xFF);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        } elseif ($uv <= self::MAX_4_BYTE) {
            $sign |= 3;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 24) & 0xFF);
            $buf->writeByte(($uv >> 16) & 0xFF);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        } elseif ($uv <= self::MAX_5_BYTE) {
            $sign |= 4;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 32) & 0xFF);
            $buf->writeByte(($uv >> 24) & 0xFF);
            $buf->writeByte(($uv >> 16) & 0xFF);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        } elseif ($uv <= self::MAX_6_BYTE) {
            $sign |= 5;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 40) & 0xFF);
            $buf->writeByte(($uv >> 32) & 0xFF);
            $buf->writeByte(($uv >> 24) & 0xFF);
            $buf->writeByte(($uv >> 16) & 0xFF);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        } elseif ($uv <= self::MAX_7_BYTE) {
            $sign |= 6;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 48) & 0xFF);
            $buf->writeByte(($uv >> 40) & 0xFF);
            $buf->writeByte(($uv >> 32) & 0xFF);
            $buf->writeByte(($uv >> 24) & 0xFF);
            $buf->writeByte(($uv >> 16) & 0xFF);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        } elseif ($uv <= self::MAX_8_BYTE) {
            $sign |= 7;
            $buf->writeByte($sign);
            $buf->writeByte(($uv >> 56) & 0xFF);
            $buf->writeByte(($uv >> 48) & 0xFF);
            $buf->writeByte(($uv >> 40) & 0xFF);
            $buf->writeByte(($uv >> 32) & 0xFF);
            $buf->writeByte(($uv >> 24) & 0xFF);
            $buf->writeByte(($uv >> 16) & 0xFF);
            $buf->writeByte(($uv >> 8) & 0xFF);
            $buf->writeByte($uv & 0xFF);
        }
    }

    public static function mergeTag(?Tag $dst, ?Tag $src): ?Tag
    {
        if ($src === null) {
            return $dst;
        }

        if ($dst === null) {
            return $src;
        }

        if ($src->isNull) {
            $dst->isNull = $src->isNull;
        }

        if ($src->example) {
            $dst->example = $src->example;
        }

        if ($src->desc !== '') {
            $dst->desc = $src->desc;
        }

        if ($src->type !== ValueType::UNKNOWN) {
            $dst->type = $src->type;
        }

        if ($src->raw) {
            $dst->raw = true;
        }

        if ($src->nullable) {
            $dst->nullable = true;
        }

        if ($src->allowEmpty) {
            $dst->allowEmpty = true;
        }

        if ($src->unique) {
            $dst->unique = true;
        }

        if ($src->defaultValue !== '') {
            $dst->defaultValue = $src->defaultValue;
        }

        if ($src->min !== '') {
            $dst->min = $src->min;
        }

        if ($src->max !== '') {
            $dst->max = $src->max;
        }

        if ($src->size !== 0) {
            $dst->size = $src->size;
        }

        if ($src->enumValues !== '') {
            $dst->enumValues = $src->enumValues;
        }

        if ($src->pattern !== '') {
            $dst->pattern = $src->pattern;
        }

        if ($src->locationHours !== 0) {
            $dst->locationHours = $src->locationHours;
        }

        if ($src->version !== self::DEFAULT_VERSION) {
            $dst->version = $src->version;
        }

        if ($src->mime !== '') {
            $dst->mime = $src->mime;
        }

        if ($src->childDesc !== '') {
            $dst->childDesc = $src->childDesc;
        }

        if ($src->childType !== ValueType::UNKNOWN) {
            $dst->childType = $src->childType;
        }

        if ($src->childRaw) {
            $dst->childRaw = true;
        }

        if ($src->childNullable) {
            $dst->childNullable = true;
        }

        if ($src->childAllowEmpty) {
            $dst->childAllowEmpty = true;
        }

        if ($src->childUnique) {
            $dst->childUnique = true;
        }

        if ($src->childDefault !== '') {
            $dst->childDefault = $src->childDefault;
        }

        if ($src->childMin !== '') {
            $dst->childMin = $src->childMin;
        }

        if ($src->childMax !== '') {
            $dst->childMax = $src->childMax;
        }

        if ($src->childSize !== 0) {
            $dst->childSize = $src->childSize;
        }

        if ($src->childEnum !== '') {
            $dst->childEnum = $src->childEnum;
        }

        if ($src->childPattern !== '') {
            $dst->childPattern = $src->childPattern;
        }

        if ($src->childLocationHours !== 0) {
            $dst->childLocationHours = $src->childLocationHours;
        }

        if ($src->childVersion !== self::DEFAULT_VERSION) {
            $dst->childVersion = $src->childVersion;
        }

        if ($src->childMime !== '') {
            $dst->childMime = $src->childMime;
        }

        return $dst;
    }

    public static function parseMMTag(string $tag): self
    {
        $r = self::newTag();
        $tag = trim($tag);
        if (str_starts_with($tag, '//')) {
            $tag = substr($tag, 2);
        }
        $tag = trim($tag);
        $tagLower = strtolower($tag);
        if (str_starts_with($tagLower, 'mm:')) {
            $tag = substr($tag, 3);
        }
        $tag = trim($tag);
        if ($tag === '') {
            return $r;
        }

        $parts = self::splitTag($tag);
        foreach ($parts as $p) {
            $p = trim($p);
            if ($p === '') {
                continue;
            }
            $k = $p;
            $v = '';
            $eqPos = strpos($p, '=');
            if ($eqPos !== false) {
                $k = trim(substr($p, 0, $eqPos));
                $v = trim(substr($p, $eqPos + 1));
            } else {
                $k = trim($p);
                $v = '';
            }

            if (strlen($v) >= 2 && $v[0] === '"' && $v[strlen($v) - 1] === '"') {
                $unquoted = substr($v, 1, -1);
                $v = $unquoted;
            }

            $lower = strtolower($k);
            switch ($lower) {
                case self::T_NAME:
                    $r->name = $v;
                    break;

                case self::T_IS_NULL:
                    $r->isNull = true;
                    $r->nullable = true;
                    break;

                case self::T_EXAMPLE:
                    $r->example = true;
                    break;

                case self::T_DESC:
                    $r->desc = $v;
                    break;

                case self::T_TYPE:
                    $r->type = ValueType::parseWireName($v);
                    break;

                case self::T_RAW:
                    $r->raw = true;
                    break;

                case self::T_NULLABLE:
                    $r->nullable = true;
                    break;

                case self::T_ALLOW_EMPTY:
                    $r->allowEmpty = true;
                    break;

                case self::T_UNIQUE:
                    $r->unique = true;
                    break;

                case self::T_DEFAULT:
                    $r->defaultValue = $v;
                    break;

                case self::T_PATTERN:
                    $r->pattern = $v;
                    break;

                case self::T_MIN:
                    $r->min = $v;
                    break;

                case self::T_MAX:
                    $r->max = $v;
                    break;

                case self::T_SIZE:
                    $r->size = (int)$v;
                    break;

                case self::T_ENUM:
                    $r->type = ValueType::ENUM;
                    $r->enumValues = $v;
                    break;

                case self::T_LOCATION:
                    $r->locationHours = (int)$v;
                    break;

                case self::T_VERSION:
                    $r->version = (int)$v;
                    break;

                case self::T_MIME:
                    $r->mime = $v;
                    break;

                case self::T_CHILD_DESC:
                    $r->childDesc = $v;
                    break;

                case self::T_CHILD_TYPE:
                    $r->childType = ValueType::parseWireName($v);
                    break;

                case self::T_CHILD_RAW:
                    $r->childRaw = true;
                    break;

                case self::T_CHILD_NULLABLE:
                    $r->childNullable = true;
                    break;

                case self::T_CHILD_ALLOW_EMPTY:
                    $r->childAllowEmpty = true;
                    break;

                case self::T_CHILD_UNIQUE:
                    $r->childUnique = true;
                    break;

                case self::T_CHILD_DEFAULT:
                    $r->childDefault = $v;
                    break;

                case self::T_CHILD_PATTERN:
                    $r->childPattern = $v;
                    break;

                case self::T_CHILD_MIN:
                    $r->childMin = $v;
                    break;

                case self::T_CHILD_MAX:
                    $r->childMax = $v;
                    break;

                case self::T_CHILD_SIZE:
                    $r->childSize = (int)$v;
                    break;

                case self::T_CHILD_ENUM:
                    $r->childEnum = $v;
                    $r->childType = ValueType::ENUM;
                    break;

                case self::T_CHILD_LOCATION:
                    $r->childLocationHours = (int)$v;
                    break;

                case self::T_CHILD_VERSION:
                    $r->childVersion = (int)$v;
                    break;

                case self::T_CHILD_MIME:
                    $r->childMime = $v;
                    break;

                default:
                    break;
            }
        }
        return $r;
    }

    private static function splitTag(string $tag): array
    {
        if ($tag === '') {
            return [];
        }

        $parts = explode(';', $tag);
        foreach ($parts as $i => $part) {
            $parts[$i] = trim($part);
        }
        return $parts;
    }

    private function quote(string $s): string
    {
        return '"' . addcslashes($s, "\\\"\n\r\t") . '"';
    }
}

class TagByteWriter
{
    private array $buf;
    private int $len;

    public function __construct()
    {
        $this->buf = array_fill(0, 64, 0);
        $this->len = 0;
    }

    public function writeByte(int ...$bs): void
    {
        foreach ($bs as $b) {
            if ($this->len >= count($this->buf)) {
                $this->buf = array_pad($this->buf, count($this->buf) * 2, 0);
            }
            $this->buf[$this->len++] = $b;
        }
    }

    public function writeAscii(string $s): void
    {
        $b = unpack('C*', $s);
        $this->writeByte(...$b);
    }

    public function toByteArray(): array
    {
        return array_slice($this->buf, 0, $this->len);
    }
}