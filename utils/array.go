package utils

// ArrayFilter returns a new array containing all the elements for which the condition returned true
func ArrayFilter[V any](input []V, cond func(V) bool) []V {
	output := make([]V, 0, len(input))
	for _, v := range input {
		if cond(v) {
			output = append(output, v)
		}
	}
	return output
}

// ArrayMap returns a new array containing for each key of the source array, the value returned by the callback function
func ArrayMap[In any, Out any](input []In, callback func(In) Out) []Out {
	output := make([]Out, len(input))
	for i, v := range input {
		output[i] = callback(v)
	}
	return output
}

// UnorderedArray implements a utility container allowing fast deletion without preserving items order
type UnorderedArray[K any] []K

// NewUnorderedArray is a quick way to get an UnorderedArray without worrying about specification and value assignation
func NewUnorderedArray[K any](values ...K) UnorderedArray[K] {
	return append(UnorderedArray[K]{}, values...)
}

// Values returns the container values
func (unorderedArray UnorderedArray[K]) Values() []K {
	return unorderedArray
}

// Insert appends elements to the end
func (unorderedArray *UnorderedArray[K]) Insert(elems ...K) {
	*unorderedArray = append(*unorderedArray, elems...)
}

// Remove removes element at index i by copying at this index the last one, then shrink the container
func (unorderedArray *UnorderedArray[K]) Remove(i int) (old K) {
	if l := len(*unorderedArray); i >= 0 && i < l {
		old = (*unorderedArray)[i]
		(*unorderedArray)[i] = (*unorderedArray)[l-1]
		*unorderedArray = (*unorderedArray)[:l-1]
	}
	return old
}

// OrderedArray implements a utility container allowing deletion preserving items order
type OrderedArray[K any] []K

// NewOrderedArray is a quick way to get an OrderedArray without worrying about specification and value assignation
func NewOrderedArray[K any](values ...K) OrderedArray[K] {
	return append(OrderedArray[K]{}, values...)
}

// Values returns the container values
func (orderedArray OrderedArray[K]) Values() []K {
	return orderedArray
}

// Insert appends elements to the end
func (orderedArray *OrderedArray[K]) Insert(elems ...K) {
	*orderedArray = append(*orderedArray, elems...)
}

// Remove removes element at index i
func (orderedArray *OrderedArray[K]) Remove(i int) (old K) {
	if l := len(*orderedArray); i >= 0 && i < l {
		old = (*orderedArray)[i]
		*orderedArray = append((*orderedArray)[:i], (*orderedArray)[i+1:]...)
	}
	return old
}

// CyclicArray implements a utility container allowing rotative access to elements
type CyclicArray[K any] struct {
	values []K
	index  int
}

// NewCyclicArray is a quick way to get an CyclicArray without worrying about specification and value assignation
func NewCyclicArray[K any](values ...K) CyclicArray[K] {
	return CyclicArray[K]{values: append([]K{}, values...)}
}

// Values returns the container values
func (cyclicArray CyclicArray[K]) Values() []K {
	return cyclicArray.values
}

// Index returns the current element index pointed
func (cyclicArray CyclicArray[K]) Index() int {
	return cyclicArray.index % len(cyclicArray.values)
}

// Get returns the value at index i and set the current index to this index
func (cyclicArray *CyclicArray[K]) Get(i int) K {
	cyclicArray.index = i % len(cyclicArray.values)
	return cyclicArray.values[cyclicArray.index]
}

// Next returns the value currently pointed, and move pointer to the next location
func (cyclicArray *CyclicArray[K]) Next() (value K) {
	value = cyclicArray.Get(cyclicArray.index)
	cyclicArray.index++
	return
}
