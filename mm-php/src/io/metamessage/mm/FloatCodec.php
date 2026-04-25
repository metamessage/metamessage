<?php

namespace io\metamessage\mm;

class FloatCodec {
    public static function parseDecimalString(string $s): array {
        if (empty($s)) {
            throw new \Exception('Empty numeric string');
        }
        $neg = $s[0] === '-';
        if ($neg) {
            $s = substr($s, 1);
            if (empty($s)) {
                throw new \Exception('Invalid numeric string: only minus sign');
            }
        }
        $expPart = null;
        $eIdx = strpos($s, 'e');
        if ($eIdx === false) {
            $eIdx = strpos($s, 'E');
        }
        if ($eIdx !== false) {
            $expPart = substr($s, $eIdx + 1);
            $s = substr($s, 0, $eIdx);
            if (empty($expPart)) {
                throw new \Exception('Missing exponent part in scientific notation');
            }
        }
        $dot = strpos($s, '.');
        if ($dot === false) {
            $intPart = empty($s) ? '0' : $s;
            $fracPart = '';
        } else {
            $intPart = substr($s, 0, $dot);
            $fracPart = substr($s, $dot + 1);
        }
        $intPartFinal = empty($intPart) ? '0' : $intPart;
        $baseExp = -strlen($fracPart);
        if ($expPart !== null) {
            $baseExp += (int)$expPart;
        }
        if ($baseExp < -128 || $baseExp > 127) {
            throw new \Exception("Final exponent out of range: $baseExp");
        }
        $mantissaStr = ltrim($intPartFinal . $fracPart, '0');
        if (empty($mantissaStr)) {
            $mantissaStr = '0';
        }
        $mantissa = (int)$mantissaStr;
        if (strlen($mantissaStr) > 19) {
            throw new \Exception("Mantissa overflow (exceeds uint64 max): $mantissaStr");
        }
        return [
            'negative' => $neg,
            'exponent' => (int)$baseExp,
            'mantissa' => $mantissa
        ];
    }

    public static function mantissaToDecimal(int $mantissa, int $exp): string {
        $numStr = (string)$mantissa;
        $decimalPos = strlen($numStr) + $exp;
        if ($decimalPos <= 0) {
            return '0.' . str_repeat('0', -$decimalPos) . $numStr;
        } elseif ($decimalPos > 0 && $decimalPos < strlen($numStr)) {
            return substr($numStr, 0, $decimalPos) . '.' . substr($numStr, $decimalPos);
        } else {
            return $numStr . str_repeat('0', max(0, $decimalPos - strlen($numStr)));
        }
    }
}
