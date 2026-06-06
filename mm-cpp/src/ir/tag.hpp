#ifndef MMCPP_IR_TAG_HPP
#define MMCPP_IR_TAG_HPP

#include "value_type.hpp"
#include <algorithm>
#include <cstdint>
#include <memory>
#include <sstream>
#include <string>
#include <vector>

namespace mmc {
namespace ir {

enum TagKey : uint8_t {
  KIsNull = 0 << 3,
  KExample = 1 << 3,

  KDeprecated = 2 << 3,
  KDesc = 3 << 3,
  KType = 4 << 3,
  KNullable = 5 << 3,
  KAllowEmpty = 6 << 3,
  KUnique = 7 << 3,
  KDefaultVal = 8 << 3,
  KMin = 9 << 3,
  KMax = 10 << 3,
  KSize = 11 << 3,
  KEnums = 12 << 3,
  KPattern = 13 << 3,
  KLocation = 14 << 3,
  KVersion = 15 << 3,
  KMime = 16 << 3,

  KChildDesc = 17 << 3,
  KChildType = 18 << 3,
  KChildNullable = 19 << 3,
  KChildAllowEmpty = 20 << 3,
  KChildUnique = 21 << 3,
  KChildDefaultVal = 22 << 3,
  KChildMin = 23 << 3,
  KChildMax = 24 << 3,
  KChildSize = 25 << 3,
  KChildEnums = 26 << 3,
  KChildPattern = 27 << 3,
  KChildLocation = 28 << 3,
  KChildVersion = 29 << 3,
  KChildMime = 30 << 3,
  KMore = 31 << 3
};

constexpr int DefaultVersion = 0;
constexpr int DefaultLocationOffset = 0;

struct Tag {
  std::string name;

  bool isNull = false;
  bool example = false;

  std::string desc;
  ValueType type = ValueType::Unknown;
  bool deprecated = false;
  bool nullable = false;
  bool allowEmpty = false;
  bool unique = false;
  std::string default_val;
  std::string min;
  std::string max;
  int size = 0;
  std::string enums;
  std::string pattern;
  int locationOffset = DefaultLocationOffset;
  int version = DefaultVersion;
  std::string mime;
  uint8_t more = 0;

  std::string childDesc;
  ValueType childType = ValueType::Unknown;
  bool childNullable = false;
  bool childAllowEmpty = false;
  bool childUnique = false;
  std::string child_default_val;
  std::string childMin;
  std::string childMax;
  int childSize = 0;
  std::string child_enums;
  std::string childPattern;
  int childLocationOffset = DefaultLocationOffset;
  int childVersion = DefaultVersion;
  std::string childMime;

  bool isInherit = false;

  static Tag create() { return Tag{}; }

  void inherit(const Tag &parent) {
    isInherit = true;
    if (!parent.childDesc.empty())
      desc = parent.childDesc;
    if (parent.childType != ValueType::Unknown)
      type = parent.childType;
    if (parent.childNullable)
      nullable = parent.childNullable;
    if (parent.childAllowEmpty)
      allowEmpty = parent.childAllowEmpty;
    if (parent.childUnique)
      unique = parent.childUnique;
    if (!parent.child_default_val.empty())
      default_val = parent.child_default_val;
    if (!parent.childMin.empty())
      min = parent.childMin;
    if (!parent.childMax.empty())
      max = parent.childMax;
    if (parent.childSize != 0)
      size = parent.childSize;
    if (!parent.child_enums.empty()) {
      enums = parent.child_enums;
      type = ValueType::Enums;
    }
    if (!parent.childPattern.empty())
      pattern = parent.childPattern;
    if (parent.childLocationOffset != DefaultLocationOffset)
      locationOffset = parent.childLocationOffset;
    if (parent.childVersion != DefaultVersion)
      version = parent.childVersion;
    if (!parent.childMime.empty()) {
      mime = parent.childMime;
      type = ValueType::Media;
    }
  }

