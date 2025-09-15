package gormite_query_builders

import "golang.org/x/exp/constraints"

// CloneTo - Clones a query builder to a new query builder.
func CloneTo[New any, Old any](oldQb *QueryBuilder[Old]) *QueryBuilder[New] {
	newQb := QueryBuilder[New](*oldQb)

	return newQb.Clone()
}

// AsPrimitiveStrSlice - Converts a slice of strings to a slice of primitive strings.
// Helpful for parameters binding.
func AsPrimitiveStrSlice[T ~string](values []T) []string {
	casted := make([]string, len(values))
	for i, value := range values {
		casted[i] = string(value)
	}
	return casted
}

// AsPrimitiveIntSlice - Converts a slice of integers to a slice of primitive integers.
// Helpful for parameters binding.
func AsPrimitiveIntSlice[T constraints.Integer](values []T) []int {
	casted := make([]int, len(values))
	for i, value := range values {
		casted[i] = int(value)
	}
	return casted
}
