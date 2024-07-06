package main

import "fmt"

type Stack[T any] struct {
	data []T
	Eq   func(a, b T) bool
}

func NewStack[T any](eq func(a, b T) bool) *Stack[T] {
	return &Stack[T]{
		Eq: eq,
	}
}

func (s *Stack[T]) Push(element T) {
	s.data = append(s.data, element)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		return s.data[0], false
	}
	element := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return element, true
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}

func (s *Stack[T]) Size() int {
	return len(s.data)
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		return s.data[0], false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Clear() {
	s.data = []T{}
}

func (s *Stack[T]) Contains(element T) bool {
	for _, e := range s.data {
		if s.Eq(element, e) {
			return true
		}
	}
	return false
}

func (s *Stack[T]) Values() []T {
	return s.data
}

func (s *Stack[T]) Clone() *Stack[T] {
	clone := NewStack[T](s.Eq)
	for _, e := range s.data {
		clone.Push(e)
	}
	return clone
}

func (s *Stack[T]) Equals(other *Stack[T]) bool {
	if s.Size() != other.Size() {
		return false
	}
	for i, e := range s.data {
		if !s.Eq(e, other.data[i]) {
			return false
		}
	}
	return true
}

func (s *Stack[T]) String() string {
	return fmt.Sprintf("%v", s.data)
}
