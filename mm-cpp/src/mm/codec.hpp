#ifndef MMCPP_CODEC_HPP
#define MMCPP_CODEC_HPP

#include <cstdint>
#include <vector>
#include <list>
#include <map>
#include <string>
#include <memory>
#include <variant>
#include <optional>
#include <functional>
#include <deque>
#include <set>

namespace mmc {

enum class ValueType : uint8_t {
    Unknown = 0,
    String = 1,
    Bool = 2,
    Int = 3,
    Int8 = 4,
    Int16 = 5,
    Int32 = 6,
    Int64 = 7,
    Uint = 8,
    Uint8 = 9,
    Uint16 = 10,
    Uint32 = 11,
    Uint64 = 12,
    Float32 = 13,
    Float64 = 14,
    Bytes = 15,
    BigInt = 16,
    Array = 17,
    Slice = 18,
    Object = 19,
    Map = 20,
    List = 21
};

struct Tag {
    ValueType type = ValueType::Unknown;
    int size = -1;
    int required = -1;
    int minValue = 0;
    int maxValue = 0;
    int minLength = 0;
    int maxLength = 0;
    int deprecated = 0;
    bool nullable = false;
};

using ValueData = std::variant<
    std::monostate,
    bool,
    int8_t,
    int16_t,
    int32_t,
    int64_t,
    uint8_t,
    uint16_t,
    uint32_t,
    uint64_t,
    float,
    double,
    std::string,
    std::vector<uint8_t>
>;

struct Encoder {
    std::vector<uint8_t> buffer;
    size_t offset = 0;

    void reset() {
        buffer.clear();
        buffer.reserve(1024);
        offset = 0;
    }

    size_t encode(const std::vector<uint8_t>& data) {
        size_t oldSize = buffer.size();
        buffer.insert(buffer.end(), data.begin(), data.end());
        return buffer.size() - oldSize;
    }

    size_t encodeBool(bool value);
    size_t encodeInt64(int64_t value);
    size_t encodeUint64(uint64_t value);
    size_t encodeFloat(double value);
    size_t encodeString(const std::string& value);
    size_t encodeBytes(const std::vector<uint8_t>& value);
    size_t encodeArray(const std::vector<std::vector<uint8_t>>& items);
    size_t encodeObject(const std::vector<std::pair<std::string, std::vector<uint8_t>>>& fields);
    size_t encodeTag(const Tag& tag, const std::vector<uint8_t>& payload);

    std::vector<uint8_t> result() const {
        return std::vector<uint8_t>(buffer.begin() + offset, buffer.end());
    }
};

struct Decoder {
    std::vector<uint8_t> data;
    size_t offset = 0;

    void setData(const std::vector<uint8_t>& input) {
        data = input;
        offset = 0;
    }

    uint8_t readByte() {
        if (offset >= data.size()) return 0;
        return data[offset++];
    }

    std::vector<uint8_t> readBytes(size_t n) {
        if (offset + n > data.size()) n = data.size() - offset;
        std::vector<uint8_t> result(data.begin() + offset, data.begin() + offset + n);
        offset += n;
        return result;
    }

    bool decodeBool();
    int64_t decodeInt64();
    uint64_t decodeUint64();
    double decodeFloat();
    std::string decodeString();
    std::vector<uint8_t> decodeBytes();
    std::vector<std::vector<uint8_t>> decodeArray();
    std::map<std::string, std::vector<uint8_t>> decodeObject();
};

template<typename T>
struct FieldHelper {
    using Type = T;
};

template<typename... Types>
struct VariantHelper {};

template<typename T>
struct is_stl_container : std::false_type {};

template<typename T, typename Alloc>
struct is_stl_container<std::vector<T, Alloc>> : std::true_type {};

template<typename T, typename Alloc>
struct is_stl_container<std::list<T, Alloc>> : std::true_type {};

template<typename T, typename Alloc>
struct is_stl_container<std::deque<T, Alloc>> : std::true_type {};

template<typename T, typename Alloc>
struct is_stl_container<std::set<T, Alloc>> : std::true_type {};

template<typename T, typename Alloc>
struct is_stl_container<std::multiset<T, Alloc>> : std::true_type {};

template<typename K, typename V, typename Comp, typename Alloc>
struct is_stl_container<std::map<K, V, Comp, Alloc>> : std::true_type {};

template<typename K, typename V, typename Comp, typename Alloc>
struct is_stl_container<std::multimap<K, V, Comp, Alloc>> : std::true_type {};

template<typename T>
inline constexpr bool is_stl_container_v = is_stl_container<T>::value;

template<typename T>
struct is_optional : std::false_type {};

template<typename T>
struct is_optional<std::optional<T>> : std::true_type {};

template<typename T>
inline constexpr bool is_optional_v = is_optional<T>::value;

template<typename T>
struct is_string : std::false_type {};

template<typename... Args>
struct is_string<std::basic_string<Args...>> : std::true_type {};

template<typename T>
inline constexpr bool is_string_v = is_string<T>::value;

template<typename T>
std::vector<uint8_t> encodeValue(const T& value) {
    Encoder enc;
    enc.reset();

    if constexpr (std::is_same_v<T, bool>) {
        enc.encodeBool(value);
    } else if constexpr (std::is_integral_v<T> && std::is_signed_v<T>) {
        enc.encodeInt64(static_cast<int64_t>(value));
    } else if constexpr (std::is_integral_v<T> && std::is_unsigned_v<T>) {
        enc.encodeUint64(static_cast<uint64_t>(value));
    } else if constexpr (std::is_floating_point_v<T>) {
        enc.encodeFloat(static_cast<double>(value));
    } else if constexpr (is_string_v<T>) {
        enc.encodeString(std::string(value));
    } else if constexpr (is_stl_container_v<T>) {
        std::vector<std::vector<uint8_t>> items;
        for (const auto& item : value) {
            items.push_back(encodeValue(item));
        }
        enc.encodeArray(items);
    } else {
        static_assert(sizeof(T) == 0, "Unsupported type");
    }

    return enc.result();
}

template<typename T>
T decodeValue(const std::vector<uint8_t>& data) {
    Decoder dec;
    dec.setData(data);

    if constexpr (std::is_same_v<T, bool>) {
        return dec.decodeBool();
    } else if constexpr (std::is_integral_v<T> && std::is_signed_v<T>) {
        return static_cast<T>(dec.decodeInt64());
    } else if constexpr (std::is_integral_v<T> && std::is_unsigned_v<T>) {
        return static_cast<T>(dec.decodeUint64());
    } else if constexpr (std::is_floating_point_v<T>) {
        return static_cast<T>(dec.decodeFloat());
    } else if constexpr (is_string_v<T>) {
        return T(dec.decodeString());
    } else {
        static_assert(sizeof(T) == 0, "Unsupported type");
    }
}

template<typename T>
struct Codec {
    static std::vector<uint8_t> encode(const T& value) {
        return encodeValue(value);
    }

    static T decode(const std::vector<uint8_t>& data) {
        return decodeValue<T>(data);
    }
};

template<typename T>
std::vector<uint8_t> encode(const T& value) {
    return Codec<T>::encode(value);
}

template<typename T>
T decode(const std::vector<uint8_t>& data) {
    return Codec<T>::decode(data);
}

}

#endif
