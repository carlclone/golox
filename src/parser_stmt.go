package main

// parse single stmt ast-tree from tokens
func (p *parser) declaration() (s Stmt) {
	defer func() {
		if e := recover(); e != nil {
			_ = e.(ParsingError) // Panic for other errors
			/*
				var a interface{}
				a=1
				v,ok:=a.(string)  //显示使用了第二个返回值 , 不 panic
				_=a.(string)  // 没使用 , panic

			*/
			p.sync()
			s = nil
		}
	}()

	// Recursive-descent parse
	// program        -> declaration* EOF ;
	//
	// declaration    -> funDecl
	//                 | lambdaCall
	//                 | varDecl
	//                 | statement ;
	if p.match(Class) {
		return p.classDeclaration()
	}
	if p.match(Fun) {
		if p.check(LeftParen) {
			return p.lambdaCall()
		}
		return p.funDecl("function")
	}
	if p.match(Var) {
		return p.varDecl()
	}
	return p.statement()
}

func (p *parser) funDecl(kind string) Stmt {
	name := p.consume(Identifier, "expected "+kind+" name")
	p.consume(LeftParen, "expected '(' after "+kind+" name")

	// parse params
	params := make([]*tokenObj, 0)
	if !p.check(RightParen) {
		for {
			if len(params) >= 255 {
				p.yerror(p.peek(), "can't have more than 255 parameters")
			}
			params = append(params, p.consume(Identifier, "expected parameter name"))
			if !p.match(Comma) {
				break
			}
		}
	}
	p.consume(RightParen, "expected ')' after parameters")
	p.consume(LeftBrace, "expected '{' after "+kind+" signature")

	body := p.block()
	return &FunStmt{name: name, params: params, body: body}
}

// var a  , var a=1
func (p *parser) varDecl() Stmt {
	name := p.consume(Identifier, "expected variable name")
	var init Expr

	if p.match(Equal) {
		init = p.expression()
	}
	p.consume(Semicolon, "expected ';' after variable declaration")
	return &VarStmt{name: name, init: init}
}

/* difference between statement and expr ?

almost the same , part of ast-tree , separate for convenient

*/

// statement      -> exprStmt
//                 | breakStmt
//                 | continueStmt
//                 | forStmt
//                 | ifStmt
//                 | printStmt
//                 | returnStmt
//                 | whileStmt
//				   | block ;
func (p *parser) statement() Stmt {
	if p.match(Break) {
		return p.breakStatement()
	}
	if p.match(Continue) {
		return p.continueStatement()
	}
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(Return) {
		return p.returnStatement()
	}
	if p.match(While) {
		return p.whileStatement()
	}
	if p.match(LeftBrace) {
		return &BlockStmt{list: p.block()}
	}
	return p.exprStatement()
}

func (p *parser) breakStatement() Stmt {
	key := p.prev()
	// sanity check
	if p.inLoop < 1 {
		p.primaryError(key, "expected inside the loop")
	}
	p.consume(Semicolon, "expected ';' after break")
	return &BreakStmt{keyword: key}
}

func (p *parser) continueStatement() Stmt {
	key := p.prev()
	if p.inLoop < 1 {
		p.primaryError(key, "expected inside the loop")
	}
	p.consume(Semicolon, "expected ';' after continue")
	return &ContinueStmt{keyword: key}
}

func (p *parser) forStatement() Stmt {
	p.consume(LeftParen, "expected '(' after 'for'")

	var initial Stmt
	switch {
	case p.match(Semicolon):
		initial = nil
	case p.match(Var):
		initial = p.varDecl()
	default:
		initial = p.exprStatement()
	}

	var cond Expr
	if !p.check(Semicolon) {
		cond = p.expression()
	}
	p.consume(Semicolon, "expected ';' after for condition")

	var incr Expr
	if !p.check(RightParen) {
		incr = p.expression()
	}
	p.consume(RightParen, "expected ')' after for clauses")

	p.inLoop += 1
	body := p.statement() //may be a block statement or other one line code
	p.inLoop -= 1

	// all of these three may not exist , nested set since  , TODO;beautiful design
	if incr != nil {
		body = &BlockStmt{list: []Stmt{
			body,
			&ExprStmt{expression: incr}}}
	}
	if cond != nil {
		body = &WhileStmt{condition: cond, body: body}
	}
	if initial != nil {
		body = &BlockStmt{list: []Stmt{
			initial,
			body}}
	}
	return body
}

//TODO;readable code standard ,  code just a tool for implement logic , what most important is logic
func (p *parser) ifStatement() Stmt {
	p.consume(LeftParen, "expected '(' after 'if'")
	//must be an expression here , becasue stmt dont produce
	e := p.expression()
	p.consume(RightParen, "expected ')' after if condition")
	a := p.statement()
	var b Stmt = nil

	//else
	if p.match(Else) {
		b = p.statement()
	}
	return &IfStmt{condition: e, block1: a, block2: b}
}

func (p *parser) printStatement() Stmt {
	// print stmt follow with expression
	e := p.expression()
	p.consume(Semicolon, "expected ';' after expression")
	return &PrintStmt{expression: e}
}

func (p *parser) returnStatement() Stmt {
	k := p.prev()
	var val Expr
	if !p.check(Semicolon) { //case : return;
		val = p.expression()
	}
	p.consume(Semicolon, "expected ';' after return value")
	return &ReturnStmt{keyword: k, value: val}
}

func (p *parser) whileStatement() Stmt {
	p.consume(LeftParen, "expected '(' after while")
	expr := p.expression()
	p.consume(RightParen, "expected ')' after while condition")
	p.inLoop += 1
	body := p.statement()
	p.inLoop -= 1
	return &WhileStmt{condition: expr, body: body}
}

// funDecl        -> "fun" function ;
// function       -> IDENTIFIER "(" parameters? ")" block ;
// block		  -> "{" declaration* "}" ;
func (p *parser) block() []Stmt {
	list := make([]Stmt, 0)
	for !p.check(RightBrace) && !p.atEnd() {
		list = append(list, p.declaration())
	}
	p.consume(RightBrace, "expected '}' after block")
	return list
}

func (p *parser) exprStatement() Stmt {
	e := p.expression()
	p.consume(Semicolon, "expected ';' after expression")
	return &ExprStmt{expression: e}
}

// lambdaCall     -> funExpr "(" arguments? ")" ";" ;      // decl and call
// funExpr        -> "fun" "(" parameters? ")" block ;    //fun decl
func (p *parser) lambdaCall() Stmt {
	expr := p.funExpr()
	for {
		if p.match(LeftParen) {
			expr = p.finishCall(expr) // case :  func() { print 1 }()
		} else {
			break
		}
	}
	p.consume(Semicolon, "expected ';' call to a function")
	return &ExprStmt{expression: expr}
}
