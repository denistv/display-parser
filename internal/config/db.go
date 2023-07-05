package config

import (
	"errors"
	"fmt"
)

type DB struct {
	User         string
	Password     string
	Hostname     string
	DBName       string
	Port         int
	PoolMaxConns int
}

func (d *DB) Validate() error {
	if d.Hostname == "" {
		return errors.New("empty hostname")
	}

	if d.User == "" {
		return errors.New("empty user")
	}

	if d.Password == "" {
		return errors.New("empty password")
	}

	if d.DBName == "" {
		return errors.New("empty dbname")
	}

	if d.Port <= 0 {
		return errors.New("port must be >= 0")
	}

	if d.PoolMaxConns <= 0 {
		return errors.New("pool-max-conns must be >= 0")
	}

	return nil
}

func (d *DB) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		d.User,
		d.Password,
		d.Hostname,
		d.Port,
		d.DBName,
		d.PoolMaxConns,
	)
}
