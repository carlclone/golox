package main

import "fmt"

type token uint

//go:generate stringer -type token -linecomment tokens.go

// preserved keyword token
var keywords = map[string]token{
	"and":      And,
	"break":    Break,
	"class":    Class,
	"continue": Continue,
	"else":     Else,
	"false":    False,
	"for":      For,
	"fun":      Fun,
	"if":       If,
	"nil":      Nil,
	"or":       Or,
	"print":    Print,
	"return":   Return,
	"super":    Super,
	"this":     This,
	"true":     True,
	"var":      Var,
	"while":    While,
}

const (
	_ token = iota
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Colon
	Question
	Slash
	Star

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Number

	And
	Break
	Class
	Continue
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True  //true
	Var   //var
	While //while

	EOF //eof

)

type tokenObj struct {
	tok     token
	lexeme  string
	line    int
	literal interface{}
}

func (t *tokenObj) String() string {
	return fmt.Sprintf("token: %v lex: %v lit: %v",
		t.tok, t.lexeme, t.literal)
}
