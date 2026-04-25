package io.metamessage.mm;

import java.math.BigInteger;

record FloatParts(boolean negative, byte exponent, long mantissa) {}

final class FloatCodec {
    private FloatCodec() {}

    static FloatParts parseDecimalString(String s) {
        if (s == null || s.isEmpty()) {
            throw new IllegalArgumentException("empty numeric string");
        }
        boolean neg = s.charAt(0) == '-';
        if (neg) {
            s = s.substring(1);
            if (s.isEmpty()) {
                throw new IllegalArgumentException("invalid numeric string: only minus sign");
            }
        }
        String expPart = null;
        int eIdx = indexOfExp(s);
        if (eIdx >= 0) {
            expPart = s.substring(eIdx + 1);
            s = s.substring(0, eIdx);
            if (expPart.isEmpty()) {
                throw new IllegalArgumentException("missing exponent part in scientific notation");
            }
        }
        int dot = s.indexOf('.');
        String intPart;
        String fracPart;
        if (dot < 0) {
            intPart = s.isEmpty() ? "0" : s;
            fracPart = "";
        } else {
            intPart = s.substring(0, dot);
            fracPart = s.substring(dot + 1);
        }
        if (intPart.isEmpty()) {
            intPart = "0";
        }
        long baseExp = -((long) fracPart.length());
        if (expPart != null) {
            baseExp += Long.parseLong(expPart);
        }
        if (baseExp < Byte.MIN_VALUE || baseExp > Byte.MAX_VALUE) {
            throw new IllegalArgumentException("final exponent out of range: " + baseExp);
        }
        String mantissaStr = stripLeadingZeros(intPart + fracPart);
        if (mantissaStr.isEmpty()) {
            mantissaStr = "0";
        }
        BigInteger bi = new BigInteger(mantissaStr);
        if (bi.bitLength() > 64) {
            throw new IllegalArgumentException("mantissa overflow (exceeds uint64 max): " + mantissaStr);
        }
        return new FloatParts(neg, (byte) baseExp, bi.longValue());
    }

    private static int indexOfExp(String s) {
        for (int i = 0; i < s.length(); i++) {
            char ch = s.charAt(i);
            if (ch == 'e' || ch == 'E') {
                return i;
            }
        }
        return -1;
    }

    private static String stripLeadingZeros(String s) {
        int i = 0;
        while (i < s.length() && s.charAt(i) == '0') {
            i++;
        }
        return s.substring(i);
    }

    static String mantissaToDecimal(long mantissa, byte exp) {
        String numStr = Long.toString(mantissa);
        int decimalPos = numStr.length() + exp;
        if (decimalPos <= 0) {
            return "0." + "0".repeat(-decimalPos) + numStr;
        }
        if (decimalPos > 0 && decimalPos < numStr.length()) {
            return numStr.substring(0, decimalPos) + "." + numStr.substring(decimalPos);
        }
        return numStr + "0".repeat(Math.max(0, decimalPos - numStr.length()));
    }
}
