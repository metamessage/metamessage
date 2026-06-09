<?php

namespace io\metamessage;

use io\metamessage\core\MetaMessage;
use io\metamessage\core\WireDecoder;
use io\metamessage\ir\Node;
use io\metamessage\ir\NodeArray;
use io\metamessage\ir\NodeObject;
use io\metamessage\ir\NodeScalar;
use function io\metamessage\jsonc\parseJsonc;

/**
 * Convert a PHP value directly to MetaMessage binary format.
 *
 * @param mixed $value The value to encode.
 * @param string $tag Optional mm tag string (e.g. 'name=foo; type=i').
 * @return int[] Array of bytes representing the encoded data.
 */
function encodeFromValue(mixed $value, string $tag = ''): array
{
    return MetaMessage::FromValue($value, $tag);
}

/**
 * Convert a JSONC string to MetaMessage binary format.
 *
 * @param string $jsonc The JSONC string to encode.
 * @return int[] Array of bytes representing the encoded data.
 */
function encodeFromJsonc(string $jsonc): array
{
    return MetaMessage::FromJSONC($jsonc);
}

/**
 * Decode MetaMessage binary format to a plain PHP value.
 *
 * @param int[] $data Array of bytes to decode.
 * @return mixed The decoded PHP value (array for objects/maps, list for arrays, scalar otherwise).
 */
function decodeToValue(array $data): mixed
{
    $decoder = new WireDecoder([]);
    $node = $decoder->decode($data);
    return _nodeData($node);
}

/**
 * Decode MetaMessage binary format to a JSONC string.
 *
 * @param int[] $data Array of bytes to decode.
 * @return string The decoded JSONC string.
 */
function decodeToJsonc(array $data): string
{
    return MetaMessage::DecodeToJsonc($data);
}

/**
 * Convert a PHP value to a JSONC string.
 *
 * @param mixed $value The value to convert.
 * @param string $tag Optional mm tag string.
 * @return string The JSONC string.
 */
function valueToJsonc(mixed $value, string $tag = ''): string
{
    return MetaMessage::ValueToJSONC($value, $tag);
}

/**
 * Convert a JSONC string to a plain PHP value.
 *
 * @param string $jsonc The JSONC string to parse.
 * @return mixed The decoded PHP value.
 */
function jsoncToValue(string $jsonc): mixed
{
    $node = parseJsonc($jsonc);
    return _nodeData($node);
}

/**
 * Convert a Node tree to a plain PHP value.
 *
 * @internal
 */
function _nodeData(Node $node): mixed
{
    if ($node instanceof NodeObject) {
        $result = [];
        foreach ($node->Fields as $field) {
            $result[$field->Key] = _nodeData($field->Value);
        }
        return $result;
    }
    if ($node instanceof NodeArray) {
        $result = [];
        foreach ($node->Items as $item) {
            $result[] = _nodeData($item);
        }
        return $result;
    }
    if ($node instanceof NodeScalar) {
        return $node->Data;
    }
    return null;
}
