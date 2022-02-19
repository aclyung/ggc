package evaluator

import (
	"main/legacy/binding"
	"main/legacy/general"
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
