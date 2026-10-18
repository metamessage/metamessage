import Foundation

/// v0.5.0 尾定界辅助（与 Go internal/core/footer.go + tailframe.go 逐字节对齐）：
///
///   1. meta_len 尾部槽编解码（footerSlotBytes / parseFooterSlot）+ 字节逆序（reverseBytes）。
///   2. 解码入口剥帧 untailify（nodeTailStart / untailifyNode / untailifyChildren）：
///      从节点尾 sign 反向定界，把尾定界 wire 还原成「sign 在前」的头部形态，再交给原有正向解码器。
///
/// 编码器各原语已直接按尾定界顺序写出（body→lenBE→sign），本文件不含编码侧换位逻辑，
/// 也不设独立 payload_len 槽：各类型 sign 低 5 位的档位/内联语义（v0.3.0 原样）同时承担定界信息。
enum Footer {
    static let slotDirectMax = 253 // 一档可直写的 meta 长度上限（尾字节 < 254）
    static let slotMarker1 = 254   // 二档标记：倒数第 2 字节为 u8 长度
    static let slotMarker2 = 255   // 三档标记：倒数第 3..2 字节为 u16 大端长度

    /// footerSlotBytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
    static func footerSlotBytes(_ metaLen: Int) -> [UInt8] {
        if metaLen <= slotDirectMax {
            return [UInt8(metaLen)]
        } else if metaLen <= 0xFF {
            return [UInt8(metaLen), UInt8(slotMarker1)]
        } else if metaLen <= 0xFFFF {
            return [UInt8((metaLen >> 8) & 0xFF), UInt8(metaLen & 0xFF), UInt8(slotMarker2)]
        } else {
            preconditionFailure("tag too long, max length: 65535, actual: \(metaLen)")
        }
    }

    /// parseFooterSlot 从节点尾端反向解析长度槽，返回 (metaLen, slotWidth)。
    static func parseFooterSlot(_ data: [UInt8], _ nodeEnd: Int) throws -> (metaLen: Int, width: Int) {
        guard nodeEnd >= 1 else { throw MMError.invalidData }
        let marker = Int(data[nodeEnd - 1])
        switch marker {
        case slotMarker1:
            guard nodeEnd >= 2 else { throw MMError.invalidData }
            return (Int(data[nodeEnd - 2]), 2)
        case slotMarker2:
            guard nodeEnd >= 3 else { throw MMError.invalidData }
            return ((Int(data[nodeEnd - 3]) << 8) | Int(data[nodeEnd - 2]), 3)
        default:
            return (marker, 1)
        }
    }

    /// reverseBytes 返回 bs 的字节逆序拷贝。
    static func reverseBytes(_ bs: [UInt8]) -> [UInt8] {
        var out = [UInt8](repeating: 0, count: bs.count)
        for i in bs.indices { out[bs.count - 1 - i] = bs[i] }
        return out
    }

    /// metaPrefixWidth 正向 meta 前缀槽宽度：与尾部槽同一套 canonical 三档，宽度恒等于 footerSlotBytes 槽宽。
    static func metaPrefixWidth(_ ml: Int) -> Int {
        if ml < 254 { return 1 } else if ml <= 255 { return 2 } else { return 3 }
    }

    /// metaPrefixBytes 还原正向 meta 前缀槽。
    static func metaPrefixBytes(_ ml: Int) -> [UInt8] {
        if ml < 254 {
            return [UInt8(ml)]
        } else if ml <= 0xFF {
            return [254, UInt8(ml & 0xFF)]
        } else if ml <= 0xFFFF {
            return [255, UInt8((ml >> 8) & 0xFF), UInt8(ml & 0xFF)]
        } else {
            preconditionFailure("meta too long: \(ml)")
        }
    }

    // ----- 各类型档位：返回 (w, inline)，w==0 表示内联档 -----

    private static func tierOf(_ sign: UInt8) throws -> (prefix: UInt8, w: Int, inline: Int) {
        let p = sign & MMConstants.prefixMask
        switch p {
        case MMPrefix.simple.rawValue:
            return (p, 0, 0)
        case MMPrefix.positiveInt.rawValue, MMPrefix.negativeInt.rawValue:
            let (eb, inl) = intLen(sign)
            return (p, eb, inl)
        case MMPrefix.prefixFloat.rawValue:
            let (eb, inl) = floatLen(sign)
            return (p, eb, inl)
        case MMPrefix.prefixString.rawValue:
            let (eb, inl) = stringLen(sign)
            return (p, eb, inl)
        case MMPrefix.prefixBytes.rawValue:
            let (eb, inl) = bytesLen(sign)
            return (p, eb, inl)
        case MMPrefix.container.rawValue:
            let (eb, inl) = containerLen(sign)
            return (p, eb, inl)
        case MMPrefix.prefixTag.rawValue:
            let (eb, inl) = tagLen(sign)
            return (p, eb, inl)
        default:
            throw MMError.invalidPrefix
        }
    }

    private static func tierBytesBE(_ length: Int, _ w: Int) -> [UInt8] {
        if w == 0 { return [] }
        var out = [UInt8](repeating: 0, count: w)
        for i in stride(from: w - 1, through: 0, by: -1) {
            out[w - 1 - i] = UInt8((length >> (8 * i)) & 0xFF)
        }
        return out
    }

