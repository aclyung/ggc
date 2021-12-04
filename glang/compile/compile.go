package compile

import (
	"almeng.com/glang/ast/tree"
	"almeng.com/glang/binding"
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/evaluator"
	"almeng.com/glang/general"
)

type Compiler struct {
	Syntax tree.Tree
}

func NewCompiler(syntax tree.Tree) Compiler {
	return Compiler{syntax}
}

func (c *Compiler) Evaluate(vars *map[general.VariableSymbol]boundNode.BoundExpression) evaluator.EvaluationResult {
	binder := binding.NewBinder(vars)
	boundExp := binder.Bind(c.Syntax.Root)
	diag := general.ConcatDiag(c.Syntax.Diagnostics, binder.Diag)
	if len(diag.Notions) > 0 {
		return evaluator.Result(diag, binding.InvalidLiteraExpression)
	}
	eval := evaluator.NewEvaluator(boundExp, vars)
	res := eval.Evaluate().(*binding.BoundLiteralExpression)
	return evaluator.Result(diag, *res)
}
