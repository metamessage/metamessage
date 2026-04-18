package pkg

import (
	"io"

	"github.com/lizongying/meta-message/internal/jsonc"
	"github.com/lizongying/meta-message/internal/jsonc/ast"
	"github.com/lizongying/meta-message/internal/mm"
)

type Encoder interface {
	Reset(io.Writer)
	EncodeStream(in any) (n int, err error)
}

func NewEncoder(w io.Writer) Encoder {
	return mm.NewEncoder(w)
}

func EncodeFromStruct(in any, tag string) (out []byte, err error) {
	return mm.FromStruct(in, tag)
}

func EncodeFromJSONC(in string) (out []byte, err error) {
	return mm.FromJSONC(in)
}

type Decoder interface {
	Reset(r io.Reader)
	DecodeStream(out any) (n int, err error)
}

func NewDecoder(r io.Reader) Decoder {
	return mm.NewDecoder(r)
}

func Decode(in []byte, out any) (err error) {
	var node ast.Node
	node, err = mm.Decode(in)
	if err != nil {
		return
	}

	return jsonc.Bind(node, out)
}

func DecodeToJSONC(in []byte) (str string, err error) {
	var node ast.Node
	node, err = mm.Decode(in)
	if err != nil {
		return
	}

	return jsonc.ToString(node), nil
}
