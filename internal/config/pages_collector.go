package config

import "errors"

type PagesCollector struct {
	Count int
	// Пересобрать модели на основе кэша страниц в базе. Если флаг взведен, не ходим во внешний сервис для сбора данных и используем имеющийся кэш страниц в БД.
	// Полезно в тех случаях, когда сайт спаршен (данные страниц сохранены в кэше в таблице pages, но сущность модели расширена дополнительным полем.
	// Чтобы не собирать все данные по новой через сеть, используем сохраненные в базу страницы и перераспаршиваем их, обновляя сущности моделей).
	PagesCache bool
}

func (p *PagesCollector) Validate() error {
	if p.Count <= 0 {
		return errors.New("count must be > 0")
	}

	return nil
}
