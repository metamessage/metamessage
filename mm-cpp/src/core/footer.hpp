#ifndef MMCPP_CORE_FOOTER_HPP
#define MMCPP_CORE_FOOTER_HPP

// v0.5.0 全尾部定界解码入口：把编码器一趟正向产出的尾定界 wire（body 先落盘、
// 定界大端长度 + sign 追加在尾部）反向剥帧，还原成「sign 在前」的头部形态，
// 再交给原有正向解码器。档位/内联语义与 v0.3.0 一致，仅物理落盘顺序尾化，
// 节点宽度不变。逐字节对齐 Go internal/core/tailframe.go。
//
// meta_len 尾槽三档（末字节为标记或直接值，反向可读）：
//   - <=253：[len]；<=255：[len u8][254]；<=65535：[hi][lo][255]

#include "constants.hpp"
#include <cstddef>
#include <cstdint>
#include <stdexcept>
#include <vector>

namespace mmc {
namespace core {

// footer_slot_bytes：编码尾部长度槽三档。
inline std::vector<uint8_t> footerSlotBytes(size_t metaLen) {
  std::vector<uint8_t> slot;
  if (metaLen <= 253) {
    slot.push_back(static_cast<uint8_t>(metaLen));
  } else if (metaLen <= 0xFF) {
    slot.push_back(static_cast<uint8_t>(metaLen & 0xFF));
    slot.push_back(254);
  } else if (metaLen <= 0xFFFF) {
    slot.push_back(static_cast<uint8_t>((metaLen >> 8) & 0xFF));
    slot.push_back(static_cast<uint8_t>(metaLen & 0xFF));
    slot.push_back(255);
  }
  return slot;
}

inline std::vector<uint8_t> reverseBytes(const uint8_t *src, size_t len) {
  std::vector<uint8_t> dst(len);
  for (size_t i = 0; i < len; i++)
    dst[i] = src[len - 1 - i];
  return dst;
}

struct Tier {
  int prefix;
  int w;
  int inlineLen;
};

// 各类型 sign 的 (档位字节宽度 w, 内联长度/值 inlineLen)。未知前缀返回 false。
inline bool tierOf(uint8_t sign, Tier &t) {
  t.w = 0;
  t.inlineLen = 0;
  Prefix p = getPrefix(sign);
  t.prefix = static_cast<int>(p);
  switch (p) {
  case Prefix::Simple:
    return true;
  case Prefix::PositiveInt:
  case Prefix::NegativeInt: {
    int l = static_cast<int>(sign & IntLenMask);
    if (l < IntLen1Byte)
      t.inlineLen = l;
    else
      t.w = l - IntLen1Byte + 1;
    return true;
  }
  case Prefix::Float: {
    int l = static_cast<int>(sign & FloatLenMask);
    if (l < FloatLen1Byte)
      t.inlineLen = l;
    else
      t.w = l - FloatLen1Byte + 1;
    return true;
  }
  case Prefix::String: {
    int l = static_cast<int>(sign & StringLenMask);
    t.inlineLen = l;
    if (l < StringLen1Byte)
      t.w = 0;
    else if (l == StringLen1Byte)
      t.w = 1;
    else
      t.w = 2;
    return true;
  }
  case Prefix::Bytes: {
    int l = static_cast<int>(sign & BytesLenMask);
    switch (l) {
    case BytesLen1Byte:
      t.w = 1;
      break;
    case BytesLen2Byte:
      t.w = 2;
      break;
    case BytesLen3Byte:
      t.w = 3;
      break;
    case BytesLen4Byte:
      t.w = 4;
      break;
    default:
      t.w = 0;
      t.inlineLen = l;
    }
    return true;
  }
  case Prefix::Container: {
    int l = static_cast<int>(sign & ContainerLenMask);
    switch (l) {
    case ContainerLen1Byte:
      t.w = 1;
      break;
    case ContainerLen2Byte:
      t.w = 2;
      break;
    case ContainerLen3Byte:
      t.w = 3;
      break;
    case ContainerLen4Byte:
      t.w = 4;
      break;
    default:
      t.w = 0;
      t.inlineLen = l;
    }
    return true;
  }
  case Prefix::Tag: {
    int l = static_cast<int>(sign & TagLenMask);
    switch (l) {
    case TagLen1Byte:
      t.w = 1;
      break;
    case TagLen2Byte:
      t.w = 2;
      break;
    case TagLen3Byte:
      t.w = 3;
      break;
    case TagLen4Byte:
      t.w = 4;
      break;
    default:
      t.w = 0;
      t.inlineLen = l;
    }
    return true;
  }
  }
  return false;
}

inline bool readTierBE(const uint8_t *data, size_t end, int w, size_t &out) {
  if (static_cast<int64_t>(end) - w < 0)
    return false;
  size_t v = 0;
  for (size_t i = end - w; i < end; i++)
    v = (v << 8) | data[i];
  out = v;
  return true;
}

inline void tierBytesBE(std::vector<uint8_t> &b, size_t length, int w) {
  switch (w) {
  case 0:
    break;
  case 1:
    b.push_back(static_cast<uint8_t>(length & 0xFF));
    break;
  case 2:
    b.push_back(static_cast<uint8_t>((length >> 8) & 0xFF));
    b.push_back(static_cast<uint8_t>(length & 0xFF));
    break;
  case 3:
    b.push_back(static_cast<uint8_t>((length >> 16) & 0xFF));
    b.push_back(static_cast<uint8_t>((length >> 8) & 0xFF));
    b.push_back(static_cast<uint8_t>(length & 0xFF));
    break;
  default:
    b.push_back(static_cast<uint8_t>((length >> 24) & 0xFF));
    b.push_back(static_cast<uint8_t>((length >> 16) & 0xFF));
    b.push_back(static_cast<uint8_t>((length >> 8) & 0xFF));
    b.push_back(static_cast<uint8_t>(length & 0xFF));
    break;
  }
}

inline bool parseFooterSlot(const uint8_t *data, size_t nodeEnd, size_t &ml,
                            size_t &slotW) {
  if (nodeEnd < 1)
    return false;
  uint8_t marker = data[nodeEnd - 1];
  if (marker == 254) {
    if (nodeEnd < 2)
      return false;
    ml = data[nodeEnd - 2];
    slotW = 2;
    return true;
  } else if (marker == 255) {
    if (nodeEnd < 3)
      return false;
    ml = (static_cast<size_t>(data[nodeEnd - 3]) << 8) | data[nodeEnd - 2];
    slotW = 3;
    return true;
  }
  ml = marker;
  slotW = 1;
  return true;
}

// metaPrefixBytes：还原正向 meta 前缀槽——与尾部槽同一套 canonical 三档，
// 宽度恒等于 footerSlotBytes 槽宽（ml<254→1，254..255→2，256..65535→3）。
inline void metaPrefixBytes(std::vector<uint8_t> &b, size_t ml) {
  if (ml < 254) {
    b.push_back(static_cast<uint8_t>(ml));
  } else if (ml <= 255) {
    b.push_back(254);
    b.push_back(static_cast<uint8_t>(ml & 0xFF));
  } else {
    b.push_back(255);
    b.push_back(static_cast<uint8_t>((ml >> 8) & 0xFF));
    b.push_back(static_cast<uint8_t>(ml & 0xFF));
  }
}

// nodeTailStart：计算 wire 中 [?,end) 单节点的起始偏移（反向定界）。
inline bool nodeTailStart(const uint8_t *wire, size_t wlen, size_t end,
                          size_t &start) {
  if (end < 1 || end > wlen)
    return false;
  uint8_t sign = wire[end - 1];
  Tier t;
  if (!tierOf(sign, t))
    return false;
  switch (static_cast<Prefix>(t.prefix)) {
  case Prefix::Simple:
    start = end - 1;
    return true;
  case Prefix::PositiveInt:
  case Prefix::NegativeInt:
    start = (t.w == 0) ? end - 1 : end - 1 - static_cast<size_t>(t.w);
    return true;
  case Prefix::Float:
    // float head [sign][exp][mant w 字节]：body 实宽 w+1 → wire [exp][mant][sign]
    start = (t.w == 0) ? end - 1 : end - 2 - static_cast<size_t>(t.w);
    return true;
  case Prefix::String:
  case Prefix::Bytes:
  case Prefix::Container: {
    if (t.w == 0) {
      if (end - 1 < static_cast<size_t>(t.inlineLen))
        return false;
      start = end - 1 - static_cast<size_t>(t.inlineLen);
    } else {
      size_t l;
      if (!readTierBE(wire, end - 1, t.w, l))
        return false;
      if (static_cast<int64_t>(end) - 1 - t.w - static_cast<int64_t>(l) < 0)
        return false;
      start = end - 1 - static_cast<size_t>(t.w) - l;
    }
    return true;
  }
  case Prefix::Tag: {
    size_t total;
    if (t.w > 0) {
      if (!readTierBE(wire, end - 1, t.w, total))
        return false;
    } else {
      total = static_cast<size_t>(t.inlineLen);
    }
    size_t slotEnd = end - 1 - static_cast<size_t>(t.w);
    size_t ml, slotW;
    if (!parseFooterSlot(wire, slotEnd, ml, slotW))
      return false;
    // 反推用实际尾槽宽 slotW：与编码侧 total=inner+ml+len(尾槽) 严格互逆，
    // 且尾槽/前缀槽已统一为同一套三档，ml>255 不再回绕。
    int64_t inner = static_cast<int64_t>(total) - static_cast<int64_t>(ml) -
                    static_cast<int64_t>(slotW);
    if (inner < 0)
      return false;
    int64_t s = static_cast<int64_t>(slotEnd) - static_cast<int64_t>(slotW) -
                static_cast<int64_t>(ml) - inner;
    if (s < 0)
      return false;
    start = static_cast<size_t>(s);
    return true;
  }
  default:
    return false;
  }
}

inline bool untailifyChildren(const uint8_t *payload, size_t plen,
                              std::vector<uint8_t> &out);

// untailifyNode：wire 形态 [start,end) → head 形态，追加到 out。
inline bool untailifyNode(const uint8_t *wire, size_t start, size_t end,
                          std::vector<uint8_t> &out) {
  if (start >= end)
    return false;
  uint8_t sign = wire[end - 1];
  Tier t;
  if (!tierOf(sign, t))
    return false;
  switch (static_cast<Prefix>(t.prefix)) {
  case Prefix::Simple:
    out.push_back(sign);
    return true;
  case Prefix::PositiveInt:
  case Prefix::NegativeInt:
  case Prefix::Float: {
    if (t.w == 0) {
      out.push_back(sign);
      return true;
    }
    out.push_back(sign);
    out.insert(out.end(), wire + start, wire + (end - 1));
    return true;
  }
  case Prefix::String:
  case Prefix::Bytes:
  case Prefix::Container: {
    size_t l;
    if (t.w > 0) {
      if (!readTierBE(wire, end - 1, t.w, l))
        return false;
    } else {
      l = static_cast<size_t>(t.inlineLen);
    }
    if (static_cast<int64_t>(end) - 1 - t.w - static_cast<int64_t>(start) < 0)
      return false;
    size_t bodyLen = end - 1 - static_cast<size_t>(t.w) - start;
    if (bodyLen != l)
      return false;
    out.push_back(sign);
    tierBytesBE(out, l, t.w);
    if (static_cast<Prefix>(t.prefix) == Prefix::Container)
      return untailifyChildren(wire + start, bodyLen, out);
    out.insert(out.end(), wire + start, wire + start + bodyLen);
    return true;
  }
  case Prefix::Tag: {
    size_t total;
    if (t.w > 0) {
      if (!readTierBE(wire, end - 1, t.w, total))
        return false;
    } else {
      total = static_cast<size_t>(t.inlineLen);
    }
    size_t slotEnd = end - 1 - static_cast<size_t>(t.w);
    size_t ml, slotW;
    if (!parseFooterSlot(wire, slotEnd, ml, slotW))
      return false;
    int64_t innerLenS = static_cast<int64_t>(total) - static_cast<int64_t>(ml) -
                        static_cast<int64_t>(slotW);
    if (innerLenS < 0)
      return false;
    size_t innerLen = static_cast<size_t>(innerLenS);
    size_t metaEnd = slotEnd - slotW;
    if (metaEnd < ml || (metaEnd - ml) < start)
      return false;
    size_t innerWireLen = metaEnd - ml - start;
    if (innerWireLen != innerLen)
      return false;

    out.push_back(sign);
    tierBytesBE(out, total, t.w);
    metaPrefixBytes(out, ml);
    // meta：wire[metaEnd-ml .. metaEnd) 是逆序区，还原成正序。
    for (size_t i = 0; i < ml; i++)
      out.push_back(wire[metaEnd - 1 - i]);
    if (innerLen > 0) {
      if (!untailifyNode(wire, start, start + innerLen, out))
        return false;
    }
    return true;
  }
  default:
    return false;
  }
}

// untailifyChildren：wire 容器 payload 内顺序 wire 子节点逐个还原，正序拼接。
inline bool untailifyChildren(const uint8_t *payload, size_t plen,
                              std::vector<uint8_t> &out) {
  std::vector<std::vector<uint8_t>> nodes;
  size_t end = plen;
  while (end > 0) {
    size_t start;
    if (!nodeTailStart(payload, plen, end, start))
      return false;
    if (start >= end)
      return false;
    std::vector<uint8_t> node;
    if (!untailifyNode(payload, start, end, node))
      return false;
    nodes.push_back(std::move(node));
    end = start;
  }
  for (size_t i = nodes.size(); i-- > 0;) {
    out.insert(out.end(), nodes[i].begin(), nodes[i].end());
  }
  return true;
}

// untailify：解码入口，整帧 wire 形态 → head 形态。
inline std::vector<uint8_t> untailify(const std::vector<uint8_t> &wire) {
  if (wire.empty())
    throw std::runtime_error("empty wire");
  size_t start;
  if (!nodeTailStart(wire.data(), wire.size(), wire.size(), start))
    throw std::runtime_error("footer boundary mismatch");
  if (start != 0)
    throw std::runtime_error("footer boundary mismatch");
  std::vector<uint8_t> out;
  if (!untailifyNode(wire.data(), start, wire.size(), out))
    throw std::runtime_error("untailify failed");
  return out;
}

} // namespace core
} // namespace mmc

#endif
