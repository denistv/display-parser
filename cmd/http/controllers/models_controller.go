package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

// parseModelQuery создает экземпляр структуры, заполняя ее данными из параметров запроса
// Прячем в подобных конструкторах всю подобную внутрянку в пакете `controllers`.
// Таким образом, слой с http-стаффом и слой с репозиториями разделены и не мешают друг другу.
// Более высокоуровневый пакет содержащий репозитории БД ничего не знает про низкоуровневый пакет с HTTP-контроллераи,
// но пакет с http-контроллерами использует структуры пакета репозиториев и знает как заполнить его структуры
// запросов данными.
func parseModelQuery(r *http.Request) (repository.ModelQuery, error) {
	var err error

	mq := repository.NewModelQuery()

	if v := r.URL.Query().Get("limit"); v != "" {
		if mq.Limit.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse limit: %w", err)
		}
		mq.Limit.Valid = true
	}

	if v := r.URL.Query().Get("brand"); v != "" {
		mq.Brand.String = v
		mq.Brand.Valid = true
	}

	if v := r.URL.Query().Get("ppi-from"); v != "" {
		if mq.PPIFrom.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse ppi-from: %w", err)
		}
		mq.PPIFrom.Valid = true
	}
	if v := r.URL.Query().Get("ppi-to"); v != "" {
		if mq.PPITo.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse ppi-to: %w", err)
		}
		mq.PPITo.Valid = true
	}

	if v := r.URL.Query().Get("year-from"); v != "" {
		if mq.YearFrom.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse year: %w", err)
		}
		mq.YearFrom.Valid = true
	}
	if v := r.URL.Query().Get("year-to"); v != "" {
		if mq.YearTo.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse year-to: %w", err)
		}
		mq.YearTo.Valid = true
	}

	if v := r.URL.Query().Get("size-from"); v != "" {
		if mq.SizeFrom.Float64, err = strconv.ParseFloat(v, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse size-from: %w", err)
		}
		mq.SizeFrom.Valid = true
	}
	if v := r.URL.Query().Get("size-to"); v != "" {
		if mq.SizeTo.Float64, err = strconv.ParseFloat(v, 64); err != nil {
			return repository.ModelQuery{}, fmt.Errorf("parse size-to: %w", err)
		}
		mq.SizeTo.Valid = true
	}

	return mq, nil
}

func (m *ModelsController) ModelsIndex(w http.ResponseWriter, r *http.Request) {
	// todo populate params
	mq, err := parseModelQuery(r)
	if err != nil {
		m.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err = mq.Validate(); err != nil {
		m.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	all, err := m.repo.All(r.Context(), mq)
	if err != nil {
		m.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(all)
	if err != nil {
		m.logger.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, _ = w.Write(data)
	w.WriteHeader(http.StatusOK)
}
