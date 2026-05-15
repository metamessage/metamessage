package core

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/metamessage/metamessage/internal/ir"
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
	var node ir.Node
	node, err = ValueToNode(v, tag)
	if err != nil {
		err = fmt.Errorf("value to node error: %w", err)
		return
	}

	// fmt.Println("FromValue", Dump(node))

	encoder := getEncoder()
	defer putEncoder(encoder)
	return encoder.Encode(node)
}

func FromJSONC(s string) (bs []byte, err error) {
	var node ir.Node
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

func Decode(data []byte) (ir.Node, error) {
	decoder := getDecoder()
	defer putDecoder(decoder)

	return decoder.Decode(data)
}

func ValueToJSONC(value any, name string) (string, error) {
	node, err := ValueToNode(value, name)
	if err != nil {
		return "", err
	}

	return jsonc.ToJSONC(node), nil
}

func BindFromJSONC(in string, out any) (err error) {
	var n ir.Node
	n, err = ParseFromJSONC(in)
	if err != nil {
		return
	}

	return Bind(n, out)
}

func GetInt(node ir.Node, path string) (int, error) {
	return 0, nil
}

func GetString(node ir.Node, path string) (string, error) {
	return "", nil
}

func GetFloat(node ir.Node, path string) (float64, error) {
	return 0, nil
}

func ParseFromJSONC(in string) (out ir.Node, err error) {
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

func PrintJSONC(n ir.Node) {
	println(jsonc.ToJSONC(n))
}

func Dump(n ir.Node) string {
	b, _ := json.MarshalIndent(n, "", "  ")
	return string(b)
}
