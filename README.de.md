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

MetaMessage (mm) ist ein strukturiertes Datenaustauschprotokoll. Es ist selbstbeschreibend, selbstbeschränkend und selbstbeispielgebend und ermöglicht verlustfreien Datenaustausch. Es ist als universelles Protokoll der nächsten Generation konzipiert, das KI, Menschen und Maschinen nativ unterstützt.

- Benutzer- und KI-freundlich
- Export/Import nach JSONC (derzeit; YAML/TOML-Unterstützung geplant)
- Geeignet für Konfigurationsdateien und Datenaustausch
- Funktioniert für traditionelle APIs und KI-Interaktionsszenarien
- Unterstützt die Konvertierung zwischen Sprachstrukturen/-klassen und MetaMessage
- Unterstützt die Codegenerierung für mehrere Sprachen
- Daten enthalten Typ, Einschränkung, Beschreibung und Beispiel ohne separate Dokumentation
- Alle Metadaten können mit den Daten aktualisiert werden, ohne zusätzliche Abstimmung
- Strukturen und Werte bleiben zwischen den Sprachen konsistent
- Keine Strukturverluste; Parser passen sich automatisch an und stürzen nicht ab
- Kann in kompaktes Binärformat serialisiert werden für schnelleres Parsen und kleinere Größe

**Gelöste Probleme**

- Unklare Typen, z. B. nicht zu wissen, ob ein Feld uint8 ist
- Unvollständige Struktur, z. B. null ohne interne Typinformation
- Keine Validierungsregeln, daher kann die Datenintegrität nicht überprüft werden
- Keine Beispiele oder Beschreibungen, was eine Abhängigkeit von externen Dokumenten erfordert
- Formatänderungen erfordern Anpassungen von Codierung/Dekodierung und erneute Dokumentensynchronisation

MetaMessage eignet sich von Natur aus für das Verständnis und die Interaktion mit KI und löst Mehrdeutigkeiten und Ungenauigkeiten. Es ersetzt traditionelle API-Dokumentation, mündliche Formatvereinbarungen und manuelle Versionssynchronisierung, indem es Daten selbsterklärend und unabhängig entwickelbar macht.

Hinweis: Derzeit in Entwicklung und Test, nicht für den Produktionseinsatz empfohlen

