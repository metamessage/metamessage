package io.github.metamessage.mm;

import java.util.HashMap;
import java.util.Locale;
import java.util.Map;

final class MimeWire {
    private static final Map<String, Integer> STR_TO_MIME = new HashMap<>();

    static {
        put("image/jpeg", 1);
        put("image/jpg", 1);
        put("image/png", 2);
        put("image/gif", 3);
        put("image/webp", 4);
        put("image/svg+xml", 5);
        put("image/avif", 6);
        put("image/bmp", 7);
        put("image/x-icon", 8);
        put("image/tiff", 9);
        put("image/heic", 10);
        put("image/heif", 11);
        put("text/plain", 12);
        put("text/html", 13);
        put("text/css", 14);
        put("text/javascript", 15);
        put("application/javascript", 15);
        put("application/json", 16);
        put("text/csv", 17);
        put("text/markdown", 18);
        put("application/pdf", 19);
        put("application/zip", 20);
        put("application/gzip", 21);
        put("application/x-tar", 22);
        put("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", 23);
        put("application/vnd.openxmlformats-officedocument.wordprocessingml.document", 24);
        put("application/octet-stream", 25);
        put("video/mp4", 26);
        put("video/webm", 27);
        put("video/mov", 28);
        put("audio/mpeg", 29);
        put("audio/wav", 30);
        put("audio/flac", 31);
        put("font/woff2", 32);
        put("font/ttf", 33);
    }

    private static void put(String k, int v) {
        STR_TO_MIME.put(k.toLowerCase(Locale.ROOT), v);
    }

    static int parse(String s) {
        if (s == null || s.isBlank()) {
            return 0;
        }
        Integer m = STR_TO_MIME.get(s.trim().toLowerCase(Locale.ROOT));
        return m != null ? m : 0;
    }

    static String toString(int code) {
        return switch (code) {
            case 1 -> "image/jpeg";
            case 2 -> "image/png";
            case 3 -> "image/gif";
            case 4 -> "image/webp";
            case 5 -> "image/svg+xml";
            case 6 -> "image/avif";
            case 7 -> "image/bmp";
            case 8 -> "image/x-icon";
            case 9 -> "image/tiff";
            case 10 -> "image/heic";
            case 11 -> "image/heif";
            case 12 -> "text/plain";
            case 13 -> "text/html";
            case 14 -> "text/css";
            case 15 -> "text/javascript";
            case 16 -> "application/json";
            case 17 -> "text/csv";
            case 18 -> "text/markdown";
            case 19 -> "application/pdf";
            case 20 -> "application/zip";
            case 21 -> "application/gzip";
            case 22 -> "application/x-tar";
            case 23 -> "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet";
            case 24 -> "application/vnd.openxmlformats-officedocument.wordprocessingml.document";
            case 25 -> "application/octet-stream";
            case 26 -> "video/mp4";
            case 27 -> "video/webm";
            case 28 -> "video/mov";
            case 29 -> "audio/mpeg";
            case 30 -> "audio/wav";
            case 31 -> "audio/flac";
            case 32 -> "font/woff2";
            case 33 -> "font/ttf";
            default -> "unknown";
        };
    }

    private MimeWire() {}
}
