#ifndef MMCPP_CORE_ENCODER_HPP
#define MMCPP_CORE_ENCODER_HPP

#include "constants.hpp"
#include "../ir/ast.hpp"
#include <vector>
#include <cstdint>
#include <string>
#include <stdexcept>

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
            case ir::NodeType::Object:
                n = encodeNodeObject(std::static_pointer_cast<ir::Object>(node));
                break;
            case ir::NodeType::Array:
                n = encodeNodeArray(std::static_pointer_cast<ir::Array>(node));
                break;
            case ir::NodeType::Value:
                n = encodeNodeValue(std::static_pointer_cast<ir::Value>(node));
                break;
            default:
                throw std::runtime_error("unsupported node type");
        }

        std::vector<uint8_t> out(buf_.begin() + offset_ - n, buf_.begin() + offset_);
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
            if (newSize > maxCap_) newSize = maxCap_;
            if (newSize < offset_ + needed) newSize = offset_ + needed;
            buf_.resize(newSize);
        }
    }

    uint32_t writeByte(uint8_t b) {
        ensure(1);
        buf_[offset_++] = b;
        return 1;
    }

    uint32_t writeBytes(const std::vector<uint8_t>& data) {
        ensure(data.size());
        std::copy(data.begin(), data.end(), buf_.begin() + offset_);
        offset_ += data.size();
        return static_cast<uint32_t>(data.size());
    }

    uint32_t writeString(const std::string& s) {
        ensure(s.size());
        std::copy(s.begin(), s.end(), buf_.begin() + offset_);
        offset_ += s.size();
        return static_cast<uint32_t>(s.size());
    }

    uint32_t encodeSimple(SimpleValue value) {
        return writeByte(static_cast<uint8_t>(Prefix::Simple) | static_cast<uint8_t>(value));
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
        } else if (uv <= 0xFFFFFFFF) {
            prefix = static_cast<uint8_t>(sign) | IntLen4Byte;
            len = 4;
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

    uint32_t encodeU64(uint64_t uv) {
        return encodeInt(Prefix::PositiveInt, uv);
    }

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

    uint32_t encodeFloat(const std::string& text) {
        double val = std::stod(text);
        uint64_t bits;
        std::memcpy(&bits, &val, sizeof(bits));
        uint8_t prefix = static_cast<uint8_t>(Prefix::Float) | 7;
        writeByte(prefix);
        for (int i = 7; i >= 0; --i)
            writeByte(static_cast<uint8_t>(bits >> (i * 8)));
        return 9;
    }

    uint32_t encodeString(const std::string& s) {
        size_t l = s.size();
        if (l <= 30) {
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

    uint32_t encodeBytes(const std::vector<uint8_t>& data) {
        size_t l = data.size();
        if (l <= 30) {
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

    uint32_t encodeArray(const std::vector<uint8_t>& data) {
        size_t l = data.size();
        if (l <= 15) {
            writeByte(static_cast<uint8_t>(Prefix::Container) |
                      ContainerArray |
                      static_cast<uint8_t>(l));
            writeBytes(data);
            return static_cast<uint32_t>(1 + l);
        } else if (l <= 255) {
            writeByte(static_cast<uint8_t>(Prefix::Container) |
                      ContainerArray |
                      ContainerLen1Byte);
            writeByte(static_cast<uint8_t>(l));
            writeBytes(data);
            return static_cast<uint32_t>(2 + l);
        } else {
            writeByte(static_cast<uint8_t>(Prefix::Container) |
                      ContainerArray |
                      ContainerLen2Byte);
            writeByte(static_cast<uint8_t>(l >> 8));
            writeByte(static_cast<uint8_t>(l));
            writeBytes(data);
            return static_cast<uint32_t>(3 + l);
        }
    }

    uint32_t encodeObject(const std::vector<uint8_t>& data) {
        size_t l = data.size();
        if (l <= 15) {
            writeByte(static_cast<uint8_t>(Prefix::Container) |
                      ContainerObject |
                      static_cast<uint8_t>(l));
            writeBytes(data);
            return static_cast<uint32_t>(1 + l);
        } else if (l <= 255) {
            writeByte(static_cast<uint8_t>(Prefix::Container) |
                      ContainerObject |
                      ContainerLen1Byte);
            writeByte(static_cast<uint8_t>(l));
            writeBytes(data);
            return static_cast<uint32_t>(2 + l);
        } else {
            writeByte(static_cast<uint8_t>(Prefix::Container) |
                      ContainerObject |
                      ContainerLen2Byte);
            writeByte(static_cast<uint8_t>(l >> 8));
            writeByte(static_cast<uint8_t>(l));
            writeBytes(data);
            return static_cast<uint32_t>(3 + l);
        }
    }

    uint32_t encodeTagPayload(uint64_t uv) {
        if (uv <= 0xFF) {
            writeByte(TagPayload1Byte);
            writeByte(static_cast<uint8_t>(uv));
            return 2;
        } else {
            writeByte(TagPayload2Byte);
            writeByte(static_cast<uint8_t>(uv >> 8));
            writeByte(static_cast<uint8_t>(uv));
            return 3;
        }
    }

    uint32_t encodeT(const std::vector<uint8_t>& tagBytes) {
        if (tagBytes.empty()) return 0;

        size_t l = tagBytes.size();
        if (l <= 30) {
            writeByte(static_cast<uint8_t>(Prefix::Tag) | static_cast<uint8_t>(l));
            writeBytes(tagBytes);
            return static_cast<uint32_t>(1 + l);
        } else if (l <= 255) {
            writeByte(static_cast<uint8_t>(Prefix::Tag) | TagLen1Byte);
            writeByte(static_cast<uint8_t>(l));
            writeBytes(tagBytes);
            return static_cast<uint32_t>(2 + l);
        } else {
            writeByte(static_cast<uint8_t>(Prefix::Tag) | TagLen2Byte);
            writeByte(static_cast<uint8_t>(l >> 8));
            writeByte(static_cast<uint8_t>(l));
            writeBytes(tagBytes);
            return static_cast<uint32_t>(3 + l);
        }
    }

    uint32_t encodeTag(const std::vector<uint8_t>& payload, const std::vector<uint8_t>& tagData) {
        std::vector<uint8_t> combined = payload;
        combined.insert(combined.end(), tagData.begin(), tagData.end());
        return encodeTagPayload(combined.size()) | writeBytes(combined);
    }

    uint32_t encodeComment(const std::vector<uint8_t>& payload, const ir::Tag* tag) {
        if (tag == nullptr) return static_cast<uint32_t>(payload.size());
        auto tagBytes = tag->bytes();
        if (tagBytes.empty()) return static_cast<uint32_t>(payload.size());

        uint32_t ns = encodeT(tagBytes);
        if (ns == 0) return static_cast<uint32_t>(payload.size());

        size_t prevOffset = offset_ - ns;
        std::vector<uint8_t> tagData(buf_.begin() + prevOffset, buf_.begin() + prevOffset + ns);
        offset_ = prevOffset;

        return encodeTag(payload, tagData);
    }

    uint32_t encodeNodeObject(std::shared_ptr<ir::Object> obj) {
        std::vector<uint8_t> bufKey, buf;

        auto* tag = obj->getTag();

        for (auto& field : obj->fields) {
            uint32_t n = 0;
            switch (field.value->getType()) {
                case ir::NodeType::Object:
                    n = encodeNodeObject(std::static_pointer_cast<ir::Object>(field.value));
                    break;
                case ir::NodeType::Array:
                    n = encodeNodeArray(std::static_pointer_cast<ir::Array>(field.value));
                    break;
                case ir::NodeType::Value:
                    n = encodeNodeValue(std::static_pointer_cast<ir::Value>(field.value));
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
        if (n1 == 0) return n;
        return n1;
    }

    uint32_t encodeNodeArray(std::shared_ptr<ir::Array> arr) {
        std::vector<uint8_t> buf;
        auto* tag = arr->getTag();

        for (auto& item : arr->items) {
            uint32_t n = 0;
            switch (item->getType()) {
                case ir::NodeType::Object:
                    n = encodeNodeObject(std::static_pointer_cast<ir::Object>(item));
                    break;
                case ir::NodeType::Array:
                    n = encodeNodeArray(std::static_pointer_cast<ir::Array>(item));
                    break;
                case ir::NodeType::Value:
                    n = encodeNodeValue(std::static_pointer_cast<ir::Value>(item));
                    break;
                default:
                    throw std::runtime_error("unsupported item type");
            }

            auto encodedSub = getEncodedBytes(n);
            buf.insert(buf.end(), encodedSub.begin(), encodedSub.end());
        }

        uint32_t n = encodeArray(buf);
        uint32_t n1 = encodeComment(getEncodedBytes(n), tag);
        if (n1 == 0) return n;
        return n1;
    }

    uint32_t encodeNodeValue(std::shared_ptr<ir::Value> val) {
        uint32_t n = 0;
        auto* tag = val->getTag();

        switch (tag->type) {
            case ir::ValueType::Datetime:
            case ir::ValueType::Date:
            case ir::ValueType::Time:
                n = encodeU64(0);
                break;

            case ir::ValueType::I:
            case ir::ValueType::I8:
            case ir::ValueType::I16:
            case ir::ValueType::I32:
            case ir::ValueType::I64:
                if (tag->isNull) n = encodeSimple(SimpleNullInt);
                else n = encodeInt64(std::stoll(val->text));
                break;

            case ir::ValueType::U:
            case ir::ValueType::U8:
            case ir::ValueType::U16:
            case ir::ValueType::U32:
            case ir::ValueType::U64:
                if (tag->isNull) n = 0;
                else n = encodeU64(std::stoull(val->text));
                break;

            case ir::ValueType::F32:
            case ir::ValueType::F64:
                if (tag->isNull) n = encodeSimple(SimpleNullFloat);
                else n = encodeFloat(val->text);
                break;

            case ir::ValueType::Decimal:
                if (tag->isNull) n = 0;
                else n = encodeFloat(val->text);
                break;

            case ir::ValueType::Bigint:
                if (tag->isNull) n = 0;
                else n = encodeU64(std::stoull(val->text));
                break;

            case ir::ValueType::Str:
            case ir::ValueType::Email:
            case ir::ValueType::Url:
                if (tag->isNull) n = encodeSimple(SimpleNullString);
                else n = encodeString(val->text);
                break;

            case ir::ValueType::Uuid:
            case ir::ValueType::Ip:
            case ir::ValueType::Bytes:
                if (tag->isNull) n = encodeSimple(SimpleNullBytes);
                else {
                    std::vector<uint8_t> bytes(val->text.begin(), val->text.end());
                    n = encodeBytes(bytes);
                }
                break;

            case ir::ValueType::Bool:
                if (tag->isNull) n = encodeSimple(SimpleNullBool);
                else n = encodeBool(val->text == "true" || val->text == "1");
                break;

            case ir::ValueType::Enum:
                if (tag->isNull) n = 0;
                else n = encodeU64(0);
                break;

            default:
                throw std::runtime_error("unsupported type: " + ir::valueTypeToString(tag->type));
        }

        uint32_t n1 = encodeComment(getEncodedBytes(n), tag);
        if (n1 == 0) return n;
        return n1;
    }

    std::vector<uint8_t> getEncodedBytes(uint32_t written) {
        return std::vector<uint8_t>(
            buf_.begin() + offset_ - written,
            buf_.begin() + offset_
        );
    }
};

} // namespace core
} // namespace mmc

#endif