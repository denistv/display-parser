package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"display_parser/internal/app"
	"display_parser/internal/repository"
	"display_parser/internal/services"
	"display_parser/internal/services/pipeline"
)

const defaultErrorCode = 255

func main() {
	cfg := app.NewConfig()
	rootCmd := newRootCommand(&cfg)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(defaultErrorCode)
	}

	// Создаем контекст с отменой для реализации graceful-shutdown и в дальнейшем передаем его в сервисы приложения.
	// Сервисы могут читать stop-канал в созданном контексте для корректного завершения своей работы, если это требует их реализация.
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,  // Прерывание приложения через Ctrl+C
		syscall.SIGTERM, // Общий сигнал завершения работы (посылаемый командой kill)
	)

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(defaultErrorCode)
	}

	conn, err := pgx.Connect(ctx, cfg.DB.DSN())
	if err != nil {
		logger.Fatal(fmt.Errorf("creating pq connector: %w", err).Error())
	}

	defer conn.Close(ctx)

	const dbDriver = "postgres"

	sqlxConn, err := sqlx.Connect(dbDriver, cfg.DB.NewSqlxDSN())
	if err != nil {
		logger.Fatal(err.Error())
	}

	goquDB := goqu.New(dbDriver, sqlxConn)

	// Repositories
	dbWrapper := repository.NewDBWrapper(conn)
	pageRepo := repository.NewPage(dbWrapper, goquDB)
	modelsRepo := repository.NewModel(goquDB)

	httpClient := services.NewDefaultHTTPClient(cfg.HTTP.Timeout)
	delayedHTTPClient := services.NewDelayedHTTPClient(ctx, cfg.HTTP.DelayPerRequest, httpClient)

	// Collectors
	brandsCollector := pipeline.NewBrandsCollector(logger, delayedHTTPClient, cancel)
	modelPagesCollector := pipeline.NewPagesCollector(logger, pageRepo, delayedHTTPClient, cfg.Pipeline.PageCollector)
	modelsURLCollector := pipeline.NewModelsURLCollector(logger, delayedHTTPClient)
	modelParser := pipeline.NewModelParser(logger, modelsRepo)

	pp := pipeline.NewPipeline(cfg.Pipeline, brandsCollector, modelPagesCollector, modelsURLCollector, modelParser, logger, pageRepo)
	pp.Run(ctx)

	<-ctx.Done()
}
