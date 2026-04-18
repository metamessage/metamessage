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

## Conversion de données

Prend en charge la sortie vers JSONC, YAML, TOML et d'autres formats texte.

**JSONC**

- Autorise les virgules finales dans les tableaux ou objets

Style de commentaire recommandé :

- Les commentaires ordinaires sont autorisés
- Les commentaires doivent être écrits au-dessus des champs
- Le tag mm doit être sur la dernière ligne
- Laisser une ligne vide entre le tag mm et les commentaires ordinaires pour une meilleure lisibilité

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

datetime : UTC par défaut 1970-01-01 00:00:00

## Tags

- is_null : indique une valeur null avec un espace réservé vide
- example : données d'exemple utilisées lorsque les tableaux ou maps sont vides
- min : capacité minimale pour les tableaux, longueur minimale pour les chaînes/octets, ou valeur minimale pour les nombres
- max : capacité maximale pour les tableaux, longueur maximale pour les chaînes/octets, ou valeur maximale pour les nombres
- size : taille fixe pour les tableaux, chaînes ou octets
- location : décalage de fuseau horaire, défaut 0, plage -12 à 14

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
go get github.com/metamessage/metamessage/pkg
```

#### Exemple

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

#### Résumé de l'API

- `NewEncoder(w io.Writer) Encoder` : crée un encodeur
- `EncodeFromStruct(in any) ([]byte, error)` : encode à partir d'une struct
- `EncodeFromJSONC(in string) ([]byte, error)` : encode à partir d'une chaîne JSONC
- `NewDecoder(r io.Reader) Decoder` : crée un décodeur
- `Decode(in []byte, out any) error` : décode vers une struct
- `DecodeToJSONC(in []byte) (string, error)` : décode vers une chaîne JSONC

### Exemples

Voir le répertoire `examples/` pour des exemples de code.
