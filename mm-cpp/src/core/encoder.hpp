#ifndef MMCPP_CORE_ENCODER_HPP
#define MMCPP_CORE_ENCODER_HPP

#include "../ir/ast.hpp"
#include "constants.hpp"
#include <algorithm>
#include <cstdint>
#include <ctime>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

namespace {

std::string base64_decode(const std::string &text) {
  if (text.empty())
    return {};

  auto base64_char = [](char c) -> uint8_t {
    if (c >= 'A' && c <= 'Z')
      return (uint8_t)(c - 'A');
    if (c >= 'a' && c <= 'z')
      return (uint8_t)(c - 'a' + 26);
    if (c >= '0' && c <= '9')
      return (uint8_t)(c - '0' + 52);
    if (c == '-' || c == '+')
      return 62;
    if (c == '_' || c == '/')
      return 63;
    return 0xFF;
  };

  size_t padding = 0;
  if (text.size() > 0 && text.back() == '=')
    padding++;
  if (text.size() > 1 && text[text.size() - 2] == '=')
    padding++;

  size_t out_len = text.size() / 4 * 3 - padding;
  std::string out(out_len, '\0');
  size_t o = 0;

  for (size_t i = 0; i < text.size(); i += 4) {
    uint8_t a = base64_char(text[i]);
    uint8_t b = (i + 1 < text.size()) ? base64_char(text[i + 1]) : 0;
    uint8_t c = (i + 2 < text.size()) ? base64_char(text[i + 2]) : 0;
    uint8_t d = (i + 3 < text.size()) ? base64_char(text[i + 3]) : 0;
    if (o < out_len)
      out[o++] = (uint8_t)((a << 2) | (b >> 4));
    if (o < out_len)
      out[o++] = (uint8_t)((b << 4) | (c >> 2));
    if (o < out_len)
      out[o++] = (uint8_t)((c << 6) | d);
  }
  return out;
}

} // anonymous namespace

namespace mmc {
namespace core {

class Encoder {
public:
  Encoder(size_t maxCap = 1024 * 1024 * 1024)
      : buf_(1024), offset_(0), maxCap_(maxCap) {}

  std::vector<uint8_t> encode(std::shared_ptr<ir::Node> node) {
    offset_ = 0;
    ensure(1024);

    uint32_t n = 0;
    switch (node->getType()) {
    case ir::NodeType::NodeObject:
      n = encodeNodeObject(std::static_pointer_cast<ir::NodeObject>(node));
      break;
    case ir::NodeType::NodeArray:
      n = encodeNodeArray(std::static_pointer_cast<ir::NodeArray>(node));
      break;
    case ir::NodeType::Value:
      n = encodeNodeValue(std::static_pointer_cast<ir::NodeScalar>(node));
      break;
    default:
      throw std::runtime_error("unsupported node type");
    }

    std::vector<uint8_t> out(buf_.begin() + offset_ - n,
                             buf_.begin() + offset_);
    offset_ = 0;
    return out;
  }

private:
  std::vector<uint8_t> buf_;
  size_t offset_;
  size_t maxCap_;

  void ensure(size_t needed) {
    if (offset_ + needed > buf_.size()) {
      size_t newSize = buf_.size() * 2;
      if (newSize > maxCap_)
        newSize = maxCap_;
      if (newSize < offset_ + needed)
        newSize = offset_ + needed;
      buf_.resize(newSize);
    }
  }

  uint32_t writeByte(uint8_t b) {
    ensure(1);
    buf_[offset_++] = b;
    return 1;
  }

  uint32_t writeBytes(const std::vector<uint8_t> &data) {
    ensure(data.size());
    std::copy(data.begin(), data.end(), buf_.begin() + offset_);
    offset_ += data.size();
    return static_cast<uint32_t>(data.size());
  }

  uint32_t writeString(const std::string &s) {
    ensure(s.size());
    std::copy(s.begin(), s.end(), buf_.begin() + offset_);
    offset_ += s.size();
    return static_cast<uint32_t>(s.size());
  }

  uint32_t encodeSimple(SimpleValue value) {
    return writeByte(static_cast<uint8_t>(Prefix::Simple) |
                     static_cast<uint8_t>(value));
  }

