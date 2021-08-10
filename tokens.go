package main

import "fmt"

type token uint

//go:generate stringer -type token -linecomment tokens.go

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
