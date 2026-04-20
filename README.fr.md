# MetaMessage

MetaMessage (mm) est un protocole d'échange de données structurées. Il est auto-descriptif, auto-contraint et auto-exemplifié, permettant un échange de données sans perte. Il est conçu comme un protocole universel de nouvelle génération prenant en charge nativement l'IA, les humains et les machines.

- Convient aux humains et à l'IA
- Exportation/importation vers JSONC (actuellement ; prise en charge YAML/TOML prévue)
- Adapté aux fichiers de configuration et à l'échange de données
- Convient aux API traditionnelles et aux scénarios d'interaction IA
- Prend en charge la conversion entre structures/classes de langages et MetaMessage
- Prend en charge la génération de code pour plusieurs langages
- Les données contiennent le type, les contraintes, la description et un exemple sans documentation séparée
- Toutes les métadonnées peuvent être mises à jour avec les données, sans coordination supplémentaire
- Les structures et les valeurs restent cohérentes entre les langages
- Aucune perte de structure ; l'analyseur s'adapte automatiquement et ne plante pas
- Peut être sérialisé en binaire compact pour un décodage plus rapide et une taille plus petite

**Problèmes résolus**

- Types inconnus, par exemple ne pas savoir si un champ est uint8
- Structure incomplète, par exemple null sans information de type interne
- Pas de règles de validation, donc impossibilité de vérifier la validité des données
- Pas d'exemples ni de descriptions, obligeant à dépendre de documentation externe
- Les modifications de format nécessitent un ajustement du codage/décodage et une resynchronisation de la documentation

MetaMessage est naturellement adapté à la compréhension et à l'interaction avec l'IA, résolvant l'ambiguïté et l'imprécision. Il remplace la documentation d'API traditionnelle, les accords verbaux de format et la synchronisation manuelle des versions en rendant les données auto-explicatives et indépendamment évolutives.

**Exemple**

```jsonc
{
    // mm: type=datetime; desc=heure de création
    "create_time": "2026-01-01 00:00:00"
}
```

[meta-message](https://github.com/metamessage/metamessage)

## Formats texte

**JSONC**

- Autorise les virgules finales dans les tableaux ou objets
- Autorise les commentaires ordinaires
- Les commentaires doivent être écrits au-dessus des champs
- Le tag mm doit être sur la dernière ligne
- Laisser une ligne vide entre le tag mm et les commentaires ordinaires pour une meilleure lisibilité

**YAML**

**TOML**

## Notes

- Il reste encore de nombreux bugs et les tests sont incomplets ; l'utilisation en production n'est pas recommandée
- Les tableaux et slices n'autorisent pas les types composites ; les clés de map doivent être des chaînes et les valeurs ne doivent pas être des types composites
- Les tableaux/slices vides insèrent automatiquement une valeur d'exemple
- Les entiers et les chaînes ne nécessitent pas d'étiquettes de type explicites
- Les structs et slices ne nécessitent pas d'étiquettes de type explicites
- Lorsque la taille du tableau est > 0, les étiquettes de type explicites ne sont pas nécessaires
- Les floats ne prennent pas en charge NaN/Inf/-0
- L'encodage prend en charge jusqu'à 65535 octets (64KB) ; cela pourrait être étendu ultérieurement
- Les littéraux à virgule flottante doivent inclure un point décimal
- Les littéraux entiers ne doivent pas inclure de point décimal

## Types de données

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
- datetime: UTC par défaut 1970-01-01 00:00:00
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

Les tags sont des annotations, étiquettes ou attributs de structures de langages de programmation, ou des commentaires dans les formats texte

- is_null : indique une valeur null avec un espace réservé vide
- desc: résumé, s'applique à tous les types. Longueur maximale 65535 bits
- type: type de données. Dans les formats texte, les chaînes, entiers (int), décimales (float64) et objets (ou structures similaires) ne nécessitent pas d'étiquettes de type explicites lorsque sans ambiguïté. Dans les langages de programmation, si les objets (ou structures similaires) et maps peuvent être déterminés, les maps ne nécessitent pas non plus d'étiquettes de type
- raw: dans certains langages de programmation, les types de données utilisent généralement des types wrapper, comme Java. Les types wrapper sont utilisés par défaut; définissez sur raw si vous ne le souhaitez pas. À déterminer, peut être supprimé dans les versions futures
- nullable: si null est autorisé, s'applique à tous les types
- allow_empty: sauf pour les types booléens, les autres types ne permettent pas le vide par défaut. Lorsque allow_empty est défini, les valeurs vides sont autorisées selon certaines règles
- unique: s'applique uniquement aux slices ou tableaux, indique que les éléments ne peuvent pas être répétés
- default: valeur par défaut, non encore activée
- example: données d'exemple utilisées lorsque les tableaux ou maps sont vides
- min: capacité minimale pour les tableaux, longueur minimale pour les chaînes/octets, ou valeur minimale pour les nombres
- max: capacité maximale pour les tableaux, longueur maximale pour les chaînes/octets, ou valeur maximale pour les nombres
- size: capacité pour les tableaux, longueur fixe pour les chaînes ou octets
- enum: quand cette étiquette est présente, la valeur est du type enum par défaut. Le type enum ici est sous forme de chaîne et n'accepte pas d'autres formes
- pattern: regex, s'applique aux chaînes
- location: décalage de fuseau horaire, défaut 0, s'applique uniquement aux types datetime, plage -12 à 14
- version: limiter la version dans uuid; dans ip peut restreindre ipv4 ou ipv6
- mime: type de document, non encore activé

## Utilisation

### Outil CLI

Ce projet fournit un outil en ligne de commande `mm` pour l'encodage, le décodage et la génération de code.

[releases](https://github.com/metamessage/metamessage/releases/latest)

#### Compilation

```bash
make
```

#### Exemples

1. Encoder JSONC vers MetaMessage

```bash
./mm -encode -in input.jsonc -out output.mm
```

Ou lire depuis stdin :

```bash
cat input.jsonc | ./mm -encode > output.mm
```

2. Décoder MetaMessage vers JSONC

```bash
./mm -decode -in input.mm -out output.jsonc
```

Ou lire depuis stdin :

```bash
cat input.mm | ./mm -decode > output.jsonc
```

3. Générer des structs et du code depuis JSONC

Prend en charge go, java, ts, kt, py, js, cs, rs, swift, php

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

#### Options

- -encode, -e : mode encodage
- -decode, -d : mode décodage
- -generate, -g : mode génération de code
- -in, -i : chemin du fichier d'entrée (vide pour stdin)
- -out, -o : chemin du fichier de sortie (vide pour stdout)
- -force, -f : écraser le fichier de sortie
- -lang, -l : langue cible de génération (go, java, ts, kt, py, js, cs, rs, swift, php)

### Utilisation de la bibliothèque

Le projet fournit une bibliothèque Go pour un usage programmatique.

#### Installation

```bash
go get github.com/metamessage/metamessage
```

#### Exemple

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

#### Résumé de l'API

- `NewEncoder(w io.Writer) Encoder` : crée un encodeur
- `EncodeFromStruct(in any) ([]byte, error)` : encode à partir d'une struct
- `EncodeFromJSONC(in string) ([]byte, error)` : encode à partir d'une chaîne JSONC
- `NewDecoder(r io.Reader) Decoder` : crée un décodeur
- `Decode(in []byte, out any) error` : décode vers une struct
- `DecodeToJSONC(in []byte) (string, error)` : décode vers une chaîne JSONC

### Exemples

Voir le répertoire `examples/` pour des exemples de code.