  std::string toString() const {
    std::ostringstream b;
    bool first = true;
    auto add = [&](const std::string &s) {
      if (!first)
        b << "; ";
      b << s;
      first = false;
    };

    if (type != ValueType::Unknown && !isInherit) {
      bool skip = (type == ValueType::Str || type == ValueType::I ||
                   type == ValueType::F64 || type == ValueType::Bool ||
                   type == ValueType::Obj || type == ValueType::Vec);
      if (!(skip || (type == ValueType::Arr && size > 0) ||
            (type == ValueType::Enums && !enums.empty()) ||
            (type == ValueType::Media && !mime.empty()))) {
        add("type=" + valueTypeToString(type));
      }
    }

    if (example)
      add("example");
    if (isNull)
      add("is_null");
    if (nullable && !isNull && !isInherit)
      add("nullable");
    if (!desc.empty() && !isInherit)
      add("desc=\"" + desc + "\"");
    if (deprecated && !isInherit)
      add("deprecated");
    if (allowEmpty && !isInherit)
      add("allow_empty");
    if (unique && !isInherit)
      add("unique");
    if (!default_val.empty() && !isInherit)
      add("default_val=" + default_val);
    if (!min.empty() && !isInherit)
      add("min=" + min);
    if (!max.empty() && !isInherit)
      add("max=" + max);
    if (size != 0 && !isInherit)
      add("size=" + std::to_string(size));
    if (!enums.empty() && !isInherit)
      add("enums=" + enums);
    if (!pattern.empty() && !isInherit)
      add("pattern=" + pattern);
    if (locationOffset != DefaultLocationOffset && !isInherit)
      add("location=" + std::to_string(locationOffset));
    if (version != DefaultVersion && !isInherit)
      add("version=" + std::to_string(version));
    if (!mime.empty() && !isInherit)
      add("mime=" + mime);
    if (!childDesc.empty())
      add("child_desc=\"" + childDesc + "\"");
    if (childType != ValueType::Unknown) {
      bool childSkip =
          (childType == ValueType::Str || childType == ValueType::I ||
           childType == ValueType::F64 || childType == ValueType::Bool ||
           childType == ValueType::Obj || childType == ValueType::Vec);
      if (!(childSkip || (childType == ValueType::Arr && childSize > 0) ||
            (childType == ValueType::Enums && !child_enums.empty()) ||
            (childType == ValueType::Media && !childMime.empty()))) {
        add("child_type=" + valueTypeToString(childType));
      }
    }
    if (childNullable)
      add("child_nullable");
    if (childAllowEmpty)
      add("child_allow_empty");
    if (childUnique)
      add("child_unique");
    if (!child_default_val.empty())
      add("child_default_val=" + child_default_val);
    if (!childMin.empty())
      add("child_min=" + childMin);
    if (!childMax.empty())
      add("child_max=" + childMax);
    if (childSize != 0)
      add("child_size=" + std::to_string(childSize));
    if (!child_enums.empty())
      add("child_enums=" + child_enums);
    if (!childPattern.empty())
      add("child_pattern=" + childPattern);
    if (childLocationOffset != DefaultLocationOffset)
      add("child_location=" + std::to_string(childLocationOffset));
    if (childVersion != DefaultVersion)
      add("child_version=" + std::to_string(childVersion));
    if (!childMime.empty())
      add("child_mime=" + childMime);

    return b.str();
  }

