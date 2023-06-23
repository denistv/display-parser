package app

import (
	"errors"
	"fmt"
	"time"

	"display_parser/internal/services/pipeline"
)

func NewConfig() Config {
	return Config{}
}

type Config struct {
	HTTP     HTTP
	DB       DB
	Pipeline pipeline.Cfg
}

func (c *Config) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("validating db config: %w", err)
	}

	if err := c.Pipeline.Validate(); err != nil {
		return fmt.Errorf("validating pipeline config: %w", err)
	}

	return nil
}

type HTTP struct {
	// Задержка между HTTP-запросами в сервис
	// Если не использовать ограничений, сервис забанит вас на какое-то время.
	Timeout         time.Duration
	DelayPerRequest time.Duration
}

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

// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
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
