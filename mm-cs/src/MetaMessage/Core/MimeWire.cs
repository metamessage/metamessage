namespace MetaMessage.Core;

public static class MimeWire
{
    private static readonly Dictionary<string, int> MimeToIndex = new(StringComparer.OrdinalIgnoreCase)
    {
        {"image/jpeg", 1},
        {"image/jpg", 1},
        {"image/png", 2},
        {"image/gif", 3},
        {"image/webp", 4},
        {"image/svg+xml", 5},
        {"image/avif", 6},
        {"image/bmp", 7},
        {"image/x-icon", 8},
        {"image/tiff", 9},
        {"image/heic", 10},
        {"image/heif", 11},
        {"text/plain", 12},
        {"text/html", 13},
        {"text/css", 14},
        {"text/javascript", 15},
        {"application/json", 16},
        {"text/csv", 17},
        {"text/markdown", 18},
        {"application/pdf", 19},
        {"application/zip", 20},
        {"application/gzip", 21},
        {"application/x-tar", 22},
        {"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", 23},
        {"application/vnd.openxmlformats-officedocument.wordprocessingml.document", 24},
        {"application/octet-stream", 25},
        {"video/mp4", 26},
        {"video/webm", 27},
        {"video/mov", 28},
        {"audio/mpeg", 29},
        {"audio/wav", 30},
        {"audio/flac", 31},
        {"font/woff2", 32},
        {"font/ttf", 33},
    };

    public static int ParseMIME(string s)
    {
        s = s.Trim().ToLowerInvariant();
        if (MimeToIndex.TryGetValue(s, out int index))
        {
            return index;
        }
        return 0;
    }
    public static string MimeForWire(string mimeType)
    {
        return mimeType switch
        {
            "application/json" => "json",
            "application/xml" => "xml",
            "text/plain" => "txt",
            "text/html" => "html",
            "text/css" => "css",
            "text/javascript" => "js",
            "image/jpeg" => "jpg",
            "image/png" => "png",
            "image/gif" => "gif",
            "image/webp" => "webp",
            "audio/mpeg" => "mp3",
            "audio/wav" => "wav",
            "video/mp4" => "mp4",
            "video/webm" => "webm",
            "application/pdf" => "pdf",
            "application/zip" => "zip",
            "application/gzip" => "gz",
            "application/tar" => "tar",
            "application/7z" => "7z",
            "text/csv" => "csv",
            "application/vnd.ms-excel" => "xls",
            "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" => "xlsx",
            "application/msword" => "doc",
            "application/vnd.openxmlformats-officedocument.wordprocessingml.document" => "docx",
            "application/vnd.ms-powerpoint" => "ppt",
            "application/vnd.openxmlformats-officedocument.presentationml.presentation" => "pptx",
            _ => mimeType
        };
    }

    public static string WireForMime(string wireMime)
    {
        return wireMime switch
        {
            "json" => "application/json",
            "xml" => "application/xml",
            "txt" => "text/plain",
            "html" => "text/html",
            "css" => "text/css",
            "js" => "text/javascript",
            "jpg" => "image/jpeg",
            "png" => "image/png",
            "gif" => "image/gif",
            "webp" => "image/webp",
            "mp3" => "audio/mpeg",
            "wav" => "audio/wav",
            "mp4" => "video/mp4",
            "webm" => "video/webm",
            "pdf" => "application/pdf",
            "zip" => "application/zip",
            "gz" => "application/gzip",
            "tar" => "application/tar",
            "7z" => "application/7z",
            "csv" => "text/csv",
            "xls" => "application/vnd.ms-excel",
            "xlsx" => "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
            "doc" => "application/msword",
            "docx" => "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            "ppt" => "application/vnd.ms-powerpoint",
            "pptx" => "application/vnd.openxmlformats-officedocument.presentationml.presentation",
            _ => wireMime
        };
    }
}