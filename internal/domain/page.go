package domain

import (
	"fmt"
	"time"
)

// PageEntity представляет собой сущность страницы с описанием монитора на сайте.
// Страница содержит Body, который можно многократно парсить (при необходимости расширения сущности ModelEntity),
// то есть, единожды сохраненную страницу можно многократно парсить без хождения в сеть, если возникла потребность
// добавить в модель монитора новые поля.
type PageEntity struct {
	URL       string    `db:"url"`
	Body      string    `db:"body"`
	EntityID  EntityID  `db:"entity_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (p *PageEntity) Validate() error {
	if err := p.EntityID.Validate(); err != nil {
		return fmt.Errorf("validating page entity ID: %w", err)
	}

	return nil
}
