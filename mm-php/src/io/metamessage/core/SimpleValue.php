<?php

namespace io\metamessage\core;

class SimpleValue
{
    public const NULL_BOOL = 0;
    public const NULL_INT = 1;
    public const NULL_FLOAT = 2;
    public const NULL_STRING = 3;
    public const NULL_BYTES = 4;
    public const FALSE = 5;
    public const TRUE = 6;

    public const CODE = 7;
    public const MESSAGE = 8;
    public const DATA = 9;
    public const SUCCESS = 10;
    public const ERROR = 11;
    public const UNKNOWN = 12;

    public const PAGE = 13;
    public const LIMIT = 14;
    public const OFFSET = 15;
    public const TOTAL = 16;
    public const ID = 17;
    public const NAME = 18;
    public const DESCRIPTION = 19;
    public const TYPE = 20;
    public const VERSION = 21;
    public const STATUS = 22;
    public const URL = 23;
    public const CREATE_TIME = 24;
    public const UPDATE_TIME = 25;
    public const DELETE_TIME = 26;
    public const ACCOUNT = 27;
    public const TOKEN = 28;
    public const EXPIRE_TIME = 29;
    public const KEY = 30;
    public const VAL = 31;

    public static function nameOf(int $sv): string
    {
        return match ($sv) {
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
