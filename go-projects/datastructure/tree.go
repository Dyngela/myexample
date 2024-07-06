package main

type Tree[T any] struct {
	Root  *Node[T]
	Less  func(a, b T) bool
	Equal func(a, b T) bool
}

type Node[T any] struct {
	Value T
	Left  *Node[T]
	Right *Node[T]
}

func NewTree[T any](less func(a, b T) bool, eq func(a, b T) bool) *Tree[T] {
	return &Tree[T]{
		Less:  less,
		Equal: eq,
	}
}

func (t *Tree[T]) Insert(value T) {
	if t.Root == nil {
		t.Root = &Node[T]{Value: value}
		return
	}

	t.insert(t.Root, value)
}

func (t *Tree[T]) insert(node *Node[T], value T) {
	if t.Less != nil && t.Less(value, node.Value) {
		if node.Left == nil {
			node.Left = &Node[T]{Value: value}
		} else {
			t.insert(node.Left, value)
		}
	} else {
		if node.Right == nil {
			node.Right = &Node[T]{Value: value}
		} else {
			t.insert(node.Right, value)
		}
	}
}

func (t *Tree[T]) InOrder() []T {
	var result []T
	t.inOrder(t.Root, &result)
	return result
}

func (t *Tree[T]) inOrder(node *Node[T], result *[]T) {
	if node == nil {
		return
	}

	t.inOrder(node.Left, result)
	*result = append(*result, node.Value)
	t.inOrder(node.Right, result)
}

func (t *Tree[T]) PreOrder() []T {
	var result []T
	t.preOrder(t.Root, &result)
	return result
}

func (t *Tree[T]) preOrder(node *Node[T], result *[]T) {
	if node == nil {
		return
	}

	*result = append(*result, node.Value)
	t.preOrder(node.Left, result)
	t.preOrder(node.Right, result)
}

func (t *Tree[T]) PostOrder() []T {
	var result []T
	t.postOrder(t.Root, &result)
	return result
}

func (t *Tree[T]) postOrder(node *Node[T], result *[]T) {
	if node == nil {
		return
	}

	t.postOrder(node.Left, result)
	t.postOrder(node.Right, result)
	*result = append(*result, node.Value)
}

func (t *Tree[T]) Search(value T) bool {
	return t.search(t.Root, value)
}

func (t *Tree[T]) search(node *Node[T], value T) bool {
	if node == nil {
		return false
	}

	if t.Equal != nil && t.Equal(node.Value, value) {
		return true
	}

	if t.Less != nil && t.Less(value, node.Value) {
		return t.search(node.Left, value)
	}

	return t.search(node.Right, value)
}

func (t *Tree[T]) Delete(value T) {
	t.Root = t.delete(t.Root, value)
}

func (t *Tree[T]) delete(node *Node[T], value T) *Node[T] {
	if node == nil {
		return nil
	}

	if t.Less != nil && t.Less(value, node.Value) {
		node.Left = t.delete(node.Left, value)
	} else if t.Less != nil && t.Less(node.Value, value) {
		node.Right = t.delete(node.Right, value)
	} else {
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}

		node.Value = t.minValue(node.Right)
		node.Right = t.delete(node.Right, node.Value)
	}

	return node
}

func (t *Tree[T]) minValue(node *Node[T]) T {
	for node.Left != nil {
		node = node.Left
	}
	return node.Value
}
