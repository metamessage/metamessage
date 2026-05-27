#!/usr/bin/env php
<?php
/**
 * MetaMessage PHP test harness - parse JSONC file and re-print to JSONC.
 */

if ($argc < 2) {
    fwrite(STDERR, "usage: harness [--encode|--decode] <file.jsonc>\n");
    exit(1);
}

// Resolve mm-php's autoloader
$autoloader = __DIR__ . '/../../../mm-php/vendor/autoload.php';
if (!file_exists($autoloader)) {
    fwrite(STDERR, "error: mm-php vendor/autoload.php not found. Run 'composer install' in mm-php/\n");
    exit(1);
}
require_once $autoloader;

use io\metamessage\jsonc\Jsonc;
use io\metamessage\core\MetaMessage;

if ($argv[1] === '--encode') {
    if ($argc < 3) {
        fwrite(STDERR, "usage: harness --encode <file.jsonc>\n");
        exit(1);
    }
    $input = file_get_contents($argv[2]);
    if ($input === false) {
        fwrite(STDERR, "read error\n");
        exit(1);
    }
    try {
        $wire = MetaMessage::FromJSONC($input);
        echo bin2hex(pack('C*', ...$wire));
    } catch (Exception $e) {
        fwrite(STDERR, "encode error: " . $e->getMessage() . "\n");
        exit(1);
    }
    return;
}

if ($argv[1] === '--decode') {
    $hexStr = trim(stream_get_contents(STDIN));
    try {
        $wire = array_values(unpack('C*', hex2bin($hexStr)));
        $output = MetaMessage::DecodeToJsonc($wire);
        echo $output;
    } catch (Exception $e) {
        fwrite(STDERR, "decode error: " . $e->getMessage() . "\n");
        exit(1);
    }
    return;
}

// Existing behavior
$input = file_get_contents($argv[1]);
if ($input === false) {
    fwrite(STDERR, "read error: cannot open {$argv[1]}\n");
    exit(1);
}

try {
    $node = Jsonc::ParseFromString($input);
    $output = Jsonc::ToJSONC($node);
    echo $output;
} catch (Exception $e) {
    fwrite(STDERR, "parse error: " . $e->getMessage() . "\n");
    exit(1);
}