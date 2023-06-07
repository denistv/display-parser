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
	"display_parser/internal/services/pipeline"
)

func main() {
	// Создаем контекст с отменой для реализации gracefull-shutdown и в дальнейшем передаем его в сервисы приложения.
	// Сервисы могут читать stop-канал в созданном контексте для корректного завершения своей работы, если это требует их реализация.
	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,  // Прерывание приложения через Ctrl+C
		syscall.SIGTERM, // Общий сигнал завершения работы (посылаемый командой kill)
	)
	cfg := app.NewConfigDev()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
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

	goquDb := goqu.New(dbDriver, sqlxConn)

	// Repositories
	dbWrapper := repository.NewDBWrapper(conn)
	pageRepo := repository.NewPage(dbWrapper, goquDb)
	modelsRepo := repository.NewModel(goquDb)

	// Collectors
	brandsCollector := pipeline.NewBrandsCollector(logger)
	modelPagesCollector := pipeline.NewPagesCollector(logger, pageRepo)
	modelsURLCollector := pipeline.NewModelsURLCollector(logger)
	modelParser := pipeline.NewModelParser(logger, modelsRepo)

	// Pipeline chains
	brandURLsChan := brandsCollector.Run(ctx)
	modelsIndexURLsChan := modelsURLCollector.Run(ctx, brandURLsChan)
	pagesChan := modelPagesCollector.Run(ctx, modelsIndexURLsChan)

	modelParser.Run(ctx, pagesChan)
	modelParser.Run(ctx, pagesChan)
	modelParser.Run(ctx, pagesChan)
	modelParser.Run(ctx, pagesChan)

	<-ctx.Done()
}
