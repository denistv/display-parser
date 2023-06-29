package controllers

import (
	"gopkg.in/guregu/null.v4"
	"net/http"
	"reflect"
	"testing"

	"display_parser/internal/repository"
)

func Test_parseModelQuery(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    repository.ModelQuery
		wantErr bool
	}{
		// Year
		{
			name: "parse year-from",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("year-from", "2020")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.YearFrom = null.NewInt(2020, true)

				return mq
			}(),
		},
		{
			name: "parse year-to",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("year-to", "2020")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.YearTo = null.NewInt(2020, true)

				return mq
			}(),
		},

		// Size
		{
			name: "parse size-from",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("size-from", "32.4")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.SizeFrom = null.NewFloat(32.4, true)

				return mq
			}(),
		},
		{
			name: "parse size-to",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("size-to", "32.4")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.SizeTo = null.NewFloat(32.4, true)

				return mq
			}(),
		},
		{
			name: "parse incorrect size-to",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("size-to", "qwe")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want:    repository.ModelQuery{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseModelQuery(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseModelQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseModelQuery() got = %v, want %v", got, tt.want)
			}
		})
	}
}
