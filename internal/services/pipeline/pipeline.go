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

type Pipeline struct {
	logger        *zap.Logger
	cfg           Cfg
	brandsColl    *BrandsCollector
	modelsURLColl *ModelsURLCollector
	pagesColl     *PagesCollector
	modelParser   *ModelParser
}

func (p *Pipeline) Run(ctx context.Context) {
	p.logger.Info("starting pipeline")

	// Pipeline chains
	brandURLsChan := p.brandsColl.Run(ctx)
	modelURLChan := p.modelsURLColl.Run(ctx, brandURLsChan)
	pageURLChan := p.pagesColl.Run(ctx, modelURLChan)

	for i := 0; i < p.cfg.ModelParserCount; i++ {
		p.logger.Info(fmt.Sprintf("starting model parser #%d of %d", i+1, p.cfg.ModelParserCount))
		p.modelParser.Run(ctx, pageURLChan)
	}
}
