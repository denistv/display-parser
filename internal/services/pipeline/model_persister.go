package pipeline

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
)

func NewModelPersister(logger *zap.Logger, modelRepo repository.ModelRepository) *ModelPersister {
	return &ModelPersister{
		logger:    logger,
		modelRepo: modelRepo,
	}
}

// ModelPersister определяет, требуется ли создать новую сущность или обновить уже существующую и вызывает
// соответствующий метод репозитория
type ModelPersister struct {
	logger    *zap.Logger
	modelRepo repository.ModelRepository
}

func (m *ModelPersister) Run(ctx context.Context, modelChan <-chan domain.ModelEntity) chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			select {
			case model, ok := <-modelChan:
				if !ok {
					return
				}

				m.logger.Debug(fmt.Sprintf("persisting model %s", model.URL))

				if model.ID == 0 {
					err := m.modelRepo.Create(ctx, model)
					if err != nil {
						m.logger.Error(fmt.Errorf("creating model: %w", err).Error())
						continue
					}
				} else {
					err := m.modelRepo.Update(ctx, model)
					if err != nil {
						m.logger.Error(fmt.Errorf("updating model: %w", err).Error())
						continue
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return done
}
