package main

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"

	"displayCrawler/internal/pipeline"
	"displayCrawler/internal/repository"
)

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
	}

	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@192.168.1.72:5432/display_crawler")
	if err != nil {
		logger.Fatal(fmt.Errorf("creating pq connector: %w", err).Error())
	}

	defer conn.Close(ctx)

	sqlxConn, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=display_crawler host=192.168.1.72 sslmode=disable")
	if err != nil {
		logger.Fatal(err.Error())
	}

	goquDb := goqu.New("postgres", sqlxConn)

	// Repositories
	dbWrapper := repository.NewDBWrapper(conn)
	docRepo := repository.NewDocument(dbWrapper, goquDb)
	modelsRepo := repository.NewModel(goquDb)

	// Collectors
	brandsCollector := pipeline.NewBrandsCollector(logger)
	docsCollector := pipeline.NewDocumentsCollector(logger, docRepo)
	modelsURLCollector := pipeline.NewModelsURLCollector(logger)
	modelParser := pipeline.NewModelParser(logger, modelsRepo)

	// Pipeline chains
	brandURLsChan := brandsCollector.Run()
	modelsIndexURLsChan := modelsURLCollector.Run(brandURLsChan)
	documentsChan := docsCollector.Run(modelsIndexURLsChan)
	modelParser.Run(documentsChan)

	<-ctx.Done()
}
