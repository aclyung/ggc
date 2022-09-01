package compiler

import (
	buitin "almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var BuiltinOpers = general.AsList([]*Operator{
	&Operator{Op: syntax.OperAdd, TypeL: buitin.Int, TypeR: buitin.Int, Operation: Add_Int_Int_Int}, // Add Int Int -> Int
	&Operator{Op: syntax.OperSub, TypeL: buitin.Int, TypeR: buitin.Int, Operation: Sub_Int_Int_Int}, // Sub Int Int -> Int
	&Operator{Op: syntax.OperMul, TypeL: buitin.Int, TypeR: buitin.Int, Operation: Mul_Int_Int_Int}, // Mul Int Int -> Int
	&Operator{Op: syntax.OperDiv, TypeL: buitin.Int, TypeR: buitin.Int, Operation: Div_Int_Int_Int}, // Div Int Int -> Int
	&Operator{Op: syntax.OperRem, TypeL: buitin.Int, TypeR: buitin.Int, Operation: Rem_Int_Int_Int}, // Rem Int Int -> Int
})

func (c *Compiler) InitOper() {
	BuiltinOpers.Each(func(o *Operator) { c.Opers.Append(o) })
}

var (
	Add_Int_Int_Int = func(c *Compiler, b *ir.Block, l, r value.Value) value.Value { return b.NewAdd(l, r) }
	Sub_Int_Int_Int = func(c *Compiler, b *ir.Block, l, r value.Value) value.Value { return b.NewSub(l, r) }
	Mul_Int_Int_Int = func(c *Compiler, b *ir.Block, l, r value.Value) value.Value { return b.NewMul(l, r) }
	Div_Int_Int_Int = func(c *Compiler, b *ir.Block, l, r value.Value) value.Value { return b.NewSDiv(l, r) }
	Rem_Int_Int_Int = func(c *Compiler, b *ir.Block, l, r value.Value) value.Value { return b.NewSRem(l, r) }
)

type Operator struct {
	Op    syntax.Token
	TypeL types.Type
	// TypeR == nil means Unary Operation
	TypeR     types.Type
	Body      *syntax.BlockStmt
	Operation func(*Compiler, *ir.Block, value.Value, value.Value) value.Value
}

func (Op *Operator) Name() string {
	return OperName(Op.Op, Op.TypeL, Op.TypeR)
}

func NewOperator(Op syntax.Token, TypeL, TypeR types.Type, body *syntax.BlockStmt) *Operator {
	op := &Operator{
		Op:    Op,
		TypeL: TypeL,
		TypeR: TypeR,
		Body:  body,
	}
	op.Operation = op._CompileOperation()
	return op
}

func (Op *Operator) _CompileOperation() func(c *Compiler, b *ir.Block, r, l value.Value) value.Value {
	return func(c *Compiler, b *ir.Block, r, l value.Value) value.Value {
		for _, v := range Op.Body.StmtList[:len(Op.Body.StmtList)-1] {
			c.CompileStmt(b, v)
		}
		if v, ok := Op.Body.StmtList[len(Op.Body.StmtList)-1].(*syntax.ReturnStmt); ok {
			return c.CompileExpr(b, v.Return)
		}
		panic("expected return value for operation name " + Op.Name())
	}
}
