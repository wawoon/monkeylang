package evaluator

import (
	"github.com/wawoon/monkeylang/ast"
	"github.com/wawoon/monkeylang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	}
	return NULL
}

func evalBangExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return NULL
	}

	rightValue := right.(*object.Integer).Value
	return &object.Integer{Value: -rightValue}
}
