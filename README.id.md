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

MetaMessage (mm) adalah protokol pertukaran data terstruktur. Itu menjelaskan dirinya sendiri, membatasi dirinya sendiri, dan memberi contoh dirinya sendiri, memungkinkan pertukaran data tanpa kehilangan. Ini dirancang sebagai protokol universal generasi berikutnya yang mendukung AI, manusia, dan mesin secara native.

- Ramah manusia dan AI
- Mendukung ekspor/impor ke JSONC (saat ini; dukungan YAML/TOML direncanakan)
- Cocok untuk file konfigurasi dan pertukaran data
- Bekerja untuk API tradisional dan skenario interaksi AI
- Mendukung konversi antara struct/kelas bahasa dan MetaMessage
- Mendukung pembuatan kode untuk banyak bahasa
- Data membawa tipe, batasan, deskripsi, dan contoh tanpa dokumentasi terpisah
- Semua metadata dapat diperbarui bersama data tanpa koordinasi tambahan
- Struktur dan nilai konsisten di antara bahasa
- Tidak ada kehilangan struktur; parser beradaptasi secara otomatis dan tidak crash
- Dapat diserialisasi ke biner kompak untuk parsing lebih cepat dan ukuran lebih kecil

**Masalah yang diselesaikan**

- Tipe tidak jelas, misalnya tidak tahu apakah field adalah uint8
- Struktur tidak lengkap, misalnya null tanpa informasi tipe internal
- Tidak ada aturan validasi, sehingga tidak dapat memeriksa validitas data
- Tidak ada contoh atau deskripsi, memaksa bergantung pada dokumentasi eksternal
- Perubahan format memerlukan penyesuaian encoding/decoding dan sinkronisasi ulang dokumentasi

MetaMessage secara alami cocok untuk pemahaman dan interaksi AI, menyelesaikan ambiguitas dan ketidakakuratan. Ini menggantikan dokumentasi API tradisional, kesepakatan format verbal, dan sinkronisasi versi manual dengan membuat data menjelaskan dirinya sendiri dan berkembang secara independen.

Catatan:

- Encoding saat ini mendukung hingga 65535 byte (64KB). Batas ini dapat diperluas setelah dukungan penuh untuk tipe dokumen diterapkan.
- Proyek ini sedang dalam pengembangan dan pengujian aktif, dan belum direkomendasikan untuk penggunaan produksi.
- API dan perilaku masih bisa berubah, jadi harap perhatikan pembaruan versi.

