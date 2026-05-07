package io.github.metamessage.jsonc;

import java.nio.charset.StandardCharsets;
import java.util.List;

/**
 * JSONC entry points (parse, print, bind) mirroring the Go {@code internal/jsonc} package.
 */
public final class Jsonc {
    private Jsonc() {}

    public static JcNode parseFromString(String s) {
        if (s == null) {
            return null;
        }
        List<JsoncToken> toks = new JsoncScanner(s).tokenizeAll();
        return new JsoncParser(toks).parseDocument();
    }

    public static JcNode parseFromBytes(byte[] b) {
        if (b == null) {
            return null;
        }
        return parseFromString(new String(b, StandardCharsets.UTF_8));
    }

    public static String toString(JcNode n) {
        return JsoncPrinter.toString(n);
    }

    public static void bind(JcNode node, Object out) {
        JsoncBinder.bind(node, out);
    }

    public static <T> T bindFromString(String s, Class<T> type) {
        try {
            T inst = type.getDeclaredConstructor().newInstance();
            bind(parseFromString(s), inst);
            return inst;
        } catch (ReflectiveOperationException e) {
            throw new JsoncException("need public no-arg constructor: " + type.getName(), e);
        }
    }

    public static <T> T bindFromBytes(byte[] b, Class<T> type) {
        return bindFromString(b == null ? null : new String(b, StandardCharsets.UTF_8), type);
    }
}
