package evaluator

import (
	"Eldrlang/object"
	"Eldrlang/parser"
	"Eldrlang/util"
)

type Evaluator struct {
	*object.Runtime
	errors  util.Errors
	srcCode *parser.Package
}

func NewEvaluator() *Evaluator {
	return &Evaluator{Runtime: object.NewRuntime()}
}

func (e *Evaluator) Evaluate(src *parser.Package) *Output {
	e.srcCode = src
	var result = Output{}
	for _, node := range e.srcCode.Nodes {
		result.Object = e.EvaluateNode(node)
	}
	result.Errors = e.errors
	e.errors.Clear()
	return &result
}

func (e *Evaluator) EvaluateNode(node parser.Node) object.Object {
	switch stat := node.(type) {
	case *parser.Boolean:
		return &object.Boolean{Value: stat.Value}
	case *parser.Integer:
		return &object.Integer{Value: stat.Value}
	case *parser.String:
		return &object.String{Value: stat.Value}
	case *parser.Identifier:
		return e.evalIdentifier(stat)
	case *parser.Variable:
		e.SetValue(stat.Literal(), &object.Null{})
	case *parser.Assign:
		e.evalAssign(stat)
	case *parser.Prefix:
		return e.evalPrefix(stat)
	case *parser.Infix:
		return e.evalInfix(stat)
	case *parser.Block:
		return e.evalBlock(stat)
	case *parser.Conditional:
		return e.evalConditional(stat)
	case *parser.Function:
		e.evalFunction(stat)
	case *parser.FuncCall:
		return e.evalFuncCall(stat)
	case *parser.Return:
		return e.EvaluateNode(stat.Exp)
	}
	return nil
}

func (e *Evaluator) evalIdentifier(ident *parser.Identifier) object.Object {
	if val := e.GetValue(ident.Value); val != nil {
		return val
	}
	if builtin, ok := object.Builtins[ident.Value]; ok {
		return builtin
	}
	err := util.NewError(ident.Token, util.IdentNotFound, ident.Value)
	e.errors.Add(err)
	return nil
}

func (e *Evaluator) evalAssign(stat *parser.Assign) {
	e.EvaluateNode(stat.Left)
	name := stat.Left.Literal()
	if e.GetValue(name) != nil {
		e.SetValue(name, e.EvaluateNode(stat.Right))
	}
}

func (e *Evaluator) evalPrefix(pref *parser.Prefix) object.Object {
	switch exp := e.EvaluateNode(pref.Right).(type) {
	case *object.Boolean:
		if pref.Operator != "!" {
			err := util.NewError(pref.Token, util.InvalidOpForO)
			e.errors.Add(err)
			return nil
		}
		return &object.Boolean{Value: !exp.Value}
	case *object.Integer:
		switch pref.Operator {
		case "+":
			return &object.Integer{Value: exp.Value}
		case "-":
			return &object.Integer{Value: -exp.Value}
		default:
			err := util.NewError(pref.Token, util.InvalidOpForO)
			e.errors.Add(err)
			return nil
		}
	}
	err := util.NewError(pref.Token, util.InvalidOpForO)
	e.errors.Add(err)
	return nil
}

