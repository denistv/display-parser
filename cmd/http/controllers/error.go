package controllers

import (
	"display_parser/internal/domain"
	"errors"
	"fmt"
)

func NewParseError(s string) error {
	return fmt.Errorf("%w: %s", ErrParseError, s)
}

var ErrParseError = errors.New("parse error")

func is400Err(err error) bool {
	return errors.Is(err, domain.ErrValidationError) || errors.Is(err, ErrParseError)
}