  uint32_t encodeBool(bool v) {
    return encodeSimple(v ? SimpleTrue : SimpleFalse);
  }

  uint32_t encodeInt(Prefix sign, uint64_t uv) {
    uint8_t prefix;
    int len = 0;

    if (uv <= 23) {
      prefix = static_cast<uint8_t>(sign) | static_cast<uint8_t>(uv);
      len = 0;
    } else if (uv <= 0xFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen1Byte;
      len = 1;
    } else if (uv <= 0xFFFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen2Byte;
      len = 2;
    } else if (uv <= 0xFFFFFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen3Byte;
      len = 3;
    } else if (uv <= 0xFFFFFFFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen4Byte;
      len = 4;
    } else if (uv <= 0xFFFFFFFFFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen5Byte;
      len = 5;
    } else if (uv <= 0xFFFFFFFFFFFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen6Byte;
      len = 6;
    } else if (uv <= 0xFFFFFFFFFFFFFF) {
      prefix = static_cast<uint8_t>(sign) | IntLen7Byte;
      len = 7;
    } else {
      prefix = static_cast<uint8_t>(sign) | IntLen8Byte;
      len = 8;
    }

    writeByte(prefix);
    for (int i = len - 1; i >= 0; --i) {
      writeByte(static_cast<uint8_t>(uv >> (i * 8)));
    }
    return static_cast<uint32_t>(1 + len);
  }

  uint32_t encodeU64(uint64_t uv) { return encodeInt(Prefix::PositiveInt, uv); }

  uint32_t encodeInt64(int64_t v) {
    Prefix sign;
    uint64_t uv;
    if (v >= 0) {
      sign = Prefix::PositiveInt;
      uv = static_cast<uint64_t>(v);
    } else {
      sign = Prefix::NegativeInt;
      uv = static_cast<uint64_t>(-v);
    }
    return encodeInt(sign, uv);
  }

  uint32_t encodeFloat(const std::string &text) {
    bool isNegative = false;
    std::string s = text;
    if (!s.empty() && s[0] == '-') {
      isNegative = true;
      s = s.substr(1);
    }

    int64_t expPart = 0;
    auto ePos = s.find_first_of("eE");
    if (ePos != std::string::npos) {
      expPart = std::stoll(s.substr(ePos + 1));
      s = s.substr(0, ePos);
    }

    auto dotPos = s.find('.');
    std::string intPart, fracPart;
    if (dotPos != std::string::npos) {
      intPart = s.substr(0, dotPos);
      fracPart = s.substr(dotPos + 1);
    } else {
      intPart = s;
    }

    if (intPart.empty())
      intPart = "0";

    int64_t baseExp = -static_cast<int64_t>(fracPart.size()) + expPart;
    if (baseExp < INT8_MIN || baseExp > INT8_MAX) {
      throw std::runtime_error("float exponent out of range: " + text);
    }
    int8_t exponent = static_cast<int8_t>(baseExp);

    std::string mantissaStr = intPart + fracPart;
    size_t firstNonZero = mantissaStr.find_first_not_of('0');
    if (firstNonZero == std::string::npos) {
      mantissaStr = "0";
    } else {
      mantissaStr = mantissaStr.substr(firstNonZero);
    }

    uint64_t mantissa = std::stoull(mantissaStr);

    uint8_t sign = static_cast<uint8_t>(Prefix::Float);
    if (isNegative)
      sign |= FloatPositiveNegativeMask;

    if (exponent == -1 && mantissa <= 7) {
      sign |= static_cast<uint8_t>(mantissa);
      return writeByte(sign);
    }

    uint8_t prefixLen;
    if (mantissa <= 0xFF) {
      prefixLen = FloatLen1Byte;
    } else if (mantissa <= 0xFFFF) {
      prefixLen = FloatLen2Byte;
    } else if (mantissa <= 0xFFFFFF) {
      prefixLen = FloatLen3Byte;
    } else if (mantissa <= 0xFFFFFFFF) {
      prefixLen = FloatLen4Byte;
    } else if (mantissa <= 0xFFFFFFFFFF) {
      prefixLen = FloatLen5Byte;
    } else if (mantissa <= 0xFFFFFFFFFFFF) {
      prefixLen = FloatLen6Byte;
    } else if (mantissa <= 0xFFFFFFFFFFFFFF) {
      prefixLen = FloatLen7Byte;
    } else {
      prefixLen = FloatLen8Byte;
    }

    sign |= prefixLen;
    writeByte(sign);
    writeByte(static_cast<uint8_t>(exponent));

    int byteCount = floatLen(sign);
    for (int i = byteCount - 1; i >= 0; --i) {
      writeByte(static_cast<uint8_t>(mantissa >> (i * 8)));
    }

    return static_cast<uint32_t>(2 + byteCount);
  }

