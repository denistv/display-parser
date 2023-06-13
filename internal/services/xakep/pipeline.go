package xakep

import (
	"context"

	"display_parser/internal/services"
)

func NewPipeline(gen *URLGenerator, pdf *PDFDownloader, httpClient services.HTTPClient) *Pipeline {
	p := Pipeline{}

	return &p
}

type Pipeline struct {
	gen *URLGenerator
	pdf *PDFDownloader
}

func (p *Pipeline) Run(_ context.Context) {

}
