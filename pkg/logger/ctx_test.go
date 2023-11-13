package logger

import (
	"context"
	"reflect"
	"testing"
)

func TestWithFields(t *testing.T) {
	type args struct {
		ctx    context.Context
		fields []Field
	}
	tests := []struct {
		name string
		args args
		want []Field
	}{
		{
			args: args{
				ctx:    context.Background(),
				fields: []Field{
					NewStringField("string-1", "value-1"),
					NewInt64Field("int64-1", int64(1)),
				},
			},
			want: []Field{
				NewStringField("string-1", "value-1"),
				NewInt64Field("int64-1", int64(1)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithFields(tt.args.ctx, tt.args.fields...)
			if got := ExtractFields(ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
