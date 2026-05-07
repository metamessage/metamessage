package io.github.metamessage.mm;

/**
 * Public entry points for the Java MM codec (binary compatible with the Go {@code mm} package).
 * Field metadata is expressed with {@link MM} instead of Go struct tags.
 */
public final class MetaMessage {

    private MetaMessage() {}

    /**
     * Encode a Java bean graph to MM bytes. Fields must be accessible (typically public or
     * package-private with same-module access); use {@link MM} on fields or types where the wire
     * type should differ from Java inference.
     */
    public static byte[] encode(Object root) throws ReflectiveOperationException {
        return ReflectMmEncoder.encode(root);
    }

    /**
     * Decode MM bytes into a new instance of {@code clazz}. Only a subset of Java field types is
     * supported; the tree must be a tagged object at the root.
     */
    public static <T> T decode(byte[] data, Class<T> clazz) throws Exception {
        MmTree tree = new WireDecoder(data).decode();
        return ReflectMmBinder.bind(tree, clazz);
    }
}
