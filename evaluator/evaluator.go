package evaluator

import (
	"Eldrlang/object"
	"Eldrlang/parser"
	"Eldrlang/util"
	"fmt"
)

type Evaluator struct {
	Context *object.Context
	Errors  util.Errors
	pkg     *parser.Package
}

func NewEvaluator(ctx *object.Context) *Evaluator {
	return &Evaluator{
		Context: ctx,
	}
}

func (e *Evaluator) Evaluate(pkg *parser.Package) object.Object {
	e.pkg = pkg
	var result object.Object
	for _, node := range pkg.Nodes {
		result = e.EvaluateNode(node)
	}
	return result
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
		e.Context.Set(e.asPkgName(stat.Name.Value), nil)
	case *parser.Prefix:
		return e.evalPrefix(stat)
	case *parser.Infix:
		return e.evalInfix(stat)
	case *parser.Block:
		return e.evalBlock(stat)
	case *parser.Conditional:
		return e.evalConditional(stat)
	case *parser.Return:
		return e.EvaluateNode(stat.Exp)
	}
	return nil
}

func (e *Evaluator) evalIdentifier(ident *parser.Identifier) object.Object {
	if val, ok := e.Context.Get(ident.Value); ok {
		return val
	}
	// TODO FIX BUILTIN
	/*if builtin, ok := builtins[ident.Value]; ok {
		return builtin
	}*/
	err := util.NewError(ident.Token, util.IdentNotFound,
		ident.Value)
	e.Errors.Add(err)
	return nil
}

func (e *Evaluator) evalPrefix(pref *parser.Prefix) object.Object {
	switch exp := e.EvaluateNode(pref.Right).(type) {
	case *object.Boolean:
		if pref.Operator != "!" {
			err := util.NewError(pref.Token, util.InvalidOpForO)
			e.Errors.Add(err)
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
			e.Errors.Add(err)
			return nil
		}
	}
	err := util.NewError(pref.Token, util.InvalidOpForO)
	e.Errors.Add(err)
	return nil
}

func (e *Evaluator) evalInfix(inf *parser.Infix) object.Object {
	right := e.EvaluateNode(inf.Right)
	switch left := e.EvaluateNode(inf.Left).(type) {
	case *object.String:
		str, valid := right.(*object.String)
		if !valid {
			err := util.NewError(inf.Token, util.InvalidOpComb)
			e.Errors.Add(err)
			return nil
		}
		switch inf.Operator {
		case "=":
			left.Value = str.Value
		case "+":
			return &object.String{
				Value: left.Value + str.Value}
		default:
			err := util.NewError(inf.Token, util.InvalidOpForO)
			e.Errors.Add(err)
			return nil
		}
	case *object.Integer:
		num, valid := right.(*object.Integer)
		if !valid {
			err := util.NewError(inf.Token, util.InvalidOpComb)
			e.Errors.Add(err)
			return nil
		}
		switch inf.Operator {
		case "=":
			left.Value = num.Value
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
			e.Errors.Add(err)
			return nil
		}
	case *object.Boolean:
		val, valid := right.(*object.Boolean)
		if !valid {
			err := util.NewError(inf.Token, util.InvalidOpComb)
			e.Errors.Add(err)
			return nil
		}
		switch inf.Operator {
		case "=":
			left.Value = val.Value
		case "==":
			return &object.Boolean{Value: left.Value == val.Value}
		case "!=":
			return &object.Boolean{Value: left.Value != val.Value}
		default:
			err := util.NewError(inf.Token, util.InvalidOpForO)
			e.Errors.Add(err)
			return nil
		}
	}
	return nil
}

func (e *Evaluator) evalBlock(block *parser.Block) object.Object {
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
	return result
}

func (e *Evaluator) evalConditional(ie *parser.Conditional) object.Object {
	condition := e.EvaluateNode(ie.Require)
	/*if isError(condition) {
		return condition
	}*/
	cond, valid := condition.(*object.Boolean)
	if !valid {
		return nil
	}
	if cond.Value {
		return e.EvaluateNode(ie.To)
	} else if ie.Else != nil {
		return e.EvaluateNode(ie.Else)
	} else {
		return &object.Null{}
	}
}

func (e *Evaluator) asPkgName(name string) string {
	return fmt.Sprintf("%s.%s", e.pkg.Namespace, name)
}
