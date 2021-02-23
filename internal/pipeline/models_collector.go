package pipeline

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
)

func NewModelsCollector(logger *zap.Logger) *ModelsCollector {
	return &ModelsCollector{logger: logger}
}

type ModelsCollector struct {
	logger *zap.Logger
}

func (d *ModelsCollector) GetItemsIndex(brandsURLChan <-chan string) chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for brand := range brandsURLChan {
			err := d.collectBrandModels(brand, out)
			if err != nil {
				d.logger.Error(fmt.Errorf("collect brand models: %w", err).Error())

				continue
			}
		}
	}()

	return out
}

func (d *ModelsCollector) collectBrandModels(brandURL string, out chan string) error {
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
