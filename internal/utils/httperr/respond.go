package httperr

import (
	"errors"
	"net/http"

	"github.com/siemasusel/companiesapi/internal/utils/apperr"
	"golang.org/x/exp/slog"

	"github.com/go-chi/render"
)

func RespondWithError(err error, w http.ResponseWriter, r *http.Request) {
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) {
		RespondWithError(apperr.NewInternal(err, "internal server error"), w, r)
		return
	}

	switch appErr.Type {
	case apperr.ErrorTypeBadRequest:
		httpRespondWithError(appErr, appErr.Message, w, r, http.StatusBadRequest)
	case apperr.ErrorTypeInvalidArgument:
		httpRespondWithError(appErr, appErr.Message, w, r, http.StatusUnprocessableEntity)
	case apperr.ErrorTypeAlreadyExists:
		httpRespondWithError(appErr, appErr.Message, w, r, http.StatusConflict)
	case apperr.ErrorTypeNotFound:
		httpRespondWithError(appErr, appErr.Message, w, r, http.StatusNotFound)
	case apperr.ErrorTypeUnauthorized:
		httpRespondWithError(appErr, appErr.Message, w, r, http.StatusUnauthorized)
	default:
		httpRespondWithError(appErr, appErr.Message, w, r, http.StatusInternalServerError)
	}
}

func httpRespondWithError(err error, msg string, w http.ResponseWriter, r *http.Request, status int) {
	slog.ErrorContext(r.Context(), "HTTP error", "err", err.Error(), "msg", msg, "status", status)

	render.Status(r, status)
	render.Respond(w, r, map[string]string{"error": msg})
}
