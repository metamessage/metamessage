# MetaMessage

MMC is a high-performance C library implementing the **MetaMessage** protocol — a compact binary serialization format with rich metadata (tags/annotations). It supports round-trip encoding/decoding between binary and JSONC (JSON with comments) representations.

## Architecture

```
mm-c/
├── src/
│   ├── core/           # Binary encoder/decoder (compact wire format)
│   ├── ir/             # Intermediate representation (AST nodes & tags)
│   ├── jsonc/          # JSONC parser, printer, and IR bridge
│   ├── mm.h            # Public API
│   └── mm.c            # Implementation
├── tests/              # Test suites
└── CMakeLists.txt
```

## Features

- **30+ value types**: `str`, `bool`, `i8`-`i64`, `u8`-`u64`, `f32`, `f64`, `datetime`, `uuid`, `email`, `url`, `ip`, `image`, `video`, etc.
- **Rich metadata**: Each value carries an extensible tag with `desc`, `min`, `max`, `size`, `nullable`, `enum`, `pattern`, `default`, `version`, `mime`, and child element constraints.
- **Binary encoding**: Compact wire format for efficient storage and transmission.
- **JSONC support**: Parse and print JSONC with comment-embedded metadata annotations.
- **Round-trip fidelity**: Encode → Decode preserves all data and tag attributes.

## Prerequisites

- CMake >= 3.10
- C99-compatible compiler (GCC, Clang, MSVC)

## Build & Install

```bash
cd mm-c
mkdir build && cd build
cmake ..
cmake --build .
```

## Usage

### Include the header

```c
#include "mm.h"
```

### Creating values with metadata

```c
// Simple values using convenience macros
mm_node_t* age   = mm_int(30, .min=0, .max=150, .desc="年龄");
mm_node_t* name  = mm_str("Alice", .desc="姓名");
mm_node_t* flag  = mm_bool(true, .desc="flag");
mm_node_t* pi    = mm_float(3.14, .desc="pi");
```

### Building objects

```c
mm_obj_t* person = mm_obj_new();
mm_obj_set(person, "name", mm_str("Alice", .desc="姓名"));
mm_obj_set(person, "age",  mm_int(30, .min=0, .max=150));
mm_obj_free(person);
```

### Building arrays

```c
mm_node_t* scores = mm_arr_new();
mm_arr_add(scores, mm_int(95));
mm_arr_add(scores, mm_int(87));
mm_arr_add(scores, mm_int(92));
```

### Usage guidance

`mm_int` / `mm_str` / `mm_bool` / `mm_float` macros already imply the default type (`i`, `str`, `bool`, `f64`). Only use `mm_value_create_str` with explicit type when you need a non-default type.

```c
// --- 错误写法 ---
// 宏已经隐含了类型，不需要再手动指定 type
mm_node_t* id   = mm_value_create_str("42",    MM_VALUE_I,    (mm_field_attr_t){ .desc = "用户ID" });
mm_node_t* name = mm_value_create_str("Alice", MM_VALUE_STR,  (mm_field_attr_t){ .desc = "用户名" });
mm_node_t* age  = mm_value_create_str("30",    MM_VALUE_U8,   (mm_field_attr_t){ .desc = "年龄", .min = 0, .max = 150 });

// --- 正确写法 ---
// 用 mm_int / mm_str 等宏省略默认类型，仅在需要 non-default 类型时显式指定
mm_node_t* id    = mm_int(42,     .desc = "用户ID");
mm_node_t* name  = mm_str("Alice", .desc = "用户名");
mm_node_t* age   = mm_int(30,      .desc = "年龄", .min = 0, .max = 150);
mm_node_t* email = mm_value_create_str("a@b.com", MM_VALUE_EMAIL, (mm_field_attr_t){ .desc = "电子邮箱" });

// 最外层可以指定 tag 字符串
mm_node_t* root = mm_obj_new();
// root->tag = "desc=用户"
```

### Binary encode / decode

