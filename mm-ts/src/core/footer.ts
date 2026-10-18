// v0.5.0 尾定界辅助：
//   1. meta_len 尾部槽（footerSlotBytes / parseFooterSlot）+ 字节逆序（reverseBytes），
//      与 Go internal/core/footer.go 逐字节对齐。
//   2. 解码入口剥帧 Untailify（nodeTailStart / untailifyNode / untailifyChildren），
//      与 Go internal/core/tailframe.go 逐字节对齐：从节点尾 sign 反向定界，把尾定界
//      wire 还原成「sign 在前」的头部形态，再交给原有正向解码器。
//
// 编码器各原语已改为直接按尾定界顺序写出（body→lenBE→sign），本文件不再包含编码侧
// 换位逻辑，也不设独立 payload_len 槽（各类型 sign 低 5 位档位/内联语义承担定界）。
//
// meta_len 尾槽三档（canonical 最小宽度，末字节为标记或直接值）：
//   - metaLen <= 253：[metaLen]                     槽宽 1
//   - metaLen <= 255：[metaLen][254]                槽宽 2
//   - metaLen <= 65535：[hi][lo][255]（u16 大端）    槽宽 3
import {
  Prefix,
  getPrefix,
  intLen,
  floatLen,
  stringLen,
  bytesLen,
  containerLen,
  tagLen,
  Max1Byte,
  Max2Byte,
} from './constants';

export const FooterSlotDirectMax = 253; // 一档可直写的 meta 长度上限（尾字节 < 254）
export const FooterSlotMarker1 = 254; // 二档标记：倒数第 2 字节为 u8 长度
export const FooterSlotMarker2 = 255; // 三档标记：倒数第 3..2 字节为 u16 大端长度

// footerSlotBytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
export function footerSlotBytes(metaLen: number): Uint8Array {
  if (metaLen <= FooterSlotDirectMax) {
    return new Uint8Array([metaLen]);
  } else if (metaLen <= Max1Byte) {
    return new Uint8Array([metaLen, FooterSlotMarker1]);
  } else if (metaLen <= Max2Byte) {
    return new Uint8Array([
      (metaLen >> 8) & 0xff,
      metaLen & 0xff,
      FooterSlotMarker2,
    ]);
  }
  throw new Error(`tag too long, max length: ${Max2Byte}, actual: ${metaLen}`);
}

// parseFooterSlot 从节点尾端反向解析长度槽。
// data 为完整 wire 缓冲，nodeEnd 为节点（含封装）的结束偏移（不含）。
// 返回 [metaLen, slotWidth]。
export function parseFooterSlot(
  data: Uint8Array,
  nodeEnd: number,
): [number, number] {
  if (nodeEnd < 1) {
    throw new Error('footer slot out of range');
  }
  const marker = data[nodeEnd - 1] ?? 0;
  switch (marker) {
    case FooterSlotMarker1:
      if (nodeEnd < 2) {
        throw new Error('footer slot truncated');
      }
      return [data[nodeEnd - 2] ?? 0, 2];
    case FooterSlotMarker2:
      if (nodeEnd < 3) {
        throw new Error('footer slot truncated');
      }
      return [((data[nodeEnd - 3] ?? 0) << 8) | (data[nodeEnd - 2] ?? 0), 3];
    default:
      return [marker, 1];
  }
}

// reverseBytes 返回 bs 的字节逆序拷贝。
export function reverseBytes(bs: Uint8Array): Uint8Array {
  const out = new Uint8Array(bs.length);
  for (let i = 0; i < bs.length; i++) {
    out[bs.length - 1 - i] = bs[i] as number;
  }
  return out;
}

// ---------------------------------------------------------------------------
// 解码剥帧：wire（尾定界）→ head（sign 在前）
// ---------------------------------------------------------------------------

// metaPrefixWidth 正向 meta 前缀槽宽度：与尾部槽同一套 canonical 三档，
// 宽度恒等于 footerSlotBytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
export function metaPrefixWidth(ml: number): number {
  if (ml < 254) return 1;
  if (ml <= 255) return 2;
  return 3;
}

