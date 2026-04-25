namespace MetaMessage.Mm;

public static class MimeWire
{
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
            "application/javascript" => "js",
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
            "js" => "application/javascript",
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