package utils

func FlipSlice[V comparable](data []V) map[V]int {
	result := make(map[V]int, len(data))

	for i, v := range data {
		result[v] = i
	}

	return result
}
