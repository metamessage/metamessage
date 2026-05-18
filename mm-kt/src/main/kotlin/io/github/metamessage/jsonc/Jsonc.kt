package io.github.metamessage.jsonc

import io.github.metamessage.ir.Node

class JsoncException(message: String) : Exception(message)

fun parseFromJsonc(jsonc: String): Node {
    val scanner = JsoncScanner(jsonc)
    val tokens = mutableListOf<JsoncToken>()
    while (true) {
        val t = scanner.nextToken()
        tokens.add(t)
        if (t.type == JsoncTokenType.EOF) break
    }
    val parser = JsoncParser(tokens)
    return parser.parse() ?: throw JsoncException("failed to parse JSONC")
}

fun printJsonc(node: Node?): String {
    if (node == null) return ""
    val sb = StringBuilder()
    writeLeadingComments(sb, node.tag, 0)
    sb.append(JsoncPrinter.toString(node))
    return sb.toString()
}

private fun writeLeadingComments(
        sb: StringBuilder,
        tag: io.github.metamessage.ir.Tag?,
        indentLevel: Int
) {
    val tagStr = tag?.toString() ?: ""
    if (tagStr.isNotEmpty()) {
        sb.append("// mm: $tagStr\n")
    }
}

fun toJsonc(node: Node?): String {
    return printJsonc(node)
}
