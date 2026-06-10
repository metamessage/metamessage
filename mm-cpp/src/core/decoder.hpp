#ifndef MMCPP_CORE_DECODER_HPP
#define MMCPP_CORE_DECODER_HPP

#include "../ir/ast.hpp"
#include "constants.hpp"
#include <cmath>
#include <cstdint>
#include <cstring>
#include <iostream>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

namespace {

std::string base64_encode(const std::vector<uint8_t> &data) {
  static const char *chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
                             "abcdefghijklmnopqrstuvwxyz"
                             "0123456789+/";
  std::string out;
  out.reserve(((data.size() + 2) / 3) * 4);
  for (size_t i = 0; i < data.size(); i += 3) {
    uint32_t val = 0;
    int cnt = 0;
    for (int j = 0; j < 3 && i + j < data.size(); ++j) {
      val = (val << 8) | data[i + j];
      cnt++;
    }
    val <<= (3 - cnt) * 8;
    out += chars[(val >> 18) & 0x3F];
    out += chars[(val >> 12) & 0x3F];
    out += (cnt >= 2) ? chars[(val >> 6) & 0x3F] : '=';
    out += (cnt >= 3) ? chars[val & 0x3F] : '=';
  }
  return out;
}

} // namespace

namespace mmc {
namespace core {

using namespace mmc::ir;

class Decoder {
public:
  Decoder() : data_(nullptr), size_(0), offset_(0) {}

  std::shared_ptr<ir::Node> decode(const std::vector<uint8_t> &encoded) {
    data_ = encoded.data();
    size_ = encoded.size();
    offset_ = 0;
    return decodeNode(nullptr);
  }

private:
  const uint8_t *data_;
  size_t size_;
  size_t offset_;

  uint8_t readByte() {
    if (offset_ >= size_)
      throw std::runtime_error("unexpected EOF");
    return data_[offset_++];
  }

  std::vector<uint8_t> readBytes(size_t n) {
    if (offset_ + n > size_)
      throw std::runtime_error("unexpected EOF");
    std::vector<uint8_t> result(data_ + offset_, data_ + offset_ + n);
    offset_ += n;
    return result;
  }

  int64_t readInt(int byteCount) {
    uint64_t val = 0;
    for (int i = 0; i < byteCount; ++i) {
      val = (val << 8) | readByte();
    }
    return static_cast<int64_t>(val);
  }

  uint64_t readUInt(int byteCount) {
    uint64_t val = 0;
    for (int i = 0; i < byteCount; ++i) {
      val = (val << 8) | readByte();
    }
    return val;
  }

  std::shared_ptr<ir::Node> decodeNode(const ir::Tag *parentTag) {
    uint8_t b = readByte();
    Prefix prefix = getPrefix(b);

    switch (prefix) {
    case Prefix::Tag:
      return decodeTagNode(b);
    case Prefix::Simple:
      return decodeSimple(b, parentTag);
    case Prefix::PositiveInt:
      return decodeInt(b, parentTag, true);
    case Prefix::NegativeInt:
      return decodeInt(b, parentTag, false);
    case Prefix::Float:
      return decodeFloat(b, parentTag);
    case Prefix::String:
      return decodeString(b, parentTag);
    case Prefix::Bytes:
      return decodeBytes(b, parentTag);
    case Prefix::Container:
      return decodeContainer(b, parentTag);
    default:
      throw std::runtime_error("invalid prefix");
    }
  }

  std::pair<ir::Tag, int> decodeTagHeader() {
    ir::Tag tag;
    uint8_t b = readByte();
    int l;

    int extraLen = tagExtraLen(b);
    int inlineLen = tagInlineLen(b);
    if (extraLen == 0) {
      l = inlineLen;
    } else if (extraLen == 1) {
      l = readByte();
    } else {
      auto bytes = readBytes(2);
      l = (bytes[0] << 8) | bytes[1];
    }

    int totalRead = 1 + extraLen;
    int payloadLen = l;

    uint8_t tb = readByte();
    int pl;
    if ((tb & BytesLenMask) < BytesLen1Byte) {
      pl = tb & BytesLenMask;
    } else {
      pl = readByte();
    }

    return {tag, totalRead + 1 + pl};
  }

