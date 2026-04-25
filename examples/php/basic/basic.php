<?php
require_once __DIR__ . '/../../../mm-php/vendor/autoload.php';

use io\metamessage\mm\MetaMessage;

class Person {
    public string $name = "Ed";
    public int $age = 30;
}

// 创建 Person 对象
$person = new Person();
echo "Original: Name={$person->name}, Age={$person->age}\n";

// 编码到 Wire 格式
$wire = MetaMessage::encode($person);
echo "Encoded: " . bin2hex(implode(array_map('chr', $wire))) . "\n";

// 从 Wire 解码
$decoded = MetaMessage::decode($wire, Person::class);
echo "Decoded: Name={$decoded->name}, Age={$decoded->age}\n";
