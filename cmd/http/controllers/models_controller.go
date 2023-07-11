package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"display_parser/internal/repository"
)

func NewModelsController(logger *zap.Logger, repo repository.ModelRepository) *ModelsController {
	j := newJSONController(logger)

	return &ModelsController{
		jsonController: j,
		logger:         logger,
		repo:           repo,
	}
}

type ModelsController struct {
	jsonController

	logger *zap.Logger
	repo   repository.ModelRepository
}

// parseModelQuery создает экземпляр структуры, заполняя ее данными из параметров запроса
// Прячем в подобных конструкторах всю подобную внутрянку в пакете `controllers`.
// Таким образом, слой с http-стаффом и слой с репозиториями разделены и не мешают друг другу.
// Более высокоуровневый пакет содержащий репозитории БД ничего не знает про низкоуровневый пакет с HTTP-контроллераи,
// но пакет с http-контроллерами использует структуры пакета репозиториев и знает как заполнить его структуры
// запросов данными.
//
//nolint:gocyclo
func parseModelQuery(r *http.Request) (repository.ModelQuery, error) {
	var err error

	mq := repository.NewModelQuery()

	if v := r.URL.Query().Get("limit"); v != "" {
		if mq.Limit.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse limit: %w", err).Error())
		}
		mq.Limit.Valid = true
	}

	if v := r.URL.Query().Get("brand"); v != "" {
		mq.Brand.String = v
		mq.Brand.Valid = true
	}

	if v := r.URL.Query().Get("ppi-from"); v != "" {
		if mq.PPIFrom.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse ppi-from: %w", err).Error())
		}
		mq.PPIFrom.Valid = true
	}
	if v := r.URL.Query().Get("ppi-to"); v != "" {
		if mq.PPITo.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse ppi-to: %w", err).Error())
		}
		mq.PPITo.Valid = true
	}

	if v := r.URL.Query().Get("year-from"); v != "" {
		if mq.YearFrom.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse year: %w", err).Error())
		}
		mq.YearFrom.Valid = true
	}
	if v := r.URL.Query().Get("year-to"); v != "" {
		if mq.YearTo.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse year-to: %w", err).Error())
		}
		mq.YearTo.Valid = true
	}

	if v := r.URL.Query().Get("size-from"); v != "" {
		if mq.SizeFrom.Float64, err = strconv.ParseFloat(v, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse size-from: %w", err).Error())
		}
		mq.SizeFrom.Valid = true
	}
	if v := r.URL.Query().Get("size-to"); v != "" {
		if mq.SizeTo.Float64, err = strconv.ParseFloat(v, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse size-to: %w", err).Error())
		}
		mq.SizeTo.Valid = true
	}

	if v := r.URL.Query().Get("panel-bit-depth"); v != "" {
		if mq.PanelBitDepth.Int64, err = strconv.ParseInt(v, 10, 64); err != nil {
			return repository.ModelQuery{}, NewParseError(fmt.Errorf("parse panel-bit-depth: %w", err).Error())
		}
		mq.PanelBitDepth.Valid = true
	}

	return mq, nil
}

func (m *ModelsController) ModelsIndex(w http.ResponseWriter, r *http.Request) {
	mq, err := parseModelQuery(r)
	if err != nil {
		m.handleError(w, err)
		return
	}

	if err = mq.Validate(); err != nil {
		m.handleError(w, err)
		return
	}

	all, err := m.repo.All(r.Context(), mq)
	if err != nil {
		m.handleError(w, err)
		return
	}

	m.writeJSONResponse(w, all)
}
