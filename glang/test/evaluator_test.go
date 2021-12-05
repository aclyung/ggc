package test

import (
	"almeng.com/glang/ast"
	"almeng.com/glang/binding/boundNode"
	"almeng.com/glang/compile"
	"almeng.com/glang/general"
	assert "almeng.com/glang/general/Assert"
	"testing"
)

var data = map[string]interface{}{
	"1":              int64(1),
	"+1":             int64(1),
	"-1":             int64(-1),
	"14+12":          int64(26),
	"12-3":           int64(9),
	"4 *  2":         int64(8),
	"9/3":            int64(3),
	"11 -3 == 8":     true,
	"true == !false": true,
	"true":           true,
	"false":          false,
	"!true":          false,
	"!false":         true,
}

func Eval(text string, value interface{}) {
	tree := ast.ParseTree(text)
	comp := compile.NewCompiler(tree)
	vars := &map[general.VariableSymbol]boundNode.BoundExpression{}
	res := comp.Evaluate(vars)
	assert.Empty(res.Diags.Notions)
	assert.Equal(value, res.Value)

}

func TestEvaluation(t *testing.T) {
	t.Run("eval", func(t *testing.T) {
		for k, v := range data {
			Eval(k, v)
		}
	})
}
