package errorhandler

type ErrorData struct {
	Code         int
	Message      string
	ErrorMessage error
}

func (e ErrorData) Error() string {
	return e.Message
}
