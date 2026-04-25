package io.metamessage.mm

enum class ValueType {
    UNKNOWN,
    DOC,
    SLICE,
    ARRAY,
    STRUCT,
    MAP,
    STRING,
    BYTES,
    BOOL,
    INT,
    INT8,
    INT16,
    INT32,
    INT64,
    UINT,
    UINT8,
    UINT16,
    UINT32,
    UINT64,
    FLOAT32,
    FLOAT64,
    BIGINT,
    DATETIME,
    DATE,
    TIME,
    UUID,
    DECIMAL,
    IP,
    URL,
    EMAIL,
    ENUM,
    IMAGE,
    VIDEO;

    fun code(): Byte = ordinal.toByte()

    fun wireName(): String = when (this) {
        UNKNOWN -> "unknown"
        DOC -> "doc"
        SLICE -> "slice"
        ARRAY -> "arr"
        STRUCT -> "struct"
        MAP -> "map"
        STRING -> "str"
        BYTES -> "bytes"
        BOOL -> "bool"
        INT -> "i"
        INT8 -> "i8"
        INT16 -> "i16"
        INT32 -> "i32"
        INT64 -> "i64"
        UINT -> "u"
        UINT8 -> "u8"
        UINT16 -> "u16"
        UINT32 -> "u32"
        UINT64 -> "u64"
        FLOAT32 -> "f32"
        FLOAT64 -> "f64"
        BIGINT -> "bi"
        DATETIME -> "datetime"
        DATE -> "date"
        TIME -> "time"
        UUID -> "uuid"
        DECIMAL -> "decimal"
        IP -> "ip"
        URL -> "url"
        EMAIL -> "email"
        ENUM -> "enum"
        IMAGE -> "image"
        VIDEO -> "video"
    }

    companion object {
        fun fromCode(b: Int): ValueType {
            val values = entries
            return if (b < 0 || b >= values.size) UNKNOWN else values[b]
        }

        private val WIRE: Map<String, ValueType> = run {
            val map = mutableMapOf<String, ValueType>()
            for (v in entries) {
                map[v.wireName()] = v
            }
            map
        }

        fun parseWireName(s: String): ValueType {
            if (s.isEmpty()) return UNKNOWN
            return WIRE[s.trim().lowercase()] ?: UNKNOWN
        }
    }
}
