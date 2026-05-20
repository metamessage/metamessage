#ifndef MMCPP_IR_VALUE_TYPE_HPP
#define MMCPP_IR_VALUE_TYPE_HPP

#include <cstdint>
#include <string>
#include <unordered_map>

namespace mmc {
namespace ir {

enum class ValueType : uint8_t {
    Unknown = 0,

    Doc  = 1,
    Vec  = 2,
    Arr  = 3,
    Obj  = 4,
    Map  = 5,

    Str   = 6,
    Bytes = 7,
    Bool  = 8,

    I   = 9,
    I8  = 10,
    I16 = 11,
    I32 = 12,
    I64 = 13,
    U   = 14,
    U8  = 15,
    U16 = 16,
    U32 = 17,
    U64 = 18,

    F32 = 19,
    F64 = 20,

    Bigint   = 21,
    Datetime = 22,
    Date     = 23,
    Time     = 24,

    Uuid    = 25,
    Decimal = 26,
    Ip      = 27,
    Url     = 28,
    Email   = 29,

    Enum  = 30,
    Image = 31,
    Video = 32
};

inline std::string valueTypeToString(ValueType vt) {
    switch (vt) {
        case ValueType::Unknown:  return "unknown";
        case ValueType::Doc:      return "doc";
        case ValueType::Vec:      return "vec";
        case ValueType::Arr:      return "arr";
        case ValueType::Obj:      return "obj";
        case ValueType::Map:      return "map";
        case ValueType::Str:      return "str";
        case ValueType::Bytes:    return "bytes";
        case ValueType::Bool:     return "bool";
        case ValueType::I:        return "i";
        case ValueType::I8:       return "i8";
        case ValueType::I16:      return "i16";
        case ValueType::I32:      return "i32";
        case ValueType::I64:      return "i64";
        case ValueType::U:        return "u";
        case ValueType::U8:       return "u8";
        case ValueType::U16:      return "u16";
        case ValueType::U32:      return "u32";
        case ValueType::U64:      return "u64";
        case ValueType::F32:      return "f32";
        case ValueType::F64:      return "f64";
        case ValueType::Bigint:   return "bigint";
        case ValueType::Datetime: return "datetime";
        case ValueType::Date:     return "date";
        case ValueType::Time:     return "time";
        case ValueType::Uuid:     return "uuid";
        case ValueType::Decimal:  return "decimal";
        case ValueType::Ip:       return "ip";
        case ValueType::Url:      return "url";
        case ValueType::Email:    return "email";
        case ValueType::Enum:     return "enum";
        case ValueType::Image:    return "image";
        case ValueType::Video:    return "video";
        default:                  return "ValueType(" + std::to_string(static_cast<int>(vt)) + ")";
    }
}

inline ValueType parseValueType(const std::string& s) {
    static const std::unordered_map<std::string, ValueType> map = {
        {"unknown",  ValueType::Unknown},
        {"doc",      ValueType::Doc},
        {"vec",      ValueType::Vec},
        {"arr",      ValueType::Arr},
        {"obj",      ValueType::Obj},
        {"map",      ValueType::Map},
        {"str",      ValueType::Str},
        {"bytes",    ValueType::Bytes},
        {"bool",     ValueType::Bool},
        {"i",        ValueType::I},
        {"i8",       ValueType::I8},
        {"i16",      ValueType::I16},
        {"i32",      ValueType::I32},
        {"i64",      ValueType::I64},
        {"u",        ValueType::U},
        {"u8",       ValueType::U8},
        {"u16",      ValueType::U16},
        {"u32",      ValueType::U32},
        {"u64",      ValueType::U64},
        {"f32",      ValueType::F32},
        {"f64",      ValueType::F64},
        {"bigint",   ValueType::Bigint},
        {"datetime", ValueType::Datetime},
        {"date",     ValueType::Date},
        {"time",     ValueType::Time},
        {"uuid",     ValueType::Uuid},
        {"decimal",  ValueType::Decimal},
        {"ip",       ValueType::Ip},
        {"url",      ValueType::Url},
        {"email",    ValueType::Email},
        {"enum",     ValueType::Enum},
        {"image",    ValueType::Image},
        {"video",    ValueType::Video}
    };

    std::string lower = s;
    for (auto& c : lower) c = static_cast<char>(std::tolower(static_cast<unsigned char>(c)));

    auto it = map.find(lower);
    if (it != map.end()) return it->second;
    return ValueType::Unknown;
}

} // namespace ir
} // namespace mmc

#endif