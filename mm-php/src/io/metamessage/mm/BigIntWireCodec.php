<?php

namespace io\metamessage\mm;

class BigIntWireCodec {
    public static function encodeSignedDecimal(string $s): array {
        if (empty($s)) {
            return [];
        }
        $neg = $s[0] === '-';
        $body = $neg ? substr($s, 1) : $s;
        $bits = [$neg ? 1 : 0];
        $i = 0;
        while ($i < strlen($body)) {
            $rem = strlen($body) - $i;
            if ($rem >= 3) {
                $num = (int)substr($body, $i, 3);
                $bits = array_merge($bits, self::toBits($num, 10));
                $i += 3;
            } elseif ($rem == 2) {
                $num = (int)substr($body, $i, 2);
                $bits = array_merge($bits, self::toBits($num, 7));
                $i += 2;
            } else {
                $num = (int)substr($body, $i, 1);
                $bits = array_merge($bits, self::toBits($num, 4));
                $i += 1;
            }
        }
        return self::bitsToBytes($bits);
    }

    public static function decodePositive(array $data, int $digitGroups): string {
        $bits = self::bytesToBits($data);
        if (empty($bits)) {
            return '';
        }
        $numStr = '';
        $n = $digitGroups;
        $idx = 0;
        while ($n > 0) {
            if ($n >= 3 && $idx + 10 <= count($bits)) {
                $num = self::fromBits($bits, $idx, 10);
                $idx += 10;
                $numStr .= str_pad((string)$num, 3, '0', STR_PAD_LEFT);
                $n -= 3;
            } elseif ($n >= 2 && $idx + 7 <= count($bits)) {
                $num = self::fromBits($bits, $idx, 7);
                $idx += 7;
                $numStr .= str_pad((string)$num, 2, '0', STR_PAD_LEFT);
                $n -= 2;
            } elseif ($n >= 1 && $idx + 4 <= count($bits)) {
                $num = self::fromBits($bits, $idx, 4);
                $idx += 4;
                $numStr .= (string)$num;
                $n -= 1;
            } else {
                break;
            }
        }
        return $numStr;
    }

    private static function toBits(int $v, int $n): array {
        $b = array_fill(0, $n, 0);
        for ($i = 0; $i < $n; $i++) {
            $b[$n - 1 - $i] = ($v >> $i) & 1;
        }
        return $b;
    }

    private static function fromBits(array $bits, int $start, int $len): int {
        $v = 0;
        for ($i = 0; $i < $len; $i++) {
            $v = ($v << 1) | $bits[$start + $i];
        }
        return $v;
    }

    private static function bitsToBytes(array $bits): array {
        if (empty($bits)) {
            return [];
        }
        $bt = 0;
        $bl = 0;
        $out = [];
        foreach ($bits as $b) {
            $bt = ($bt << 1) | $b;
            $bl++;
            if ($bl == 8) {
                $out[] = $bt;
                $bt = 0;
                $bl = 0;
            }
        }
        if ($bl > 0) {
            $bt = $bt << (8 - $bl);
            $out[] = $bt;
        }
        return $out;
    }

    private static function bytesToBits(array $data): array {
        $bits = [];
        foreach ($data as $bt) {
            for ($i = 7; $i >= 0; $i--) {
                $bits[] = ($bt >> $i) & 1;
            }
        }
        return $bits;
    }

    public static function digitCount(string $s): int {
        $t = $s[0] === '-' ? substr($s, 1) : $s;
        return strlen($t);
    }
}