[meta-message](https://github.com/metamessage/metamessage)

## Textformate

### JSONC

- Erlaubt abschließende Kommas in Arrays oder Objekten
- Normale Kommentare sind erlaubt
- Kommentare sollten über den Feldern stehen
- Das mm-Tag muss in der letzten Zeile stehen
- Lasse eine Leerzeile zwischen mm-Tag und normalen Kommentaren für bessere Lesbarkeit

**Beispiel**

```jsonc
{
    // mm: type=datetime; desc=Erstellungszeit
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## Datentypen

- doc: Die Kodierung unterstützt bis zu 65535 Bytes (64KB). Diese Grenze kann nach vollständiger Unterstützung von Dokumenttypen überschritten werden
- slice: Arrays und Slices erlauben keine zusammengesetzten Typen
- array: arr
- struct:
- map: Map-Schlüssel müssen Strings sein und Werte dürfen keine zusammengesetzten Typen sein
- string: str
- bytes:
- bool:
- int: i; Ganzzahl-Literale dürfen keinen Dezimalpunkt enthalten
- int8: i8
- int16: i16
- int32: i32
- int64: i64
- uint: u
- uint8: u8
- uint16: u16
- uint32: u32
- uint64: u64
- float32: f32; Floats unterstützen NaN/Inf/-0 nicht; Float-Literale müssen einen Dezimalpunkt enthalten, z.B. 0.0
- float64: f64
- bigint: bi
- datetime: Standard UTC 1970-01-01 00:00:00
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

Tags sind Annotationen, Labels oder Attribute von Programmiersprachen-Strukturen oder Kommentare in Textformaten

- is_null: zeigt einen null-Wert mit einem leeren Platzhalter an

- desc: Zusammenfassung, gilt für alle Typen. Maximale Länge 65535 Bits

- type: Datentyp. In Textformaten erfordern Strings, Integer (int), Dezimalzahlen (float64), Slices, Objekte (oder ähnliche Strukturen) keine expliziten Typ-Tags, wenn diese eindeutig sind, z.B. wenn die Array-Größe > 0. In Programmiersprachen benötigen Arrays, Maps und andere Typen auch keine Typ-Tags, wenn diese ermittelt werden können

- raw: In einigen Programmiersprachen verwenden Datentypen üblicherweise Wrapper-Typen wie Java. Wrapper-Typen werden standardmäßig verwendet; setzen Sie auf raw, falls nicht gewünscht. Zu bestimmen, kann in zukünftigen Versionen entfernt werden

- nullable: ob null zulässig ist, gilt für alle Typen

- allow_empty: außer für boolesche Typen erlauben andere Typen standardmäßig keine leeren Werte. Wenn allow_empty gesetzt ist, werden leere Werte nach bestimmten Regeln zugelassen

- unique: gilt nur für Slices oder Arrays, zeigt an, dass Elemente nicht wiederholt werden können

- default: Standardwert, noch nicht aktiviert

- example: Beispielwert, der verwendet wird, wenn Arrays, Slices oder Maps leer sind, automatisch werden leere Wert-Beispiele generiert

- min: minimale Kapazität für Arrays, minimale Länge für Strings/Byte-Arrays oder minimaler Wert für Zahlen (Integer, Dezimalzahlen, bigint)

- max: maximale Kapazität für Arrays, maximale Länge für Strings/Byte-Arrays oder maximaler Wert für Zahlen (Integer, Dezimalzahlen, bigint)

- size: Kapazität für Arrays, feste Länge für Strings oder Byte-Arrays

- enum: wenn dieses Tag vorhanden ist, ist der Wert standardmäßig vom Typ enum. Der enum-Typ hier ist in Stringform und akzeptiert keine anderen Formen

- pattern: Regex, gilt für Strings

- location: Zeitzonenoffset, Standard 0, gilt nur für datetime-Typen, Bereich -12 bis 14

- version: Version in uuid begrenzen; in ip kann ipv4 oder ipv6 beschränkt werden

- mime: Dokumenttyp, noch nicht aktiviert

## Verwendung

### CLI-Tool

Dieses Projekt bietet ein Kommandozeilenwerkzeug `mm` zum Codieren, Dekodieren und Generieren von Code.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Build

```bash
make
```

#### Beispiele

1. JSONC nach MetaMessage codieren

```bash
./mm -encode -in input.jsonc -out output.mm
```

Oder aus stdin lesen:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. MetaMessage nach JSONC dekodieren

```bash
./mm -decode -in input.mm -out output.jsonc
```

Oder aus stdin lesen:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Strukturen und Code aus JSONC generieren

Unterstützt go, java, ts, kt, py, js, cs, rs, swift, php

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

#### Optionen

- -encode, -e: Codierungsmodus
- -decode, -d: Dekodierungsmodus
- -generate, -g: Codegenerierungsmodus
- -in, -i: Eingabedateipfad (leer = stdin)
- -out, -o: Ausgabedateipfad (leer = stdout)
- -force, -f: Ausgabedatei überschreiben
- -lang, -l: Zielsprache für Generierung (go, java, ts, kt, py, js, cs, rs, swift, php)

### Bibliotheksnutzung

Das Projekt bietet eine Go-Bibliothek für die programmgesteuerte Nutzung.

#### Installation

```bash
go get github.com/metamessage/metamessage
```

#### Beispiel

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

#### API-Übersicht

- `NewEncoder(w io.Writer) Encoder`: erstellt einen Encoder
- `EncodeFromStruct(in any) ([]byte, error)`: codiert aus einer Struct
- `EncodeFromJSONC(in string) ([]byte, error)`: codiert aus einer JSONC-Zeichenkette
- `NewDecoder(r io.Reader) Decoder`: erstellt einen Decoder
- `Decode(in []byte, out any) error`: dekodiert in eine Struct
- `DecodeToJSONC(in []byte) (string, error)`: dekodiert in eine JSONC-Zeichenkette

### Beispiele

Siehe das Verzeichnis `examples/` für Beispielcode.
