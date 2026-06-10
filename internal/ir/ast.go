package ir

type NodeType uint8

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeDoc
	NodeTypeObject
	NodeTypeArray
	NodeTypeScalar
	NodeTypeNull

	LabelNodeTypeUnknown = "unknown"
	LabelNodeTypeDoc     = "doc"
	LabelNodeTypeObject  = "object"
	LabelNodeTypeArray   = "array"
	LabelNodeTypeScalar  = "scalar"
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
	case NodeTypeScalar:
		return LabelNodeTypeScalar
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
	case LabelNodeTypeScalar:
		return NodeTypeScalar
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

type NodeObject struct {
	Fields []*Field
	Tag    *Tag
	Path   string
}

func (o *NodeObject) GetPath() string { return o.Path }

func (o *NodeObject) SetPath(path string) { o.Path = path }

func (o *NodeObject) GetType() NodeType { return NodeTypeObject }

func (o *NodeObject) GetTag() *Tag {
	if o == nil {
		return nil
	}
	if o.Tag != nil {
		return o.Tag
	}
	return nil
}

type NodeArray struct {
	Items []Node
	Tag   *Tag
	Path  string
}

func (a *NodeArray) GetPath() string { return a.Path }

func (a *NodeArray) SetPath(path string) { a.Path = path }

func (a *NodeArray) GetType() NodeType { return NodeTypeArray }

func (a *NodeArray) GetTag() *Tag {
	if a == nil {
		return nil
	}
	if a.Tag != nil {
		return a.Tag
	}
	return nil
}

type NodeScalar struct {
	Data any
	Text string
	Tag  *Tag
	Path string
}

func (v *NodeScalar) GetPath() string { return v.Path }

func (v *NodeScalar) SetPath(path string) { v.Path = path }

func (v *NodeScalar) GetType() NodeType { return NodeTypeScalar }

func (v *NodeScalar) GetTag() *Tag {
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

type NodeNull struct {
	Tag *Tag
}

func (d *NodeNull) GetPath() string { return "" }

func (d *NodeNull) SetPath(path string) {}

func (d *NodeNull) GetType() NodeType { return NodeTypeNull }

func (d *NodeNull) GetTag() *Tag {
	if d == nil {
		return nil
	}
	if d.Tag != nil {
		return d.Tag
	}
	return nil
}
