<?php

namespace io\metamessage\core;

use io\metamessage\ir\Node;
use io\metamessage\ir\Tag;
use io\metamessage\ir\Object_;
use io\metamessage\ir\Array_;
use io\metamessage\ir\Value;
use io\metamessage\ir\Field;
use io\metamessage\ir\ValueType;
use io\metamessage\ir\Constants;
use io\metamessage\jsonc\Jsonc;
use io\metamessage\jsonc\JsoncParser;
use io\metamessage\jsonc\JsoncScanner;
use io\metamessage\jsonc\JsoncToken;
use io\metamessage\jsonc\JsoncTokenType;
use ReflectionClass;
use ReflectionProperty;
use ReflectionNamedType;

class MetaMessage
{
    private const MAX_DEPTH = 32;

    private function __construct() {}

    public static function FromValue(mixed $v, string $tag): array
    {
        $tagObj = Tag::parseMMTag($tag);
        $node = self::valueToNode($v, $tagObj, 0, '');
        $encoder = new WireEncoder();
        return $encoder->encode($node);
    }

    public static function FromJSONC(string $s): array
    {
        $node = self::ParseFromJSONC($s);
        $encoder = new WireEncoder();
        return $encoder->encode($node);
    }

    public static function DecodeWire(array $data): Node
    {
        $decoder = new WireDecoder([]);
        return $decoder->decode($data);
    }

    public static function ValueToJSONC(mixed $value, string $name): string
    {
        $tag = Tag::parseMMTag($name);
        $node = self::valueToNode($value, $tag, 0, '');
        return Jsonc::ToJSONC($node);
    }

    public static function BindFromJSONC(string $in, mixed $out): void
    {
        $n = self::ParseFromJSONC($in);
        self::bind($n, $out);
    }

    public static function GetInt(Node $node, string $path): int
    {
        return 0;
    }

    public static function GetString(Node $node, string $path): string
    {
        return '';
    }

    public static function GetFloat(Node $node, string $path): float
    {
        return 0;
    }

    public static function ParseFromJSONC(string $in): Node
    {
        $scanner = new JsoncScanner($in);
        $toks = [];
        while (true) {
            $t = $scanner->nextToken();
            $toks[] = $t;
            if ($t->type === JsoncTokenType::EOF) {
                break;
            }
        }

        $parser = new JsoncParser($toks);
        return $parser->parse();
    }

    public static function PrintJSONC(Node $n): void
    {
        echo Jsonc::ToJSONC($n) . "\n";
    }

    public static function Dump(Node $n): string
    {
        return json_encode($n, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE);
    }

    public static function encode(object $obj): array
    {
        $encoder = new WireEncoder();
        if ($obj instanceof Node) {
            return $encoder->encode($obj);
        }
        $node = self::valueToNode($obj, null, 0, '');
        return $encoder->encode($node);
    }

    public static function decode(array $wire, string $className): object
    {
        $decoder = new WireDecoder([]);
        $node = $decoder->decode($wire);
        $obj = self::newInstance($className);
        self::bind($node, $obj);
        return $obj;
    }

    private static function newInstance(string $className): object
    {
        $ref = new ReflectionClass($className);
        $constructor = $ref->getConstructor();
        if ($constructor !== null && $constructor->getNumberOfRequiredParameters() === 0) {
            return $ref->newInstance();
        }
        return $ref->newInstanceWithoutConstructor();
    }

    private static function valueToNode(mixed $v, ?Tag $tag, int $depth, string $path, ?\ReflectionNamedType $propType = null): Node
    {
        if ($depth > self::MAX_DEPTH) {
            throw new \RuntimeException('max depth exceeded: ' . self::MAX_DEPTH);
        }

        if ($tag === null) {
            $tag = Tag::newTag();
        }

        if ($v === null) {
            if ($tag->type === ValueType::UNKNOWN && $propType !== null) {
                $typeName = $propType->getName();
                $tag->type = self::phpTypeToValueType($typeName);
            }
            if ($tag->type === ValueType::UNKNOWN) {
                throw new \RuntimeException('invalid input: v is untyped nil (no concrete type/value)');
            }
            if (!$tag->nullable) {
                $tag->nullable = true;
            }
            $tag->isNull = true;
            return new Value(null, Constants::NULL, $tag, $path);
        }

        if (is_object($v)) {
            return self::objectToNode($v, $tag, $depth, $path);
        }

        if (is_array($v)) {
            return self::arrayToNode($v, $tag, $depth, $path);
        }

        return self::scalarToNode($v, $tag, $path);
    }

    private static function phpTypeToValueType(string $typeName): ValueType
    {
        return match ($typeName) {
            'string' => ValueType::STR,
            'int' => ValueType::I,
            'float' => ValueType::F64,
            'bool' => ValueType::BOOL,
            'array' => ValueType::VEC,
            default => ValueType::UNKNOWN,
        };
    }

