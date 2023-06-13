package pipeline

import (
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
	"display_parser/mocks"
)

func TestModelParser_parsePPI(t *testing.T) {
	pageRaw, err := os.ReadFile("./test_data/html/en/model/4e4a322f.html")
	if err != nil {
		t.Error(err)
	}

	page := domain.PageEntity{
		Body: string(pageRaw),
	}

	modelRepo := mocks.NewModelRepository(t)

	type fields struct {
		logger     *zap.Logger
		modelsRepo repository.ModelRepository
	}
	type args struct {
		page domain.PageEntity
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "parse correct PPI",
			fields: fields{
				logger:     zap.NewNop(),
				modelsRepo: modelRepo,
			},
			args:    args{page: page},
			want:    93,
			wantErr: false,
		},
		{
			name: "parse incorrect PPI -- want error",
			fields: fields{
				logger:     zap.NewNop(),
				modelsRepo: modelRepo,
			},
			args:    args{page: domain.PageEntity{Body: ""}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelParser{
				logger:     tt.fields.logger,
				modelsRepo: tt.fields.modelsRepo,
			}
			got, err := m.parsePPI(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parsePPI() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelParser_parse(t *testing.T) {
	pageRaw, err := os.ReadFile("./test_data/html/en/model/4e4a322f.html")
	if err != nil {
		t.Error(err)
	}

	page := domain.PageEntity{
		URL:  "https://example.com/page",
		Body: string(pageRaw),
	}

	modelRepo := mocks.NewModelRepository(t)

	type fields struct {
		logger     *zap.Logger
		modelsRepo repository.ModelRepository
	}
	type args struct {
		page domain.PageEntity
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    domain.ModelEntity
		wantErr bool
	}{
		{
			name: "full page parse",
			fields: fields{
				logger:     zap.NewNop(),
				modelsRepo: modelRepo,
			},
			args: args{page: page},
			want: domain.ModelEntity{
				URL:    page.URL,
				Brand:  "Acer",
				Series: "XZ2",
				Name:   "XZ323QU X3",
				Year:   2023,
				PPI:    93,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelParser{
				logger:     tt.fields.logger,
				modelsRepo: tt.fields.modelsRepo,
			}
			got, err := m.parse(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
