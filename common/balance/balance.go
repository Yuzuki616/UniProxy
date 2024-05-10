package balance

import "math/rand"

type List[T any] struct {
	balance  string // round otr random
	index    int
	elements []T
}

func New[T any](balance string, elements []T) *List[T] {
	return &List[T]{
		balance:  balance,
		index:    0,
		elements: elements,
	}
}

func (l *List[T]) Next() T {
	var temp T
	if len(l.elements) == 0 {
		return temp
	}
	if len(l.elements) == 1 {
		return l.elements[0]
	}
	switch l.balance {
	case "round":
		temp = l.elements[l.index]
		l.index = (l.index) % len(l.elements)
	case "random":
		temp = l.elements[rand.Intn(len(l.elements))]
	}
	return temp
}