    private static function objectToNode(object $v, Tag $tag, int $depth, string $path): Node
    {
        if ($v instanceof \DateTime) {
            if ($tag->type === ValueType::UNKNOWN) {
                $tag->type = ValueType::DATETIME;
            }
            switch ($tag->type) {
                case ValueType::DATETIME:
                    $text = $v->format('Y-m-d H:i:s');
                    return new Value($v, $text, $tag, $path);
                case ValueType::DATE:
                    $text = $v->format('Y-m-d');
                    return new Value($v, $text, $tag, $path);
                case ValueType::TIME:
                    $text = $v->format('H:i:s');
                    return new Value($v, $text, $tag, $path);
                default:
                    $text = $v->format('Y-m-d H:i:s');
                    return new Value($v, $text, $tag, $path);
            }
        }

        return self::structToNode($v, $tag, $depth, $path);
    }

    private static function structToNode(object $v, Tag $tag, int $depth, string $path): Node
    {
        $depth++;
        $ref = new ReflectionClass($v);

        $className = $ref->getShortName();
        $tagName = self::camelToSnake($className);

        if ($tag->name === '') {
            $tag->name = $tagName;
        }

        if ($path === '') {
            $path = $tag->name;
        } else {
            $path = $path . '.' . $tag->name;
        }

        $tag->type = ValueType::OBJ;

        $objNode = new Object_();
        $objNode->Tag = $tag;
        $objNode->Path = $path;

        foreach ($ref->getProperties(ReflectionProperty::IS_PUBLIC) as $prop) {
            if ($prop->isStatic()) {
                continue;
            }

            $fieldKey = self::camelToSnake($prop->getName());

            $ann = null;
            $mmAttributes = $prop->getAttributes(MM::class);
            if (!empty($mmAttributes)) {
                $ann = $mmAttributes[0]->newInstance();
            }

            $tagField = null;
            if ($ann !== null) {
                $tagField = Tag::fromAnnotation($ann);
                if ($tagField->name !== '') {
                    if ($tagField->name === '-') {
                        continue;
                    }
                    $fieldKey = $tagField->name;
                }
            }

            if ($tagField === null) {
                $tagField = Tag::newTag();
            }

            if ($tagField->name === '') {
                $tagField->name = $fieldKey;
            }

            $fieldValue = $prop->getValue($v);

            $propType = $prop->getType();
            if ($propType instanceof ReflectionNamedType && $propType->getName() === 'array' && $tagField->type === ValueType::UNKNOWN && $ann !== null && $ann->childType !== ValueType::UNKNOWN) {
                $tagField->childType = $ann->childType;
            }

            $subPath = $path . '.' . $fieldKey;
            $fieldNode = self::valueToNode($fieldValue, $tagField, $depth, $subPath, $propType instanceof \ReflectionNamedType ? $propType : null);

            $field = new Field();
            $field->Key = $fieldKey;
            $field->Value = $fieldNode;
            $objNode->Fields[] = $field;
        }

        return $objNode;
    }

    private static function arrayToNode(array $v, Tag $tag, int $depth, string $path): Node
    {
        $depth++;
        if (array_is_list($v)) {
            $tag->type = ValueType::VEC;
            $node = new Array_();
            $node->Tag = $tag;
            $node->Path = $path;

            foreach ($v as $i => $item) {
                $tagItem = Tag::newTag();
                $tagItem->inheritFromArrayParent($tag);
                $subPath = $path . '[' . $i . ']';
                $itemNode = self::valueToNode($item, $tagItem, $depth, $subPath);
                $node->Items[] = $itemNode;
            }

            return $node;
        }

        $tag->type = ValueType::MAP;
        $node = new Object_();
        $node->Tag = $tag;
        $node->Path = $path;

        foreach ($v as $key => $val) {
            $keyStr = (string)$key;
            $keyStr = self::camelToSnake($keyStr);

            $tagItem = Tag::newTag();
            $tagItem->inherit($tag);
            $tagItem->name = $keyStr;

            $subPath = $path . '[' . $keyStr . ']';
            $valNode = self::valueToNode($val, $tagItem, $depth, $subPath);

            $field = new Field();
            $field->Key = $keyStr;
            $field->Value = $valNode;
            $node->Fields[] = $field;
        }

        return $node;
    }

    private static function scalarToNode(mixed $v, Tag $tag, string $path): Value
    {
        if (is_bool($v)) {
            if ($tag->type === ValueType::UNKNOWN) {
                $tag->type = ValueType::BOOL;
            }
            return new Value($v, $v ? Constants::TRUE : Constants::FALSE, $tag, $path);
        }

        if (is_int($v)) {
            if ($tag->type === ValueType::UNKNOWN) {
                $tag->type = ValueType::I;
            }
            return new Value($v, (string)$v, $tag, $path);
        }

        if (is_float($v)) {
            if ($tag->type === ValueType::UNKNOWN) {
                $tag->type = ValueType::F64;
            }
            return new Value($v, (string)$v, $tag, $path);
        }

        if (is_string($v)) {
            if ($tag->type === ValueType::UNKNOWN) {
                $tag->type = ValueType::STR;
            }
            if ($tag->type === ValueType::ENUM) {
                return new Value((int)$v, $v, $tag, $path);
            }
            return new Value($v, $v, $tag, $path);
        }

        throw new \RuntimeException('unsupported scalar type: ' . gettype($v));
    }

