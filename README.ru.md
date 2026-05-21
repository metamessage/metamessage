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

MetaMessage (mm) — это протокол обмена структурированными данными. Он самодокументируемый, самодостаточный и самоиллюстрируемый, обеспечивающий безпотерянный обмен данными. Это протокол следующего поколения, который нативно поддерживает ИИ, людей и машины.

- Удобен для людей и ИИ
- Поддерживает экспорт/импорт в JSONC (в настоящее время; планируется поддержка YAML/TOML)
- Подходит для конфигурационных файлов и обмена данными
- Работает для традиционных API и сценариев взаимодействия с ИИ
- Поддерживает преобразование между структурами/классами языков и MetaMessage
- Поддерживает генерацию кода для нескольких языков
- Данные содержат тип, ограничения, описание и пример без отдельной документации
- Все метаданные могут обновляться вместе с данными без дополнительной синхронизации
- Структуры и значения остаются согласованными между языками
- Нет потери структуры; парсер адаптируется автоматически и не падает
- Может сериализовываться в компактный бинарный формат для более быстрого парсинга и меньшего размера

**Решаемые проблемы**

- Неопределенные типы, например непонятно, является ли поле uint8
- Неполная структура, например null без информации о внутреннем типе
- Отсутствие правил проверки, поэтому нельзя проверить корректность данных
- Отсутствие примеров и описаний, требуется внешняя документация
- Изменения формата требуют корректировки кодировщика/декодировщика и повторной синхронизации документации

MetaMessage естественно подходит для понимания и взаимодействия с ИИ, решая проблемы неоднозначности и неточности. Он заменяет традиционную документацию API, устные соглашения по формату и ручную синхронизацию версий, делая данные самообъясняющимися и независимыми в развитии.

Примечание:

- Кодирование в настоящее время поддерживает до 65535 байт (64КБ). Этот предел может быть расширен после реализации полной поддержки типов документов.
- Проект находится в активной разработке и тестировании, и пока не рекомендуется для использования в продакшене.
- API и поведение могут еще изменяться, поэтому следите за обновлениями версий.

