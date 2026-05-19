#include "codec.hpp"
#include <cstring>

#if defined(__APPLE__) || defined(__linux__)
#include <endian.h>
#elif defined(_WIN32)
#include <stdlib.h>
#endif

namespace mmc {

static inline bool isLittleEndian() {
    const uint16_t n = 1;
    return *reinterpret_cast<const uint8_t*>(&n) == 1;
}

static constexpr uint8_t PrefixMask = 0xC0;
static constexpr uint8_t LengthMask = 0x3F;

static constexpr uint8_t PrefixSimple = 0x00;
static constexpr uint8_t PrefixPositiveInt = 0x40;
static constexpr uint8_t PrefixNegativeInt = 0x80;
static constexpr uint8_t PrefixFloat = 0xC0;
static constexpr uint8_t PrefixString = 0x18;
static constexpr uint8_t PrefixBytes = 0x10;
static constexpr uint8_t PrefixContainer = 0x20;
static constexpr uint8_t PrefixTag = 0x04;

static constexpr uint8_t SimpleTrue = 0x01;
static constexpr uint8_t SimpleFalse = 0x02;

size_t Encoder::encodeBool(bool value) {
    buffer.push_back(value ? SimpleTrue : SimpleFalse);
    return 1;
}

size_t Encoder::encodeInt64(int64_t value) {
    if (value >= 0) {
        return encodeU64(static_cast<uint64_t>(value));
    } else {
        buffer.push_back(PrefixNegativeInt);
        uint64_t uv = static_cast<uint64_t>(-(value + 1)) + 1;
        if (uv < 0x100) {
            buffer.push_back(static_cast<uint8_t>(uv));
            return 2;
        } else if (uv < 0x10000) {
            buffer.push_back(0x01);
            buffer.push_back(static_cast<uint8_t>(uv & 0xFF));
            buffer.push_back(static_cast<uint8_t>((uv >> 8) & 0xFF));
            return 3;
        } else if (uv < 0x1000000) {
            buffer.push_back(0x02);
            buffer.push_back(static_cast<uint8_t>(uv & 0xFF));
            buffer.push_back(static_cast<uint8_t>((uv >> 8) & 0xFF));
            buffer.push_back(static_cast<uint8_t>((uv >> 16) & 0xFF));
            return 4;
        } else {
            buffer.push_back(0x03);
            for (int i = 0; i < 8; i++) {
                buffer.push_back(static_cast<uint8_t>((uv >> (i * 8)) & 0xFF));
            }
            return 9;
        }
    }
}

size_t Encoder::encodeU64(uint64_t value) {
    buffer.push_back(PrefixPositiveInt);
    if (value < 0x100) {
        buffer.push_back(static_cast<uint8_t>(value));
        return 2;
    } else if (value < 0x10000) {
        buffer.push_back(0x01);
        buffer.push_back(static_cast<uint8_t>(value & 0xFF));
        buffer.push_back(static_cast<uint8_t>((value >> 8) & 0xFF));
        return 3;
    } else if (value < 0x1000000) {
        buffer.push_back(0x02);
        buffer.push_back(static_cast<uint8_t>(value & 0xFF));
        buffer.push_back(static_cast<uint8_t>((value >> 8) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((value >> 16) & 0xFF));
        return 4;
    } else if (value < 0x100000000ULL) {
        buffer.push_back(0x03);
        for (int i = 0; i < 4; i++) {
            buffer.push_back(static_cast<uint8_t>((value >> (i * 8)) & 0xFF));
        }
        return 5;
    } else {
        buffer.push_back(0x07);
        for (int i = 0; i < 8; i++) {
            buffer.push_back(static_cast<uint8_t>((value >> (i * 8)) & 0xFF));
        }
        return 9;
    }
}

size_t Encoder::encodeFloat(double value) {
    buffer.push_back(PrefixFloat);
    uint64_t bits;
    if (isLittleEndian()) {
        std::memcpy(&bits, &value, sizeof(double));
    } else {
        std::memcpy(&bits, &value, sizeof(double));
        bits = __builtin_bswap64(bits);
    }
    for (int i = 0; i < 8; i++) {
        buffer.push_back(static_cast<uint8_t>((bits >> (i * 8)) & 0xFF));
    }
    return 9;
}

size_t Encoder::encodeString(const std::string& value) {
    size_t len = value.size();
    buffer.push_back(PrefixString);
    if (len < 0x100) {
        buffer.push_back(static_cast<uint8_t>(len));
    } else if (len < 0x10000) {
        buffer.push_back(0x01);
        buffer.push_back(static_cast<uint8_t>(len & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 8) & 0xFF));
    } else {
        buffer.push_back(0x02);
        buffer.push_back(static_cast<uint8_t>(len & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 8) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 16) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 24) & 0xFF));
    }
    buffer.insert(buffer.end(), value.begin(), value.end());
    return 1 + (len < 0x100 ? 1 : len < 0x10000 ? 2 : 4) + len;
}

