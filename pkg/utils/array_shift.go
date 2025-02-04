package utils

func ArrayShift[T interface{}](s []T) *T {
	if len(s) == 0 {
		return nil
	}

	t := s[0]

	s = s[1:]

	return &t
}
