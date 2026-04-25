<?php

namespace io\metamessage\mm;

class MetaMessage {
    private function __construct() {}

    /**
     * Encode a PHP object to MM bytes. Fields must be accessible;
     * use {@link MM} attribute on fields or types where the wire
     * type should differ from PHP inference.
     */
    public static function encode($root): array {
        return ReflectMmEncoder::encode($root);
    }

    /**
     * Decode MM bytes into a new instance of $clazz.
     * Only a subset of PHP field types is supported;
     * the tree must be a tagged object at the root.
     */
    public static function decode(array $data, string $clazz) {
        $tree = (new WireDecoder($data))->decode();
        return ReflectMmBinder::bind($tree, $clazz);
    }
}
