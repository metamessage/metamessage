package io.metamessage.mm

object Prefix {
    const val SIMPLE = 0b000 shl 5
    const val POSITIVE_INT = 0b001 shl 5
    const val NEGATIVE_INT = 0b010 shl 5
    const val FLOAT = 0b011 shl 5
    const val STRING = 0b100 shl 5
    const val BYTES = 0b101 shl 5
    const val CONTAINER = 0b110 shl 5
    const val TAG = 0b111 shl 5

    const val PREFIX_MASK = 0b11100000
    const val SUFFIX_MASK = 0b00011111

    fun of(b: Byte): Int = b.toInt() and PREFIX_MASK

    fun suffix(b: Byte): Int = b.toInt() and SUFFIX_MASK
}
