package evaluator

import (
	"fmt"
	"math"

	"github.com/makramkd/go-monkey/ast"
	"github.com/makramkd/go-monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}

	NULL = &object.Null{}
)

func Eval(root ast.Node, env *object.Env) object.Object {
	switch node := root.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.ArrayAccessExpression:
		return evalArrayAccessExpression(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		// Both errors and returns should stop execution of the program
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Env) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			t := result.Type()
			if t == object.RETURN_VALUE || t == object.ERROR {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBoolean(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func evalIdentifier(ident *ast.Identifier, env *object.Env) object.Object {
	if v, ok := env.Get(ident.Value); ok {
		return v
	}

	if builtin, ok := builtins[ident.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", ident.Value)
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return nativeBoolToBoolean(leftVal && rightVal)
	case "||":
		return nativeBoolToBoolean(leftVal || rightVal)
	case "==":
		return nativeBoolToBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	// Arithmetic operators
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "**":
		return &object.Integer{Value: int64(math.Pow(float64(leftVal), float64(rightVal)))}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}

	// Comparison operators
	case "<":
		return nativeBoolToBoolean(leftVal < rightVal)
	case "<=":
		return nativeBoolToBoolean(leftVal <= rightVal)
	case ">":
		return nativeBoolToBoolean(leftVal > rightVal)
	case ">=":
		return nativeBoolToBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolean(leftVal != rightVal)

	// Boolean operators
	case "&&":
		return nativeBoolToBoolean((leftVal != 0) && (rightVal != 0))
	case "||":
		return nativeBoolToBoolean((leftVal != 0) || (rightVal != 0))

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: l + r}
	case "==":
		return nativeBoolToBoolean(l == r)
	case "!=":
		return nativeBoolToBoolean(l != r)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalNegativeOperator(right)
	default:
		return NULL
	}
}

func evalCallExpression(call *ast.CallExpression, env *object.Env) object.Object {
	v := Eval(call.Function, env)
	if isError(v) {
		return v
	}

	evaluatedArgs := []object.Object{}
	for _, arg := range call.Arguments {
		v := Eval(arg, env)
		if isError(v) {
			return v
		}
		evaluatedArgs = append(evaluatedArgs, v)
	}

	switch f := v.(type) {
	case *object.Function:
		fEnv := object.NewScopedEnv(f.Env)
		// evaluate arguments left to right, propagating errors as necessary
		// and set the appropriate values in the environment.
		for i, arg := range evaluatedArgs {
			fEnv.Set(f.Parameters[i].Value, arg)
		}

		// eval the function body with this new environment
		ret := Eval(f.Body, fEnv)

		return unwrapReturnValue(ret)
	case *object.Builtin:
		return f.F(evaluatedArgs...)
	default:
		return newError("not a function: %s", f.Type())
	}

}

func evalArrayLiteral(array *ast.ArrayLiteral, env *object.Env) object.Object {
	vals := []object.Object{}
	for _, exp := range array.Elements {
		v := Eval(exp, env)
		if isError(v) {
			return v
		}
		vals = append(vals, v)
	}

	return &object.Array{Values: vals}
}

func evalArrayAccessExpression(expr *ast.ArrayAccessExpression, env *object.Env) object.Object {
	// Evaluate index to make sure it's integral and >= 0
	index := Eval(expr.Index, env)

	if index.Type() != object.INTEGER {
		return newError("array access index must be integral, got: %s", index.Type())
	}

	idx := index.(*object.Integer)
	if idx.Value < 0 {
		return newError("array access index must be greater than zero, got: %d", idx.Value)
	}

	// Evaluate the array in the access expression
	var array *object.Array
	switch a := expr.Array.(type) {
	case *ast.ArrayLiteral:
		lit := Eval(a, env)
		if isError(lit) {
			return lit
		}
		array = lit.(*object.Array)
	case *ast.Identifier:
		i := Eval(a, env)
		if isError(i) {
			return i
		}
		if i.Type() != object.ARRAY {
			return newError("object is not of type array: %s, actual type: %s", a.Value, i.Type())
		}
		array = i.(*object.Array)
	}

	// Access the index at idx or return out of bounds error
	if int(idx.Value) >= len(array.Values) {
		return newError("out of bounds error: index %d is out of range for array", idx.Value)
	}

	return array.Values[int(idx.Value)]
}

func unwrapReturnValue(obj object.Object) object.Object {
	if rValue, ok := obj.(*object.ReturnValue); ok {
		return rValue.Value
	}
	return obj
}

func evalBangOperator(right object.Object) *object.Boolean {
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

func evalNegativeOperator(right object.Object) object.Object {
	switch e := right.(type) {
	case *object.Integer:
		return &object.Integer{
			Value: e.Value * -1,
		}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression, env *object.Env) object.Object {
	condition := Eval(exp.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(exp.Consequence, env)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(o object.Object) bool {
	switch o {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(message string, format ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(message, format...),
	}
}

func isError(o object.Object) bool {
	if o != nil {
		return o.Type() == object.ERROR
	}
	return false
}
