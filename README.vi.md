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

MetaMessage (mm) là một giao thức trao đổi dữ liệu có cấu trúc. Nó tự mô tả, tự ràng buộc và tự minh họa, cho phép trao đổi dữ liệu không mất mát. Nó được thiết kế như một giao thức phổ quát thế hệ mới hỗ trợ bản địa cho AI, con người và máy móc.

- Thân thiện với con người và AI
- Hỗ trợ xuất/nhập JSONC (hiện tại; dự định hỗ trợ YAML/TOML)
- Phù hợp với tệp cấu hình và trao đổi dữ liệu
- Hoạt động với cả API truyền thống và các kịch bản tương tác AI
- Hỗ trợ chuyển đổi giữa struct/lớp ngôn ngữ và MetaMessage
- Hỗ trợ tạo mã cho nhiều ngôn ngữ
- Dữ liệu tự chứa kiểu, ràng buộc, mô tả và ví dụ mà không cần tài liệu riêng
- Tất cả siêu dữ liệu có thể được cập nhật cùng với dữ liệu mà không cần phối hợp thêm
- Cấu trúc và giá trị nhất quán giữa các ngôn ngữ
- Không mất cấu trúc; bộ phân tích tự thích ứng và không bị sập
- Có thể tuần tự hóa sang nhị phân nhỏ gọn để phân tích nhanh hơn và kích thước nhỏ hơn

**Các vấn đề được giải quyết**

- Loại không rõ, ví dụ không biết trường có phải uint8 hay không
- Cấu trúc không đầy đủ, ví dụ null không có thông tin kiểu bên trong
- Không có quy tắc xác thực, nên không thể kiểm tra tính hợp lệ của dữ liệu
- Không có ví dụ hoặc mô tả, buộc phải phụ thuộc vào tài liệu bên ngoài
- Thay đổi định dạng yêu cầu điều chỉnh mã hóa/giải mã và đồng bộ tài liệu lại

MetaMessage tự nhiên phù hợp với hiểu biết và tương tác AI, giải quyết sự mơ hồ và thiếu chính xác. Nó thay thế tài liệu API truyền thống, thỏa thuận định dạng bằng lời và đồng bộ phiên bản thủ công bằng cách làm cho dữ liệu tự giải thích và phát triển độc lập.

Lưu ý: Hiện đang trong giai đoạn phát triển và thử nghiệm, không khuyến nghị sử dụng cho môi trường sản xuất.

