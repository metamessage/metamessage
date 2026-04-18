package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/metamessage/metamessage/internal/gen"
	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/mm"
)

// go run ./cmd/mm -generate -lang go -in example.jsonc / cat example.jsonc | go run ./cmd/mm -generate -lang go
func main() {
	validLangs := []string{"go", "java", "ts", "kt", "py", "js", "cs", "rs", "swift", "php"}

	encode := flag.Bool("encode", false, "encode mode: jsonc -> MetaMessage")
	flag.BoolVar(encode, "e", false, "shorthand for -encode")

	decode := flag.Bool("decode", false, "decode mode: MetaMessage -> jsonc")
	flag.BoolVar(decode, "d", false, "shorthand for -decode")

	generate := flag.Bool("generate", false, "generate mode: jsonc -> struct (support multi language)")
	flag.BoolVar(generate, "g", false, "shorthand for -gen")

	in := flag.String("in", "", "input file path (empty = read from stdin)")
	flag.StringVar(in, "i", "", "shorthand for -in")

	out := flag.String("out", "", "output file path (empty = write to stdout)")
	flag.StringVar(out, "o", "", "shorthand for -out")

	force := flag.Bool("force", false, "overwrite output file if it already exists (default: false)")
	flag.BoolVar(force, "f", false, "shorthand for -force")

	lang := flag.String("lang", "none", fmt.Sprintf("generate target language (only valid for -gen, default: none, support: %s)", strings.Join(validLangs, ", ")))
	flag.StringVar(lang, "l", "none", "shorthand for -lang")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "\nMode (mutually exclusive, choose one):")
		fmt.Fprintln(os.Stderr, "  -encode, -e        Encode JSONC to MetaMessage format")
		fmt.Fprintln(os.Stderr, "  -decode, -d        Decode MetaMessage to JSONC format")
		fmt.Fprintln(os.Stderr, "  -generate, -g      Generate struct code from JSONC")
		fmt.Fprintln(os.Stderr, "\nCommon Options:")
		fmt.Fprintln(os.Stderr, "  -in, -i string     Input file path (empty = read from stdin)")
		fmt.Fprintln(os.Stderr, "  -out, -o string    Output file path (empty = write to stdout)")
		fmt.Fprintln(os.Stderr, "  -force, -f         Overwrite output file if it exists (default: false)")
		fmt.Fprintln(os.Stderr, "\nGenerate Options (only for -gen):")
		fmt.Fprintln(os.Stderr, "  -lang, -l string   Target language (default: none, support: go, java, ts, kt, py, js, cs, rs, swift, php)")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  # Encode JSONC to MetaMessage (stdin -> stdout)")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-encode -in input.jsonc -out output.MetaMessage")
		fmt.Fprintln(os.Stderr, "  # Decode MetaMessage to JSONC (stdin -> stdout)")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-decode < input.MetaMessage > output.jsonc")
		fmt.Fprintln(os.Stderr, "  # Generate Go struct from JSONC")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-gen -lang go -in input.jsonc -out output.go")
		fmt.Fprintln(os.Stderr, "  # Generate C# class from JSONC")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-gen -lang cs -in input.jsonc -out output.cs")
		fmt.Fprintln(os.Stderr, "  # Generate Rust struct from JSONC")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-gen -lang rs -in input.jsonc -out output.rs")
		fmt.Fprintln(os.Stderr, "  # Generate Swift struct from JSONC")
		fmt.Fprintln(os.Stderr, "  ", os.Args[0], "-gen -lang swift -in input.jsonc -out output.swift")
	}
	flag.Parse()

	modeCount := 0
	if *encode {
		modeCount++
	}
	if *decode {
		modeCount++
	}
	if *generate {
		modeCount++
	}

	if modeCount == 0 {
		fmt.Fprintln(os.Stderr, "Error: A mode must be specified! Valid options: -encode / -decode / -generate")
		flag.Usage()
		os.Exit(1)
	}

	if modeCount > 1 {
		fmt.Fprintln(os.Stderr, "Error: Only one mode can be selected! Valid options: -encode / -decode / -generate")
		flag.Usage()
		os.Exit(1)
	}

	if *generate {
		if !slices.Contains(validLangs, *lang) {
			fmt.Fprintf(os.Stderr, "Error: Unsupported language %s! Valid options: %s, all\n", *lang, strings.Join(validLangs, ", "))
			os.Exit(1)
		}
	}

	var data []byte
	var err error
	if *in == "" {
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read stdin: %v\n", err)
			os.Exit(2)
		}
	} else {
		data, err = os.ReadFile(*in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read file: %v\n", err)
			os.Exit(2)
		}
	}

	switch {
	case *encode:
		fmt.Printf("Encoding Mode, Input: %s, Output: %s\n", *in, *out)
		output, e := mm.FromJSONCBytes(data)
		if e != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(2)
		}

		if *out == "" {
			fmt.Println(base64.StdEncoding.EncodeToString(output))
		} else {
			outDir := filepath.Dir(*out)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Create directory failed: %v\n", err)
				os.Exit(1)
			}

			if fileInfo, err := os.Stat(*out); err == nil {
				if fileInfo.Mode().IsRegular() && !*force {
					fmt.Fprintf(os.Stderr, "Error: Output file '%s' exists. Use -force to overwrite.\n", *out)
					os.Exit(1)
				} else if *force {
					fmt.Fprintf(os.Stderr, "Warning: Overwriting '%s'\n", *out)
				}
			} else if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Check output file failed: %v\n", err)
				os.Exit(1)
			}

			if err = os.WriteFile(*out, output, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write file: %v\n", err)
				os.Exit(2)
			}
		}

	case *decode:
		fmt.Printf("Decoding Mode, Input: %s, Output:%s\n", *in, *out)
		outputStr, e := mm.ToJSONC(data)
		if e != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(2)
		}

		if *out == "" {
			fmt.Println(outputStr)
		} else {
			outDir := filepath.Dir(*out)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Create directory failed: %v\n", err)
				os.Exit(1)
			}

			if fileInfo, err := os.Stat(*out); err == nil {
				if fileInfo.Mode().IsRegular() && !*force {
					fmt.Fprintf(os.Stderr, "Error: Output file '%s' exists. Use -force to overwrite.\n", *out)
					os.Exit(1)
				} else if *force {
					fmt.Fprintf(os.Stderr, "Warning: Overwriting '%s'\n", *out)
				}
			} else if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Check output file failed: %v\n", err)
				os.Exit(1)
			}

			if err = os.WriteFile(*out, []byte(outputStr), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write file: %v\n", err)
				os.Exit(2)
			}
		}

	case *generate:
		fmt.Printf("Generation Mode, Input: %s, Output: %s, Target Language: %s\n", *in, *out, *lang)
		node, err := jsonc.ParseFromBytes(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			os.Exit(2)
		}

		var outputStr string
		switch *lang {
		case "go":
			outputStr = gen.ToGo(node)
		case "java":
			outputStr = gen.ToJava(node)
		case "ts":
			outputStr = gen.ToTS(node)
		case "kt":
			outputStr = gen.ToKotlin(node)
		case "py":
			outputStr = gen.ToPy(node)
		case "js":
			outputStr = gen.ToJS(node)
		case "cs":
			outputStr = gen.ToCSharp(node)
		case "rs":
			outputStr = gen.ToRust(node)
		case "swift":
			outputStr = gen.ToSwift(node)
		case "php":
			outputStr = gen.ToPHP(node)
		default:
			fmt.Fprintf(os.Stderr, "unsupported language: %s\n", *lang)
			fmt.Fprintf(os.Stderr, "supported languages: %s\n", strings.Join(validLangs, ", "))
			os.Exit(2)
		}

		if *out == "" {
			fmt.Println(outputStr)
		} else {
			outDir := filepath.Dir(*out)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error: Create directory failed: %v\n", err)
				os.Exit(1)
			}

			if fileInfo, err := os.Stat(*out); err == nil {
				if fileInfo.Mode().IsRegular() && !*force {
					fmt.Fprintf(os.Stderr, "Error: Output file '%s' exists. Use -force to overwrite.\n", *out)
					os.Exit(1)
				} else if *force {
					fmt.Fprintf(os.Stderr, "Warning: Overwriting '%s'\n", *out)
				}
			} else if !os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Check output file failed: %v\n", err)
				os.Exit(1)
			}

			if err = os.WriteFile(*out, []byte(outputStr), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "write file: %v\n", err)
				os.Exit(2)
			}
		}
	}
}
