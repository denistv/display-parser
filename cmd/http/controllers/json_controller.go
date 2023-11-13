package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"display_parser/pkg/logger"
)

func newJSONController(l logger.Logger) jsonController {
	return jsonController{logger: l}
}

// jsonController можно встроить в структуру контроллера для придания ему нужного поведения
type jsonController struct {
	logger logger.Logger
}

func (j *jsonController) writeHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

// handleError решает, как будет сформирована ошибка и отправлена клиенту.
// Если ошибка относится к 400 http-кодам, то будет взято ее текстовое описание и отправлено клиенту.
// Другие ошибки будут отправлены клиенту как обезличенные "internal server error", чтобы не светить лишнего в целях безопасности.
func (j *jsonController) handleError(w http.ResponseWriter, err error) {
	if is400Err(err) {
		j.logger.Debug(fmt.Errorf("bad request: %w", err).Error())
		j.writeJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	j.logger.Error(err.Error())
	j.writeJSONError(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func newErrorResponse(err string) errorResponse {
	return errorResponse{Error: err}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (j *jsonController) writeJSONError(w http.ResponseWriter, err string, code int) {
	j.writeHeaders(w)
	w.WriteHeader(code)

	e := newErrorResponse(err)
	if err := json.NewEncoder(w).Encode(e); err != nil {
		j.logger.Error(fmt.Errorf("error while encoding json response: %w", err).Error())
		return
	}
}

func (j *jsonController) writeJSON(w http.ResponseWriter, v any) {
	j.writeHeaders(w)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		j.logger.Error(fmt.Errorf("error while encoding json response: %w", err).Error())
		return
	}
}
