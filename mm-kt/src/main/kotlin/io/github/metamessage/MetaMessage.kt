package io.github.metamessage

import io.github.metamessage.core.Binder
import io.github.metamessage.core.Decoder
import io.github.metamessage.core.Encoder
import io.github.metamessage.core.Encoder.rootTagForClass
import io.github.metamessage.core.valueToNode
import io.github.metamessage.ir.Node
import io.github.metamessage.jsonc.parseFromJsonc
import io.github.metamessage.jsonc.printJsonc
import io.github.metamessage.jsonc.toJsonc

object MetaMessage {

    @JvmStatic
    fun encodeFromValue(value: Any): ByteArray {
        return Encoder.encode(value)
    }

    @JvmStatic
    fun encodeFromJsonc(jsonc: String): ByteArray {
        val node = parseFromJsonc(jsonc)
        return Encoder.encodeNode(node)
    }

    @JvmStatic
    fun <T> decodeToValue(wire: ByteArray, clazz: Class<T>): T {
        val node = Decoder().decode(wire)
        return Binder.bind(node, clazz)
    }

    @JvmStatic
    fun decodeToJsonc(wire: ByteArray): String {
        val node = Decoder().decode(wire)
        return toJsonc(node)
    }

    @JvmStatic
    fun <T> jsoncToValue(jsonc: String, clazz: Class<T>): T {
        val node = parseFromJsonc(jsonc)
        return Binder.bind(node, clazz)
    }

    @JvmStatic
    fun valueToJsonc(value: Any): String {
        val node = valueToNode(value, rootTagForClass(value.javaClass), "")
        return toJsonc(node)
    }

    @JvmStatic
    fun dump(node: Node): String {
        return io.github.metamessage.ir.dump(node)
    }

    @JvmStatic
    fun printJsonc(node: Node) {
        println(io.github.metamessage.jsonc.printJsonc(node))
    }
}
