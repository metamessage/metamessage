package ir

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

const (
	mimeStringUnknown = "unknown"

	// Image
	mimeStringJpeg = "image/jpeg"
	mimeStringPng  = "image/png"
	mimeStringGif  = "image/gif"
	mimeStringWebp = "image/webp"
	mimeStringSvg  = "image/svg+xml"
	mimeStringAvif = "image/avif"
	mimeStringBmp  = "image/bmp"
	mimeStringIco  = "image/x-icon"
	mimeStringTiff = "image/tiff"
	mimeStringHeic = "image/heic"
	mimeStringHeif = "image/heif"

	// Text
	mimeStringTextPlain  = "text/plain"
	mimeStringHtml       = "text/html"
	mimeStringCss        = "text/css"
	mimeStringJavaScript = "text/javascript"
	mimeStringJson       = "application/json"
	mimeStringCsv        = "text/csv"
	mimeStringMarkdown   = "text/markdown"

	// Application
	mimeStringPdf         = "application/pdf"
	mimeStringZip         = "application/zip"
	mimeStringGzip        = "application/gzip"
	mimeStringTar         = "application/x-tar"
	mimeStringXlsx        = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	mimeStringDocx        = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	mimeStringOctetStream = "application/octet-stream"

	// Video
	mimeStringMp4  = "video/mp4"
	mimeStringWebm = "video/webm"
	mimeStringMov  = "video/quicktime"

	// Audio
	mimeStringMp3  = "audio/mpeg"
	mimeStringWav  = "audio/wav"
	mimeStringFlac = "audio/flac"

	// Font
	mimeStringWoff2 = "font/woff2"
	mimeStringTtf   = "font/ttf"
)

func (m MIME) String() string {
	switch m {
	// Image
	case MIMEJpeg:
		return mimeStringJpeg
	case MIMEPng:
		return mimeStringPng
	case MIMEGif:
		return mimeStringGif
	case MIMEWebp:
		return mimeStringWebp
	case MIMESvg:
		return mimeStringSvg
	case MIMEAvif:
		return mimeStringAvif
	case MIMEBmp:
		return mimeStringBmp
	case MIMEIco:
		return mimeStringIco
	case MIMETiff:
		return mimeStringTiff
	case MIMEHeic:
		return mimeStringHeic
	case MIMEHeif:
		return mimeStringHeif

	// Text
	case MIMETextPlain:
		return mimeStringTextPlain
	case MIMEHtml:
		return mimeStringHtml
	case MIMECss:
		return mimeStringCss
	case MIMEJavaScript:
		return mimeStringJavaScript
	case MIMEJson:
		return mimeStringJson
	case MIMECsv:
		return mimeStringCsv
	case MIMEMarkdown:
		return mimeStringMarkdown

	// Application
	case MIMEPdf:
		return mimeStringPdf
	case MIMEZip:
		return mimeStringZip
	case MIMEGzip:
		return mimeStringGzip
	case MIMETar:
		return mimeStringTar
	case MIMEXlsx:
		return mimeStringXlsx
	case MIMEDocx:
		return mimeStringDocx
	case MIMEOctetStream:
		return mimeStringOctetStream

	// Video
	case MIMEMp4:
		return mimeStringMp4
	case MIMEWebm:
		return mimeStringWebm
	case MIMEMov:
		return mimeStringMov

	// Audio
	case MIMEMp3:
		return mimeStringMp3
	case MIMEWav:
		return mimeStringWav
	case MIMEFlac:
		return mimeStringFlac

	// Font
	case MIMEWoff2:
		return mimeStringWoff2
	case MIMETtf:
		return mimeStringTtf

	default:
		return mimeStringUnknown
	}
}

var strToMIME = map[string]MIME{
	mimeStringJpeg:        MIMEJpeg,
	mimeStringPng:         MIMEPng,
	mimeStringGif:         MIMEGif,
	mimeStringWebp:        MIMEWebp,
	mimeStringSvg:         MIMESvg,
	mimeStringAvif:        MIMEAvif,
	mimeStringBmp:         MIMEBmp,
	mimeStringIco:         MIMEIco,
	mimeStringTiff:        MIMETiff,
	mimeStringHeic:        MIMEHeic,
	mimeStringHeif:        MIMEHeif,
	mimeStringTextPlain:   MIMETextPlain,
	mimeStringHtml:        MIMEHtml,
	mimeStringCss:         MIMECss,
	mimeStringJavaScript:  MIMEJavaScript,
	mimeStringJson:        MIMEJson,
	mimeStringCsv:         MIMECsv,
	mimeStringMarkdown:    MIMEMarkdown,
	mimeStringPdf:         MIMEPdf,
	mimeStringZip:         MIMEZip,
	mimeStringGzip:        MIMEGzip,
	mimeStringTar:         MIMETar,
	mimeStringXlsx:        MIMEXlsx,
	mimeStringDocx:        MIMEDocx,
	mimeStringOctetStream: MIMEOctetStream,
	mimeStringMp4:         MIMEMp4,
	mimeStringWebm:        MIMEWebm,
	mimeStringMov:         MIMEMov,
	mimeStringMp3:         MIMEMp3,
	mimeStringWav:         MIMEWav,
	mimeStringFlac:        MIMEFlac,
	mimeStringWoff2:       MIMEWoff2,
	mimeStringTtf:         MIMETtf,
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
