<?php

namespace io\metamessage\mm;

class ReflectMmBinder {
    public static function bind(MmTree $tree, string $clazz) {
        if (!($tree instanceof MmTree\MmObject)) {
            throw new \Exception('Root must be object');
        }
        $inst = new $clazz();
        $byKey = self::fieldsByJsonKey($clazz);
        foreach ($tree->fields as $field) {
            list($key, $value) = $field;
            if (!isset($byKey[$key])) {
                continue;
            }
            $f = $byKey[$key];
            $f->setAccessible(true);
            $f->setValue($inst, self::materialize($f, $value));
        }
        return $inst;
    }

    private static function fieldsByJsonKey(string $className): array {
        $m = [];
        $reflection = new \ReflectionClass($className);
        foreach ($reflection->getProperties() as $f) {
            if ($f->isStatic()) {
                continue;
            }
            $ft = TypeInference::forProperty($f);
            $ann = $f->getAttributes(MM::class)[0] ?? null;
            if ($ann) {
                $mm = $ann->newInstance();
                if ($mm->name === '-') {
                    continue;
                }
            }
            $key = self::fieldKey($f, $ft, $ann);
            $m[$key] = $f;
        }
        return $m;
    }

    private static function fieldKey(\ReflectionProperty $f, MmTag $ft, ?\ReflectionAttribute $ann): string {
        if ($ann) {
            $mm = $ann->newInstance();
            if (!empty($mm->name) && $mm->name !== '-') {
                return $mm->name;
            }
        }
        if (!empty($ft->name)) {
            return $ft->name;
        }
        return CamelToSnake::convert($f->getName());
    }

    private static function materialize(\ReflectionProperty $f, MmTree $node) {
        $ft = $f->getType();
        $typeName = $ft instanceof \ReflectionNamedType ? $ft->getName() : 'mixed';
        switch (get_class($node)) {
            case MmTree\MmScalar::class:
                return self::scalarToField($typeName, $node);
            case MmTree\MmArray::class:
                return self::listFrom($node, $f);
            case MmTree\MmObject::class:
                return self::bind($node, $typeName);
            default:
                return null;
        }
    }

    private static function listFrom(MmTree\MmArray $arr, \ReflectionProperty $f): array {
        $type = $f->getType();
        $elemType = 'string';
        if ($type instanceof \ReflectionNamedType && $type->getName() === 'array') {
            $docComment = $f->getDocComment();
            if (preg_match('/@var\s+array<([^>]+)>/', $docComment, $matches)) {
                $elemType = $matches[1];
            }
        }
        $out = [];
        foreach ($arr->items as $ch) {
            if ($ch instanceof MmTree\MmScalar) {
                $out[] = self::scalarToField($elemType, $ch);
            } elseif ($ch instanceof MmTree\MmObject) {
                $sub = new $elemType();
                $byKey = self::fieldsByJsonKey($elemType);
                foreach ($ch->fields as $field) {
                    list($key, $value) = $field;
                    if (isset($byKey[$key])) {
                        $sf = $byKey[$key];
                        $sf->setAccessible(true);
                        $sf->setValue($sub, self::materialize($sf, $value));
                    }
                }
                $out[] = $sub;
            } else {
                $out[] = null;
            }
        }
        return $out;
    }

    private static function scalarToField(string $ft, MmTree\MmScalar $sc) {
        $d = $sc->data;
        $tag = $sc->tag;
        
        // 处理 nullable 类型的 null 值
        if (str_starts_with($ft, '?') && ($tag->isNull || $d === '' || $d === null || $sc->text === '' || $sc->text === 'null')) {
            return null;
        }
        // 处理非 nullable 类型的 null 值
        if ($tag->isNull) {
            return null;
        }
        
        switch ($ft) {
            case 'int':
            case '?int':
                return is_numeric($d) ? (int)$d : (int)$sc->text;
            case 'float':
            case '?float':
                return is_numeric($d) ? (float)$d : (float)$sc->text;
            case 'bool':
            case '?bool':
                return $d === true;
            case 'string':
            case '?string':
                return (string)$sc->text;
            case '\DateTime':
            case 'DateTime':
            case '?\DateTime':
            case '?DateTime':
                if ($d instanceof \DateTime) {
                    return $d;
                }
                if (is_numeric($sc->text)) {
                    return new \DateTime('@' . $sc->text, new \DateTimeZone('UTC'));
                }
                if (!empty($sc->text) && $sc->text !== 'true' && $sc->text !== 'false') {
                    return new \DateTime($sc->text, new \DateTimeZone('UTC'));
                }
                // 如果 $sc->text 为空或不是有效的日期时间字符串，返回当前时间
                return new \DateTime('now', new \DateTimeZone('UTC'));
            case 'array':
            case '?array':
                return is_array($d) ? $d : [];
            default:
                return $d;
        }
    }
}
