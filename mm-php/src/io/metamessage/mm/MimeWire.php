<?php

namespace io\metamessage\mm;

class MimeWire {
    private static array $STR_TO_MIME = [
        'image/jpeg' => 1,
        'image/jpg' => 1,
        'image/png' => 2,
        'image/gif' => 3,
        'image/webp' => 4,
        'image/svg+xml' => 5,
        'image/avif' => 6,
        'image/bmp' => 7,
        'image/x-icon' => 8,
        'image/tiff' => 9,
        'image/heic' => 10,
        'image/heif' => 11,
        'text/plain' => 12,
        'text/html' => 13,
        'text/css' => 14,
        'text/javascript' => 15,
        'application/javascript' => 15,
        'application/json' => 16,
        'text/csv' => 17,
        'text/markdown' => 18,
        'application/pdf' => 19,
        'application/zip' => 20,
        'application/gzip' => 21,
        'application/x-tar' => 22,
        'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' => 23,
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document' => 24,
        'application/octet-stream' => 25,
        'video/mp4' => 26,
        'video/webm' => 27,
        'video/mov' => 28,
        'audio/mpeg' => 29,
        'audio/wav' => 30,
        'audio/flac' => 31,
        'font/woff2' => 32,
        'font/ttf' => 33
    ];

    public static function parse(?string $s): int {
        if (empty($s)) {
            return 0;
        }
        $key = strtolower(trim($s));
        return self::$STR_TO_MIME[$key] ?? 0;
    }

    public static function toString(int $code): string {
        return match($code) {
            1 => 'image/jpeg',
            2 => 'image/png',
            3 => 'image/gif',
            4 => 'image/webp',
            5 => 'image/svg+xml',
            6 => 'image/avif',
            7 => 'image/bmp',
            8 => 'image/x-icon',
            9 => 'image/tiff',
            10 => 'image/heic',
            11 => 'image/heif',
            12 => 'text/plain',
            13 => 'text/html',
            14 => 'text/css',
            15 => 'text/javascript',
            16 => 'application/json',
            17 => 'text/csv',
            18 => 'text/markdown',
            19 => 'application/pdf',
            20 => 'application/zip',
            21 => 'application/gzip',
            22 => 'application/x-tar',
            23 => 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
            24 => 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
            25 => 'application/octet-stream',
            26 => 'video/mp4',
            27 => 'video/webm',
            28 => 'video/mov',
            29 => 'audio/mpeg',
            30 => 'audio/wav',
            31 => 'audio/flac',
            32 => 'font/woff2',
            33 => 'font/ttf',
            default => 'unknown'
        };
    }
}
