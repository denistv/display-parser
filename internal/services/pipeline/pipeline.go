package pipeline

import (
	"context"
	"display_parser/internal/domain"
	"display_parser/internal/repository"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Cfg struct {
	ModelParserCount int
	PageCollector    PagesCollectorCfg
}

func (c *Cfg) Validate() error {
	if c.ModelParserCount <= 0 {
		return errors.New("parser count must be leather than 0")
	}

	if err := c.PageCollector.Validate(); err != nil {
		return fmt.Errorf("validating page collector config: %w", err)
	}

	return nil
}

func NewPipeline(cfg Cfg, brandsCollector *BrandsCollector, pagesColl *PageCollector, modelURLColl *ModelsURLCollector, modelParser *ModelParser, logger *zap.Logger, pageRepo *repository.Page, modelPersister *ModelPersister) *Pipeline {
	return &Pipeline{
		cfg:           cfg,
		logger:        logger,
		brandsColl:    brandsCollector,
		modelsURLColl: modelURLColl,
		pagesColl:     pagesColl,
		modelParser:   modelParser,
		pageRepo:      pageRepo,
		modelPersister: modelPersister,
	}
}

// Pipeline представляет собой сущность, которая связывает шаги пайплайна и централизовано управляет его жизненным циклом.
type Pipeline struct {
	logger        *zap.Logger
	cfg           Cfg
	brandsColl    *BrandsCollector
	modelsURLColl *ModelsURLCollector
	pagesColl     *PageCollector
	modelParser   *ModelParser
	modelPersister *ModelPersister
	pageRepo      *repository.Page
}

// Run связывает этапы пайплайна и запускает его.
// В зависимости от настройки UseStoredPagesCacheOnly конфигурируются требуемые шаги.
func (p *Pipeline) Run(ctx context.Context) chan struct{}{
	p.logger.Info("starting pipeline")

	pageCh := make([]<-chan domain.PageEntity, 0, 0)

	if p.cfg.PageCollector.UseStoredPagesOnly {
		// используем кэш страниц в базе. Подходит для второго и последующих запусков или когда у сущности модели
		// появился новый параметр, который необходимо быстро перепарсить без хождения в сеть
		pageCh = p.loadPagesFromCache(ctx)
	} else {
		//В этом случае не используем кэш страниц, хранящийся в базе и получаем все данные из интернета.
		// Подходит для первого запуска.
		brandURLsChan := p.brandsColl.Run(ctx)
		modelURLChan := p.modelsURLColl.Run(ctx, brandURLsChan)
		pageCh = p.spawnPageCollectors(ctx, modelURLChan)
	}

	modelsCh := p.spawnParsers(ctx, mergeCh(pageCh...))
	done := p.modelPersister.Run(ctx, mergeCh(modelsCh...))

	return done
}

// spawnPageCollectors запускает нужное число воркеров PageCollector
func (p *Pipeline) spawnPageCollectors(ctx context.Context, modelURLChan <-chan string) []<-chan domain.PageEntity {
	pageCh := make([]<-chan domain.PageEntity, 0, p.cfg.ModelParserCount)

	for i := 0; i < p.cfg.PageCollector.Count; i++ {
		p.logger.Info(fmt.Sprintf("starting model parser #%d of %d", i+1, p.cfg.ModelParserCount))

		ch := p.pagesColl.Run(ctx, modelURLChan)
		pageCh = append(pageCh, ch)
	}

	return pageCh
}

// spawnParsers запускает нужное число ModelParser
func (p *Pipeline) spawnParsers(ctx context.Context, pageCh chan domain.PageEntity) []<-chan domain.ModelEntity {
	modelsCh := make([]<-chan domain.ModelEntity, 0, p.cfg.ModelParserCount)

	for i := 0; i < p.cfg.ModelParserCount; i++ {
		p.logger.Info(fmt.Sprintf("starting model parser #%d of %d", i+1, p.cfg.ModelParserCount))

		ch := p.modelParser.Run(ctx, pageCh)
		modelsCh = append(modelsCh, ch)
	}

	return modelsCh
}

func mergeCh[T any](in ...<-chan T) chan T {
	out := make(chan T, len(in))
	var wg sync.WaitGroup
	wg.Add(len(in))

	// стартуем горутины, которые читают из входных каналов и пересылают результат в один выходной
	for _, c := range in {
		go func(c <-chan T){
			for v := range c {
				out <-v
			}

			wg.Done()
		}(c)
	}

	go func() {
		// Закрываем выходной канал в том случае, если все входные так же закрыты
		wg.Wait()
		close(out)
	}()

	return out
}

func (p *Pipeline) loadPagesFromCache(ctx context.Context) []<-chan domain.PageEntity {
	out := make(chan domain.PageEntity)

	go func() {
		pages, err := p.pageRepo.All(ctx) // помним про то, что в настоящем проекте так делать не нужно
		if err != nil {
			p.logger.Error(err.Error())
			//todo cancel
		}

		for _, page := range pages {
			out <- page
		}

		close(out)
	}()

	return []<-chan domain.PageEntity{out}
}