  std::shared_ptr<ir::Node> decodeTagNode(uint8_t prefix) {
    int extraLen = tagExtraLen(prefix);
    int inlineLen = tagInlineLen(prefix);
    int l;

    if (extraLen == 0) {
      l = inlineLen;
    } else if (extraLen == 1) {
      l = readByte();
    } else {
      auto bytes = readBytes(2);
      l = (bytes[0] << 8) | bytes[1];
    }

    size_t frameEnd = offset_ + l;

    ir::Tag tag;

    uint8_t bodyLenByte = readByte();
    int tagBodyLen;
    if (bodyLenByte < 254) {
      tagBodyLen = bodyLenByte;
    } else if (bodyLenByte == 254) {
      tagBodyLen = readByte();
    } else {
      auto bytes = readBytes(2);
      tagBodyLen = (bytes[0] << 8) | bytes[1];
    }

    size_t tagBodyEnd = offset_ + tagBodyLen;

    while (offset_ < tagBodyEnd) {
      decodeTagBytes(tag);
    }

    auto node = decodeNode(&tag);

    offset_ = frameEnd;

    return node;
  }

  int decodeTagBytes(ir::Tag &tag) {
    uint8_t b = readByte();
    uint8_t key = b & 0xF8;
    int payload = static_cast<int>(b & 0x07);

    switch (static_cast<ir::TagKey>(key)) {
    case ir::KIsNull:
      tag.isNull = (payload & 1) == 1;
      if (tag.isNull) {
        tag.nullable = true;
      }
      return 1;

    case ir::KExample:
      tag.example = (payload & 1) == 1;
      return 1;

    case ir::KDeprecated:
      tag.deprecated = (payload & 1) == 1;
      return 1;

    case ir::KNullable:
      tag.nullable = (payload & 1) == 1;
      return 1;

    case ir::KAllowEmpty:
      tag.allowEmpty = (payload & 1) == 1;
      return 1;

    case ir::KUnique:
      tag.unique = (payload & 1) == 1;
      return 1;

    case ir::KChildNullable:
      tag.childNullable = (payload & 1) == 1;
      return 1;

    case ir::KChildAllowEmpty:
      tag.childAllowEmpty = (payload & 1) == 1;
      return 1;

    case ir::KChildUnique:
      tag.childUnique = (payload & 1) == 1;
      return 1;

    case ir::KType: {
      uint8_t tb = readByte();
      tag.type = static_cast<ir::ValueType>(tb);
      return 2;
    }

    case ir::KChildType: {
      uint8_t tb = readByte();
      tag.childType = static_cast<ir::ValueType>(tb);
      return 2;
    }

    case ir::KSize: {
      uint64_t v = readUInt(payload + 1);
      tag.size = static_cast<int>(v);
      return 2 + payload;
    }

    case ir::KChildSize: {
      uint64_t v = readUInt(payload + 1);
      tag.childSize = static_cast<int>(v);
      return 2 + payload;
    }

    case ir::KVersion: {
      uint64_t v = readUInt(payload + 1);
      tag.version = static_cast<int>(v);
      return 2 + payload;
    }

    case ir::KChildVersion: {
      uint64_t v = readUInt(payload + 1);
      tag.childVersion = static_cast<int>(v);
      return 2 + payload;
    }

    case ir::KMore: {
      uint64_t v = readUInt(payload + 1);
      tag.more = static_cast<uint8_t>(v);
      return 2 + payload;
    }

    case ir::KDesc: {
      tag.desc = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.desc.size());
    }

    case ir::KDefaultVal: {
      tag.default_val = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.default_val.size());
    }

    case ir::KMin: {
      tag.min = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.min.size());
    }

    case ir::KMax: {
      tag.max = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.max.size());
    }

    case ir::KEnums: {
      tag.enums = decodeTagString(payload);
      tag.type = ir::ValueType::Enums;
      return decodeTagStrConsumed(payload, tag.enums.size());
    }

