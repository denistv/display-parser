package pipeline

import (
	"context"
	"displayCrawler/internal/domain"
	"displayCrawler/internal/storage"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

func NewSaver(logger *zap.Logger, stor *storage.Storage) *Saver {
	return &Saver{
		logger: logger,
		stor:   stor,
	}
}

type Saver struct {
	logger *zap.Logger
	stor   *storage.Storage
}

func (s *Saver) Persist(ctx context.Context, in <-chan domain.Device) {
	for device := range in {
		sqlQuery, args, err := goqu.Insert("devices").Rows(device).ToSQL()
		if err != nil {
			s.logger.Error(fmt.Errorf("building sql query: %w", err).Error())
			continue
		}

		_, err = s.stor.Conn.Query(ctx, sqlQuery, args...)
		if err != nil {
			s.logger.Error(fmt.Errorf("inserting device to db: %w", err).Error())
			continue
		}
	}
}
