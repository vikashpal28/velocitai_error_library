package myerrors

import "net/http"

type Code string

const (
	CodeUnknown      Code = "UNKNOWN"
	CodeValidation   Code = "VALIDATION_FAILED"
	CodeNotFound     Code = "NOT_FOUND"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeForbidden    Code = "FORBIDDEN"
	CodeConflict     Code = "CONFLICT"
	CodeRateLimit    Code = "RATE_LIMIT"
	CodeInternal     Code = "INTERNAL_ERROR"
	CodeUnavailable  Code = "SERVICE_UNAVAILABLE"
	CodeTimeout      Code = "TIMEOUT"
	CodeBadRequest   Code = "BAD_REQUEST"
)

var statusByCode = map[Code]int{
	CodeUnknown:      http.StatusInternalServerError,
	CodeValidation:   http.StatusBadRequest,
	CodeNotFound:     http.StatusNotFound,
	CodeUnauthorized: http.StatusUnauthorized,
	CodeForbidden:    http.StatusForbidden,
	CodeConflict:     http.StatusConflict,
	CodeRateLimit:    http.StatusTooManyRequests,
	CodeInternal:     http.StatusInternalServerError,
	CodeUnavailable:  http.StatusServiceUnavailable,
	CodeTimeout:      http.StatusRequestTimeout,
	CodeBadRequest:   http.StatusBadRequest,
}

func (c Code) HTTPStatus() int {
	if status , ok := statusByCode[c]; ok {
		return status
	}
	return http.StatusInternalServerError
}	
