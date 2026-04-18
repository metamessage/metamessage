# MetaMessage

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

**Ví dụ**

```jsonc
{
    // mm: type=datetime; desc=thời gian tạo
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## Chuyển đổi dữ liệu

Hỗ trợ xuất sang JSONC, YAML, TOML và các định dạng văn bản khác.

**JSONC**

- Cho phép dấu phẩy cuối cùng trong mảng hoặc đối tượng

Phong cách nhận xét được khuyến nghị:

- Cho phép nhận xét thông thường
- Nên viết nhận xét trên trường
- Thẻ mm phải ở dòng cuối cùng
- Giữ một dòng trống giữa thẻ mm và nhận xét thông thường để dễ đọc

## Lưu ý

- Vẫn còn nhiều lỗi và kiểm thử chưa hoàn chỉnh; không khuyến nghị sử dụng trong sản xuất
- Mảng và slice không cho phép kiểu hợp chất; khóa map phải là chuỗi và giá trị không được là kiểu hợp chất
- Mảng/slice rỗng tự động chèn giá trị ví dụ
- Số nguyên và chuỗi không cần thẻ kiểu rõ ràng
- Struct và slice không cần thẻ kiểu rõ ràng
- Khi kích thước mảng > 0, không cần thẻ kiểu rõ ràng
- Float không hỗ trợ NaN/Inf/-0
- Mã hóa hỗ trợ đến 65535 byte (64KB); có thể mở rộng sau
- Hằng số float phải có dấu thập phân
- Hằng số nguyên không được có dấu thập phân

## Kiểu dữ liệu

datetime: mặc định UTC 1970-01-01 00:00:00

## Thẻ

- is_null: giá trị null với trình giữ chỗ trống
- example: dữ liệu ví dụ dùng khi mảng hoặc map rỗng
- min: dung lượng tối thiểu cho mảng, độ dài tối thiểu cho chuỗi/byte array, hoặc giá trị tối thiểu cho số
- max: dung lượng tối đa cho mảng, độ dài tối đa cho chuỗi/byte array, hoặc giá trị tối đa cho số
- size: độ dài cố định cho mảng, chuỗi hoặc byte array
- location: độ lệch múi giờ, mặc định 0, phạm vi -12 đến 14

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
go get github.com/metamessage/metamessage/pkg
```

#### Ví dụ

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

#### Tóm tắt API

- `NewEncoder(w io.Writer) Encoder`: tạo encoder
- `EncodeFromStruct(in any) ([]byte, error)`: mã hóa từ struct
- `EncodeFromJSONC(in string) ([]byte, error)`: mã hóa từ chuỗi JSONC
- `NewDecoder(r io.Reader) Decoder`: tạo decoder
- `Decode(in []byte, out any) error`: giải mã vào struct
- `DecodeToJSONC(in []byte) (string, error)`: giải mã thành chuỗi JSONC

### Ví dụ

Xem thư mục `examples/` để biết mã ví dụ.
