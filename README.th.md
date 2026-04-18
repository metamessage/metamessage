# MetaMessage

MetaMessage (mm) เป็นโปรโตคอลแลกเปลี่ยนข้อมูลแบบมีโครงสร้าง มันอธิบายตัวเอง กำหนดข้อจำกัดตัวเอง และยกตัวอย่างตัวเอง ทำให้การแลกเปลี่ยนข้อมูลไม่สูญเสียข้อมูล ถูกออกแบบเป็นโปรโตคอลสากลยุคถัดไปที่รองรับ AI มนุษย์ และเครื่องจักรโดยตรง

- เป็นมิตรกับมนุษย์และ AI
- รองรับการส่งออก/นำเข้าไปยัง JSONC (ปัจจุบัน; วางแผนรองรับ YAML/TOML)
- เหมาะสำหรับไฟล์การตั้งค่าและการแลกเปลี่ยนข้อมูล
- ทำงานได้ทั้งกับ API แบบดั้งเดิมและสถานการณ์โต้ตอบ AI
- รองรับการแปลงระหว่าง struct/class ของภาษาและ MetaMessage
- รองรับการสร้างโค้ดสำหรับหลายภาษา
- ข้อมูลมีประเภท ข้อจำกัด คำอธิบาย และตัวอย่างในตัวโดยไม่ต้องมีเอกสารแยกต่างหาก
- เมตาดาต้าทั้งหมดสามารถอัปเดตพร้อมกับข้อมูลโดยไม่ต้องประสานเพิ่มเติม
- โครงสร้างและค่าคงที่สอดคล้องกันระหว่างภาษา
- ไม่มีการสูญเสียโครงสร้าง ตัวถอดรหัสปรับตัวเองได้โดยอัตโนมัติและไม่ล่ม
- สามารถแปลงเป็นไบนารีขนาดกะทัดรัดเพื่อการวิเคราะห์ที่เร็วขึ้นและขนาดที่เล็กลง

**ปัญหาที่แก้ไข**

- ประเภทไม่ชัดเจน เช่น ไม่รู้ว่าฟิลด์เป็น uint8 หรือไม่
- โครงสร้างไม่สมบูรณ์ เช่น null ไม่มีข้อมูลประเภทภายใน
- ไม่มีข้อกำหนดการตรวจสอบความถูกต้อง จึงไม่สามารถตรวจสอบความถูกต้องของข้อมูลได้
- ไม่มีตัวอย่างหรือคำอธิบาย จึงต้องพึ่งพาเอกสารภายนอก
- การเปลี่ยนแปลงรูปแบบต้องปรับโค้ดเข้ารหัส/ถอดรหัสและซิงค์เอกสารใหม่

MetaMessage เหมาะสำหรับการเข้าใจและโต้ตอบกับ AI โดยธรรมชาติ แก้ปัญหาความกำกวมและความไม่แม่นยำ มันแทนที่เอกสาร API แบบเดิม การตกลงรูปแบบด้วยวาจา และการซิงค์เวอร์ชันด้วยตนเอง โดยทำให้ข้อมูลอธิบายตัวเองได้และพัฒนาได้อย่างอิสระ

**ตัวอย่าง**

