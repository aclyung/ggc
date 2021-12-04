package evaluator

import (
	"almeng.com/glang/binding"
	"almeng.com/glang/general"
)

type EvaluationResult struct {
	Diags general.Diags
	binding.BoundLiteralExpression
}

func Result(diag general.Diags, val binding.BoundLiteralExpression) EvaluationResult {
	return EvaluationResult{diag, val}
}
