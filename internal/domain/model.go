package domain

import "time"

// ModelEntity сущность, представляющая разобранную модель монитора.
type ModelEntity struct {
	// URL страницы с монитором
	URL string `db:"url"`

	// Название бренда
	Brand string `db:"brand"`

	// Линейка монитора (может отсутствовать)
	Series string `db:"series"`

	// Название модели
	Name string `db:"name"`

	// Год выпуска
	Year int64 `db:"year"`

	// Диагональ
	Size int64 `db:"size"`

	// Число точек на дюйм
	PPI int64 `db:"ppi"`

	// Время сохранения модели в БД
	CreatedAt time.Time
}
