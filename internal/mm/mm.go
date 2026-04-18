package mm

import (
	"fmt"
	"sync"

	"github.com/lizongying/meta-message/internal/jsonc"
	"github.com/lizongying/meta-message/internal/jsonc/ast"
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

func FromStruct(v any, tag string) (bs []byte, err error) {
	var node ast.Node
	node, err = jsonc.StructToJSONC(v, tag)
	if err != nil {
		return
	}
	fmt.Println("FromStruct", jsonc.Json(node))

	encoder := getEncoder()
	defer putEncoder(encoder)
	return encoder.Encode(node)
}

func FromJSONC(s string) (bs []byte, err error) {
	var node ast.Node
	node, err = jsonc.ParseFromString(s)
	if err != nil {
		return
	}

	encoder := getEncoder()
	defer putEncoder(encoder)
	return encoder.Encode(node)
}

func FromJSONCBytes(bs []byte) ([]byte, error) {
	return FromJSONC(string(bs))
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

func ToJSONC(data []byte) (str string, err error) {
	var node ast.Node
	node, err = Decode(data)
	if err != nil {
		return
	}

	return jsonc.ToString(node), nil
}

func ToJSONCBytes(data []byte) (bs []byte, err error) {
	str, err := ToJSONC(data)
	if err != nil {
		return
	}

	return []byte(str), nil
}

func Bind(in []byte, out any) (err error) {
	var n ast.Node
	n, err = Decode(in)
	if err != nil {
		return
	}

	return jsonc.Bind(n, out)
}
