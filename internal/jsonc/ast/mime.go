package ast

import (
	"fmt"
	"strings"
)

type MIME uint8

const (
	MIMEUnknown MIME = iota

	// Image
	MIMEJpeg
	MIMEPng
	MIMEGif
	MIMEWebp
	MIMESvg
	MIMEAvif
	MIMEBmp
	MIMEIco
	MIMETiff
	MIMEHeic
	MIMEHeif

	// Text
	MIMETextPlain
	MIMEHtml
	MIMECss
	MIMEJavaScript
	MIMEJson
	MIMECsv
	MIMEMarkdown

	// Application
	MIMEPdf
	MIMEZip
	MIMEGzip
	MIMETar
	MIMEXlsx
	MIMEDocx
	MIMEOctetStream

	// Video
	MIMEMp4
	MIMEWebm
	MIMEMov

	// Audio
	MIMEMp3
	MIMEWav
	MIMEFlac

	// Font
	MIMEWoff2
	MIMETtf
)

func (m MIME) String() string {
	switch m {
	// Image
	case MIMEJpeg:
		return "image/jpeg"
	case MIMEPng:
		return "image/png"
	case MIMEGif:
		return "image/gif"
	case MIMEWebp:
		return "image/webp"
	case MIMESvg:
		return "image/svg+xml"
	case MIMEAvif:
		return "image/avif"
	case MIMEBmp:
		return "image/bmp"
	case MIMEIco:
		return "image/x-icon"
	case MIMETiff:
		return "image/tiff"
	case MIMEHeic:
		return "image/heic"
	case MIMEHeif:
		return "image/heif"

	// Text
	case MIMETextPlain:
		return "text/plain"
	case MIMEHtml:
		return "text/html"
	case MIMECss:
		return "text/css"
	case MIMEJavaScript:
		return "text/javascript"
	case MIMEJson:
		return "application/json"
	case MIMECsv:
		return "text/csv"
	case MIMEMarkdown:
		return "text/markdown"

	// Application
	case MIMEPdf:
		return "application/pdf"
	case MIMEZip:
		return "application/zip"
	case MIMEGzip:
		return "application/gzip"
	case MIMETar:
		return "application/x-tar"
	case MIMEXlsx:
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case MIMEDocx:
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case MIMEOctetStream:
		return "application/octet-stream"

	// Video
	case MIMEMp4:
		return "video/mp4"
	case MIMEWebm:
		return "video/webm"
	case MIMEMov:
		return "video/mov"

	// Audio
	case MIMEMp3:
		return "audio/mpeg"
	case MIMEWav:
		return "audio/wav"
	case MIMEFlac:
		return "audio/flac"

	// Font
	case MIMEWoff2:
		return "font/woff2"
	case MIMETtf:
		return "font/ttf"

	default:
		return "unknown"
	}
}

var strToMIME = map[string]MIME{
	"image/jpeg":             MIMEJpeg,
	"image/jpg":              MIMEJpeg,
	"image/png":              MIMEPng,
	"image/gif":              MIMEGif,
	"image/webp":             MIMEWebp,
	"image/svg+xml":          MIMESvg,
	"image/avif":             MIMEAvif,
	"image/bmp":              MIMEBmp,
	"image/x-icon":           MIMEIco,
	"image/tiff":             MIMETiff,
	"image/heic":             MIMEHeic,
	"image/heif":             MIMEHeif,
	"text/plain":             MIMETextPlain,
	"text/html":              MIMEHtml,
	"text/css":               MIMECss,
	"text/javascript":        MIMEJavaScript,
	"application/javascript": MIMEJavaScript,
	"application/json":       MIMEJson,
	"text/csv":               MIMECsv,
	"text/markdown":          MIMEMarkdown,
	"application/pdf":        MIMEPdf,
	"application/zip":        MIMEZip,
	"application/gzip":       MIMEGzip,
	"application/x-tar":      MIMETar,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       MIMEXlsx,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": MIMEDocx,
	"application/octet-stream": MIMEOctetStream,
	"video/mp4":                MIMEMp4,
	"video/webm":               MIMEWebm,
	"video/mov":                MIMEMov,
	"audio/mpeg":               MIMEMp3,
	"audio/wav":                MIMEWav,
	"audio/flac":               MIMEFlac,
	"font/woff2":               MIMEWoff2,
	"font/ttf":                 MIMETtf,
}

func ParseMIME(s string) (MIME, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if m, ok := strToMIME[s]; ok {
		return m, nil
	}
	return MIMEUnknown, fmt.Errorf("Invalid MIME string: %s", s)
}

// IsImage
func (m MIME) IsImage() bool {
	switch m {
	case MIMEJpeg, MIMEPng, MIMEGif, MIMEWebp, MIMESvg, MIMEAvif, MIMEBmp, MIMEIco, MIMETiff, MIMEHeic, MIMEHeif:
		return true
	default:
		return false
	}
}

// Ext
func (m MIME) Ext() string {
	switch m {
	case MIMEJpeg:
		return "jpeg"
	case MIMEPng:
		return "png"
	case MIMEGif:
		return "gif"
	case MIMEWebp:
		return "webp"
	case MIMESvg:
		return "svg"
	case MIMEAvif:
		return "avif"
	case MIMEBmp:
		return "bmp"
	case MIMEIco:
		return "ico"
	case MIMETiff:
		return "tiff"
	case MIMEHeic:
		return "heic"
	case MIMEHeif:
		return "heif"
	case MIMETextPlain:
		return "txt"
	case MIMEHtml:
		return "html"
	case MIMECss:
		return "css"
	case MIMEJavaScript:
		return "js"
	case MIMEJson:
		return "json"
	case MIMECsv:
		return "csv"
	case MIMEMarkdown:
		return "md"
	case MIMEPdf:
		return "pdf"
	case MIMEZip:
		return "zip"
	case MIMEGzip:
		return "gz"
	case MIMETar:
		return "tar"
	case MIMEXlsx:
		return "xlsx"
	case MIMEDocx:
		return "docx"
	case MIMEMp4:
		return "mp4"
	case MIMEWebm:
		return "webm"
	case MIMEMov:
		return "mov"
	case MIMEMp3:
		return "mp3"
	case MIMEWav:
		return "wav"
	case MIMEFlac:
		return "flac"
	case MIMEWoff2:
		return "woff2"
	case MIMETtf:
		return "ttf"
	default:
		return ""
	}
}
