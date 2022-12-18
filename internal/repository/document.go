package repository

import (
	"context"
	"database/sql"
	"displayCrawler/internal/domain"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

func NewDocument(db *DBWrapper, goquDB *goqu.Database) *Document {
	return &Document{
		dbw:   db,
		goquDB: goquDB,
		table: "documents",
	}
}

type Document struct {
	dbw   *DBWrapper
	goquDB *goqu.Database
	table string
}

func (d *Document) Find(ctx context.Context, modelURL string) (domain.DocumentEntity, bool, error) {
	var doc domain.DocumentEntity

	ok, err := d.goquDB.
		From("documents").
		Where(goqu.C("url").Eq(modelURL)).
		ScanStructContext(ctx, &doc)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.DocumentEntity{}, false, nil
		}

		return domain.DocumentEntity{}, false, fmt.Errorf("exec query: %w", err)
	}

	return doc, ok, nil
}

func (d *Document) Create(ctx context.Context, doc domain.DocumentEntity) error {
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
