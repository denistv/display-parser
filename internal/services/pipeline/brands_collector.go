package pipeline

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"

	"display_parser/internal/services"
)

func NewBrandsCollector(logger *zap.Logger, httpClient services.HTTPClient, cancel context.CancelFunc) *BrandsCollector {
	return &BrandsCollector{
		logger:     logger,
		sourceURL:  "https://www.displayspecifications.com", // TODO вынести в конфиг
		httpClient: httpClient,
		cancel:     cancel,
	}
}

type BrandsCollector struct {
	logger     *zap.Logger
	sourceURL  string
	httpClient services.HTTPClient

	// Останавливаем через канал отмены работу программы в случаях, когда дальнейшая работа не имеет смысла
	cancel context.CancelFunc
}

func (b *BrandsCollector) Run(ctx context.Context) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		b.logger.Info("getting brands...")

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.sourceURL, http.NoBody)
		if err != nil {
			b.logger.Error(fmt.Errorf("creating http req: %w", err).Error())
			b.cancel()
			return
		}

		res, err := b.httpClient.Do(req)
		if err != nil {
			b.logger.Error(fmt.Errorf("getting brands: %w", err).Error())
			b.cancel()
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			b.logger.Error("non-200 status code")
			b.cancel()

			return
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			b.logger.Error(fmt.Errorf("reading document: %w from remote server", err).Error())
			b.cancel()
			return
		}

		if doc.Text() == "" {
			b.logger.Error("empty response from remote server")
			b.cancel()
			return
		}

		urls := make([]string, 0)

		doc.
			Find(".brand-listing-container-frontpage").
			Each(func(i int, s *goquery.Selection) {
				s.Find("a").Each(func(i int, s *goquery.Selection) {
					href, _ := s.Attr("href")

					if href == "" {
						b.logger.Warn("empty href")

						return
					}

					b.logger.Debug(fmt.Sprintf("brand collected: %s", href))

					urls = append(urls, href)
				})
			})

		for _, v := range urls {
			out <- v
		}

		b.logger.Info("brands collected")
	}()

	return out
}
