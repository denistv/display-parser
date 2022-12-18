package pipeline

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

func NewModelsURLCollector(logger *zap.Logger) *ModelsURLCollector {
	return &ModelsURLCollector{logger: logger}
}

// Собирает URL моделей устройств, чтобы затем обойти их и загрузить страницы с описаниями устройств
type ModelsURLCollector struct {
	logger *zap.Logger
}

func (d *ModelsURLCollector) Run(brandsURLChan <-chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for brand := range brandsURLChan {
			err := d.collect(brand, out)
			if err != nil {
				d.logger.Error(fmt.Errorf("collect brand models: %w", err).Error())

				continue
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