[meta-message](https://github.com/metamessage/metamessage)

## Текстовые форматы

### JSONC

- Разрешены конечные запятые в массивах или объектах
- Обычные комментарии разрешены
- Комментарии следует писать над полями
- Тег mm должен быть в последней строке
- Оставляйте пустую строку между тегом mm и обычными комментариями для читабельности

**Пример**

```jsonc
{
    // mm: type=datetime; desc=время создания
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## Типы данных

Определяйте типы данных с помощью тега `type=`. Формат: `type=typeIdentifier`, например `type=i` означает целое число.

- doc: Кодирование поддерживает до 65535 байт (64КБ). Эта грань может быть превышена после полной поддержки типов документов
- vec: динамические массивы/срезы, не допускают составные типы
- arr: массив фиксированной длины, не допускает составные типы
- obj: объект/структура, составная структура, соответствует мульти-язычному struct/object
- map: Ключи map должны быть строками, а значения не должны быть составными типами
- str: string
- bytes: массив байтов
- bool: boolean
- i: i; целочисленные литералы не должны содержать десятичную точку
- i8: i8
- i16: i16
- i32: i32
- i64: i64
- u: u
- u8: u8
- u16: u16
- u32: u32
- u64: u64
- f32: float32; числа с плавающей запятой не поддерживают NaN/Inf/-0; литералы с плавающей запятой должны содержать десятичную точку, например 0.0
- f64: float64; числа с плавающей запятой не поддерживают NaN/Inf/-0; литералы с плавающей запятой должны содержать десятичную точку, например 0.0
- bigint: bigint
- datetime: по умолчанию UTC 1970-01-01 00:00:00
- date: 1970-01-01
- time: 00:00:00
- uuid: уникальный идентификатор
- decimal: десятичное, должно передаваться как строка
- ip: IP, поддерживает IPv4/IPv6
- url: URL, должен быть валидным URL
- email: email, должен быть валидным
- enum: enum, значения — строки, разделенные |
- image: изображение, в основе bytes
- video: видео, в основе bytes

## Теги

Теги — это аннотации, метки или атрибуты структур языков программирования или комментарии в текстовых форматах

- is_null: значение null с пустым заполнителем

- desc: краткое описание, применяется ко всем типам. Максимальная длина 65535 бит

- type: тип данных. В текстовых форматах строки, целые числа (int), десятичные числа (float64), срезы, объекты (или похожие структуры) не требуют явных тегов типа, если однозначны, например когда размер массива > 0. В языках программирования, если массивы, maps и другие типы можно определить, то теги типа также не требуются

- raw: в некоторых языках программирования типы данных обычно используют типы-обертки, такие как Java. По умолчанию используются типы-обертки; установите raw, если вы не хотите их использовать. Предварительно решение, может быть удалено в будущих версиях

- nullable: допускается ли null, применяется ко всем типам

- allow_empty: кроме логических типов, другие типы по умолчанию не допускают пустые значения. Когда устанавливается allow_empty, пустые значения разрешены в соответствии с определенными правилами

- unique: применяется только к срезам или массивам, указывает, что элементы не могут повторяться

- default: значение по умолчанию, еще не включено

- example: пример данных, используемый, когда массив, срезы или map пусты, автоматически генерируется пример пустого значения

- min: минимальная емкость для массивов, минимальная длина для строк/байтовых массивов или минимальное значение для чисел (целые числа, десятичные числа, bigint)

- max: максимальная емкость для массивов, максимальная длина для строк/байтовых массивов или максимальное значение для чисел (целые числа, десятичные числа, bigint)

- size: емкость для массивов, фиксированная длина для строк или байтовых массивов

- enum: когда присутствует этот тег, значение по умолчанию имеет тип enum. Тип enum здесь представлен в виде строки и не принимает другие формы

- pattern: регулярное выражение, применяется к строкам

- location: смещение часового пояса, по умолчанию 0, применяется только к типам datetime, диапазон -12 до 14

- version: ограничить версию в uuid; в ip может ограничить ipv4 или ipv6

- mime: тип документа, еще не включено

## Использование

### CLI-инструмент

Этот проект предоставляет инструмент командной строки `mm` для кодирования, декодирования и генерации кода.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Сборка

```bash
make
```

#### Примеры

1. Кодирование JSONC в MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

Или считывать из stdin:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. Декодирование MetaMessage в JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

Или считывать из stdin:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Генерация структур и кода из JSONC

Поддерживаются go, java, ts, kt, py, js, cs, rs, swift, php

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

#### Опции

- -encode, -e: режим кодирования
- -decode, -d: режим декодирования
- -generate, -g: режим генерации кода
- -in, -i: путь к входному файлу (пусто = stdin)
- -out, -o: путь к выходному файлу (пусто = stdout)
- -force, -f: перезаписать выходной файл
- -lang, -l: целевой язык генерации (go, java, ts, kt, py, js, cs, rs, swift, php)

### Использование библиотеки

Проект предоставляет библиотеку Go для программного использования.

#### Установка

```bash
go get github.com/metamessage/metamessage
```

#### Пример

Та же логика применима и к другим языкам

```go
// Неправильный подход
// ID уже имеет тип int64, не нужно указывать type=i64
// В mm не используйте json теги вместе. mm автоматически обрабатывает именование; при необходимости используйте name= в mm
// Age должен использовать нативный тип uint8, поэтому type=u8 можно опустить
type User struct {
	ID       int64  `mm:"type=i64;desc=User ID" json:"id"`
	Name     string `mm:"type=str;desc=User Name;min=1;max=50" json:"name"`
	Email    string `mm:"type=email;desc=Email" json:"email"`
	Age      int    `mm:"type=u8;desc=Age;min=0;max=150" json:"age"`
	IsActive bool   `mm:"type=bool;desc=Is Active" json:"is_active"`
}

// Правильный подход
// Email не имеет нативного типа, поэтому требуется type=email
type User struct {
	ID       int64  `mm:"desc=User ID"`
	Name     string `mm:"desc=User Name;min=1;max=50"`
	Email    string `mm:"type=email;desc=Email"`
	Age      uint8  `mm:"desc=Age;min=0;max=150"`
	IsActive bool   `mm:"desc=Is Active"`
}

user := User{}

// Тег корневого уровня можно указать здесь
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

#### Краткий обзор API

- `NewEncoder(w io.Writer) Encoder`: создает энкодер
- `EncodeFromValue(in any) ([]byte, error)`: кодирует из структуры
- `EncodeFromJsonc(in string) ([]byte, error)`: кодирует из строки JSONC
- `NewDecoder(r io.Reader) Decoder`: создает декодер
- `DecodeToValue(in []byte, out any) error`: декодирует в структуру
- `DecodeToJsonc(in []byte) (string, error)`: декодирует в строку JSONC

### Примеры на других языках

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

### Примеры

См. каталог `examples/` для примеров кода.
