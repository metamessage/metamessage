<?php

namespace io\metamessage\core;

/**
 * v0.5.0 尾定界辅助（与 Go internal/core/footer.go + tailframe.go 逐字节对齐）：
 *
 *   1. meta_len 尾部槽编解码（footerSlotBytes / parseFooterSlot）+ 字节逆序
 *      （reverseBytes）。
 *   2. 解码入口剥帧 untailify（nodeTailStart / untailifyNode / untailifyChildren）：
 *      从节点尾 sign 反向定界，把尾定界 wire 还原成「sign 在前」的头部形态，
 *      再交给原有正向解码器。
 *
 * 编码器各原语已直接按尾定界顺序写出（body→lenBE→sign），本文件不含编码侧换位
 * 逻辑，也不设独立 payload_len 槽：各类型 sign 低 5 位的档位/内联语义（v0.3.0
 * 原样）同时承担定界信息。
 *
 * meta_len 尾槽三档（canonical 最小宽度，末字节为标记或直接值）：
 *   - metaLen <= 253：[metaLen]                     槽宽 1
 *   - metaLen <= 255：[metaLen][254]                槽宽 2
 *   - metaLen <= 65535：[hi][lo][255]（u16 大端）    槽宽 3
 */
class Footer
{
    public const SLOT_DIRECT_MAX = 253; // 一档可直写的 meta 长度上限（尾字节 < 254）
    public const SLOT_MARKER1 = 254;    // 二档标记：倒数第 2 字节为 u8 长度
    public const SLOT_MARKER2 = 255;    // 三档标记：倒数第 3..2 字节为 u16 大端长度

    // footerSlotBytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
    public static function footerSlotBytes(int $metaLen): array
    {
        if ($metaLen <= self::SLOT_DIRECT_MAX) {
            return [$metaLen];
        } elseif ($metaLen <= WireConstants::MAX_1) {
            return [$metaLen, self::SLOT_MARKER1];
        } elseif ($metaLen <= WireConstants::MAX_2) {
            return [(($metaLen >> 8) & 0xFF), ($metaLen & 0xFF), self::SLOT_MARKER2];
        }
        throw new MmDecodeException('tag too long, max length: ' . WireConstants::MAX_2 . ', actual: ' . $metaLen);
    }

    // parseFooterSlot 从节点尾端反向解析长度槽，返回 [metaLen, slotWidth]。
    public static function parseFooterSlot(array $data, int $nodeEnd): array
    {
        if ($nodeEnd < 1) {
            throw new MmDecodeException('footer slot out of range');
        }
        $marker = $data[$nodeEnd - 1];
        if ($marker === self::SLOT_MARKER1) {
            if ($nodeEnd < 2) {
                throw new MmDecodeException('footer slot truncated');
            }
            return [$data[$nodeEnd - 2], 2];
        } elseif ($marker === self::SLOT_MARKER2) {
            if ($nodeEnd < 3) {
                throw new MmDecodeException('footer slot truncated');
            }
            return [(($data[$nodeEnd - 3] << 8) | $data[$nodeEnd - 2]), 3];
        }
        return [$marker, 1];
    }

    // reverseBytes 返回 bs 的字节逆序拷贝。
    public static function reverseBytes(array $bs): array
    {
        return array_reverse($bs);
    }

    // metaPrefixWidth 正向 meta 前缀槽宽度：与尾部槽同一套 canonical 三档，
    // 宽度恒等于 footerSlotBytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
    public static function metaPrefixWidth(int $ml): int
    {
        if ($ml < 254) {
            return 1;
        }
        if ($ml <= 255) {
            return 2;
        }
        return 3;
    }

    // metaPrefixBytes 还原正向 meta 前缀槽。
    public static function metaPrefixBytes(int $ml): array
    {
        if ($ml < 254) {
            return [$ml];
        }
        if ($ml <= WireConstants::MAX_1) {
            return [254, $ml];
        }
        if ($ml <= WireConstants::MAX_2) {
            return [255, (($ml >> 8) & 0xFF), ($ml & 0xFF)];
        }
        throw new MmDecodeException('meta too long: ' . $ml);
    }

    // ----- 各类型档位：返回 [w, inline]，w==0 表示内联档 -----

    private static function intLen(int $b): array
    {
        $l = $b & WireConstants::INT_LEN_MASK;
        if ($l < WireConstants::INT_LEN_1) {
            return [0, $l];
        }
        return [$l - WireConstants::INT_LEN_1 + 1, 0];
    }