  static Tag parse(const std::string &tagStr) {
    Tag r;
    std::string s = tagStr;
    s.erase(0, s.find_first_not_of(" \t\r\n"));
    if (s.size() >= 2 && s[0] == '/' && s[1] == '/') {
      s = s.substr(2);
      s.erase(0, s.find_first_not_of(" \t\r\n"));
    }
    if (s.size() >= 3 && s.substr(0, 3) == "mm:") {
      s = s.substr(3);
      s.erase(0, s.find_first_not_of(" \t\r\n"));
    }
    if (s.empty())
      return r;

    auto parts = splitTag(s);
    for (const auto &p : parts) {
      auto trimmed = trim(p);
      if (trimmed.empty())
        continue;

      std::string k, v;
      auto eqPos = trimmed.find('=');
      if (eqPos != std::string::npos) {
        k = trim(trimmed.substr(0, eqPos));
        v = trim(trimmed.substr(eqPos + 1));
      } else {
        k = trim(trimmed);
      }

      if (v.size() >= 2 && v.front() == '"' && v.back() == '"') {
        v = unquote(v);
      }

      std::string lower = k;
      for (auto &c : lower)
        c = static_cast<char>(std::tolower(static_cast<unsigned char>(c)));

      if (lower == "name") {
        r.name = v;
      } else if (lower == "is_null") {
        r.isNull = true;
        r.nullable = true;
      } else if (lower == "example") {
        r.example = true;
      } else if (lower == "desc") {
        r.desc = v;
      } else if (lower == "type") {
        r.type = parseValueType(v);
      } else if (lower == "deprecated") {
        r.deprecated = true;
      } else if (lower == "nullable") {
        r.nullable = true;
      } else if (lower == "allow_empty") {
        r.allowEmpty = true;
      } else if (lower == "unique") {
        r.unique = true;
      } else if (lower == "default_val") {
        r.default_val = v;
      } else if (lower == "min") {
        r.min = v;
      } else if (lower == "max") {
        r.max = v;
      } else if (lower == "size") {
        r.size = std::stoi(v);
      } else if (lower == "enums") {
        r.type = ValueType::Enums;
        r.enums = v;
      } else if (lower == "pattern") {
        r.pattern = v;
      } else if (lower == "location") {
        r.locationOffset = std::stoi(v);
      } else if (lower == "version") {
        r.version = std::stoi(v);
      } else if (lower == "mime") {
        r.mime = v;
        r.type = ValueType::Media;
      } else if (lower == "child_desc") {
        r.childDesc = v;
      } else if (lower == "child_type") {
        r.childType = parseValueType(v);
      } else if (lower == "child_raw") {
      } else if (lower == "child_nullable") {
        r.childNullable = true;
      } else if (lower == "child_allow_empty") {
        r.childAllowEmpty = true;
      } else if (lower == "child_unique") {
        r.childUnique = true;
      } else if (lower == "child_default_val") {
        r.child_default_val = v;
      } else if (lower == "child_min") {
        r.childMin = v;
      } else if (lower == "child_max") {
        r.childMax = v;
      } else if (lower == "child_size") {
        r.childSize = std::stoi(v);
      } else if (lower == "child_enums") {
        r.child_enums = v;
        r.childType = ValueType::Enums;
      } else if (lower == "child_pattern") {
        r.childPattern = v;
      } else if (lower == "child_location") {
        r.childLocationOffset = std::stoi(v);
      } else if (lower == "child_version") {
        r.childVersion = std::stoi(v);
      } else if (lower == "child_mime") {
        r.childMime = v;
        r.childType = ValueType::Media;
      }
    }
    return r;
  }

