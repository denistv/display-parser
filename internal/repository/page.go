package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"display_parser/internal/domain"
)

func NewPage(goquDB *goqu.Database) *Page {
	return &Page{
		goquDB: goquDB,
		table:  "pages",
	}
}

type Page struct {
	goquDB *goqu.Database
	table  string
}

func (p *Page) All(ctx context.Context) ([]domain.PageEntity, error) {
	var pages []domain.PageEntity

	err := p.goquDB.From(p.table).ScanStructsContext(ctx, &pages)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return pages, nil
}

func (p *Page) Find(ctx context.Context, pageURL string) (domain.PageEntity, bool, error) {
	var page domain.PageEntity

	ok, err := p.goquDB.
		From("pages").
		Where(goqu.C("url").Eq(pageURL)).
		ScanStructContext(ctx, &page)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PageEntity{}, false, nil
		}

		return domain.PageEntity{}, false, fmt.Errorf("exec query: %w", err)
	}

	return page, ok, nil
}

func (p *Page) Create(ctx context.Context, doc domain.PageEntity) error {
	_, err := goqu.
		Insert(p.table).
		Rows(doc).
		OnConflict(goqu.DoUpdate("url", doc)).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("building sql query: %w", err)
	}

	return nil
}