    private static function floatLen(int $b): array
    {
        $l = $b & WireConstants::FLOAT_LEN_MASK;
        if ($l < WireConstants::FLOAT_LEN_1) {
            return [0, $l];
        }
        return [$l - WireConstants::FLOAT_LEN_1 + 1, 0];
    }

    private static function stringLen(int $b): array
    {
        $l = $b & WireConstants::STRING_LEN_MASK;
        if ($l < WireConstants::STRING_LEN_1) {
            return [0, $l];
        } elseif ($l === WireConstants::STRING_LEN_1) {
            return [1, 0];
        }
        return [2, 0];
    }

    private static function bytesLen(int $b): array
    {
        $l = $b & WireConstants::BYTES_LEN_MASK;
        if ($l < WireConstants::BYTES_LEN_1) {
            return [0, $l];
        }
        return [$l - WireConstants::BYTES_LEN_1 + 1, 0];
    }

    private static function containerLen(int $b): array
    {
        $l = $b & WireConstants::CONTAINER_LEN_MASK;
        if ($l < WireConstants::CONTAINER_LEN_1) {
            return [0, $l];
        }
        return [$l - WireConstants::CONTAINER_LEN_1 + 1, 0];
    }

    private static function tagLen(int $b): array
    {
        $l = $b & WireConstants::TAG_LEN_MASK;
        if ($l < WireConstants::TAG_LEN_1) {
            return [0, $l];
        }
        return [$l - WireConstants::TAG_LEN_1 + 1, 0];
    }

    // tierOf 依 sign 的 3 位类型前缀分派，返回 [prefix, w, inline]。
    private static function tierOf(int $sign): array
    {
        $p = $sign & Prefix::PREFIX_MASK;
        switch ($p) {
            case Prefix::SIMPLE:
                return [$p, 0, 0];
            case Prefix::POSITIVE_INT:
            case Prefix::NEGATIVE_INT:
                return array_merge([$p], self::intLen($sign));
            case Prefix::FLOAT:
                return array_merge([$p], self::floatLen($sign));
            case Prefix::STRING:
                return array_merge([$p], self::stringLen($sign));
            case Prefix::BYTES:
                return array_merge([$p], self::bytesLen($sign));
            case Prefix::CONTAINER:
                return array_merge([$p], self::containerLen($sign));
            case Prefix::TAG:
                return array_merge([$p], self::tagLen($sign));
        }
        throw new MmDecodeException('invalid prefix ' . $p);
    }

    private static function tierBytesBE(int $length, int $w): array
    {
        if ($w === 0) {
            return [];
        }
        $out = [];
        for ($i = $w - 1; $i >= 0; $i--) {
            $out[] = (($length >> (8 * $i)) & 0xFF);
        }
        return $out;
    }

    // readTierBE 从 [end-w, end) 读大端档位长度。
    private static function readTierBE(array $data, int $end, int $w): int
    {
        if ($end - $w < 0) {
            throw new MmDecodeException('truncated: tier bytes out of range');
        }
        $v = 0;
        for ($i = $end - $w; $i < $end; $i++) {
            $v = ($v * 256) + $data[$i];
        }
        return $v;
    }

    // nodeTailStart 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
    private static function nodeTailStart(array $wire, int $end): int
    {
        if ($end < 1 || $end > count($wire)) {
            throw new MmDecodeException('truncated: node end ' . $end . ' out of range');
        }
        $sign = $wire[$end - 1];
        [$p, $w, $inline] = self::tierOf($sign);
        switch ($p) {
            case Prefix::SIMPLE:
                return $end - 1;
            case Prefix::POSITIVE_INT:
            case Prefix::NEGATIVE_INT:
                return $w === 0 ? $end - 1 : $end - 1 - $w;
            case Prefix::FLOAT:
                // float body 实宽 = w + 1（含 exp）→ wire [exp][mant][sign]
                return $w === 0 ? $end - 1 : $end - 2 - $w;
            case Prefix::STRING:
            case Prefix::BYTES:
            case Prefix::CONTAINER:
                if ($w === 0) {
                    return $end - 1 - $inline;
                }
                return $end - 1 - $w - self::readTierBE($wire, $end - 1, $w);
            case Prefix::TAG:
                $total = $w === 0 ? $inline : self::readTierBE($wire, $end - 1, $w);
                $slotEnd = $end - 1 - $w;
                [$ml, $slotW] = self::parseFooterSlot($wire, $slotEnd);
                // 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
                // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
                $inner = $total - $ml - $slotW;
                if ($inner < 0) {
                    throw new MmDecodeException('tag total underflow');
                }
                return $slotEnd - $slotW - $ml - $inner;
        }
        throw new MmDecodeException('invalid prefix');
    }

