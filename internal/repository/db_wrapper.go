package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func NewDBWrapper(conn *pgx.Conn) *DBWrapper {
	return &DBWrapper{
		Conn: conn,
	}
}

type DBWrapper struct {
	Conn *pgx.Conn
}

func (d *DBWrapper) Close() error {
	return d.Conn.Close(context.Background())
}

/*
Begin() (*sql.Tx, error)
BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
 */