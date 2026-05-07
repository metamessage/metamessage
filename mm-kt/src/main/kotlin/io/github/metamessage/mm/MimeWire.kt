package io.github.metamessage.mm

object MimeWire {
    private val STR_TO_MIME = mutableMapOf<String, Int>().apply {
        put("image/jpeg", 1)
        put("image/jpg", 1)
        put("image/png", 2)
        put("image/gif", 3)
        put("image/webp", 4)
        put("image/svg+xml", 5)
        put("image/avif", 6)
        put("image/bmp", 7)
        put("image/x-icon", 8)
        put("image/tiff", 9)
        put("image/heic", 10)
        put("image/heif", 11)
        put("text/plain", 12)
        put("text/html", 13)
        put("text/css", 14)
        put("text/javascript", 15)
        put("application/javascript", 15)
        put("application/json", 16)
        put("text/csv", 17)
        put("text/markdown", 18)
        put("application/pdf", 19)
        put("application/zip", 20)
        put("application/gzip", 21)
        put("application/x-tar", 22)
        put("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", 23)
        put("application/vnd.openxmlformats-officedocument.wordprocessingml.document", 24)
        put("application/octet-stream", 25)
        put("video/mp4", 26)
        put("video/webm", 27)
        put("video/mov", 28)
        put("audio/mpeg", 29)
        put("audio/wav", 30)
        put("audio/flac", 31)
        put("font/woff2", 32)
        put("font/ttf", 33)
    }

    fun parse(s: String?): Int {
        if (s.isNullOrBlank()) return 0
        return STR_TO_MIME[s.trim().lowercase()] ?: 0
    }

    fun toString(code: Int): String = when (code) {
        1 -> "image/jpeg"
        2 -> "image/png"
        3 -> "image/gif"
        4 -> "image/webp"
        5 -> "image/svg+xml"
        6 -> "image/avif"
        7 -> "image/bmp"
        8 -> "image/x-icon"
        9 -> "image/tiff"
        10 -> "image/heic"
        11 -> "image/heif"
        12 -> "text/plain"
        13 -> "text/html"
        14 -> "text/css"
        15 -> "text/javascript"
        16 -> "application/json"
        17 -> "text/csv"
        18 -> "text/markdown"
        19 -> "application/pdf"
        20 -> "application/zip"
        21 -> "application/gzip"
        22 -> "application/x-tar"
        23 -> "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        24 -> "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
        25 -> "application/octet-stream"
        26 -> "video/mp4"
        27 -> "video/webm"
        28 -> "video/mov"
        29 -> "audio/mpeg"
        30 -> "audio/wav"
        31 -> "audio/flac"
        32 -> "font/woff2"
        33 -> "font/ttf"
        else -> "unknown"
    }
}
