package pipeline

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type Cfg struct {
	ModelParserCount int
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
}

// Run связывает этапы пайплайна и запускает его
func (p *Pipeline) Run(ctx context.Context) {
	p.logger.Info("starting pipeline")

	brandURLsChan := p.brandsColl.Run(ctx)
	modelURLChan := p.modelsURLColl.Run(ctx, brandURLsChan)
	pageURLChan := p.pagesColl.Run(ctx, modelURLChan)

	// Запускаем требуемое число парсеров. С практической точки зрения, в данной задаче запускать большое число парсеров
	// на небольших наборах данных особого смысла не имеет.
	// Просто для демонстрации паралеллизма. Здесь пайплайн ветвится -- канал с URL страниц читает множество парсеров
	// TODO сделать объединение результатов работы этапа парсинга в этап сохранения модели в рамках отдельного шага
	for i := 0; i < p.cfg.ModelParserCount; i++ {
		p.logger.Info(fmt.Sprintf("starting model parser #%d of %d", i+1, p.cfg.ModelParserCount))
		p.modelParser.Run(ctx, pageURLChan)
	}
}
