package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
	"display_parser/internal/services"
)

func NewPagesCollector(logger *zap.Logger, docRepo *repository.Page, httpClient services.HTTPClient, useStoredPagesOnly bool) *PagesCollector {
	return &PagesCollector{
		logger:             logger,
		pageRepo:           docRepo,
		httpClient:         httpClient,
		useStoredPagesOnly: useStoredPagesOnly,
	}
}

// Слушает канал с URL моделей устройств и для каждого URL загружает документ с описанием модели
type PagesCollector struct {
	logger             *zap.Logger
	pageRepo           *repository.Page
	httpClient         services.HTTPClient
	useStoredPagesOnly bool
}

func (d *PagesCollector) Run(ctx context.Context, in <-chan string) chan domain.PageEntity {
	out := make(chan domain.PageEntity)

	go func() {
		defer close(out)

		for {
			select {
			case pageURL, ok := <-in:
				if !ok {
					return
				}

				page, isExists, err := d.pageRepo.Find(ctx, pageURL)
				if err != nil {
					d.logger.Error("checking model is exists: " + err.Error())

					continue
				}

				if !d.useStoredPagesOnly && !isExists {
					body, err := d.download(ctx, pageURL)
					if err != nil {
						d.logger.Error(fmt.Errorf("downloading document: %w", err).Error())

						continue
					}

					page = domain.PageEntity{
						URL:  pageURL,
						Body: body,
					}

					err = d.pageRepo.Create(ctx, page)
					if err != nil {
						d.logger.Error(fmt.Errorf("creating page for %s: %w", pageURL, err).Error())
						continue
					}
				}

				out <- page
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

	d.logger.Debug(fmt.Sprintf("getting model page for %s", pageURL))

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
