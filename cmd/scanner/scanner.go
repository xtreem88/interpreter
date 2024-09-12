package scanner

import (
	"fmt"
	"os"
)

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	STAR        TokenType = "STAR"
	EOF         TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	hadError bool
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source, line: 1}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{Type: EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	default:
		s.error(fmt.Sprintf("Unexpected character: %c", c))
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: nil, Line: s.line})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) error(message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", s.line, message)
	s.hadError = true
}

func (s *Scanner) HadError() bool {
	return s.hadError
}