    private static function bind(Node $node, mixed &$out): void
    {
        if ($node instanceof Object_) {
            if ($node->Tag !== null && $node->Tag->type === ValueType::OBJ) {
                self::bindStruct($node, $out);
            } else {
                self::bindMap($node, $out);
            }
        } elseif ($node instanceof Array_) {
            self::bindArray($node, $out);
        } elseif ($node instanceof Value) {
            self::bindValue($node, $out);
        } else {
            throw new \RuntimeException('unsupported node type: ' . get_class($node));
        }
    }

    private static function bindStruct(Object_ $obj, object &$out): void
    {
        $ref = new ReflectionClass($out);

        foreach ($obj->Fields as $field) {
            $fieldKey = $field->Key;
            $propName = self::snakeToCamel($fieldKey);

            if (!$ref->hasProperty($propName)) {
                continue;
            }

            $prop = $ref->getProperty($propName);
            if (!$prop->isPublic()) {
                continue;
            }

            $propValue = $prop->getValue($out);

            if ($field->Value instanceof Object_) {
                if ($propValue === null || !is_object($propValue)) {
                    $propType = $prop->getType();
                    if ($propType instanceof ReflectionNamedType && !$propType->isBuiltin()) {
                        $className = $propType->getName();
                        $propValue = self::newInstance($className);
                    }
                }
                self::bind($field->Value, $propValue);
            } elseif ($field->Value instanceof Array_) {
                self::bind($field->Value, $propValue);
            } elseif ($field->Value instanceof Value) {
                $tag = $field->Value->getTag();
                if ($tag !== null && $tag->type === ValueType::DATETIME && $field->Value->Data instanceof \DateTime) {
                    $propValue = $field->Value->Data;
                } else {
                    $propValue = $field->Value->Data;
                }
                if ($tag !== null && $tag->isNull && $tag->nullable) {
                    $propValue = null;
                }
            } else {
                $propValue = $field->Value;
            }

            $prop->setValue($out, $propValue);
        }
    }

    private static function bindMap(Object_ $obj, object &$out): void {}

    private static function bindArray(Array_ $arr, mixed &$out): void
    {
        if (!is_array($out)) {
            $out = [];
        }

        $result = [];
        foreach ($arr->Items as $item) {
            if ($item instanceof Value) {
                $tag = $item->getTag();
                if ($tag !== null && $tag->isNull && $tag->nullable) {
                    $result[] = null;
                } else {
                    $result[] = $item->Data;
                }
            } elseif ($item instanceof Object_) {
                $itemObj = new \stdClass();
                self::bind($item, $itemObj);
                $result[] = $itemObj;
            } elseif ($item instanceof Array_) {
                $subArr = [];
                self::bind($item, $subArr);
                $result[] = $subArr;
            }
        }

        $out = $result;
    }

    private static function bindValue(Value $val, mixed &$out): void
    {
        $tag = $val->getTag();
        if ($tag === null) {
            $out = $val->Data;
            return;
        }

        if ($tag->isNull && $tag->nullable) {
            $out = null;
            return;
        }

        $data = $val->Data;

        switch ($tag->type) {
            case ValueType::DATETIME:
            case ValueType::DATE:
            case ValueType::TIME:
                if ($data instanceof \DateTime) {
                    $out = $data;
                } else {
                    $out = $val->Text;
                }
                break;

            case ValueType::STR:
            case ValueType::EMAIL:
            case ValueType::URL:
            case ValueType::UUID:
            case ValueType::DECIMAL:
            case ValueType::ENUM:
            case ValueType::IP:
                $out = $val->Text;
                break;

            case ValueType::BYTES:
                if (is_array($data)) {
                    $out = $data;
                } else {
                    $out = $val->Text;
                }
                break;

            case ValueType::BOOL:
                $out = (bool)$data;
                break;

            case ValueType::I:
            case ValueType::I8:
            case ValueType::I16:
            case ValueType::I32:
            case ValueType::I64:
            case ValueType::U:
            case ValueType::U8:
            case ValueType::U16:
            case ValueType::U32:
            case ValueType::U64:
                $out = (int)$data;
                break;

            case ValueType::F32:
            case ValueType::F64:
                $out = (float)$data;
                break;

            default:
                $out = $data;
                break;
        }
    }

    private static function camelToSnake(string $input): string
    {
        $result = '';
        $len = strlen($input);
        for ($i = 0; $i < $len; $i++) {
            $ch = $input[$i];
            if (ctype_upper($ch)) {
                if ($i > 0) {
                    $result .= '_';
                }
                $result .= strtolower($ch);
            } else {
                $result .= $ch;
            }
        }
        return $result;
    }

    private static function snakeToCamel(string $input): string
    {
        $parts = explode('_', $input);
        $result = '';
        foreach ($parts as $i => $part) {
            if ($part !== '') {
                if ($i === 0) {
                    $result .= $part;
                } else {
                    $result .= ucfirst($part);
                }
            }
        }
        return $result;
    }
}
