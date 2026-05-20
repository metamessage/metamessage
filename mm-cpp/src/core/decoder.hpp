#ifndef MMCPP_CORE_DECODER_HPP
#define MMCPP_CORE_DECODER_HPP

#include "constants.hpp"
#include "../ir/ast.hpp"
#include <vector>
#include <cstdint>
#include <string>
#include <cstring>
#include <stdexcept>
#include <iostream>

namespace mmc {
namespace core {

using namespace mmc::ir;

class Decoder {
public:
    Decoder() : data_(nullptr), size_(0), offset_(0) {}

    std::shared_ptr<ir::Node> decode(const std::vector<uint8_t>& encoded) {
        data_ = encoded.data();
        size_ = encoded.size();
        offset_ = 0;
        return decodeNode(nullptr);
    }

private:
    const uint8_t* data_;
    size_t size_;
    size_t offset_;

    uint8_t readByte() {
        if (offset_ >= size_) throw std::runtime_error("unexpected EOF");
        return data_[offset_++];
    }

    std::vector<uint8_t> readBytes(size_t n) {
        if (offset_ + n > size_) throw std::runtime_error("unexpected EOF");
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

    std::shared_ptr<ir::Node> decodeNode(const ir::Tag* parentTag) {
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

        ir::Tag tag;
        tag.isInherit = true;

        int remaining = l;
        while (remaining > 0) {
            int n = decodeTagBytes(tag);
            if (n == 0 || n > remaining) break;
            remaining -= n;
        }

        auto node = decodeNode(&tag);
        return node;
    }

    int decodeTagBytes(ir::Tag& tag) {
        uint8_t b = readByte();
        uint8_t suffix = getSuffix(b);

        // Boolean flags have suffix == 1 (KKey | 1)
        // String values have suffix indicating length
        // Type key has suffix == 0, followed by type byte

        switch (suffix) {
            case 1: {
                if (b == (KIsNull | 1)) {
                    tag.isNull = true;
                    return 1;
                }
                if (b == (KExample | 1)) {
                    tag.example = true;
                    return 1;
                }
                if (b == (KNullable | 1)) {
                    tag.nullable = true;
                    return 1;
                }
                if (b == (KRaw | 1)) {
                    tag.raw = true;
                    return 1;
                }
                if (b == (KAllowEmpty | 1)) {
                    tag.allowEmpty = true;
                    return 1;
                }
                if (b == (KUnique | 1)) {
                    tag.unique = true;
                    return 1;
                }
                if (b == (KChildRaw | 1)) {
                    tag.childRaw = true;
                    return 1;
                }
                if (b == (KChildNullable | 1)) {
                    tag.childNullable = true;
                    return 1;
                }
                if (b == (KChildAllowEmpty | 1)) {
                    tag.childAllowEmpty = true;
                    return 1;
                }
                if (b == (KChildUnique | 1)) {
                    tag.childUnique = true;
                    return 1;
                }
                return 1;
            }
            case 0: {
                // Type key: KType, KChildType
                if (b == (KType)) {
                    tag.type = static_cast<ir::ValueType>(readByte());
                    return 2;
                }
                if (b == (KChildType)) {
                    tag.childType = static_cast<ir::ValueType>(readByte());
                    return 2;
                }
                // Size keys with length encoding
                if (b >= KSize && b < KSize + 8) {
                    tag.size = static_cast<int>(decodeTagU64());
                    return 2;
                }
                if (b >= KChildSize && b < KChildSize + 8) {
                    tag.childSize = static_cast<int>(decodeTagU64());
                    return 2;
                }
                if (b >= KVersion && b < KVersion + 8) {
                    tag.version = static_cast<int>(decodeTagU64());
                    return 2;
                }
                if (b >= KChildVersion && b < KChildVersion + 8) {
                    tag.childVersion = static_cast<int>(decodeTagU64());
                    return 2;
                }
                return 1;
            }
            default: {
                if (b >= KDesc && b < KDesc + 8) {
                    tag.desc = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.desc.size());
                }
                if (b >= KDefault && b < KDefault + 8) {
                    tag.defaultVal = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.defaultVal.size());
                }
                if (b >= KMin && b < KMin + 8) {
                    tag.min = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.min.size());
                }
                if (b >= KMax && b < KMax + 8) {
                    tag.max = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.max.size());
                }
                if (b >= KEnum && b < KEnum + 8) {
                    tag.enumVal = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.enumVal.size());
                }
                if (b >= KPattern && b < KPattern + 8) {
                    tag.pattern = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.pattern.size());
                }
                if (b >= KLocation && b < KLocation + 8) {
                    tag.locationOffset = std::stoi(decodeTagString(suffix));
                    return static_cast<int>(1 + 1);
                }
                if (b >= KMime && b < KMime + 8) {
                    tag.mime = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.mime.size());
                }
                if (b >= KChildDesc && b < KChildDesc + 8) {
                    tag.childDesc = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childDesc.size());
                }
                if (b >= KChildDefault && b < KChildDefault + 8) {
                    tag.childDefault = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childDefault.size());
                }
                if (b >= KChildMin && b < KChildMin + 8) {
                    tag.childMin = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childMin.size());
                }
                if (b >= KChildMax && b < KChildMax + 8) {
                    tag.childMax = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childMax.size());
                }
                if (b >= KChildEnum && b < KChildEnum + 8) {
                    tag.childEnum = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childEnum.size());
                }
                if (b >= KChildPattern && b < KChildPattern + 8) {
                    tag.childPattern = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childPattern.size());
                }
                if (b >= KChildLocation && b < KChildLocation + 8) {
                    tag.childLocationOffset = std::stoi(decodeTagString(suffix));
                    return static_cast<int>(1 + 1);
                }
                if (b >= KChildMime && b <= 255) {
                    tag.childMime = decodeTagString(suffix);
                    return static_cast<int>(1 + tag.childMime.size());
                }
                return 1;
            }
        }
    }

    std::string decodeTagString(uint8_t suffix) {
        int len = suffix;
        if (len > 5) len = readByte();
        auto bytes = readBytes(static_cast<size_t>(len));
        return std::string(bytes.begin(), bytes.end());
    }

    uint64_t decodeTagU64() {
        uint8_t b = readByte();
        int byteLen = intLen(b);
        if (byteLen == 0) return 0;
        return readUInt(byteLen);
    }

    std::shared_ptr<ir::Node> decodeSimple(uint8_t b, const ir::Tag* parentTag) {
        auto val = ir::makeValue();
        auto* tag = val->getTag();

        if (parentTag) {
            *tag = *parentTag;
            tag->inherit(*parentTag);
        }

        SimpleValue sv = static_cast<SimpleValue>(getSuffix(b));

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
                break;
            case SimpleNullInt:
                tag->type = ir::ValueType::I;
                tag->isNull = true;
                break;
            case SimpleNullFloat:
                tag->type = ir::ValueType::F64;
                tag->isNull = true;
                break;
            case SimpleNullString:
                tag->type = ir::ValueType::Str;
                tag->isNull = true;
                break;
            case SimpleNullBytes:
                tag->type = ir::ValueType::Bytes;
                tag->isNull = true;
                break;
            default:
                break;
        }
        return val;
    }

    std::shared_ptr<ir::Node> decodeInt(uint8_t b, const ir::Tag* parentTag, bool positive) {
        auto val = ir::makeValue();
        auto* tag = val->getTag();

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

        if (positive && parentTag) {
            tag->type = parentTag->type;
            if (tag->type == ir::ValueType::Unknown)
                tag->type = ir::ValueType::I;
        }
        if (!positive && parentTag) {
            tag->type = parentTag->type;
            if (tag->type == ir::ValueType::Unknown)
                tag->type = ir::ValueType::I;
        }

        if (tag->type == ir::ValueType::Unknown)
            tag->type = ir::ValueType::I;

        val->text = std::to_string(positive ? uv : -static_cast<int64_t>(uv));
        return val;
    }

    std::shared_ptr<ir::Node> decodeFloat(uint8_t b, const ir::Tag* parentTag) {
        auto val = ir::makeValue();
        auto* tag = val->getTag();

        if (parentTag) {
            *tag = *parentTag;
            tag->inherit(*parentTag);
        }

        int byteLen = floatLen(b);
        uint64_t bits = 0;
        if (byteLen == 0) {
            bits = getSuffix(b);
        } else {
            bits = readUInt(byteLen);
        }

        double d;
        std::memcpy(&d, &bits, sizeof(d));
        val->text = std::to_string(d);
        tag->type = ir::ValueType::F64;
        return val;
    }

    std::shared_ptr<ir::Node> decodeString(uint8_t b, const ir::Tag* parentTag) {
        auto val = ir::makeValue();
        auto* tag = val->getTag();

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
        tag->type = ir::ValueType::Str;
        return val;
    }

    std::shared_ptr<ir::Node> decodeBytes(uint8_t b, const ir::Tag* parentTag) {
        auto val = ir::makeValue();
        auto* tag = val->getTag();

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
        val->text = std::string(rawBytes.begin(), rawBytes.end());
        tag->type = ir::ValueType::Bytes;
        return val;
    }

    std::shared_ptr<ir::Node> decodeContainer(uint8_t b, const ir::Tag* parentTag) {
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

    std::shared_ptr<ir::Node> decodeContainerArray(size_t len, const ir::Tag* parentTag) {
        auto arr = ir::makeArray();
        size_t startOffset = offset_;

        // Decode key array first (to know all item lengths)
        // Then decode key-value pairs
        // Simple approach: just recurse
        while (offset_ - startOffset < len) {
            auto item = decodeNode(parentTag);
            arr->items.push_back(item);
        }

        if (parentTag) {
            *arr->getTag() = *parentTag;
        }

        return arr;
    }

    std::shared_ptr<ir::Node> decodeContainerObject(size_t len, const ir::Tag* parentTag) {
        auto obj = ir::makeObject();
        size_t startOffset = offset_;
        size_t endOffset = startOffset + len;

        // First decode the key array
        auto keyArrayNode = decodeNode(parentTag);
        auto keyArray = std::dynamic_pointer_cast<ir::Array>(keyArrayNode);

        // Then decode values and pair them with keys
        if (keyArray) {
            for (auto& keyItem : keyArray->items) {
                if (offset_ >= endOffset) break;
                auto keyVal = std::dynamic_pointer_cast<ir::Value>(keyItem);
                auto valNode = decodeNode(parentTag);
                if (keyVal && !keyVal->text.empty()) {
                    ir::Field field(keyVal->text, valNode);
                    obj->fields.push_back(std::move(field));
                }
            }
        } else {
            // Fallback: decode remaining as key-value pairs
            while (offset_ < endOffset) {
                auto keyNode = decodeNode(parentTag);
                if (offset_ >= endOffset) break;
                auto valNode = decodeNode(parentTag);
                auto key = std::dynamic_pointer_cast<ir::Value>(keyNode);
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

        return obj;
    }
};

} // namespace core
} // namespace mmc

#endif