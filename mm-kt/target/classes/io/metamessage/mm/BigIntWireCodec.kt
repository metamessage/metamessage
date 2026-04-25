package io.metamessage.mm

object BigIntWireCodec {
    fun encodeSignedDecimal(s: String?): ByteArray {
        if (s.isNullOrEmpty()) return ByteArray(0)
        val neg = s[0] == '-'
        val body = if (neg) s.substring(1) else s
        val bits = mutableListOf<Int>()
        bits.add(if (neg) 1 else 0)
        var i = 0
        while (i < body.length) {
            val rem = body.length - i
            when {
                rem >= 3 -> {
                    val num = atoi(body.substring(i, i + 3))
                    bits.addAll(toBits(num, 10))
                    i += 3
                }
                rem == 2 -> {
                    val num = atoi(body.substring(i, i + 2))
                    bits.addAll(toBits(num, 7))
                    i += 2
                }
                else -> {
                    val num = atoi(body.substring(i, i + 1))
                    bits.addAll(toBits(num, 4))
                    i += 1
                }
            }
        }
        return bitsToBytes(bits)
    }

    fun decodePositive(data: ByteArray, digitGroups: Int): String {
        val bits = bytesToBits(data)
        if (bits.isEmpty()) return ""
        val numStr = StringBuilder()
        var n = digitGroups
        var idx = 0
        while (n > 0) {
            when {
                n >= 3 && idx + 10 <= bits.size -> {
                    val num = fromBits(bits, idx, 10)
                    idx += 10
                    numStr.append(String.format("%03d", num))
                    n -= 3
                }
                n >= 2 && idx + 7 <= bits.size -> {
                    val num = fromBits(bits, idx, 7)
                    idx += 7
                    numStr.append(String.format("%02d", num))
                    n -= 2
                }
                n >= 1 && idx + 4 <= bits.size -> {
                    val num = fromBits(bits, idx, 4)
                    idx += 4
                    numStr.append(num)
                    n -= 1
                }
                else -> break
            }
        }
        return numStr.toString()
    }

    private fun atoi(s: String): Int {
        var v = 0
        for (i in s.indices) {
            v = v * 10 + (s[i] - '0')
        }
        return v
    }

    private fun toBits(v: Int, n: Int): List<Int> {
        val b = MutableList(n) { 0 }
        for (i in 0 until n) {
            b[n - 1 - i] = (v shr i) and 1
        }
        return b
    }

    private fun fromBits(bits: List<Int>, start: Int, len: Int): Int {
        var v = 0
        for (i in 0 until len) {
            v = (v shl 1) or bits[start + i]
        }
        return v
    }

    private fun bitsToBytes(bits: List<Int>): ByteArray {
        if (bits.isEmpty()) return ByteArray(0)
        var bt: Byte = 0
        var bl = 0
        val out = mutableListOf<Byte>()
        for (b in bits) {
            bt = ((bt.toInt()) shl 1 or b).toByte()
            bl++
            if (bl == 8) {
                out.add(bt)
                bt = 0
                bl = 0
            }
        }
        if (bl > 0) {
            bt = (bt.toInt() shl (8 - bl)).toByte()
            out.add(bt)
        }
        return out.toByteArray()
    }

    private fun bytesToBits(data: ByteArray): List<Int> {
        val bits = mutableListOf<Int>()
        for (bt in data) {
            for (i in 7 downTo 0) {
                bits.add((bt.toInt() shr i) and 1)
            }
        }
        return bits
    }

    fun digitCount(s: String): Int {
        val t = if (s.startsWith("-")) s.substring(1) else s
        return t.length
    }
}
