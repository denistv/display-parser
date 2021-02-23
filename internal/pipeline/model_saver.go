package pipeline

import (
	"context"
	"displayCrawler/internal/domain"
	"displayCrawler/internal/respository"
	"fmt"
	"go.uber.org/zap"
)

func NewModelSaver(logger *zap.Logger, deviceRepo *respository.Model) *ModelSaver {
	return &ModelSaver{
		logger:     logger,
		deviceRepo: deviceRepo,
	}
}

type ModelSaver struct {
	logger     *zap.Logger
	deviceRepo *respository.Model
}

func (s *ModelSaver) Persist(ctx context.Context, in <-chan domain.Model) {
	for device := range in {
		err := s.deviceRepo.Create(ctx, &device)
		if err != nil {
			s.logger.Error(fmt.Errorf("creating device in db: %w", err).Error())
			continue
		}
	}
}
