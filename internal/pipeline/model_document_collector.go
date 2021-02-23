package pipeline

import (
	"displayCrawler/internal/domain"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"time"

	"displayCrawler/internal/respository"
)

func NewModelDocumentCollector(logger *zap.Logger, docRepo *respository.Document) *ModelDocumentCollector {
	return &ModelDocumentCollector{
		logger:  logger,
		docRepo: docRepo,
	}
}

type ModelDocumentCollector struct {
	logger  *zap.Logger
	docRepo *respository.Document
}

func (d *ModelDocumentCollector) Collect(in <-chan string) chan domain.ModelDocument {
	out := make(chan domain.ModelDocument)

	go func() {
		for docURL := range in {
			body, err := d.download(docURL)
			if err != nil {
				d.logger.Error(fmt.Errorf("downloading document: %w", err).Error())

				continue
			}

			doc := domain.ModelDocument{
				URL:  docURL,
				Body: body,
			}

			out <- doc

			time.Sleep(100 * time.Millisecond)
		}

		close(out)
	}()

	return out
}

func (d *ModelDocumentCollector) download(docURL string) (string, error) {
	if docURL == "" {
		return "", errors.New("url cannot be empty")
	}

	res, err := http.Get(docURL)
	if err != nil {
		return "", fmt.Errorf("getting model: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New("non-200 status code")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("reading document body: %w", err)
	}

	return string(body), nil
}
