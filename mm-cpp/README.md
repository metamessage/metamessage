# MetaMessage

MetaMessage (mm) is a structured data exchange protocol. It is self-describing, self-constraining, and self-exemplifying, enabling lossless data exchange. It is designed as a next-generation universal protocol that natively supports AI, humans, and machines.

- Human and AI friendly
- Export/import to JSONC (currently; YAML/TOML support planned)
- Suitable for configuration files and data exchange
- Works for traditional APIs and AI interaction scenarios
- Supports conversion between language structs/classes and MetaMessage
- Supports code generation for multiple languages
- Data carries type, constraint, description, and example without separate documentation
- All metadata can be updated with the data itself, without extra coordination
- Structures and values stay consistent across languages
- No structural loss; parsers adapt automatically and do not crash
- Can serialize to compact binary for faster decoding and smaller size

**Problems solved**

- Unknown types, such as not knowing whether a field is uint8
- Incomplete structure, such as null without inner type information
- No validation rules, so data legality cannot be checked
- No examples or descriptions, forcing reliance on separate docs
- Format changes require protocol adjustment and documentation resync

MetaMessage is naturally suited for AI understanding and interaction, solving ambiguity and imprecision in data. It replaces traditional API docs, verbal format agreements, and manual version sync by making data self-explanatory and independently evolvable.