    // readTierBE 从 [end-w, end) 读大端档位长度（用 Int 防 0xFFFFFFFF 溢出）。
    private static func readTierBE(_ data: [UInt8], _ end: Int, _ w: Int) throws -> Int {
        guard end - w >= 0 else { throw MMError.invalidData }
        var v = 0
        for i in (end - w)..<end { v = v * 256 + Int(data[i]) }
        return v
    }

    // nodeTailStart 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
    private static func nodeTailStart(_ wire: [UInt8], _ end: Int) throws -> Int {
        guard end >= 1 && end <= wire.count else { throw MMError.invalidData }
        let sign = wire[end - 1]
        let t = try tierOf(sign)
        let p = t.prefix
        let w = t.w
        let inline = t.inline
        switch p {
        case MMPrefix.simple.rawValue:
            return end - 1
        case MMPrefix.positiveInt.rawValue, MMPrefix.negativeInt.rawValue:
            return w == 0 ? end - 1 : end - 1 - w
        // float body 实宽 = w + 1（含 exp）→ wire [exp][mant][sign]
        case MMPrefix.prefixFloat.rawValue:
            return w == 0 ? end - 1 : end - 2 - w
        case MMPrefix.prefixString.rawValue, MMPrefix.prefixBytes.rawValue, MMPrefix.container.rawValue:
            if w == 0 { return end - 1 - inline }
            return end - 1 - w - (try readTierBE(wire, end - 1, w))
        case MMPrefix.prefixTag.rawValue:
            let total = w == 0 ? inline : (try readTierBE(wire, end - 1, w))
            let slotEnd = end - 1 - w
            let slot = try parseFooterSlot(wire, slotEnd)
            let ml = slot.metaLen
            let slotW = slot.width
            // 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
            // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
            let inner = total - ml - slotW
            if inner < 0 { throw MMError.invalidTag }
            return slotEnd - slotW - ml - inner
        default:
            throw MMError.invalidPrefix
        }
    }

    // untailifyNode：wire 形态 [start,end) → head 形态字节。
    private static func untailifyNode(_ wire: [UInt8], _ start: Int, _ end: Int) throws -> [UInt8] {
        guard start >= 0 && end <= wire.count && start < end else { throw MMError.invalidData }
        let sign = wire[end - 1]
        let t = try tierOf(sign)
        let p = t.prefix
        let w = t.w
        let inline = t.inline

        if p == MMPrefix.simple.rawValue { return [sign] }
        if p == MMPrefix.positiveInt.rawValue || p == MMPrefix.negativeInt.rawValue || p == MMPrefix.prefixFloat.rawValue {
            if w == 0 { return [sign] }
            let body = Array(wire[start..<end - 1])
            return [sign] + body
        }
        if p == MMPrefix.prefixString.rawValue || p == MMPrefix.prefixBytes.rawValue || p == MMPrefix.container.rawValue {
            let l = w == 0 ? inline : (try readTierBE(wire, end - 1, w))
            let body = Array(wire[start..<end - 1 - w])
            if body.count != l { throw MMError.invalidData }
            let head = [sign] + tierBytesBE(l, w)
            if p == MMPrefix.container.rawValue {
                let child = try untailifyChildren(body)
                return head + child
            }
            return head + body
        }
        if p == MMPrefix.prefixTag.rawValue {
            let total = w == 0 ? inline : (try readTierBE(wire, end - 1, w))
            let slotEnd = end - 1 - w
            let slot = try parseFooterSlot(wire, slotEnd)
            let ml = slot.metaLen
            let slotW = slot.width
            let innerLen = total - ml - slotW
            if innerLen < 0 { throw MMError.invalidTag }
            let metaEnd = slotEnd - slotW
            let meta = reverseBytes(Array(wire[metaEnd - ml..<metaEnd]))
            let innerWire = Array(wire[start..<metaEnd - ml])
            if innerWire.count != innerLen { throw MMError.invalidTag }
            var inner: [UInt8] = []
            if innerLen > 0 {
                inner = try untailifyNode(innerWire, 0, innerWire.count)
            }
            var out = [sign] + tierBytesBE(total, w)
            out += metaPrefixBytes(ml)
            out += meta
            out += inner
            return out
        }
        throw MMError.invalidData
    }

    // untailifyChildren：wire 容器 payload 内顺序子节点逐个还原（从末端反向扫描，正序拼接）。
    private static func untailifyChildren(_ payload: [UInt8]) throws -> [UInt8] {
        var nodes: [[UInt8]] = []
        var end = payload.count
        while end > 0 {
            let start = try nodeTailStart(payload, end)
            if start < 0 || start >= end { throw MMError.invalidData }
            nodes.append(try untailifyNode(payload, start, end))
            end = start
        }
        var out: [UInt8] = []
        for i in stride(from: nodes.count - 1, through: 0, by: -1) { out += nodes[i] }
        return out
    }

    // untailify 解码入口：整帧 wire 形态 → head 形态。
    static func untailify(_ wire: [UInt8]) throws -> [UInt8] {
        if wire.isEmpty { throw MMError.invalidData }
        let start = try nodeTailStart(wire, wire.count)
        if start != 0 { throw MMError.invalidData }
        return try untailifyNode(wire, start, wire.count)
    }
}
