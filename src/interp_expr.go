package main

import (
	"fmt"
)

// ------------------------------------------
// Expression Eval

func (e *BinaryExpr) eval(env *Env) value {
	switch e.operator.tok {
	case Plus:
		x := e.left.eval(env)
		xval, xok := x.(float64)
		if xok {
			y := e.right.eval(env)
			yval, yok := y.(float64)
			if yok {
				return xval + yval
			}
			runtimeErr(e.operator, "expected number as right operand")
		}
		if xval, xok := x.(string); xok {
			if yval, yok := e.right.eval(env).(string); yok {
				return xval + yval
			}
			runtimeErr(e.operator, "expected string as right operand")
		}
		runtimeErr(e.operator, "operands must be two numbers or two strings")
	case Minus:
		xval, yval := e.evalFloats(env)
		return xval - yval
	case Slash:
		xval, yval := e.evalFloats(env)
		if yval == 0 {
			runtimeErr(e.operator, "division by zero")
		}
		return xval / yval
	case Star:
		xval, yval := e.evalFloats(env)
		return xval * yval
	case Greater:
		xval, yval := e.evalFloats(env)
		return xval > yval
	case GreaterEqual:
		xval, yval := e.evalFloats(env)
		return xval >= yval
	case Less:
		xval, yval := e.evalFloats(env)
		return xval < yval
	case LessEqual:
		xval, yval := e.evalFloats(env)
		return xval <= yval
	case EqualEqual:
		return e.equal(env)
	case BangEqual:
		return !e.equal(env)
	}
	return nil // Unreachable?
}

func (e *BinaryExpr) evalFloats(env *Env) (float64, float64) {
	x, ok := e.left.eval(env).(float64)
	if !ok {
		runtimeErr(e.operator, "left operand must be a number")
	}
	y, ok := e.right.eval(env).(float64)
	if !ok {
		runtimeErr(e.operator, "left operand must be a number")
	}
	return x, y
}

func (e *BinaryExpr) equal(env *Env) bool {
	x := e.left.eval(env)
	y := e.right.eval(env)
	if x == nil && y == nil {
		return true
	}
	if x == nil {
		return false
	}
	return x == y
}

//TODO; something that is callable ? , nested call ?
func (e *CallExpr) eval(env *Env) value {
	callee := e.callee.eval(env)
	args := make([]value, 0)
	for _, a := range e.args {
		args = append(args, a.eval(env))
	}
	if fn, ok := callee.(Callable); ok {
		if len(args) != fn.arity() {
			runtimeErr(e.paren,
				fmt.Sprintf("expected %v arguments but got %v", fn.arity(), len(args)))
		}
		return fn.call(env, args)
	} else {
		err := fmt.Sprintf("'%v' is not a function or class", callee)
		runtimeErr(e.paren, err)
		return nil
	}
}

// closure produce/eval a FunAnon object which is able to be execute
func (s *FunExpr) eval(env *Env) value {
	fn := &FunAnon{decl: s, closure: NewEnv(env)}
	return fn
}

// TODO; () ?
func (e *GroupingExpr) eval(env *Env) value {
	return e.expression.eval(env)
}

// like string or number , produce itself
func (e *LiteralExpr) eval(env *Env) value {
	return e.value
}

func (e *LogicalExpr) eval(env *Env) value {
	left := e.left.eval(env)
	if e.operator.tok == Or {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}
	return e.right.eval(env) //!false
}

func (e *GetExpr) eval(env *Env) value {
	object := e.object.eval(env)
	if o, ok := object.(*LoxInstance); ok {
		return o.get(e.name)
	}
	runtimeErr(e.name, "Only instance have properties")
	return nil
}

func (e *SetExpr) eval(env *Env) value {
	obj := e.object.eval(env)

	if o, ok := obj.(*LoxInstance); ok {
		vlue := e.vlue.eval(env)
		o.set(e.name, vlue)
		return vlue
	} else {
		//TODO
		panic(e.name.lexeme + "Only instances have fields.")
	}
}

func (e *ThisExpr) eval(env *Env) value {
	return env.lookUpVariable(e.keyword, e)
}

func (e *UnaryExpr) eval(env *Env) value {
	val := e.right.eval(env)
	switch e.operator.tok {
	case Minus:
		f, ok := val.(float64)
		if !ok {
			// TODO: handle this as error
			panic("not a float")
		}
		return -f
	case Bang:
		return !isTruthy(val)
	}
	// unreachable?
	return nil
}

// produce variable name
func (e *VarExpr) eval(env *Env) value {
	return env.lookUpVariable(e.name, e)
	//return env.get(e.name)
}

func (e *AssignExpr) eval(env *Env) value {
	/*

		public Object visitAssignExpr(Expr.Assign expr) {
		    Object value = evaluate(expr.value);
		lox/Interpreter.java
		in visitAssignExpr()
		replace 1 line

		    Integer distance = locals.get(expr);
		    if (distance != null) {
		      environment.assignAt(distance, expr.name, value);
		    } else {
		      globals.assign(expr.name, value);
		    }

		    return value;
	*/
	v := e.value.eval(env)
	distance, ok := locals.get(e)
	if ok {
		env.assignAt(distance, e.name, v)
	} else {
		env.globals.assign(e.name, v)
	}
	//env.assign(e.name, v)
	return v
}

// false and nil are the only falsey values
func isTruthy(v value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}
