package storage

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func NewStorage(conn *pgx.Conn) *Storage {
	return &Storage{
		Conn: conn,
	}
}

type Storage struct {
	Conn *pgx.Conn
}

func (d *Storage) Close() error {
	return d.Conn.Close(context.Background())
}
