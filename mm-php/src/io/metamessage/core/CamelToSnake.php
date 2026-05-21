<?php

namespace io\metamessage\core;

class CamelToSnake
{
    public static function convert(string $input): string
    {
        if (empty($input)) {
            return '';
        }
        $result = strtolower(preg_replace('/([a-z])([A-Z])/', '$1_$2', $input));
        return $result;
    }
}