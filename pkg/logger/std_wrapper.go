package logger

import (
	"fmt"
	"log"
	"strings"
)

type LogLevel string

const (
	Debug LogLevel = "Debug"
	Info  LogLevel = "Info"
	Warn  LogLevel = "Warn"
	Error LogLevel = "Error"
	Panic LogLevel = "Panic"
	Fatal LogLevel = "Fatal"
)

func NewStdWrapper() *STDWrapper {
	return &STDWrapper{
		logger: log.Default(),
	}
}

// STDWrapper обертка для логгера из стандартной библиотеки Go
type STDWrapper struct {
	logger *log.Logger
}

func (s *STDWrapper) Debug(msg string, fields ...Field) {
	s.logger.Printf(newMsg(Debug, msg, fields...))
}

func (s *STDWrapper) Info(msg string, fields ...Field) {
	s.logger.Printf(newMsg(Info, msg, fields...))
}

func (s *STDWrapper) Warn(msg string, fields ...Field) {
	s.logger.Printf(newMsg(Warn, msg, fields...))
}

func (s *STDWrapper) Error(msg string, fields ...Field) {
	s.logger.Printf(newMsg(Error, msg, fields...))
}

func (s *STDWrapper) Panic(msg string, fields ...Field) {
	s.logger.Panic(newMsg(Panic, msg, fields...))
}

func (s *STDWrapper) Fatal(msg string, fields ...Field) {
	s.logger.Fatal(newMsg(Fatal, msg, fields...))
}

func newMsg(level LogLevel, msg string, fields ...Field) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("[%s] %s", level, msg))

	if len(fields) > 0 {
		sb.WriteString(" ")
		for _, v := range fields {
			sb.WriteString(fmt.Sprintf(`{key: "%s", value: %q}`, v.Key, v.Value))
		}
	}

	return sb.String()
}
