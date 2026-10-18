// v0.5.0 尾定界解码剥帧（与 Go internal/core/footer.go + tailframe.go 逐字节对齐）：
//
//   1. meta_len 尾部槽编码（footer_slot_bytes）+ 反向解析（parse_footer_slot）+ 字节逆序。
//   2. 解码入口 untailify（node_tail_start / untailify_node / untailify_children）：
//      从节点尾 sign 反向定界，把尾定界 wire 还原成「sign 在前」的头部形态，再交给原有正向解码器。
//
// 编码器各原语已直接按尾定界顺序写出（body→lenBE→sign），本文件不含编码侧换位逻辑，
// 也不设独立 payload_len 槽：各类型 sign 低 5 位的档位/内联语义（v0.3.0 原样）同时承担定界信息。
use crate::core::constants::*;
use crate::core::prefix::*;

const SLOT_DIRECT_MAX: usize = 253; // 一档可直写的 meta 长度上限（尾字节 < 254）
const SLOT_MARKER1: u8 = 254; // 二档标记：倒数第 2 字节为 u8 长度
const SLOT_MARKER2: u8 = 255; // 三档标记：倒数第 3..2 字节为 u16 大端长度

fn err<T>(msg: &str) -> Result<T, std::io::Error> {
    Err(std::io::Error::new(std::io::ErrorKind::InvalidData, msg))
}

/// footer_slot_bytes 编码尾部长度槽，返回按正序写入节点尾部的槽字节。
pub fn footer_slot_bytes(meta_len: usize) -> Vec<u8> {
    if meta_len <= SLOT_DIRECT_MAX {
        vec![meta_len as u8]
    } else if meta_len <= 0xFF {
        vec![(meta_len & 0xFF) as u8, SLOT_MARKER1]
    } else if meta_len <= 0xFFFF {
        vec![((meta_len >> 8) & 0xFF) as u8, (meta_len & 0xFF) as u8, SLOT_MARKER2]
    } else {
        vec![]
    }
}

/// parse_footer_slot 从节点尾端反向解析长度槽，返回 (meta_len, slot_width)。
fn parse_footer_slot(data: &[u8], node_end: usize) -> Result<(usize, usize), std::io::Error> {
    if node_end < 1 {
        return err("footer slot out of range");
    }
    let marker = data[node_end - 1];
    match marker {
        SLOT_MARKER1 => {
            if node_end < 2 {
                return err("footer slot truncated");
            }
            Ok((data[node_end - 2] as usize, 2))
        }
        SLOT_MARKER2 => {
            if node_end < 3 {
                return err("footer slot truncated");
            }
            Ok(((data[node_end - 3] as usize) << 8 | (data[node_end - 2] as usize), 3))
        }
        _ => Ok((marker as usize, 1)),
    }
}

/// reverse_bytes 返回 bs 的字节逆序拷贝。
fn reverse_bytes(bs: &[u8]) -> Vec<u8> {
    bs.iter().rev().cloned().collect()
}

/// meta_prefix_bytes 还原正向 meta 前缀槽：与尾部槽同一套 canonical 三档，
/// 宽度恒等于 footer_slot_bytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
fn meta_prefix_bytes(ml: usize) -> Vec<u8> {
    if ml < 254 {
        vec![ml as u8]
    } else if ml <= 0xFF {
        vec![254, (ml & 0xFF) as u8]
    } else {
        vec![255, ((ml >> 8) & 0xFF) as u8, (ml & 0xFF) as u8]
    }
}

// ----- 各类型档位：返回 (w, inline)，w==0 表示内联档 -----

fn int_len(sign: u8) -> (usize, usize) {
    let l = (sign & INT_LEN_MASK) as usize;
    if l < INT_LEN_1 as usize {
        (0, l)
    } else {
        (l - INT_LEN_1 as usize + 1, 0)
    }
}

fn float_len(sign: u8) -> (usize, usize) {
    let l = (sign & FLOAT_LEN_MASK) as usize;
    if l < FLOAT_LEN_1 as usize {
        (0, l)
    } else {
        (l - FLOAT_LEN_1 as usize + 1, 0)
    }
}

fn string_len(sign: u8) -> (usize, usize) {
    let l = (sign & SUFFIX_MASK) as usize;
    if l < STRING_LEN_1 as usize {
        (0, l)
    } else if l == STRING_LEN_1 as usize {
        (1, l)
    } else {
        (2, l)
    }
}

