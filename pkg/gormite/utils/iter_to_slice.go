package utils

import "iter"

func IterToSlice[T any](iter iter.Seq[T]) []T {
	var result []T
	for v := range iter {
		result = append(result, v)
	}
	return result
}
