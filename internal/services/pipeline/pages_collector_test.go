package pipeline

import (
	"bytes"
	"context"
	"display_parser/internal/config"
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
	c.On("Do", mock.IsType(&http.Request{})).Return(&res, nil)

	pageRepo := mocks.NewPageRepository(t)
	pageRepo.
		On("Find", mock.AnythingOfType("*context.emptyCtx"), mock.IsType(domain.EntityID(""))).
		Return(domain.PageEntity{URL: "https://example.com/en/model/4e4a322f", Body: string(page), EntityID: "4e4a322f"}, false, nil)
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
		cfg        config.PagesCollector
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
				cfg:        config.PagesCollector{},
			},
			want: []domain.PageEntity{
				{
					URL:      "https://example.com/en/model/4e4a322f",
					Body:     string(page),
					EntityID: domain.EntityID("4e4a322f"),
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
			in <- "https://example.com/en/model/4e4a322f"
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