    // untailifyNode：wire 形态 [start,end) → head 形态字节。
    private static function untailifyNode(array $wire, int $start, int $end): array
    {
        if ($start < 0 || $end > count($wire) || $start >= $end) {
            throw new MmDecodeException('untailify: bad range [' . $start . ',' . $end . ')');
        }
        $sign = $wire[$end - 1];
        [$p, $w, $inline] = self::tierOf($sign);

        if ($p === Prefix::SIMPLE) {
            return [$sign];
        }
        if ($p === Prefix::POSITIVE_INT || $p === Prefix::NEGATIVE_INT || $p === Prefix::FLOAT) {
            if ($w === 0) {
                return [$sign];
            }
            $body = array_slice($wire, $start, $end - 1 - $start);
            return array_merge([$sign], $body);
        }
        if ($p === Prefix::STRING || $p === Prefix::BYTES || $p === Prefix::CONTAINER) {
            $l = $w === 0 ? $inline : self::readTierBE($wire, $end - 1, $w);
            $body = array_slice($wire, $start, $end - 1 - $w - $start);
            if (count($body) !== $l) {
                throw new MmDecodeException('untailify: payload ' . count($body) . ', tier says ' . $l);
            }
            $out = array_merge([$sign], self::tierBytesBE($l, $w));
            if ($p === Prefix::CONTAINER) {
                return array_merge($out, self::untailifyChildren($body));
            }
            return array_merge($out, $body);
        }
        if ($p === Prefix::TAG) {
            $total = $w === 0 ? $inline : self::readTierBE($wire, $end - 1, $w);
            $slotEnd = $end - 1 - $w;
            [$ml, $slotW] = self::parseFooterSlot($wire, $slotEnd);
            $innerLen = $total - $ml - $slotW;
            if ($innerLen < 0) {
                throw new MmDecodeException('untailify tag: total underflow');
            }
            $metaEnd = $slotEnd - $slotW;
            $meta = self::reverseBytes(array_slice($wire, $metaEnd - $ml, $ml));
            $innerWire = array_slice($wire, $start, $metaEnd - $ml - $start);
            if (count($innerWire) !== $innerLen) {
                throw new MmDecodeException('untailify tag: inner ' . count($innerWire) . ', want ' . $innerLen);
            }
            $inner = [];
            if ($innerLen > 0) {
                $inner = self::untailifyNode($innerWire, 0, count($innerWire));
            }
            $out = array_merge([$sign], self::tierBytesBE($total, $w));
            $out = array_merge($out, self::metaPrefixBytes($ml));
            $out = array_merge($out, $meta);
            return array_merge($out, $inner);
        }
        throw new MmDecodeException('untailify: unreachable');
    }

    // untailifyChildren：wire 容器 payload 内顺序 wire 子节点逐个还原（从末端反向扫描，正序拼接）。
    private static function untailifyChildren(array $payload): array
    {
        $nodes = [];
        $end = count($payload);
        while ($end > 0) {
            $start = self::nodeTailStart($payload, $end);
            if ($start < 0 || $start >= $end) {
                throw new MmDecodeException('untailify children: bad boundary ' . $start . '..' . $end);
            }
            $nodes[] = self::untailifyNode($payload, $start, $end);
            $end = $start;
        }
        $out = [];
        for ($i = count($nodes) - 1; $i >= 0; $i--) {
            $out = array_merge($out, $nodes[$i]);
        }
        return $out;
    }

    // untailify 解码入口：整帧 wire 形态 → head 形态。
    public static function untailify(array $wire): array
    {
        if (count($wire) === 0) {
            throw new MmDecodeException('empty wire');
        }
        $start = self::nodeTailStart($wire, count($wire));
        if ($start !== 0) {
            throw new MmDecodeException('footer boundary mismatch: node start ' . $start . ', want 0');
        }
        return self::untailifyNode($wire, $start, count($wire));
    }
}
