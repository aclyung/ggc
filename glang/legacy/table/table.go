package table

import (
	boundNode "almeng.com/glang/legacy/binding/boundNode"
	"almeng.com/glang/legacy/general"
)

type Table struct {
	Map map[general.VariableSymbol]boundNode.BoundExpression
}

func (t *Table) Get(v general.VariableSymbol) {

}
