package models

type ErrorResponseK struct {
	Status  int
	Message string
	Error   interface{}
}

type SuccessResponse struct {
	Status  int
	Message string
	Data    interface{}
}
