package evaluator

import (
	"almeng.com/glang/legacy/binding"
	"almeng.com/glang/legacy/general"
)

type EvaluationResult struct {
	Diags general.Diags
	*binding.BoundLiteralExpression
}

// Invalid Result constant
var InvalidResult = &EvaluationResult{}

func Result(diag general.Diags, val *binding.BoundLiteralExpression) EvaluationResult {
	return EvaluationResult{diag, val}
}
