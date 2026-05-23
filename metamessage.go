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

// EncodeFromValue encodes the given value into MetaMessage wire format.
//
// The tag parameter provides additional metadata (e.g., "desc=...").
//
// # Pointer and Nullable Behavior
//
// If value is a pointer,
// The encoded output will have the "nullable" tag automatically set.
//
// If value is a non-pointer,
// The encoded output will NOT have the "nullable" tag.
//
// Examples:
//
//		// Non-pointer (no nullable)
//		req := CreateUserRequest{Name: "David"}
//		wire, _ := mm.EncodeFromValue(req, "")
//		// jsonc:
//	    //
//	    {"name": "David"}
//
//		// Pointer to struct (with nullable)
//		req := &CreateUserRequest{Name: "David"}
//		wire, _ := mm.EncodeFromValue(req, "")
//		// jsonc:
//	    //
//		// mm: nullable
//		{"name": "David"}
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
