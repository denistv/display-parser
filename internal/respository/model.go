package respository

import (
	"context"
	"fmt"

	"displayCrawler/internal/domain"
	"github.com/doug-martin/goqu/v9"
)

func NewModel(db *DBWrapper) *Model {
	return &Model{
		dbw:   db,
		table: "models",
	}
}

type Model struct {
	dbw   *DBWrapper
	table string
}

func (d *Model) Create(ctx context.Context, item *domain.Model) error {
	sqlQuery, args, err := goqu.
		Insert(d.table).
		Rows(item).
		ToSQL()
	if err != nil {
		return fmt.Errorf("building sql query: %w", err)
	}

	rows, err := d.dbw.Conn.Query(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("inserting device to db: %w", err)
	}

	rows.Close()

	return nil
}
