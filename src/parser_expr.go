package main

// expression     -> funExpr
//                 | assignment ;
// funExpr        -> "fun" "(" parameters? ")" block ;
// assignment     -> IDENTIFIER "=" assignment
//				   | logicOr ;
// logicOr        -> logicAnd ( "or" logicAnd )* ;
// logicAnd       -> equality ( "and" equality )* ;
// equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           -> factor ( ( "-" | "+" ) factor )* ;
// factor         -> unary ( ( "/" | "*" ) unary )* ;
// unary          -> ( "!" | "-" ) unary | call ;
// call			  -> primary ( "(" arguments? ")" )* ;
// arguments      -> expression ( "," expression )* ;
// primary        -> NUMBER | STRING | "true" | "false" | "nil"
//                 | "(" expression ")"
//                 | IDENTIFIER ;

//priority related design , BNF method
func (p *parser) assignment() Expr {
	expr := p.or()
	if p.match(Equal) {
		equals := p.prev()
		value := p.assignment()
		if ev, ok := expr.(*VarExpr); ok {
			name := ev.name
			return &AssignExpr{name: name, value: value}
		} else if get, ok := expr.(*GetExpr); ok {
			return &SetExpr{
				name:   get.name,
				object: get.object,
				vlue:   value,
			}
		}
		p.yerror(equals, "invalid assignment target")
	}
	return expr
}

func (p *parser) or() Expr {
	expr := p.and() //也是优先级的体现
	for p.match(Or) {
		op := p.prev()
		right := p.and() //优先级的体现 ,  a&&b || c&&d , avoid exec b||c
		expr = &LogicalExpr{operator: op, left: expr, right: right}
	}
	return expr
}

func (p *parser) and() Expr { //a==b && c==d , avoid b&&c
	expr := p.equality()
	for p.match(And) {
		op := p.prev()
		right := p.equality()
		expr = &LogicalExpr{operator: op, left: expr, right: right}
	}
	return expr
}

// equality -> comparison ( ( "!=" | "==" ) comparison )* ;
func (p *parser) equality() Expr { // a>=c == b<=d , avoid c==b
	expr := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		op := p.prev()
		right := p.comparison()
		expr = &BinaryExpr{operator: op, left: expr, right: right}
	}
	return expr
}

// comparison -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *parser) comparison() Expr { // 1+2 >= 3-4
	expr := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		op := p.prev()
		right := p.term()
		expr = &BinaryExpr{operator: op, left: expr, right: right}
	}
	return expr
}

// term ->  factor ( ( "-" | "+" ) factor )* ;
//what is term 这里的term是项的意思，是指在加减法里运算符左右的两个项
func (p *parser) term() Expr { //(1+2) +(3+4)
	expr := p.factor()
	for p.match(Plus, Minus) {
		op := p.prev()
		right := p.factor()
		expr = &BinaryExpr{operator: op, left: expr, right: right}
	}
	return expr
}

// factor -> unary ( ( "/" | "*" ) unary )* ;
//括号的运算优先级是比乘除法还要高的，所以我们新增一个非终结符factor（因子）
func (p *parser) factor() Expr {
	expr := p.unary() //  -1 * -2
	for p.match(Slash, Star) {
		op := p.prev()
		right := p.unary()
		expr = &BinaryExpr{operator: op, left: expr, right: right}
	}
	return expr
}

// unary -> ( "!" | "-" ) unary
//        | primary ;
func (p *parser) unary() Expr {
	if p.match(Bang, Minus) { //--2  !!true
		op := p.prev()
		right := p.unary()
		return &UnaryExpr{operator: op, right: right}
	}
	return p.call()
}

// call			  -> primary ( "(" arguments? ")" )* ;
// arguments      -> expression ( "," expression )* ;
// primary        -> NUMBER | STRING | "true" | "false" | "nil"
//                 | "(" expression ")"
//                 | IDENTIFIER ;
func (p *parser) call() Expr {
	expr := p.primary() //主表达式 , 理解为一个值 , 或者产生值的 主体
	for {
		if p.match(LeftParen) {
			expr = p.finishCall(expr)
		} else if p.match(Dot) {
			name := p.consume(Identifier, "Expect property name after '.'.")
			expr = &GetExpr{
				name:   name,
				object: expr,
			}
		} else {
			break
		}
	}
	return expr
}

//immediately call the funExpr produce
func (p *parser) finishCall(expr Expr) Expr {
	args := make([]Expr, 0)
	if !p.check(RightParen) {
		// parse params
		for {
			if len(args) >= 255 {
				p.yerror(p.peek(), "can't have more than 255 arguments")
			}
			args = append(args, p.expression())
			if !p.match(Comma) {
				break
			}
		}
	}
	// for error display
	paren := p.consume(RightParen, "expected ')' after arguments")

	//callee : 被 call 的人 , 先被 call 产生值作为参数
	return &CallExpr{callee: expr, paren: paren, args: args}
}

// primary -> NUMBER | STRING | "true" | "false" | "nil"
//          | "(" expression ")" ;
func (p *parser) primary() Expr {
	switch {
	case p.match(False):
		return &LiteralExpr{value: false}
	case p.match(True):
		return &LiteralExpr{value: true}
	case p.match(Nil):
		return &LiteralExpr{value: nil}
	case p.match(Number, String):
		return &LiteralExpr{value: p.prev().literal}
	case p.match(Identifier):
		return &VarExpr{name: p.prev()}
	case p.match(LeftParen):
		expr := p.expression()
		p.consume(RightParen, "expected enclosing ')' after expression")
		return &GroupingExpr{expression: expr}
	}
	p.primaryError(p.peek(), "expected expression")
	return nil
}

// parse expr part

// expression     -> funExpr
//                 | assignment ;
// funExpr        -> "fun" "(" parameters? ")" block ;
// assignment     -> IDENTIFIER "=" assignment
//				   | logicOr ;
// logicOr        -> logicAnd ( "or" logicAnd )* ;
// logicAnd       -> equality ( "and" equality )* ;
// equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
// comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
// term           -> factor ( ( "-" | "+" ) factor )* ;
// factor         -> unary ( ( "/" | "*" ) unary )* ;
// unary          -> ( "!" | "-" ) unary | call ;
// call			  -> primary ( "(" arguments? ")" )* ;
// arguments      -> expression ( "," expression )* ;
// primary        -> NUMBER | STRING | "true" | "false" | "nil"
//                 | "(" expression ")"
//                 | IDENTIFIER ;
func (p *parser) expression() Expr {
	if p.match(Fun) {
		return p.funExpr()
	}
	return p.assignment()
}

//// funExpr        -> "fun" "(" parameters? ")" block ;
// is also an expression , it produce what ? a closure
//TODO; code just a tool for implement logic , what most important is logic
func (p *parser) funExpr() Expr {
	p.consume(LeftParen, "expected '(' after 'fun'")
	params := make([]*tokenObj, 0)
	if !p.check(RightParen) {
		// parse param name
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
	p.consume(LeftBrace, "expected '{' after anonymous function signature")
	// parse block
	body := p.block()
	return &FunExpr{params: params, body: body}
}
