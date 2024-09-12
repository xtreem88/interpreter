package scanner

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"

	COMMA     TokenType = "COMMA"
	DOT       TokenType = "DOT"
	MINUS     TokenType = "MINUS"
	PLUS      TokenType = "PLUS"
	SEMICOLON TokenType = "SEMICOLON"
	STAR      TokenType = "STAR"

	EOF TokenType = "EOF"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
}

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: source}
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{Type: EOF, Lexeme: "", Literal: nil})
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
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: nil})
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
