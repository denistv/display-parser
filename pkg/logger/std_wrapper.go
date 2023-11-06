package logger

import (
	"fmt"
	"log"
	"strings"
)

type logLevel string

const (
	debugLevel logLevel = "Debug"
	infoLevel  logLevel = "Info"
	warnLevel  logLevel = "Warn"
	errorLevel logLevel = "Error"
	panicLevel logLevel = "Panic"
	fatalLevel logLevel = "Fatal"
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
	s.logger.Printf(newMsg(debugLevel, msg, fields...))
}

func (s *STDWrapper) Info(msg string, fields ...Field) {
	s.logger.Printf(newMsg(infoLevel, msg, fields...))
}

func (s *STDWrapper) Warn(msg string, fields ...Field) {
	s.logger.Printf(newMsg(warnLevel, msg, fields...))
}

func (s *STDWrapper) Error(msg string, fields ...Field) {
	s.logger.Printf(newMsg(errorLevel, msg, fields...))
}

func (s *STDWrapper) Panic(msg string, fields ...Field) {
	s.logger.Panic(newMsg(panicLevel, msg, fields...))
}

func (s *STDWrapper) Fatal(msg string, fields ...Field) {
	s.logger.Fatal(newMsg(fatalLevel, msg, fields...))
}

func newMsg(level logLevel, msg string, fields ...Field) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("[%s] %s", level, msg))

	if len(fields) > 1 {
		sb.WriteString(" ")

		for i, v := range fields {
			sb.WriteString(fmt.Sprintf(`%v`, v))

			if i == len(fields)-1 {
				continue
			}

			// Разделяем поля запятой и пробелом
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
