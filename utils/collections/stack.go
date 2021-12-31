package collections

// stack implements basic LIFO container
type stack[T any] struct {
	elems []T
}

// NewStack creates a new LIFO container containing the list of elements
func NewStack[T any](elems ...T) *stack[T] {
	return &stack[T]{append([]T{}, elems...)}
}

// Push insert elements at the end
func (stack *stack[T]) Push(elems ...T) {
	stack.elems = append(stack.elems, elems...)
}

// Peek access the top element, and return whether it has been found or not
func (stack *stack[T]) Peek() (elem T, found bool) {
	l := len(stack.elems)
	if l == 0 {
		return elem, false
	}
	elem = stack.elems[l-1]
	return elem, true
}

// Pop returns and removes the top element, and return whether it has been found or not
func (stack *stack[T]) Pop() (elem T, found bool) {
	l := len(stack.elems)
	if l == 0 {
		return elem, false
	}
	elem, stack.elems = stack.elems[l-1], stack.elems[:l-1]
	return elem, true
}

// IsEmpty checks whether the underlying container is empty
func (stack *stack[T]) IsEmpty() bool {
	return len(stack.elems) == 0
}
