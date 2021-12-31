package collections

// queue implements basic FIFO container
type queue[T any] struct {
	elems []T
}

// NewQueue creates a new FIFO container containing the list of elements
func NewQueue[T any](elems ...T) *queue[T] {
	return &queue[T]{append([]T{}, elems...)}
}

// Enqueue insert elements at the end
func (queue *queue[T]) Enqueue(elems ...T) {
	queue.elems = append(queue.elems, elems...)
}

// Peek access the top element, and return whether it has been found or not
func (queue *queue[T]) Peek() (elem T, found bool) {
	if queue.IsEmpty() {
		return elem, false
	}
	elem = queue.elems[0]
	return elem, true
}

// Dequeue returns and removes the top element, and return whether it has been found or not
func (queue *queue[T]) Dequeue() (elem T, found bool) {
	if queue.IsEmpty() {
		return elem, false
	}
	elem, queue.elems = queue.elems[0], queue.elems[1:]
	return elem, true
}

// IsEmpty checks whether the underlying container is empty
func (queue *queue[T]) IsEmpty() bool {
	return len(queue.elems) == 0
}
