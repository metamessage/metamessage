package io.github.metamessage.core

object SnakeToCamel {
    fun convert(s: String): String {
        if (s.isEmpty()) return ""
        val result = StringBuilder(s.length)
        var capitalizeNext = false
        for (i in s.indices) {
            val c = s[i]
            if (c == '_') {
                capitalizeNext = true
            } else if (capitalizeNext) {
                result.append(c.uppercaseChar())
                capitalizeNext = false
            } else {
                result.append(c)
            }
        }
        return result.toString()
    }
}
