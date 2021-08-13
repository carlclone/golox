package main

type (
	Stmt interface {
		aStmt()
		execute(*Env)
	}

	stmt struct{}

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
)

func (*stmt) aStmt() {}

//func (*stmt) execute(*Env) {}
