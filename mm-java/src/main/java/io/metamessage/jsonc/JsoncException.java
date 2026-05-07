package io.github.metamessage.jsonc;

public final class JsoncException extends RuntimeException {
    public JsoncException(String message) {
        super(message);
    }

    public JsoncException(String message, Throwable cause) {
        super(message, cause);
    }
}
