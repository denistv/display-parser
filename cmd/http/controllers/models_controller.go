package controllers

import (
	"display_parser/internal/repository"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func NewModelsController(logger *zap.Logger, repo repository.ModelRepository) *ModelsController {
	return &ModelsController{
		logger: logger,
		repo: repo,
	}
}

type ModelsController struct {
	logger *zap.Logger
	repo repository.ModelRepository
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

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

