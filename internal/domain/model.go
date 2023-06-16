package domain

import "time"

// ModelEntity сущность, представляющая разобранную модель монитора.
type ModelEntity struct {
	ID int64 `db:"id" goqu:"defaultifempty"`
	// URL страницы с монитором
	URL string `db:"url"`
	Brand string `db:"brand"`
	Series string `db:"series"`
	Name string `db:"name"`
	Year int64 `db:"year"`
	// Диагональ (опционально),может быть дробным числом
	Size float64 `db:"size"`
	PPI int64 `db:"ppi"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
