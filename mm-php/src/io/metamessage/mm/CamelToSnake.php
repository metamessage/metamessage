<?php

namespace io\metamessage\mm;

class CamelToSnake {
    public static function convert(string $s): string {
        if (empty($s)) {
            return '';
        }
        $result = '';
        for ($i = 0; $i < strlen($s); $i++) {
            $c = $s[$i];
            if (ctype_upper($c)) {
                if ($i > 0) {
                    $prev = $s[$i - 1];
                    $prevUpper = ctype_upper($prev);
                    $nextUpper = $i + 1 < strlen($s) && ctype_upper($s[$i + 1]);
                    if (!$prevUpper || ($i + 1 < strlen($s) && !$nextUpper)) {
                        $result .= '_';
                    }
                }
                $result .= strtolower($c);
            } else {
                $result .= $c;
            }
        }
        return $result;
    }
}
