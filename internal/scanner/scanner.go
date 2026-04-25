package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/islml/neptune/internal/token"
)

var keywords = map[string]token.TokenType{
	"and":    token.TokenType_And,
	"class":  token.TokenType_Class,
	"else":   token.TokenType_Else,
	"false":  token.TokenType_False,
	"for":    token.TokenType_For,
	"fun":    token.TokenType_Fun,
	"if":     token.TokenType_If,
	"nil":    token.TokenType_Nil,
	"or":     token.TokenType_Or,
	"print":  token.TokenType_Print,
	"return": token.TokenType_Return,
	"super":  token.TokenType_Super,
	"this":   token.TokenType_This,
	"true":   token.TokenType_True,
	"var":    token.TokenType_Var,
	"while":  token.TokenType_While,
}

type ScanError struct {
	Line    int
	Message string
}

type Scanner struct {
	Source []rune
	Tokens []token.Token
	Errors []ScanError

	start   int
	current int
	line    int
}

func New(source string) *Scanner {
	return &Scanner{Source: []rune(source)}
}

func (s *Scanner) ScanTokens() ([]token.Token, []ScanError) {
	s.start = 0
	s.current = 0
	s.line = 1

	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, token.Token{
		Type:    token.TokenType_EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})

	return s.Tokens, s.Errors
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	// Single-character tokens.
	case '(':
		s.addToken(token.TokenType_LeftParen)
	case ')':
		s.addToken(token.TokenType_RightParen)
	case '{':
		s.addToken(token.TokenType_LeftBrace)
	case '}':
		s.addToken(token.TokenType_RightBrace)
	case ',':
		s.addToken(token.TokenType_Comma)
	case '.':
		s.addToken(token.TokenType_Dot)
	case '-':
		s.addToken(token.TokenType_Minus)
	case '+':
		s.addToken(token.TokenType_Plus)
	case ';':
		s.addToken(token.TokenType_Semicolon)
	case '*':
		s.addToken(token.TokenType_Star)
	case '!':
		if s.match('=') {
			s.addToken(token.TokenType_BangEqual)
		} else {
			s.addToken(token.TokenType_Bang)
		}

	// One or two character tokens.
	case '=':
		if s.match('=') {
			s.addToken(token.TokenType_EqualEqual)
		} else {
			s.addToken(token.TokenType_Equal)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.TokenType_LessEqual)
		} else {
			s.addToken(token.TokenType_Less)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.TokenType_GreaterEqual)
		} else {
			s.addToken(token.TokenType_Greater)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.TokenType_Slash)
		}

	// Ignored characters.
	case ' ', '\r', '\t':
		return
	case '\n':
		s.line++

	// Strings.
	case '"':
		s.stringLiteral()

	// Numbers & Identifiers.
	default:
		if unicode.IsDigit(c) {
			s.number()
		} else if unicode.IsLetter(c) || c == '_' {
			s.identifier()
		} else {
			s.addError(s.line, fmt.Sprintf("Unexpected character %q.", c))
		}
	}
}

func (s *Scanner) advance() rune {
	temp := s.Source[s.current]
	s.current++
	return temp
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	str := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, token.Token{
		Type:    tokenType,
		Lexeme:  string(str),
		Literal: nil,
		Line:    s.line,
	})
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal any) {
	s.Tokens = append(s.Tokens, token.Token{
		Type:    tokenType,
		Lexeme:  string(s.Source[s.start:s.current]),
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) addError(line int, message string) {
	s.Errors = append(s.Errors, ScanError{Line: line, Message: message})
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.Source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.Source[s.current]
}

func (s *Scanner) nextPeek() rune {
	if s.current+1 >= len(s.Source) {
		return 0
	}

	return s.Source[s.current+1]
}

func (s *Scanner) stringLiteral() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.addError(s.line, "Unterminated string.")
		return
	}

	s.advance()
	value := string(s.Source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(token.TokenType_String, value)
}

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.nextPeek()) {
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(string(s.Source[s.start:s.current]), 64)
	if err != nil {
		s.addError(s.line, fmt.Sprintf("Invalid number %q.", string(s.Source[s.start:s.current])))
		return
	}

	s.addTokenWithLiteral(token.TokenType_Number, num)
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) || s.peek() == '_' {
		s.advance()
	}

	if t, ok := keywords[string(s.Source[s.start:s.current])]; ok {
		s.addToken(t)
		return
	}

	s.addToken(token.TokenType_Identifier)
}
