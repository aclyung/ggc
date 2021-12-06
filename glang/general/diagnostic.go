package general

import (
	"almeng.com/glang/general/TextSpan"
	"fmt"
)

type Diags struct {
	Notions []Diag
}

func NewDiag() Diags {
	return Diags{}
}

func ConcatDiag(diagA Diags, diagB Diags) Diags {
	notions := make([]Diag, 0)
	notions = append(notions, diagA.Notions...)
	notions = append(notions, diagB.Notions...)
	return Diags{notions}
}

func (d *Diags) Diagnose(span TextSpan.TextSpan, text string, l Level) {
	diag := Diag{span, text, l}
	d.Notions = append(d.Notions, diag)
}

func (d *Diags) InvalidNumber(span TextSpan.TextSpan, text string, numType string) {
	d.Diagnose(span, fmt.Sprint("the number ", text, " is not valid ", numType), ERROR)
}

func (d *Diags) BadCharacter(span TextSpan.TextSpan, c string) {
	d.Diagnose(span, "Illegal character '"+c+"'", ERROR)
}

func (d *Diags) UnexpectedToken(span TextSpan.TextSpan, wanted string, got string) {
	d.Diagnose(span, " Expected <"+wanted+">, got <"+got+">", WARN)
}

func (d *Diags) UndefinedBinaryOperator(span TextSpan.TextSpan, oper string, left string, right string) {
	d.Diagnose(span, "Binary Operator '"+oper+"' is not defined for types "+left+" and "+right, ERROR)
}

func (d *Diags) VariableTypeMisMatch(span TextSpan.TextSpan, name string, varType string, expType string) {
	d.Diagnose(span, "Variable '"+name+"' is "+varType+", not "+expType, ERROR)
}

func (d *Diags) UndefinedIdentifier(span TextSpan.TextSpan, name string) {
	d.Diagnose(span, "Undefined identifier '"+name+"'", ERROR)
}

func (d *Diags) UndefinedUnaryOperator(span TextSpan.TextSpan, oper string, operand string) {
	d.Diagnose(span, "Unary Operator '"+oper+"' is not defined for type "+operand, ERROR)
}

type Diag struct {
	Span TextSpan.TextSpan
	Msg  string
	Lev  Level
}
