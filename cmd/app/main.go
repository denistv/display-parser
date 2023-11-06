package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"display_parser/internal/config"
	"display_parser/internal/repository"
	"display_parser/internal/services"
	"display_parser/internal/services/pipeline"
	"display_parser/pkg/logger"
)

func main() {
	cfg := config.NewCmdApp()
	rootCmd := newRootCommand(&cfg)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(config.UNIXDefaultErrCode)
	}

	// Создаем контекст с отменой для реализации graceful-shutdown и в дальнейшем передаем его в сервисы приложения.
	// Сервисы могут читать stop-канал в созданном контексте для корректного завершения своей работы, если это требует их реализация.
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,  // Прерывание приложения через Ctrl+C
		syscall.SIGTERM, // Общий сигнал завершения работы (посылаемый командой kill)
	)

	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(config.UNIXDefaultErrCode)
	}

	wrappedLogger := logger.NewZapWrapper(zapLogger)

	dbpool, err := pgxpool.New(context.Background(), cfg.DB.ConnString())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// Repositories
	pageRepo := repository.NewPage(dbpool)
	modelsRepo := repository.NewModel(dbpool)

	httpClient := services.NewDefaultHTTPClient(cfg.HTTP.Timeout)
	delayedHTTPClient := services.NewDelayedHTTPClient(ctx, cfg.HTTP.DelayPerRequest, httpClient)

	// Collectors
	brandsCollector := pipeline.NewBrandsCollector(wrappedLogger, delayedHTTPClient, cancel)
	modelPagesCollector := pipeline.NewPageCollector(wrappedLogger, pageRepo, delayedHTTPClient, cfg.Pipeline.PageCollector)
	modelsURLCollector := pipeline.NewModelsURLCollector(wrappedLogger, delayedHTTPClient)
	modelParser := pipeline.NewModelParser(wrappedLogger, modelsRepo)
	modelPersister := pipeline.NewModelPersister(wrappedLogger, modelsRepo)

	pp := pipeline.NewPipeline(
		cfg.Pipeline,
		brandsCollector,
		modelPagesCollector,
		modelsURLCollector,
		modelParser,
		wrappedLogger,
		pageRepo,
		modelPersister,
	)

	<-pp.Run(ctx)

	zapLogger.Info("exiting")
}
