package pipeline

import (
	"context"
	"display_parser/internal/services"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
)

func NewPagesCollector(logger *zap.Logger, docRepo *repository.Page, httpClient services.HTTPClient) *PagesCollector {
	return &PagesCollector{
		logger:     logger,
		docRepo:    docRepo,
		httpClient: httpClient,
	}
}

// Слушает канал с URL моделей устройств и для каждого URL загружает документ с описанием модели
type PagesCollector struct {
	logger     *zap.Logger
	docRepo    *repository.Page
	httpClient services.HTTPClient
}

func (d *PagesCollector) Run(ctx context.Context, in <-chan string) chan domain.PageEntity {
	out := make(chan domain.PageEntity)

	go func() {
		defer close(out)

		for {
			select {
			case docURL, ok := <-in:
				if !ok {
					return
				}

				doc, isExists, err := d.docRepo.Find(ctx, docURL)
				if err != nil {
					d.logger.Error("checking model is exists: " + err.Error())

					continue
				}

				if !isExists {
					body, err := d.download(ctx, docURL)
					if err != nil {
						d.logger.Error(fmt.Errorf("downloading document: %w", err).Error())

						continue
					}

					doc = domain.PageEntity{
						URL:  docURL,
						Body: body,
					}
				}

				out <- doc
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (d *PagesCollector) download(ctx context.Context, pageURL string) (string, error) {
	if pageURL == "" {
		return "", errors.New("url cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("creating http req: %w", err)
	}

	res, err := d.httpClient.Do(req)
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