size_t Encoder::encodeBytes(const std::vector<uint8_t>& value) {
    size_t len = value.size();
    buffer.push_back(PrefixBytes);
    if (len < 0x100) {
        buffer.push_back(static_cast<uint8_t>(len));
    } else if (len < 0x10000) {
        buffer.push_back(0x01);
        buffer.push_back(static_cast<uint8_t>(len & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 8) & 0xFF));
    } else {
        buffer.push_back(0x02);
        buffer.push_back(static_cast<uint8_t>(len & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 8) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 16) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((len >> 24) & 0xFF));
    }
    buffer.insert(buffer.end(), value.begin(), value.end());
    return 1 + (len < 0x100 ? 1 : len < 0x10000 ? 2 : 4) + len;
}

size_t Encoder::encodeArray(const std::vector<std::vector<uint8_t>>& items) {
    size_t payloadSize = 0;
    for (const auto& item : items) {
        payloadSize += item.size();
    }

    buffer.push_back(PrefixContainer);
    if (payloadSize < 0x100) {
        buffer.push_back(static_cast<uint8_t>(payloadSize));
    } else if (payloadSize < 0x10000) {
        buffer.push_back(0x01);
        buffer.push_back(static_cast<uint8_t>(payloadSize & 0xFF));
        buffer.push_back(static_cast<uint8_t>((payloadSize >> 8) & 0xFF));
    } else {
        buffer.push_back(0x02);
        buffer.push_back(static_cast<uint8_t>(payloadSize & 0xFF));
        buffer.push_back(static_cast<uint8_t>((payloadSize >> 8) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((payloadSize >> 16) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((payloadSize >> 24) & 0xFF));
    }

    for (const auto& item : items) {
        buffer.insert(buffer.end(), item.begin(), item.end());
    }

    return 1 + (payloadSize < 0x100 ? 1 : payloadSize < 0x10000 ? 2 : 4) + payloadSize;
}

size_t Encoder::encodeObject(const std::vector<std::pair<std::string, std::vector<uint8_t>>>& fields) {
    std::vector<std::vector<uint8_t>> keyItems;
    size_t totalKeysSize = 0;
    for (const auto& [key, _] : fields) {
        Encoder keyEnc;
        keyEnc.reset();
        keyEnc.encodeString(key);
        keyItems.push_back(keyEnc.result());
        totalKeysSize += keyItems.back().size();
    }

    std::vector<std::vector<uint8_t>> valueItems;
    for (const auto& [_, value] : fields) {
        valueItems.push_back(value);
    }

    size_t totalValuesSize = 0;
    for (const auto& v : valueItems) {
        totalValuesSize += v.size();
    }

    Encoder arrayEnc;
    arrayEnc.reset();
    arrayEnc.encodeArray(keyItems);

    size_t headerSize = 1 + (totalValuesSize < 0x100 ? 1 : totalValuesSize < 0x10000 ? 2 : 4);
    size_t objectSize = headerSize + totalValuesSize;

    buffer.push_back(PrefixContainer);
    if (objectSize < 0x100) {
        buffer.push_back(static_cast<uint8_t>(objectSize));
    } else if (objectSize < 0x10000) {
        buffer.push_back(0x01);
        buffer.push_back(static_cast<uint8_t>(objectSize & 0xFF));
        buffer.push_back(static_cast<uint8_t>((objectSize >> 8) & 0xFF));
    } else {
        buffer.push_back(0x02);
        buffer.push_back(static_cast<uint8_t>(objectSize & 0xFF));
        buffer.push_back(static_cast<uint8_t>((objectSize >> 8) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((objectSize >> 16) & 0xFF));
        buffer.push_back(static_cast<uint8_t>((objectSize >> 24) & 0xFF));
    }

    buffer.insert(buffer.end(), arrayEnc.buffer.begin(), arrayEnc.buffer.end());

    for (const auto& v : valueItems) {
        buffer.insert(buffer.end(), v.begin(), v.end());
    }

    return 1 + (objectSize < 0x100 ? 1 : objectSize < 0x10000 ? 2 : 4) + objectSize;
}

size_t Encoder::encodeTag(const Tag& tag, const std::vector<uint8_t>& payload) {
    if (tag.type == ValueType::Unknown) {
        return 0;
    }
    return 0;
}

bool Decoder::decodeBool() {
    uint8_t b = readByte();
    return (b & 0x01) != 0;
}

int64_t Decoder::decodeInt64() {
    uint8_t prefix = readByte();
    if ((prefix & 0xC0) == PrefixPositiveInt) {
        uint64_t v = 0;
        size_t len = prefix & 0x3F;
        if (len == 0) {
            v = prefix & 0x3F;
        } else {
            auto bytes = readBytes(len);
            for (size_t i = 0; i < len; i++) {
                v |= (static_cast<uint64_t>(bytes[i]) << (i * 8));
            }
        }
        return static_cast<int64_t>(v);
    } else if ((prefix & 0xC0) == PrefixNegativeInt) {
        uint64_t v = 0;
        size_t len = prefix & 0x3F;
        if (len == 0) {
            v = prefix & 0x3F;
        } else {
            auto bytes = readBytes(len);
            for (size_t i = 0; i < len; i++) {
                v |= (static_cast<uint64_t>(bytes[i]) << (i * 8));
            }
        }
        if (v == 0) return 0;
        return -static_cast<int64_t>(v - 1) - 1;
    }
    return 0;
}

uint64_t Decoder::decodeUint64() {
    uint8_t prefix = readByte();
    if ((prefix & 0xC0) == PrefixPositiveInt) {
        uint64_t v = 0;
        size_t len = prefix & 0x3F;
        if (len == 0) {
            v = prefix & 0x3F;
        } else {
            auto bytes = readBytes(len);
            for (size_t i = 0; i < len; i++) {
                v |= (static_cast<uint64_t>(bytes[i]) << (i * 8));
            }
        }
        return v;
    }
    return 0;
}

double Decoder::decodeFloat() {
    readByte();
    uint64_t bits = 0;
    auto bytes = readBytes(8);
    for (size_t i = 0; i < 8; i++) {
        bits |= (static_cast<uint64_t>(bytes[i]) << (i * 8));
    }
    double value;
    if (isLittleEndian()) {
        std::memcpy(&value, &bits, sizeof(double));
    } else {
        uint64_t swapped = __builtin_bswap64(bits);
        std::memcpy(&value, &swapped, sizeof(double));
    }
    return value;
}

std::string Decoder::decodeString() {
    uint8_t prefix = readByte();
    size_t len = 0;
    if ((prefix & 0xC0) == PrefixString) {
        size_t lenLen = (prefix & 0x30) >> 4;
        if (lenLen == 0) {
            len = prefix & 0x0F;
        } else if (lenLen == 1) {
            auto bytes = readBytes(1);
            len = bytes[0];
        } else if (lenLen == 2) {
            auto bytes = readBytes(2);
            len = static_cast<size_t>(bytes[0]) | (static_cast<size_t>(bytes[1]) << 8);
        }
    }
    auto bytes = readBytes(len);
    return std::string(bytes.begin(), bytes.end());
}

std::vector<uint8_t> Decoder::decodeBytes() {
    uint8_t prefix = readByte();
    size_t len = 0;
    if ((prefix & 0xC0) == PrefixBytes) {
        size_t lenLen = (prefix & 0x30) >> 4;
        if (lenLen == 0) {
            len = prefix & 0x0F;
        } else if (lenLen == 1) {
            auto bytes = readBytes(1);
            len = bytes[0];
        } else if (lenLen == 2) {
            auto bytes = readBytes(2);
            len = static_cast<size_t>(bytes[0]) | (static_cast<size_t>(bytes[1]) << 8);
        }
    }
    return readBytes(len);
}

std::vector<std::vector<uint8_t>> Decoder::decodeArray() {
    uint8_t prefix = readByte();
    size_t payloadSize = 0;
    if ((prefix & 0xC0) == PrefixContainer) {
        size_t lenLen = (prefix & 0x30) >> 4;
        if (lenLen == 0) {
            payloadSize = prefix & 0x0F;
        } else if (lenLen == 1) {
            auto bytes = readBytes(1);
            payloadSize = bytes[0];
        } else if (lenLen == 2) {
            auto bytes = readBytes(2);
            payloadSize = static_cast<size_t>(bytes[0]) | (static_cast<size_t>(bytes[1]) << 8);
        }
    }

    size_t endOffset = offset + payloadSize;
    std::vector<std::vector<uint8_t>> items;

    while (offset < endOffset && offset < data.size()) {
        size_t itemStart = offset;
        uint8_t peek = data[offset];

        if ((peek & 0xC0) == PrefixPositiveInt || (peek & 0xC0) == PrefixNegativeInt) {
            decodeInt64();
        } else if ((peek & 0xC0) == PrefixFloat) {
            decodeFloat();
        } else if ((peek & 0xC0) == PrefixString) {
            decodeString();
        } else if ((peek & 0xC0) == PrefixBytes) {
            decodeBytes();
        } else if ((peek & 0xC0) == PrefixContainer) {
            decodeArray();
        } else {
            offset++;
        }

        if (offset <= itemStart) {
            offset = itemStart + 1;
        }
        if (offset >= endOffset || offset >= data.size()) break;
    }

    return items;
}

std::map<std::string, std::vector<uint8_t>> Decoder::decodeObject() {
    std::map<std::string, std::vector<uint8_t>> result;
    uint8_t prefix = readByte();
    size_t payloadSize = 0;
    if ((prefix & 0xC0) == PrefixContainer) {
        size_t lenLen = (prefix & 0x30) >> 4;
        if (lenLen == 0) {
            payloadSize = prefix & 0x0F;
        } else if (lenLen == 1) {
            auto bytes = readBytes(1);
            payloadSize = bytes[0];
        } else if (lenLen == 2) {
            auto bytes = readBytes(2);
            payloadSize = static_cast<size_t>(bytes[0]) | (static_cast<size_t>(bytes[1]) << 8);
        }
    }
    return result;
}

}
