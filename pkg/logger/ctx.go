package logger

import (
	"context"
)

type ctxKey int

const ctxFieldsKey ctxKey = iota

func WithFields(ctx context.Context, fields ...Field) context.Context {
	s, ok := ctx.Value(ctxFieldsKey).(*ctxFields)
	if !ok || s == nil{
		s = newCtxFields(ctx)
	}

	s.Push(fields...)

	return context.WithValue(ctx, ctxFieldsKey, s)

}

func ExtractFields(ctx context.Context, fields ...Field) []Field {
	s, ok := ctx.Value(ctxFieldsKey).(*ctxFields)
	if !ok {
		return nil
	}

	ctxFieldsAll := s.All()

	if len(fields) != 0 {
		out := make([]Field, 0, len(fields) + len(ctxFieldsAll))
		out = append(out, ctxFieldsAll...)
		out = append(out, fields...)
		return out
	}

	return ctxFieldsAll
}


func newCtxFields(ctx context.Context) *ctxFields {
	return &ctxFields{
		parent: ctx,
		fields: make([]Field, 0),
	}
}

type ctxFields struct {
	parent context.Context

	fields []Field
}

func (c *ctxFields) Push(fields ...Field) {
	for _, f := range fields {
		c.fields = append(c.fields, f)
	}
}

func (c *ctxFields) All() []Field {
	return c.fields
}