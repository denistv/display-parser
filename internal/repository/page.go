package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"display_parser/internal/domain"
)

func NewPage(db *DBWrapper, goquDB *goqu.Database) *Page {
	return &Page{
		dbw:    db,
		goquDB: goquDB,
		table:  "documents",
	}
}

type Page struct {
	dbw    *DBWrapper
	goquDB *goqu.Database
	table  string
}

func (d *Page) Find(ctx context.Context, pageURL string) (domain.PageEntity, bool, error) {
	var doc domain.PageEntity

	ok, err := d.goquDB.
		From("pages").
		Where(goqu.C("url").Eq(pageURL)).
		ScanStructContext(ctx, &doc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PageEntity{}, false, nil
		}

		return domain.PageEntity{}, false, fmt.Errorf("exec query: %w", err)
	}

	return doc, ok, nil
}

func (d *Page) Create(ctx context.Context, doc domain.PageEntity) error {
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
