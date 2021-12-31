package utils

import (
	"maps"
)

// CopyMap returns a new map containing copy of elements created by calling V.Copy()
func CopyMap[M ~map[K]V, K comparable, V any](input M) M {
	return maps.Clone(input)
}

// MapKeys return a new array containing all map keys
func MapKeys[M ~map[K]V, K comparable, V any](input M) []K {
	keys := make([]K, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

// MapValues return a new array containing all map values (no check for duplicates)
func MapValues[M ~map[K]V, K comparable, V any](input M) []V {
	values := make([]V, 0, len(input))
	for _, v := range input {
		values = append(values, v)
	}
	return values
}

// MapFilter returns a new map containing all the elements for which the condition returned true
func MapFilter[M ~map[K]V, K comparable, V any](input M, cond func(K, V) bool) M {
	output := make(M, len(input))
	for k, v := range input {
		if cond(k, v) {
			output[k] = v
		}
	}
	return output
}

// MapMap returns a new map containing for each key of the source map, the value returned by the callback function
func MapMap[M1 ~map[K1]V1, K1 comparable, V1 any, M2 map[K2]V2, K2 comparable, V2 any](input M1, callback func(K1, V1) (K2, V2)) M2 {
	output := make(M2, len(input))
	for ki, vi := range input {
		ko, vo := callback(ki, vi)
		output[ko] = vo
	}
	return output
}
