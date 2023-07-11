package domain

import (
	"database/sql/driver"
	"errors"
	"regexp"
)

var entityIDRegexp = regexp.MustCompile("^(?P<schema>http|https)://(?P<domain>[A-z0-9.]+)/(?P<lang>[a-z]+)/model/(?P<id>[a-z0-9]+)$")

// NewEntityIDFromURL Идентификатор сущности модели во внешней системе. Извлекаем его из URL
func NewEntityIDFromURL(pageURL string) (EntityID, error) {
	matches := entityIDRegexp.FindStringSubmatch(pageURL)
	if len(matches) == 0 {
		return "", errors.New("incorrect URL: no matches")
	}

	i := entityIDRegexp.SubexpIndex("id")

	if i == -1 {
		return "", errors.New("cannot get ID from page URL")
	}

	return EntityID(matches[i]), nil
}

type EntityID string

func (e *EntityID) Validate() error {
	if e.String() == "" {
		return NewValidationError("entity ID cannot be empty")
	}

	return nil
}

func (e *EntityID) Value() (driver.Value, error) {
	return e.String(), nil
}

func (e *EntityID) Scan(src any) error {
	v, ok := src.(string)
	if !ok {
		return errors.New("value in db is not string")
	}

	ei := EntityID(v)
	if err := ei.Validate(); err != nil {
		return err
	}

	*e = ei

	return nil
}

func (e *EntityID) String() string {
	return string(*e)
}
