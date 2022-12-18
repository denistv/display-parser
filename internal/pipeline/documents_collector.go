package pipeline

import (
	"context"
	"displayCrawler/internal/domain"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"

	"displayCrawler/internal/repository"
)

func NewDocumentsCollector(logger *zap.Logger, docRepo *repository.Document) *DocumentsCollector {
	return &DocumentsCollector{
		logger:  logger,
		docRepo: docRepo,
	}
}

// Слушает канал с URL моделей устройств и для каждого URL загружает документ с описанием модели
type DocumentsCollector struct {
	logger  *zap.Logger
	docRepo *repository.Document
}

func (d *DocumentsCollector) Run(in <-chan string) chan domain.DocumentEntity {
	out := make(chan domain.DocumentEntity)

	go func() {
		for docURL := range in {
			doc, isExists, err := d.docRepo.Find(context.Background(), docURL)
			if err != nil {
				d.logger.Error("checking model is exists: " + err.Error())

				continue
			}

			if !isExists {
				body, err := d.download(docURL)
				if err != nil {
					d.logger.Error(fmt.Errorf("downloading document: %w", err).Error())

					continue
				}

				doc = domain.DocumentEntity{
					URL:  docURL,
					Body: body,
				}

				time.Sleep(50 * time.Millisecond)
			}

			out <- doc
		}

		close(out)
	}()

	return out
}

func (d *DocumentsCollector) download(docURL string) (string, error) {
	if docURL == "" {
		return "", errors.New("url cannot be empty")
	}

	res, err := http.Get(docURL)
	if err != nil {
		return "", fmt.Errorf("getting model: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("non-200 status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading document body: %w", err)
	}

	return string(body), nil
}
