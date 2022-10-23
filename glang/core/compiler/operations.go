package compiler

import (
	"almeng.com/glang/core/builtin"
	"almeng.com/glang/core/syntax"
	"github.com/almenglee/general"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

var (
	Int  = builtin.Int
	Bool = builtin.Bool
)

var BuiltinOpers = general.AsList([]*Operator{
	{Op: syntax.Add, OperType: syntax.OperAdd, RtnType: Int, TypeL: Int, TypeR: Int, Operation: Add_Int_Int_Int},   // Add Int Int -> Int
	{Op: syntax.Sub, OperType: syntax.OperSub, RtnType: Int, TypeL: Int, TypeR: Int, Operation: Sub_Int_Int_Int},   // Sub Int Int -> Int
	{Op: syntax.Mul, OperType: syntax.OperMul, RtnType: Int, TypeL: Int, TypeR: Int, Operation: Mul_Int_Int_Int},   // Mul Int Int -> Int
	{Op: syntax.Div, OperType: syntax.OperDiv, RtnType: Int, TypeL: Int, TypeR: Int, Operation: Div_Int_Int_Int},   // Div Int Int -> Int
	{Op: syntax.Rem, OperType: syntax.OperRem, RtnType: Int, TypeL: Int, TypeR: Int, Operation: Rem_Int_Int_Int},   // Rem Int Int -> Int
	{Op: syntax.Eql, OperType: syntax.OperEql, RtnType: Bool, TypeL: Int, TypeR: Int, Operation: Eql_Int_Int_Bool}, // Eql Int Int -> Bool
	{Op: syntax.Gtr, OperType: syntax.OperGtr, RtnType: Bool, TypeL: Int, TypeR: Int, Operation: Gtr_Int_Int_Bool}, // Gtr Int Int -> Bool

	// not operator for bool
	{Op: syntax.Not, OperType: syntax.OperNot, RtnType: Bool, TypeL: Bool, TypeR: nil, Operation: Not_Bool_Bool}, // Not Bool Bool -> Bool
})

func (c *Compiler) InitOper() {
	BuiltinOpers.Each(func(o *Operator) { c.Opers.Append(o) })
}

var (
	Add_Int_Int_Int  = func(c *Compiler, ctx *Context, l, r value.Value) value.Value { return ctx.Block.NewAdd(l, r) }
	Sub_Int_Int_Int  = func(c *Compiler, ctx *Context, l, r value.Value) value.Value { return ctx.Block.NewSub(l, r) }
	Mul_Int_Int_Int  = func(c *Compiler, ctx *Context, l, r value.Value) value.Value { return ctx.Block.NewMul(l, r) }
	Div_Int_Int_Int  = func(c *Compiler, ctx *Context, l, r value.Value) value.Value { return ctx.Block.NewSDiv(l, r) }
	Rem_Int_Int_Int  = func(c *Compiler, ctx *Context, l, r value.Value) value.Value { return ctx.Block.NewSRem(l, r) }
	Eql_Int_Int_Bool = func(c *Compiler, ctx *Context, l, r value.Value) value.Value {
		return ctx.Block.NewICmp(enum.IPredEQ, l, r)
	}
	Gtr_Int_Int_Bool = func(c *Compiler, ctx *Context, l, r value.Value) value.Value {
		return ctx.Block.NewICmp(enum.IPredSGT, l, r)
	}
	Not_Bool_Bool = func(c *Compiler, ctx *Context, l, r value.Value) value.Value {
		return ctx.Block.NewICmp(enum.IPredEQ, l, ctx.localize(ctx.lookupVariable("false")))
	}
)

type Operator struct {
	Op       syntax.Operator
	OperType syntax.Token
	RtnType  types.Type
	TypeL    types.Type
	// TypeR == nil means Unary Operation
	TypeR     types.Type
	Body      *syntax.BlockStmt
	_cntr     map[string]int
	Operation func(*Compiler, *Context, value.Value, value.Value) value.Value
}

// Name For debugging conveniences
func (Op *Operator) Name() string {
	name := "oper."
	oper := Op.OperType
	if Op.TypeR == nil {
		name += "unary." + oper.String() + "._" + Op.TypeL.Name() + "._" + Op.RtnType.Name()
		return name
	}
	name += "binary."
	nameL, nameR := "_"+Op.TypeL.Name(), "_"+Op.TypeR.Name()
	if oper.IsReversedOper() {
		name += (oper - (syntax.Reversed_oper - syntax.Operator_beg)).String() + "."
		name += nameR + nameL
	} else {
		name += oper.String() + "."
		name += nameL + nameR
	}
	name += "._" + Op.RtnType.Name()
	return name
}

func NewOperator(Op *syntax.OperDecl, TypeL, TypeR, RtnType types.Type) *Operator {
	return (&Operator{
		Op:       Op.Oper.OperTokenToOperator(),
		OperType: Op.Oper,
		TypeL:    TypeL,
		TypeR:    TypeR,
		RtnType:  RtnType,
		Body:     Op.Body,
	})._CompileOperation(Op)
}

func (Op *Operator) _CompileOperation(decl *syntax.OperDecl) *Operator {
	Op.Operation = func(c *Compiler, ctx *Context, r, l value.Value) value.Value {
		// to make sure that the operation is under the 'space' context
		// create sub context of the current 'space' context, whose Block is the given context's Block

		origin := c.CurrentSpace.Block
		opctx := c.CurrentSpace.NewContext(ctx.Block.Parent.NewBlock(Op.Name() + ".entry"))
		block := ctx.Block
		block.NewBr(opctx.Block)
		leave := ir.NewBlock(Op.Name() + ".leave")
		opctx.leave = leave
		opctx.vars[decl.TypeL.Name.Value] = l
		opctx.vars[decl.TypeR.Name.Value] = r
		var operRtn value.Value = opctx.NewAlloca(Op.RtnType)
		for _, v := range Op.Body.StmtList {
			c.CompileStmt(opctx, v, true, &operRtn)
		}
		block.Parent.Blocks = append(block.Parent.Blocks, leave)
		leave.Parent = block.Parent
		ctx.Block = leave
		rtn := ctx.NewLoad(Op.RtnType, operRtn)
		//if opctx.Block != block {
		//	ctx.Block = opctx.Block
		//}
		c.CurrentSpace.Block = origin

		return rtn
		//if v, ok := Op.Body.StmtList[len(Op.Body.StmtList)-1].(*syntax.ReturnStmt); ok {
		//	res := c.CompileExpr(opctx, v.Return)
		//	if opctx.Block != block {
		//		ctx.Block = opctx.Block
		//	}
		//	c.CurrentSpace.Block = origin
		//	return res
		//}
		//panic("expected return value for operation name " + Op.Name())
	}
	return Op
}