    case ir::KPattern: {
      tag.pattern = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.pattern.size());
    }

    case ir::KLocation: {
      tag.locationOffset = std::stoi(decodeTagString(payload));
      return decodeTagStrConsumed(payload, 0);
    }

    case ir::KMime: {
      uint64_t v = readUInt(payload + 1);
      tag.mime = ir::Tag::mimeToString(static_cast<int>(v));
      tag.type = ir::ValueType::Media;
      return 2 + payload;
    }

    case ir::KChildDesc: {
      tag.childDesc = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.childDesc.size());
    }

    case ir::KChildDefaultVal: {
      tag.child_default_val = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.child_default_val.size());
    }

    case ir::KChildMin: {
      tag.childMin = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.childMin.size());
    }

    case ir::KChildMax: {
      tag.childMax = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.childMax.size());
    }

    case ir::KChildEnums: {
      tag.child_enums = decodeTagString(payload);
      tag.childType = ir::ValueType::Enums;
      return decodeTagStrConsumed(payload, tag.child_enums.size());
    }

    case ir::KChildPattern: {
      tag.childPattern = decodeTagString(payload);
      return decodeTagStrConsumed(payload, tag.childPattern.size());
    }

    case ir::KChildLocation: {
      tag.childLocationOffset = std::stoi(decodeTagString(payload));
      return decodeTagStrConsumed(payload, 0);
    }

    case ir::KChildMime: {
      uint64_t v = readUInt(payload + 1);
      tag.childMime = ir::Tag::mimeToString(static_cast<int>(v));
      tag.childType = ir::ValueType::Media;
      return 2 + payload;
    }

    default:
      return 1;
    }
  }

  std::string decodeTagString(int payload) {
    int len = payload;
    if (payload == 6) {
      len = readByte();
    } else if (payload == 7) {
      auto bytes = readBytes(2);
      len = (bytes[0] << 8) | bytes[1];
    }
    auto bytes = readBytes(static_cast<size_t>(len));
    return std::string(bytes.begin(), bytes.end());
  }

  int decodeTagStrConsumed(int payload, size_t strLen) {
    if (payload == 6)
      return 2 + static_cast<int>(strLen);
    if (payload == 7)
      return 3 + static_cast<int>(strLen);
    return 1 + payload;
  }

  uint64_t decodeTagU64() {
    uint8_t b = readByte();
    int byteLen = intLen(b);
    if (byteLen == 0)
      return 0;
    return readUInt(byteLen);
  }

  std::shared_ptr<ir::Node> decodeSimple(uint8_t b, const ir::Tag *parentTag) {
    SimpleValue sv = static_cast<SimpleValue>(getSuffix(b));

    if (sv == SimpleNull) {
      auto node = ir::makeNodeNull();
      if (parentTag) {
        *node->getTag() = *parentTag;
        node->getTag()->inherit(*parentTag);
      }
      if (node->getTag()->type != ir::ValueType::Unknown) {
        throw std::runtime_error("unsupported value types for SimpleNull");
      }
      return node;
    }

    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();

    if (parentTag) {
      *tag = *parentTag;
      tag->inherit(*parentTag);
    }

    switch (sv) {
    case SimpleTrue:
      tag->type = ir::ValueType::Bool;
      val->text = "true";
      break;
    case SimpleFalse:
      tag->type = ir::ValueType::Bool;
      val->text = "false";
      break;
    case SimpleNullBool:
      tag->type = ir::ValueType::Bool;
      tag->isNull = true;
      val->text = "false";
      break;
    case SimpleNullInt:
      tag->type = ir::ValueType::I;
      tag->isNull = true;
      val->text = "0";
      break;
    case SimpleNullFloat:
      tag->type = ir::ValueType::F64;
      tag->isNull = true;
      val->text = "0.0";
      break;
    case SimpleNullString:
      tag->type = ir::ValueType::Str;
      tag->isNull = true;
      val->text = "";
      break;
    case SimpleNullBytes:
      tag->type = ir::ValueType::Bytes;
      tag->isNull = true;
      val->text = "";
      break;
    default:
      break;
    }
    return val;
  }

  std::shared_ptr<ir::Node> decodeInt(uint8_t b, const ir::Tag *parentTag,
                                      bool positive) {
    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();

    if (parentTag) {
      *tag = *parentTag;
      tag->inherit(*parentTag);
    }

    int byteLen = intLen(b);
    uint64_t uv;
    if (byteLen == 0) {
      uv = getSuffix(b);
    } else {
      uv = readUInt(byteLen);
    }

    if (parentTag) {
      if (parentTag->childType == ir::ValueType::Unknown) {
        tag->type = parentTag->type;
      }
    }
    if (tag->type == ir::ValueType::Unknown)
      tag->type = ir::ValueType::I;

    int64_t intValue =
        positive ? static_cast<int64_t>(uv) : -static_cast<int64_t>(uv);
    val->data = intValue;
    val->text = std::to_string(intValue);

    if (tag->type == ir::ValueType::Datetime) {
      int64_t ts =
          positive ? static_cast<int64_t>(uv) : -static_cast<int64_t>(uv);
      int64_t loc_offset = static_cast<int64_t>(tag->locationOffset);
      int64_t adjusted_ts = ts + loc_offset * 3600;
      time_t t = static_cast<time_t>(adjusted_ts);
      struct tm tm;
      gmtime_r(&t, &tm);
      char buf[64];
      strftime(buf, sizeof(buf), "%Y-%m-%d %H:%M:%S", &tm);
      val->text = buf;
    } else if (tag->type == ir::ValueType::Date) {
      int64_t days =
          positive ? static_cast<int64_t>(uv) : -static_cast<int64_t>(uv);
      time_t t = static_cast<time_t>(days * 86400);
      struct tm tm;
      gmtime_r(&t, &tm);
      char buf[64];
      strftime(buf, sizeof(buf), "%Y-%m-%d", &tm);
      val->text = buf;
    } else if (tag->type == ir::ValueType::Time) {
      int64_t secs =
          positive ? static_cast<int64_t>(uv) : -static_cast<int64_t>(uv);
      int hours = static_cast<int>(secs / 3600);
      int mins = static_cast<int>((secs % 3600) / 60);
      int sec = static_cast<int>(secs % 60);
      char buf[16];
      snprintf(buf, sizeof(buf), "%02d:%02d:%02d", hours, mins, sec);
      val->text = buf;
    } else if (tag->type == ir::ValueType::Enums && !tag->enums.empty()) {
      int64_t idx =
          positive ? static_cast<int64_t>(uv) : -static_cast<int64_t>(uv);
      if (idx >= 0) {
        const std::string &enums = tag->enums;
        size_t pos = 0;
        int current = 0;
        size_t start = 0;
        while (pos <= enums.size()) {
          if (pos == enums.size() || enums[pos] == '|') {
            if (current == idx) {
              val->text = enums.substr(start, pos - start);
              break;
            }
            current++;
            start = pos + 1;
          }
          pos++;
        }
      }
    }

    return val;
  }

  std::shared_ptr<ir::Node> decodeFloat(uint8_t b, const ir::Tag *parentTag) {
    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();

    if (parentTag) {
      *tag = *parentTag;
      tag->inherit(*parentTag);
    }

    int mantissaBytes = floatLen(b);
    bool isNegative = (b & FloatPositiveNegativeMask) != 0;
    int8_t exponent = 0;
    uint64_t mantissa = 0;

    if (mantissaBytes == 0) {
      mantissa = b & FloatLenMask;
      exponent = -1;
    } else {
      exponent = static_cast<int8_t>(readByte());
      mantissa = readUInt(mantissaBytes);
    }

    double result = static_cast<double>(mantissa) * std::pow(10.0, exponent);
    if (isNegative)
      result = -result;

    std::ostringstream oss;
    oss << result;
    val->text = oss.str();
    if (tag->type == ir::ValueType::Unknown)
      tag->type = ir::ValueType::F64;
    return val;
  }

  std::shared_ptr<ir::Node> decodeString(uint8_t b, const ir::Tag *parentTag) {
    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();

    if (parentTag) {
      *tag = *parentTag;
      tag->inherit(*parentTag);
    }

    int extraLen = stringExtraLen(b);
    int inlineLen = stringInlineLen(b);
    size_t len;

    if (extraLen == 0) {
      len = static_cast<size_t>(inlineLen);
    } else if (extraLen == 1) {
      len = readByte();
    } else {
      auto bytes = readBytes(2);
      len = static_cast<size_t>((bytes[0] << 8) | bytes[1]);
    }

    auto bytes = readBytes(len);
    val->text = std::string(bytes.begin(), bytes.end());
    if (tag->type == ir::ValueType::Unknown)
      tag->type = ir::ValueType::Str;
    return val;
  }

  std::string decodeBigInt(const std::vector<uint8_t> &data) {
    int digitLen = data[0];
    if (digitLen <= 0)
      return "0";

    // Read bits from the bit-packed data starting from data[1]
    const uint8_t *bits = data.data() + 1;
    size_t bitsLen = data.size() - 1;
    size_t bitPos = 0;

    auto readBit = [&]() -> int {
      if (bitPos >= bitsLen * 8)
        return 0;
      int byteIdx = static_cast<int>(bitPos / 8);
      int bitInByte = 7 - static_cast<int>(bitPos % 8);
      int b = (bits[byteIdx] >> bitInByte) & 1;
      bitPos++;
      return b;
    };

    auto readBits = [&](int n) -> int {
      int val = 0;
      for (int i = 0; i < n; i++) {
        val = (val << 1) | readBit();
      }
      return val;
    };

    int sign = readBit();
    bool neg = (sign == 1);

    std::string result;
    int remaining = digitLen;
    while (remaining > 0) {
      if (remaining >= 3) {
        int val = readBits(10);
        char buf[4];
        snprintf(buf, sizeof(buf), "%03d", val);
        result += buf;
        remaining -= 3;
      } else if (remaining == 2) {
        int val = readBits(7);
        char buf[3];
        snprintf(buf, sizeof(buf), "%02d", val);
        result += buf;
        remaining -= 2;
      } else {
        int val = readBits(4);
        result += std::to_string(val);
        remaining -= 1;
      }
    }

    // Trim leading zeros
    size_t firstNonZero = result.find_first_not_of('0');
    std::string trimmed;
    if (firstNonZero == std::string::npos) {
      trimmed = "0";
    } else {
      trimmed = result.substr(firstNonZero);
    }

    if (neg && trimmed != "0")
      return "-" + trimmed;
    return trimmed;
  }

  std::shared_ptr<ir::Node> decodeBytes(uint8_t b, const ir::Tag *parentTag) {
    auto val = ir::makeNodeScalar();
    auto *tag = val->getTag();

    if (parentTag) {
      *tag = *parentTag;
      tag->inherit(*parentTag);
    }

    int extraLen = bytesExtraLen(b);
    int inlineLen = bytesInlineLen(b);
    size_t len;

    if (extraLen == 0) {
      len = static_cast<size_t>(inlineLen);
    } else if (extraLen == 1) {
      len = readByte();
    } else {
      auto bytes = readBytes(2);
      len = static_cast<size_t>((bytes[0] << 8) | bytes[1]);
    }

    auto rawBytes = readBytes(len);

    if (tag->type == ir::ValueType::Bigint && rawBytes.size() > 1) {
      val->text = decodeBigInt(rawBytes);
    } else if (tag->type == ir::ValueType::Uuid && rawBytes.size() == 16) {
      char buf[64];
      snprintf(buf, sizeof(buf),
               "%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%"
               "02x%02x",
               rawBytes[0], rawBytes[1], rawBytes[2], rawBytes[3], rawBytes[4],
               rawBytes[5], rawBytes[6], rawBytes[7], rawBytes[8], rawBytes[9],
               rawBytes[10], rawBytes[11], rawBytes[12], rawBytes[13],
               rawBytes[14], rawBytes[15]);
      val->text = buf;
    } else {
      val->text = base64_encode(rawBytes);
    }

    if (tag->type == ir::ValueType::Unknown)
      tag->type = ir::ValueType::Bytes;
    return val;
  }

  std::shared_ptr<ir::Node> decodeContainer(uint8_t b,
                                            const ir::Tag *parentTag) {
    bool isArray = (b & ContainerMask) == ContainerArray;
    int extraLen = containerExtraLen(b);
    int inlineLen = containerInlineLen(b);
    size_t len;

    if (extraLen == 0) {
      len = static_cast<size_t>(inlineLen);
    } else if (extraLen == 1) {
      len = readByte();
    } else {
      auto bytes = readBytes(2);
      len = static_cast<size_t>((bytes[0] << 8) | bytes[1]);
    }

    if (isArray) {
      return decodeContainerArray(len, parentTag);
    } else {
      return decodeContainerObject(len, parentTag);
    }
  }

  std::shared_ptr<ir::Node> decodeContainerArray(size_t len,
                                                 const ir::Tag *parentTag) {
    auto arr = ir::makeNodeArray();
    size_t startOffset = offset_;

    // Create a clean item tag that only inherits child_* properties
    ir::Tag itemTag;
    if (parentTag) {
      itemTag.inherit(*parentTag);
      // If no child_type specified, propagate parent's type to items
      if (parentTag->childType == ir::ValueType::Unknown) {
        if (itemTag.type == ir::ValueType::Unknown) {
          itemTag.type = parentTag->type;
        }
        if (itemTag.enums.empty() && !parentTag->enums.empty()) {
          itemTag.enums = parentTag->enums;
        }
        if (itemTag.mime.empty() && !parentTag->mime.empty()) {
          itemTag.mime = parentTag->mime;
        }
        if (itemTag.version == ir::DefaultVersion &&
            parentTag->version != ir::DefaultVersion) {
          itemTag.version = parentTag->version;
        }
        if (itemTag.locationOffset == ir::DefaultLocationOffset &&
            parentTag->locationOffset != ir::DefaultLocationOffset) {
          itemTag.locationOffset = parentTag->locationOffset;
        }
      }
    }

    while (offset_ - startOffset < len) {
      auto item = decodeNode(parentTag ? &itemTag : nullptr);
      arr->items.push_back(item);
    }

    if (parentTag) {
      *arr->getTag() = *parentTag;
    }
    if (arr->getTag()->type == ir::ValueType::Unknown) {
      if (arr->getTag()->size > 0) {
        arr->getTag()->type = ir::ValueType::Arr;
      } else {
        arr->getTag()->type = ir::ValueType::Vec;
      }
    }

    return arr;
  }

  std::shared_ptr<ir::Node> decodeContainerObject(size_t len,
                                                  const ir::Tag *parentTag) {
    auto obj = ir::makeNodeObject();
    size_t startOffset = offset_;
    size_t endOffset = startOffset + len;

    // Create a clean item tag that only inherits child_* properties
    ir::Tag itemTag;
    if (parentTag) {
      itemTag.inherit(*parentTag);
      if (parentTag->childType == ir::ValueType::Unknown) {
        if (itemTag.type == ir::ValueType::Unknown) {
          itemTag.type = parentTag->type;
        }
        if (itemTag.enums.empty() && !parentTag->enums.empty()) {
          itemTag.enums = parentTag->enums;
        }
        if (itemTag.mime.empty() && !parentTag->mime.empty()) {
          itemTag.mime = parentTag->mime;
        }
        if (itemTag.version == ir::DefaultVersion &&
            parentTag->version != ir::DefaultVersion) {
          itemTag.version = parentTag->version;
        }
        if (itemTag.locationOffset == ir::DefaultLocationOffset &&
            parentTag->locationOffset != ir::DefaultLocationOffset) {
          itemTag.locationOffset = parentTag->locationOffset;
        }
      }
    }
    const ir::Tag *childTag = parentTag ? &itemTag : nullptr;

    // First decode the key array
    auto keyArrayNode = decodeNode(childTag);
    auto keyArray = std::dynamic_pointer_cast<ir::NodeArray>(keyArrayNode);

    // Then decode values and pair them with keys
    if (keyArray) {
      for (auto &keyItem : keyArray->items) {
        if (offset_ >= endOffset)
          break;
        auto keyVal = std::dynamic_pointer_cast<ir::NodeScalar>(keyItem);
        auto valNode = decodeNode(childTag);
        if (keyVal && !keyVal->text.empty()) {
          ir::Field field(keyVal->text, valNode);
          obj->fields.push_back(std::move(field));
        }
      }
    } else {
      // Fallback: decode remaining as key-value pairs
      while (offset_ < endOffset) {
        auto keyNode = decodeNode(childTag);
        if (offset_ >= endOffset)
          break;
        auto valNode = decodeNode(childTag);
        auto key = std::dynamic_pointer_cast<ir::NodeScalar>(keyNode);
        if (key && !key->text.empty()) {
          ir::Field field(key->text, valNode);
          obj->fields.push_back(std::move(field));
        }
      }
    }

    offset_ = endOffset;

    if (parentTag) {
      *obj->getTag() = *parentTag;
    }
    if (obj->getTag()->type == ir::ValueType::Unknown) {
      obj->getTag()->type = ir::ValueType::Obj;
    }

    return obj;
  }
};

} // namespace core
} // namespace mmc

#endif