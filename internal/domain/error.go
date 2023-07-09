package domain

import (
	"errors"
	"fmt"
)

func NewValidationError(s string) error {
	return fmt.Errorf("%w: %s", ErrValidationError, s)
}

var ErrValidationError = errors.New("validation error")
