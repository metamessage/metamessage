# MetaMessage

- [README 中文](README.md)
- [README English](README.en.md)
- [README 日本語](README.ja.md)
- [README 한국어](README.ko.md)
- [README Español](README.es.md)
- [README Français](README.fr.md)
- [README Deutsch](README.de.md)
- [README Русский](README.ru.md)
- [README Tiếng Việt](README.vi.md)
- [README Bahasa Indonesia](README.id.md)
- [README ไทย](README.th.md)

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

Note: Currently in development and testing, not recommended for production use

[meta-message](https://github.com/metamessage/metamessage)

## Text Formats

### JSONC

- Allows trailing commas in arrays or objects
- Allows ordinary comments
- Comments should be written above fields
- The mm tag must be on the last line
- Leave an empty line between mm tag and ordinary comments for readability

**Example**

```jsonc
{
    // mm: type=datetime; desc=creation time
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## Data Types

- doc: Encoding supports up to 65535 bytes (64KB). This limit may be extended after full support for document types
- slice: Arrays and slices do not allow composite types
- array: arr
- struct:
- map: Map keys must be strings, and values must not be composite types
- string: str
- bytes:
- bool:
- int: i; integer literals must not include a decimal point
- int8: i8
- int16: i16
- int32: i32
- int64: i64
- uint: u
- uint8: u8
- uint16: u16
- uint32: u32
- uint64: u64
- float32: f32; floats do not support NaN / Inf / -0; float literals must include a decimal point, e.g., 0.0
- float64: f64
- bigint: bi
- datetime: default UTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid
- decimal
- ip
- url
- email
- enum
- image
- video

## Tags

Tags are annotations, labels, or attributes of programming language structs, or comments in text formats

- is_null: value is null using an empty placeholder

- desc: summary, applies to all types. Maximum length 65535 bits

- type: data type. In text formats, strings, integers (int), decimals (float64), slices, objects (or similar structures) do not require explicit type tags when unambiguous, such as when array size > 0. In programming languages, if arrays, maps, and other types can be determined, type tags are also not needed

- raw: in some programming languages, data types typically use wrapper types, such as Java. Wrapper types are used by default; set to raw if not desired. To be determined, may be removed in future versions

- nullable: whether null is allowed, applies to all types

- allow_empty: except for boolean types, other types do not allow empty by default. When allow_empty is set, empty values are allowed following certain rules

- unique: applies only to slices or arrays, indicates elements cannot be repeated

- default: default value, not yet enabled

- example: sample data used when arrays, slices, or maps are empty, automatically generating an empty value example

- min: for arrays, indicates minimum capacity; for strings and byte arrays, indicates minimum length; for numeric types (integers, decimals, bigint), indicates minimum value

- max: for arrays, indicates maximum capacity; for strings and byte arrays, indicates maximum length; for numeric types (integers, decimals, bigint), indicates maximum value

- size: for arrays, indicates capacity; for strings and byte arrays, indicates fixed length

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
go get github.com/metamessage/metamessage
```

#### Example

```go
package main

import (
    "fmt"
    mm "github.com/metamessage/metamessage"
)

func main() {
    type Person struct {
        Name string
        Age  int
    }

    p := Person{Name: "Alice", Age: 30}
    data, err := mm.EncodeFromObject(p)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %x\n", data)

    var decoded Person
    err = mm.Decode(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+v\n", decoded)

    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := mm.EncodeFromJSONC(jsoncStr)
    if err != nil {
        panic(err)
    }

    jsoncOut, err := mm.DecodeToJSONC(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### API Summary

- `NewEncoder(w io.Writer) Encoder`: create encoder
- `EncodeFromObject(in any) ([]byte, error)`: encode from struct
- `EncodeFromJSONC(in string) ([]byte, error)`: encode from JSONC string
- `NewDecoder(r io.Reader) Decoder`: create decoder
- `Decode(in []byte, out any) error`: decode to struct
- `DecodeToJSONC(in []byte) (string, error)`: decode to JSONC string

### Other Language Examples

#### Java

```java
import io.metamessage.mm.MetaMessage;
import io.metamessage.mm.MM;

@MM
class Person {
    public String name = "Ed";
    public int age = 30;
}

public class Example {
    public static void main(String[] args) throws Exception {
        Person person = new Person();
        byte[] wire = MetaMessage.encode(person);
        Person decoded = MetaMessage.decode(wire, Person.class);
    }
}
```

#### Kotlin

```kotlin
import io.metamessage.mm.MetaMessage
import io.metamessage.mm.MM

@MM
class Person(var name: String = "Ed", var age: Int = 30)

fun main() {
    val person = Person()
    val wire = MetaMessage.encode(person)
    val decoded = MetaMessage.decode(wire, Person::class.java)
}
```

#### TypeScript

```typescript
import { encode, decode } from '@metamessage/ts';

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
```

#### Python

```python
from metamessage import encode, decode

person = {"name": "Ed", "age": 30}
wire = encode(person)
decoded = decode(wire)
```

#### JavaScript

```javascript
const { encode, decode } = require('metamessage');

const person = { name: "Ed", age: 30 };
const wire = encode(person);
const decoded = decode(wire);
```

#### C\#

```csharp
using MetaMessage;

var person = new Person { Name = "Ed", Age = 30 };
byte[] wire = MetaMessage.Encode(person);
var decoded = MetaMessage.Decode<Person>(wire);
```

#### Rust

```rust
use metamessage::{encode, decode, Node};

let person = Node::Object(/* ... */);
let wire = encode(&person);
let decoded = decode(&wire).unwrap();
```

#### Swift

```swift
import MetaMessage

let person = Person(name: "Ed", age: 30)
let wire = MetaMessage.encode(person)
let decoded = try MetaMessage.decode(wire)
```

#### PHP

```php
<?php
use io\metamessage\mm\MetaMessage;

$person = new Person();
$wire = MetaMessage::encode($person);
$decoded = MetaMessage::decode($wire, Person::class);
```

### Examples

See the `examples/` directory for sample code.
