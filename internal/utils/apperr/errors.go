package apperr

func NewInternal(parent error, format string, a ...interface{}) *AppError {
	return newAppError(ErrorTypeInternal, parent, format, a...)
}

func NewBadRequest(parent error, format string, a ...interface{}) *AppError {
	return newAppError(ErrorTypeBadRequest, parent, format, a...)
}

func NewNotFound(parent error, format string, a ...interface{}) *AppError {
	return newAppError(ErrorTypeNotFound, parent, format, a...)
}

func NewUnauthorized(parent error, format string, a ...interface{}) *AppError {
	return newAppError(ErrorTypeUnauthorized, parent, format, a...)
}

func NewInvalidArgument(parent error, format string, a ...interface{}) *AppError {
	return newAppError(ErrorTypeInvalidArgument, parent, format, a...)
}

func NewAlreadyExists(parent error, format string, a ...interface{}) *AppError {
	return newAppError(ErrorTypeAlreadyExists, parent, format, a...)
}
