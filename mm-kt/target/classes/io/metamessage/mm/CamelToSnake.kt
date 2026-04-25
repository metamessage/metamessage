package io.metamessage.mm

object CamelToSnake {
    fun convert(s: String): String {
        if (s.isEmpty()) return ""
        val result = StringBuilder(s.length + 4)
        for (i in s.indices) {
            val c = s[i]
            if (c.isUpperCase()) {
                if (i > 0) {
                    val prev = s[i - 1]
                    val prevUpper = prev.isUpperCase()
                    val nextUpper = i + 1 < s.length && s[i + 1].isUpperCase()
                    if (!prevUpper || (i + 1 < s.length && !nextUpper)) {
                        result.append('_')
                    }
                }
                result.append(c.lowercaseChar())
            } else {
                result.append(c)
            }
        }
        return result.toString()
    }
}
