package table

import (
	boundNode "main/legacy/binding/boundNode"
	"main/legacy/general"
)

type Table struct {
	Map map[general.VariableSymbol]boundNode.BoundExpression
}

func (t *Table) Get(v general.VariableSymbol) {

}
