package pipeline

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"display_parser/internal/services"
	"display_parser/mocks"
)

func TestBrandsCollector_Run(t *testing.T) {
	brandURLsRaw, err := os.ReadFile("./test_data/html/index_urls.txt")
	if err != nil {
		t.Error(err)
	}

	brandURLs := strings.Split(strings.ReplaceAll(string(brandURLsRaw), "\r\n", "\n"), "\n")

	indexPageHTML, err := os.ReadFile("./test_data/html/index.html")
	if err != nil {
		t.Error(err)
	}

	res := http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(bytes.NewBuffer(indexPageHTML)),
		ContentLength: -1,
	}

	httpClient := mocks.NewHTTPClient(t)
	httpClient.On("Do", mock.AnythingOfType("*http.Request")).Return(&res, nil)

	type fields struct {
		logger     *zap.Logger
		sourceURL  string
		httpClient services.HTTPClient
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "parsing index page for collect brand URLs",
			fields: fields{
				logger:     zap.NewNop(),
				httpClient: httpClient,
			},
			args: args{ctx: context.Background()},
			want: brandURLs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)

			b := &BrandsCollector{
				logger:     tt.fields.logger,
				sourceURL:  tt.fields.sourceURL,
				httpClient: tt.fields.httpClient,
				cancel:     cancel,
			}

			ch := b.Run(ctx)
			urls := make([]string, 0)

			for v := range ch {
				urls = append(urls, v)
			}

			assert.Equal(t, tt.want, urls)
		})
	}
}
