package io.github.metamessage.ir

import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

class MimeTest {

    @Test
    fun parseApplicationJson() {
        assertEquals(16, Mime.parse("application/json"))
    }

    @Test
    fun parseTextHtml() {
        assertEquals(13, Mime.parse("text/html"))
    }

    @Test
    fun parseImagePng() {
        assertEquals(2, Mime.parse("image/png"))
    }

    @Test
    fun parseImageJpeg() {
        assertEquals(1, Mime.parse("image/jpeg"))
    }

    @Test
    fun parseImageJpg() {
        assertEquals(1, Mime.parse("image/jpg"))
    }

    @Test
    fun parseAudioWav() {
        assertEquals(30, Mime.parse("audio/wav"))
    }

    @Test
    fun parseAudioMpeg() {
        assertEquals(29, Mime.parse("audio/mpeg"))
    }

    @Test
    fun parseApplicationJavascript() {
        assertEquals(15, Mime.parse("application/javascript"))
    }

    @Test
    fun parseTextJavascript() {
        assertEquals(15, Mime.parse("text/javascript"))
    }

    @Test
    fun parseVideoMp4() {
        assertEquals(26, Mime.parse("video/mp4"))
    }

    @Test
    fun parseFontWoff2() {
        assertEquals(32, Mime.parse("font/woff2"))
    }

    @Test
    fun parseApplicationOctetStream() {
        assertEquals(25, Mime.parse("application/octet-stream"))
    }

    @Test
    fun parseApplicationZip() {
        assertEquals(20, Mime.parse("application/zip"))
    }

    @Test
    fun parseApplicationPdf() {
        assertEquals(19, Mime.parse("application/pdf"))
    }

    @Test
    fun parseNullReturnsZero() {
        assertEquals(0, Mime.parse(null))
    }

    @Test
    fun parseBlankReturnsZero() {
        assertEquals(0, Mime.parse(""))
        assertEquals(0, Mime.parse("   "))
    }

    @Test
    fun parseInvalidReturnsZero() {
        assertEquals(0, Mime.parse("invalid/type"))
        assertEquals(0, Mime.parse("unknown/mime"))
    }

    @Test
    fun parseCaseInsensitive() {
        assertEquals(16, Mime.parse("Application/JSON"))
        assertEquals(13, Mime.parse("TEXT/HTML"))
        assertEquals(2, Mime.parse("Image/PNG"))
    }

    @Test
    fun parseTrimmed() {
        assertEquals(16, Mime.parse(" application/json "))
        assertEquals(13, Mime.parse("\ttext/html\n"))
    }

    @Test
    fun toStringKnownType() {
        assertEquals("image/jpeg", Mime.toString(1))
        assertEquals("image/png", Mime.toString(2))
        assertEquals("image/gif", Mime.toString(3))
        assertEquals("image/webp", Mime.toString(4))
        assertEquals("image/svg+xml", Mime.toString(5))
        assertEquals("image/avif", Mime.toString(6))
        assertEquals("application/json", Mime.toString(16))
        assertEquals("application/pdf", Mime.toString(19))
        assertEquals("video/mp4", Mime.toString(26))
        assertEquals("audio/wav", Mime.toString(30))
        assertEquals("font/woff2", Mime.toString(32))
        assertEquals("font/ttf", Mime.toString(33))
    }

    @Test
    fun toStringUnknownType() {
        assertEquals("unknown", Mime.toString(0))
        assertEquals("unknown", Mime.toString(999))
    }

    @Test
    fun roundtrip() {
        val mimeTypes = listOf(
            "image/jpeg", "image/png", "image/gif", "image/webp",
            "text/plain", "text/html", "application/json",
            "video/mp4", "audio/mpeg", "font/woff2"
        )
        for (m in mimeTypes) {
            val code = Mime.parse(m)
            assertTrue(code > 0, "MIME '$m' should parse to a positive code")
            assertEquals(m, Mime.toString(code), "MIME '$m' roundtrip failed")
        }
    }
}