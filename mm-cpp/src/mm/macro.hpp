#ifndef MMCPP_MM_MACRO_HPP
#define MMCPP_MM_MACRO_HPP

#include "../ir/ast.hpp"
#include "../ir/tag.hpp"
#include "../ir/value_type.hpp"
#include <string>
#include <vector>
#include <memory>
#include <functional>
#include <type_traits>

#define MM_CONCAT2(a, b) a##b
#define MM_CONCAT(a, b) MM_CONCAT2(a, b)

#define MM_FIELD(field, type, ...) \
    { #field, #type, { __VA_ARGS__ } }

#define MM_OBJECT(structName, ...) \
    struct FieldDescriptor { \
        const char* name; \
        const char* type; \
        struct TagInit { \
            bool nullable = false; \
            bool raw = false; \
            bool allow_empty = false; \
            bool unique = false; \
            const char* desc = ""; \
            const char* default_val = ""; \
            int min = -1; \
            int max = -1; \
            int size = 0; \
            const char* enum_val = ""; \
            const char* pattern = ""; \
            int version = 0; \
            const char* mime = ""; \
        } tag; \
    }; \
    \
    static const FieldDescriptor MM_CONCAT(_mm_fields_, structName)[] = { \
        __VA_ARGS__ \
    }; \
    \
    static constexpr size_t MM_CONCAT(_mm_field_count_, structName) = \
        sizeof(MM_CONCAT(_mm_fields_, structName)) / sizeof(FieldDescriptor); \
    \
    inline ir::Tag MM_CONCAT(_mm_build_field_tag_, structName)(const FieldDescriptor& fd) { \
        ir::Tag tag = ir::Tag::create(); \
        tag.type = ir::parseValueType(fd.type); \
        if (fd.tag.nullable) tag.nullable = true; \
        if (fd.tag.raw) tag.raw = true; \
        if (fd.tag.allow_empty) tag.allowEmpty = true; \
        if (fd.tag.unique) tag.unique = true; \
        if (fd.tag.desc[0]) tag.desc = fd.tag.desc; \
        if (fd.tag.default_val[0]) tag.defaultVal = fd.tag.default_val; \
        if (fd.tag.min >= 0) tag.min = std::to_string(fd.tag.min); \
        if (fd.tag.max >= 0) tag.max = std::to_string(fd.tag.max); \
        if (fd.tag.size > 0) tag.size = fd.tag.size; \
        if (fd.tag.enum_val[0]) tag.enumVal = fd.tag.enum_val; \
        if (fd.tag.pattern[0]) tag.pattern = fd.tag.pattern; \
        if (fd.tag.version > 0) tag.version = fd.tag.version; \
        if (fd.tag.mime[0]) tag.mime = fd.tag.mime; \
        return tag; \
    } \
    \
    inline std::shared_ptr<ir::Object> MM_CONCAT(_mm_to_node_, structName)(const structName& obj) { \
        auto node = ir::makeObject(); \
        const auto& fields = MM_CONCAT(_mm_fields_, structName); \
        for (size_t i = 0; i < MM_CONCAT(_mm_field_count_, structName); ++i) { \
            auto val = ir::makeValue(); \
            val->tag = MM_CONCAT(_mm_build_field_tag_, structName)(fields[i]); \
            val->tag.type = ir::parseValueType(fields[i].type); \
            mmc::ir::detail::setFieldValue(val, obj, fields[i].name, ir::parseValueType(fields[i].type)); \
            node->fields.emplace_back(fields[i].name, val); \
        } \
        return node; \
    } \
    \
    inline structName MM_CONCAT(_mm_from_node_, structName)(std::shared_ptr<ir::Object> node) { \
        structName obj; \
        for (auto& field : node->fields) { \
            auto val = std::dynamic_pointer_cast<ir::Value>(field.value); \
            if (val) { \
                mmc::ir::detail::getFieldValue(obj, field.key, val->text, val->tag.type); \
            } \
        } \
        return obj; \
    }

namespace mmc {
namespace ir {
namespace detail {

template<typename T>
inline void setFieldValue(std::shared_ptr<Value> val, const T& obj,
                          const std::string& name, ValueType type) {
    (void)obj;
    val->text = "0";
}

template<>
inline void setFieldValue<std::string>(std::shared_ptr<Value> val, const std::string& obj,
                                       const std::string& name, ValueType type) {
    (void)obj;
    (void)name;
    val->text = obj;
    val->tag.type = type;
}

template<typename T>
inline void getFieldValue(T& obj, const std::string& name,
                          const std::string& text, ValueType type) {
    (void)obj;
    (void)name;
    (void)text;
    (void)type;
}

} // namespace detail
} // namespace ir
} // namespace mmc

#endif