```c
// Encode to binary buffer
mm_buffer_t* buf = mm_encode(node);

// Decode from binary buffer
mm_node_t* decoded = mm_decode(buf);

mm_buffer_free(buf);
mm_node_free(decoded);
```

### JSONC serialization / deserialization

```c
// Convert node tree to JSONC string
char* jsonc = mm_to_jsonc(node);

// Parse JSONC string back to node tree
mm_node_t* parsed = mm_from_jsonc(jsonc_str);

mm_string_free(jsonc);
mm_node_free(parsed);
```

### Available tag attributes

| Attribute           | Type      | Description                             |
| ------------------- | --------- | --------------------------------------- |
| `desc`              | `char*`   | Field description                       |
| `default_val`       | `char*`   | Default value                           |
| `min`               | `int64_t` | Minimum value                           |
| `max`               | `int64_t` | Maximum value                           |
| `size`              | `int`     | Size constraint                         |
| `enums`             | `char*`   | Enum values (pipe-separated)            |
| `pattern`           | `char*`   | Regex pattern                           |
| `location`          | `int`     | Timezone offset hours [-12, +14]        |
| `nullable`          | `bool`    | Whether null is allowed                 |
| `raw`               | `bool`    | Raw mode flag                           |
| `allow_empty`       | `bool`    | Whether empty value is allowed          |
| `unique`            | `bool`    | Uniqueness constraint                   |
| `version`           | `int`     | Version number                          |
| `mime`              | `char*`   | MIME type                               |
| `child_desc`        | `char*`   | Child element description               |
| `child_min`         | `int64_t` | Child element minimum value             |
| `child_max`         | `int64_t` | Child element maximum value             |
| `child_size`        | `int`     | Child element size constraint           |
| `child_nullable`    | `bool`    | Child element nullable                  |
| `child_raw`         | `bool`    | Child element raw mode                  |
| `child_allow_empty` | `bool`    | Child element allow empty               |
| `child_unique`      | `bool`    | Child element uniqueness                |
| `child_default`     | `char*`   | Child element default value             |
| `child_enum`        | `char*`   | Child element enum values               |
| `child_pattern`     | `char*`   | Child element regex pattern             |
| `child_location`    | `int`     | Child element timezone offset [-12,+14] |
| `child_version`     | `int`     | Child element version                   |
| `child_mime`        | `char*`   | Child element MIME type                 |
| `child_type`        | `char*`   | Child element type (e.g. "str", "i")    |

## Testing

### Run all tests

```bash
cd mm-c/build
cmake --build .    # build first
ctest              # run all tests
```

### Run a specific test suite

```bash
# JSONC parser build test
./mmc_test

# Comprehensive round-trip & API tests
./mmc_comprehensive_test
```

### Test coverage

The comprehensive test suite covers:

- Single value encode/decode round-trips (`int`, `str`, `bool`, `float`, negative values)
- Object round-trips (empty, simple, nested objects)
- Array round-trips (simple, mixed types, arrays in objects)
- Tag attribute preservation (all `mm_field_attr_t` fields)
- Byte-level binary format compatibility
- Value type parsing and string conversion
- Tag parsing, merging, and inheritance
- AST node creation and manipulation

## Supported Value Types

| Type                                        | Alias              |
| ------------------------------------------- | ------------------ |
| `str`                                       | `string`           |
| `bool`                                      | `boolean`          |
| `i`                                         | `int`              |
| `i8`-`i64`                                  | `int8`-`int64`     |
| `u`-`u64`                                   | `uint`-`uint64`    |
| `f32`                                       | `float32`          |
| `f64`                                       | `float`, `float64` |
| `bytes`, `uuid`, `datetime`, `date`, `time` |                    |
| `url`, `email`, `ip`                        |                    |
| `bigint`, `decimal`, `enum`                 |                    |
| `image`, `video`                            |                    |
| `doc`, `vec`, `arr`, `obj`, `map`           |                    |

## License

[MetaMessage](https://github.com/metamessage/metamessage)

## TODO

- default_val
- enums
- is_nulls
- name
