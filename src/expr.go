package main

type (
	value interface{} // alias for readability

	Expr interface { // interface declaration for type check
		aExpr()
		//eval(*Env) value
	}

	expr struct{} // think as a parent class

	AssignExpr struct {
		name  *tokenObj
		value Expr
		expr  //extend parent class and its method
	}

	BinaryExpr struct {
		operator    *tokenObj
		left, right Expr
		expr
	}

	CallExpr struct { // asd(1,2) , no env related struct
		callee Expr
		paren  *tokenObj
		args   []Expr
		expr
	}

	FunExpr struct { //fun decl , difference between FunAnon , FunObj , no env related
		params []*tokenObj
		body   []Stmt
		expr
	}

	GroupingExpr struct {
		e Expr //mostly expression surround by ()
		// design for general method , like extend in oop
		expr
	}

	LiteralExpr struct {
		value interface{}
		expr
	}

	LogicalExpr struct {
		operator    *tokenObj
		left, right Expr
		expr
	}

	UnaryExpr struct {
		operator *tokenObj
		right    Expr
		expr
	}

	VarExpr struct {
		name *tokenObj
		expr
	}
)

func (*expr) aExpr() {} //add a empty method to distinct with other interface that has eval
//func (*expr) eval(*Env) value { return nil }