  std::vector<uint8_t> bytes() const {
    std::vector<uint8_t> bs;
    auto writeByte = [&](uint8_t b) { bs.push_back(b); };
    auto writeStr = [&](const std::string &s) {
      bs.insert(bs.end(), s.begin(), s.end());
    };

    if (example)
      writeByte(static_cast<uint8_t>(KExample | 1));
    if (isNull)
      writeByte(static_cast<uint8_t>(KIsNull | 1));
    if (nullable && !isInherit && !isNull)
      writeByte(static_cast<uint8_t>(KNullable | 1));

    if (!desc.empty() && !isInherit)
      encodeString(&bs, KDesc, desc);
    if (type != ValueType::Unknown && !isInherit) {
      bool skip = (type == ValueType::Str || type == ValueType::Bytes ||
                   type == ValueType::I || type == ValueType::F64 ||
                   type == ValueType::Bool || type == ValueType::Obj ||
                   type == ValueType::Vec);
      if (!(skip || (type == ValueType::Arr && size > 0) ||
            (type == ValueType::Enums && !enums.empty()) ||
            (type == ValueType::Media && !mime.empty()))) {
        writeByte(static_cast<uint8_t>(KType));
        writeByte(static_cast<uint8_t>(type));
      }
    }
    if (deprecated && !isInherit)
      writeByte(static_cast<uint8_t>(KDeprecated | 1));
    if (allowEmpty && !isInherit)
      writeByte(static_cast<uint8_t>(KAllowEmpty | 1));
    if (unique && !isInherit)
      writeByte(static_cast<uint8_t>(KUnique | 1));
    if (!default_val.empty() && !isInherit)
      encodeString(&bs, KDefaultVal, default_val);
    if (!min.empty() && !isInherit)
      encodeString(&bs, KMin, min);
    if (!max.empty() && !isInherit)
      encodeString(&bs, KMax, max);
    if (size != 0 && !isInherit)
      encodeU64(&bs, KSize, static_cast<uint64_t>(size));
    if (!enums.empty() && !isInherit)
      encodeString(&bs, KEnums, enums);
    if (!pattern.empty() && !isInherit)
      encodeString(&bs, KPattern, pattern);
    if (locationOffset != 0 && !isInherit) {
      std::string v = std::to_string(locationOffset);
      writeByte(static_cast<uint8_t>(KLocation) |
                static_cast<uint8_t>(v.size()));
      writeStr(v);
    }
    if (version != DefaultVersion && !isInherit)
      encodeU64(&bs, KVersion, static_cast<uint64_t>(version));
    if (!mime.empty() && !isInherit)
      encodeU64(&bs, KMime, static_cast<uint64_t>(parseMime(mime)));

    if (!childDesc.empty())
      encodeString(&bs, KChildDesc, childDesc);
    if (childType != ValueType::Unknown) {
      bool childSkip =
          (childType == ValueType::Str || childType == ValueType::I ||
           childType == ValueType::F64 || childType == ValueType::Bool ||
           childType == ValueType::Obj || childType == ValueType::Vec);
      if (!(childSkip || (childType == ValueType::Arr && childSize > 0) ||
            (childType == ValueType::Enums && !child_enums.empty()) ||
            (childType == ValueType::Media && !childMime.empty()))) {
        writeByte(static_cast<uint8_t>(KChildType));
        writeByte(static_cast<uint8_t>(childType));
      }
    }
    if (childNullable)
      writeByte(static_cast<uint8_t>(KChildNullable | 1));
    if (childAllowEmpty)
      writeByte(static_cast<uint8_t>(KChildAllowEmpty | 1));
    if (childUnique)
      writeByte(static_cast<uint8_t>(KChildUnique | 1));
    if (!child_default_val.empty())
      encodeString(&bs, KChildDefaultVal, child_default_val);
    if (!childMin.empty())
      encodeString(&bs, KChildMin, childMin);
    if (!childMax.empty())
      encodeString(&bs, KChildMax, childMax);
    if (childSize != 0)
      encodeU64(&bs, KChildSize, static_cast<uint64_t>(childSize));
    if (!child_enums.empty())
      encodeString(&bs, KChildEnums, child_enums);
    if (!childPattern.empty())
      encodeString(&bs, KChildPattern, childPattern);
    if (childLocationOffset != DefaultLocationOffset) {
      std::string v = std::to_string(childLocationOffset);
      writeByte(static_cast<uint8_t>(KChildLocation) |
                static_cast<uint8_t>(v.size()));
      writeStr(v);
    }
    if (childVersion != DefaultVersion)
      encodeU64(&bs, KChildVersion, static_cast<uint64_t>(childVersion));
    if (!childMime.empty())
      encodeU64(&bs, KChildMime, static_cast<uint64_t>(parseMime(childMime)));

    if (more != 0)
      encodeU64(&bs, KMore, static_cast<uint64_t>(more));

    return bs;
  }

private:
  static std::vector<std::string> splitTag(const std::string &tag) {
    std::vector<std::string> parts;
    std::string current;
    bool inQuote = false;
    for (size_t i = 0; i < tag.size(); ++i) {
      char c = tag[i];
      if (c == '"')
        inQuote = !inQuote;
      if (c == ';' && !inQuote) {
        parts.push_back(current);
        current.clear();
      } else {
        current.push_back(c);
      }
    }
    if (!current.empty())
      parts.push_back(current);
    return parts;
  }

  static std::string trim(const std::string &s) {
    size_t start = s.find_first_not_of(" \t\r\n");
    if (start == std::string::npos)
      return "";
    size_t end = s.find_last_not_of(" \t\r\n");
    return s.substr(start, end - start + 1);
  }

  static std::string unquote(const std::string &s) {
    if (s.size() < 2)
      return s;
    std::string result;
    for (size_t i = 1; i < s.size() - 1; ++i) {
      if (s[i] == '\\' && i + 1 < s.size() - 1) {
        switch (s[i + 1]) {
        case 'n':
          result.push_back('\n');
          ++i;
          break;
        case 'r':
          result.push_back('\r');
          ++i;
          break;
        case 't':
          result.push_back('\t');
          ++i;
          break;
        case '"':
          result.push_back('"');
          ++i;
          break;
        case '\\':
          result.push_back('\\');
          ++i;
          break;
        default:
          result.push_back(s[i]);
          break;
        }
      } else {
        result.push_back(s[i]);
      }
    }
    return result;
  }

