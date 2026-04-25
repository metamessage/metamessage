<?php

namespace io\metamessage\mm;

class TypeInference {
    public static function forProperty(
        \ReflectionProperty $property
    ): MmTag {
        $ann = $property->getAttributes(MM::class)[0] ?? null;
        if ($ann) {
            $mm = $ann->newInstance();
            $t = MmTag::fromAnnotation($mm);
            if (self::isListType($property)) {
                if ($t->type === ValueType::UNKNOWN) {
                    $t->type = ValueType::SLICE;
                }
            }
            return $t;
        }
        $t = MmTag::empty();
        self::applyPropertyType($t, $property);
        $type = $property->getType();
        if ($type instanceof \ReflectionNamedType) {
            $genericType = $type->getName();
            if ($genericType === 'array') {
                $t->type = ValueType::SLICE;
                $t->childType = ValueType::STRING;
            }
        }
        return $t;
    }

    public static function valueTypeForComponent(string $type): ValueType {
        $t = MmTag::empty();
        self::applyType($t, $type);
        return $t->type;
    }

    private static function applyPropertyType(MmTag $t, \ReflectionProperty $property): void {
        $type = $property->getType();
        if ($type instanceof \ReflectionNamedType) {
            self::applyType($t, $type->getName());
        } elseif ($type instanceof \ReflectionUnionType) {
            $types = $type->getTypes();
            foreach ($types as $unionType) {
                if (!$unionType->isBuiltin() || $unionType->getName() !== 'null') {
                    self::applyType($t, $unionType->getName());
                    break;
                }
            }
        }
    }

    private static function applyType(MmTag $t, string $type): void {
        switch ($type) {
            case 'string':
                $t->type = ValueType::STRING;
                break;
            case 'bool':
                $t->type = ValueType::BOOL;
                break;
            case 'int':
                $t->type = ValueType::INT;
                break;
            case 'float':
                $t->type = ValueType::FLOAT64;
                break;
            case 'array':
                $t->type = ValueType::SLICE;
                break;
            case 'NULL':
                $t->type = ValueType::STRING;
                break;
            case '\DateTime':
            case 'DateTime':
                $t->type = ValueType::DATETIME;
                break;
            case '\DateTimeImmutable':
            case 'DateTimeImmutable':
                $t->type = ValueType::DATETIME;
                break;
            default:
                $t->type = ValueType::STRUCT;
                break;
        }
    }

    private static function isListType(\ReflectionProperty $property): bool {
        $type = $property->getType();
        if ($type instanceof \ReflectionNamedType) {
            return $type->getName() === 'array';
        }
        return false;
    }
}
