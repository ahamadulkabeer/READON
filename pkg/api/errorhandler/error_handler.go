package errorhandler

import (
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ErrorData struct {
	Code         int
	Message      string
	ErrorMessage error
}

func (e ErrorData) Error() string {
	return e.Message
}

func HandleDatabaseError(err error) (int, error) {
	// log the database error
	logError(err)

	switch {
	// 4xx Client Errors
	case errors.Is(err, gorm.ErrRecordNotFound):
		return 404, errors.New("resource not found")
	case errors.Is(err, gorm.ErrInvalidData):
		return 400, errors.New("invalid data provided. please check your input and try again")
	case errors.Is(err, gorm.ErrMissingWhereClause):
		return 400, errors.New("a where clause is required for this operation")
	case errors.Is(err, logger.ErrRecordNotFound):
		return 404, errors.New("log record not found")

	// 5xx Server Errors
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return 500, errors.New("internal server error. please try again later")
	case errors.Is(err, gorm.ErrUnsupportedDriver):
		return 500, errors.New("database driver is not supported")
	case errors.Is(err, gorm.ErrRegistered):
		return 500, errors.New("resource is already registered")
	case errors.Is(err, gorm.ErrNotImplemented):
		return 501, errors.New("functionality is not implemented yet")

	// Default case for any other errors
	default:
		return 500, errors.New("unexpected internal server error. please try again later")
	}
}

func logError(err error) {
	log.Printf("Error occurred: %v", err)
}

// func logError(err error) {
// 	log.WithFields(logrus.Fields{
// 		"error": err,
// 	}).Error("Database error occurred")
// }
