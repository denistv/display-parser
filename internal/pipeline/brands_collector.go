package pipeline

import (
	"displayCrawler/internal/domain"
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

func (b *BrandsCollector) Run() (<-chan domain.Brand, error){
	res, err := http.Get(b.sourceURL)
	if err != nil {
		return nil, fmt.Errorf("getting brands: %w", err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading document: %w", err)
	}

	out := make(chan domain.Brand)

	go func(){
		doc.
			Find(".brand-listing-container-frontpage").
			Each(func(i int, s *goquery.Selection) {
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					b := domain.Brand{
						Name: s.Text(),
					}
					href, _ := s.Attr("href")
					b.Href = href

					out <- b
				})
			})

		close(out)
	}()

	return out, nil
}
