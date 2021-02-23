package respository

import (
	"context"
	"displayCrawler/internal/domain"
	"fmt"
	"github.com/doug-martin/goqu/v9"
)

func NewModelDocument(db *DBWrapper) *Document {
	return &Document{
		dbw:   db,
		table: "documents",
	}
}


type Document struct {
	dbw   *DBWrapper
	table string
}

func (d *Document) Create(ctx context.Context, doc domain.ModelDocument) error {
	sqlQuery, args, err := goqu.
		Insert(d.table).
		Rows(doc).
		OnConflict(goqu.DoUpdate("url", doc)).
		ToSQL()
	if err != nil {
		return fmt.Errorf("building sql query: %w", err)
	}

	rows, err := d.dbw.Conn.Query(ctx, sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("inserting item to db: %w", err)
	}

	rows.Close()

	return nil
}
