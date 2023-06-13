package pipeline

import (
	"context"
	"display_parser/internal/domain"
	"display_parser/internal/repository"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type Cfg struct {
	// Пересобрать модели на основе кэша страниц в базе. Если флаг взведен, не ходим во внешний сервис для сбора данных и используем имеющийся кэш страниц в БД.
	// Полезно в тех случаях, когда сайт спаршен (данные страниц сохранены в кэше в таблице pages, но сущность модели расширена дополнительным полем.
	// Чтобы не собирать все данные по новой через сеть, используем сохраненные в базу страницы и перераспаршиваем их, обновляя сущности моделей).
	UseStoredPagesCacheOnly bool
	ModelParserCount        int
}

func (c *Cfg) Validate() error {
	if c.ModelParserCount <= 0 {
		return errors.New("parser count must be leather than 0")
	}

	return nil
}

func NewPipeline(cfg Cfg, brandsCollector *BrandsCollector, pagesColl *PagesCollector, modelURLColl *ModelsURLCollector, modelParser *ModelParser, logger *zap.Logger) *Pipeline {
	return &Pipeline{
		cfg:           cfg,
		logger:        logger,
		brandsColl:    brandsCollector,
		modelsURLColl: modelURLColl,
		pagesColl:     pagesColl,
		modelParser:   modelParser,
	}
}

// Pipeline представляет собой сущность, которая связывает шаги пайплайна и централизовано управляет его жизненным циклом.
type Pipeline struct {
	logger        *zap.Logger
	cfg           Cfg
	brandsColl    *BrandsCollector
	modelsURLColl *ModelsURLCollector
	pagesColl     *PagesCollector
	modelParser   *ModelParser
	pageRepo      *repository.Page
}

// Run связывает этапы пайплайна и запускает его.
// В зависимости от настройки UseStoredPagesCacheOnly конфигурируются требуемые шаги.
func (p *Pipeline) Run(ctx context.Context) {
	p.logger.Info("starting pipeline")

	pageChan := make(chan domain.PageEntity)

	if p.cfg.UseStoredPagesCacheOnly {
		// используем кэш страниц в базе. Подходит для второго и последующих запусков или когда у сущности модели
		// появился новый параметр, который необходимо быстро перепарсить без хождения в сеть
		p.loadPagesFromCache(ctx, pageChan)
	} else {
		//В этом случае не используем кэш страниц, хранящийся в базе и получаем все данные из интернета.
		// Подходит для первого запуска.
		brandURLsChan := p.brandsColl.Run(ctx)
		modelURLChan := p.modelsURLColl.Run(ctx, brandURLsChan)
		pageChan = p.pagesColl.Run(ctx, modelURLChan)
	}

	// Запускаем требуемое число парсеров. С практической точки зрения, в данной задаче запускать большое число парсеров
	// на небольших наборах данных особого смысла не имеет.
	// Просто для демонстрации паралеллизма. Здесь пайплайн ветвится -- канал с URL страниц читает множество парсеров
	// TODO сделать объединение результатов работы этапа парсинга в этап сохранения модели в рамках отдельного шага
	for i := 0; i < p.cfg.ModelParserCount; i++ {
		p.logger.Info(fmt.Sprintf("starting model parser #%d of %d", i+1, p.cfg.ModelParserCount))
		p.modelParser.Run(ctx, pageChan)
	}
}

func (p *Pipeline) loadPagesFromCache(ctx context.Context, pageChan chan domain.PageEntity) {
	pages, err := p.pageRepo.All(ctx)
	if err != nil {
		p.logger.Error(err.Error())
		//todo cancel
	}

	for _, page := range pages {
		pageChan <- page
	}
}
