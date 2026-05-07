package io.github.metamessage.mm;

public final class Prefix {
    public static final int SIMPLE = 0b000 << 5;
    public static final int POSITIVE_INT = 0b001 << 5;
    public static final int NEGATIVE_INT = 0b010 << 5;
    public static final int FLOAT = 0b011 << 5;
    public static final int STRING = 0b100 << 5;
    public static final int BYTES = 0b101 << 5;
    public static final int CONTAINER = 0b110 << 5;
    public static final int TAG = 0b111 << 5;

    public static final int PREFIX_MASK = 0b11100000;
    public static final int SUFFIX_MASK = 0b00011111;

    private Prefix() {}

    public static int of(byte b) {
        return b & PREFIX_MASK;
    }

    public static int suffix(byte b) {
        return b & SUFFIX_MASK;
    }
}
