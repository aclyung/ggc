package compile

import (
	"almeng.com/glang/legacy/ast/tree"
	binding2 "almeng.com/glang/legacy/binding"
	"almeng.com/glang/legacy/binding/boundNode"
	evaluator2 "almeng.com/glang/legacy/evaluator"
	general2 "almeng.com/glang/legacy/general"
)

type Compiler struct {
	Syntax tree.Tree
}

func NewCompiler(syntax tree.Tree) Compiler {
	return Compiler{syntax}
}

func (c *Compiler) Evaluate(vars *map[general2.VariableSymbol]boundNode.BoundExpression) evaluator2.EvaluationResult {
	defer func() {
		recover()

	}()
	globalScope := binding2.BindBoundGlobalScope(c.Syntax.Root)
	diag := general2.ConcatDiag(c.Syntax.Diagnostics, globalScope.Diags())
	//boundExp := globalScope.Bind(c.Syntax.Root.Expression)
	if len(diag.Notions) > 0 {
		return evaluator2.Result(diag, binding2.InvalidLiteralExpression)
	}
	eval := evaluator2.NewEvaluator(globalScope.Expression(), vars)
	res := eval.Evaluate()
	if res != nil {
		return evaluator2.Result(diag, res.(*binding2.BoundLiteralExpression))
	}
	return evaluator2.Result(diag, binding2.InvalidLiteralExpression)
}
