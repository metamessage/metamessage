package io.metamessage.mm;

final class WireConstants {
    static final long MAX_1 = 0xFFL;
    static final long MAX_2 = 0xFFFFL;
    static final long MAX_3 = 0xFFFFFFL;
    static final long MAX_4 = 0xFFFFFFFFL;
    static final long MAX_5 = 0xFFFFFFFFFFL;
    static final long MAX_6 = 0xFFFFFFFFFFFFL;
    static final long MAX_7 = 0xFFFFFFFFFFFFFFL;
    static final long MAX_8 = 0xFFFFFFFFFFFFFFFFL;

    static final int INT_LEN_MASK = 0b11111;
    static final int INT_LEN_1 = INT_LEN_MASK - 7;
    static final int INT_LEN_2 = INT_LEN_MASK - 6;
    static final int INT_LEN_3 = INT_LEN_MASK - 5;
    static final int INT_LEN_4 = INT_LEN_MASK - 4;
    static final int INT_LEN_5 = INT_LEN_MASK - 3;
    static final int INT_LEN_6 = INT_LEN_MASK - 2;
    static final int INT_LEN_7 = INT_LEN_MASK - 1;
    static final int INT_LEN_8 = INT_LEN_MASK;

    static final int FLOAT_NEG_MASK = 0b10000;
    static final int FLOAT_LEN_MASK = 0b01111;
    static final int FLOAT_LEN_1 = FLOAT_LEN_MASK - 7;
    static final int FLOAT_LEN_2 = FLOAT_LEN_MASK - 6;
    static final int FLOAT_LEN_3 = FLOAT_LEN_MASK - 5;
    static final int FLOAT_LEN_4 = FLOAT_LEN_MASK - 4;
    static final int FLOAT_LEN_5 = FLOAT_LEN_MASK - 3;
    static final int FLOAT_LEN_6 = FLOAT_LEN_MASK - 2;
    static final int FLOAT_LEN_7 = FLOAT_LEN_MASK - 1;
    static final int FLOAT_LEN_8 = FLOAT_LEN_MASK;

    static final int STRING_LEN_MASK = 0b11111;
    static final int STRING_LEN_1 = STRING_LEN_MASK - 1;
    static final int STRING_LEN_2 = STRING_LEN_MASK;

    static final int BYTES_LEN_MASK = 0b11111;
    static final int BYTES_LEN_1 = BYTES_LEN_MASK - 1;
    static final int BYTES_LEN_2 = BYTES_LEN_MASK;

    static final int CONTAINER_MASK = 0b10000;
    static final int CONTAINER_MAP = 0b00000;
    static final int CONTAINER_ARRAY = 0b10000;
    static final int CONTAINER_LEN_MASK = 0b01111;
    static final int CONTAINER_LEN_1 = CONTAINER_LEN_MASK - 1;
    static final int CONTAINER_LEN_2 = CONTAINER_LEN_MASK;

    static final int TAG_LEN_MASK = 0b11111;
    static final int TAG_LEN_1 = TAG_LEN_MASK - 1;
    static final int TAG_LEN_2 = TAG_LEN_MASK;

    private WireConstants() {}
}
