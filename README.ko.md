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

MetaMessage(mm)은 구조화된 데이터 교환 프로토콜입니다. 자기 설명, 자기 제약, 자기 예시를 갖추어 손실 없는 데이터 교환을 제공합니다. AI, 사람, 기계에 네이티브로 적합한 차세대 범용 프로토콜입니다.

- 사람과 AI에 친화적
- JSONC로 내보내기/가져오기 지원(현재는 JSONC, 추후 YAML/TOML 등 지원 예정)
- 구성 파일과 데이터 교환에 적합
- 기존 API 및 AI 상호작용 시나리오에 모두 적합
- 언어별 구조체/클래스와 MetaMessage 간 상호 변환 지원
- 여러 언어의 구조체 및 데이터 코드 자동 생성 지원
- 데이터 자체에 타입, 제약, 설명, 예시가 포함되어 별도 문서가 필요 없음
- 타입, 제약, 설명, 예시는 데이터와 함께 업데이트되어 별도 동기화가 필요 없음
- 언어 간 데이터 구조와 값이 완전히 일치
- 구조 손실이 없고 파서는 자동 적응하며 절대 크래시하지 않음
- 컴팩트한 바이너리로 직렬화하여 더 빠르게 파싱하고 더 작은 크기를 유지

**해결하는 문제**

- uint8인지 알 수 없는 유형 불명확성
- null 내부 구조를 알 수 없어 구조가 불완전해짐
- 제약 규칙 부재로 데이터 유효성을 검증할 수 없음
- 예시나 설명이 없어 외부 문서에 의존해야 함
- 형식 변경에 민감하여 인코딩/디코딩과 문서 재동기화가 필요함

MetaMessage는 AI의 이해와 상호작용에 자연스럽게 맞으며, 모호성 및 부정확성을 해결합니다. 전통적인 인터페이스 문서, 구두 포맷 합의, 수동 버전 동기화를 대체하여 데이터 자체가 자가 설명적이고 독립적으로 진화할 수 있도록 합니다.

참고: 현재 개발 및 테스트 중이며, 프로덕션 사용을 권장하지 않습니다

