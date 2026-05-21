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

Hinweis:

- Die Codierung unterstützt derzeit bis zu 65535 Bytes (64KB). Diese Grenze kann erweitert werden, sobald der vollständige Support für Dokumenttypen umgesetzt ist.
- Das Projekt befindet sich in aktiver Entwicklung und Testphase und wird für den Produktionseinsatz noch nicht empfohlen.
- APIs und Verhalten können sich noch ändern, daher bitte die Versionshinweise beachten.

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

Definieren Sie Datentypen mit dem Tag `type=`. Das Format ist `type=typeIdentifier`, z. B. `type=i` bedeutet Integer.

- doc: Die Kodierung unterstützt bis zu 65535 Bytes (64KB). Diese Grenze kann nach vollständiger Unterstützung von Dokumenttypen überschritten werden
- vec: dynamische Arrays/Slices, erlauben keine zusammengesetzten Typen
- arr: Array, feste Länge, erlaubt keine zusammengesetzten Typen
- obj: Objekt/Struct, zusammengesetzte Struktur, entspricht mehrsprachigem struct/object
- map: Map-Schlüssel müssen Strings sein und Werte dürfen keine zusammengesetzten Typen sein
- str: string
- bytes: Byte-Array
- bool: boolean
- i: int; Ganzzahl-Literale dürfen keinen Dezimalpunkt enthalten
- i8: int8
- i16: int16
- i32: int32
- i64: int64
- u: uint
- u8: uint8
- u16: uint16
- u32: uint32
- u64: uint64
- f32: float32; Floats unterstützen NaN/Inf/-0 nicht; Float-Literale müssen einen Dezimalpunkt enthalten, z.B. 0.0
- f64: float64; Floats unterstützen NaN/Inf/-0 nicht; Float-Literale müssen einen Dezimalpunkt enthalten, z.B. 0.0
- bigint: bigint
- datetime: Standard UTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: eindeutige Kennung
- decimal: Dezimalzahl, muss als Zeichenfolge übergeben werden
- ip: IP, unterstützt IPv4/IPv6
- url: URL, muss eine gültige URL sein
- email: E-Mail, muss gültig sein
- enum: enum, Werte sind Strings, getrennt durch |
- image: Bild, intern bytes
- video: Video, intern bytes

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

Gleiche Logik gilt für andere Sprachen

```go
// Falscher Ansatz
// ID hat bereits den Typ int64, keine Typangabe type=i64 nötig
// In mm keine json-Tags zusammen verwenden. mm handhabt die Benennung automatisch; falls nötig name= in mm verwenden
// Age sollte nativen Typ uint8 verwenden, daher kann type=u8 weggelassen werden
type User struct {
	ID       int64  `mm:"type=i64;desc=User ID" json:"id"`
	Name     string `mm:"type=str;desc=User Name;min=1;max=50" json:"name"`
	Email    string `mm:"type=email;desc=Email" json:"email"`
	Age      int    `mm:"type=u8;desc=Age;min=0;max=150" json:"age"`
	IsActive bool   `mm:"type=bool;desc=Is Active" json:"is_active"`
}

// Richtiger Ansatz
// Email hat keinen nativen Typ, daher ist type=email erforderlich
type User struct {
	ID       int64  `mm:"desc=User ID"`
	Name     string `mm:"desc=User Name;min=1;max=50"`
	Email    string `mm:"type=email;desc=Email"`
	Age      uint8  `mm:"desc=Age;min=0;max=150"`
	IsActive bool   `mm:"desc=Is Active"`
}

user := User{}

// Tag auf oberster Ebene kann hier angegeben werden
tag := "desc=User"
_, _ = EncodeFromValue(user, tag)
```

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

#### API-Übersicht

- `NewEncoder(w io.Writer) Encoder`: erstellt einen Encoder
- `EncodeFromValue(in any) ([]byte, error)`: codiert aus einer Struct
- `EncodeFromJsonc(in string) ([]byte, error)`: codiert aus einer JSONC-Zeichenkette
- `NewDecoder(r io.Reader) Decoder`: erstellt einen Decoder
- `DecodeToValue(in []byte, out any) error`: dekodiert in eine Struct
- `DecodeToJsonc(in []byte) (string, error)`: dekodiert in eine JSONC-Zeichenkette

### Beispiele in anderen Sprachen

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

### Beispiele

Siehe das Verzeichnis `examples/` für Beispielcode.
