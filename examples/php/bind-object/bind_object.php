<?php
require_once __DIR__ . '/../../../mm-php/vendor/autoload.php';

use io\metamessage\jsonc\Jsonc;
use io\metamessage\jsonc\JsoncBinder;

class User {
    public string $name;
    public int $age;
    public bool $active;
    public array $scores;
}

// JSONC 字符串
$jsonc = '{
    // mm: type=str; desc=姓名
    "name": "Alice",
    // mm: type=i; desc=年龄
    "age": 25,
    // mm: type=bool; desc=是否激活
    "active": true,
    // mm: type=array; child_type=i; desc=分数
    "scores": [95, 87, 92]
}';

echo "Input JSONC:\n";
echo $jsonc . "\n";

// 解析 JSONC
$node = Jsonc::parseFromString($jsonc);

// 绑定到 User 对象
$user = new User();
JsoncBinder::bind($node, $user);

echo "\nBound to object:\n";
echo "Name: {$user->name}\n";
echo "Age: {$user->age}\n";
echo "Active: " . ($user->active ? "true" : "false") . "\n";
echo "Scores: " . implode(", ", $user->scores) . "\n";
