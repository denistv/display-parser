package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
)

func NewPagesCollector(logger *zap.Logger, docRepo *repository.Page) *PagesCollector {
	return &PagesCollector{
		logger:  logger,
		docRepo: docRepo,
	}
}

// Слушает канал с URL моделей устройств и для каждого URL загружает документ с описанием модели
type PagesCollector struct {
	logger  *zap.Logger
	docRepo *repository.Page
}

func (d *PagesCollector) Run(in <-chan string) chan domain.PageEntity {
	out := make(chan domain.PageEntity)

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

				doc = domain.PageEntity{
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

func (d *PagesCollector) download(pageURL string) (string, error) {
	if pageURL == "" {
		return "", errors.New("url cannot be empty")
	}

	res, err := http.Get(pageURL)
	if err != nil {
		return "", fmt.Errorf("getting model: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("non-200 status code")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading document body: %w", err)
	}

	return string(body), nil
}
