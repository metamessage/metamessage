package io.github.metamessage.core

object WireConstants {
    val MAX_1 = 0xFFL
    val MAX_2 = 0xFFFFL
    val MAX_3 = 0xFFFFFFL
    val MAX_4 = 0xFFFFFFFFL
    val MAX_5 = 0xFFFFFFFFFFL
    val MAX_6 = 0xFFFFFFFFFFFFL
    val MAX_7 = 0xFFFFFFFFFFFFFFL
    val MAX_8 = -1L

    const val INT_LEN_MASK = 0b11111
    const val INT_LEN_1 = INT_LEN_MASK - 7
    const val INT_LEN_2 = INT_LEN_MASK - 6
    const val INT_LEN_3 = INT_LEN_MASK - 5
    const val INT_LEN_4 = INT_LEN_MASK - 4
    const val INT_LEN_5 = INT_LEN_MASK - 3
    const val INT_LEN_6 = INT_LEN_MASK - 2
    const val INT_LEN_7 = INT_LEN_MASK - 1
    const val INT_LEN_8 = INT_LEN_MASK

    const val FLOAT_NEG_MASK = 0b10000
    const val FLOAT_LEN_MASK = 0b01111
    const val FLOAT_LEN_1 = FLOAT_LEN_MASK - 7
    const val FLOAT_LEN_2 = FLOAT_LEN_MASK - 6
    const val FLOAT_LEN_3 = FLOAT_LEN_MASK - 5
    const val FLOAT_LEN_4 = FLOAT_LEN_MASK - 4
    const val FLOAT_LEN_5 = FLOAT_LEN_MASK - 3
    const val FLOAT_LEN_6 = FLOAT_LEN_MASK - 2
    const val FLOAT_LEN_7 = FLOAT_LEN_MASK - 1
    const val FLOAT_LEN_8 = FLOAT_LEN_MASK

    const val STRING_LEN_MASK = 0b11111
    const val STRING_LEN_1 = STRING_LEN_MASK - 1
    const val STRING_LEN_2 = STRING_LEN_MASK

    // Bytes length（v0.3.0：4 档；内联 0..27，档位标记 28/29/30/31）
    const val BYTES_LEN_MASK = 0b11111
    const val BYTES_LEN_1 = BYTES_LEN_MASK - 3
    const val BYTES_LEN_2 = BYTES_LEN_MASK - 2
    const val BYTES_LEN_3 = BYTES_LEN_MASK - 1
    const val BYTES_LEN_4 = BYTES_LEN_MASK

    const val CONTAINER_MASK = 0b10000
    const val CONTAINER_MAP = 0b00000
    const val CONTAINER_ARRAY = 0b10000
    // Container length（v0.3.0：4 档；内联 0..11，档位标记 12/13/14/15）
    const val CONTAINER_LEN_MASK = 0b01111
    const val CONTAINER_LEN_1 = CONTAINER_LEN_MASK - 3
    const val CONTAINER_LEN_2 = CONTAINER_LEN_MASK - 2
    const val CONTAINER_LEN_3 = CONTAINER_LEN_MASK - 1
    const val CONTAINER_LEN_4 = CONTAINER_LEN_MASK

    // Tag length（v0.3.0：4 档；tag+payload 总槽；内联 0..27）
    const val TAG_LEN_MASK = 0b11111
    const val TAG_LEN_1 = TAG_LEN_MASK - 3
    const val TAG_LEN_2 = TAG_LEN_MASK - 2
    const val TAG_LEN_3 = TAG_LEN_MASK - 1
    const val TAG_LEN_4 = TAG_LEN_MASK
}
