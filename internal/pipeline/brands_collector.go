package pipeline

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
	"net/http"
)

func NewBrandsCollector(logger *zap.Logger) *BrandsCollector {
	return &BrandsCollector{
		logger:    logger,
		sourceURL: "https://www.displayspecifications.com",
	}
}

type BrandsCollector struct {
	logger    *zap.Logger
	sourceURL string
}

func (b *BrandsCollector) BrandURLs() <- chan string{
	out := make(chan string)

	go func(){
		defer close(out)

		res, err := http.Get(b.sourceURL)
		if err != nil {
			b.logger.Error(fmt.Errorf("getting brands: %w", err).Error())
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			b.logger.Error("non-200 status code")

			return
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			b.logger.Error(fmt.Errorf("reading document: %w", err).Error())

			return
		}

		if doc.Text() == "" {
			b.logger.Error("empty response")
			return
		}

		doc.
			Find(".brand-listing-container-frontpage").
			Each(func(i int, s *goquery.Selection) {
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					href, _ := s.Attr("href")

					if href == "" {
						b.logger.Warn("empty href")

						return
					}

					out <- href
				})
			})
	}()

	return out
}
