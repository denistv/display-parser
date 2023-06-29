package config

import (
	"fmt"

	"display_parser/internal/config/service_cfg"
)

// UNIXDefaultErrCode unexpected error
const UNIXDefaultErrCode = 255

func NewAppConfig() AppConfig {
	return AppConfig{}
}

type AppConfig struct {
	HTTP     service_cfg.HTTP
	DB       service_cfg.DB
	Pipeline service_cfg.Pipeline
}

func (c *AppConfig) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("validating db config: %w", err)
	}

	if err := c.Pipeline.Validate(); err != nil {
		return fmt.Errorf("validating pipeline config: %w", err)
	}

	return nil
}
