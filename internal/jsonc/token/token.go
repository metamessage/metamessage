package token

type Type int

const (
	EOF Type = iota

	// symbols
	LBrace   // {
	RBrace   // }
	LBracket // [
	RBracket // ]
	Colon    // :
	Comma    // ,

	// literals
	String
	Number
	True
	False
	Null

	// comments
	Comment
)

func (t Type) String() string {
	switch t {
	case EOF:
		return "EOF"
	case LBrace:
		return "{"
	case RBrace:
		return "}"
	case LBracket:
		return "["
	case RBracket:
		return "]"
	case Colon:
		return ":"
	case Comma:
		return ","
	case String:
		return "String"
	case Number:
		return "Number"
	case True:
		return "True"
	case False:
		return "False"
	case Null:
		return "Null"
	case Comment:
		return "Comment"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	switch t.Type {
	case EOF:
		return "EOF"
	case LBrace:
		return "{"
	case RBrace:
		return "}"
	case LBracket:
		return "["
	case RBracket:
		return "]"
	case Colon:
		return ":"
	case Comma:
		return ","
	case String:
		return "String(\"" + t.Literal + "\")"
	case Number:
		return "Number(" + t.Literal + ")"
	case True:
		return "True"
	case False:
		return "False"
	case Null:
		return "Null"
	case Comment:
		return "Comment(" + t.Literal + ")"
	default:
		return "Unknown"
	}
}
