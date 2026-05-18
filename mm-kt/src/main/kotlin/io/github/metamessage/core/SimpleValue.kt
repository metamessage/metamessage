package io.github.metamessage.core

object SimpleValue {
    const val NULL_BOOL = 0
    const val NULL_INT = 1
    const val NULL_FLOAT = 2
    const val NULL_STRING = 3
    const val NULL_BYTES = 4
    const val FALSE = 5
    const val TRUE = 6
    const val CODE = 7
    const val MESSAGE = 8
    const val DATA = 9
    const val SUCCESS = 10
    const val ERROR = 11
    const val UNKNOWN = 12
    const val PAGE = 13
    const val LIMIT = 14
    const val OFFSET = 15
    const val TOTAL = 16
    const val ID = 17
    const val NAME = 18
    const val DESCRIPTION = 19
    const val TYPE = 20
    const val VERSION = 21
    const val STATUS = 22
    const val URL = 23
    const val CREATE_TIME = 24
    const val UPDATE_TIME = 25
    const val DELETE_TIME = 26
    const val ACCOUNT = 27
    const val TOKEN = 28
    const val EXPIRE_TIME = 29
    const val KEY = 30
    const val VAL = 31

    fun toString(v: Int): String = when (v) {
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
        VAL -> "value"
        else -> "SimpleValue($v)"
    }

    fun isValid(v: Int): Boolean = v < 32
}