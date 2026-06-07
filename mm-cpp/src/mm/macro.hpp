#ifndef MMCPP_MM_MACRO_HPP
#define MMCPP_MM_MACRO_HPP

#include "../ir/ast.hpp"
#include "../ir/tag.hpp"
#include "../ir/value_type.hpp"
#include <cstdint>
#include <memory>
#include <string>
#include <vector>

#define MM_CONCAT2(a, b) a##b
#define MM_CONCAT(a, b) MM_CONCAT2(a, b)

// ---------------------------------------------------------------------------
// toString / fromString helpers for all primitive types
// ---------------------------------------------------------------------------
namespace mmc {
namespace ir {
namespace detail {

inline std::string toString(const std::string &v) { return v; }
inline std::string toString(int32_t v) { return std::to_string(v); }
inline std::string toString(int64_t v) { return std::to_string(v); }
inline std::string toString(uint32_t v) { return std::to_string(v); }
inline std::string toString(uint64_t v) { return std::to_string(v); }
inline std::string toString(uint8_t v) {
  return std::to_string(static_cast<unsigned>(v));
}
inline std::string toString(uint16_t v) { return std::to_string(v); }
inline std::string toString(int8_t v) {
  return std::to_string(static_cast<int>(v));
}
inline std::string toString(int16_t v) { return std::to_string(v); }
inline std::string toString(bool v) { return v ? "true" : "false"; }
inline std::string toString(float v) { return std::to_string(v); }
inline std::string toString(double v) { return std::to_string(v); }

template <typename T> inline T fromString(const std::string &s);
template <> inline std::string fromString<std::string>(const std::string &s) {
  return s;
}
template <> inline int32_t fromString<int32_t>(const std::string &s) {
  return std::stoi(s);
}
template <> inline int64_t fromString<int64_t>(const std::string &s) {
  return std::stoll(s);
}
template <> inline uint32_t fromString<uint32_t>(const std::string &s) {
  return static_cast<uint32_t>(std::stoul(s));
}
template <> inline uint64_t fromString<uint64_t>(const std::string &s) {
  return std::stoull(s);
}
template <> inline uint8_t fromString<uint8_t>(const std::string &s) {
  return static_cast<uint8_t>(std::stoul(s));
}
template <> inline uint16_t fromString<uint16_t>(const std::string &s) {
  return static_cast<uint16_t>(std::stoul(s));
}
template <> inline int8_t fromString<int8_t>(const std::string &s) {
  return static_cast<int8_t>(std::stoi(s));
}
template <> inline int16_t fromString<int16_t>(const std::string &s) {
  return static_cast<int16_t>(std::stoi(s));
}
template <> inline bool fromString<bool>(const std::string &s) {
  return s == "true" || s == "1";
}
template <> inline float fromString<float>(const std::string &s) {
  return std::stof(s);
}
template <> inline double fromString<double>(const std::string &s) {
  return std::stod(s);
}

} // namespace detail
} // namespace ir
} // namespace mmc

// ---------------------------------------------------------------------------
// Type-string to C++ type mapping (used by MM_CPP_TYPE(type))
// ---------------------------------------------------------------------------
#define MM_CPP_TYPE_str std::string
#define MM_CPP_TYPE_bytes std::string
#define MM_CPP_TYPE_i int32_t
#define MM_CPP_TYPE_i8 int8_t
#define MM_CPP_TYPE_i16 int16_t
#define MM_CPP_TYPE_i32 int32_t
#define MM_CPP_TYPE_i64 int64_t
#define MM_CPP_TYPE_u uint32_t
#define MM_CPP_TYPE_u8 uint8_t
#define MM_CPP_TYPE_u16 uint16_t
#define MM_CPP_TYPE_u32 uint32_t
#define MM_CPP_TYPE_u64 uint64_t
#define MM_CPP_TYPE_f32 float
#define MM_CPP_TYPE_f64 double
#define MM_CPP_TYPE_bool bool

