package compile

import (
	"main/legacy/ast/tree"
	binding2 "main/legacy/binding"
	"main/legacy/binding/boundNode"
	evaluator2 "main/legacy/evaluator"
	general2 "main/legacy/general"
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
