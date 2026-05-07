package mm

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/metamessage/metamessage/internal/ast"
	"github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/jsonc/parser"
	"github.com/metamessage/metamessage/internal/jsonc/scanner"
	"github.com/metamessage/metamessage/internal/jsonc/token"
)

var encoderPool = sync.Pool{
	New: func() any {
		return NewEncoder(nil)
	},
}

func getEncoder() Encoder {
	return encoderPool.Get().(Encoder)
}

func putEncoder(encoder Encoder) {
	encoder.Reset(nil)
	encoderPool.Put(encoder)
}

func FromValue(v any, tag string) (bs []byte, err error) {
	var node ast.Node
	node, err = ValueToMM(v, tag)
	if err != nil {
		err = fmt.Errorf("struct to jsonc error: %w", err)
		return
	}

	// fmt.Println("FromValue", jsonc.Json(node))

	encoder := getEncoder()
	defer putEncoder(encoder)
	return encoder.Encode(node)
}

func FromJSONC(s string) (bs []byte, err error) {
	var node ast.Node
	node, err = ParseFromJSONC(s)
	if err != nil {
		err = fmt.Errorf("parse error: %w", err)
		return
	}

	encoder := getEncoder()
	defer putEncoder(encoder)
	return encoder.Encode(node)
}

var decoderPool = sync.Pool{
	New: func() any {
		return NewDecoder(nil)
	},
}

func getDecoder() Decoder {
	return decoderPool.Get().(Decoder)
}

func putDecoder(decoder Decoder) {
	decoderPool.Put(decoder)
}

func Decode(data []byte) (ast.Node, error) {
	decoder := getDecoder()
	defer putDecoder(decoder)

	return decoder.Decode(data)
}

func ValueToJSONC(value any, name string) (string, error) {
	node, err := ValueToMM(value, name)
	if err != nil {
		return "", err
	}

	return jsonc.ToJSONC(node), nil
}

func BindFromJSONC(in string, out any) (err error) {
	var n ast.Node
	n, err = ParseFromJSONC(in)
	if err != nil {
		return
	}

	return Bind(n, out)
}

func GetInt(node ast.Node, path string) (int, error) {
	return 0, nil
}

func GetString(node ast.Node, path string) (string, error) {
	return "", nil
}

func GetFloat(node ast.Node, path string) (float64, error) {
	return 0, nil
}

func ParseFromJSONC(in string) (out ast.Node, err error) {
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

func PrintJSONC(n ast.Node) {
	println(jsonc.ToJSONC(n))
}

func Dump(n ast.Node) string {
	b, _ := json.MarshalIndent(n, "", "  ")
	return string(b)
}
