package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/islml/neptune/token"
)

type Scanner struct {
	Source []rune
	Tokens []token.Token

	start int
	current int
	line int
}

func (s *Scanner) ScanTokens() []token.Token {
	s.start   = 0
	s.current = 0
	s.line    = 1
	
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, token.Token{
		Type: token.TokenType_EOF,
		Lexeme: "",
		Literal: nil,
		Line: s.line,
	})
	return s.Tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.Source)
}

func (s *Scanner) scanToken() {
	var c rune = s.advance()

	switch c {
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
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		s.sstring()
	default:
		if unicode.IsDigit(c) {
			s.number()
		} else if unicode.IsLetter(c) {
			s.identifier()
		} else {
			panic(fmt.Sprintf("Unexpected character at line %d", s.line))
		}
	}
}

func (s *Scanner) advance() rune {
	temp := s.Source[s.current]
	s.current++
	return temp
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	str := s.Source[s.start : s.current]
	s.Tokens = append(s.Tokens, token.Token{
		Type: tokenType,
		Lexeme: string(str),
		Literal: nil,
		Line: s.line,
	})
}

func (s *Scanner) addTokenWithLiteral(tokenType token.TokenType, literal any) {
	s.Tokens = append(s.Tokens, token.Token{
		Type: tokenType,
		Lexeme: fmt.Sprintf("%v", literal),
		Literal: literal,
		Line: s.line,
	})
}

func (s *Scanner) match(expected rune) bool {
	if (s.isAtEnd()) {
		return false
	}
	if (s.Source[s.current] != expected) {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\n'
	} else {
		return s.Source[s.current]
	}
}

func (s *Scanner) nextPeek() rune {
	if s.current + 1 >= len(s.Source) {
		return 0
	}
  
	return s.Source[s.current + 1];
}

func (s *Scanner) sstring() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		panic(fmt.Sprintf("Unterminated string at line %d", s.line))
	}

	s.advance()
	value := s.Source[s.start + 1 : s.current - 1]
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

	num, _ := strconv.ParseFloat(string(s.Source[s.start : s.current]), 64)
	s.addTokenWithLiteral(token.TokenType_Number, num)
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if t, ok := keywords[string(s.Source[s.start:s.current])]; ok {
		s.addToken(t)
		return
	}
	s.addToken(token.TokenType_Identifier)
}

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