#ifndef MMCPP_CORE_CONSTANTS_HPP
#define MMCPP_CORE_CONSTANTS_HPP

#include <cstdint>
#include <string>

namespace mmc {
namespace core {

enum class Prefix : uint8_t {
    Simple      = 0b000 << 5,
    PositiveInt = 0b001 << 5,
    NegativeInt = 0b010 << 5,
    Float       = 0b011 << 5,
    String      = 0b100 << 5,
    Bytes       = 0b101 << 5,
    Container   = 0b110 << 5,
    Tag         = 0b111 << 5
};

enum SimpleValue : uint8_t {
    SimpleNullBool   = 0,
    SimpleNullInt    = 1,
    SimpleNullFloat  = 2,
    SimpleNullString = 3,
    SimpleNullBytes  = 4,

    SimpleFalse = 5,
    SimpleTrue  = 6,

    SimpleCode        = 7,
    SimpleMessage     = 8,
    SimpleData        = 9,
    SimpleSuccess     = 10,
    SimpleError       = 11,
    SimpleUnknown     = 12,
    SimplePage        = 13,
    SimpleLimit       = 14,
    SimpleOffset      = 15,
    SimpleTotal       = 16,
    SimpleId          = 17,
    SimpleName        = 18,
    SimpleDescription = 19,
    SimpleType        = 20,
    SimpleVersion     = 21,
    SimpleStatus      = 22,
    SimpleUrl         = 23,
    SimpleCreateTime  = 24,
    SimpleUpdateTime  = 25,
    SimpleDeleteTime  = 26,
    SimpleAccount     = 27,
    SimpleToken       = 28,
    SimpleExpireTime  = 29,
    SimpleKey         = 30,
    SimpleVal         = 31
};

constexpr uint8_t PrefixMask = 0b11100000;
constexpr uint8_t SuffixMask = 0b00011111;

inline Prefix getPrefix(uint8_t b) {
    return static_cast<Prefix>(b & PrefixMask);
}

inline uint8_t getSuffix(uint8_t b) {
    return b & SuffixMask;
}

// Int length
constexpr uint8_t IntLenMask  = 0b11111;
constexpr uint8_t IntLen1Byte = IntLenMask - 7;
constexpr uint8_t IntLen2Byte = IntLenMask - 6;
constexpr uint8_t IntLen3Byte = IntLenMask - 5;
constexpr uint8_t IntLen4Byte = IntLenMask - 4;
constexpr uint8_t IntLen5Byte = IntLenMask - 3;
constexpr uint8_t IntLen6Byte = IntLenMask - 2;
constexpr uint8_t IntLen7Byte = IntLenMask - 1;
constexpr uint8_t IntLen8Byte = IntLenMask;

inline int intLen(uint8_t b) {
    int l = static_cast<int>(b & IntLenMask);
    if (l < IntLen1Byte) return 0; // embedded value = l
    return l - IntLen1Byte + 1;
}

// Float length
constexpr uint8_t FloatPositiveNegativeMask = 0b10000;
constexpr uint8_t FloatLenMask  = 0b01111;
constexpr uint8_t FloatLen1Byte = FloatLenMask - 7;
constexpr uint8_t FloatLen2Byte = FloatLenMask - 6;
constexpr uint8_t FloatLen3Byte = FloatLenMask - 5;
constexpr uint8_t FloatLen4Byte = FloatLenMask - 4;
constexpr uint8_t FloatLen5Byte = FloatLenMask - 3;
constexpr uint8_t FloatLen6Byte = FloatLenMask - 2;
constexpr uint8_t FloatLen7Byte = FloatLenMask - 1;
constexpr uint8_t FloatLen8Byte = FloatLenMask;

inline int floatLen(uint8_t b) {
    int l = static_cast<int>(b & FloatLenMask);
    if (l < FloatLen1Byte) return 0;
    return l - FloatLen1Byte + 1;
}

// String length
constexpr uint8_t StringLenMask  = 0b11111;
constexpr uint8_t StringLen1Byte = StringLenMask - 1;
constexpr uint8_t StringLen2Byte = StringLenMask;

inline int stringExtraLen(uint8_t b) {
    int l = static_cast<int>(b & StringLenMask);
    if (l < StringLen1Byte) return 0;
    if (l == StringLen1Byte) return 1;
    return 2;
}

inline int stringInlineLen(uint8_t b) {
    int l = static_cast<int>(b & StringLenMask);
    if (l < StringLen1Byte) return l;
    return 0;
}

// Bytes length
constexpr uint8_t BytesLenMask  = 0b11111;
constexpr uint8_t BytesLen1Byte = BytesLenMask - 1;
constexpr uint8_t BytesLen2Byte = BytesLenMask;

inline int bytesExtraLen(uint8_t b) {
    int l = static_cast<int>(b & BytesLenMask);
    if (l < BytesLen1Byte) return 0;
    if (l == BytesLen1Byte) return 1;
    return 2;
}

inline int bytesInlineLen(uint8_t b) {
    int l = static_cast<int>(b & BytesLenMask);
    if (l < BytesLen1Byte) return l;
    return 0;
}

// Container
constexpr uint8_t ContainerMask   = 0b10000;
constexpr uint8_t ContainerObject = 0b00000;
constexpr uint8_t ContainerArray  = 0b10000;
constexpr uint8_t ContainerLenMask  = 0b01111;
constexpr uint8_t ContainerLen1Byte = ContainerLenMask - 1;
constexpr uint8_t ContainerLen2Byte = ContainerLenMask;

inline int containerExtraLen(uint8_t b) {
    int l = static_cast<int>(b & ContainerLenMask);
    if (l < ContainerLen1Byte) return 0;
    if (l == ContainerLen1Byte) return 1;
    return 2;
}

inline int containerInlineLen(uint8_t b) {
    int l = static_cast<int>(b & ContainerLenMask);
    if (l < ContainerLen1Byte) return l;
    return 0;
}

inline bool isArrayContainer(uint8_t b) {
    return (b & ContainerMask) == ContainerArray;
}

// Tag
constexpr uint8_t TagLenMask  = 0b11111;
constexpr uint8_t TagLen1Byte = TagLenMask - 1;
constexpr uint8_t TagLen2Byte = TagLenMask;

inline int tagExtraLen(uint8_t b) {
    int l = static_cast<int>(b & TagLenMask);
    if (l < TagLen1Byte) return 0;
    if (l == TagLen1Byte) return 1;
    return 2;
}

inline int tagInlineLen(uint8_t b) {
    int l = static_cast<int>(b & TagLenMask);
    if (l < TagLen1Byte) return l;
    return 0;
}

// TagPayload
constexpr uint8_t TagPayload1Byte = BytesLenMask - 1;
constexpr uint8_t TagPayload2Byte = BytesLenMask;

} // namespace core
} // namespace mmc

#endif