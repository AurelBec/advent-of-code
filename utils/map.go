package utils

// CopyMap returns a new map containing copy of elements created by calling V.Copy()
func CopyMap[K comparable, V any](input map[K]V) map[K]V {
	output := make(map[K]V, len(input))
	for k, v := range input {
		output[k] = v
	}
	return output
}

// MapKeys return a new array containing all map keys
func MapKeys[K comparable, V any](input map[K]V) []K {
	keys := make([]K, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	return keys
}

// MapValues return a new array containing all map values (no check for duplicates)
func MapValues[K comparable, V any](input map[K]V) []V {
	values := make([]V, 0, len(input))
	for _, v := range input {
		values = append(values, v)
	}
	return values
}

// MapFilter returns a new map containing all the elements for which the condition returned true
func MapFilter[K comparable, V any](input map[K]V, cond func(K, V) bool) map[K]V {
	output := make(map[K]V, len(input))
	for k, v := range input {
		if cond(k, v) {
			output[k] = v
		}
	}
	return output
}

// MapMap returns a new map containing for each key of the source map, the value returned by the callback function
func MapMap[Ki comparable, Vi any, Ko comparable, Vo any](input map[Ki]Vi, callback func(Ki, Vi) (Ko, Vo)) map[Ko]Vo {
	output := make(map[Ko]Vo, len(input))
	for ki, vi := range input {
		ko, vo := callback(ki, vi)
		output[ko] = vo
	}
	return output
}
