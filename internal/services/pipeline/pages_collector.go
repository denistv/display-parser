package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/denistv/wdlogger"

	"display_parser/internal/config"
	"display_parser/internal/domain"
	"display_parser/internal/iface"
	"display_parser/internal/iface/db"
)

func NewPageCollector(l logger.Logger, pageRepo db.PageRepository, httpClient iface.HTTPClient, cfg config.PagesCollector) *PageCollector {
	return &PageCollector{
		logger:     l,
		pageRepo:   pageRepo,
		httpClient: httpClient,
		cfg:        cfg,
	}
}

// PageCollector Слушает канал с URL моделей устройств и для каждого URL загружает документ с описанием модели
type PageCollector struct {
	logger     logger.Logger
	pageRepo   db.PageRepository
	httpClient iface.HTTPClient
	cfg        config.PagesCollector
}

func (d *PageCollector) Run(ctx context.Context, in <-chan string) <-chan domain.PageEntity {
	out := make(chan domain.PageEntity, d.cfg.Count)

	go func() {
		defer close(out)

		for {
			select {
			case pageURL, ok := <-in:
				if !ok {
					return
				}

				entityID, err := domain.NewEntityIDFromURL(pageURL)
				if err != nil {
					d.logger.Error(fmt.Errorf("parsing entity ID: %w", err).Error())
					continue
				}

				page, isExists, err := d.pageRepo.Find(ctx, entityID)
				if err != nil {
					d.logger.Error("checking model is exists: " + err.Error())

					continue
				}

				// скачиваем страницу только в случае, если выключен кеш и ее нет в базе
				if !d.cfg.PagesCache && !isExists {
					d.logger.Debug(fmt.Sprintf("downloading page: %s", pageURL))

					body, err := d.download(ctx, pageURL)
					if err != nil {
						d.logger.Error(fmt.Errorf("downloading document: %w", err).Error())

						continue
					}

					page = domain.PageEntity{
						URL:  pageURL,
						Body: body,
					}

					page.EntityID = entityID

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

func (d *PageCollector) download(ctx context.Context, pageURL string) (string, error) {
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
