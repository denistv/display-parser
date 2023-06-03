package dbwrapper

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// DBWrapper нужен для того, чтобы абстрагироваться от конкретной реализации драйвера для работы с БД и при необходимости,
// иметь возможность быстро поменять его на другой без переписывания всего приложения.
// Приложение завязывается на общий интерфейс враппера который не меняется и адаптируется под конкретную библиотеку.
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
