use std::collections::HashMap;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
#[allow(non_camel_case_types)]
pub enum Mime {
    Unknown = 0,

    Jpeg = 1,
    Png = 2,
    Gif = 3,
    Webp = 4,
    Svg = 5,
    Avif = 6,
    Bmp = 7,
    Ico = 8,
    Tiff = 9,
    Heic = 10,
    Heif = 11,

    TextPlain = 12,
    Html = 13,
    Css = 14,
    JavaScript = 15,
    Json = 16,
    Csv = 17,
    Markdown = 18,

    Pdf = 19,
    Zip = 20,
    Gzip = 21,
    Tar = 22,
    Xlsx = 23,
    Docx = 24,
    OctetStream = 25,

    Mp4 = 26,
    Webm = 27,
    Mov = 28,

    Mp3 = 29,
    Wav = 30,
    Flac = 31,

    Woff2 = 32,
    Ttf = 33,
}

pub fn parse_mime(s: &str) -> u8 {
    let s = s.to_lowercase();
    let s = s.trim();
    let map = mime_str_to_id();
    map.get(s).copied().unwrap_or(Mime::Unknown as u8)
}

pub fn mime_to_str(id: u8) -> &'static str {
    match id {
        1 => "image/jpeg",
        2 => "image/png",
        3 => "image/gif",
        4 => "image/webp",
        5 => "image/svg+xml",
        6 => "image/avif",
        7 => "image/bmp",
        8 => "image/x-icon",
        9 => "image/tiff",
        10 => "image/heic",
        11 => "image/heif",
        12 => "text/plain",
        13 => "text/html",
        14 => "text/css",
        15 => "text/javascript",
        16 => "application/json",
        17 => "text/csv",
        18 => "text/markdown",
        19 => "application/pdf",
        20 => "application/zip",
        21 => "application/gzip",
        22 => "application/x-tar",
        23 => "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        24 => "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        25 => "application/octet-stream",
        26 => "video/mp4",
        27 => "video/webm",
        28 => "video/mov",
        29 => "audio/mpeg",
        30 => "audio/wav",
        31 => "audio/flac",
        32 => "font/woff2",
        33 => "font/ttf",
        _ => "unknown",
    }
}

fn mime_str_to_id() -> HashMap<&'static str, u8> {
    let mut m = HashMap::new();
    m.insert("image/jpeg", Mime::Jpeg as u8);
    m.insert("image/jpg", Mime::Jpeg as u8);
    m.insert("image/png", Mime::Png as u8);
    m.insert("image/gif", Mime::Gif as u8);
    m.insert("image/webp", Mime::Webp as u8);
    m.insert("image/svg+xml", Mime::Svg as u8);
    m.insert("image/avif", Mime::Avif as u8);
    m.insert("image/bmp", Mime::Bmp as u8);
    m.insert("image/x-icon", Mime::Ico as u8);
    m.insert("image/tiff", Mime::Tiff as u8);
    m.insert("image/heic", Mime::Heic as u8);
    m.insert("image/heif", Mime::Heif as u8);
    m.insert("text/plain", Mime::TextPlain as u8);
    m.insert("text/html", Mime::Html as u8);
    m.insert("text/css", Mime::Css as u8);
    m.insert("text/javascript", Mime::JavaScript as u8);
    m.insert("application/json", Mime::Json as u8);
    m.insert("text/csv", Mime::Csv as u8);
    m.insert("text/markdown", Mime::Markdown as u8);
    m.insert("application/pdf", Mime::Pdf as u8);
    m.insert("application/zip", Mime::Zip as u8);
    m.insert("application/gzip", Mime::Gzip as u8);
    m.insert("application/x-tar", Mime::Tar as u8);
    m.insert(
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        Mime::Xlsx as u8,
    );
    m.insert(
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
        Mime::Docx as u8,
    );
    m.insert("application/octet-stream", Mime::OctetStream as u8);
    m.insert("video/mp4", Mime::Mp4 as u8);
    m.insert("video/webm", Mime::Webm as u8);
    m.insert("video/mov", Mime::Mov as u8);
    m.insert("audio/mpeg", Mime::Mp3 as u8);
    m.insert("audio/wav", Mime::Wav as u8);
    m.insert("audio/flac", Mime::Flac as u8);
    m.insert("font/woff2", Mime::Woff2 as u8);
    m.insert("font/ttf", Mime::Ttf as u8);
    m
}