// metaPrefixBytes 还原正向 meta 前缀槽。
export function metaPrefixBytes(ml: number): Uint8Array {
  if (ml < 254) return new Uint8Array([ml]);
  if (ml <= Max1Byte) return new Uint8Array([254, ml]);
  if (ml <= Max2Byte)
    return new Uint8Array([255, (ml >> 8) & 0xff, ml & 0xff]);
  throw new Error(`meta too long: ${ml}`);
}

// tierOf 依 sign 的 3 位类型前缀分派各档位函数，返回 [prefix, w, inline]。
// w==0 表示内联档（无独立档位长度字节）。
function tierOf(sign: number): [Prefix, number, number] {
  const p = getPrefix(sign);
  switch (p) {
    case Prefix.Simple:
      return [p, 0, 0];
    case Prefix.PositiveInt:
    case Prefix.NegativeInt: {
      const [w, l] = intLen(sign);
      return [p, w, l];
    }
    case Prefix.Float: {
      const [w, l] = floatLen(sign);
      return [p, w, l];
    }
    case Prefix.String: {
      const [w, l] = stringLen(sign);
      return [p, w, l];
    }
    case Prefix.Bytes: {
      const [w, l] = bytesLen(sign);
      return [p, w, l];
    }
    case Prefix.Container: {
      const [w, l] = containerLen(sign);
      return [p, w, l];
    }
    case Prefix.Tag: {
      const [w, l] = tagLen(sign);
      return [p, w, l];
    }
    default:
      throw new Error(`invalid prefix ${p}`);
  }
}

function tierBytesBE(length: number, w: number): number[] {
  switch (w) {
    case 0:
      return [];
    case 1:
      return [length & 0xff];
    case 2:
      return [(length >> 8) & 0xff, length & 0xff];
    case 3:
      return [
        (length >>> 16) & 0xff,
        (length >>> 8) & 0xff,
        length & 0xff,
      ];
    default:
      return [
        Math.floor(length / 0x1000000) & 0xff,
        (length >>> 16) & 0xff,
        (length >>> 8) & 0xff,
        length & 0xff,
      ];
  }
}

// readTierBE 从 [end-w, end) 读大端档位长度。
function readTierBE(data: Uint8Array, end: number, w: number): number {
  if (end - w < 0) {
    throw new Error('truncated: tier bytes out of range');
  }
  let v = 0;
  for (let i = end - w; i < end; i++) {
    v = v * 256 + (data[i] ?? 0);
  }
  return v;
}

// nodeTailStart 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
function nodeTailStart(wire: Uint8Array, end: number): number {
  if (end < 1 || end > wire.length) {
    throw new Error(`truncated: node end ${end} out of range`);
  }
  const sign = wire[end - 1] ?? 0;
  const [p, w, inline] = tierOf(sign);
  switch (p) {
    case Prefix.Simple:
      return end - 1;
    case Prefix.PositiveInt:
    case Prefix.NegativeInt:
      return w === 0 ? end - 1 : end - 1 - w;
    case Prefix.Float:
      // float body 实宽 = w + 1（含 exp）→ wire [exp][mant][sign]
      return w === 0 ? end - 1 : end - 2 - w;
    case Prefix.String:
    case Prefix.Bytes:
    case Prefix.Container:
      if (w === 0) {
        // 内联档：head [sign][body inline] → wire [body inline][sign]
        return end - 1 - inline;
      }
      return end - 1 - w - readTierBE(wire, end - 1, w);
    case Prefix.Tag: {
      const total =
        w === 0 ? inline : readTierBE(wire, end - 1, w);
      const slotEnd = end - 1 - w;
      const [ml, slotW] = parseFooterSlot(wire, slotEnd);
      // 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
      // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
      const inner = total - ml - slotW;
      if (inner < 0) {
        throw new Error('tag total underflow');
      }
      return slotEnd - slotW - ml - inner;
    }
  }
  throw new Error('invalid prefix');
}