#define MM_CPP_TYPE(type) MM_CONCAT(MM_CPP_TYPE_, type)

// ---------------------------------------------------------------------------
// X-macro helpers — defined at file scope, used by MM_OBJECT in three passes
//   Pass 1 (descriptor):   MM_DESC_FIELD / MM_DESC_ARRAY_FIELD
//   Pass 2 (to_node):      MM_TO_FIELD   / MM_TO_ARRAY_FIELD
//   Pass 3 (from_node):    MM_FROM_FIELD / MM_FROM_ARRAY_FIELD
//
// These reference local variables from the enclosing function scope:
//   _fields, _idx, _node, obj, _f, _mm_build_field_tag
// ---------------------------------------------------------------------------

// --- Pass 1 helpers: generate FieldDescriptor entries ---

#define MM_DESC_FIELD(field, type, ...) {#field, #type, {__VA_ARGS__}},

#define MM_DESC_ARRAY_FIELD(field, childType, ...)                             \
  {#field, "arr", {.child_type = #childType, __VA_ARGS__}},

// --- Pass 2 helpers: generate _mm_to_node_ body ---

#define MM_TO_FIELD(field, type, ...)                                          \
  do {                                                                         \
    auto _v = ir::makeNodeScalar();                                            \
    _v->tag = _mm_build_field_tag(_fields[_idx]);                              \
    _v->text = mmc::ir::detail::toString(obj.field);                           \
    _node->fields.emplace_back(#field, _v);                                    \
    ++_idx;                                                                    \
  } while (0);

#define MM_TO_ARRAY_FIELD(field, childType, ...)                               \
  do {                                                                         \
    auto _a = ir::makeNodeArray();                                             \
    _a->tag = _mm_build_field_tag(_fields[_idx]);                              \
    for (auto &_item : obj.field) {                                            \
      auto _iv = ir::makeNodeScalar();                                         \
      _iv->text = mmc::ir::detail::toString(_item);                            \
      _a->items.push_back(_iv);                                                \
    }                                                                          \
    _node->fields.emplace_back(#field, _a);                                    \
    ++_idx;                                                                    \
  } while (0);

// --- Pass 3 helpers: generate _mm_from_node_ if-else chain ---

#define MM_FROM_FIELD(field, type, ...)                                        \
  if (_f.key == #field) {                                                      \
    auto _v = std::dynamic_pointer_cast<ir::NodeScalar>(_f.value);             \
    if (_v)                                                                    \
      obj.field = mmc::ir::detail::fromString<MM_CPP_TYPE(type)>(_v->text);    \
  } else

#define MM_FROM_ARRAY_FIELD(field, childType, ...)                             \
  if (_f.key == #field) {                                                      \
    auto _a = std::dynamic_pointer_cast<ir::NodeArray>(_f.value);              \
    if (_a) {                                                                  \
      for (auto &_item : _a->items) {                                          \
        auto _iv = std::dynamic_pointer_cast<ir::NodeScalar>(_item);           \
        if (_iv)                                                               \
          obj.field.push_back(                                                 \
              mmc::ir::detail::fromString<MM_CPP_TYPE(childType)>(_iv->text)); \
      }                                                                        \
    }                                                                          \
  } else

// ---------------------------------------------------------------------------
// MM_OBJECT — three-pass X-macro pattern for plain object structs (no child_*)
//
// Usage:
//   #define PersonFields(F) \
//       F(name, str, .desc="姓名", .min=1, .max=64) \
//       F(age, u8, .desc="年龄", .min=0, .max=150)
//   MM_OBJECT(Person, PersonFields)
//
//   Pass 1: generate FieldDescriptor array via MM_DESC_FIELD
//   Pass 2: generate _mm_to_node_    via MM_TO_FIELD
//   Pass 3: generate _mm_from_node_  via MM_FROM_FIELD
// ---------------------------------------------------------------------------
#define MM_OBJECT(structName, fieldsMacro)                                     \
  /* ---- FieldDescriptor struct (no child_* fields) ---- */                   \
  struct MM_CONCAT(FieldDescriptor, structName) {                              \
    const char *name;                                                          \
    const char *type;                                                          \
    struct TagInit {                                                           \
      bool nullable = false;                                                   \
      bool deprecated = false;                                                 \
      bool allow_empty = false;                                                \
      bool unique = false;                                                     \
      const char *desc = "";                                                   \
      const char *default_val = "";                                            \
      int min = -1;                                                            \
      int max = -1;                                                            \
      int size = 0;                                                            \
      const char *enum_val = "";                                               \
      const char *pattern = "";                                                \
      int version = 0;                                                         \
      const char *mime = "";                                                   \
    } tag;                                                                     \
  };                                                                           \
                                                                               \
  /* ---- Pass 1: field descriptor array (scalar fields only) ---- */          \
  static const MM_CONCAT(FieldDescriptor, structName)                          \
      MM_CONCAT(_mm_fields_, structName)[] = {fieldsMacro(MM_DESC_FIELD)};     \
                                                                               \
  static constexpr size_t MM_CONCAT(_mm_field_count_, structName) =            \
      sizeof(MM_CONCAT(_mm_fields_, structName)) /                             \
      sizeof(MM_CONCAT(FieldDescriptor, structName));                          \
                                                                               \
  /* ---- _mm_build_field_tag_ helper (no child_*) ---- */                     \
  inline ir::Tag MM_CONCAT(_mm_build_field_tag_, structName)(                  \
      const MM_CONCAT(FieldDescriptor, structName) & fd) {                     \
    ir::Tag tag = ir::Tag::create();                                           \
    tag.type = ir::parseValueType(fd.type);                                    \
    if (fd.tag.nullable)                                                       \
      tag.nullable = true;                                                     \
    if (fd.tag.deprecated)                                                     \
      tag.deprecated = true;                                                   \
    if (fd.tag.allow_empty)                                                    \
      tag.allowEmpty = true;                                                   \
    if (fd.tag.unique)                                                         \
      tag.unique = true;                                                       \
    if (fd.tag.desc[0])                                                        \
      tag.desc = fd.tag.desc;                                                  \
    if (fd.tag.default_val[0])                                                 \
      tag.default_val = fd.tag.default_val;                                    \
    if (fd.tag.min >= 0)                                                       \
      tag.min = std::to_string(fd.tag.min);                                    \
    if (fd.tag.max >= 0)                                                       \
      tag.max = std::to_string(fd.tag.max);                                    \
    if (fd.tag.size > 0)                                                       \
      tag.size = fd.tag.size;                                                  \
    if (fd.tag.enum_val[0])                                                    \
      tag.enums = fd.tag.enum_val;                                             \
    if (fd.tag.pattern[0])                                                     \
      tag.pattern = fd.tag.pattern;                                            \
    if (fd.tag.version > 0)                                                    \
      tag.version = fd.tag.version;                                            \
    if (fd.tag.mime[0])                                                        \
      tag.mime = fd.tag.mime;                                                  \
    return tag;                                                                \
  }                                                                            \
                                                                               \
  /* ---- Pass 2: _mm_to_node_ ---- */                                         \
  inline std::shared_ptr<ir::NodeObject> MM_CONCAT(_mm_to_node_, structName)(  \
      const structName &obj) {                                                 \
    auto _node = ir::makeNodeObject();                                         \
    const auto &_fields = MM_CONCAT(_mm_fields_, structName);                  \
    size_t _idx = 0;                                                           \
    auto _mm_build_field_tag =                                                 \
        [](const MM_CONCAT(FieldDescriptor, structName) & fd) -> ir::Tag {     \
      return MM_CONCAT(_mm_build_field_tag_, structName)(fd);                  \
    };                                                                         \
    fieldsMacro(MM_TO_FIELD) return _node;                                     \
  }                                                                            \
                                                                               \
  /* ---- Pass 3: _mm_from_node_ ---- */                                       \
  inline structName MM_CONCAT(_mm_from_node_, structName)(                     \
      std::shared_ptr<ir::NodeObject> node) {                                  \
    structName obj{};                                                          \
    for (auto &_f : node->fields) {                                            \
      fieldsMacro(MM_FROM_FIELD) {}                                            \
    }                                                                          \
    return obj;                                                                \
  }

// ---------------------------------------------------------------------------
// MM_MAP — three-pass X-macro pattern for map/container structs (with child_*)
//
// Usage:
//   #define MapFields(F, A) \
//       F(key, str, .desc="键") \
//       A(values, str, .child_desc="值列表", .child_min=1, .child_max=32)
//   MM_MAP(MyMap, MapFields)
//
//   Pass 1: generate FieldDescriptor array via MM_DESC_FIELD /
//   MM_DESC_ARRAY_FIELD Pass 2: generate _mm_to_node_    via MM_TO_FIELD   /
//   MM_TO_ARRAY_FIELD Pass 3: generate _mm_from_node_  via MM_FROM_FIELD /
//   MM_FROM_ARRAY_FIELD
// ---------------------------------------------------------------------------
#define MM_MAP(structName, fieldsMacro)                                        \
  /* ---- FieldDescriptor struct (with child_* fields) ---- */                 \
  struct MM_CONCAT(FieldDescriptor, structName) {                              \
    const char *name;                                                          \
    const char *type;                                                          \
    struct TagInit {                                                           \
      bool nullable = false;                                                   \
      bool deprecated = false;                                                 \
      bool allow_empty = false;                                                \
      bool unique = false;                                                     \
      const char *desc = "";                                                   \
      const char *default_val = "";                                            \
      int min = -1;                                                            \
      int max = -1;                                                            \
      int size = 0;                                                            \
      const char *enum_val = "";                                               \
      const char *pattern = "";                                                \
      int version = 0;                                                         \
      const char *mime = "";                                                   \
      const char *child_type = "";                                             \
      const char *child_desc = "";                                             \
      const char *child_default_val = "";                                      \
      int child_min = -1;                                                      \
      int child_max = -1;                                                      \
      int child_size = 0;                                                      \
      const char *child_enum_val = "";                                         \
      const char *child_pattern = "";                                          \
      bool child_nullable = false;                                             \
      bool child_allow_empty = false;                                          \
      bool child_unique = false;                                               \
      int child_version = 0;                                                   \
      const char *child_mime = "";                                             \
    } tag;                                                                     \
  };                                                                           \
                                                                               \
  /* ---- Pass 1: field descriptor array (scalar + array) ---- */              \
  static const MM_CONCAT(FieldDescriptor, structName)                          \
      MM_CONCAT(_mm_fields_, structName)[] = {                                 \
          fieldsMacro(MM_DESC_FIELD, MM_DESC_ARRAY_FIELD)};                    \
                                                                               \
  static constexpr size_t MM_CONCAT(_mm_field_count_, structName) =            \
      sizeof(MM_CONCAT(_mm_fields_, structName)) /                             \
      sizeof(MM_CONCAT(FieldDescriptor, structName));                          \
                                                                               \
  /* ---- _mm_build_field_tag_ helper (with child_*) ---- */                   \
  inline ir::Tag MM_CONCAT(_mm_build_field_tag_, structName)(                  \
      const MM_CONCAT(FieldDescriptor, structName) & fd) {                     \
    ir::Tag tag = ir::Tag::create();                                           \
    tag.type = ir::parseValueType(fd.type);                                    \
    if (fd.tag.nullable)                                                       \
      tag.nullable = true;                                                     \
    if (fd.tag.deprecated)                                                     \
      tag.deprecated = true;                                                   \
    if (fd.tag.allow_empty)                                                    \
      tag.allowEmpty = true;                                                   \
    if (fd.tag.unique)                                                         \
      tag.unique = true;                                                       \
    if (fd.tag.desc[0])                                                        \
      tag.desc = fd.tag.desc;                                                  \
    if (fd.tag.default_val[0])                                                 \
      tag.default_val = fd.tag.default_val;                                    \
    if (fd.tag.min >= 0)                                                       \
      tag.min = std::to_string(fd.tag.min);                                    \
    if (fd.tag.max >= 0)                                                       \
      tag.max = std::to_string(fd.tag.max);                                    \
    if (fd.tag.size > 0)                                                       \
      tag.size = fd.tag.size;                                                  \
    if (fd.tag.enum_val[0])                                                    \
      tag.enums = fd.tag.enum_val;                                             \
    if (fd.tag.pattern[0])                                                     \
      tag.pattern = fd.tag.pattern;                                            \
    if (fd.tag.version > 0)                                                    \
      tag.version = fd.tag.version;                                            \
    if (fd.tag.mime[0])                                                        \
      tag.mime = fd.tag.mime;                                                  \
    if (fd.tag.child_type[0])                                                  \
      tag.childType = ir::parseValueType(fd.tag.child_type);                   \
    if (fd.tag.child_desc[0])                                                  \
      tag.childDesc = fd.tag.child_desc;                                       \
    if (fd.tag.child_default_val[0])                                           \
      tag.child_default_val = fd.tag.child_default_val;                        \
    if (fd.tag.child_min >= 0)                                                 \
      tag.childMin = std::to_string(fd.tag.child_min);                         \
    if (fd.tag.child_max >= 0)                                                 \
      tag.childMax = std::to_string(fd.tag.child_max);                         \
    if (fd.tag.child_size > 0)                                                 \
      tag.childSize = fd.tag.child_size;                                       \
    if (fd.tag.child_enum_val[0])                                              \
      tag.child_enums = fd.tag.child_enum_val;                                 \
    if (fd.tag.child_pattern[0])                                               \
      tag.childPattern = fd.tag.child_pattern;                                 \
    if (fd.tag.child_nullable)                                                 \
      tag.childNullable = true;                                                \
    if (fd.tag.child_allow_empty)                                              \
      tag.childAllowEmpty = true;                                              \
    if (fd.tag.child_unique)                                                   \
      tag.childUnique = true;                                                  \
    if (fd.tag.child_version > 0)                                              \
      tag.childVersion = fd.tag.child_version;                                 \
    if (fd.tag.child_mime[0])                                                  \
      tag.childMime = fd.tag.child_mime;                                       \
    return tag;                                                                \
  }                                                                            \
                                                                               \
  /* ---- Pass 2: _mm_to_node_ ---- */                                         \
  inline std::shared_ptr<ir::NodeObject> MM_CONCAT(_mm_to_node_, structName)(  \
      const structName &obj) {                                                 \
    auto _node = ir::makeNodeObject();                                         \
    const auto &_fields = MM_CONCAT(_mm_fields_, structName);                  \
    size_t _idx = 0;                                                           \
    auto _mm_build_field_tag =                                                 \
        [](const MM_CONCAT(FieldDescriptor, structName) & fd) -> ir::Tag {     \
      return MM_CONCAT(_mm_build_field_tag_, structName)(fd);                  \
    };                                                                         \
    fieldsMacro(MM_TO_FIELD, MM_TO_ARRAY_FIELD) return _node;                  \
  }                                                                            \
                                                                               \
  /* ---- Pass 3: _mm_from_node_ ---- */                                       \
  inline structName MM_CONCAT(_mm_from_node_, structName)(                     \
      std::shared_ptr<ir::NodeObject> node) {                                  \
    structName obj{};                                                          \
    for (auto &_f : node->fields) {                                            \
      fieldsMacro(MM_FROM_FIELD, MM_FROM_ARRAY_FIELD) {}                       \
    }                                                                          \
    return obj;                                                                \
  }

#endif