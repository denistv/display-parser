package config

import (
	"errors"
	"fmt"
	"net/url"
)

func NewCmdHTTP() CmdHTTP {
	return CmdHTTP{}
}

type CmdHTTP struct {
	DB                DB
	ListenPort        int
	CORSAllowedOrigin string
}

func (h *CmdHTTP) Validate() error {
	if err := h.DB.Validate(); err != nil {
		return fmt.Errorf("validating db config: %w", err)
	}

	if h.ListenPort <= 0 {
		return errors.New("listen port must be > 0")
	}

	if h.CORSAllowedOrigin != "" {
		if _, err := url.Parse(h.CORSAllowedOrigin); err != nil {
			return fmt.Errorf("parsing cors allowed origin: %w", err)
		}
	}

	return nil
}
