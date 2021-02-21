package pipeline

import "go.uber.org/zap"

func NewSaver(logger zap.Logger) *Saver{
	return &Saver{
		logger: logger,
	}
}

type Saver struct {
	logger zap.Logger
}
