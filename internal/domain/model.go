package domain

import "errors"

type Model struct {
	ID               int64     `db:"id" goqu:"defaultifempty"`
	ExternalID       string    `db:"external_id"`
	URL              string    `db:"url"`
	Brand            Brand     `db:"-"`
	Name             string    `db:"name"`
	Kind             ModelKind `db:"-"`
	ModelYear        int       `db:"-"`
	Size             float32   `db:"-"`
	PanelType        string    `db:"-"`
	Resolution       string    `db:"-"`
	PixelDensity     int       `db:"-"`
	PixelDensityText string
}

func (m *Model) Validate() error {
	if m.ExternalID == "" {
		return errors.New("ExternalID cannot be empty")
	}

	if m.URL == "" {
		return errors.New("URL cannot be empty")
	}

	if m.Name == "" {
		return errors.New("Name cannot be empty")
	}

	return nil
}
