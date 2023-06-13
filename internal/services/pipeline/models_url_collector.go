package pipeline

import (
	"context"
	"display_parser/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func NewModelsURLCollector(logger *zap.Logger, httpClient services.HTTPClient) *ModelsURLCollector {
	return &ModelsURLCollector{
		logger:     logger,
		httpClient: httpClient,
	}
}

// ModelsURLCollector Собирает URL моделей устройств, чтобы затем обойти их и загрузить страницы с описаниями устройств
type ModelsURLCollector struct {
	logger     *zap.Logger
	httpClient services.HTTPClient
}

func (d *ModelsURLCollector) Run(ctx context.Context, brandURLsChan <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for {
			select {
			case brandURL, ok := <-brandURLsChan:
				if !ok {
					return
				}

				err := d.collect(ctx, brandURL, out)
				if err != nil {
					d.logger.Error(fmt.Errorf("collect brand models: %w", err).Error())

					continue
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func (d *ModelsURLCollector) collect(ctx context.Context, brandURL string, out chan string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, brandURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("creating http req: %w", err)
	}

	res, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("getting brand: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("non-200 status code")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return fmt.Errorf("creating document from reader: %w", err)
	}

	doc.
		Find(".model-listing-container-80 > div[id^='model_'] > a[href]").
		Each(func(i int, s *goquery.Selection) {
			href, ok := s.Attr("href")
			if !ok {
				// todo log warn
				return
			}

			if href == "" {
				return
			}

			out <- href
		})

	return nil
}
