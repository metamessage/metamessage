<?php

namespace io\metamessage\core;

use io\metamessage\ir\ValueType;
use PHPUnit\Framework\TestCase;

/**
 * v0.5.0 尾定界模式专项测试。
 * golden hex / 尾部槽档位 / 篡改防御，期望字节与 Go internal/core/footer_test.go 完全一致。
 */
class FooterTest extends TestCase
{
    private static function hex(array $bs): string
    {
        return implode(' ', array_map(fn($b) => sprintf('%02X', $b & 0xFF), $bs));
    }

    public function testSlotTiersCanonicalWidths(): void
    {
        $cases = [[1, 1], [253, 1], [254, 2], [255, 2], [256, 3], [65535, 3]];
        foreach ($cases as [$metaLen, $width]) {
            $slot = Footer::footerSlotBytes($metaLen);
            $this->assertCount($width, $slot, "metaLen=$metaLen slot width");
            $data = array_merge([0x00, 0x11, 0x22, 0x33], $slot);
            [$gotLen, $gotW] = Footer::parseFooterSlot($data, count($data));
            $this->assertSame($metaLen, $gotLen);
            $this->assertSame($width, $gotW);
        }
    }

    public function testSlotTooLongRejected(): void
    {
        $this->expectException(\Exception::class);
        Footer::footerSlotBytes(65536);
    }

    public function testGoldenSingleFieldDesc(): void
    {
        // v0.5.0 尾序：inner 节点 [68 69 82]（body 'hi' + sign 0x80|2 内联长度）；
        // 包装段 [reverse(meta) 62 61 1A][meta_len 槽 03][tag sign E7=0xE0|total7]。
        $this->assertSame('68 69 82 62 61 1A 03 E7', self::hex(MetaMessage::FromValue('hi', 'type=str;desc=ab')));
    }

    public function testGoldenMultiFieldOrder(): void
    {
        $wire = MetaMessage::FromValue('x', 'type=str;example;nullable;desc=ab');
        $this->assertSame('78 81 62 61 1A 29 09 05 E8', self::hex($wire));
        $tag = MetaMessage::DecodeWire($wire)->Tag;
        $this->assertTrue($tag->example);
        $this->assertTrue($tag->nullable);
        $this->assertSame('ab', $tag->desc);
    }

    public function testGoldenRefMetaAtTail(): void
    {
        $wire = MetaMessage::FromValue('r', 'type=str;ref_to=2');
        // meta 正序: [KMore|1 = 0xF9][bitmap 0x02][ref_to 0x02] → 逆序入尾。
        $this->assertSame('72 81 02 02 F9 03 E6', self::hex($wire));
        $this->assertSame(2, MetaMessage::DecodeWire($wire)->Tag->refTo);
    }

    public function testReverseBytesRoundtrip(): void
    {
        $bs = [1, 2, 3, 4, 5];
        $this->assertSame($bs, Footer::reverseBytes(Footer::reverseBytes($bs)));
        $this->assertSame([5, 4, 3, 2, 1], Footer::reverseBytes($bs));
    }

    public function testUnicodeInheritEnumsRoundtripStable(): void
    {
        // footer 编解码往返稳定：encode→decode→re-encode 逐字节还原。
        $srcs = [
            "{\n  // mm: desc=数据; pattern=^\\d$\n  \"v\": \"1\"\n}",
            "{\n  // mm: child_type=i8\n  \"a\": [1, 2]\n}",
            "{\n  // mm: enums=a|b|c\n  \"v\": \"b\"\n}",
        ];
        foreach ($srcs as $src) {
            $wire = MetaMessage::FromJSONC($src);
            $jsonc = MetaMessage::DecodeToJsonc($wire);
            $wire2 = MetaMessage::FromJSONC($jsonc);
            $this->assertSame(self::hex($wire), self::hex($wire2), "roundtrip drift for $src");
        }
    }

    public function testDecodeDefenseTamperedFooter(): void
    {
        $golden = MetaMessage::FromValue('hi', 'type=str;desc=ab'); // 68 69 82 62 61 1A 03 E7

        // 截断尾部（去掉 tag sign 0xE7）→ 定界不再对齐
        $this->expectException(MmDecodeException::class);
        MetaMessage::DecodeWire(array_slice($golden, 0, count($golden) - 1));
    }

    public function testDecodeDefenseEmptyMeta(): void
    {
        $golden = MetaMessage::FromValue('hi', 'type=str;desc=ab');
        $golden[count($golden) - 2] = 0x7F; // meta_len 尾槽声称巨大 meta → inner 下溢
        $this->expectException(MmDecodeException::class);
        MetaMessage::DecodeWire($golden);
    }

    public function testDecodeDefenseOverstatedSlot(): void
    {
        $golden = MetaMessage::FromValue('hi', 'type=str;desc=ab');
        $golden[count($golden) - 2] = 0x20; // meta_len 32 > 封装内可用空间 → tag total underflow
        $this->expectException(MmDecodeException::class);
        MetaMessage::DecodeWire($golden);
    }

    public function testDecodeDefenseInflatedTotal(): void
    {
        $golden = MetaMessage::FromValue('hi', 'type=str;desc=ab');
        $golden[count($golden) - 1] = 0xEB; // tag sign 篡改为 bogus 档位 → 定界越界
        $this->expectException(MmDecodeException::class);
        MetaMessage::DecodeWire($golden);
    }
}
