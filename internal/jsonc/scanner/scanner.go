package scanner

import (
	"strings"
	"unicode"

	"github.com/lizongying/meta-message/internal/jsonc/ast"
	"github.com/lizongying/meta-message/internal/jsonc/token"
)

type Scanner struct {
	src     []rune
	pos     int
	line    int
	col     int
	newLine bool
}

func New(input string) *Scanner {
	return &Scanner{
		src:  []rune(input),
		line: 1,
		col:  1,
	}
}

func (s *Scanner) peek() rune {
	if s.pos >= len(s.src) {
		return 0
	}
	return s.src[s.pos]
}

func (s *Scanner) next() rune {
	if s.pos >= len(s.src) {
		return 0
	}
	ch := s.src[s.pos]
	s.pos++

	if ch == '\n' {
		s.newLine = true
		s.line++
		s.col = 1
	} else {
		s.col++
	}
	return ch
}

func (s *Scanner) skipWhitespace() {
	for unicode.IsSpace(s.peek()) {
		s.next()
	}
}

func (s *Scanner) NextToken() token.Token {
	s.skipWhitespace()

	ch := s.peek()
	if ch == 0 {
		return token.Token{Type: token.EOF, Line: s.line, Column: s.col}
	}

	startLine, startCol := s.line, s.col

	switch ch {
	case '{':
		s.next()
		return token.Token{Type: token.LBrace, Line: startLine, Column: startCol}
	case '}':
		s.next()
		return token.Token{Type: token.RBrace, Line: startLine, Column: startCol}
	case '[':
		s.next()
		return token.Token{Type: token.LBracket, Line: startLine, Column: startCol}
	case ']':
		s.next()
		return token.Token{Type: token.RBracket, Line: startLine, Column: startCol}
	case ':':
		s.next()
		s.newLine = false
		return token.Token{Type: token.Colon, Line: startLine, Column: startCol}
	case ',':
		s.next()
		s.newLine = false
		return token.Token{Type: token.Comma, Line: startLine, Column: startCol}
	case '"':
		return s.scanString()
	case '/':
		return s.scanComment()
	default:
		return s.scanLiteral()
	}
}

func (s *Scanner) scanString() token.Token {
	startLine, startCol := s.line, s.col
	s.next() // "

	var buf strings.Builder
	for {
		ch := s.next()
		if ch == 0 || ch == '\n' {
			break
		}
		if ch == '"' {
			break
		}
		if ch == '\\' {
			buf.WriteRune(ch)
			buf.WriteRune(s.next())
			continue
		}
		buf.WriteRune(ch)
	}

	return token.Token{
		Type:    token.String,
		Literal: buf.String(),
		Line:    startLine,
		Column:  startCol,
	}
}

func (s *Scanner) scanComment() token.Token {
	startLine, startCol := s.line, s.col
	s.next()

	if s.peek() == '/' {
		c := token.LeadingComment
		if !s.newLine {
			c = token.TrailingComment
		}
		s.next()
		var buf strings.Builder
		for {
			ch := s.peek()
			if ch == '\n' || ch == 0 {
				break
			}
			buf.WriteRune(s.next())
		}
		return token.Token{
			Type:    c,
			Literal: strings.TrimSpace(buf.String()),
			Line:    startLine,
			Column:  startCol,
		}
	}

	if s.peek() == '*' {
		c := token.LeadingComment
		if !s.newLine {
			c = token.TrailingComment
		}
		s.next()
		var buf strings.Builder
		for {
			if s.peek() == 0 {
				break
			}
			if s.peek() == '*' && s.pos+1 < len(s.src) && s.src[s.pos+1] == '/' {
				s.next()
				s.next()
				break
			}
			buf.WriteRune(s.next())
		}
		return token.Token{
			Type:    c,
			Literal: strings.TrimSpace(buf.String()),
			Line:    startLine,
			Column:  startCol,
		}
	}

	return token.Token{Type: token.EOF}
}

func (s *Scanner) scanLiteral() token.Token {
	startLine, startCol := s.line, s.col
	var buf strings.Builder

	for {
		ch := s.peek()
		if ch == 0 || strings.ContainsRune(" \t\r\n,:{}[]", ch) {
			break
		}
		buf.WriteRune(s.next())
	}

	lit := buf.String()
	switch lit {
	case ast.True:
		return token.Token{Type: token.True, Line: startLine, Column: startCol}
	case ast.False:
		return token.Token{Type: token.False, Line: startLine, Column: startCol}
	case ast.Null:
		return token.Token{Type: token.Null, Line: startLine, Column: startCol}
	default:
		return token.Token{
			Type:    token.Number,
			Literal: lit,
			Line:    startLine,
			Column:  startCol,
		}
	}
}
