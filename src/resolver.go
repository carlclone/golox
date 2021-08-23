package main

import (
	"fmt"
)

type FunctionType uint
type ClassType uint

const (
	FT_NONE = iota
	FT_FUNCTION
	FT_INITIALIZER
	FT_METHOD
)

const (
	CT_NONE = iota
	CT_CLASS
	CT_SUBCLASS
)

func NewResolver() *Resolver {
	return &Resolver{
		scopes:          make([]map[string]bool, 0),
		currentFunction: 0,
		currentClass:    0,
	}
}

type Resolver struct {
	scopes          []map[string]bool
	currentFunction FunctionType
	currentClass    ClassType
}

func (r *Resolver) resolve(stmts []Stmt) {
	for _, stmt := range stmts {
		r.resolveStmt(stmt)
	}
}

func (r *Resolver) resolveStmt(stmt Stmt) {
	stmt.accept(r)
}

func (r *Resolver) visitBlockStmt(s *BlockStmt) {
	r.beginScope()
	r.resolve(s.list)
	r.endScope()
	return
}

//TODO
func (r *Resolver) visitClassStmt(s *ClassStmt) {
	r.declare(s.name)

	for _, method := range s.methods {
		r.resolveFunction(method.(*FunStmt), FT_METHOD)
	}
	r.define(s.name)
	return
}

func (r *Resolver) visitExpressionStmt(s *ExprStmt) {
	r.resolveExpr(s.expression)
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) resolveExpr(expression Expr) {
	expression.accept(r)
}

func (r *Resolver) visitFunctionStmt(s *FunStmt) {
	// separate for edge case
	/*
		var a = 0
		{
			var a = a
		}
	*/
	r.declare(s.name)
	r.define(s.name)

	r.resolveFunction(s, FT_FUNCTION)
}

func (r *Resolver) visitIfStmt(s *IfStmt) {
	r.resolveExpr(s.condition)
	r.resolveStmt(s.block1)
	if s.block2 != nil {
		r.resolveStmt(s.block2)
	}
}

func (r *Resolver) visitPrintStmt(s *PrintStmt) {
	r.resolveExpr(s.expression)
}

func (r *Resolver) visitReturnStmt(s *ReturnStmt) {
	//TODO
	if r.currentFunction == FT_NONE {
		errorAtToken(s.keyword, "Can't return from top-level code.")
	}

	//TODO
	if s.value != nil {
		if r.currentFunction == FT_INITIALIZER {
			errorAtToken(s.keyword, "Can't return a value from an initializer.")
		}

		r.resolveExpr(s.value)
	}
	return
}

func (r *Resolver) visitVarStmt(s *VarStmt) {
	r.declare(s.name)
	if s.init != nil {
		r.resolveExpr(s.init)
	}
	r.define(s.name)
	return
}

func (r *Resolver) visitWhileStmt(s *WhileStmt) {
	r.resolveExpr(s.condition)
	r.resolveStmt(s.body) //block stmt
	return
}

func (r *Resolver) visitAssignExpr(e *AssignExpr) {
	r.resolveExpr(e.value)
	r.resolveLocal(e, e.name)
	return
}

func (r *Resolver) visitBinaryExpr(e *BinaryExpr) {
	r.resolveExpr(e.left)
	r.resolveExpr(e.right)
	return
}

func (r *Resolver) visitCallExpr(e *CallExpr) {
	r.resolveExpr(e.callee)
	for _, argument := range e.args {
		r.resolveExpr(argument)
	}
	return
}

func (r *Resolver) visitGetExpr(e *GetExpr) {
	r.resolveExpr(e.object)
}
func (r *Resolver) visitSetExpr(e *SetExpr) {
	r.resolveExpr(e.object)
	r.resolveExpr(e.vlue)
}

func (r *Resolver) visitGroupingExpr(e *GroupingExpr) {
	r.resolveExpr(e.expression)
	return
}

func (r *Resolver) visitLiteralExpr(e *LiteralExpr) {
	return
}

func (r *Resolver) visitLogicalExpr(e *LogicalExpr) {
	r.resolveExpr(e.left)
	r.resolveExpr(e.right)
	return
}

//todo
//func (r *Resolver) visitSetExpr(e *LiteralExpr) {
//	r.resolveExpr(e.value)
//	r.resolveExpr(e.object)
//	return
//}

//todo
func (r *Resolver) visitSuperExpr(e *LiteralExpr) {
	return
}

func (r *Resolver) visitThisExpr(e *LiteralExpr) {
	return
}
func (r *Resolver) visitUnaryExpr(e *UnaryExpr) {
	r.resolveExpr(e.right)
	return
}
func (r *Resolver) visitVariableExpr(e *VarExpr) {
	if len(r.scopes) != 0 && r.scopePeek()[e.name.lexeme] == false {
		errorAtToken(e.name, "Can't read local variable in its own initializer.")
	}
	r.resolveLocal(e, e.name)
	return
}

func (r *Resolver) declare(name *tokenObj) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	_, ok := scope[name.lexeme]
	if ok {
		errorAtToken(name, "Already a variable with this name in this scope.")
	}
	scope[name.lexeme] = false
}

func (r *Resolver) define(name *tokenObj) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.lexeme] = true
}

func (r *Resolver) resolveFunction(s *FunStmt, typee FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typee

	r.beginScope()

	//resolve params
	for _, tok := range s.params {
		r.declare(tok)
		r.define(tok)
	}
	//resolve body
	r.resolve(s.body)

	r.endScope()

	r.currentFunction = enclosingFunction
}

//TODO
func (r *Resolver) resolveLocal(id Expr, name *tokenObj) {
	fmt.Println("resolveLocal called")
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		if containKey(scope, name.lexeme) {
			locals.put(id, len(r.scopes)-1-i)
			return
		}
	}

}

func (r *Resolver) scopePeek() map[string]bool {
	return r.scopes[len(r.scopes)-1]
}

func containKey(m map[string]bool, key string) bool {
	_, ok := m[key]
	return ok
}

func (s *LiteralExpr) reslove(r *Resolver) {
	r.visitLiteralExpr(s)
}
