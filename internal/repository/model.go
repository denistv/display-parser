package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"

	"display_parser/internal/db"
	"display_parser/internal/domain"
)

type ModelRepository interface {
	Find(ctx context.Context, url string) (domain.ModelEntity, bool, error)
	Create(ctx context.Context, item domain.ModelEntity) error
	Update(ctx context.Context, item domain.ModelEntity) error
	All(ctx context.Context) ([]domain.ModelEntity, error)
}

func NewModel(db db.SQLDatabase) *Model {
	return &Model{
		db:    db,
		table: "models",
	}
}

type Model struct {
	db    db.SQLDatabase
	table string
}

func selectFields() []any {
	return []any{"id", "url", "brand", "series", "name", "year", "size", "ppi", "created_at", "updated_at"}
}

func (d *Model) Find(ctx context.Context, url string) (domain.ModelEntity, bool, error) {
	var model domain.ModelEntity

	query, params, err := goqu.
		Select(selectFields()...).
		From(d.table).
		Where(goqu.C("url").Eq(url)).
		ToSQL()
	if err != nil {
		return domain.ModelEntity{}, false, fmt.Errorf("make query: %w", err)
	}

	row := d.db.QueryRowContext(ctx, query, params...)
	if row.Err() != nil {
		return domain.ModelEntity{}, false, fmt.Errorf("exec query: %w", err)
	}

	err = row.Scan(&model.ID, &model.URL, &model.Brand, &model.Series, &model.Name, &model.Year, &model.Size, &model.PPI, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return domain.ModelEntity{}, false, fmt.Errorf("scan: %w", err)
	}

	return model, false, nil
}

//nolint:gocritic
func (d *Model) Create(ctx context.Context, item domain.ModelEntity) error {
	item.CreatedAt = time.Now()
	item.UpdatedAt = item.CreatedAt

	query, params, err := goqu.Insert(d.table).Rows(item).ToSQL()
	if err != nil {
		return fmt.Errorf("make query: %w", err)
	}

	if _, err = d.db.ExecContext(ctx, query, params...); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (d *Model) Update(ctx context.Context, item domain.ModelEntity) error {
	item.UpdatedAt = time.Now()

	query, params, err := goqu.
		Update(d.table).
		Where(
			goqu.C("url").Eq(item.URL),
		).ToSQL()
	if err != nil {
		return fmt.Errorf("making sql query: %w", err)
	}

	_, err = d.db.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (d *Model) All(ctx context.Context) ([]domain.ModelEntity, error) {
	var out []domain.ModelEntity

	query, params, err := goqu.Select(selectFields()...).From(d.table).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("make query: %w", err)
	}

	rows, err := d.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		item := domain.ModelEntity{}
		err = rows.Scan(&item.ID, &item.URL, &item.Brand, &item.Series, &item.Name, &item.Year, &item.Size, &item.PPI, &item.UpdatedAt, &item.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning item: %w", err)
		}

		out = append(out, item)
	}

	return out, nil
}
