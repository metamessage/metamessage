# MetaMessage

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

**Ejemplo**

```jsonc
{
    // mm: type=datetime; desc=fecha de creación
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## Conversión de datos

Admite salida a JSONC, YAML, TOML y otros formatos de texto.

**JSONC**

- Permite comas finales en arreglos u objetos

Estilo de comentario recomendado:

- Se permiten comentarios normales
- Los comentarios deben escribirse encima de los campos
- La etiqueta mm debe ir en la última línea
- Deje una línea vacía entre la etiqueta mm y los comentarios normales para mejorar la legibilidad

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

datetime: UTC por defecto 1970-01-01 00:00:00

## Etiquetas

- is_null: valor nulo con marcador de posición vacío
- example: datos de ejemplo usados cuando un arreglo o map está vacío
- min: capacidad mínima en arreglos, longitud mínima en cadenas/byte arrays, valor mínimo en números
- max: capacidad máxima en arreglos, longitud máxima en cadenas/byte arrays, valor máximo en números
- size: longitud fija para arreglos, cadenas o byte arrays
- location: desplazamiento de zona horaria, valor predeterminado 0, rango -12 a 14

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
go get github.com/metamessage/metamessage/pkg
```

#### Ejemplo

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

#### Resumen de la API

- `NewEncoder(w io.Writer) Encoder`: crea un codificador
- `EncodeFromStruct(in any) ([]byte, error)`: codificar desde struct
- `EncodeFromJSONC(in string) ([]byte, error)`: codificar desde cadena JSONC
- `NewDecoder(r io.Reader) Decoder`: crea un decodificador
- `Decode(in []byte, out any) error`: decodificar a struct
- `DecodeToJSONC(in []byte) (string, error)`: decodificar a cadena JSONC

### Ejemplos

Consulta el directorio `examples/` para ver ejemplos de código.
