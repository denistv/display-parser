package pipeline

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"display_parser/internal/domain"
	"display_parser/internal/repository"
	"display_parser/mocks"
	"display_parser/pkg/logger"
)

func TestModelPersister_Run(t *testing.T) {
	type fields struct {
		logger    *zap.Logger
		modelRepo func(t *testing.T) repository.ModelRepository
	}
	type args struct {
		in func() <-chan domain.ModelEntity
	}
	tests := []struct {
		name        string
		modelEntity domain.ModelEntity
		fields      fields
		args        args
	}{
		{
			name: "create new entity",
			args: args{
				in: func() <-chan domain.ModelEntity {
					ch := make(chan domain.ModelEntity)
					go func() {
						ch <- domain.ModelEntity{
							ID: 0,
						}
						close(ch)
					}()

					return ch
				},
			},
			fields: fields{
				logger: zap.NewNop(),
				modelRepo: func(t *testing.T) repository.ModelRepository {
					m := mocks.NewModelRepository(t)
					m.On(
						"Create",
						mock.AnythingOfType("*context.emptyCtx"),
						mock.AnythingOfType("domain.ModelEntity"),
					).Return(nil)

					return m
				},
			},
		},
		{
			name: "update existing entity",
			args: args{
				in: func() <-chan domain.ModelEntity {
					ch := make(chan domain.ModelEntity, 1)
					go func() {
						ch <- domain.ModelEntity{
							ID: 1,
						}
						close(ch)
					}()
					return ch
				},
			},
			fields: fields{
				logger: zap.NewNop(),
				modelRepo: func(t *testing.T) repository.ModelRepository {
					m := mocks.NewModelRepository(t)
					m.On(
						"Update",
						mock.AnythingOfType("*context.emptyCtx"),
						mock.AnythingOfType("domain.ModelEntity"),
					).Return(nil)

					return m
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelRepo := tt.fields.modelRepo(t)
			m := NewModelPersister(logger.NewNopWrapper(), modelRepo)
			done := m.Run(context.Background(), tt.args.in())

			<-done
		})
	}
}