fn bytes_len(sign: u8) -> (usize, usize) {
    let l = (sign & SUFFIX_MASK) as usize;
    match l {
        _ if l == BYTES_LEN_1 as usize => (1, 0),
        _ if l == BYTES_LEN_2 as usize => (2, 0),
        _ if l == BYTES_LEN_3 as usize => (3, 0),
        _ if l == BYTES_LEN_4 as usize => (4, 0),
        _ => (0, l),
    }
}

fn container_len(sign: u8) -> (usize, usize) {
    let l = (sign & CONTAINER_LEN_MASK) as usize;
    match l {
        _ if l == CONTAINER_LEN_1 as usize => (1, 0),
        _ if l == CONTAINER_LEN_2 as usize => (2, 0),
        _ if l == CONTAINER_LEN_3 as usize => (3, 0),
        _ if l == CONTAINER_LEN_4 as usize => (4, 0),
        _ => (0, l),
    }
}

fn tag_len(sign: u8) -> (usize, usize) {
    let l = (sign & SUFFIX_MASK) as usize;
    match l {
        _ if l == TAG_LEN_1 as usize => (1, 0),
        _ if l == TAG_LEN_2 as usize => (2, 0),
        _ if l == TAG_LEN_3 as usize => (3, 0),
        _ if l == TAG_LEN_4 as usize => (4, 0),
        _ => (0, l),
    }
}

/// tier_of：依 sign 高 3 位前缀分派，返回 (w, inline)。
fn tier_of(sign: u8) -> (usize, usize) {
    match sign & PREFIX_MASK {
        PREFIX_SIMPLE => (0, 0),
        PREFIX_POSITIVE_INT | PREFIX_NEGATIVE_INT => int_len(sign),
        PREFIX_FLOAT => float_len(sign),
        PREFIX_STRING => string_len(sign),
        PREFIX_BYTES => bytes_len(sign),
        PREFIX_CONTAINER => container_len(sign),
        PREFIX_TAG => tag_len(sign),
        _ => (0, 0),
    }
}

fn tier_bytes_be(length: usize, w: usize) -> Vec<u8> {
    match w {
        0 => vec![],
        1 => vec![(length & 0xFF) as u8],
        2 => vec![(length >> 8) as u8, length as u8],
        3 => vec![(length >> 16) as u8, (length >> 8) as u8, length as u8],
        _ => vec![
            (length >> 24) as u8,
            (length >> 16) as u8,
            (length >> 8) as u8,
            length as u8,
        ],
    }
}

// read_tier_be 从 [end-w, end) 读大端档位长度。
fn read_tier_be(data: &[u8], end: usize, w: usize) -> Result<usize, std::io::Error> {
    if end < w {
        return err("truncated: tier bytes out of range");
    }
    let mut v: usize = 0;
    for i in (end - w)..end {
        v = (v << 8) | data[i] as usize;
    }
    Ok(v)
}

// node_tail_start 依尾 sign 反推 wire 中结束于 end 的单节点起始偏移。
fn node_tail_start(wire: &[u8], end: usize) -> Result<usize, std::io::Error> {
    if end < 1 || end > wire.len() {
        return err("truncated: node end out of range");
    }
    let sign = wire[end - 1];
    let (w, inline) = tier_of(sign);
    match sign & PREFIX_MASK {
        PREFIX_SIMPLE => Ok(end - 1),
        PREFIX_POSITIVE_INT | PREFIX_NEGATIVE_INT => {
            if w == 0 {
                Ok(end - 1)
            } else {
                Ok(end - 1 - w)
            }
        }
        // float head [sign][exp][mant w 字节]：body 实宽 w+1 → wire [exp][mant][sign]
        PREFIX_FLOAT => {
            if w == 0 {
                Ok(end - 1)
            } else {
                Ok(end - 2 - w)
            }
        }
        PREFIX_STRING | PREFIX_BYTES | PREFIX_CONTAINER => {
            if w == 0 {
                Ok(end - 1 - inline)
            } else {
                let l = read_tier_be(wire, end - 1, w)?;
                Ok(end - 1 - w - l)
            }
        }
        PREFIX_TAG => {
            let total = if w > 0 { read_tier_be(wire, end - 1, w)? } else { inline };
            let slot_end = end - 1 - w;
            let (ml, slot_w) = parse_footer_slot(wire, slot_end)?;
            // 反推用实际尾槽宽 slot_w：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
            // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
            let inner = total as isize - ml as isize - slot_w as isize;
            if inner < 0 {
                return err("tag total underflow");
            }
            Ok((slot_end as isize - slot_w as isize - ml as isize - inner) as usize)
        }
        _ => err("invalid prefix"),
    }
}

