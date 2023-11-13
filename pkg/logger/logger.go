package logger

import "context"

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

// CtxLogger по сути тот же интерфейс Logger, только методы CtxLogger первым аргументом принимают контекст, где находятся поля логгера
type CtxLogger interface {
	DebugCtx(ctx context.Context, msg string, fields ...Field)
	InfoCtx(ctx context.Context, msg string, fields ...Field)
	WarnCtx(ctx context.Context, msg string, fields ...Field)
	ErrorCtx(ctx context.Context, msg string, fields ...Field)
	PanicCtx(ctx context.Context, msg string, fields ...Field)
	FatalCtx(ctx context.Context, msg string, fields ...Field)
}
