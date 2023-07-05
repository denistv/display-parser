package config

import (
	"fmt"
)

// UNIXDefaultErrCode unexpected error
const UNIXDefaultErrCode = 255

func NewCmdApp() CmdApp {
	return CmdApp{}
}

type CmdApp struct {
	HTTP     HTTP
	DB       DB
	Pipeline Pipeline
}

func (c *CmdApp) Validate() error {
	if err := c.DB.Validate(); err != nil {
		return fmt.Errorf("validating db config: %w", err)
	}

	if err := c.Pipeline.Validate(); err != nil {
		return fmt.Errorf("validating pipeline config: %w", err)
	}

	return nil
}
