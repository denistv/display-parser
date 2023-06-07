package app

import (
	"errors"
	"fmt"
)

func NewConfigDev() Config {

	return Config{
		DB: DB{
			User:     "postgres",
			Password: "postgres",
			Hostname: "localhost",
			DBName:   "display_parser",
			Port:     5432,
		},
	}
}

func NewConfigFromEnv() Config {
	return Config{}
}

func NewConfigFromJSONFile(path string) (Config, error) {
	return Config{}, nil
}

type Config struct {
	// Количество запускаемых парсеров страниц (слишком мало - будет копится очередь страниц на парсинг, слишком много - будут простаивать впустую)
	PageParserCount int
	DB              DB
}

func (c Config) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("validating db config: %w", err)
	}

	return nil
}

type DB struct {
	User     string
	Password string
	Hostname string
	DBName   string
	Port     int
}

func (d DB) Validate() error {
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

	if d.Port == 0 {
		return errors.New("empty port")
	}

	return nil
}

func (d DB) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s", // schema://user:password@host:port/db
		d.User,
		d.Password,
		d.Hostname,
		d.Port,
		d.DBName,
	)
}

func (d DB) NewSqlxDSN() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		d.User,
		d.Password,
		d.DBName,
		d.Hostname,
	)
}
