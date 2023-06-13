package main

import (
	"context"
	"display_parser/internal/services"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"display_parser/internal/app"
	"display_parser/internal/repository"
	"display_parser/internal/services/pipeline"
)

const defaultErrorCode = 255

func newRootCommand(cfg *app.Config) *cobra.Command {
	rootCmd := cobra.Command{}

	// Common flags
	rootCmd.PersistentFlags().DurationVar(&cfg.HTTP.DelayPerRequest, "http-delay-per-request", 1000*time.Millisecond, "use golang time.Duration string format. Example: 1m30s500ms")
	rootCmd.PersistentFlags().DurationVar(&cfg.HTTP.Timeout, "http-timeout", 10*time.Second, "use golang time.Duration string format. Example: 1m30s500ms")

	rootCmd.PersistentFlags().IntVar(&cfg.Pipeline.ModelParserCount, "pipeline-model-parser-count", 5, "")

	// Database
	rootCmd.PersistentFlags().StringVar(&cfg.DB.DBName, "db-name", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.User, "db-user", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.Password, "db-password", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.Hostname, "db-hostname", "localhost", "")
	rootCmd.PersistentFlags().IntVar(&cfg.DB.Port, "db-port", 5432, "")

	return &rootCmd
}

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
	ctx, _ := signal.NotifyContext(
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
	brandsCollector := pipeline.NewBrandsCollector(logger, delayedHTTPClient)
	modelPagesCollector := pipeline.NewPagesCollector(logger, pageRepo, delayedHTTPClient)
	modelsURLCollector := pipeline.NewModelsURLCollector(logger, delayedHTTPClient)
	modelParser := pipeline.NewModelParser(logger, modelsRepo)

	pp := pipeline.NewPipeline(cfg.Pipeline, brandsCollector, modelPagesCollector, modelsURLCollector, modelParser, logger)
	pp.Run(ctx)

	<-ctx.Done()
}
