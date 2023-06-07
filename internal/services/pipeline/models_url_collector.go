package pipeline

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func NewModelsURLCollector(logger *zap.Logger) *ModelsURLCollector {
	return &ModelsURLCollector{logger: logger}
}

// ModelsURLCollector Собирает URL моделей устройств, чтобы затем обойти их и загрузить страницы с описаниями устройств
type ModelsURLCollector struct {
	logger *zap.Logger
}

func (d *ModelsURLCollector) Run(ctx context.Context, brandURLsChan <-chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for {
			select {
			case brandURL := <-brandURLsChan:
				err := d.collect(brandURL, out)
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

func (d *ModelsURLCollector) collect(brandURL string, out chan string) error {
	res, err := http.Get(brandURL)
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
