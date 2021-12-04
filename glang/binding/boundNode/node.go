package boundNode

import "reflect"

type BoundNode interface {
	Kind() BoundNodeKind
	Type() reflect.Type
}

type BoundExpression = BoundNode
