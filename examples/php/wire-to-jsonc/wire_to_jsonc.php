<?php
require_once __DIR__ . '/../../../mm-php/vendor/autoload.php';

use io\metamessage\jsonc\Jsonc;
use io\metamessage\mm\MetaMessage;

// JSONC 字符串
$jsonc = '{
    // mm: type=datetime; desc=创建时间
    "create_time": "2026-01-01 00:00:00",
    // mm: type=str; desc=用户名称
    "user_name": "Alice",
    // mm: type=bool; desc=是否激活
    "is_active": true,
    // mm: type=array; child_type=i
    "scores": [95, 87, 92]
}';

echo "Original JSONC:\n";
echo $jsonc . "\n";

// 解析并编码
$node = Jsonc::parseFromString($jsonc);
$wire = MetaMessage::encode($node);

echo "\nEncoded Wire:\n";
echo bin2hex(implode(array_map('chr', $wire))) . "\n";

// 从 Wire 解码到 JSONC
$decodedNode = MetaMessage::decode($wire, 'io\metamessage\jsonc\JcNode');
$outputJsonc = Jsonc::toString($decodedNode);

echo "\nDecoded to JSONC:\n";
echo $outputJsonc . "\n";
