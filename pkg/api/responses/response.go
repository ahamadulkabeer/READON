package responses

func ClientReponse(statusCode int, message string, err any, data any) Response {
	return Response{
		StatusCode: statusCode,
		Message:    message,
		Error:      err,
		Data:       data,
	}
}

type Response struct {
	StatusCode int
	Message    string
	Error      interface{}
	Data       interface{}
}
