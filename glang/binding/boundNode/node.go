package boundNode

import "reflect"

type BoundNode interface {
	Kind() BoundNodeKind
	Type() reflect.Kind
}

type BoundExpression = BoundNode