  uint32_t encodeString(const std::string &s) {
    size_t l = s.size();
    if (l < StringLen1Byte) {
      writeByte(static_cast<uint8_t>(Prefix::String) | static_cast<uint8_t>(l));
      writeString(s);
      return static_cast<uint32_t>(1 + l);
    } else if (l <= 255) {
      writeByte(static_cast<uint8_t>(Prefix::String) | StringLen1Byte);
      writeByte(static_cast<uint8_t>(l));
      writeString(s);
      return static_cast<uint32_t>(2 + l);
    } else {
      writeByte(static_cast<uint8_t>(Prefix::String) | StringLen2Byte);
      writeByte(static_cast<uint8_t>(l >> 8));
      writeByte(static_cast<uint8_t>(l));
      writeString(s);
      return static_cast<uint32_t>(3 + l);
    }
  }

  uint32_t encodeBytes(const std::vector<uint8_t> &data) {
    size_t l = data.size();
    if (l < BytesLen1Byte) {
      writeByte(static_cast<uint8_t>(Prefix::Bytes) | static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(1 + l);
    } else if (l <= 255) {
      writeByte(static_cast<uint8_t>(Prefix::Bytes) | BytesLen1Byte);
      writeByte(static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(2 + l);
    } else {
      writeByte(static_cast<uint8_t>(Prefix::Bytes) | BytesLen2Byte);
      writeByte(static_cast<uint8_t>(l >> 8));
      writeByte(static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(3 + l);
    }
  }

  uint32_t encodeArray(const std::vector<uint8_t> &data) {
    size_t l = data.size();
    if (l < ContainerLen1Byte) {
      writeByte(static_cast<uint8_t>(Prefix::Container) | ContainerArray |
                static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(1 + l);
    } else if (l <= 255) {
      writeByte(static_cast<uint8_t>(Prefix::Container) | ContainerArray |
                ContainerLen1Byte);
      writeByte(static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(2 + l);
    } else {
      writeByte(static_cast<uint8_t>(Prefix::Container) | ContainerArray |
                ContainerLen2Byte);
      writeByte(static_cast<uint8_t>(l >> 8));
      writeByte(static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(3 + l);
    }
  }

  uint32_t encodeObject(const std::vector<uint8_t> &data) {
    size_t l = data.size();
    if (l < ContainerLen1Byte) {
      writeByte(static_cast<uint8_t>(Prefix::Container) | ContainerObject |
                static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(1 + l);
    } else if (l <= 255) {
      writeByte(static_cast<uint8_t>(Prefix::Container) | ContainerObject |
                ContainerLen1Byte);
      writeByte(static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(2 + l);
    } else {
      writeByte(static_cast<uint8_t>(Prefix::Container) | ContainerObject |
                ContainerLen2Byte);
      writeByte(static_cast<uint8_t>(l >> 8));
      writeByte(static_cast<uint8_t>(l));
      writeBytes(data);
      return static_cast<uint32_t>(3 + l);
    }
  }

  uint32_t encodeT(const std::vector<uint8_t> &tagBytes) {
    // Writes: [tag_byte_count] [tag_fields]
    if (tagBytes.empty())
      return 0;

    size_t l = tagBytes.size();
    if (l < 254) {
      writeByte(static_cast<uint8_t>(l));
    } else {
      writeByte(254);
      writeByte(static_cast<uint8_t>(l));
    }
    writeBytes(tagBytes);
    return static_cast<uint32_t>((l < 254 ? 1 : 2) + l);
  }

  uint32_t encodeTag(const std::vector<uint8_t> &payload,
                     const std::vector<uint8_t> &tagData) {
    // tagData = [tag_byte_count] [tag_fields]
    // Writes: [PrefixTag | totalLen] [tag_byte_count] [tag_fields] [payload]
    if (tagData.empty())
      return static_cast<uint32_t>(payload.size());

    size_t totalLen = tagData.size() + payload.size();
    uint8_t sign = static_cast<uint8_t>(Prefix::Tag);

    size_t headerLen;
    if (totalLen < TagLen1Byte) {
      sign |= static_cast<uint8_t>(totalLen);
      headerLen = 1;
    } else if (totalLen <= 0xFF) {
      sign |= TagLen1Byte;
      headerLen = 2;
    } else {
      sign |= TagLen2Byte;
      headerLen = 3;
    }

    writeByte(sign);
    if (headerLen == 2) {
      writeByte(static_cast<uint8_t>(totalLen));
    } else if (headerLen == 3) {
      writeByte(static_cast<uint8_t>(totalLen >> 8));
      writeByte(static_cast<uint8_t>(totalLen));
    }

    writeBytes(tagData);
    writeBytes(payload);

    return static_cast<uint32_t>(headerLen + totalLen);
  }

  uint32_t encodeComment(const std::vector<uint8_t> &payload,
                         const ir::Tag *tag) {
    if (tag == nullptr)
      return static_cast<uint32_t>(payload.size());
    auto tagBytes = tag->bytes();
    if (tagBytes.empty())
      return static_cast<uint32_t>(payload.size());

    uint32_t ns = encodeT(tagBytes);
    if (ns == 0)
      return static_cast<uint32_t>(payload.size());

    size_t prevOffset = offset_ - ns;
    std::vector<uint8_t> tagData(buf_.begin() + prevOffset,
                                 buf_.begin() + prevOffset + ns);
    offset_ = prevOffset;

    return encodeTag(payload, tagData);
  }

  uint32_t encodeNodeObject(std::shared_ptr<ir::NodeObject> obj) {
    std::vector<uint8_t> bufKey, buf;

    auto *tag = obj->getTag();

    for (auto &field : obj->fields) {
      uint32_t n = 0;
      switch (field.value->getType()) {
      case ir::NodeType::NodeObject:
        n = encodeNodeObject(std::static_pointer_cast<ir::NodeObject>(field.value));
        break;
      case ir::NodeType::NodeArray:
        n = encodeNodeArray(
            std::static_pointer_cast<ir::NodeArray>(field.value));
        break;
      case ir::NodeType::Value:
        n = encodeNodeValue(
            std::static_pointer_cast<ir::NodeScalar>(field.value));
        break;
      default:
        throw std::runtime_error("unsupported field type");
      }

      auto encodedSub = getEncodedBytes(n);
      buf.insert(buf.end(), encodedSub.begin(), encodedSub.end());

      uint32_t ns = encodeString(field.key);
      auto encodedKey = getEncodedBytes(ns);
      bufKey.insert(bufKey.end(), encodedKey.begin(), encodedKey.end());
    }

    uint32_t nk = encodeArray(bufKey);
    auto encodedKeyArray = getEncodedBytes(nk);
    std::vector<uint8_t> bufAll;
    bufAll.insert(bufAll.end(), encodedKeyArray.begin(), encodedKeyArray.end());
    bufAll.insert(bufAll.end(), buf.begin(), buf.end());

    uint32_t n = encodeObject(bufAll);
    uint32_t n1 = encodeComment(getEncodedBytes(n), tag);
    if (n1 == 0)
      return n;
    return n1;
  }

  uint32_t encodeNodeArray(std::shared_ptr<ir::NodeArray> arr) {
    std::vector<uint8_t> buf;
    auto *tag = arr->getTag();

    for (auto &item : arr->items) {
      uint32_t n = 0;
      switch (item->getType()) {
      case ir::NodeType::NodeObject:
        n = encodeNodeObject(std::static_pointer_cast<ir::NodeObject>(item));
        break;
      case ir::NodeType::NodeArray:
        n = encodeNodeArray(std::static_pointer_cast<ir::NodeArray>(item));
        break;
      case ir::NodeType::Value:
        n = encodeNodeValue(std::static_pointer_cast<ir::NodeScalar>(item));
        break;
      default:
        throw std::runtime_error("unsupported item type");
      }

      auto encodedSub = getEncodedBytes(n);
      buf.insert(buf.end(), encodedSub.begin(), encodedSub.end());
    }

    uint32_t n = encodeArray(buf);
    uint32_t n1 = encodeComment(getEncodedBytes(n), tag);
    if (n1 == 0)
      return n;
    return n1;
  }

  uint32_t encodeBigInt(const std::string &s) {
    if (s.empty())
      return 0;

    bool neg = false;
    std::string body = s;
    if (body[0] == '-') {
      neg = true;
      body = body.substr(1);
    }

    std::vector<int> bits;
    bits.push_back(neg ? 1 : 0);

    size_t i = 0;
    while (i < body.size()) {
      size_t rem = body.size() - i;
      if (rem >= 3) {
        int num = (body[i] - '0') * 100 + (body[i + 1] - '0') * 10 +
                  (body[i + 2] - '0');
        toBitsVec(num, 10, bits);
        i += 3;
      } else if (rem == 2) {
        int num = (body[i] - '0') * 10 + (body[i + 1] - '0');
        toBitsVec(num, 7, bits);
        i += 2;
      } else {
        int num = body[i] - '0';
        toBitsVec(num, 4, bits);
        i += 1;
      }
    }

    std::vector<uint8_t> packed = bitsToBytes(bits);
    std::vector<uint8_t> inner;
    inner.push_back(static_cast<uint8_t>(s.size()));
    inner.insert(inner.end(), packed.begin(), packed.end());

    return encodeBytes(inner);
  }

  static void toBitsVec(int v, int n, std::vector<int> &bits) {
    for (int i = n - 1; i >= 0; --i) {
      bits.push_back((v >> i) & 1);
    }
  }

  static std::vector<uint8_t> bitsToBytes(const std::vector<int> &bits) {
    std::vector<uint8_t> out;
    uint8_t bt = 0;
    int bl = 0;
    for (int b : bits) {
      bt = static_cast<uint8_t>((bt << 1) | b);
      bl++;
      if (bl == 8) {
        out.push_back(bt);
        bt = 0;
        bl = 0;
      }
    }
    if (bl > 0) {
      bt <<= (8 - bl);
      out.push_back(bt);
    }
    return out;
  }

  uint32_t encodeNodeValue(std::shared_ptr<ir::NodeScalar> val) {
    uint32_t n = 0;
    auto *tag = val->getTag();

    switch (tag->type) {
    case ir::ValueType::Datetime: {
      struct tm tm = {};
      if (sscanf(val->text.c_str(), "%d-%d-%d %d:%d:%d", &tm.tm_year,
                 &tm.tm_mon, &tm.tm_mday, &tm.tm_hour, &tm.tm_min,
                 &tm.tm_sec) == 6) {
        tm.tm_year -= 1900;
        tm.tm_mon -= 1;
        tm.tm_isdst = -1;
        int64_t epoch = static_cast<int64_t>(timegm(&tm));
        epoch -= static_cast<int64_t>(tag->locationOffset) * 3600;
        n = encodeInt64(epoch);
      } else {
        n = encodeInt64(0);
      }
      break;
    }

    case ir::ValueType::Date: {
      struct tm tm = {};
      if (sscanf(val->text.c_str(), "%d-%d-%d", &tm.tm_year, &tm.tm_mon,
                 &tm.tm_mday) >= 3) {
        tm.tm_year -= 1900;
        tm.tm_mon -= 1;
        tm.tm_isdst = -1;
        time_t t = timegm(&tm);
        int64_t days = static_cast<int64_t>(t) / 86400;
        n = encodeInt64(days);
      } else {
        n = encodeInt64(0);
      }
      break;
    }

    case ir::ValueType::Time: {
      int hour = 0, minute = 0, sec = 0;
      if (sscanf(val->text.c_str(), "%d:%d:%d", &hour, &minute, &sec) >= 2) {
        int total_secs = hour * 3600 + minute * 60 + sec;
        n = encodeInt64(total_secs);
      } else {
        n = encodeInt64(0);
      }
      break;
    }

    case ir::ValueType::I:
    case ir::ValueType::I8:
    case ir::ValueType::I16:
    case ir::ValueType::I32:
    case ir::ValueType::I64:
      if (tag->isNull)
        n = encodeSimple(SimpleNullInt);
      else
        n = encodeInt64(val->data != 0 || val->text == "0" || val->text.empty()
                            ? val->data
                            : std::stoll(val->text));
      break;

    case ir::ValueType::U:
    case ir::ValueType::U8:
    case ir::ValueType::U16:
    case ir::ValueType::U32:
    case ir::ValueType::U64:
      if (tag->isNull)
        n = 0;
      else
        n = encodeU64(static_cast<uint64_t>(
            val->data != 0 || val->text == "0" || val->text.empty()
                ? val->data
                : static_cast<int64_t>(std::stoull(val->text))));
      break;

    case ir::ValueType::F32:
    case ir::ValueType::F64:
      if (tag->isNull)
        n = encodeSimple(SimpleNullFloat);
      else
        n = encodeFloat(val->text);
      break;

    case ir::ValueType::Decimal:
      if (tag->isNull)
        n = 0;
      else
        n = encodeFloat(val->text);
      break;

    case ir::ValueType::Bigint:
      if (tag->isNull)
        n = 0;
      else
        n = encodeBigInt(val->text);
      break;

    case ir::ValueType::Str:
    case ir::ValueType::Email:
    case ir::ValueType::Url:
      if (tag->isNull)
        n = encodeSimple(SimpleNullString);
      else
        n = encodeString(val->text);
      break;

    case ir::ValueType::Uuid:
      if (tag->isNull)
        n = encodeSimple(SimpleNullBytes);
      else {
        std::string u = val->text;
        u.erase(std::remove(u.begin(), u.end(), '-'), u.end());
        std::vector<uint8_t> bytes;
        for (size_t i = 0; i < u.size(); i += 2) {
          unsigned int byteVal;
          std::string byteStr = u.substr(i, 2);
          sscanf(byteStr.c_str(), "%2x", &byteVal);
          bytes.push_back(static_cast<uint8_t>(byteVal));
        }
        n = encodeBytes(bytes);
      }
      break;

    case ir::ValueType::Bytes:
      if (tag->isNull)
        n = encodeSimple(SimpleNullBytes);
      else {
        std::string decoded = base64_decode(val->text);
        std::vector<uint8_t> bytes(decoded.begin(), decoded.end());
        n = encodeBytes(bytes);
      }
      break;

    case ir::ValueType::Ip:
      if (tag->isNull)
        n = encodeSimple(SimpleNullString);
      else
        n = encodeString(val->text);
      break;

    case ir::ValueType::Image:
    case ir::ValueType::Video:
    case ir::ValueType::Media:
      if (tag->isNull)
        n = encodeSimple(SimpleNullBytes);
      else {
        std::string decoded = base64_decode(val->text);
        std::vector<uint8_t> bytes(decoded.begin(), decoded.end());
        n = encodeBytes(bytes);
      }
      break;

    case ir::ValueType::Bool:
      if (tag->isNull)
        n = encodeSimple(SimpleNullBool);
      else
        n = encodeBool(val->text == "true" || val->text == "1");
      break;

    case ir::ValueType::Enums:
      if (tag->isNull)
        n = 0;
      else {
        int enumIndex = 0;
        if (!tag->enums.empty()) {
          std::istringstream enumStream(tag->enums);
          std::string token;
          int idx = 0;
          while (std::getline(enumStream, token, '|')) {
            token.erase(0, token.find_first_not_of(" "));
            token.erase(token.find_last_not_of(" ") + 1);
            if (token == val->text) {
              enumIndex = idx;
              break;
            }
            idx++;
          }
        }
        n = encodeU64(static_cast<uint64_t>(enumIndex));
      }
      break;

    default:
      throw std::runtime_error("unsupported type: " +
                               ir::valueTypeToString(tag->type));
    }

    uint32_t n1 = encodeComment(getEncodedBytes(n), tag);
    if (n1 == 0)
      return n;
    return n1;
  }

  std::vector<uint8_t> getEncodedBytes(uint32_t written) {
    return std::vector<uint8_t>(buf_.begin() + offset_ - written,
                                buf_.begin() + offset_);
  }
};

} // namespace core
} // namespace mmc

#endif