package io.github.metamessage.jsonc

enum class JsoncTokenType {
    EOF,
    LBrace,
    RBrace,
    LBracket,
    RBracket,
    Colon,
    Comma,
    String,
    Number,
    True,
    False,
    Null,
    LeadingComment,
    TrailingComment
}

data class JsoncToken(
    val type: JsoncTokenType,
    val literal: String = "",
    val line: Int = 0,
    val column: Int = 0
)