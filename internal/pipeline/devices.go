package pipeline

import "go.uber.org/zap"

func NewDevices(logger zap.Logger) *Devices {
	return &Devices{logger: logger}
}

type Devices struct {
	logger zap.Logger
}