func (e *Evaluator) evalInfix(inf *parser.Infix) object.Object {
	right := e.EvaluateNode(inf.Right)
	switch left := e.EvaluateNode(inf.Left).(type) {
	case *object.String:
		str, valid := right.(*object.String)
		if !valid {
			err := util.NewError(inf.Token, util.InvalidOpComb)
			e.errors.Add(err)
			return nil
		}
		switch inf.Operator {
		case "+":
			return &object.String{
				Value: left.Value + str.Value}
		default:
			err := util.NewError(inf.Token, util.InvalidOpForO)
			e.errors.Add(err)
			return nil
		}
	case *object.Integer:
		num, valid := right.(*object.Integer)
		if !valid {
			err := util.NewError(inf.Token, util.InvalidOpComb)
			e.errors.Add(err)
			return nil
		}
		switch inf.Operator {
		case "+":
			return &object.Integer{Value: left.Value + num.Value}
		case "-":
			return &object.Integer{Value: left.Value - num.Value}
		case "*":
			return &object.Integer{Value: left.Value * num.Value}
		case "/":
			return &object.Integer{Value: left.Value / num.Value}
		case "<":
			return &object.Boolean{Value: left.Value < num.Value}
		case ">":
			return &object.Boolean{Value: left.Value > num.Value}
		case "==":
			return &object.Boolean{Value: left.Value == num.Value}
		case "!=":
			return &object.Boolean{Value: left.Value != num.Value}
		default:
			err := util.NewError(inf.Token, util.InvalidOpForO)
			e.errors.Add(err)
			return nil
		}
	case *object.Boolean:
		val, valid := right.(*object.Boolean)
		if !valid {
			err := util.NewError(inf.Token, util.InvalidOpComb)
			e.errors.Add(err)
			return nil
		}
		switch inf.Operator {
		case "==":
			return &object.Boolean{Value: left.Value == val.Value}
		case "!=":
			return &object.Boolean{Value: left.Value != val.Value}
		default:
			err := util.NewError(inf.Token, util.InvalidOpForO)
			e.errors.Add(err)
			return nil
		}
	}
	return nil
}

func (e *Evaluator) evalBlock(block *parser.Block) object.Object {
	e.PushChild()
	var result object.Object
	for _, statement := range block.Nodes {
		result = e.EvaluateNode(statement)
		if result != nil {
			if _, ok := statement.(*parser.Return); ok {
				return result
			}
			if _, ok := statement.(*parser.Break); ok {
				return nil
			}
		}
	}
	e.PopChild()
	return result
}

func (e *Evaluator) evalConditional(ie *parser.Conditional) object.Object {
	condition := e.EvaluateNode(ie.Require)
	cond, valid := condition.(*object.Boolean)
	if !valid {
		err := util.NewError(ie.Token, util.ExpectedCondV)
		e.errors.Add(err)
		return nil
	}
	if cond.Value {
		return e.EvaluateNode(ie.To)
	} else if ie.Else != nil {
		return e.EvaluateNode(ie.Else)
	} else {
		return nil
	}
}

func (e *Evaluator) evalFunction(f *parser.Function) {
	fun := &object.Function{Parameters: f.Params, Body: f.Body}
	e.SetValue(f.Name.Literal(), fun)
}

func (e *Evaluator) evalFuncCall(fc *parser.FuncCall) object.Object {
	ident := fc.Function.(*parser.Identifier)
	storedFun := e.EvaluateNode(ident)
	if storedFun == nil {
		err := util.NewError(ident.Token, util.IdentNotFound, ident.Value)
		e.errors.Add(err)
		return nil
	}
	var params []object.Object
	for _, a := range fc.Arguments {
		evaluated := e.EvaluateNode(a)
		params = append(params, evaluated)
	}
	switch fun := storedFun.(type) {
	case *object.Function:
		if len(fun.Parameters) != len(params) {
			err := util.NewError(ident.Token, util.ExpectedFuncP, len(fun.Parameters))
			e.errors.Add(err)
			return nil
		}
		return e.exeFuncCall(fun, params)
	case *object.Builtin:
		if fun.Size != len(params) {
			err := util.NewError(ident.Token, util.ExpectedFuncP, fun.Size)
			e.errors.Add(err)
			return nil
		}
		return fun.Fun(params...)
	default:
		err := util.NewError(ident.Token, util.IdentNotAFunc, ident.Literal())
		e.errors.Add(err)
		return nil
	}
}

func (e *Evaluator) exeFuncCall(fun *object.Function, params []object.Object) object.Object {
	e.PushChild()
	for index, param := range fun.Parameters {
		e.SetValue(param.Literal(), params[index])
	}
	result := e.EvaluateNode(fun.Body)
	e.PopChild()
	return result
}

func (e *Evaluator) exeBuiltin(fun *object.Function, params []object.Object) object.Object {
	e.PushChild()
	for index, param := range fun.Parameters {
		e.SetValue(param.Literal(), params[index])
	}
	return e.EvaluateNode(fun.Body)
}