[meta-message](https://github.com/metamessage/metamessage)

## Định dạng văn bản

### JSONC

- Cho phép dấu phẩy cuối cùng trong mảng hoặc đối tượng
- Cho phép nhận xét thông thường
- Nên viết nhận xét trên trường
- Thẻ mm phải ở dòng cuối cùng
- Giữ một dòng trống giữa thẻ mm và nhận xét thông thường để dễ đọc

**Ví dụ**

```jsonc
{
    // mm: type=datetime; desc=thời gian tạo
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## Kiểu dữ liệu

- doc: Mã hóa hỗ trợ đến 65535 byte (64KB). Giới hạn này có thể được mở rộng sau khi hỗ trợ đầy đủ các loại tài liệu
- slice: Mảng và slice không cho phép kiểu hợp chất
- array: arr
- struct:
- map: Khóa map phải là chuỗi và giá trị không được là kiểu hợp chất
- string: str
- bytes:
- bool:
- int: i; hằng số nguyên không được có dấu thập phân
- int8: i8
- int16: i16
- int32: i32
- int64: i64
- uint: u
- uint8: u8
- uint16: u16
- uint32: u32
- uint64: u64
- float32: f32; float không hỗ trợ NaN/Inf/-0; hằng số float phải có dấu thập phân, ví dụ 0.0
- float64: f64
- bigint: bi
- datetime: mặc định UTC 1970-01-01 00:00:00
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

## Thẻ

Thẻ là chú thích, nhãn hoặc thuộc tính của cấu trúc ngôn ngữ lập trình hoặc chú thích trong định dạng văn bản

- is_null: giá trị null với trình giữ chỗ trống

- desc: tóm tắt, áp dụng cho tất cả các loại. Độ dài tối đa 65535 bit

- type: kiểu dữ liệu. Trong định dạng văn bản, chuỗi, số nguyên (int), số thập phân (float64), slices, đối tượng (hoặc cấu trúc tương tự) không yêu cầu thẻ kiểu rõ ràng khi không có sự nhập nhằng, chẳng hạn khi kích thước mảng > 0. Trong ngôn ngữ lập trình, nếu có thể xác định được mảng, maps và các loại khác, thì thẻ kiểu cũng không cần thiết

- raw: trong một số ngôn ngữ lập trình, kiểu dữ liệu thường sử dụng kiểu bao bọc, chẳng hạn như Java. Kiểu bao bọc được sử dụng theo mặc định; đặt thành raw nếu không mong muốn. Chưa xác định, có thể sẽ bị xóa trong các phiên bản sau

- nullable: có cho phép null hay không, áp dụng cho tất cả các loại

- allow_empty: ngoại trừ kiểu boolean, các kiểu khác không cho phép trống theo mặc định. Khi đặt allow_empty, các giá trị trống được phép theo một số quy tắc

- unique: chỉ áp dụng cho slices hoặc mảng, chỉ ra rằng các phần tử không thể lặp lại

- default: giá trị mặc định, chưa được bật

- example: dữ liệu ví dụ dùng khi mảng, slices hoặc map rỗng, tự động tạo ví dụ giá trị trống

- min: dung lượng tối thiểu cho mảng, độ dài tối thiểu cho chuỗi/byte array, hoặc giá trị tối thiểu cho số (số nguyên, số thập phân, bigint)

- max: dung lượng tối đa cho mảng, độ dài tối đa cho chuỗi/byte array, hoặc giá trị tối đa cho số (số nguyên, số thập phân, bigint)

- size: dung lượng cho mảng, độ dài cố định cho chuỗi hoặc byte array

- enum: khi có thẻ này, giá trị mặc định thuộc kiểu enum. Kiểu enum ở đây ở dạng chuỗi và không chấp nhận các hình thức khác

- pattern: biểu thức chính quy, áp dụng cho chuỗi

- location: độ lệch múi giờ, mặc định 0, chỉ áp dụng cho các kiểu datetime, phạm vi -12 đến 14

- version: giới hạn phiên bản trong uuid; trong ip có thể hạn chế ipv4 hoặc ipv6

- mime: loại tài liệu, chưa được bật

## Sử dụng

### Công cụ CLI

Dự án này cung cấp công cụ dòng lệnh `mm` để mã hóa, giải mã và tạo mã.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Xây dựng

```bash
make
```

#### Ví dụ

1. Mã hóa JSONC sang MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

Hoặc đọc từ stdin:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. Giải mã MetaMessage sang JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

Hoặc đọc từ stdin:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Tạo struct và mã từ JSONC

Hỗ trợ go, java, ts, kt, py, js, cs, rs, swift, php

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

#### Tùy chọn

- -encode, -e: chế độ mã hóa
- -decode, -d: chế độ giải mã
- -generate, -g: chế độ tạo mã
- -in, -i: đường dẫn tệp đầu vào (trống = stdin)
- -out, -o: đường dẫn tệp đầu ra (trống = stdout)
- -force, -f: ghi đè tệp đầu ra
- -lang, -l: ngôn ngữ mục tiêu cho tạo mã (go, java, ts, kt, py, js, cs, rs, swift, php)

### Sử dụng thư viện

Dự án cung cấp một thư viện Go để sử dụng trong chương trình.

#### Cài đặt

```bash
go get github.com/metamessage/metamessage
```

#### Ví dụ

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

#### Tóm tắt API

- `NewEncoder(w io.Writer) Encoder`: tạo encoder
- `EncodeFromObject(in any) ([]byte, error)`: mã hóa từ struct
- `EncodeFromJSONC(in string) ([]byte, error)`: mã hóa từ chuỗi JSONC
- `NewDecoder(r io.Reader) Decoder`: tạo decoder
- `Decode(in []byte, out any) error`: giải mã vào struct
- `DecodeToJSONC(in []byte) (string, error)`: giải mã thành chuỗi JSONC

### Ví dụ về các ngôn ngữ khác

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

### Ví dụ

Xem thư mục `examples/` để biết mã ví dụ.
