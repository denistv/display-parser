package pipeline

import (
	"bytes"
	"context"
	"displayCrawler/internal/domain"
	"displayCrawler/internal/repository"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
)

// Разбирает документ и парсит его описание, вытаскивая полезную иформацию, сохраняя ее в сущность модели

func NewModelParser(logger *zap.Logger, modelsRepo *repository.Model) *ModelParser {
	ppiRe := regexp.MustCompile("[0-9]+ ppi")

	return &ModelParser{
		logger:     logger,
		modelsRepo: modelsRepo,
		ppiRe:      ppiRe,
	}
}

type ModelParser struct {
	logger     *zap.Logger
	modelsRepo *repository.Model

	ppiRe *regexp.Regexp
}

func (m *ModelParser) Run(in <-chan domain.DocumentEntity) {
	go func() {
		for document := range in {
			model, ok, err := m.modelsRepo.Find(context.Background(), document.URL)
			if err != nil {
				m.logger.Error(fmt.Errorf("find model for document: %w", err).Error())
				continue
			}

			if ok {
				m.logger.Info("model exists, skipping " + document.URL)
				continue
			}

			model, err = m.parse(document)
			if err != nil {
				m.logger.Error(fmt.Errorf("parsing document: %w", err).Error())
				continue
			}

			err = m.modelsRepo.Create(context.Background(), model)
			if err != nil {
				m.logger.Error(fmt.Errorf("creating model: %w", err).Error())
				continue
			}
		}
	}()
}

func (m *ModelParser) parse(doc domain.DocumentEntity) (domain.ModelEntity, error) {
	ppiInt, err := m.parsePPI(doc)
	if err != nil {
		return domain.ModelEntity{}, fmt.Errorf("parsing ppi: %w", err)
	}

	model := domain.ModelEntity{
		URL: doc.URL,
		PPI: ppiInt,
	}

	buf := bytes.NewBufferString(doc.Body)

	htmlDoc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return domain.ModelEntity{}, fmt.Errorf("creating doc from reader: %w", err)
	}

	err = m.parseBrandSeriesModel(htmlDoc, &model)
	if err != nil {
		return domain.ModelEntity{}, fmt.Errorf("parsing model: %w", err)
	}

	err = m.parseDisplay(htmlDoc, &model)
	if err != nil {
		return domain.ModelEntity{}, fmt.Errorf("parsing model: %w", err)
	}

	return model, nil
}

func (m *ModelParser) parsePPI(doc domain.DocumentEntity) (int64, error) {
	ppiRaw := m.ppiRe.FindAllString(doc.Body, 1)
	if len(ppiRaw) == 0 {
		return 0, errors.New("cannot find ppi value")
	}

	ppi := ppiRaw[0]
	ppi = strings.Trim(ppi, " ppi")

	ppiInt, err := strconv.ParseInt(ppi, 10, 64)
	if err != nil {
		return 0, errors.New("cannot parse ppi")
	}

	return ppiInt, nil
}

func (m *ModelParser) parseBrandSeriesModel(doc *goquery.Document, model *domain.ModelEntity) error {
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
				year, _ := strconv.ParseInt(value, 10, 64)
				model.Year = year
			}
		})

	return nil
}

func (m *ModelParser) parseDisplay(doc *goquery.Document, model *domain.ModelEntity) error {
	doc.Find("#main > div:nth-child(6) > table:nth-child(4) > tbody > tr").
		Each(func(i int, s *goquery.Selection) {
			label := s.Find("td:nth-child(1)").Text()
			value := s.Find("td:nth-child(2)").Text()

			switch label {
			case "Size classSize class of the display as declared by the manufacturer. Often this is the rounded value of the actual size of the diagonal in inches.":
				sizeRaw := strings.Trim(value, "in (inches)")
				size, _ := strconv.ParseInt(sizeRaw, 10, 64)
				model.Size = size
			}
		})

	return nil
}