```jsonc
{
    // mm: type=datetime; desc=เวลาเริ่มต้นสร้าง
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## การแปลงข้อมูล

รองรับการส่งออกเป็น JSONC, YAML, TOML และรูปแบบข้อความอื่นๆ

**JSONC**

- อนุญาตให้มีเครื่องหมายจุลภาคท้ายในอาเรย์หรือวัตถุ

สไตล์คอมเมนต์ที่แนะนำ:

- อนุญาตคอมเมนต์ปกติ
- ควรเขียนคอมเมนต์ไว้ด้านบนฟิลด์
- แท็ก mm ต้องอยู่บรรทัดสุดท้าย
- เว้นบรรทัดว่างระหว่างแท็ก mm และคอมเมนต์ปกติเพื่อให้อ่านง่าย

## หมายเหตุ

- ยังมีบั๊กอยู่มากและการทดสอบยังไม่สมบูรณ์; ไม่แนะนำให้ใช้ในสภาพแวดล้อมการผลิต
- อาเรย์และสไลซ์ไม่อนุญาตให้ใช้ประเภทประกอบ; คีย์ map ต้องเป็นสตริงและค่าต้องไม่ใช่ประเภทประกอบ
- อาเรย์/สไลซ์ว่างจะใส่ค่าตัวอย่างโดยอัตโนมัติ
- จำนวนเต็มและสตริงไม่ต้องมีแท็กประเภทชัดเจน
- struct และสไลซ์ไม่ต้องมีแท็กประเภทชัดเจน
- เมื่อขนาดอาเรย์ > 0 ไม่ต้องมีแท็กประเภทชัดเจน
- float ไม่รองรับ NaN/Inf/-0
- การเข้ารหัสรองรับถึง 65535 ไบต์ (64KB); อาจขยายได้ในภายหลัง
- ลิตเติ้ล float ต้องมีจุดทศนิยม
- ลิตเติ้ล integer ต้องไม่มีจุดทศนิยม

## ประเภทข้อมูล

datetime: ดีฟอลต์ UTC 1970-01-01 00:00:00

## แท็ก

- is_null: ค่า null พร้อมตัวแทนว่าง
- example: ตัวอย่างข้อมูลที่ใช้เมื่ออาเรย์หรือ map ว่าง
- min: ความจุขั้นต่ำสำหรับอาเรย์ ความยาวขั้นต่ำสำหรับสตริง/ไบต์อาเรย์ หรือค่าต่ำสุดสำหรับตัวเลข
- max: ความจุสูงสุดสำหรับอาเรย์ ความยาวสูงสุดสำหรับสตริง/ไบต์อาเรย์ หรือค่าสูงสุดสำหรับตัวเลข
- size: ความยาวคงที่สำหรับอาเรย์ สตริง หรือไบต์อาเรย์
- location: ออฟเซ็ตโซนเวลา ค่าเริ่มต้น 0 ช่วง -12 ถึง 14

## การใช้งาน

### เครื่องมือ CLI

โครงการนี้มีเครื่องมือบรรทัดคำสั่ง `mm` สำหรับการเข้ารหัส การถอดรหัส และการสร้างโค้ด

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### สร้าง

```bash
make
```

#### ตัวอย่าง

1. เข้ารหัส JSONC เป็น MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

หรืออ่านจาก stdin:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. ถอดรหัส MetaMessage เป็น JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

หรืออ่านจาก stdin:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. สร้าง struct และโค้ดจาก JSONC

รองรับ go, java, ts, kt, py, js, cs, rs, swift, php

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

#### ตัวเลือก

- -encode, -e: โหมดเข้ารหัส
- -decode, -d: โหมดถอดรหัส
- -generate, -g: โหมดสร้างโค้ด
- -in, -i: เส้นทางไฟล์นำเข้า (ว่าง = stdin)
- -out, -o: เส้นทางไฟล์ส่งออก (ว่าง = stdout)
- -force, -f: เขียนทับไฟล์ส่งออก
- -lang, -l: ภาษาที่ต้องการสำหรับสร้างโค้ด (go, java, ts, kt, py, js, cs, rs, swift, php)

### การใช้งานไลบรารี

โครงการนี้มีไลบรารี Go สำหรับใช้ในโปรแกรม

#### ติดตั้ง

```bash
go get github.com/metamessage/metamessage/pkg
```

#### ตัวอย่าง

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

#### ภาพรวม API

- `NewEncoder(w io.Writer) Encoder`: สร้าง encoder
- `EncodeFromStruct(in any) ([]byte, error)`: เข้ารหัสจาก struct
- `EncodeFromJSONC(in string) ([]byte, error)`: เข้ารหัสจากสตริง JSONC
- `NewDecoder(r io.Reader) Decoder`: สร้าง decoder
- `Decode(in []byte, out any) error`: ถอดรหัสเป็น struct
- `DecodeToJSONC(in []byte) (string, error)`: ถอดรหัสเป็นสตริง JSONC

### ตัวอย่าง

ดูโฟลเดอร์ `examples/` สำหรับตัวอย่างโค้ด
