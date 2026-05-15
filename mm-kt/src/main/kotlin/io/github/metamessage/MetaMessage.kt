package io.github.metamessage

import io.github.metamessage.jsonc.Jsonc as jc
import io.github.metamessage.core.Encoder
import io.github.metamessage.core.Decoder
import io.github.metamessage.core.Binder
import io.github.metamessage.jsonc.toJsonc

// interface Encoder {
//     fun reset(w: Writer)
//     fun encodeStream(value: Any): Int
// }

// interface Decoder {
//     fun reset(r: Reader)
//     fun decodeStream(value: Any): Int
// }

// fun newEncoder(w: Writer): Encoder = mm.newEncoder(w)
// fun newDecoder(r: Reader): Decoder = mm.newDecoder(r)


 // companion object {
    //     private val encoderPool = LinkedBlockingQueue<WireEncoder>()

    //     private fun getEncoder(): WireEncoder {
    //         return encoderPool.poll() ?: WireEncoder()
    //     }

    //     private fun putEncoder(encoder: WireEncoder) {
    //         encoder.reset()
    //         encoderPool.offer(encoder)
    //     }

    //     @JvmStatic
    //     fun fromValue(v: Any, tag: String): ByteArray {
    //         val node = valueToNode(v, tag)
    //         val encoder = getEncoder()
    //         return try {
    //             encoder.encode(node)
    //         } finally {
    //             putEncoder(encoder)
    //         }
    //     }
    // }

object MetaMessage {

    @JvmStatic
    fun encodeFromValue(value: Any): ByteArray {
        return Encoder.encode(value)
    }

    @JvmStatic
    fun encodeFromJsonc(jsonc: String): ByteArray {
        val node = jc.parseFromString(jsonc) ?: throw RuntimeException("Failed to parse JSONC")
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
}