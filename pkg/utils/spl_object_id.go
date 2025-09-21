package utils

import "fmt"

func SplObjectID(v any) string {
	return fmt.Sprintf("%p", v)
}