// untailifyNode：wire 形态 [start,end) → head 形态字节。
function untailifyNode(
  wire: Uint8Array,
  start: number,
  end: number,
): Uint8Array {
  if (start < 0 || end > wire.length || start >= end) {
    throw new Error(`untailify: bad range [${start},${end})`);
  }
  const sign = wire[end - 1] ?? 0;
  const [p, w, inline] = tierOf(sign);

  switch (p) {
    case Prefix.Simple:
      return new Uint8Array([sign]);
    case Prefix.PositiveInt:
    case Prefix.NegativeInt:
    case Prefix.Float: {
      if (w === 0) {
        return new Uint8Array([sign]);
      }
      const body = wire.slice(start, end - 1);
      return concat(new Uint8Array([sign]), body);
    }
    case Prefix.String:
    case Prefix.Bytes:
    case Prefix.Container: {
      const l = w === 0 ? inline : readTierBE(wire, end - 1, w);
      const body = wire.slice(start, end - 1 - w);
      if (body.length !== l) {
        throw new Error(`untailify: payload ${body.length}, tier says ${l}`);
      }
      const out = concat(
        new Uint8Array([sign]),
        new Uint8Array(tierBytesBE(l, w)),
      );
      if (p === Prefix.Container) {
        return concat(out, untailifyChildren(body));
      }
      return concat(out, body);
    }
    case Prefix.Tag: {
      const total = w === 0 ? inline : readTierBE(wire, end - 1, w);
      const slotEnd = end - 1 - w;
      const [ml, slotW] = parseFooterSlot(wire, slotEnd);
      const innerLen = total - ml - slotW;
      if (innerLen < 0) {
        throw new Error('untailify tag: total underflow');
      }
      const metaEnd = slotEnd - slotW;
      const meta = reverseBytes(wire.slice(metaEnd - ml, metaEnd));
      const innerWire = wire.slice(start, metaEnd - ml);
      if (innerWire.length !== innerLen) {
        throw new Error(
          `untailify tag: inner ${innerWire.length}, want ${innerLen}`,
        );
      }
      let inner: Uint8Array = new Uint8Array(0);
      if (innerLen > 0) {
        inner = untailifyNode(innerWire, 0, innerWire.length);
      }
      const prefix = new Uint8Array(metaPrefixBytes(ml));
      let out = concat(
        new Uint8Array([sign]),
        new Uint8Array(tierBytesBE(total, w)),
      );
      out = concat(out, prefix);
      out = concat(out, meta);
      return concat(out, inner);
    }
  }
  throw new Error('untailify: unreachable');
}

// untailifyChildren：wire 容器 payload 内顺序 wire 子节点逐个还原。
// wire 子节点以尾 sign 结尾，从区域末端反向定界扫描，还原后按正序拼接。
function untailifyChildren(payload: Uint8Array): Uint8Array {
  const nodes: Uint8Array[] = [];
  for (let end = payload.length; end > 0; ) {
    const start = nodeTailStart(payload, end);
    if (start < 0 || start >= end) {
      throw new Error(`untailify children: bad boundary ${start}..${end}`);
    }
    nodes.push(untailifyNode(payload, start, end));
    end = start;
  }
  let out: Uint8Array = new Uint8Array(0);
  for (let i = nodes.length - 1; i >= 0; i--) {
    out = concat(out, nodes[i] as Uint8Array);
  }
  return out;
}

function concat(a: Uint8Array, b: Uint8Array): Uint8Array {
  const out = new Uint8Array(a.length + b.length);
  out.set(a, 0);
  out.set(b, a.length);
  return out;
}

// Untailify 解码入口：整帧 wire 形态 → head 形态。
export function Untailify(wire: Uint8Array): Uint8Array {
  if (wire.length === 0) {
    throw new Error('empty wire');
  }
  const start = nodeTailStart(wire, wire.length);
  if (start !== 0) {
    throw new Error(`footer boundary mismatch: node start ${start}, want 0`);
  }
  return untailifyNode(wire, start, wire.length);
}
