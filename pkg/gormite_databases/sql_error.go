package gormite_databases

type SqlError struct {
	err  error
	sql  string
	args []any
}

func (e *SqlError) Error() string {
	return e.err.Error()
}

func (e *SqlError) GetExtraData() map[string]any {
	return map[string]any{
		"sql":  e.sql,
		"args": e.args,
	}
}
