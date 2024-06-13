package responses

import (
	"readon/pkg/models"
)

func RespondWithError(status int, message string, err interface{}) models.ErrorResponseK {
	return models.ErrorResponseK{
		Status:  status,
		Message: message,
		Error:   err,
	}

}

func RespondWithSuccess(status int, message string, data interface{}) models.SuccessResponse {
	return models.SuccessResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

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
