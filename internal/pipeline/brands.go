package pipeline

import "go.uber.org/zap"

func NewBrand(logger zap.Logger) *Brand {
	return &Brand{logger: logger}
}

type Brand struct {
	logger zap.Logger
}
