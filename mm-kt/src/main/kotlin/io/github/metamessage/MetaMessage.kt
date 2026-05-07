package io.github.metamessage

object MetaMessage {
    @JvmStatic
    fun encode(value: Any): ByteArray {
        // 内部调用 jsonc / mm 包
    }

    @JvmStatic
    fun <T> decode(wire: ByteArray, clazz: Class<T>): T {
    }

    @JvmStatic
    fun encodeFromJsonc(jsonc: String): ByteArray {
    }

    @JvmStatic
    fun decodeToJsonc(wire: ByteArray): String {
    }
}