# MetaMessage

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

**Contoh**

```jsonc
{
    // mm: type=datetime; desc=waktu pembuatan
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## Konversi data

Mendukung output ke JSONC, YAML, TOML, dan format teks lainnya.

**JSONC**

- Mengizinkan koma terakhir dalam array atau objek

Gaya komentar yang disarankan:

- Komentar biasa diizinkan
- Komentar harus ditulis di atas field
- tag mm harus pada baris terakhir
- Sisakan baris kosong antara tag mm dan komentar biasa agar lebih mudah dibaca

## Catatan

- Masih ada banyak bug dan pengujian belum lengkap; tidak disarankan untuk produksi
- Array dan slice tidak mengizinkan tipe komposit; kunci map harus string dan nilai tidak boleh tipe komposit
- Array/slice kosong secara otomatis menyisipkan nilai contoh
- Integer dan string tidak memerlukan tag tipe eksplisit
- Struct dan slice tidak memerlukan tag tipe eksplisit
- Saat ukuran array > 0, tag tipe eksplisit tidak diperlukan
- Float tidak mendukung NaN/Inf/-0
- Encoding mendukung hingga 65535 byte (64KB); ini mungkin diperluas nanti
- Literal float harus menyertakan titik desimal
- Literal integer tidak boleh menyertakan titik desimal

## Tipe data

datetime: default UTC 1970-01-01 00:00:00

## Tag

- is_null: nilai null dengan placeholder kosong
- example: data contoh yang digunakan ketika array atau map kosong
- min: kapasitas minimum untuk array, panjang minimum untuk string/byte array, atau nilai minimum untuk angka
- max: kapasitas maksimum untuk array, panjang maksimum untuk string/byte array, atau nilai maksimum untuk angka
- size: panjang tetap untuk array, string, atau byte array
- location: offset zona waktu, default 0, rentang -12 hingga 14

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
go get github.com/metamessage/metamessage/pkg
```

#### Contoh

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

#### Ringkasan API

- `NewEncoder(w io.Writer) Encoder`: buat encoder
- `EncodeFromStruct(in any) ([]byte, error)`: encoding dari struct
- `EncodeFromJSONC(in string) ([]byte, error)`: encoding dari string JSONC
- `NewDecoder(r io.Reader) Decoder`: buat decoder
- `Decode(in []byte, out any) error`: decode ke struct
- `DecodeToJSONC(in []byte) (string, error)`: decode ke string JSONC

### Contoh

Lihat direktori `examples/` untuk kode contoh.
