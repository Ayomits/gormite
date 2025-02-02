package g_err

type QueryException struct {
	message string
}

func (e *QueryException) Error() string {
	return e.message
}

func NewQueryException(message string) *QueryException {
	return &QueryException{
		message: message,
	}
}
