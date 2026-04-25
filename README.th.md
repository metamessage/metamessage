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

หมายเหตุ: ขณะนี้อยู่ในระหว่างการพัฒนาและทดสอบ ไม่แนะนำให้ใช้งานในระบบผลิต

[meta-message](https://github.com/metamessage/metamessage)

## รูปแบบข้อความ

**JSONC**

- อนุญาตให้มีเครื่องหมายจุลภาคท้ายในอาเรย์หรือวัตถุ
- อนุญาตคอมเมนต์ปกติ
- ควรเขียนคอมเมนต์ไว้ด้านบนฟิลด์
- แท็ก mm ต้องอยู่บรรทัดสุดท้าย
- เว้นบรรทัดว่างระหว่างแท็ก mm และคอมเมนต์ปกติเพื่อให้อ่านง่าย

**YAML**

**TOML**

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

- doc:
- slice:
- array: arr
- struct:
- map:
- string: str
- bytes:
- bool:
- int: i
- int8: i8
- int16: i16
- int32: i32
- int64: i64
- uint: u
- uint8: u8
- uint16: u16
- uint32: u32
- uint64: u64
- float32: f32
- float64: f64
- bigint: bi
- datetime: ค่าเริ่มต้น UTC 1970-01-01 00:00:00
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
    data, err := mm.EncodeFromStruct(p)
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

#### ภาพรวม API

- `NewEncoder(w io.Writer) Encoder`: สร้าง encoder
- `EncodeFromStruct(in any) ([]byte, error)`: เข้ารหัสจาก struct
- `EncodeFromJSONC(in string) ([]byte, error)`: เข้ารหัสจากสตริง JSONC
- `NewDecoder(r io.Reader) Decoder`: สร้าง decoder
- `Decode(in []byte, out any) error`: ถอดรหัสเป็น struct
- `DecodeToJSONC(in []byte) (string, error)`: ถอดรหัสเป็นสตริง JSONC

### ตัวอย่างในภาษาอื่น

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
const { encode, decode } = require('@metamessage/js');

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

### ตัวอย่าง

ดูโฟลเดอร์ `examples/` สำหรับตัวอย่างโค้ด
