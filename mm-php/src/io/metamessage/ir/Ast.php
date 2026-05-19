<?php

namespace io\metamessage\ir;

class NodeType
{
    public const UNKNOWN = 0;
    public const OBJECT = 1;
    public const ARRAY = 2;
    public const VALUE = 3;
    public const DOC = 4;

    public const LABEL_UNKNOWN = 'unknown';
    public const LABEL_OBJECT = 'object';
    public const LABEL_ARRAY = 'array';
    public const LABEL_VALUE = 'value';
    public const LABEL_DOC = 'doc';

    public static function toString(int $nt): string
    {
        return match ($nt) {
            self::OBJECT => self::LABEL_OBJECT,
            self::ARRAY => self::LABEL_ARRAY,
            self::VALUE => self::LABEL_VALUE,
            self::DOC => self::LABEL_DOC,
            default => self::LABEL_UNKNOWN,
        };
    }

    public static function parse(string $s): int
    {
        return match ($s) {
            self::LABEL_OBJECT => self::OBJECT,
            self::LABEL_ARRAY => self::ARRAY,
            self::LABEL_VALUE => self::VALUE,
            self::LABEL_DOC => self::DOC,
            default => self::UNKNOWN,
        };
    }
}

interface Node
{
    public function getTag(): ?Tag;

    public function getType(): int;

    public function getPath(): string;

    public function setPath(string $path): void;
}

class Field
{
    public string $Key;
    public Node $Value;

    public function __construct(string $key = '', ?Node $value = null)
    {
        $this->Key = $key;
        if ($value !== null) {
            $this->Value = $value;
        }
    }
}

class Object_ implements Node
{
    public array $Fields = [];
    public ?Tag $Tag = null;
    public string $Path = '';

    public function getPath(): string
    {
        return $this->Path;
    }

    public function setPath(string $path): void
    {
        $this->Path = $path;
    }

    public function getType(): int
    {
        return NodeType::OBJECT;
    }

    public function getTag(): ?Tag
    {
        return $this->Tag;
    }
}

class Array_ implements Node
{
    public array $Items = [];
    public ?Tag $Tag = null;
    public string $Path = '';

    public function getPath(): string
    {
        return $this->Path;
    }

    public function setPath(string $path): void
    {
        $this->Path = $path;
    }

    public function getType(): int
    {
        return NodeType::ARRAY;
    }

    public function getTag(): ?Tag
    {
        return $this->Tag;
    }
}

class Value implements Node
{
    public mixed $Data = null;
    public string $Text = '';
    public ?Tag $Tag = null;
    public string $Path = '';

    public function __construct(mixed $data = null, string $text = '', ?Tag $tag = null, string $path = '')
    {
        $this->Data = $data;
        $this->Text = $text;
        $this->Tag = $tag;
        $this->Path = $path;
    }

    public function getPath(): string
    {
        return $this->Path;
    }

    public function setPath(string $path): void
    {
        $this->Path = $path;
    }

    public function getType(): int
    {
        return NodeType::VALUE;
    }

    public function getTag(): ?Tag
    {
        return $this->Tag;
    }
}

class Doc implements Node
{
    public array $Fields = [];
    public ?Tag $Tag = null;
    public string $Path = '';

    public function getPath(): string
    {
        return $this->Path;
    }

    public function setPath(string $path): void
    {
        $this->Path = $path;
    }

    public function getType(): int
    {
        return NodeType::DOC;
    }

    public function getTag(): ?Tag
    {
        return $this->Tag;
    }
}