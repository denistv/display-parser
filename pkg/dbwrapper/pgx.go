package dbwrapper

import "github.com/jackc/pgx/v5"

func NewPgx(conn *pgx.Conn) *DBWrapper {
	return &DBWrapper{
		Conn: conn,
	}
}
