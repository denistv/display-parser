package db

import (
	"context"
	"database/sql"
)

// SQLDatabase Внутри приложения завязываемся на этот интерфейс. Если понадобится сменить либу, работающую с базой,
// делаем адаптер, реализующий этот интерфейс и прокидывающий вызовы в соответствующую либу.
type SQLDatabase interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
