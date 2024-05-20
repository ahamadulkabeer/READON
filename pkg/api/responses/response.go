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

func ClientReponse(status int, message string, err error, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Error:   err,
		Data:    data,
	}
}

type Response struct {
	Status  int
	Message string
	Error   error
	Data    interface{}
}
