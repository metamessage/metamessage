package jsonc

import (
	"encoding/json"
	"regexp"

	"github.com/metamessage/metamessage/internal/jsonc/ast"
	"github.com/metamessage/metamessage/internal/jsonc/parser"
	"github.com/metamessage/metamessage/internal/jsonc/scanner"
	"github.com/metamessage/metamessage/internal/jsonc/token"
)

var Email = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var Decimal = regexp.MustCompile(`^-?\d+\.\d+$`)

func GetInt(node ast.Node, path string) (int, error) {
	return 0, nil
}

func GetString(node ast.Node, path string) (string, error) {
	return "", nil
}

func GetFloat(node ast.Node, path string) (float64, error) {
	return 0, nil
}

func Print(n ast.Node) {
	println(ToString(n))
}

func Json(n ast.Node) string {
	b, _ := json.MarshalIndent(n, "", "  ")
	return string(b)
}

func ParseFromString(in string) (out ast.Node, err error) {
	sc := scanner.New(in)
	var toks []token.Token
	for {
		t := sc.NextToken()
		toks = append(toks, t)
		if t.Type == token.EOF {
			break
		}
	}

	p := parser.New(toks)
	out, err = p.Parse()
	if err != nil {
		return
	}
	return
}

func ParseFromBytes(in []byte) (out ast.Node, err error) {
	return ParseFromString(string(in))
}

func BindFromString(in string, out any) (err error) {
	var n ast.Node
	n, err = ParseFromString(in)
	if err != nil {
		return
	}

	return Bind(n, out)
}

func BindFromBytes(in []byte, out any) (err error) {
	var n ast.Node
	n, err = ParseFromBytes(in)
	if err != nil {
		return
	}

	return Bind(n, out)
}

func StructToJSONCString(value any, name string) (string, error) {
	node, err := StructToJSONC(value, name)
	if err != nil {
		return "", err
	}

	return ToString(node), nil
}
