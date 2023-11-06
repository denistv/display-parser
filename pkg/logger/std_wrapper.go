package logger

import "log"

func NewStdWrapper() *STDWrapper{
	return &STDWrapper{
		logger: log.Default(),
	}
}

// STDWrapper обертка для логгера из стандартной библиотеки Go
type STDWrapper struct {
	logger *log.Logger
}
