package io.github.metamessage.ir

enum class ValueType {
    UNKNOWN,
    DOC,
    VEC,
    ARR,
    OBJ,
    MAP,
    STR,
    BYTES,
    BOOL,
    I,
    I8,
    I16,
    I32,
    I64,
    U,
    U8,
    U16,
    U32,
    U64,
    F32,
    F64,
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

    override fun toString(): String =
            when (this) {
                UNKNOWN -> "unknown"
                DOC -> "doc"
                VEC -> "vec"
                ARR -> "arr"
                OBJ -> "obj"
                MAP -> "map"
                STR -> "str"
                BYTES -> "bytes"
                BOOL -> "bool"
                I -> "i"
                I8 -> "i8"
                I16 -> "i16"
                I32 -> "i32"
                I64 -> "i64"
                U -> "u"
                U8 -> "u8"
                U16 -> "u16"
                U32 -> "u32"
                U64 -> "u64"
                F32 -> "f32"
                F64 -> "f64"
                BIGINT -> "bigint"
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
                map[v.toString()] = v
            }
            map
        }

        fun parseWireName(s: String): ValueType {
            if (s.isEmpty()) return UNKNOWN
            return WIRE[s.trim().lowercase()] ?: UNKNOWN
        }
    }
}
