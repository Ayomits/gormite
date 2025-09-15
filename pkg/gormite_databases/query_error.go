package gormite_databases

type queryError struct {
	err   error
	query *Query
}

func (e *queryError) Error() string {
	return e.err.Error()
}
