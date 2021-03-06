package main

import "fmt"

func errorAt(line int, where string, msg string) string {
	return fmt.Sprintf("[line %v] error%v: %v", line, where, msg)
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlphaNum(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func errorAtToken(t *tokenObj, msg string) string {
	var e string
	if t.tok == EOF {
		e = errorAt(t.line, " at end", msg)
	} else {
		e = errorAt(t.line, " at '"+t.lexeme+"'", msg)
	}
	return e
}

func printExprAST(e Expr) string {
	switch o := e.(type) {
	case *BinaryExpr:
		return fmt.Sprintf("(%v %v %v)",
			o.operator.tok, printExprAST(o.left), printExprAST(o.right))
	//case *TernaryExpr:
	//	return fmt.Sprintf("(%v %v %v %v)",
	//		o.operator.tok, printExprAST(o.op1), printExprAST(o.op2), printExprAST(o.op3))
	case *UnaryExpr:
		return fmt.Sprintf("(%v %v)",
			o.operator.tok, printExprAST(o.right))
	case *GroupingExpr:
		return fmt.Sprintf("(group %v)", printExprAST(o.expression))
	case *LiteralExpr:
		return fmt.Sprintf("%v", o.value)
	default:
		panic("unexpected type of expr")
	}
}