[meta-message](https://github.com/metamessage/metamessage)

## 텍스트 형식

### JSONC

- 배열 또는 객체의 끝에 쉼표를 허용
- 일반 주석 허용
- 필드 위에 주석 작성 권장
- mm 태그는 마지막 줄에 있어야 함
- 읽기 쉽도록 mm 태그와 일반 주석 사이에 빈 줄을 둠

**예제**

```jsonc
{
    // mm: type=datetime; desc=생성 시간
    "create_time": "2026-01-01 00:00:00"
}
```

### YAML

### TOML

## 데이터 타입

- doc: 인코딩은 최대 65535바이트(64KB)까지 지원합니다. 문서 타입의 완전한 지원 후 이 제한을 초과할 수 있습니다
- slice: 배열과 슬라이스는 복합 타입을 허용하지 않습니다
- array: arr
- struct:
- map: 맵 키는 문자열이어야 하며 값은 복합 타입이어서는 안 됩니다
- string: str
- bytes:
- bool:
- int: i; 정수 리터럴에는 소수점을 포함해서는 안 됩니다
- int8: i8
- int16: i16
- int32: i32
- int64: i64
- uint: u
- uint8: u8
- uint16: u16
- uint32: u32
- uint64: u64
- float32: f32; 부동 소수점은 NaN/Inf/-0을 지원하지 않습니다. 부동 소수점 리터럴에는 소수점이 있어야 합니다(예: 0.0)
- float64: f64
- bigint: bi
- datetime: 기본 UTC 1970-01-01 00:00:00
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

## 태그

태그는 프로그래밍 언어 구조체의 주석, 레이블 또는 속성, 또는 텍스트 형식의 주석입니다

- is_null: null 값을 빈 자리 표시자로 표현

- desc: 요약으로 모든 타입에 적용. 최대 길이 65535 비트

- type: 데이터 타입. 텍스트 형식에서는 문자열, 정수(int), 소수(float64), 슬라이스, 객체(또는 유사한 구조)가 모호하지 않을 때는 명시적 타입 태그가 필요 없습니다(예: 배열 크기 > 0). 프로그래밍 언어에서 배열, map 등을 판별할 수 있으면 타입 태그가 필요 없습니다

- raw: 일부 프로그래밍 언어에서는 Java와 같이 데이터 타입이 일반적으로 래퍼 타입을 사용합니다. 기본적으로 래퍼 타입을 사용하며, 사용하지 않으려면 raw로 설정합니다. 미정, 향후 삭제될 수 있습니다

- nullable: null 가능 여부로 모든 타입에 적용

- allow_empty: 부울 타입을 제외하고 다른 타입은 기본적으로 빈 값을 허용하지 않습니다. allow_empty가 설정되면 일부 규칙에 따라 빈 값이 허용됩니다

- unique: 슬라이스 또는 배열에만 적용되며 요소가 중복될 수 없음을 나타냅니다

- default: 기본값으로 아직 활성화되지 않음

- example: 배열, 슬라이스, map이 비어 있을 때 사용하는 샘플 데이터로 자동으로 빈 값의 예시가 생성됩니다

- min: 배열의 최소 용량, 문자열/바이트 배열의 최소 길이, 숫자(정수, 소수, bigint)의 최소값

- max: 배열의 최대 용량, 문자열/바이트 배열의 최대 길이, 숫자(정수, 소수, bigint)의 최대값

- size: 배열의 용량, 문자열/바이트 배열의 고정 길이

- enum: 이 태그가 있으면 값이 기본적으로 enum 타입입니다. 여기서 enum 타입은 문자열 형태로 표현되며 다른 형태를 허용하지 않습니다

- pattern: 정규표현식으로 문자열에 적용

- location: 시간대 오프셋으로 기본값 0, datetime 타입에만 적용, 범위 -12에서 14

- version: uuid에서 버전을 제한하며 ip에서는 ipv4 또는 ipv6을 제한할 수 있습니다

- mime: 문서 타입으로 아직 활성화되지 않음

## 사용 방법

### CLI 도구

이 프로젝트는 인코딩, 디코딩 및 코드 생성을 위한 명령줄 도구 `mm`를 제공합니다.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### 빌드

```bash
make
```

#### 사용 예시

1. JSONC를 MetaMessage로 인코딩

```bash
./mm -encode -in input.jsonc -out output.mm
```

또는 stdin에서 읽기:

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. MetaMessage를 JSONC로 디코딩

```bash
./mm -decode -in input.mm -out output.jsonc
```

또는 stdin에서 읽기:

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. JSONC에서 구조체 및 코드 생성

지원 언어: go, java, ts, kt, py, js, cs, rs, swift, php

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

#### 옵션

- -encode, -e: 인코딩 모드
- -decode, -d: 디코딩 모드
- -generate, -g: 코드 생성 모드
- -in, -i: 입력 파일 경로(비워두면 stdin 사용)
- -out, -o: 출력 파일 경로(비워두면 stdout 사용)
- -force, -f: 출력 파일 덮어쓰기
- -lang, -l: 생성 대상 언어(go, java, ts, kt, py, js, cs, rs, swift, php)

### 라이브러리 사용

프로그램에서 사용하기 위한 Go 라이브러리를 제공합니다.

#### 설치

```bash
go get github.com/metamessage/metamessage
```

#### 예시 코드

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

#### API 요약

- `NewEncoder(w io.Writer) Encoder`: 인코더 생성
- `EncodeFromStruct(in any) ([]byte, error)`: 구조체에서 인코딩
- `EncodeFromJSONC(in string) ([]byte, error)`: JSONC 문자열에서 인코딩
- `NewDecoder(r io.Reader) Decoder`: 디코더 생성
- `Decode(in []byte, out any) error`: 구조체로 디코딩
- `DecodeToJSONC(in []byte) (string, error)`: JSONC 문자열로 디코딩

### 예시

`examples/` 디렉터리의 샘플 코드를 참조하세요.
