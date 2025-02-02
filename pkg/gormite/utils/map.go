package utils

func Map[T1 comparable, T2, V any](ts map[T1]T2, fn func(T1, T2) V) map[T1]V {
	result := make(map[T1]V)
	for k, v := range ts {
		result[k] = fn(k, v)
	}
	return result
}
