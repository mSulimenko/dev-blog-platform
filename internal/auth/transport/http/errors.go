package httphandler

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	ErrCodeValidation  = "VALIDATION_ERROR"
	ErrCodeNotFound    = "NOT_FOUND"
	ErrCodeInternal    = "INTERNAL_ERROR"
	ErrCodeConflict    = "CONFLICT"
	ErrCodeInvalidJSON = "INVALID_JSON"
)

func (h *Handler) sendError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}