  static void encodeString(std::vector<uint8_t> *bs, TagKey sign,
                           const std::string &val) {
    size_t l = val.size();
    if (l <= 5) {
      bs->push_back(static_cast<uint8_t>(sign) | static_cast<uint8_t>(l));
      bs->insert(bs->end(), val.begin(), val.end());
    } else if (l <= 255) {
      bs->push_back(static_cast<uint8_t>(sign) | 6);
      bs->push_back(static_cast<uint8_t>(l));
      bs->insert(bs->end(), val.begin(), val.end());
    } else {
      bs->push_back(static_cast<uint8_t>(sign) | 7);
      bs->push_back(static_cast<uint8_t>(l >> 8));
      bs->push_back(static_cast<uint8_t>(l));
      bs->insert(bs->end(), val.begin(), val.end());
    }
  }

  static void encodeU64(std::vector<uint8_t> *bs, TagKey sign, uint64_t uv) {
    constexpr uint64_t Max1Byte = 0xFF;
    constexpr uint64_t Max2Byte = 0xFFFF;
    constexpr uint64_t Max3Byte = 0xFFFFFF;
    constexpr uint64_t Max4Byte = 0xFFFFFFFF;
    constexpr uint64_t Max5Byte = 0xFFFFFFFFFF;
    constexpr uint64_t Max6Byte = 0xFFFFFFFFFFFF;
    constexpr uint64_t Max7Byte = 0xFFFFFFFFFFFFFF;

    if (uv <= Max1Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 0);
      bs->push_back(static_cast<uint8_t>(uv));
    } else if (uv <= Max2Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 1);
      bs->push_back(static_cast<uint8_t>(uv >> 8));
      bs->push_back(static_cast<uint8_t>(uv));
    } else if (uv <= Max3Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 2);
      bs->push_back(static_cast<uint8_t>(uv >> 16));
      bs->push_back(static_cast<uint8_t>(uv >> 8));
      bs->push_back(static_cast<uint8_t>(uv));
    } else if (uv <= Max4Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 3);
      bs->push_back(static_cast<uint8_t>(uv >> 24));
      bs->push_back(static_cast<uint8_t>(uv >> 16));
      bs->push_back(static_cast<uint8_t>(uv >> 8));
      bs->push_back(static_cast<uint8_t>(uv));
    } else if (uv <= Max5Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 4);
      bs->push_back(static_cast<uint8_t>(uv >> 32));
      bs->push_back(static_cast<uint8_t>(uv >> 24));
      bs->push_back(static_cast<uint8_t>(uv >> 16));
      bs->push_back(static_cast<uint8_t>(uv >> 8));
      bs->push_back(static_cast<uint8_t>(uv));
    } else if (uv <= Max6Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 5);
      bs->push_back(static_cast<uint8_t>(uv >> 40));
      bs->push_back(static_cast<uint8_t>(uv >> 32));
      bs->push_back(static_cast<uint8_t>(uv >> 24));
      bs->push_back(static_cast<uint8_t>(uv >> 16));
      bs->push_back(static_cast<uint8_t>(uv >> 8));
      bs->push_back(static_cast<uint8_t>(uv));
    } else if (uv <= Max7Byte) {
      bs->push_back(static_cast<uint8_t>(sign) | 6);
      bs->push_back(static_cast<uint8_t>(uv >> 48));
      bs->push_back(static_cast<uint8_t>(uv >> 40));
      bs->push_back(static_cast<uint8_t>(uv >> 32));
      bs->push_back(static_cast<uint8_t>(uv >> 24));
      bs->push_back(static_cast<uint8_t>(uv >> 16));
      bs->push_back(static_cast<uint8_t>(uv >> 8));
      bs->push_back(static_cast<uint8_t>(uv));
    } else {
      bs->push_back(static_cast<uint8_t>(sign) | 7);
      for (int i = 7; i >= 0; --i)
        bs->push_back(static_cast<uint8_t>(uv >> (i * 8)));
    }
  }

  static int parseMime(const std::string &mime) {
    std::string key = mime;
    std::transform(key.begin(), key.end(), key.begin(), ::tolower);
    key.erase(0, key.find_first_not_of(" \t"));
    key.erase(key.find_last_not_of(" \t") + 1);

    if (key == "image/jpeg" || key == "image/jpg")
      return 1;
    if (key == "image/png")
      return 2;
    if (key == "image/gif")
      return 3;
    if (key == "image/webp")
      return 4;
    if (key == "image/svg+xml")
      return 5;
    if (key == "image/avif")
      return 6;
    if (key == "image/bmp")
      return 7;
    if (key == "image/x-icon")
      return 8;
    if (key == "image/tiff")
      return 9;
    if (key == "image/heic")
      return 10;
    if (key == "image/heif")
      return 11;
    if (key == "text/plain")
      return 12;
    if (key == "text/html")
      return 13;
    if (key == "text/css")
      return 14;
    if (key == "text/javascript")
      return 15;
    if (key == "application/json")
      return 16;
    if (key == "text/csv")
      return 17;
    if (key == "text/markdown")
      return 18;
    if (key == "application/pdf")
      return 19;
    if (key == "application/zip")
      return 20;
    if (key == "application/gzip")
      return 21;
    if (key == "application/x-tar")
      return 22;
    if (key ==
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
      return 23;
    if (key == "application/"
               "vnd.openxmlformats-officedocument.wordprocessingml.document")
      return 24;
    if (key == "application/octet-stream")
      return 25;
    if (key == "video/mp4")
      return 26;
    if (key == "video/webm")
      return 27;
    if (key == "video/mov")
      return 28;
    if (key == "audio/mpeg")
      return 29;
    if (key == "audio/wav")
      return 30;
    if (key == "audio/flac")
      return 31;
    if (key == "font/woff2")
      return 32;
    if (key == "font/ttf")
      return 33;
    return 0;
  }
};