[meta-message](https://github.com/metamessage/metamessage)

## Format teks

### JSONC

- Mengizinkan koma terakhir dalam array atau objek
- Komentar biasa diizinkan
- Komentar harus ditulis di atas field
- tag mm harus pada baris terakhir
- Sisakan baris kosong antara tag mm dan komentar biasa agar lebih mudah dibaca

**Contoh**

```jsonc
{
    // mm: type=datetime; desc=waktu pembuatan
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## Tipe data

Definisikan tipe data menggunakan tag `type=`. Formatnya adalah `type=typeIdentifier`, misalnya `type=i` berarti integer.

- doc: Encoding mendukung hingga 65535 byte (64KB). Batasan ini mungkin diperluas setelah dukungan penuh untuk tipe dokumen
- vec: array/slice dinamis, tidak mengizinkan tipe komposit
- arr: array panjang tetap, tidak mengizinkan tipe komposit
- obj: objek/struct, struktur komposit, sesuai dengan struct/object lintas-bahasa
- map: Kunci map harus string dan nilai tidak boleh tipe komposit
- str: string
- bytes: array byte
- bool: boolean
- i: i; literal integer tidak boleh menyertakan titik desimal
- i8: i8
- i16: i16
- i32: i32
- i64: i64
- u: u
- u8: u8
- u16: u16
- u32: u32
- u64: u64
- f32: float32; float tidak mendukung NaN/Inf/-0; literal float harus menyertakan titik desimal, misalnya 0.0
- f64: float64; float tidak mendukung NaN/Inf/-0; literal float harus menyertakan titik desimal, misalnya 0.0
- bigint: bigint
- datetime: default UTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: pengenal unik
- decimal: desimal, harus dikirim sebagai string
- ip: IP, mendukung IPv4/IPv6
- url: URL, valid
- email: email, valid
- enum: enum, nilai adalah string yang dipisahkan oleh |
- image: gambar, bytes di bawahnya
- video: video, bytes di bawahnya

## Tag

Tag adalah anotasi, label, atau atribut dari struktur bahasa pemrograman, atau komentar dalam format teks

- is_null: nilai null dengan placeholder kosong

- desc: ringkasan, berlaku untuk semua tipe. Panjang maksimum 65535 bit

- type: tipe data. Dalam format teks, string, integer (int), desimal (float64), slices, objek (atau struktur serupa) tidak memerlukan tag tipe eksplisit ketika tidak ambigu, misalnya ketika ukuran array > 0. Dalam bahasa pemrograman, jika array, maps, dan tipe lain dapat ditentukan, tag tipe juga tidak diperlukan

- raw: dalam beberapa bahasa pemrograman, tipe data biasanya menggunakan tipe wrapper, seperti Java. Tipe wrapper digunakan secara default; atur ke raw jika tidak diinginkan. Untuk ditentukan, mungkin akan dihapus di versi mendatang

- nullable: apakah null diizinkan, berlaku untuk semua tipe

- allow_empty: kecuali untuk tipe boolean, tipe lain tidak mengizinkan kosong secara default. Ketika allow_empty diatur, nilai kosong diizinkan mengikuti aturan tertentu

- unique: hanya berlaku untuk slices atau array, menunjukkan elemen tidak dapat diulang

- default: nilai default, belum diaktifkan

- example: data contoh yang digunakan ketika array, slices, atau map kosong, secara otomatis menghasilkan contoh nilai kosong

- min: kapasitas minimum untuk array, panjang minimum untuk string/byte array, atau nilai minimum untuk angka (integer, desimal, bigint)

- max: kapasitas maksimum untuk array, panjang maksimum untuk string/byte array, atau nilai maksimum untuk angka (integer, desimal, bigint)

- size: kapasitas untuk array, panjang tetap untuk string atau byte array

- enum: ketika tag ini ada, nilai secara default bertipe enum. Tipe enum di sini dalam bentuk string dan tidak menerima bentuk lain

- pattern: regex, berlaku untuk string

- location: offset zona waktu, default 0, hanya berlaku untuk tipe datetime, rentang -12 hingga 14

- version: batasi versi dalam uuid; di ip dapat membatasi ipv4 atau ipv6

- mime: tipe dokumen, belum diaktifkan

## Penggunaan

### Alat CLI

Proyek ini menyediakan alat baris perintah `mm` untuk encoding, decoding, dan pembuatan kode.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Bangun

```bash
make
```

#### Contoh

1. Encode JSONC ke MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

Atau baca dari stdin:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. Decode MetaMessage ke JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

Atau baca dari stdin:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Hasilkan struct dan kode dari JSONC

Mendukung go, java, ts, kt, py, js, cs, rs, swift, php

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

#### Opsi

- -encode, -e: mode encoding
- -decode, -d: mode decoding
- -generate, -g: mode pembuatan kode
- -in, -i: jalur file masukan (kosong = stdin)
- -out, -o: jalur file keluaran (kosong = stdout)
- -force, -f: timpa file keluaran
- -lang, -l: bahasa target untuk pembuatan kode (go, java, ts, kt, py, js, cs, rs, swift, php)

### Penggunaan pustaka

Proyek ini menyediakan pustaka Go untuk penggunaan programatik.

#### Instalasi

```bash
go get github.com/metamessage/metamessage
```

#### Contoh

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

#### Ringkasan API

- `NewEncoder(w io.Writer) Encoder`: buat encoder
- `EncodeFromValue(in any) ([]byte, error)`: encoding dari struct
- `EncodeFromJsonc(in string) ([]byte, error)`: encoding dari string JSONC
- `NewDecoder(r io.Reader) Decoder`: buat decoder
- `DecodeToValue(in []byte, out any) error`: decode ke struct
- `DecodeToJsonc(in []byte) (string, error)`: decode ke string JSONC

### Contoh dalam bahasa lain

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

### Contoh

Lihat direktori `examples/` untuk kode contoh.
