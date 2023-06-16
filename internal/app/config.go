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

func (d DB) ConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		d.User,
		d.Password,
		d.Hostname,
		d.Port,
		d.DBName,
	)
}

func (d DB) ConnStringSQLX() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s sslmode=disable",
		d.User,
		d.Password,
		d.DBName,
		d.Hostname,
	)
}
