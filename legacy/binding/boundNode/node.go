package boundNode

import "reflect"

// Node interface

type BoundNode interface {
	Kind() Kind
	Type() reflect.Kind
}

// BoundExpression is sub structure of BoundNode

type BoundExpression = BoundNode
