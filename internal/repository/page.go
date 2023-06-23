package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"display_parser/internal/domain"
	"display_parser/internal/iface/db"
)

func NewPage(dbConn db.Pool) *Page {
	return &Page{
		dbConn: dbConn,
		table:  "pages",
	}
}

type Page struct {
	dbConn db.Pool
	table  string
}

func (p *Page) All(ctx context.Context) ([]domain.PageEntity, error) {
	var pages []domain.PageEntity

	query, params, err := goqu.From(p.table).Select("url", "body", "created_at").ToSQL()
	if err != nil {
		return nil, fmt.Errorf("making query: %w", err)
	}

	rows, err := p.dbConn.Query(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		p := domain.PageEntity{}

		err = rows.Scan(&p.URL, &p.Body, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning struct: %w", err)
		}

		pages = append(pages, p)
	}

	return pages, nil
}

func (p *Page) Find(ctx context.Context, pageURL string) (domain.PageEntity, bool, error) {
	query, params, err := goqu.
		From(p.table).
		Where(goqu.C("url").Eq(pageURL)).
		Limit(1).
		ToSQL()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PageEntity{}, false, nil
		}

		return domain.PageEntity{}, false, fmt.Errorf("exec query: %w", err)
	}

	item := domain.PageEntity{}
	row := p.dbConn.QueryRow(ctx, query, params...)
	if err = row.Scan(&item.URL, &item.Body, &item.CreatedAt); err != nil {
		return domain.PageEntity{}, false, fmt.Errorf("scan item: %w", err)
	}

	return item, true, nil
}

func (p *Page) Create(ctx context.Context, page domain.PageEntity) error {
	query, params, err := goqu.
		Insert(p.table).
		Rows(page).OnConflict(goqu.DoUpdate("url", page)).
		ToSQL()
	if err != nil {
		return fmt.Errorf("make query: %w", err)
	}

	if _, err = p.dbConn.Exec(ctx, query, params...); err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}
