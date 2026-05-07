package io.github.metamessage.mm

object MetaMessage {
    fun encode(root: Any): ByteArray {
        return ReflectMmEncoder.encode(root)
    }

    fun <T> decode(data: ByteArray, clazz: Class<T>): T {
        val tree = WireDecoder(data).decode()
        return ReflectMmBinder.bind(tree, clazz)
    }
}
