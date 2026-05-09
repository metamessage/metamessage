package metamessage

import (
	"io"

	"github.com/metamessage/metamessage/internal/ast"
	jc "github.com/metamessage/metamessage/internal/jsonc"
	"github.com/metamessage/metamessage/internal/mm"
)

type Encoder interface {
	Reset(io.Writer)
	EncodeStream(value any) (n int, err error)
}

func NewEncoder(w io.Writer) Encoder {
	return mm.NewEncoder(w)
}

func EncodeFromValue(value any, tag string) (wire []byte, err error) {
	return mm.FromValue(value, tag)
}

func EncodeFromJSONC(jsonc string) (wire []byte, err error) {
	return mm.FromJSONC(jsonc)
}

type Decoder interface {
	Reset(r io.Reader)
	DecodeStream(value any) (n int, err error)
}

func NewDecoder(r io.Reader) Decoder {
	return mm.NewDecoder(r)
}

func DecodeToValue(wire []byte, value any) (err error) {
	var node ast.Node
	node, err = mm.Decode(wire)
	if err != nil {
		return
	}

	return mm.Bind(node, value)
}

func DecodeToJSONC(wire []byte) (jsonc string, err error) {
	var node ast.Node
	node, err = mm.Decode(wire)
	if err != nil {
		return
	}

	return jc.ToJSONC(node), nil
}
