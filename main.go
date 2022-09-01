package main

import (
	"almeng.com/glang/core/compiler"
	"flag"
)

//	"github.com/c-bata/go-prompt"
//)
//
//const config = "./config.json"
//
//func completer(d prompt.Document) []prompt.Suggest {
//	return prompt.FilterHasPrefix(nil, d.GetWordBeforeCursor(), true)
//}

func main() {
	debug := flag.Bool("d", false, "enable debug diagnosis")
	verbose := flag.Bool("v", false, "enable verbose debug diagnosis")
	*debug = *debug || *verbose
	flag.Parse()
	path := "./main.gg"
	compiler.Compile(path, *debug, *verbose)
}

//type Node[T any] struct {
//	Element T
//	Next    *Node[T]
//}
//
//type LinkList[T any] struct {
//	Root *Node[T]
//	len  int
//}
//
//func emptyRoot[T any]() *Node[T] {
//	return &Node[T]{
//		Next: nil,
//	}
//}
//
//func (l *LinkList[T]) Len() int {
//	return l.len
//}
//
//func NewList[T any]() LinkList[T] {
//	return LinkList[T]{Root: emptyRoot[T](), len: 0}
//}
//
//func (l LinkList[T]) String() string {
//	str := "[ "
//	cur := l.Root
//	if l.len == 0 {
//		return "[]"
//	}
//	for cur.Next != nil {
//		str += fmt.Sprint(cur.Element) + " "
//		cur = cur.Next
//	}
//	str += fmt.Sprint(cur.Element) + " ]"
//	return str
//}
//
//func (l *LinkList[T]) Set(i int, e T) error {
//	n := l.getNode(i)
//	n.Element = e
//	return nil
//}
//
//func (l *LinkList[T]) Add(e T) {
//	node := &Node[T]{
//		Element: e,
//		Next:    nil,
//	}
//	last := l.lastNode()
//	if l.len == 0 {
//		*last = *node
//		goto rtn
//	}
//	last.Next = node
//rtn:
//	l.len++
//	return
//}
//
//func (l *LinkList[T]) Remove(i int) T {
//	if l.len == 0 {
//		panic("the list is empty")
//	}
//	rtn := l.getNode(i).Element
//	if i == 0 {
//		*(l.Root) = *(l.Root.Next)
//	} else if l.len-1 == i {
//		n := l.getNode(i - 1)
//		n.Next = nil
//	} else {
//		c := l.getNode(i - 1)
//		n := l.getNode(i + 1)
//		*(c.Next) = *n
//	}
//	l.len--
//	return rtn
//}
//
//func (l *LinkList[T]) lastNode() *Node[T] {
//	cur := l.Root
//	for cur.Next != nil {
//		cur = cur.Next
//	}
//	return cur
//}
//
//func (l *LinkList[T]) getNode(i int) *Node[T] {
//	if l.len-1 < i {
//		panic("out of bound Exception")
//	}
//	cur := l.Root
//	for index := 0; index < i; index++ {
//		cur = cur.Next
//	}
//	return cur
//}
//
//func (l *LinkList[T]) Last() T {
//	if l.len == 0 {
//		panic("list is empty")
//	}
//	cur := l.Root
//	for cur.Next != nil {
//		cur = cur.Next
//	}
//	return cur.Element
//}
//
//func (l *LinkList[T]) Get(i int) T {
//	if l.len-1 < i {
//		panic("out of bound Exception")
//	}
//	cur := l.Root
//	for index := 0; index < i; index++ {
//		cur = cur.Next
//	}
//	return cur.Element
//}

//l := NewList[int]()
//l.Add(10)
//l.Add(20)
//l.Add(10)
//l.Add(10)
//l.Add(20)
//l.Add(10)
//l.Add(10)
//l.Add(20)
//l.Add(10)
//l.Add(10)
//l.Add(20)
//l.Add(10)
//err := l.Set(2, 300)
//if err != nil {
//	return
//}
//fmt.Println(l.Remove(1))
//fmt.Println(l.Len())
