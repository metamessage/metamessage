package io.github.metamessage.jsonc;

public record JsoncToken(JsoncTokenType type, String literal, int line, int column) {}
