package main

import "fmt"

// ScanError define new error type for scan error
type ScanError string

func (e ScanError) Error() string {
	return string(e)
}

type Scanner struct {
	source  string      //input source code string
	tokens  []*tokenObj //parsed tokens
	start   int         //token lexeme start
	current int
	line    int
	err     error // if there is an error , stop scan
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: make([]*tokenObj, 0),
		line:   1,
	}
}

func (s *Scanner) scan() ([]*tokenObj, error) {
	//keep scan like a sliding window
	for !s.atEnd() && s.err == nil {
		s.start = s.current
		s.scanToken()
	}

	//put an EOF to indicate token end
	if s.err == nil {
		s.tokens = append(s.tokens, &tokenObj{tok: EOF, line: s.line})
	}
	return s.tokens, s.err
}

//scan single token
func (s *Scanner) scanToken() {
	ch := s.advance()
	switch ch {
	case '(':
		s.token(LeftParen)
	case ')':
		s.token(RightParen)
	case '{':
		s.token(LeftBrace)
	case ',':
		s.token(Comma)
	case ':':
		s.token(Colon)
	case '.':
		s.token(Dot)
	case '-':
		s.token(Minus)
	case '?':
		s.token(Question)
	case '+':
		s.token(Plus)
	case ';':
		s.token(Semicolon)
	case '*':
		s.token(Star)
	case '!':
		if s.match('=') {
			s.token(BangEqual)
		} else {
			s.token(Bang)
		}
	case '=':
		if s.match('=') {
			s.token(EqualEqual)
		} else {
			s.token(Equal)
		}
	case '<':
		if s.match('=') {
			s.token(LessEqual)
		} else {
			s.token(Less)
		}
	case '>':
		if s.match('=') {
			s.token(GreaterEqual)
		} else {
			s.token(Greater)
		}
	case '/':
		//  match "//............."
		if s.match('/') {
			for s.peek() != '\n' && !s.atEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.fullComment()
		} else {
			//match "/"
			s.token(Slash)
		}
	case ' ', '\r', '\t':
		//ignore
		break
	case '\n':
		s.line++
	case '"':
		s.stringLit()
	default:
		if isDigit(ch) {
			s.number()
		} else if isAlpha(ch) {
			s.identifier()
		} else {
			s.report(fmt.Sprintf("unexpected character '%c", ch))
		}
	}
}

func (s *Scanner) identifier() {
	for isAlphaNum(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	var t token
	//case , keyword
	if tok, ok := keywords[text]; ok {
		t = tok
	} else {
		t = Identifier
	}
	s.token(t)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) && !s.atEnd() {
		s.advance()
	}
	lit := s.source[s.start:s.current]
	s.literal(Number, lit)
}

//match "* */" , important cases , atEnd , \n , not terminated
func (s *Scanner) fullComment() {
	for !(s.peek() == '*' && s.peekNext() == '/') && !s.atEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.atEnd() {
		s.report("unterminated /**/ comment")
		return
	}
	s.advance()
	s.advance()
}

func (s *Scanner) token(t token) {
	s.literal(t, nil)
}

func (s *Scanner) literal(t token, literal interface{}) {
	lex := s.source[s.start:s.current]
	s.tokens = append(s.tokens, &tokenObj{
		tok:     t,
		lexeme:  lex,
		line:    s.line,
		literal: literal,
	})
}

//parse string
func (s *Scanner) stringLit() {
	for s.peek() != '"' && !s.atEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.atEnd() {
		s.report("unterminated string")
		return
	}
	s.advance()
	//add string token , literal without two "
	lit := s.source[s.start+1 : s.current-1]
	s.literal(String, lit)
}

//scanner helpers
func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return byte(0)
	}
	return s.source[s.current+1]
}

// look ahead , if match then consume
func (s *Scanner) match(ch byte) bool {
	//if !s.atEnd() && s.source[s.current] == ch {
	//	s.current++
	//	return true
	//}
	//return false
	if s.peek() != ch {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.atEnd() {
		return byte(0)
	}
	return s.source[s.current]
}
func (s *Scanner) atEnd() bool {
	return s.current >= len(s.source)
}

//retrieve an token an advance
func (s *Scanner) advance() byte {
	i := s.current
	s.current++
	return s.source[i]
}

func (s *Scanner) report(msg string) {
	s.err = ScanError(errorAt(s.line, "", msg))
}
