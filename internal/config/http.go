package config

import (
	"display_parser/internal/config/service_cfg"
	"errors"
	"fmt"
)

func NewHTTPConfig() HTTPConfig {
	return HTTPConfig{}
}

type HTTPConfig struct {
	DB         service_cfg.DB
	ListenPort int
}

func (h *HTTPConfig) Validate() error {
	if err := h.DB.Validate(); err != nil {
		return fmt.Errorf("validating db config: %w", err)
	}

	if h.ListenPort <= 0 {
		return errors.New("listen port must be > 0")
	}

	return nil
}
