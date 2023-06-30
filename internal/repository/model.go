package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"display_parser/internal/domain"
	"display_parser/internal/iface/db"
	"github.com/doug-martin/goqu/v9"
	"gopkg.in/guregu/null.v4"
)

type ModelRepository interface {
	Find(ctx context.Context, url string) (domain.ModelEntity, bool, error)
	Create(ctx context.Context, item domain.ModelEntity) error
	Update(ctx context.Context, item domain.ModelEntity) error
	All(ctx context.Context, modelQuery ModelQuery) ([]domain.ModelEntity, error)
}

func NewModel(dbConn db.Pool) *Model {
	return &Model{
		db:    dbConn,
		table: "models",
	}
}

type Model struct {
	db    db.Pool
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

	row := d.db.QueryRow(ctx, query, params...)
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

	if _, err = d.db.Exec(ctx, query, params...); err != nil {
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

	_, err = d.db.Exec(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

const selectLimitDefault = 100

func NewModelQuery() ModelQuery {
	return ModelQuery{
		Limit: null.NewInt(selectLimitDefault, true),
	}
}

// ModelQuery Позволяет кастомизовать запрос в репозиторий, например, из http-контроллера
type ModelQuery struct {
	Limit null.Int
	Brand null.String

	SizeFrom null.Float
	SizeTo   null.Float

	PPIFrom null.Int
	PPITo   null.Int

	YearFrom null.Int
	YearTo   null.Int
}

func (m *ModelQuery) Validate() error {
	if m.Limit.Valid && m.Limit.Int64 <= 0 {
		return errors.New("limit must be > 0")
	}

	if m.Brand.Valid && m.Brand.String == "" {
		return errors.New("brand cannot be empty string")
	}

	// Size
	if m.SizeFrom.Valid && m.SizeFrom.Float64 <= 0 {
		return errors.New("size-from must be > 0.0")
	}
	if m.SizeTo.Valid && m.SizeTo.Float64 <= 0 {
		return errors.New("size-to must be > 0.0")
	}
	if (m.SizeFrom.Valid && m.SizeTo.Valid) && m.SizeTo.Float64 < m.SizeFrom.Float64 {
		return errors.New("size-to must greater than size-from")
	}

	// Year
	if m.YearFrom.Valid && m.YearFrom.Int64 <= 0 {
		return errors.New("year-from must be > 0")
	}
	if m.YearTo.Valid && m.YearTo.Int64 <= 0 {
		return errors.New("year-to must be > 0")
	}
	if (m.YearFrom.Valid && m.YearTo.Valid) && m.YearTo.Int64 < m.YearFrom.Int64 {
		return errors.New("year-to must greater than year-from")
	}

	// PPI
	if m.PPIFrom.Valid && m.PPIFrom.Int64 <= 0 {
		return errors.New("ppi-from must be > 0")
	}
	if m.PPITo.Valid && m.PPITo.Int64 <= 0 {
		return errors.New("ppi-to must be > 0")
	}
	if (m.PPIFrom.Valid && m.PPITo.Valid) && m.PPITo.Int64 < m.PPIFrom.Int64 {
		return errors.New("ppi-to must greater than ppi-from")
	}

	return nil
}

func (d *Model) All(ctx context.Context, mq ModelQuery) ([]domain.ModelEntity, error) {
	if err := mq.Validate(); err != nil {
		return nil, fmt.Errorf("validating model query: %w", err)
	}

	out := make([]domain.ModelEntity, 0)

	q := goqu.
		Select(selectFields()...).
		From(d.table)
	if mq.Limit.Valid {
		q = q.Limit(uint(mq.Limit.Int64))
	}
	if mq.Brand.Valid {
		q = q.Where(goqu.C("brand").Eq(mq.Brand.String))
	}

	if mq.YearFrom.Valid {
		q = q.Where(goqu.C("year").Gte(mq.YearFrom.Int64))
	}
	if mq.YearTo.Valid {
		q = q.Where(goqu.C("year").Lte(mq.YearTo.Int64))
	}

	if mq.SizeFrom.Valid {
		q = q.Where(goqu.C("size").Gte(mq.SizeFrom.Float64))
	}
	if mq.SizeTo.Valid {
		q = q.Where(goqu.C("size").Lte(mq.SizeTo.Float64))
	}

	if mq.PPIFrom.Valid {
		q = q.Where(goqu.C("ppi").Gte(mq.PPIFrom.Int64))
	}
	if mq.PPITo.Valid {
		q = q.Where(goqu.C("ppi").Lte(mq.PPITo.Int64))
	}

	query, params, err := q.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("make query: %w", err)
	}

	rows, err := d.db.Query(ctx, query, params...)
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
