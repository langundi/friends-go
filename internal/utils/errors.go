package utils

import (
	"log/slog"
	"net/http"
)

// 400
func BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("bad request", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusBadRequest, err)
}

func InvalidPayloadError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("invalid request payload", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusBadRequest, err)
}

// 401
func UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("unauthorized", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusUnauthorized, err)
}

// 403
func ForbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("forbidden", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusForbidden, err)
}

// 404
func NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("not found", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusNotFound, err)
}

// 409
func ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("conflict", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusConflict, err)
}

// 429
func TooManyRequestsError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Warn("too many requests", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusTooManyRequests, err)
}

// 500
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal server", "error", err, "path", r.URL.Path)
	WriteError(w, http.StatusInternalServerError, err)
}
