package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"

	"display_parser/internal/domain"
)

type ModelRepository interface {
	Find(ctx context.Context, url string) (domain.ModelEntity, bool, error)
	Create(ctx context.Context, item domain.ModelEntity) error
}

func NewModel(goduDB *goqu.Database) *Model {
	return &Model{
		goquDB: goduDB,
		table:  "models",
	}
}

type Model struct {
	goquDB *goqu.Database
	table  string
}

func (d *Model) Find(_ context.Context, url string) (domain.ModelEntity, bool, error) {
	var model domain.ModelEntity

	ok, err := d.goquDB.
		From(d.table).
		Where(goqu.C("url").Eq(url)).
		ScanStruct(&model)
	if err != nil {
		return domain.ModelEntity{}, false, fmt.Errorf("scanning struct: %w", err)
	}

	return model, ok, nil
}

//nolint:gocritic
func (d *Model) Create(_ context.Context, item domain.ModelEntity) error {
	item.CreatedAt = time.Now()

	_, err := d.goquDB.
		Insert(d.table).
		Rows(item).
		Executor().
		Exec()
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}