[github.com](https://github.com/metamessage/metamessage)

## Architecture

```
mm-cpp/
├── src/
│   ├── core/           # Binary encoder/decoder (compact wire format)
│   │   ├── constants.hpp
│   │   ├── encoder.hpp
│   │   └── decoder.hpp
│   ├── ir/             # Intermediate representation (AST nodes & tags)
│   │   ├── ast.hpp
│   │   ├── tag.hpp
│   │   └── value_type.hpp
│   ├── jsonc/          # JSONC scanner, parser, and printer
│   │   ├── scanner.hpp / scanner.cpp
│   │   ├── token.hpp
│   │   ├── parser.hpp
│   │   └── printer.hpp
│   └── mm/
│       ├── mm.hpp      # High-level convenience API
│       └── macro.hpp   # MM_OBJECT / MM_FIELD declarative macros
├── tests/              # Test suites
└── CMakeLists.txt
```

## Features

- **Header-heavy design**: Most logic is in `.hpp` headers; only the scanner has a `.cpp` file.
- **30+ value types**: `str`, `bool`, `i8`-`i64`, `u8`-`u64`, `f32`, `f64`, `datetime`, `uuid`, `email`, `url`, `ip`, `media`, `bytes`, `bigint`, `decimal`, `enums`, etc.
- **Rich metadata tags**: `desc`, `min`, `max`, `size`, `nullable`, `raw`, `allowEmpty`, `unique`, `default_val`, `enums`, `pattern`, `location`, `version`, `mime` — plus child element variants.
- **JSONC with inline** **`mm:`** **annotations**: Parse comments like `// mm: type=u8; min=0; max=150` to attach metadata.
- **Declarative macro system**: `MM_OBJECT` / `MM_FIELD` for compile-time schema definition and auto-generated serializer/deserializer.
- **Binary encoding**: Compact wire format with prefix-based encoding for efficient storage and transmission.
- **Round-trip fidelity**: Encode → Decode preserves all data and tag attributes.

## Prerequisites

- CMake >= 3.10
- C++17-compatible compiler (GCC >= 8, Clang >= 7, MSVC >= 2017)

## Build & Test

```bash
cd mm-cpp
mkdir build && cd build
cmake .. -DCMAKE_CXX_COMPILER=/usr/bin/clang++ -DCMAKE_C_COMPILER=/usr/bin/clang
cmake --build .
```

### Run all tests

```bash
cd mm-cpp/build
cmake --build .       # build first
ctest                 # run all tests
```

### Run a specific test suite

```bash
# JSONC parser build test
./mmcpp_test

# Comprehensive tests (round-trip, macros, encoding, etc.)
./mmcpp_comprehensive_test
```

### Test coverage

The comprehensive test suite covers:

- **Value type** parsing and string conversion
- **Tag** creation, toString, parse from `mm:` comments, inheritance, merging, and binary encoding
- **AST** node creation (Object, Array, Value, Doc) and field manipulation
- **JSONC scanner/parser** — parse plain JSON and tagged JSONC with `mm:` annotations
- **JSONC printer** — round-trip JSON → parse → print
- **Encoder/Decoder** — binary encode/decode of objects with type preservation
- **Macro tag system** — `MM_OBJECT` field descriptors, tag building, field count
- **Full round-trip** — JSONC → binary → JSONC
- **Nested objects** — encode/decode nested structures

## Usage

### Include headers

```cpp
#include "mm/mm.hpp"    // High-level API: mm::fromJSONC, mm::toJSONC, etc.
#include "mm/macro.hpp" // Declarative MM_OBJECT / MM_FIELD macros
```

### Low-level AST API

```cpp
using namespace mmc;

// Create values
auto nameVal = ir::makeValue();
nameVal->text = "alice";
nameVal->tag.type = ir::ValueType::Str;

auto ageVal = ir::makeValue();
ageVal->text = "30";
ageVal->tag.type = ir::ValueType::U8;
ageVal->tag.min = "0";
ageVal->tag.max = "150";
ageVal->tag.desc = "年龄";

// Build an object
auto obj = ir::makeObject();
obj->fields.emplace_back("name", nameVal);
obj->fields.emplace_back("age", ageVal);

// Build an array
auto arr = ir::makeArray();
arr->items.push_back(makeValue("item1"));
arr->items.push_back(makeValue("item2"));
```

### Binary encode / decode

```cpp
#include "core/encoder.hpp"
#include "core/decoder.hpp"

using namespace mmc;

// Encode to binary
core::Encoder encoder;
auto encoded = encoder.encode(obj);

// Decode from binary
core::Decoder decoder;
auto decoded = decoder.decode(encoded);

// Cast to specific type
auto result = std::static_pointer_cast<ir::Object>(decoded);
```

### JSONC parse / print

```cpp
#include "jsonc/scanner.hpp"
#include "jsonc/parser.hpp"
#include "jsonc/printer.hpp"
#include "mm/mm.hpp"

using namespace mmc;

// Parse JSONC
auto scanner = jsonc::Scanner(R"({"name":"alice","age":30})");
auto tokens = scanner.scanAll();
jsonc::Parser parser(tokens);
auto node = parser.parse();

// Parse JSONC with mm: tag annotations
auto scanner2 = jsonc::Scanner(R"({
    // mm: desc="person object"
    "name": "alice",
    "age": 30  // mm: type=u8; min=0; max=150
})");
auto tokens2 = scanner2.scanAll();
jsonc::Parser parser2(tokens2);
auto node2 = parser2.parse();
parser2.applyTags(node2); // apply mm: annotations to node tags

// Print to JSONC string
std::string output = jsonc::toJSONC(node);
```

### High-level convenience API

```cpp
#include "mm/mm.hpp"

using namespace mmc;

// JSONC string → binary
auto bytes = mm::fromJSONC(R"({"name":"bob","age":25})");

// Binary → JSONC string
auto jsonc = mm::toJSONC(bytes);

// JSONC string → AST node
auto node = mm::parseJSONC(R"({"key":"value"})");

// JSONC string → AST node (with mm: tags applied)
auto tagged = mm::parseTaggedJSONC(R"({
    // mm: desc="root"
    "key": "value"  // mm: type=str; min=1; max=100
})");

// AST node → binary
auto encoded = mm::fromNode(node);

// Binary → AST node
auto decoded = mm::toNode(encoded);

// AST node → JSONC
auto jsoncStr = mm::toJSONCFromNode(node);

// Unified API (language-agnostic naming):

// node → binary (with optional tag string)
auto data = mm::encodeFromValue(node, "desc=root");

// JSONC string → binary
auto data = mm::encodeFromJsonc(jsoncStr);

// binary → node
auto node = mm::decodeToValue(data);

// binary → JSONC string
auto jsonc = mm::decodeToJsonc(data);

// node → JSONC string (with optional tag string)
auto jsonc = mm::valueToJsonc(node, "desc=root");

// JSONC string → node
auto node = mm::jsoncToValue(jsoncStr);
```

### Declarative MM_OBJECT macro

```cpp
#include "mm/macro.hpp"

using namespace mmc;

// Define a struct
struct Person {
    std::string name;
    uint8_t age;
};

// Declare schema with field metadata
MM_OBJECT(Person,
    MM_FIELD(name, str, .desc="姓名", .min=1, .max=64),
    MM_FIELD(age, u8, .desc="年龄", .min=0, .max=150)
);

// Auto-generated APIs:
//   _mm_field_count_Person   — field count
//   _mm_fields_Person[]      — array of FieldDescriptor
//   _mm_build_field_tag_Person(fd) → ir::Tag
//   _mm_to_node_Person(obj) → std::shared_ptr<ir::Object>
//   _mm_from_node_Person(node) → Person
```

**`MM_FIELD`** **tag initializer attributes:**

| Attribute     | Type          | Description                    |
| ------------- | ------------- | ------------------------------ |
| `desc`        | `const char*` | Field description              |
| `default_val` | `const char*` | Default value                  |
| `min`         | `int`         | Minimum value                  |
| `max`         | `int`         | Maximum value                  |
| `size`        | `int`         | Size constraint                |
| `enums`       | `const char*` | Enum values (pipe-separated)   |
| `pattern`     | `const char*` | Regex pattern                  |
| `nullable`    | `bool`        | Whether null is allowed        |
| `raw`         | `bool`        | Raw mode flag                  |
| `allow_empty` | `bool`        | Whether empty value is allowed |
| `unique`      | `bool`        | Uniqueness constraint          |
| `version`     | `int`         | Version number                 |
| `mime`        | `const char*` | MIME type                      |

### Usage guidance

`MM_FIELD(field, type, ...)` — the second parameter already specifies the type. Do not redundantly set `type=` in the tag string or via the low-level API for simple types.

```cpp
// --- 错误写法 ---
// ID 已有 int64_t 原生类型，不需要再标注 type=i64
// Email 没有原生类型，才需要用 type=email 指定
struct User {
    int64_t  id;
    std::string name;
    std::string email;
    uint8_t  age;
    bool     isActive;
};

// ❌ 冗余指定了默认类型
auto idVal = ir::makeValue();
idVal->text = "42";
idVal->tag.type = ir::ValueType::I64;  // 不需要：int64_t 对应 I

auto nameVal = ir::makeValue();
nameVal->text = "Alice";
nameVal->tag.type = ir::ValueType::Str;  // 不需要：std::string 对应 Str

auto ageVal = ir::makeValue();
ageVal->text = "30";
ageVal->tag.type = ir::ValueType::U8;  // 不需要：uint8_t 对应 U8

// --- 正确写法 ---
MM_OBJECT(User,
    MM_FIELD(id,       i64,    .desc="用户ID"),
    MM_FIELD(name,     str,    .desc="用户名", .min=1, .max=50),
    MM_FIELD(email,    email,  .desc="电子邮箱"),   // email 是特殊类型，需要显式指定
    MM_FIELD(age,      u8,     .desc="年龄", .min=0, .max=150),
    MM_FIELD(isActive, bool,   .desc="是否激活")
);

// 最外层可以指定 tag 字符串
User user{};
auto node = _mm_to_node_User(user);
node->tag.desc = "用户";
```

### Tag API

The `ir::Tag` struct supports the full tag protocol defined by [tag.go](/Users/lizongying/IdeaProjects/meta-message/internal/ir/tag.go):

| Field                      | Type        | Description                  |
| -------------------------- | ----------- | ---------------------------- |
| `desc`                     | `string`    | Field description            |
| `type`                     | `ValueType` | Value type                   |
| `raw`                      | `bool`      | Raw mode                     |
| `nullable`                 | `bool`      | Nullable                     |
| `allowEmpty`               | `bool`      | Allow empty value            |
| `unique`                   | `bool`      | Uniqueness                   |
| `defaultVal`               | `string`    | Default value                |
| `min` / `max`              | `string`    | Value bounds                 |
| `size`                     | `int`       | Size constraint              |
| `enumVal`                  | `string`    | Enum values (pipe-separated) |
| `pattern`                  | `string`    | Regex pattern                |
| `locationOffset`           | `int`       | Timezone offset \[-12, +14]  |
| `version`                  | `int`       | Version number               |
| `mime`                     | `string`    | MIME type                    |
| `isNull` / `example`       | `bool`      | Special flags                |
| `childDesc` .. `childMime` | —           | Child element constraints    |

Key methods:

- `toString()` — serialize to `mm:` annotation string
- `parse(str)` — parse from `mm:` annotation string
- `bytes()` — encode to compact binary
- `inherit(parent)` — inherit child attributes from parent array/object

## Supported Value Types

| Type                              | Alias | Category    |
| --------------------------------- | ----- | ----------- |
| `str`                             | —     | String      |
| `bytes`                           | —     | Binary      |
| `bool`                            | —     | Boolean     |
| `i`-`i64`                         | —     | Integer     |
| `u`-`u64`                         | —     | Unsigned    |
| `f32`                             | —     | Float       |
| `f64`                             | —     | Float       |
| `bigint`                          | —     | Big integer |
| `decimal`                         | —     | Decimal     |
| `datetime`, `date`, `time`        | —     | Temporal    |
| `uuid`, `ip`, `url`, `email`      | —     | Identifier  |
| `enums`                           | —     | Enumeration |
| `media`                           | —     | Media       |
| `doc`, `vec`, `arr`, `obj`, `map` | —     | Container   |
