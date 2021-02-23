package respository

import (
	"context"
	"github.com/jackc/pgx/v4"
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
