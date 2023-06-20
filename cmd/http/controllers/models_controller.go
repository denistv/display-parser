package controllers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"display_parser/internal/repository"
)

func NewModelsController(logger *zap.Logger, repo repository.ModelRepository) *ModelsController {
	return &ModelsController{
		logger: logger,
		repo:   repo,
	}
}

type ModelsController struct {
	logger *zap.Logger
	repo   repository.ModelRepository
}

func (m *ModelsController) ModelsIndex(w http.ResponseWriter, r *http.Request) {
	all, err := m.repo.All(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(all)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(data)
	w.WriteHeader(http.StatusOK)
}
