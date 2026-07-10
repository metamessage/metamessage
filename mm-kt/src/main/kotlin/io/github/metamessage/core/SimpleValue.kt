package io.github.metamessage.core

object SimpleValue {
    const val SIMPLE_NULL = 0
    const val NULL_BOOL = 1
    const val NULL_INT = 2
    const val NULL_FLOAT = 3
    const val NULL_STRING = 4
    const val NULL_BYTES = 5
    const val FALSE = 6
    const val TRUE = 7
    const val CODE = 8
    const val MESSAGE = 9
    const val DATA = 10
    const val SUCCESS = 11
    const val ERROR = 12
    const val UNKNOWN = 13
    const val PAGE = 14
    const val LIMIT = 15
    const val OFFSET = 16
    const val TOTAL = 17
    const val ID = 18
    const val NAME = 19
    const val DESCRIPTION = 20
    const val TYPE = 21
    const val VERSION = 22
    const val STATUS = 23
    const val URL = 24
    const val CREATE_TIME = 25
    const val UPDATE_TIME = 26
    const val DELETE_TIME = 27
    const val ACCOUNT = 28
    const val TOKEN = 29
    const val EXPIRE_TIME = 30
    const val KEY = 31

    fun toString(v: Int): String =
            when (v) {
                SIMPLE_NULL -> "simple_null"
                NULL_BOOL -> "null_bool"
                NULL_INT -> "null_int"
                NULL_FLOAT -> "null_float"
                NULL_STRING -> "null_string"
                NULL_BYTES -> "null_bytes"
                FALSE -> "false"
                TRUE -> "true"
                CODE -> "code"
                MESSAGE -> "message"
                DATA -> "data"
                SUCCESS -> "success"
                ERROR -> "error"
                UNKNOWN -> "unknown"
                PAGE -> "page"
                LIMIT -> "limit"
                OFFSET -> "offset"
                TOTAL -> "total"
                ID -> "id"
                NAME -> "name"
                DESCRIPTION -> "description"
                TYPE -> "type"
                VERSION -> "version"
                STATUS -> "status"
                URL -> "url"
                CREATE_TIME -> "create_time"
                UPDATE_TIME -> "update_time"
                DELETE_TIME -> "delete_time"
                ACCOUNT -> "account"
                TOKEN -> "token"
                EXPIRE_TIME -> "expire_time"
                KEY -> "key"
                else -> "SimpleValue($v)"
            }

    fun isValid(v: Int): Boolean = v < 32
}
