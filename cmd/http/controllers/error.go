package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"display_parser/internal/domain"
)

func NewParseError(s string) error {
	return fmt.Errorf("%w: %s", ErrParseError, s)
}

var ErrParseError = errors.New("parse error")

func newErrorResponse(err string) errorResponse {
	return errorResponse{Error: err}
}

type errorResponse struct {
	Error string `json:"error"`
}

func is400Err(err error) bool {
	return errors.Is(err, domain.ErrValidationError) || errors.Is(err, ErrParseError)
}

// handleError решает, как будет сформирована ошибка и отправлена клиенту.
// Если ошибка относится к 400 http-кодам, то будет взято ее текстовое описание и отправлено клиенту.
// Другие ошибки будут отправлены клиенту как обезличенные "internal server error", чтобы не светить лишнего в целях безопасности.
func handleError(w http.ResponseWriter, err error) {
	if is400Err(err) {
		writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func writeJSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	e := newErrorResponse(err)
	_ = json.NewEncoder(w).Encode(e)
}
