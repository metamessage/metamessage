package io.github.metamessage.mm;

import java.util.ArrayList;
import java.util.List;

final class BigIntWireCodec {
    private BigIntWireCodec() {}

    static byte[] encodeSignedDecimal(String s) {
        if (s == null || s.isEmpty()) {
            return new byte[0];
        }
        boolean neg = s.charAt(0) == '-';
        String body = neg ? s.substring(1) : s;
        List<Integer> bits = new ArrayList<>();
        bits.add(neg ? 1 : 0);
        int n = body.length();
        int i = 0;
        for (; i < n; i += 3) {
            int rem = n - i;
            if (rem >= 3) {
                int num = atoi(body.substring(i, i + 3));
                bits.addAll(toBits(num, 10));
            } else if (rem == 2) {
                int num = atoi(body.substring(i, i + 2));
                bits.addAll(toBits(num, 7));
            } else {
                int num = atoi(body.substring(i, i + 1));
                bits.addAll(toBits(num, 4));
            }
        }
        return bitsToBytes(bits);
    }

    static String decodePositive(byte[] data, int digitGroups) {
        List<Integer> bits = bytesToBits(data);
        if (bits.isEmpty()) {
            return "";
        }
        StringBuilder numStr = new StringBuilder();
        int n = digitGroups;
        int idx = 0;
        while (n > 0) {
            if (n >= 3 && idx + 10 <= bits.size()) {
                int num = fromBits(bits, idx, 10);
                idx += 10;
                numStr.append(String.format("%03d", num));
                n -= 3;
            } else if (n >= 2 && idx + 7 <= bits.size()) {
                int num = fromBits(bits, idx, 7);
                idx += 7;
                numStr.append(String.format("%02d", num));
                n -= 2;
            } else if (n >= 1 && idx + 4 <= bits.size()) {
                int num = fromBits(bits, idx, 4);
                idx += 4;
                numStr.append(num);
                n -= 1;
            } else {
                break;
            }
        }
        return numStr.toString();
    }

    private static int atoi(String s) {
        int v = 0;
        for (int i = 0; i < s.length(); i++) {
            v = v * 10 + (s.charAt(i) - '0');
        }
        return v;
    }

    private static List<Integer> toBits(int v, int n) {
        List<Integer> b = new ArrayList<>(n);
        for (int i = 0; i < n; i++) {
            b.add(0);
        }
        for (int i = 0; i < n; i++) {
            b.set(n - 1 - i, (v >> i) & 1);
        }
        return b;
    }

    private static int fromBits(List<Integer> bits, int start, int len) {
        int v = 0;
        for (int i = 0; i < len; i++) {
            v = (v << 1) | bits.get(start + i);
        }
        return v;
    }

    private static byte[] bitsToBytes(List<Integer> bits) {
        if (bits.isEmpty()) {
            return new byte[0];
        }
        byte bt = 0;
        int bl = 0;
        List<Byte> out = new ArrayList<>();
        for (int b : bits) {
            bt = (byte) ((bt << 1) | b);
            bl++;
            if (bl == 8) {
                out.add(bt);
                bt = 0;
                bl = 0;
            }
        }
        if (bl > 0) {
            bt = (byte) (bt << (8 - bl));
            out.add(bt);
        }
        byte[] arr = new byte[out.size()];
        for (int i = 0; i < out.size(); i++) {
            arr[i] = out.get(i);
        }
        return arr;
    }

    private static List<Integer> bytesToBits(byte[] data) {
        List<Integer> bits = new ArrayList<>(data.length * 8);
        for (byte bt : data) {
            for (int i = 7; i >= 0; i--) {
                bits.add((bt >> i) & 1);
            }
        }
        return bits;
    }

    static int digitCount(String s) {
        String t = s.startsWith("-") ? s.substring(1) : s;
        return t.length();
    }
}
