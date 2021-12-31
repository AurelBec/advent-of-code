package utils

import (
	"slices"
)

// ArrayIota fills an array with n values sequentially increasing starting with start
func ArrayIota(start, n int) []int {
	array := make([]int, n)
	for i := 0; i < n; i++ {
		array[i] = start + i
	}
	return array
}

// ArrayDelete returns a new array containing all the elements for which the condition returned false
func ArrayDelete[S ~[]E, E any](input S, cond func(E) bool) S {
	return slices.DeleteFunc(input, cond)
}

// ArrayMove returns a new array with element moved from i to j
func ArrayMove[S ~[]E, E any](s S, i, j int) S {
	value := s[i]
	return slices.Insert(slices.Delete(s, i, i+1), j, value)
}

// ArrayTranspose returns the transposed input array
func ArrayTranspose[S ~[][]E, E any](s S) S {
	y := len(s)
	if y == 0 {
		return s
	}

	x := len(s[0])
	result := make(S, x)
	for i := range result {
		result[i] = make([]E, y)
	}
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			result[i][j] = s[j][i]
		}
	}
	return result
}

// ArrayMap returns a new array containing for each key of the source array, the value returned by the callback function
func ArrayMap[S ~[]In, In any, Out any](input S, callback func(In) Out) []Out {
	output := make([]Out, len(input))
	for i, v := range input {
		output[i] = callback(v)
	}
	return output
}

// UnorderedArray implements a utility container allowing fast deletion without preserving items order
type UnorderedArray[E any] []E

// NewUnorderedArray is a quick way to get an UnorderedArray without worrying about specification and value assignation
func NewUnorderedArray[E any](v ...E) UnorderedArray[E] {
	return append(UnorderedArray[E]{}, v...)
}

// Values returns the container values
func (unorderedArray UnorderedArray[E]) Values() []E {
	return unorderedArray
}

// Insert appends elements to the end
func (unorderedArray *UnorderedArray[E]) Insert(v ...E) {
	*unorderedArray = append(*unorderedArray, v...)
}

// Remove removes element at index i by copying at this index the last one, then shrink the container
func (unorderedArray *UnorderedArray[E]) Remove(i int) (old E) {
	if l := len(*unorderedArray); i >= 0 && i < l {
		old = (*unorderedArray)[i]
		(*unorderedArray)[i] = (*unorderedArray)[l-1]
		*unorderedArray = (*unorderedArray)[:l-1]
	}
	return old
}

// OrderedArray implements a utility container allowing deletion preserving items order
type OrderedArray[E any] []E

// NewOrderedArray is a quick way to get an OrderedArray without worrying about specification and value assignation
func NewOrderedArray[E any](v ...E) OrderedArray[E] {
	return append(OrderedArray[E]{}, v...)
}

// Values returns the container values
func (orderedArray OrderedArray[E]) Values() []E {
	return orderedArray
}

// Insert appends elements to the end
func (orderedArray *OrderedArray[E]) Insert(v ...E) {
	*orderedArray = append(*orderedArray, v...)
}

// Remove removes element at index i
func (orderedArray *OrderedArray[E]) Remove(i int) (old E) {
	if l := len(*orderedArray); i >= 0 && i < l {
		old = (*orderedArray)[i]
		*orderedArray = append((*orderedArray)[:i], (*orderedArray)[i+1:]...)
	}
	return old
}

// CyclicArray implements a utility container allowing rotative access to elements
type CyclicArray[E any] struct {
	values []E
	index  int
}

// NewCyclicArray is a quick way to get an CyclicArray without worrying about specification and value assignation
func NewCyclicArray[E any](v ...E) CyclicArray[E] {
	return CyclicArray[E]{values: append([]E{}, v...)}
}

// Values returns the container values
func (cyclicArray CyclicArray[E]) Values() []E {
	return cyclicArray.values
}

// Index returns the current element index pointed
func (cyclicArray CyclicArray[E]) Index() int {
	return cyclicArray.index % len(cyclicArray.values)
}

// Get returns the value at index i and set the current index to this index
func (cyclicArray *CyclicArray[E]) Get(i int) E {
	cyclicArray.index = i % len(cyclicArray.values)
	return cyclicArray.values[cyclicArray.index]
}

// Next returns the value currently pointed, and move pointer to the next location
func (cyclicArray *CyclicArray[E]) Next() (value E) {
	value = cyclicArray.Get(cyclicArray.index)
	cyclicArray.index++
	return
}
