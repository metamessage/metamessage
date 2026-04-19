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

**Example**

```jsonc
{
    // mm: type=datetime; desc=creation time
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## Data Conversion

Supports output to JSONC, YAML, TOML and other text formats.

**JSONC**

- Allows trailing commas in arrays or objects

Recommended comment style:

- Ordinary comments are allowed
- Comments should be written above fields
- The mm tag must be on the last line
- Leave an empty line between mm tag and ordinary comments for readability

## Notes

- There are still many bugs and incomplete tests; production use is not recommended
- Arrays and slices do not allow composite types; map keys must be strings and values must not be composite types
- Empty arrays/slices automatically insert an example value
- Integers and strings do not require explicit type tags
- Structs and slices do not require explicit type tags
- When array size > 0, explicit type tags are not needed
- Floats do not support NaN/Inf/-0
- Encoding supports up to 65535 bytes (64KB); this may be extended later
- Float literals must include a decimal point
- Integer literals must not include a decimal point

## Data Types

datetime: default UTC 1970-01-01 00:00:00

## Tags

- is_null: indicates a null value using an empty placeholder
- desc: summary, applies to all types. Maximum length 65535 bits
- type: data type. In text formats, strings, integers (int), decimals (float64), objects (or similar structures) do not require explicit type tags when unambiguous. In programming languages, if objects (or similar structures) and maps can be determined, maps also do not require type tags
- raw: in some programming languages, data types typically use wrapper types, such as Java. Wrapper types are used by default; set to raw if not desired. To be determined, may be removed in future versions
- nullable: whether null is allowed, applies to all types
- allow_empty: except for boolean types, other types do not allow empty by default. When allow_empty is set, empty values are allowed following certain rules
- unique: applies only to slices or arrays, indicates elements cannot be repeated
- default: default value, not yet enabled
- example: sample data used when arrays or maps are empty
- min: minimum capacity for arrays, minimum length for strings/byte arrays, or minimum value for numbers
- max: maximum capacity for arrays, maximum length for strings/byte arrays, or maximum value for numbers
- size: capacity for arrays, fixed length for strings or byte arrays
- enum: when this tag is present, the value defaults to enum type. Enum type here is in string form and does not accept other forms
- pattern: regex, applies to strings
- location: timezone offset, default 0, applies only to datetime types, range -12 to 14
- version: limit version in uuid; in ip can restrict ipv4 or ipv6
- mime: document type, not yet enabled

## Usage

### CLI Tool

This project provides a command-line tool `mm` for encoding, decoding, and code generation.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Build

```bash
make
```

#### Examples

1. Encode JSONC to MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

Or read from stdin:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. Decode MetaMessage to JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

Or read from stdin:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Generate structs and code from JSONC

Supports go, java, ts, kt, py, js, cs, rs, swift, php

```bash
./mm -generate -lang go -in input.jsonc -out output.go
```

```bash
./mm -generate -lang java -in input.jsonc -out output.java
```

```bash
./mm -generate -lang ts -in input.jsonc -out output.ts
```

```bash
./mm -generate -lang kt -in input.jsonc -out output.kt
```

```bash
./mm -generate -lang py -in input.jsonc -out output.py
```

```bash
./mm -generate -lang js -in input.jsonc -out output.js
```

```bash
./mm -generate -lang cs -in input.jsonc -out output.cs
```

```bash
./mm -generate -lang rs -in input.jsonc -out output.rs
```

```bash
./mm -generate -lang swift -in input.jsonc -out output.swift
```

```bash
./mm -generate -lang php -in input.jsonc -out output.php
```

#### Options

- -encode, -e: encode mode
- -decode, -d: decode mode
- -generate, -g: generate code mode
- -in, -i: input file path (if empty, read from stdin)
- -out, -o: output file path (if empty, write to stdout)
- -force, -f: overwrite output file
- -lang, -l: target language for code generation (go, java, ts, kt, py, js, cs, rs, swift, php)

### Library Usage

The project provides a Go library for programmatic use.

#### Install

```bash
go get github.com/metamessage/metamessage/pkg
```

#### Example

```go
package main

import (
    "fmt"
    "github.com/metamessage/metamessage/pkg"
)

func main() {
    type Person struct {
        Name string
        Age  int
    }

    p := Person{Name: "Alice", Age: 30}
    data, err := pkg.EncodeFromStruct(p)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %x\n", data)

    var decoded Person
    err = pkg.Decode(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+v\n", decoded)

    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := pkg.EncodeFromJSONC(jsoncStr)
    if err != nil {
        panic(err)
    }

    jsoncOut, err := pkg.DecodeToJSONC(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API Summary

- `NewEncoder(w io.Writer) Encoder`: create encoder
- `EncodeFromStruct(in any) ([]byte, error)`: encode from struct
- `EncodeFromJSONC(in string) ([]byte, error)`: encode from JSONC string
- `NewDecoder(r io.Reader) Decoder`: create decoder
- `Decode(in []byte, out any) error`: decode to struct
- `DecodeToJSONC(in []byte) (string, error)`: decode to JSONC string

### Examples

See the `examples/` directory for sample code.
