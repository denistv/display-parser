package pipeline

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/iface"
	"display_parser/internal/iface/db"
	"display_parser/mocks"
)

func TestPageCollector_Run(t *testing.T) {
	page, err := os.ReadFile("./test_data/html/en/model/4e4a322f.html")
	if err != nil {
		t.Error(err)
	}
	res := http.Response{
		Body:       io.NopCloser(bytes.NewBuffer(page)),
		Status:     http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
	}

	c := mocks.NewHTTPClient(t)
	c.On("Do", mock.AnythingOfType("*http.Request")).Return(&res, nil)

	pageRepo := mocks.NewPageRepository(t)
	pageRepo.
		On("Find", mock.AnythingOfType("*context.emptyCtx"), mock.AnythingOfType("string")).
		Return(domain.PageEntity{URL: "https://example.com", Body: string(page)}, false, nil)
	pageRepo.On(
		"Create",
		mock.AnythingOfType("*context.emptyCtx"),
		mock.AnythingOfType("domain.PageEntity"),
	).
		Return(nil)

	type fields struct {
		logger     *zap.Logger
		pageRepo   db.PageRepository
		httpClient iface.HTTPClient
		cfg        PagesCollectorCfg
	}

	tests := []struct {
		name   string
		fields fields
		want   []domain.PageEntity
	}{
		{
			name: "create new page",
			fields: fields{
				logger:     zap.NewNop(),
				pageRepo:   pageRepo,
				httpClient: c,
				cfg:        PagesCollectorCfg{},
			},
			want: []domain.PageEntity{
				{
					URL:  "https://example.com",
					Body: string(page),
				},
			},
		},
		// todo: existing page
		// {
		//	name: "update existing page",
		//	fields: fields{
		//		logger:     zap.NewNop(),
		//		pageRepo:  pageRepo ,
		//		httpClient: c,
		//		cfg:        PagesCollectorCfg{},
		//	},
		//	want: []domain.PageEntity{
		//		{
		//			URL:       "https://example.com",
		//			Body:      string(page),
		//		},
		//	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &PageCollector{
				logger:     tt.fields.logger,
				pageRepo:   tt.fields.pageRepo,
				httpClient: tt.fields.httpClient,
				cfg:        tt.fields.cfg,
			}

			in := make(chan string, 1)
			in <- "https://example.com"
			close(in)

			out := make([]domain.PageEntity, 0)
			pageChan := d.Run(context.Background(), in)

			for v := range pageChan {
				out = append(out, v)
			}

			assert.Equal(t, tt.want, out)
		})
	}
}
