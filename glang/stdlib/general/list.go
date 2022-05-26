package general

type List[T any] []T

func join[T any](a List[T], b List[T]) *List[T] {
	rtn := a
	b.Each(func(e T) {
		rtn.Append(e)
	})
	return &rtn
}

func (l List[T]) Len() int {
	return len([]T(l))
}

func ToList[T any](l []T) *List[T] {
	return NewList(l...)
}

func NewList[T any](e ...T) *List[T] {
	r := List[T](e)
	return &r
}

func (l *List[T]) Each(f func(T)) {
	for _, t := range *l {
		f(t)
	}
}

func (l *List[T]) Iter(f func(int, T)) {
	for i, t := range *l {
		f(i, t)
	}
}

func (l *List[T]) Append(e T) {
	*l = List[T](append([]T(*l), e))
}

func (l *List[T]) Filter(cmp func(int, T) bool) *List[T] {
	//var rtn []T
	rtn := new(List[T])
	for i, v := range *l {
		if cmp(i, v) {
			rtn.Append(v)
		}
	}
	return rtn
}

func (l *List[T]) Take(n int) *List[T] {
	rtn := new(List[T])
	if len([]T(*l)) < n {
		panic("index out of bound")
	}
	for i := 0; i < n; i++ {
		rtn.Append((*l)[i])
	}
	return rtn
}
