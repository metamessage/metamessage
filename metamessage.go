package metamessage

import (
	"io"

	"github.com/metamessage/metamessage/internal/core"
	"github.com/metamessage/metamessage/internal/ir"
	jc "github.com/metamessage/metamessage/internal/jsonc"
)

type Encoder interface {
	Reset(io.Writer)
	EncodeStream(value any) (n int, err error)
}

func NewEncoder(w io.Writer) Encoder {
	return core.NewEncoder(w)
}

func EncodeFromValue(value any, tag string) (wire []byte, err error) {
	return core.FromValue(value, tag)
}

func EncodeFromJsonc(jsonc string) (wire []byte, err error) {
	return core.FromJSONC(jsonc)
}

type Decoder interface {
	Reset(r io.Reader)
	DecodeStream(value any) (n int, err error)
}

func NewDecoder(r io.Reader) Decoder {
	return core.NewDecoder(r)
}

func DecodeToValue(wire []byte, value any) (err error) {
	var node ir.Node
	node, err = core.Decode(wire)
	if err != nil {
		return
	}

	return core.Bind(node, value)
}

func DecodeToJsonc(wire []byte) (jsonc string, err error) {
	var node ir.Node
	node, err = core.Decode(wire)
	if err != nil {
		return
	}

	return jc.ToJSONC(node), nil
}

func ValueToJsonc(value any, tag string) (jsonc string, err error) {
	var node ir.Node
	node, err = core.ValueToNode(value, tag)
	if err != nil {
		return
	}

	return jc.ToJSONC(node), nil
}

func JsoncToValue(jsonc string, value any) (err error) {
	var node ir.Node
	node, err = core.ParseFromJSONC(jsonc)
	if err != nil {
		return
	}

	return core.Bind(node, value)
}
