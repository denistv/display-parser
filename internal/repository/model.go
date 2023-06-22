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
	Update(ctx context.Context, item domain.ModelEntity) error
	All(ctx context.Context) ([]domain.ModelEntity, error)
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

func (d *Model) Find(ctx context.Context, url string) (domain.ModelEntity, bool, error) {
	var model domain.ModelEntity

	ok, err := d.goquDB.
		From(d.table).
		Where(goqu.C("url").Eq(url)).
		ScanStructContext(ctx, &model)
	if err != nil {
		return domain.ModelEntity{}, false, fmt.Errorf("scanning struct: %w", err)
	}

	return model, ok, nil
}

//nolint:gocritic
func (d *Model) Create(ctx context.Context, item domain.ModelEntity) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = item.CreatedAt

	_, err := d.goquDB.
		Insert(d.table).
		Rows(item).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("executing query: %w", err)
	}

	return nil
}

func (d *Model) Update(ctx context.Context, item domain.ModelEntity) error {
	item.UpdatedAt = time.Now()

	_, err := d.goquDB.
		Update(d.table).
		Where(
			goqu.C("url").Eq(item.URL),
		).
		Set(item).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("executing update query: %w", err)
	}

	return nil
}

func (d *Model) All(ctx context.Context) ([]domain.ModelEntity, error) {
	var out []domain.ModelEntity

	err := d.goquDB.Select().From(d.table).ScanStructsContext(ctx, &out)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}

	return out, nil
}
