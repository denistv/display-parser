package logger

import "testing"

func Test_newMsg(t *testing.T) {
	type args struct {
		level  logLevel
		msg    string
		fields []Field
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "debug message without fields",
			args: args{
				level:  debugLevel,
				msg:    "debug message",
				fields: nil,
			},
			want: "[Debug] debug message",
		},
		{
			name: "debug message with one field (no comma)",
			args: args{
				level:  debugLevel,
				msg:    "debug message",
				fields: []Field{
					NewInt64Field("int64_field", int64(12345)),
				},
			},
			want: "[Debug] debug message {int64_field 12345}",
		},
		{
			name: "debug message with two fields (with comma)",
			args: args{
				level:  debugLevel,
				msg:    "debug message",
				fields: []Field{
					NewInt64Field("int64_field", int64(12345)),
					NewStringField("string_field", "string_value"),
				},
			},
			want: "[Debug] debug message {int64_field 12345}, {string_field string_value}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newMsg(tt.args.level, tt.args.msg, tt.args.fields...); got != tt.want {
				t.Errorf("newMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
