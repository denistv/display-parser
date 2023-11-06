package logger

import "go.uber.org/zap"

func NewZapWrapper(zap *zap.Logger) *ZapWrapper {
	return &ZapWrapper{
		zap: zap,
	}
}

// ZapWrapper обертка для Zap
type ZapWrapper struct{
	zap *zap.Logger
}
