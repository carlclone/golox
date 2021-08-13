package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type RuntimeError string

func (e RuntimeError) Error() string {
	return string(e)
}

//helper
func runtimeErr(t *tokenObj, msg string) error {
	panic(RuntimeError(
		fmt.Sprintf("[line %v] runtime error: %v", t.line, msg)))
}

type ReturnHack value //TODO; what return hack
type BreakErr struct{ t *tokenObj }
type ContinueErr struct{ t *tokenObj }

type Callable interface {
	arity() int //arity 参数个数
	call(*Env, []value) value
}

// ------------------------------------------
// env

// env contains bindings for variables.
type Env struct {
	values map[string]value // decl but not inited

	// init means that variable was properly initialized
	init map[string]bool

	enclosing *Env // most close upper env
	globals   *Env // always points to the root of enclosures
}

func NewEnv(enclosing *Env) *Env {
	e := &Env{make(map[string]value), make(map[string]bool), enclosing, nil}
	if enclosing == nil {
		// means that this created env is the root, that is global env
		e.globals = e
	} else {
		e.globals = enclosing.globals
	}
	return e
}

func (e *Env) defineInit(name string, v value) {
	e.values[name] = v
	e.init[name] = true
}

func (e *Env) define(name string) {
	e.values[name] = nil
}

func (e *Env) get(name *tokenObj) value {
	if v, ok := e.values[name.lexeme]; ok {
		if _, ok := e.init[name.lexeme]; !ok {
			runtimeErr(name, "variable '"+name.lexeme+"' should be initialized first")
		}
		return v
	}
	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	runtimeErr(name, "undefined variable '"+name.lexeme+"'")
	return nil
}

//
func (e *Env) assign(name *tokenObj, v value) {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = v
		e.init[name.lexeme] = true
		return
	}
	if e.enclosing != nil {
		e.enclosing.assign(name, v)
		return
	}

	runtimeErr(name, "undefined variable '"+name.lexeme+"'")
	return
}

// ------------------------------------------
// interpret

func interpret(stmt []Stmt, env *Env) (err error) {
	env.defineInit("clock", clockFn{})

	//handle panic and output , all kinds of interpret err
	defer func() {
		if e := recover(); e != nil {
			if b, ok := e.(BreakErr); ok {
				fmt.Printf("b.t = %+v\n", b.t.line)
				s := fmt.Sprintf("expected a while loop to break from at line %v ", b.t.line)
				err = errors.New(s)
				return
			}
			err = e.(RuntimeError)
		}
	}()
	for _, s := range stmt {
		s.execute(env)
	}
	return nil
}

// ------------------------------------------
// clockFn is a closure variable

type clockFn struct{}

func (c clockFn) arity() int {
	return 0
}

func (c clockFn) call(_ *Env, _ []value) value { // does not care about args , but to obey the interface standard
	return float64(time.Now().UnixNano())
}

// ------------------------------------------
// Function

type FunObj struct {
	decl    *FunStmt
	closure *Env
}

func (f *FunObj) arity() int {
	return len(f.decl.params)
}

func (f *FunObj) call(e *Env, args []value) (v value) {
	env := NewEnv(f.closure)          //create an env for function call
	for i, p := range f.decl.params { //args adds into env
		env.defineInit(p.lexeme, args[i])
	}

	defer func() {
		if e := recover(); e != nil {
			// return whatever value is being panicked at us from return stmt
			v = e.(ReturnHack) //TODO
		}
	}()
	execBlock(f.decl.body, env) //exec the func body with its env
	return nil
}

//stringfy fn , just fun name for now
func (f *FunObj) String() string {
	return fmt.Sprintf("<fn %v>", f.decl.name.lexeme)
}

// closure , anonymous function
type FunAnon struct {
	decl    *FunExpr
	closure *Env
}

func (f *FunAnon) arity() int {
	return len(f.decl.params)
}

func (f *FunAnon) call(e *Env, args []value) (v value) {
	// Should it use env that is passed by e?
	env := NewEnv(f.closure) // only difference between function
	for i, p := range f.decl.params {
		env.defineInit(p.lexeme, args[i])
	}

	defer func() {
		if e := recover(); e != nil {
			// return whatever value is being panicked at us from return stmt
			v = e.(ReturnHack)
		}
	}()
	execBlock(f.decl.body, env)
	return nil
}

func (f *FunAnon) String() string {
	s := []string{}
	for _, p := range f.decl.params {
		s = append(s, p.lexeme)
	}
	return fmt.Sprintf("<lambda (%v)>", strings.Join(s, ","))
}
