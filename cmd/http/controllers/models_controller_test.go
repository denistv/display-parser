package controllers

import (
	"net/http"
	"reflect"
	"testing"

	"gopkg.in/guregu/null.v4"

	"display_parser/internal/repository"
)

func Test_parseModelQuery(t *testing.T) {
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    *repository.ModelQuery
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
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.YearFrom = null.NewInt(2020, true)

				return &mq
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
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.YearTo = null.NewInt(2020, true)

				return &mq
			}(),
		},
		{
			name: "parse incorrect year-from",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("year-from", "incorrect")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse incorrect year-to",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("year-to", "incorrect")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want:    nil,
			wantErr: true,
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
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.SizeFrom = null.NewFloat(32.4, true)

				return &mq
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
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.SizeTo = null.NewFloat(32.4, true)

				return &mq
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
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse incorrect size-from",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("size-from", "qwe")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want:    nil,
			wantErr: true,
		},

		// PPI
		{
			name: "parse ppi-from",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("ppi-from", "100")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.PPIFrom = null.NewInt(100, true)

				return &mq
			}(),
		},
		{
			name: "parse ppi-to",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("ppi-to", "100")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.PPITo = null.NewInt(100, true)

				return &mq
			}(),
		},
		{
			name: "parse incorrect ppi-to",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("ppi-to", "incorrect")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parse incorrect ppi-from",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()
					q.Set("ppi-from", "incorrect")
					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want:    nil,
			wantErr: true,
		},

		// Full fields
		{
			name: "full query",
			args: args{
				r: func() *http.Request {
					r, err := http.NewRequest(http.MethodGet, "https://example.com", http.NoBody)
					if err != nil {
						t.Error(err)
					}

					q := r.URL.Query()

					q.Set("limit", "100")

					q.Set("year-from", "2020")
					q.Set("year-to", "2021")

					q.Set("ppi-from", "100")
					q.Set("ppi-to", "200")

					q.Set("size-from", "32.1")
					q.Set("size-to", "35.1")

					q.Set("brand", "Apple")

					r.URL.RawQuery = q.Encode()

					return r
				}(),
			},
			want: func() *repository.ModelQuery {
				mq := repository.NewModelQuery()

				mq.Limit = null.NewInt(100, true)

				mq.YearFrom = null.NewInt(2020, true)
				mq.YearTo = null.NewInt(2021, true)

				mq.PPIFrom = null.NewInt(100, true)
				mq.PPITo = null.NewInt(200, true)

				mq.SizeFrom = null.NewFloat(32.1, true)
				mq.SizeTo = null.NewFloat(35.1, true)

				mq.Brand = null.NewString("Apple", true)

				return &mq
			}(),
			wantErr: false,
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