// untailify_node：wire 形态 [start,end) → head 形态字节。
fn untailify_node(wire: &[u8], start: usize, end: usize) -> Result<Vec<u8>, std::io::Error> {
    if start >= end || end > wire.len() {
        return err("untailify: bad range");
    }
    let sign = wire[end - 1];
    let (w, inline) = tier_of(sign);
    match sign & PREFIX_MASK {
        PREFIX_SIMPLE => Ok(vec![sign]),
        PREFIX_POSITIVE_INT | PREFIX_NEGATIVE_INT | PREFIX_FLOAT => {
            if w == 0 {
                return Ok(vec![sign]);
            }
            let mut out = vec![sign];
            out.extend_from_slice(&wire[start..end - 1]);
            Ok(out)
        }
        PREFIX_STRING | PREFIX_BYTES | PREFIX_CONTAINER => {
            let l = if w > 0 { read_tier_be(wire, end - 1, w)? } else { inline };
            let body = &wire[start..end - 1 - w];
            if body.len() != l {
                return err("untailify: payload length mismatch");
            }
            let mut out = vec![sign];
            out.extend_from_slice(&tier_bytes_be(l, w));
            if sign & PREFIX_MASK == PREFIX_CONTAINER {
                let child = untailify_children(body)?;
                out.extend_from_slice(&child);
                return Ok(out);
            }
            out.extend_from_slice(body);
            Ok(out)
        }
        PREFIX_TAG => {
            let total = if w > 0 { read_tier_be(wire, end - 1, w)? } else { inline };
            let slot_end = end - 1 - w;
            let (ml, slot_w) = parse_footer_slot(wire, slot_end)?;
            let inner_len = total as isize - ml as isize - slot_w as isize;
            if inner_len < 0 {
                return err("untailify tag: total underflow");
            }
            let inner_len = inner_len as usize;
            let meta_end = slot_end - slot_w;
            if meta_end < ml {
                return err("untailify tag: meta out of range");
            }
            let meta = reverse_bytes(&wire[meta_end - ml..meta_end]);
            let inner_wire = &wire[start..meta_end - ml];
            if inner_wire.len() != inner_len {
                return err("untailify tag: inner length mismatch");
            }
            let mut inner: Vec<u8> = Vec::new();
            if inner_len > 0 {
                inner = untailify_node(inner_wire, 0, inner_wire.len())?;
            }
            let mut out = vec![sign];
            out.extend_from_slice(&tier_bytes_be(total, w));
            out.extend_from_slice(&meta_prefix_bytes(ml));
            out.extend_from_slice(&meta);
            out.extend_from_slice(&inner);
            Ok(out)
        }
        _ => err("untailify: unreachable"),
    }
}

// untailify_children：wire 容器 payload 内顺序子节点逐个还原（从末端反向扫描，正序拼接）。
fn untailify_children(payload: &[u8]) -> Result<Vec<u8>, std::io::Error> {
    let mut nodes: Vec<Vec<u8>> = Vec::new();
    let mut end = payload.len();
    while end > 0 {
        let start = node_tail_start(payload, end)?;
        if start >= end {
            return err("untailify children: bad boundary");
        }
        nodes.push(untailify_node(payload, start, end)?);
        end = start;
    }
    let mut out: Vec<u8> = Vec::new();
    for node in nodes.iter().rev() {
        out.extend_from_slice(node);
    }
    Ok(out)
}

/// untailify 解码入口：整帧 wire 形态 → head 形态。
pub fn untailify(wire: &[u8]) -> Result<Vec<u8>, std::io::Error> {
    if wire.is_empty() {
        return err("empty wire");
    }
    let start = node_tail_start(wire, wire.len())?;
    if start != 0 {
        return err("footer boundary mismatch");
    }
    untailify_node(wire, start, wire.len())
}