inline Tag mergeTag(const Tag *dst, const Tag *src) {
  Tag r;
  if (src == nullptr && dst == nullptr)
    return r;
  if (src == nullptr)
    return *dst;
  if (dst == nullptr)
    return *src;

  r = *dst;

  if (src->isNull)
    r.isNull = src->isNull;
  if (src->example)
    r.example = src->example;
  if (!src->desc.empty())
    r.desc = src->desc;
  if (src->type != ValueType::Unknown)
    r.type = src->type;
  if (src->deprecated)
    r.deprecated = true;
  if (src->nullable)
    r.nullable = true;
  if (src->allowEmpty)
    r.allowEmpty = true;
  if (src->unique)
    r.unique = true;
  if (!src->default_val.empty())
    r.default_val = src->default_val;
  if (!src->min.empty())
    r.min = src->min;
  if (!src->max.empty())
    r.max = src->max;
  if (src->size != 0)
    r.size = src->size;
  if (!src->enums.empty())
    r.enums = src->enums;
  if (!src->pattern.empty())
    r.pattern = src->pattern;
  if (src->locationOffset != 0)
    r.locationOffset = src->locationOffset;
  if (src->version != DefaultVersion)
    r.version = src->version;
  if (!src->mime.empty())
    r.mime = src->mime;
  if (!src->childDesc.empty())
    r.childDesc = src->childDesc;
  if (src->childType != ValueType::Unknown)
    r.childType = src->childType;
  if (src->childNullable)
    r.childNullable = true;
  if (src->childAllowEmpty)
    r.childAllowEmpty = true;
  if (src->childUnique)
    r.childUnique = true;
  if (!src->child_default_val.empty())
    r.child_default_val = src->child_default_val;
  if (!src->childMin.empty())
    r.childMin = src->childMin;
  if (!src->childMax.empty())
    r.childMax = src->childMax;
  if (src->childSize != 0)
    r.childSize = src->childSize;
  if (!src->child_enums.empty())
    r.child_enums = src->child_enums;
  if (!src->childPattern.empty())
    r.childPattern = src->childPattern;
  if (src->childLocationOffset != DefaultLocationOffset)
    r.childLocationOffset = src->childLocationOffset;
  if (src->childVersion != DefaultVersion)
    r.childVersion = src->childVersion;
  if (!src->childMime.empty())
    r.childMime = src->childMime;

  return r;
}

} // namespace ir
} // namespace mmc

#endif