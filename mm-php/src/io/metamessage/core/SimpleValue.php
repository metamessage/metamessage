<?php

namespace io\metamessage\core;

class SimpleValue
{
    public const SIMPLE_NULL = 0;
    public const NULL_BOOL = 1;
    public const NULL_INT = 2;
    public const NULL_FLOAT = 3;
    public const NULL_STRING = 4;
    public const NULL_BYTES = 5;
    public const FALSE = 6;
    public const TRUE = 7;

    public const CODE = 8;
    public const MESSAGE = 9;
    public const DATA = 10;
    public const SUCCESS = 11;
    public const ERROR = 12;
    public const UNKNOWN = 13;

    public const PAGE = 14;
    public const LIMIT = 15;
    public const OFFSET = 16;
    public const TOTAL = 17;
    public const ID = 18;
    public const NAME = 19;
    public const DESCRIPTION = 20;
    public const TYPE = 21;
    public const VERSION = 22;
    public const STATUS = 23;
    public const URL = 24;
    public const CREATE_TIME = 25;
    public const UPDATE_TIME = 26;
    public const DELETE_TIME = 27;
    public const ACCOUNT = 28;
    public const TOKEN = 29;
    public const EXPIRE_TIME = 30;
    public const KEY = 31;
    public const VAL = 32;

    public static function nameOf(int $sv): string
    {
        return match ($sv) {
            self::SIMPLE_NULL => 'simple_null',
            self::NULL_BOOL => 'null_bool',
            self::NULL_INT => 'null_int',
            self::NULL_FLOAT => 'null_float',
            self::NULL_STRING => 'null_string',
            self::NULL_BYTES => 'null_bytes',
            self::FALSE => 'false',
            self::TRUE => 'true',
            self::CODE => 'code',
            self::MESSAGE => 'message',
            self::DATA => 'data',
            self::SUCCESS => 'success',
            self::ERROR => 'error',
            self::UNKNOWN => 'unknown',
            self::PAGE => 'page',
            self::LIMIT => 'limit',
            self::OFFSET => 'offset',
            self::TOTAL => 'total',
            self::ID => 'id',
            self::NAME => 'name',
            self::DESCRIPTION => 'description',
            self::TYPE => 'type',
            self::VERSION => 'version',
            self::STATUS => 'status',
            self::URL => 'url',
            self::CREATE_TIME => 'create_time',
            self::UPDATE_TIME => 'update_time',
            self::DELETE_TIME => 'delete_time',
            self::ACCOUNT => 'account',
            self::TOKEN => 'token',
            self::EXPIRE_TIME => 'expire_time',
            self::KEY => 'key',
            self::VAL => 'value',
            default => 'unknown_simple_' . $sv,
        };
    }
}
