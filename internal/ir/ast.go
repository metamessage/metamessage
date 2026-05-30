package ir

type NodeType uint8

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeDoc
	NodeTypeObject
	NodeTypeArray
	NodeTypeValue
	NodeTypeNull

	LabelNodeTypeUnknown = "unknown"
	LabelNodeTypeDoc     = "doc"
	LabelNodeTypeObject  = "object"
	LabelNodeTypeArray   = "array"
	LabelNodeTypeValue   = "value"
	LabelNodeTypeNull    = "null"
)

func (nt NodeType) String() string {
	switch nt {
	case NodeTypeDoc:
		return LabelNodeTypeDoc
	case NodeTypeObject:
		return LabelNodeTypeObject
	case NodeTypeArray:
		return LabelNodeTypeArray
	case NodeTypeValue:
		return LabelNodeTypeValue
	case NodeTypeNull:
		return LabelNodeTypeNull
	default:
		return LabelNodeTypeUnknown
	}
}

func ParseNodeType(s string) NodeType {
	switch s {
	case LabelNodeTypeDoc:
		return NodeTypeDoc
	case LabelNodeTypeObject:
		return NodeTypeObject
	case LabelNodeTypeArray:
		return NodeTypeArray
	case LabelNodeTypeValue:
		return NodeTypeValue
	case LabelNodeTypeNull:
		return NodeTypeNull
	default:
		return NodeTypeUnknown
	}
}

type Node interface {
	GetTag() *Tag
	GetType() NodeType

	GetPath() string
	SetPath(path string)
}

type Field struct {
	Key   string
	Value Node
}

type Object struct {
	Fields []*Field
	Tag    *Tag
	Path   string
}

func (o *Object) GetPath() string { return o.Path }

func (o *Object) SetPath(path string) { o.Path = path }

func (o *Object) GetType() NodeType { return NodeTypeObject }

func (o *Object) GetTag() *Tag {
	if o == nil {
		return nil
	}
	if o.Tag != nil {
		return o.Tag
	}
	return nil
}

type Array struct {
	Items []Node
	Tag   *Tag
	Path  string
}

func (a *Array) GetPath() string { return a.Path }

func (a *Array) SetPath(path string) { a.Path = path }

func (a *Array) GetType() NodeType { return NodeTypeArray }

func (a *Array) GetTag() *Tag {
	if a == nil {
		return nil
	}
	if a.Tag != nil {
		return a.Tag
	}
	return nil
}

type Value struct {
	Data any
	Text string
	Tag  *Tag
	Path string
}

func (v *Value) GetPath() string { return v.Path }

func (v *Value) SetPath(path string) { v.Path = path }

func (v *Value) GetType() NodeType { return NodeTypeValue }

func (v *Value) GetTag() *Tag {
	if v == nil {
		return nil
	}
	if v.Tag != nil {
		return v.Tag
	}
	return nil
}

type Doc struct {
	Fields []*Field
	Tag    *Tag
	Path   string
}

func (d *Doc) GetPath() string { return d.Path }

func (d *Doc) SetPath(path string) { d.Path = path }

func (d *Doc) GetType() NodeType { return NodeTypeDoc }

func (d *Doc) GetTag() *Tag {
	if d == nil {
		return nil
	}
	if d.Tag != nil {
		return d.Tag
	}
	return nil
}

type NodeNull struct{}

func (d *NodeNull) GetPath() string { return "" }

func (d *NodeNull) SetPath(path string) {}

func (d *NodeNull) GetType() NodeType { return NodeTypeNull }

func (d *NodeNull) GetTag() *Tag {
	return nil
}
