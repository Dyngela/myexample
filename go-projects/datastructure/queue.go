package main

type Queue[T any] struct {
	elements []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(element any) {
	q.elements = append(q.elements, element)
}

func (q *Queue[T]) Dequeue() (T, bool) {
	var element T
	if len(q.elements) == 0 {
		return element, false
	}
	element = q.elements[0]
	q.elements = q.elements[1:]
	return element, true
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue[T]) Size() int {
	return len(q.elements)
}

func (q *Queue[T]) Peek() (T, bool) {
	var element T
	if len(q.elements) == 0 {
		return element, false
	}
	return q.elements[0], true
}
