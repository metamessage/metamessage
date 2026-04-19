# MetaMessage

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

**Beispiel**

```jsonc
{
    // mm: type=datetime; desc=Erstellungszeit
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## Datenkonvertierung

Unterstützt Ausgabe nach JSONC, YAML, TOML und andere Textformate.

**JSONC**

- Erlaubt abschließende Kommas in Arrays oder Objekten

Empfohlener Kommentarstil:

- Normale Kommentare sind erlaubt
- Kommentare sollten über den Feldern stehen
- Das mm-Tag muss in der letzten Zeile stehen
- Lasse eine Leerzeile zwischen mm-Tag und normalen Kommentaren für bessere Lesbarkeit

## Hinweise

- Es gibt noch viele Fehler und unvollständige Tests; Produktionseinsatz wird nicht empfohlen
- Arrays und Slices erlauben keine zusammengesetzten Typen; Map-Schlüssel müssen Strings sein und Werte dürfen keine zusammengesetzten Typen sein
- Leere Arrays/Slices fügen automatisch einen Beispielwert ein
- Ganzzahlen und Strings benötigen keine expliziten Typ-Tags
- Structs und Slices benötigen keine expliziten Typ-Tags
- Wenn die Array-Größe > 0 ist, sind explizite Typ-Tags nicht erforderlich
- Floats unterstützen kein NaN/Inf/-0
- Kodierung unterstützt bis zu 65535 Bytes (64KB); dies kann später erweitert werden
- Float-Literale müssen einen Dezimalpunkt enthalten
- Ganzzahl-Literale dürfen keinen Dezimalpunkt enthalten

## Datentypen

datetime: Standard UTC 1970-01-01 00:00:00

## Tags

- is_null: zeigt einen null-Wert mit einem leeren Platzhalter an
- desc: Zusammenfassung, gilt für alle Typen. Maximale Länge 65535 Bits
- type: Datentyp. In Textformaten erfordern Strings, Integer (int), Dezimalzahlen (float64) und Objekte (oder ähnliche Strukturen) keine expliziten Typ-Tags, wenn diese eindeutig sind. In Programmiersprachen benötigen Maps auch keine Typ-Tags, wenn Objekte (oder ähnliche Strukturen) und Maps ermittelt werden können
- raw: In einigen Programmiersprachen verwenden Datentypen üblicherweise Wrapper-Typen wie Java. Wrapper-Typen werden standardmäßig verwendet; setzen Sie auf raw, falls nicht gewünscht. Zu bestimmen, kann in zukünftigen Versionen entfernt werden
- nullable: ob null zulässig ist, gilt für alle Typen
- allow_empty: außer für boolesche Typen erlauben andere Typen standardmäßig keine leeren Werte. Wenn allow_empty gesetzt ist, werden leere Werte nach bestimmten Regeln zugelassen
- unique: gilt nur für Slices oder Arrays, zeigt an, dass Elemente nicht wiederholt werden können
- default: Standardwert, noch nicht aktiviert
- example: Beispielwert, der verwendet wird, wenn Arrays oder Maps leer sind
- min: minimale Kapazität für Arrays, minimale Länge für Strings/Byte-Arrays oder minimaler Wert für Zahlen
- max: maximale Kapazität für Arrays, maximale Länge für Strings/Byte-Arrays oder maximaler Wert für Zahlen
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
go get github.com/metamessage/metamessage/pkg
```

#### Beispiel

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

#### API-Übersicht

- `NewEncoder(w io.Writer) Encoder`: erstellt einen Encoder
- `EncodeFromStruct(in any) ([]byte, error)`: codiert aus einer Struct
- `EncodeFromJSONC(in string) ([]byte, error)`: codiert aus einer JSONC-Zeichenkette
- `NewDecoder(r io.Reader) Decoder`: erstellt einen Decoder
- `Decode(in []byte, out any) error`: dekodiert in eine Struct
- `DecodeToJSONC(in []byte) (string, error)`: dekodiert in eine JSONC-Zeichenkette

### Beispiele

Siehe das Verzeichnis `examples/` für Beispielcode.
