package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func newJSONController(logger *zap.Logger) jsonController {
	return jsonController{logger: logger}
}

type jsonController struct {
	logger *zap.Logger
}

func (j *jsonController) writeJSONResponse(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		j.logger.Error(fmt.Errorf("error while encoding json response: %w", err).Error())
		return
	}
}

func newErrorResponse(err string) errorResponse {
	return errorResponse{Error: err}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (j *jsonController) writeJSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	e := newErrorResponse(err)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		j.logger.Error(fmt.Errorf("error while encoding json response: %w", err).Error())
		return
	}
}

// handleError решает, как будет сформирована ошибка и отправлена клиенту.
// Если ошибка относится к 400 http-кодам, то будет взято ее текстовое описание и отправлено клиенту.
// Другие ошибки будут отправлены клиенту как обезличенные "internal server error", чтобы не светить лишнего в целях безопасности.
func (j *jsonController) handleError(w http.ResponseWriter, err error) {
	j.logger.Error(err.Error())

	if is400Err(err) {
		j.writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	j.writeJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
