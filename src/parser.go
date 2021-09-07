package main

import "fmt"

// Recursive-descent parser
//
// program        -> declaration* EOF ;
//
// declaration    -> classDecl
//				   | funDecl
//                 | lambdaCall
//                 | varDecl
//                 | statement ;
//
// funDecl        -> "fun" function ;
// function       -> IDENTIFIER "(" parameters? ")" block ;
// parameters     -> IDENTIFIER ( "," IDENTIFIER )* ;
//
// lambdaCall     -> funExpr "(" arguments? ")" ";" ;
//
// varDecl        -> "var" IDENTIFIER ( "=" expression )? ";" ;
//
// statement      -> exprStmt
//                 | breakStmt
//                 | continueStmt
//                 | forStmt
//                 | ifStmt
//                 | printStmt
//                 | returnStmt
//                 | whileStmt
//				   | block ;
//
// block		  -> "{" declaration* "}" ;
// breakStmt      -> "break" ";" ;
// continueStmt   -> "continue" ";" ;
// exprStmt       -> expression ";" ;
// forStmt        -> "for" "(" ( varDecl | exprStmt | ";" )
//                   expression? ";"
//                   expression? ")" statement ;
// ifStmt         -> "if" "(" expression ")" statement ( "else" statement )? ;
// printStmt      -> "print" expression ";" ;
// returnStmt     -> "return" expression? ";" ;
// whileStmt      -> "while" "(" expression ")" statement ;
//
// expression     -> funExpr
//                 | assignment ;
// funExpr        -> "fun" "(" parameters? ")" block ;
//assignment     → ( call "." )? IDENTIFIER "=" assignment
//               | logic_or ;
// logicOr        -> logicAnd ( "or" logicAnd )* ;
// logicAnd       -> equality ( "and" equality )* ;
// equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           -> factor ( ( "-" | "+" ) factor )* ;
// factor         -> unary ( ( "/" | "*" ) unary )* ;
// unary          -> ( "!" | "-" ) unary | call ;
// call           → primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
// arguments      -> expression ( "," expression )* ;
// primary        -> NUMBER | STRING | "true" | "false" | "nil"
//                 | "(" expression ")"
//                 | IDENTIFIER ;
//

type parser struct {
	tokens  []*tokenObj
	current int
	errs    []error //multiple error support
	inLoop  int     // correctness check for statement which should be inside loop
}

func NewParser(tokens []*tokenObj) *parser {
	p := &parser{tokens, 0, make([]error, 0), 0}
	return p
}

// match advances pointer to the next token if current token matches
// any of toks and returns true
func (p *parser) match(toks ...token) bool {
	for _, t := range toks {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *parser) advance() *tokenObj {
	if !p.atEnd() {
		p.current++
	}
	return p.prev()
}

func (p *parser) atEnd() bool {
	return p.peek().tok == EOF
}

func (p *parser) peek() *tokenObj {
	return p.tokens[p.current]
}

func (p *parser) prev() *tokenObj {
	return p.tokens[p.current-1]
}

//check current token obj is equal to input token type
func (p *parser) check(tok token) bool {
	if p.atEnd() {
		return false
	}
	// fmt.Printf("p.peek() = %+v\n", p.peek())
	return p.peek().tok == tok
}

func (p *parser) consume(expected token, msg string) *tokenObj {
	if p.check(expected) {
		return p.advance()
	}
	p.primaryError(p.peek(), msg)
	return nil
}

//error type for parsing
type ParsingError string

func (e ParsingError) Error() string {
	return string(e)
}

//primary error that stop parse immediately and panic
func (p *parser) primaryError(t *tokenObj, msg string) {
	e := ParsingError(errorAtToken(t, msg))
	p.errs = append(p.errs, e)
	panic(e)
}

//errors that dont stop parse
func (p *parser) yerror(t *tokenObj, msg string) {
	e := ParsingError(errorAtToken(t, msg))
	p.errs = append(p.errs, e)
}

//https://craftinginterpreters.com/parsing-expressions.html#panic-mode-error-recovery
//https://craftinginterpreters.com/parsing-expressions.html#synchronizing-a-recursive-descent-parser
func (p *parser) sync() {
	fmt.Println("sync")
	p.advance()
	for !p.atEnd() {
		if p.prev().tok == Semicolon { //pass 这条解析错误的语句
			return
		}
		switch p.peek().tok { // or any of these start keyword
		case Class, Fun, Var, For, If, While, Print, Return:
			return
		}
		p.advance()
	}
}

// ---------------------------------------------------------
//

// parse returns an AST of parsed tokens, if it cannot parse then it returns
// the error.
func (p *parser) parse() (s []Stmt, errs []error) {
	s = make([]Stmt, 0)
	for !p.atEnd() {
		s = append(s, p.declaration())
	}

	return s, p.errs
}

func (p *parser) classDeclaration() Stmt {
	name := p.consume(Identifier, "Expect class name.")
	p.consume(LeftBrace, "Expect '{' before class body")

	methods := []*FunStmt{}
	for !p.check(RightBrace) && !p.atEnd() {
		methods = append(methods, p.funDecl("method").(*FunStmt))
	}
	p.consume(RightBrace, "Expect '}' after class body")
	return &ClassStmt{
		name:       name,
		methods:    methods,
		superClass: nil,
	}
}
