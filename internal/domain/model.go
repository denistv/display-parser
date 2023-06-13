package domain

import "time"

// ModelEntity сущность, представляющая разобранную модель монитора.
type ModelEntity struct {
	ID int64 `db:"id" goqu:"defaultifempty"`
	// URL страницы с монитором
	URL string `db:"url"`

	// Название бренда (обязательно)
	Brand string `db:"brand"`

	// Линейка монитора (опционально)
	Series string `db:"series"`

	// Название модели
	Name string `db:"name"`

	// Год выпуска (опционально)
	Year int64 `db:"year"`

	// Диагональ (опционально)
	Size int64 `db:"size"`

	// Число точек на дюйм (опционально)
	PPI int64 `db:"ppi"`

	// Время сохранения модели в БД
	CreatedAt time.Time `db:"created_at"`
}
