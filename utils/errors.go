package utils

import (
	"errors"
	"net/http"
)

func PermissionDenied(w http.ResponseWriter) {
	WriteError(w, http.StatusForbidden, errors.New("permission denied"))
}

func BadRequest(w http.ResponseWriter) {
	WriteError(w, http.StatusBadRequest, errors.New("bad request"))
}

func InternalServerError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, errors.New("internal server error"))
}

func NotFound(w http.ResponseWriter) {
	WriteError(w, http.StatusNotFound, errors.New("not found"))
}

func Unauthorized(w http.ResponseWriter) {
	WriteError(w, http.StatusUnauthorized, errors.New("unauthorized"))
}

func MethodNotAllowed(w http.ResponseWriter) {
	WriteError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
}

func NotImplemented(w http.ResponseWriter) {
	WriteError(w, http.StatusNotImplemented, errors.New("not implemented"))
}

func ServiceUnavailable(w http.ResponseWriter) {
	WriteError(w, http.StatusServiceUnavailable, errors.New("service unavailable"))
}

func GatewayTimeout(w http.ResponseWriter) {
	WriteError(w, http.StatusGatewayTimeout, errors.New("gateway timeout"))
}

func RateLimited(w http.ResponseWriter) {
	WriteError(w, http.StatusTooManyRequests, errors.New("rate limited"))
}
