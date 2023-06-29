package service_cfg

import (
	"errors"
	"fmt"
)

type Pipeline struct {
	ModelParserCount int
	PageCollector    PagesCollectorCfg
}

func (c *Pipeline) Validate() error {
	if c.ModelParserCount <= 0 {
		return errors.New("parser count must be leather than 0")
	}

	if err := c.PageCollector.Validate(); err != nil {
		return fmt.Errorf("validating page collector config: %w", err)
	}

	return nil
}
