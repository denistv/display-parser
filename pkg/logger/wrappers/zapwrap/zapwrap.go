package zapwrap

import (
	"go.uber.org/zap"

	"display_parser/pkg/logger"
)

func NewZapWrapper(zap *zap.Logger) *ZapWrapper {
	return &ZapWrapper{
		zap: zap,
	}
}

// ZapWrapper обертка для Zap
type ZapWrapper struct {
	zap *zap.Logger
}

func (z *ZapWrapper) Debug(msg string, fields ...logger.Field) {
	z.zap.Debug(msg, newZapFields(fields...)...)
}

func (z *ZapWrapper) Info(msg string, fields ...logger.Field) {
	z.zap.Info(msg, newZapFields(fields...)...)
}

func (z *ZapWrapper) Warn(msg string, fields ...logger.Field) {
	z.zap.Warn(msg, newZapFields(fields...)...)
}

func (z *ZapWrapper) Error(msg string, fields ...logger.Field) {
	z.zap.Error(msg, newZapFields(fields...)...)
}

func (z *ZapWrapper) Panic(msg string, fields ...logger.Field) {
	z.zap.Panic(msg, newZapFields(fields...)...)
}

func (z *ZapWrapper) Fatal(msg string, fields ...logger.Field) {
	z.zap.Fatal(msg, newZapFields(fields...)...)
}

func (z *ZapWrapper) Sync() error {
	return z.zap.Sync()
}

// NewZapField реализована не самым удачным образом. Когда-нибудь вернусь к этому вопросу (или не вернусь). Сейчас это не принципиальный вопрос.
func newZapField(f logger.Field) zap.Field {
	return zap.Field{
		Key:       f.Key,
		Interface: f.Value,
	}
}

func newZapFields(fs ...logger.Field) []zap.Field {
	out := make([]zap.Field, 0, len(fs))

	for _, f := range fs {
		out = append(out, newZapField(f))
	}

	return out
}