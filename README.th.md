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

หมายเหตุ:

- การเข้ารหัสขณะนี้รองรับได้สูงสุด 65535 ไบต์ (64KB) ขีดจำกัดนี้อาจขยายได้เมื่อรองรับประเภทเอกสารอย่างเต็มรูปแบบ
- โครงการอยู่ระหว่างการพัฒนาและทดสอบอย่างต่อเนื่อง และยังไม่แนะนำให้ใช้ในสภาพแวดล้อมการผลิต
- API และพฤติกรรมอาจยังเปลี่ยนแปลงได้ โปรดติดตามการอัปเดตเวอร์ชัน

[meta-message](https://github.com/metamessage/metamessage)

## ประเภทข้อมูล

กำหนดประเภทข้อมูลด้วยแท็ก `type=` รูปแบบคือ `type=typeIdentifier` เช่น `type=i` หมายถึงจำนวนเต็ม

- doc: การเข้ารหัสรองรับถึง 65535 ไบต์ (64KB); อาจขยายได้ในภายหลังหลังจากรองรับประเภทเอกสารเต็มรูปแบบ
- vec: อาร์เรย์/slice แบบไดนามิก, ไม่อนุญาตให้ใช้ประเภทประกอบ
- arr: อาร์เรย์ความยาวคงที่, ไม่อนุญาตให้ใช้ประเภทประกอบ
- obj: อ็อบเจกต์/struct, โครงสร้างประกอบ, สอดคล้องกับ struct/object หลายภาษา
- map: คีย์ map ต้องเป็นสตริงและค่าไม่ควรเป็นประเภทประกอบ
- str: string
- bytes: อาร์เรย์ไบต์
- bool: boolean
- i: i; ลิตเติ้ลจำนวนเต็มต้องไม่มีจุดทศนิยม
- i8: i8
- i16: i16
- i32: i32
- i64: i64
- u: u
- u8: u8
- u16: u16
- u32: u32
- u64: u64
- f32: float32; float ไม่รองรับ NaN/Inf/-0; ลิตเติ้ล float ต้องมีจุดทศนิยม เช่น 0.0
- f64: float64; float ไม่รองรับ NaN/Inf/-0; ลิตเติ้ล float ต้องมีจุดทศนิยม เช่น 0.0
- bigint: bigint
- datetime: ค่าเริ่มต้น UTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: ตัวระบุเฉพาะ
- decimal: ทศนิยม ต้องส่งเป็นสตริง
- ip: IP, รองรับ IPv4/IPv6
- url: URL, ถูกต้องตามมาตรฐาน
- email: อีเมล, ถูกต้องตามมาตรฐาน
- enum: enum, ค่าคือสตริง แยกด้วย |
- image: รูปภาพ, เก็บเป็น bytes
- video: วิดีโอ, เก็บเป็น bytes

## แท็ก

แท็กคือคำอธิบายประกอบ ป้ายชื่อ หรือคุณลักษณะของโครงสร้างภาษาการเขียนโปรแกรม หรือคำอธิบายในรูปแบบข้อความ

- is_null: ค่า null พร้อมตัวแทนว่าง
- desc: สรุป ใช้ได้กับทุกประเภท ความยาวสูงสุด 65535 บิต
- type: ประเภทข้อมูล ในรูปแบบข้อความ สตริง จำนวนเต็ม (int) ทศนิยม (float64) และอ็อบเจกต์ (หรือโครงสร้างที่คล้ายกัน) ไม่ต้องการแท็กประเภทที่ชัดเจนเมื่อไม่มีความกำกวม ในภาษาการเขียนโปรแกรม หาก object (หรือโครงสร้างที่คล้ายกัน) และ map สามารถระบุได้ map ก็ไม่จำเป็นต้องมีแท็กประเภท
- raw: ในบางภาษาการเขียนโปรแกรม ประเภทข้อมูลมักใช้ประเภท wrapper เช่น Java โดยค่าเริ่มต้น ใช้ประเภท wrapper หากไม่ต้องการ ให้ตั้งเป็น raw ยังไม่ได้พิจารณา อาจจะลบในเวอร์ชันในอนาคต
- nullable: ว่าอนุญาต null หรือไม่ ใช้ได้กับทุกประเภท
- allow_empty: ยกเว้นประเภทบูลีน ประเภทอื่นจะไม่อนุญาตให้ว่างโดยค่าเริ่มต้น เมื่อตั้ง allow_empty ค่าว่างจะได้รับอนุญาตตามกฎบางประการ
- unique: ใช้ได้กับ slices หรือ arrays เท่านั้น ระบุว่าองค์ประกอบไม่สามารถทำซ้ำได้
- default: ค่าเริ่มต้น ยังไม่เปิดใช้งาน
- example: ตัวอย่างข้อมูลที่ใช้เมื่ออาเรย์หรือ map ว่าง
- min: ความจุขั้นต่ำสำหรับอาเรย์ ความยาวขั้นต่ำสำหรับสตริง/ไบต์อาเรย์ หรือค่าต่ำสุดสำหรับตัวเลข
- max: ความจุสูงสุดสำหรับอาเรย์ ความยาวสูงสุดสำหรับสตริง/ไบต์อาเรย์ หรือค่าสูงสุดสำหรับตัวเลข
- size: ความจุสำหรับอาเรย์ ความยาวคงที่สำหรับสตริงหรือไบต์อาเรย์
- enum: เมื่อมีแท็กนี้ ค่าเป็นประเภท enum โดยค่าเริ่มต้น ประเภท enum ที่นี่อยู่ในรูปสตริง และไม่ยอมรับรูปแบบอื่น
- pattern: regex ใช้ได้กับสตริง
- location: ออฟเซ็ตโซนเวลา ค่าเริ่มต้น 0 ใช้ได้กับประเภท datetime เท่านั้น ช่วง -12 ถึง 14
- version: จำกัดเวอร์ชันใน uuid ใน ip สามารถจำกัด ipv4 หรือ ipv6 ได้
- mime: ประเภทเอกสาร ยังไม่เปิดใช้งาน

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
go get github.com/metamessage/metamessage
```

#### ตัวอย่าง

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
    data, err := mm.EncodeFromValue(p)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %x\n", data)

    var decoded Person
    err = mm.DecodeToValue(data, &decoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Decoded: %+v\n", decoded)

    jsoncStr := `{"name": "Bob", "age": 25}`
    data2, err := mm.EncodeFromJsonc(jsoncStr)
    if err != nil {
        panic(err)
    }

    jsoncOut, err := mm.DecodeToJsonc(data2)
    if err != nil {
        panic(err)
    }
    fmt.Println("JSONC:", jsoncOut)
}
```

#### ภาพรวม API

- `NewEncoder(w io.Writer) Encoder`: สร้าง encoder
- `EncodeFromValue(in any) ([]byte, error)`: เข้ารหัสจาก struct
- `EncodeFromJsonc(in string) ([]byte, error)`: เข้ารหัสจากสตริง JSONC
- `NewDecoder(r io.Reader) Decoder`: สร้าง decoder
- `DecodeToValue(in []byte, out any) error`: ถอดรหัสเป็น struct
- `DecodeToJsonc(in []byte) (string, error)`: ถอดรหัสเป็นสตริง JSONC

### ตัวอย่างในภาษาอื่น

#### Java

[jitpack.io](https://jitpack.io/#metamessage/metamessage/)

```java
import io.github.metamessage.MetaMessage;
import io.github.metamessage.MM;

@MM
class Person {
    public String name = "Ed";
    public int age = 30;
}

public class Example {
    public static void main(String[] args) throws Exception {
        Person person = new Person();
        byte[] wire = MetaMessage.encodeFromValue(person);
        Person decoded = MetaMessage.decodeToValue(wire, Person.class);
    }
}
```

#### Kotlin

[jitpack.io](https://jitpack.io/#metamessage/metamessage/)

```kotlin
import io.github.metamessage.MetaMessage
import io.github.metamessage.MM

@MM(desc="person")
class Person(var name: String = "Ed", var age: Uint8 = 30.toUint8())

fun main() {
    val person = Person()

    val wire = MetaMessage.encodeFromValue(person)

    val person = MetaMessage.decodeToValue(wire, Person::class.java)

    val jsonc = MetaMessage.valueToJsonc(person)

    val person = MetaMessage.jsoncToValue(jsoncOutput, Person::class.java)

    val wire = MetaMessage.encodeFromJsonc(jsonc)

    val jsonc = MetaMessage.decodeToJsonc(wire)
}
```

#### TypeScript

```typescript
import { encodeFromValue, decodeToValue, mm, ValueType } from 'metamessage';

@mm({ desc: '' })
class Person {
    @mm({ desc: '' })
    name: string = ''
    @mm({ desc: '' })
    age: number = 0
}
const person = { name: "Ed", age: 30 };
const wire = encodeFromValue(person);
const decoded = decodeToValue(wire, Person);
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
let wire = MetaMessage.encodeFromValue(person)
let decoded = try MetaMessage.decodeToValue(wire)
```

#### PHP

```php
<?php
use io\metamessage\mm\MetaMessage;

$person = new Person();
$wire = MetaMessage::encode($person);
$decoded = MetaMessage::decode($wire, Person::class);
```

### ตัวอย่าง

ดูโฟลเดอร์ `examples/` สำหรับตัวอย่างโค้ด
