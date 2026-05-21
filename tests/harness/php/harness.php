#!/usr/bin/env php
<?php
/**
 * MetaMessage PHP test harness - parse JSONC file and re-print to JSONC.
 */

if ($argc < 2) {
    fwrite(STDERR, "usage: harness <file.jsonc>\n");
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