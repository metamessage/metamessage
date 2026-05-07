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

MetaMessage (mm) es un protocolo de intercambio de datos estructurados. Es autodescriptivo, autocontenible y autoejemplificador, lo que permite un intercambio de datos sin pérdidas. Está diseñado como un protocolo universal de próxima generación que admite de forma nativa a IA, humanos y máquinas.

- Amigable para humanos y IA
- Exportación/Importación a JSONC (actualmente; se planea soporte para YAML/TOML)
- Adecuado para archivos de configuración e intercambio de datos
- Funciona para APIs tradicionales y escenarios de interacción con IA
- Admite conversión entre structs/clases de distintos lenguajes y MetaMessage
- Admite generación de código para múltiples lenguajes
- Los datos contienen tipo, restricciones, descripción y ejemplo sin documentación adicional
- Toda la metadata puede actualizarse con los datos, sin necesidad de sincronización externa
- Estructuras y valores permanecen consistentes entre lenguajes
- Sin pérdida de estructura; el analizador se adapta automáticamente y no falla
- Se puede serializar a binario compacto para un análisis más rápido y menor tamaño

**Problemas resueltos**

- Tipos inciertos, como no saber si un campo es uint8
- Estructuras incompletas, como null sin información interna de tipo
- Falta de reglas de validación, por lo que no se puede comprobar la validez de los datos
- Falta de ejemplos o descripciones, obligando a depender de documentación externa
- Cambios de formato que requieren ajustar el codificador/decodificador y re-sincronizar la documentación

MetaMessage es naturalmente adecuado para la comprensión e interacción de IA, resolviendo ambigüedades e imprecisiones. Reemplaza la documentación de API tradicional, los acuerdos de formato verbales y la sincronización manual de versiones, al hacer que los datos sean autoexplicativos y evolucionen de manera independiente.

Nota: Actualmente en desarrollo y pruebas, no se recomienda para uso en producción

[meta-message](https://github.com/metamessage/metamessage)

## Formatos de texto

**JSONC**

- Permite comas finales en arreglos u objetos
- Permite comentarios ordinarios
- Los comentarios deben escribirse encima de los campos
- La etiqueta mm debe ir en la última línea
- Deje una línea vacía entre la etiqueta mm y los comentarios normales para mejorar la legibilidad

**YAML**

**TOML**

## Notas

- Todavía hay muchos errores y pruebas incompletas; no se recomienda su uso en producción
- Los arreglos y slices no permiten tipos compuestos; las claves de map deben ser cadenas y los valores no deben ser tipos compuestos
- Los arreglos/slices vacíos insertan automáticamente un valor de ejemplo
- Los enteros y cadenas no necesitan etiquetas de tipo explícitas
- Las estructuras y slices no necesitan etiquetas de tipo explícitas
- Cuando el tamaño del arreglo es > 0, no se necesitan etiquetas de tipo explícitas
- Los flotantes no admiten NaN/Inf/-0
- Codificación admite hasta 65535 bytes (64KB); esto puede ampliarse en el futuro
- Los literales de punto flotante deben incluir un punto decimal
- Los literales enteros no deben incluir un punto decimal

## Tipos de datos

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
- datetime: UTC por defecto 1970-01-01 00:00:00
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

## Etiquetas

Las etiquetas son anotaciones, etiquetas o atributos de estructuras de lenguajes de programación, o comentarios en formatos de texto

- is_null: valor nulo con marcador de posición vacío
- desc: resumen, se aplica a todos los tipos. Longitud máxima 65535 bits
- type: tipo de dato. En formatos de texto, las cadenas, enteros (int), decimales (float64) y objetos (o estructuras similares) no requieren etiquetas de tipo explícitas cuando no son ambiguas. En lenguajes de programación, si se pueden determinar objetos (o estructuras similares) y maps, los maps tampoco requieren etiquetas de tipo
- raw: en algunos lenguajes de programación, los tipos de datos generalmente usan tipos envueltos, como Java. Se usan tipos envueltos por defecto; establezca en raw si no lo desea. Por determinar, puede ser eliminado en futuras versiones
- nullable: si se permite null, se aplica a todos los tipos
- allow_empty: excepto para tipos booleanos, otros tipos no permiten vacío por defecto. Cuando se establece allow_empty, se permiten valores vacíos siguiendo ciertas reglas
- unique: se aplica solo a slices o arreglos, indica que los elementos no pueden repetirse
- default: valor predeterminado, no está habilitado aún
- example: datos de ejemplo usados cuando un arreglo o map está vacío
- min: capacidad mínima en arreglos, longitud mínima en cadenas/byte arrays, valor mínimo en números
- max: capacidad máxima en arreglos, longitud máxima en cadenas/byte arrays, valor máximo en números
- size: capacidad en arreglos, longitud fija para cadenas o byte arrays
- enum: cuando esta etiqueta está presente, el valor es de tipo enum por defecto. El tipo enum aquí se presenta en forma de cadena y no acepta otras formas
- pattern: expresión regular, se aplica a cadenas
- location: desplazamiento de zona horaria, valor predeterminado 0, se aplica solo a tipos datetime, rango -12 a 14
- version: limita versión en uuid; en ip puede restringir ipv4 o ipv6
- mime: tipo de documento, no está habilitado aún

## Uso

### Herramienta CLI

Este proyecto proporciona una herramienta de línea de comandos `mm` para codificar, decodificar y generar código.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Compilar

```bash
make
```

#### Ejemplos

1. Codificar JSONC a MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

O leer desde stdin:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. Decodificar MetaMessage a JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

O leer desde stdin:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Generar structs y código desde JSONC

Admite go, java, ts, kt, py, js, cs, rs, swift, php

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

#### Opciones

- -encode, -e: modo de codificación
- -decode, -d: modo de decodificación
- -generate, -g: modo de generación de código
- -in, -i: ruta del archivo de entrada (vacío usa stdin)
- -out, -o: ruta del archivo de salida (vacío usa stdout)
- -force, -f: sobrescribir archivo de salida
- -lang, -l: lenguaje objetivo para generación de código (go, java, ts, kt, py, js, cs, rs, swift, php)

### Uso de la librería

El proyecto proporciona una librería Go para uso programático.

#### Instalar

```bash
go get github.com/metamessage/metamessage
```

#### Ejemplo

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

#### Resumen de la API

- `NewEncoder(w io.Writer) Encoder`: crea un codificador
- `EncodeFromValue(in any) ([]byte, error)`: codificar desde struct
- `EncodeFromJSONC(in string) ([]byte, error)`: codificar desde cadena JSONC
- `NewDecoder(r io.Reader) Decoder`: crea un decodificador
- `DecodeToValue(in []byte, out any) error`: decodificar a struct
- `DecodeToJSONC(in []byte) (string, error)`: decodificar a cadena JSONC

### Ejemplos de otros lenguajes

#### Java

```java
import io.github.metamessage.mm.MetaMessage;
import io.github.metamessage.mm.MM;

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
import io.github.metamessage.mm.MetaMessage
import io.github.metamessage.mm.MM

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
import { encode, decode } from 'metamessage';

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

### Ejemplos

Consulta el directorio `examples/` para ver ejemplos de código.
