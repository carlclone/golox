package main

type (
	value interface{} // alias for readability

	Expr interface { // interface declaration for type check
		//aExpr()
		eval(*Env) value
		accept(resolver *Resolver)
	}

	expr struct {
		id int
	} // think as a parent class

	GetExpr struct {
		name   *tokenObj
		object Expr

		expr
	}

	SetExpr struct {
		name   *tokenObj
		object Expr
		vlue   Expr

		expr
	}

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
		expression Expr //mostly expression surround by ()
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

//func (*expr) aExpr()                    {} //add a empty method to distinct with other interface that has eval
func (*expr) eval(*Env) value           { return nil }
func (*expr) accept(resolver *Resolver) {}

func (s *VarExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitVariableExpr(s)
}

func (s *UnaryExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitUnaryExpr(s)
}
func (s *LogicalExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitLogicalExpr(s)
}
func (s *LiteralExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitLiteralExpr(s)
}
func (s *GroupingExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitGroupingExpr(s)
}
func (s *FunExpr) accept(r *Resolver) {
	s.id = GetId()

	//r.visitF(s)
}
func (s *CallExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitCallExpr(s)
}
func (s *BinaryExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitBinaryExpr(s)
}
func (s *AssignExpr) accept(r *Resolver) {
	s.id = GetId()

	r.visitAssignExpr(s)
}
