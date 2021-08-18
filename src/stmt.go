package main

type (
	Stmt interface {
		aStmt()
		execute(*Env)
		accept(*Resolver)
	}

	stmt struct {
		id int
	}

	BlockStmt struct {
		list []Stmt
		stmt
	}

	BreakStmt struct {
		keyword *tokenObj
		stmt
	}

	ContinueStmt struct {
		keyword *tokenObj
		stmt
	}

	ExprStmt struct {
		expression Expr
		stmt
	}

	FunStmt struct { //declaration of a function
		name   *tokenObj
		params []*tokenObj
		body   []Stmt
		stmt   // something like extend
	}

	// if structure
	IfStmt struct {
		condition      Expr
		block1, block2 Stmt
		stmt
	}

	PrintStmt struct {
		expression Expr
		stmt
	}

	ReturnStmt struct {
		keyword *tokenObj
		value   Expr
		stmt
	}

	VarStmt struct {
		name *tokenObj
		init Expr
		stmt
	}

	WhileStmt struct {
		condition Expr
		body      Stmt
		stmt
	}
	ClassStmt struct {
		name       *tokenObj
		methods    []Stmt
		superClass *VarExpr //todo

		stmt
	}
)

func (*stmt) aStmt()                    {}
func (*stmt) accept(resolver *Resolver) {}
func (*stmt) execute(*Env)              {}

func (s *VarStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitVarStmt(s)
}

func (s *WhileStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitWhileStmt(s)
}
func (s *ReturnStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitReturnStmt(s)
}

func (s *PrintStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitPrintStmt(s)
}
func (s *IfStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitIfStmt(s)
}

func (s *FunStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitFunctionStmt(s)
}
func (s *ExprStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitExpressionStmt(s)
}
func (s *ClassStmt) accept(r *Resolver) {
	r.visitClassStmt(s)
}
func (s *ContinueStmt) accept(r *Resolver) {
	s.id = GetId()
	//todo
	//r.visitC(s)
}
func (s *BreakStmt) accept(r *Resolver) {
	s.id = GetId()
	//todo
	//r.visibre(s)
}
func (s *BlockStmt) accept(r *Resolver) {
	s.id = GetId()
	r.visitBlockStmt(s)
}
