from typing import Optional
from enum import IntEnum


class MIME(IntEnum):
    Unknown = 0

    Jpeg = 1
    Png = 2
    Gif = 3
    Webp = 4
    Svg = 5
    Avif = 6
    Bmp = 7
    Ico = 8
    Tiff = 9
    Heic = 10
    Heif = 11

    TextPlain = 12
    Html = 13
    Css = 14
    JavaScript = 15
    Json = 16
    Csv = 17
    Markdown = 18

    Pdf = 19
    Zip = 20
    Gzip = 21
    Tar = 22
    Xlsx = 23
    Docx = 24
    OctetStream = 25

    Mp4 = 26
    Webm = 27
    Mov = 28

    Mp3 = 29
    Wav = 30
    Flac = 31

    Woff2 = 32
    Ttf = 33

    def __str__(self) -> str:
        mapping = {
            MIME.Jpeg: "image/jpeg",
            MIME.Png: "image/png",
            MIME.Gif: "image/gif",
            MIME.Webp: "image/webp",
            MIME.Svg: "image/svg+xml",
            MIME.Avif: "image/avif",
            MIME.Bmp: "image/bmp",
            MIME.Ico: "image/x-icon",
            MIME.Tiff: "image/tiff",
            MIME.Heic: "image/heic",
            MIME.Heif: "image/heif",
            MIME.TextPlain: "text/plain",
            MIME.Html: "text/html",
            MIME.Css: "text/css",
            MIME.JavaScript: "text/javascript",
            MIME.Json: "application/json",
            MIME.Csv: "text/csv",
            MIME.Markdown: "text/markdown",
            MIME.Pdf: "application/pdf",
            MIME.Zip: "application/zip",
            MIME.Gzip: "application/gzip",
            MIME.Tar: "application/x-tar",
            MIME.Xlsx: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
            MIME.Docx: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
            MIME.OctetStream: "application/octet-stream",
            MIME.Mp4: "video/mp4",
            MIME.Webm: "video/webm",
            MIME.Mov: "video/mov",
            MIME.Mp3: "audio/mpeg",
            MIME.Wav: "audio/wav",
            MIME.Flac: "audio/flac",
            MIME.Woff2: "font/woff2",
            MIME.Ttf: "font/ttf",
        }
        return mapping.get(self, "unknown")

    def is_image(self) -> bool:
        return self in (
            MIME.Jpeg, MIME.Png, MIME.Gif, MIME.Webp, MIME.Svg, MIME.Avif,
            MIME.Bmp, MIME.Ico, MIME.Tiff, MIME.Heic, MIME.Heif,
        )

    def ext(self) -> str:
        mapping = {
            MIME.Jpeg: "jpeg",
            MIME.Png: "png",
            MIME.Gif: "gif",
            MIME.Webp: "webp",
            MIME.Svg: "svg",
            MIME.Avif: "avif",
            MIME.Bmp: "bmp",
            MIME.Ico: "ico",
            MIME.Tiff: "tiff",
            MIME.Heic: "heic",
            MIME.Heif: "heif",
            MIME.TextPlain: "txt",
            MIME.Html: "html",
            MIME.Css: "css",
            MIME.JavaScript: "js",
            MIME.Json: "json",
            MIME.Csv: "csv",
            MIME.Markdown: "md",
            MIME.Pdf: "pdf",
            MIME.Zip: "zip",
            MIME.Gzip: "gz",
            MIME.Tar: "tar",
            MIME.Xlsx: "xlsx",
            MIME.Docx: "docx",
            MIME.Mp4: "mp4",
            MIME.Webm: "webm",
            MIME.Mov: "mov",
            MIME.Mp3: "mp3",
            MIME.Wav: "wav",
            MIME.Flac: "flac",
            MIME.Woff2: "woff2",
            MIME.Ttf: "ttf",
        }
        return mapping.get(self, "")


_str_to_mime = {
    "image/jpeg": MIME.Jpeg,
    "image/jpg": MIME.Jpeg,
    "image/png": MIME.Png,
    "image/gif": MIME.Gif,
    "image/webp": MIME.Webp,
    "image/svg+xml": MIME.Svg,
    "image/avif": MIME.Avif,
    "image/bmp": MIME.Bmp,
    "image/x-icon": MIME.Ico,
    "image/tiff": MIME.Tiff,
    "image/heic": MIME.Heic,
    "image/heif": MIME.Heif,
    "text/plain": MIME.TextPlain,
    "text/html": MIME.Html,
    "text/css": MIME.Css,
    "text/javascript": MIME.JavaScript,
    "application/javascript": MIME.JavaScript,
    "application/json": MIME.Json,
    "text/csv": MIME.Csv,
    "text/markdown": MIME.Markdown,
    "application/pdf": MIME.Pdf,
    "application/zip": MIME.Zip,
    "application/gzip": MIME.Gzip,
    "application/x-tar": MIME.Tar,
    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": MIME.Xlsx,
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document": MIME.Docx,
    "application/octet-stream": MIME.OctetStream,
    "video/mp4": MIME.Mp4,
    "video/webm": MIME.Webm,
    "video/mov": MIME.Mov,
    "audio/mpeg": MIME.Mp3,
    "audio/wav": MIME.Wav,
    "audio/flac": MIME.Flac,
    "font/woff2": MIME.Woff2,
    "font/ttf": MIME.Ttf,
}


def ParseMIME(s: str) -> int:
    s = s.lower().strip()
    m = _str_to_mime.get(s, MIME.Unknown)
    return int(m)