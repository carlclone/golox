package main

import "fmt"

// --------------------------------------------------------
// Statements

func (s *ExprStmt) execute(env *Env) {
	s.expression.eval(env)
}

func (s *FunStmt) execute(env *Env) {
	fn := &FunObj{decl: s, closure: env, isInitializer: false}
	env.defineInit(s.name.lexeme, fn) // add fun decl as env variable
}

func (s *PrintStmt) execute(env *Env) {
	v := s.expression.eval(env)
	fmt.Printf("%v\n", v)
}

func (s *VarStmt) execute(env *Env) {
	// make distinction between uninitialized value and nil-value
	if s.init != nil {
		v := s.init.eval(env)
		env.defineInit(s.name.lexeme, v)
	} else {
		env.define(s.name.lexeme)
	}
}

func (s *BlockStmt) execute(env *Env) {
	execBlock(s.list, NewEnv(env))
}

func execBlock(list []Stmt, env *Env) {
	for _, s := range list {
		s.execute(env)
	}
}

func (s *ClassStmt) execute(env *Env) {
	env.define(s.name.lexeme)

	methods := make(map[string]*FunObj)
	for _, method := range s.methods {
		function := &FunObj{method, env, method.name.lexeme == "init"}
		methods[method.name.lexeme] = function
	}

	klass := &LoxClass{name: s.name.lexeme, methods: methods}
	env.assign(s.name, klass)
}

type LoxClass struct {
	name    string
	methods map[string]*FunObj
}

func (l *LoxClass) String() string {
	return l.name
}

func (l *LoxClass) findMethod(name string) *FunObj {
	if method, ok := l.methods[name]; ok {
		return method
	}
	return nil
}

func (l *LoxClass) call(e *Env, arg []value) value {
	instance := &LoxInstance{
		klass:  l,
		fields: make(map[string]value),
	}
	initializer := l.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(e, arg)
	}
	return instance
}

func (l *LoxClass) arity() int {
	/*
		LoxFunction initializer = findMethod("init");
		    if (initializer == null) return 0;
		    return initializer.arity();
	*/
	initializer := l.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

type LoxInstance struct {
	klass  *LoxClass
	fields map[string]value
}

func (l *LoxInstance) get(name *tokenObj) value {
	v, ok := l.fields[name.lexeme]
	if ok {
		return v
	}

	method := l.klass.findMethod(name.lexeme)
	if method != nil {
		return method.bind(l)
	}
	runtimeErr(name, "Undefined property '"+name.lexeme+"'.")
	return nil
}

func (l *LoxInstance) set(name *tokenObj, v value) {
	l.fields[name.lexeme] = v
}

func (l *LoxInstance) String() string {
	return l.klass.name + " instance"
}

func (s *IfStmt) execute(env *Env) {
	if isTruthy(s.condition.eval(env)) {
		s.block1.execute(env)
	} else if s.block2 != nil {
		s.block2.execute(env)
	}
}

func (s *ReturnStmt) execute(env *Env) {
	var v value
	if s.value != nil {
		v = s.value.eval(env)
	}
	// Ugly hack, panic to unwind the stack back to the call
	panic(ReturnHack(v))
}

func (s *BreakStmt) execute(env *Env) {
	panic(BreakErr{t: s.keyword})
}

func (s *ContinueStmt) execute(env *Env) {
	panic(ContinueErr{t: s.keyword})
}

func (s *WhileStmt) execute(env *Env) {
	for !s.isDone(env) { //todo
	}
}

// isDone returns false when the loop was continued,
// when loop is done returns true
func (s *WhileStmt) isDone(env *Env) (done bool) { //todo
	defer func() {
		if e := recover(); e != nil {
			switch e.(type) { //  error type case
			case ContinueErr:
				done = false
				return
			case BreakErr:
				done = true
				return
			default:
				panic(e)
			}
		}
	}()
	for isTruthy(s.condition.eval(env)) {
		s.body.execute(env)
	}
	return true
}
