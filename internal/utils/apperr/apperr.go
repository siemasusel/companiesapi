package apperr

import "fmt"

type ErrorType string

var (
	ErrorTypeInternal        ErrorType = "internal"
	ErrorTypeBadRequest      ErrorType = "bad-request"
	ErrorTypeInvalidArgument ErrorType = "invalid-argument"
	ErrorTypeNotFound        ErrorType = "not-found"
	ErrorTypeUnauthorized    ErrorType = "unauthrozied"
	ErrorTypeAlreadyExists   ErrorType = "already-exists"
)

type AppError struct {
	Type    ErrorType
	Parent  error
	Message string
}

func newAppError(typ ErrorType, parent error, format string, a ...interface{}) *AppError {
	return &AppError{
		Type:    typ,
		Parent:  parent,
		Message: fmt.Sprintf(format, a...),
	}
}

func (err AppError) Error() string {
	msg := fmt.Sprintf("Message=(%s)", err.Message)

	if err.Parent != nil {
		msg += fmt.Sprintf(" Parent=(%v)", err.Parent)
	}

	return msg
}

func (err AppError) Unwrap() error {
	return err.Parent
}
