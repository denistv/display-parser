package repository

import (
	"gopkg.in/guregu/null.v4"
	"testing"
)

func TestModelQuery_Validate(t *testing.T) {
	type fields struct {
		Limit    null.Int
		Brand    null.String
		SizeFrom null.Float
		SizeTo   null.Float
		PPIFrom  null.Int
		PPITo    null.Int
		YearFrom null.Int
		YearTo   null.Int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// Year fields
		{
			name: "correct year-from/year-to",
			fields: fields{
				YearFrom: null.NewInt(2023, true),
				YearTo:   null.NewInt(2024, true),
			},
			wantErr: false,
		},
		{
			name: "incorrect year-from > year-to",
			fields: fields{
				YearFrom: null.NewInt(2024, true),
				YearTo:   null.NewInt(2023, true),
			},
			wantErr: true,
		},
		{
			name: "year-from < 0",
			fields: fields{
				YearFrom: null.NewInt(-1, true),
			},
			wantErr: true,
		},
		{
			name: "year-to < 0",
			fields: fields{
				YearTo: null.NewInt(-1, true),
			},
			wantErr: true,
		},

		// PPI fields
		{
			name: "correct ppi-from/ppi-to",
			fields: fields{
				PPIFrom: null.NewInt(100, true),
				PPITo:   null.NewInt(200, true),
			},
			wantErr: false,
		},
		{
			name: "incorrect ppi-from > ppi-to",
			fields: fields{
				PPIFrom: null.NewInt(200, true),
				PPITo:   null.NewInt(100, true),
			},
			wantErr: true,
		},
		{
			name: "ppi-from < 0",
			fields: fields{
				PPIFrom: null.NewInt(-1, true),
			},
			wantErr: true,
		},
		{
			name: "ppi-to < 0",
			fields: fields{
				PPITo: null.NewInt(-1, true),
			},
			wantErr: true,
		},

		//Size fields
		{
			name: "correct size-from/size-to",
			fields: fields{
				SizeFrom: null.NewFloat(27.3, true),
				SizeTo:   null.NewFloat(32.5, true),
			},
			wantErr: false,
		},
		{
			name: "size-from > size-to",
			fields: fields{
				SizeFrom: null.NewFloat(32.5, true),
				SizeTo:   null.NewFloat(27.3, true),
			},
			wantErr: true,
		},
		{
			name: "size-from < 0",
			fields: fields{
				SizeFrom: null.NewFloat(-1, true),
			},
			wantErr: true,
		},
		{
			name: "size-to < 0",
			fields: fields{
				SizeTo: null.NewFloat(-1, true),
			},
			wantErr: true,
		},

		// Brand
		{
			name: "brand passed, valid=true",
			fields: fields{
				Brand: null.NewString("Apple", true),
			},
			wantErr: false,
		},
		{
			name: "empty brand passed, valid=true",
			fields: fields{
				Brand: null.NewString("", true),
			},
			wantErr: true,
		},
		// Limit
		{
			name: "correct limit passed",
			fields: fields{
				Limit: null.NewInt(100, true),
			},
			wantErr: false,
		},
		{
			name: "limit < 0 passed",
			fields: fields{
				Limit: null.NewInt(-1, true),
			},
			wantErr: true,
		},
		{
			name: "limit == 0 passed",
			fields: fields{
				Limit: null.NewInt(0, true),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelQuery{
				Limit:    tt.fields.Limit,
				Brand:    tt.fields.Brand,
				SizeFrom: tt.fields.SizeFrom,
				SizeTo:   tt.fields.SizeTo,
				PPIFrom:  tt.fields.PPIFrom,
				PPITo:    tt.fields.PPITo,
				YearFrom: tt.fields.YearFrom,
				YearTo:   tt.fields.YearTo,
			}
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
