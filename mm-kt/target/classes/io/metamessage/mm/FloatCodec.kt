package io.metamessage.mm

import java.math.BigInteger

object FloatCodec {
    fun parseDecimalString(s: String): FloatParts {
        require(s.isNotEmpty()) { "empty numeric string" }
        var str = s
        val neg = str[0] == '-'
        if (neg) {
            str = str.substring(1)
            require(str.isNotEmpty()) { "invalid numeric string: only minus sign" }
        }
        var expPart: String? = null
        val eIdx = indexOfExp(str)
        if (eIdx >= 0) {
            expPart = str.substring(eIdx + 1)
            str = str.substring(0, eIdx)
            require(expPart.isNotEmpty()) { "missing exponent part in scientific notation" }
        }
        val dot = str.indexOf('.')
        val intPart: String
        val fracPart: String
        if (dot < 0) {
            intPart = if (str.isEmpty()) "0" else str
            fracPart = ""
        } else {
            intPart = str.substring(0, dot)
            fracPart = str.substring(dot + 1)
        }
        val intPartFinal = if (intPart.isEmpty()) "0" else intPart
        var baseExp = -fracPart.length.toLong()
        if (expPart != null) {
            baseExp += expPart.toLong()
        }
        require(baseExp >= Byte.MIN_VALUE.toLong() && baseExp <= Byte.MAX_VALUE.toLong()) {
            "final exponent out of range: $baseExp"
        }
        val mantissaStr = stripLeadingZeros(intPartFinal + fracPart).ifEmpty { "0" }
        val bi = BigInteger(mantissaStr)
        require(bi.bitLength() <= 64) { "mantissa overflow (exceeds uint64 max): $mantissaStr" }
        return FloatParts(neg, baseExp.toByte(), bi.toLong())
    }

    private fun indexOfExp(s: String): Int {
        for (i in s.indices) {
            val ch = s[i]
            if (ch == 'e' || ch == 'E') {
                return i
            }
        }
        return -1
    }

    private fun stripLeadingZeros(s: String): String {
        var i = 0
        while (i < s.length && s[i] == '0') {
            i++
        }
        return s.substring(i)
    }

    fun mantissaToDecimal(mantissa: Long, exp: Byte): String {
        val numStr = mantissa.toString()
        val decimalPos = numStr.length + exp.toInt()
        return if (decimalPos <= 0) {
            "0." + "0".repeat(-decimalPos) + numStr
        } else if (decimalPos > 0 && decimalPos < numStr.length) {
            numStr.substring(0, decimalPos) + "." + numStr.substring(decimalPos)
        } else {
            numStr + "0".repeat(maxOf(0, decimalPos - numStr.length))
        }
    }
}
