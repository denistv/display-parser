package db

import (
	"context"

	"display_parser/internal/domain"
)

type PageRepository interface {
	All(ctx context.Context) ([]domain.PageEntity, error)
	Find(ctx context.Context, entityID domain.EntityID) (domain.PageEntity, bool, error)
	Create(ctx context.Context, page domain.PageEntity) error
	PageIsExists(entityID domain.EntityID) (domain.PageEntity, bool)
}
