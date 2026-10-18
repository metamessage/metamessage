package io.github.metamessage.core

import io.github.metamessage.MetaMessage
import io.github.metamessage.ir.NodeScalar
import io.github.metamessage.ir.Tag
import io.github.metamessage.ir.ValueType
import org.junit.jupiter.api.Assertions.*
import org.junit.jupiter.api.Test

/**
 * v0.5.0 尾定界模式专项测试：尾部槽档位 / 字节逆序 / meta 前缀槽 / 往返稳定 / 篡改防御。
 * 逐字节对齐已由 tests/compare_with_go.sh kt（encode 36/36、decode 35/35）保证，
 * 这里聚焦 footer 侧单元逻辑与防御。
 */
class FooterTest {

    private fun hex(bs: ByteArray): String =
            bs.joinToString(" ") { (it.toInt() and 0xFF).toString(16).padStart(2, '0') }

    @Test
    fun slotTiersCanonicalWidths() {
        val cases = listOf(1 to 1, 253 to 1, 254 to 2, 255 to 2, 256 to 3, 65535 to 3)
        for ((metaLen, width) in cases) {
            val slot = Footer.footerSlotBytes(metaLen)
            assertEquals(width, slot.size, "metaLen=$metaLen slot width")
            val data = byteArrayOf(0x00, 0x11, 0x22, 0x33) + slot
            val parsed = Footer.parseFooterSlot(data, data.size)
            assertEquals(metaLen, parsed[0])
            assertEquals(width, parsed[1])
        }
    }

    @Test
    fun slotTooLongRejected() {
        assertThrows(MmDecodeException::class.java) { Footer.footerSlotBytes(65536) }
    }

    @Test
    fun reverseBytesRoundtrip() {
        val bs = byteArrayOf(1, 2, 3, 4, 5)
        assertArrayEquals(bs, Footer.reverseBytes(Footer.reverseBytes(bs)))
        assertArrayEquals(byteArrayOf(5, 4, 3, 2, 1), Footer.reverseBytes(bs))
    }

    @Test
    fun metaPrefixWidthMatchesSlotWidth() {
        // metaPrefixWidth 与 footerSlotBytes 宽度在常规档位一致（Go 参考行为原样对齐）。
        for (ml in listOf(1, 100, 253, 254, 255)) {
            assertEquals(Footer.footerSlotBytes(ml).size, Footer.metaPrefixWidth(ml), "ml=$ml")
        }
        // metaPrefixBytes 三档还原
        assertArrayEquals(byteArrayOf(10), Footer.metaPrefixBytes(10))
        assertArrayEquals(byteArrayOf(254.toByte(), 255.toByte()), Footer.metaPrefixBytes(255))
    }

    @Test
    fun roundtripStable() {
        val srcs =
                listOf(
                        "{\n  // mm: desc=数据\n  \"v\": \"1\"\n}",
                        "{\n  // mm: child_type=i8\n  \"a\": [1, 2]\n}",
                        "{\n  // mm: enums=a|b|c\n  \"v\": \"b\"\n}",
                        "{\n  \"nested\": {\"x\": [true, false], \"y\": \"str\"}\n}",
                        "{\n  // mm: type=i64\n  \"n\": -1234567890123\n}",
                )
        for (src in srcs) {
            val wire = MetaMessage.encodeFromJsonc(src)
            val jsonc = MetaMessage.decodeToJsonc(wire)
            val wire2 = MetaMessage.encodeFromJsonc(jsonc)
            assertEquals(hex(wire), hex(wire2), "roundtrip drift for $src")
        }
    }

    @Test
    fun decodeDefenseTruncatedTail() {
        val wire = MetaMessage.encodeFromJsonc("{\n  // mm: desc=ab\n  \"v\": \"hi\"\n}")
        // 截断尾字节 → 尾定界不再对齐
        assertThrows(MmDecodeException::class.java) {
            Decoder().decode(wire.copyOfRange(0, wire.size - 1))
        }
    }

    @Test
    fun decodeDefenseInflatedSign() {
        val wire = MetaMessage.encodeFromJsonc("{\n  // mm: desc=ab\n  \"v\": \"hi\"\n}")
        // 篡改对象 sign（末字节）为非法档位 → 定界越界
        val tampered = wire.copyOf()
        tampered[tampered.size - 1] = 0x00
        assertThrows(MmDecodeException::class.java) { Decoder().decode(tampered) }
    }

    // ---- golden 字节布局：裸标量 tag 封装（与 Go / mm-php / mm-py 逐字节一致）----

    private fun strTag() = Tag(type = ValueType.STR)

    @Test
    fun goldenSingleFieldDesc() {
        val tag = strTag()
        tag.desc = "ab"
        // inner [68 69 82]（body 'hi' + 内联长度 sign）；封装 [reverse(meta 62 61 1A)][槽 03][sign E7]。
        assertEquals("68 69 82 62 61 1a 03 e7", hex(Encoder.encodeNode(NodeScalar("hi", "hi", tag))))
    }

    @Test
    fun goldenMultiFieldOrder() {
        val tag = strTag()
        tag.example = true
        tag.nullable = true
        tag.desc = "ab"
        val wire = Encoder.encodeNode(NodeScalar("x", "x", tag))
        assertEquals("78 81 62 61 1a 29 09 05 e8", hex(wire))
        val dt = Decoder().decode(wire) as NodeScalar
        assertTrue(dt.tag!!.example)
        assertTrue(dt.tag!!.nullable)
        assertEquals("ab", dt.tag!!.desc)
    }

    @Test
    fun goldenRefMetaAtTail() {
        val tag = strTag()
        tag.refTo = 2
        val wire = Encoder.encodeNode(NodeScalar("r", "r", tag))
        assertEquals("72 81 02 02 f9 03 e6", hex(wire))
        assertEquals(2, (Decoder().decode(wire) as NodeScalar).tag!!.refTo)
    }
}
