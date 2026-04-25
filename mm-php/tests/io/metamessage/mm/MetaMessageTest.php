<?php

use io\metamessage\mm\MM;
use io\metamessage\mm\MetaMessage;
use io\metamessage\mm\ValueType;
use PHPUnit\Framework\TestCase;

#[MM]
class Person {
    public string $name = 'Ada';
    public int $age = 40;
}

#[MM]
class Team {
    public string $teamName = 'core';
    #[MM(childType: ValueType::STRING)]
    public array $members = ['a', 'b'];
}

#[MM]
class Clock {
    #[MM(type: ValueType::DATETIME)]
    public \DateTime $when;

    public function __construct() {
        $this->when = new \DateTime('2024-06-01 12:00:00', new \DateTimeZone('UTC'));
    }
}

#[MM]
class AllTypes {
    public string $stringField = 'test';
    public int $intField = 42;
    public float $floatField = 3.14;
    public bool $boolField = true;
    #[MM(type: ValueType::DATETIME)]
    public \DateTime $datetimeField;
    #[MM(type: ValueType::ENUM, enumValues: 'red|green|blue')]
    public int $enumField = 1;
    public array $arrayField = [1, 2, 3];

    public function __construct() {
        $this->datetimeField = new \DateTime('2024-01-01 00:00:00', new \DateTimeZone('UTC'));
    }
}

#[MM]
class NestedObject {
    public string $name = 'parent';
    public ChildObject $child;

    public function __construct() {
        $this->child = new ChildObject();
    }
}

#[MM]
class ChildObject {
    public string $name = 'child';
    public int $value = 100;
}

#[MM]
class NullableTypes {
    #[MM(nullable: true)]
    public ?string $nullableString = null;
    #[MM(nullable: true)]
    public ?int $nullableInt = null;
}

class MetaMessageTest extends TestCase {

    public function testRoundtripSimpleStruct() {
        $p = new Person();
        $wire = MetaMessage::encode($p);
        $out = MetaMessage::decode($wire, Person::class);
        $this->assertEquals($p->name, $out->name);
        $this->assertEquals($p->age, $out->age);
    }

    public function testRoundtripListField() {
        $t = new Team();
        $wire = MetaMessage::encode($t);
        $out = MetaMessage::decode($wire, Team::class);
        $this->assertEquals($t->teamName, $out->teamName);
        $this->assertEquals($t->members, $out->members);
    }

    public function testRoundtripDateTime() {
        $c = new Clock();
        $wire = MetaMessage::encode($c);
        $out = MetaMessage::decode($wire, Clock::class);
        $this->assertEquals($c->when->format('Y-m-d H:i:s'), $out->when->format('Y-m-d H:i:s'));
    }

    public function testRoundtripAllTypes() {
        $obj = new AllTypes();
        $wire = MetaMessage::encode($obj);
        $out = MetaMessage::decode($wire, AllTypes::class);
        $this->assertEquals($obj->stringField, $out->stringField);
        $this->assertEquals($obj->intField, $out->intField);
        $this->assertEquals($obj->floatField, $out->floatField);
        $this->assertEquals($obj->boolField, $out->boolField);
        $this->assertEquals($obj->datetimeField->format('Y-m-d H:i:s'), $out->datetimeField->format('Y-m-d H:i:s'));
        $this->assertEquals($obj->enumField, $out->enumField);
        $this->assertEquals($obj->arrayField, $out->arrayField);
    }

    public function testRoundtripNestedObject() {
        $obj = new NestedObject();
        $wire = MetaMessage::encode($obj);
        $out = MetaMessage::decode($wire, NestedObject::class);
        $this->assertEquals($obj->name, $out->name);
        $this->assertEquals($obj->child->name, $out->child->name);
        $this->assertEquals($obj->child->value, $out->child->value);
    }

    public function testRoundtripNullableTypes() {
        $obj = new NullableTypes();
        $wire = MetaMessage::encode($obj);
        $out = MetaMessage::decode($wire, NullableTypes::class);
        $this->assertNull($out->nullableString);
        $this->assertNull($out->nullableInt);
    }

    public function testRoundtripNullableTypesWithValues() {
        $obj = new NullableTypes();
        $obj->nullableString = 'test';
        $obj->nullableInt = 42;
        $wire = MetaMessage::encode($obj);
        $out = MetaMessage::decode($wire, NullableTypes::class);
        $this->assertEquals('test', $out->nullableString);
        $this->assertEquals(42, $out->nullableInt);
    }

    public function testEmptyArray() {
        $t = new Team();
        $t->members = [];
        $wire = MetaMessage::encode($t);
        $out = MetaMessage::decode($wire, Team::class);
        $this->assertEquals([], $out->members);
    }

    public function testLargeArray() {
        $t = new Team();
        $t->members = array_fill(0, 10, 'member');
        $wire = MetaMessage::encode($t);
        $out = MetaMessage::decode($wire, Team::class);
        $this->assertEquals(10, count($out->members));
        $this->assertEquals('member', $out->members[0]);
    }
}
