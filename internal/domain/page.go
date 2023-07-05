package domain

import "time"

// PageEntity представляет собой сущность страницы с описанием монитора на сайте.
// Страница содержит Body, который можно многократно парсить (при необходимости расширения сущности ModelEntity),
// то есть, единожды сохраненную страницу можно многократно парсить без хождения в сеть, если возникла потребность
// добавить в модель монитора новые поля.
type PageEntity struct {
	URL       string    `db:"url"`
	Body      string    `db:"body"`
	EntityID  string    `db:"entity_id"`
	CreatedAt time.Time `db:"created_at"`
}
