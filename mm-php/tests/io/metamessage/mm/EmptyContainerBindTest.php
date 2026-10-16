<?php

use io\metamessage\core\MetaMessage;
use io\metamessage\ir\Field;
use io\metamessage\ir\NodeArray;
use io\metamessage\ir\NodeObject;
use io\metamessage\ir\NodeScalar;
use io\metamessage\ir\Tag;
use io\metamessage\ir\ValueType;
use function io\metamessage\decodeToValue;

// 空容器 example 哨兵与 is_null 绑定行为测试（与 Go zzrepro_test.go 对齐）
class ZzBindTarget
{
    public array $vec = [9];
    public ?string $ptr = 'preset';
}

class EmptyContainerBindTest extends PHPUnit\Framework\TestCase
{
    private static function exampleIntScalar(): NodeScalar
    {
        $tag = Tag::newTag();
        $tag->example = true;
        $tag->type = ValueType::I;
        return new NodeScalar(0, '0', $tag);
    }

    private static function intScalar(int $v): NodeScalar
    {
        $tag = Tag::newTag();
        $tag->type = ValueType::I;
        return new NodeScalar($v, (string)$v, $tag);
    }

    private static function rootWith(string $key, NodeObject|NodeArray|NodeScalar $value): NodeObject
    {
        $rootTag = Tag::newTag();
        $rootTag->type = ValueType::OBJ;
        $root = new NodeObject();
        $root->Tag = $rootTag;
        $field = new Field();
        $field->Key = $key;
        $field->Value = $value;
        $root->Fields[] = $field;
        return $root;
    }

    private static function arr(array $items): NodeArray
    {
        $tag = Tag::newTag();
        $tag->type = ValueType::VEC;
        $arr = new NodeArray();
        $arr->Tag = $tag;
        foreach ($items as $item) {
            $arr->Items[] = $item;
        }
        return $arr;
    }

    public function testEmptyVecWithOnlyExampleBindsToEmptyArray(): void
    {
        $wire = MetaMessage::encode(self::rootWith('vec', self::arr([self::exampleIntScalar()])));

        $obj = MetaMessage::decode($wire, ZzBindTarget::class);

        $this->assertIsArray($obj->vec);
        $this->assertSame([], $obj->vec);
    }

    public function testExampleItemsSkippedAmongRealItems(): void
    {
        $wire = MetaMessage::encode(self::rootWith('vec', self::arr([
            self::intScalar(1),
            self::exampleIntScalar(),
            self::intScalar(2),
        ])));

        $obj = MetaMessage::decode($wire, ZzBindTarget::class);

        $this->assertSame([1, 2], $obj->vec);
    }

    public function testIsNullScalarBindsToNull(): void
    {
        $nullTag = Tag::newTag();
        $nullTag->isNull = true;
        $nullTag->nullable = true;
        $nullTag->type = ValueType::STR;
        $nullScalar = new NodeScalar(null, 'null', $nullTag);

        $rootTag = Tag::newTag();
        $rootTag->type = ValueType::OBJ;
        $root = new NodeObject();
        $root->Tag = $rootTag;
        $field = new Field();
        $field->Key = 'ptr';
        $field->Value = $nullScalar;
        $root->Fields[] = $field;

        $wire = MetaMessage::encode($root);
        $obj = MetaMessage::decode($wire, ZzBindTarget::class);

        $this->assertNull($obj->ptr);
    }

    public function testDecodeToValueSkipsExampleAndMapsNull(): void
    {
        // 动态取值路径：example 哨兵跳过、is_null → null
        $mapTag = Tag::newTag();
        $mapTag->type = ValueType::MAP;
        $map = new NodeObject();
        $map->Tag = $mapTag;
        $kf = new Field();
        $kf->Key = 'k';
        $kf->Value = self::exampleIntScalar();
        $map->Fields[] = $kf;

        $nullTag = Tag::newTag();
        $nullTag->isNull = true;
        $nullTag->nullable = true;
        $nullTag->type = ValueType::STR;
        $nullScalar = new NodeScalar(null, 'null', $nullTag);

        $rootTag = Tag::newTag();
        $rootTag->type = ValueType::MAP;
        $root = new NodeObject();
        $root->Tag = $rootTag;
        foreach ([['vec', self::arr([self::exampleIntScalar()])], ['m', $map], ['ptr', $nullScalar]] as [$key, $node]) {
            $field = new Field();
            $field->Key = $key;
            $field->Value = $node;
            $root->Fields[] = $field;
        }

        $wire = MetaMessage::encode($root);
        $value = decodeToValue($wire);

        $this->assertIsArray($value);
        $this->assertSame([], $value['vec']);
        $this->assertSame([], $value['m']);
        $this->assertNull($value['ptr']);
    }
}
