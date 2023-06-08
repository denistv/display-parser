package pipeline

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/PuerkitoBio/goquery"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
)

var modelParserPPIRegexp = regexp.MustCompile(`\d+ ppi`)

// NewModelParser Разбирает страницу с описанием монитора, сохраняя его свойства в сущность модели
func NewModelParser(logger *zap.Logger, modelsRepo repository.ModelRepository) *ModelParser {
	return &ModelParser{
		logger:     logger,
		modelsRepo: modelsRepo,
	}
}

type ModelParser struct {
	logger     *zap.Logger
	modelsRepo repository.ModelRepository
}

// Run запускает часть пайплайна, отвечающую за парсинг страниц.
func (m *ModelParser) Run(ctx context.Context, in <-chan domain.PageEntity) {
	go func() {
		for {
			select {
			case page := <-in:
				// TODO поправить N+1, не выглядит оптимально
				_, ok, err := m.modelsRepo.Find(context.Background(), page.URL)
				if err != nil {
					m.logger.Error(fmt.Errorf("find model for page: %w", err).Error())
					continue
				}

				if ok {
					// Пока не реализована возможность перезаписи сущестувющей модели.
					// В моем кейсе проще транкейтнуть таблицу (или удалить конкретные модели, которые нужно перепарсить) руками и перезапустить сервис
					m.logger.Info("model exists, skipping " + page.URL)
					continue
				}

				model, err := m.parse(page)
				if err != nil {
					m.logger.Error(fmt.Errorf("parsing page: %w", err).Error())
					continue
				}

				err = m.modelsRepo.Create(context.Background(), model)
				if err != nil {
					m.logger.Error(fmt.Errorf("creating model: %w", err).Error())
					continue
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Общий метод, где вызываем специализированные методы разбора нужных нам атрибутов у монитора.
// Здесь не держим какой-то конкретной логики по парсингу того или иного свойства монитора, только вызываем специализированыые методы,
// чтобы не смешивать все вместе превращая код в дробленку, сохраняя читаемость и прозрачность происходящего.
func (m *ModelParser) parse(page domain.PageEntity) (domain.ModelEntity, error) {
	ppiInt, err := m.parsePPI(page)
	if err != nil {
		return domain.ModelEntity{}, fmt.Errorf("parsing ppi: %w", err)
	}

	model := domain.ModelEntity{
		URL: page.URL,
		PPI: ppiInt,
	}

	buf := bytes.NewBufferString(page.Body)

	htmlDoc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return domain.ModelEntity{}, fmt.Errorf("creating html document from reader: %w", err)
	}

	m.parseBrandSeriesModel(htmlDoc, &model)
	m.parseDisplay(htmlDoc, &model)

	return model, nil
}

func (m *ModelParser) parsePPI(page domain.PageEntity) (int64, error) {
	ppiRaw := modelParserPPIRegexp.FindAllString(page.Body, 1)
	if len(ppiRaw) == 0 {
		return 0, errors.New("cannot find ppi value")
	}

	ppi := ppiRaw[0]
	ppi = strings.TrimSuffix(ppi, " ppi")

	ppiInt, err := strconv.ParseInt(ppi, 10, 64)
	if err != nil {
		return 0, errors.New("cannot parse ppi")
	}

	return ppiInt, nil
}

func (m *ModelParser) parseBrandSeriesModel(doc *goquery.Document, model *domain.ModelEntity) {
	doc.Find("#main > div:nth-child(6) > table:nth-child(2) > tbody > tr").
		Each(func(i int, s *goquery.Selection) {
			label := s.Find("td:nth-child(1)").Text()
			value := s.Find("td:nth-child(2)").Text()

			switch label {
			case "BrandName of the company-manufacturer.":
				model.Brand = value
			case "SeriesName of the series, which the model belongs to.":
				model.Series = value
			case "ModelDesignation of the model.":
				model.Name = value
			case "Model yearThe year in which this model was announced.":
				// не для всех моделей может существовать год, поэтому игнорируем ошибку парсинга и оставляем в этом случае дефолтное значение
				year, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					model.Year = 0
				} else {
					model.Year = year
				}
			}
		})
}

func (m *ModelParser) parseDisplay(doc *goquery.Document, model *domain.ModelEntity) {
	doc.Find("#main > div:nth-child(6) > table:nth-child(4) > tbody > tr").
		Each(func(i int, s *goquery.Selection) {
			label := s.Find("td:nth-child(1)").Text()
			value := s.Find("td:nth-child(2)").Text()

			//nolint
			const sizeExistsPattern = "Size classSize class of the display as declared by the manufacturer. Often this is the rounded value of the actual size of the diagonal in inches."

			//nolint:gocritic
			switch label {
			case sizeExistsPattern:
				sizeRaw := strings.TrimSuffix(value, "in (inches)")
				size, _ := strconv.ParseInt(sizeRaw, 10, 64)
				model.Size = size
			}
		})
}
