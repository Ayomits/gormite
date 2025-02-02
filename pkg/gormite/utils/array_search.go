package utils

func ArraySearch[T comparable](needle T, haystack []T) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}
