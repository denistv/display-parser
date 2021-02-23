package main

import (
	"context"
	"displayCrawler/internal/respository"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"

	"displayCrawler/internal/pipeline"
)
import "go.uber.org/zap"

func main() {
	ctx := context.Background()

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
	}

	conn, err := pgx.Connect(ctx, "postgres://displayspecs:displayspecs@localhost:49153/displayspecs")
	if err != nil {
		logger.Fatal(fmt.Errorf("creating pq connector: %w", err).Error())
	}

	defer conn.Close(ctx)

	// Repositories
	dbWrapper := respository.NewDBWrapper(conn)
	modelDocRepo := respository.NewModelDocument(dbWrapper)

	// Collectors
	brandsCollector := pipeline.NewBrandsCollector(logger)
	modelDocsColector := pipeline.NewModelDocumentCollector(logger, modelDocRepo)
	modelsCollector := pipeline.NewModelsCollector(logger)

	// Pipeline chains
	brandURLsChan := brandsCollector.BrandURLs()
	modelsIndexURLsChan := modelsCollector.GetItemsIndex(brandURLsChan)
	modelDocsChan := modelDocsColector.Collect(modelsIndexURLsChan)

	for modelDoc := range modelDocsChan {
		err := modelDocRepo.Create(ctx, modelDoc)
		if err != nil {
			logger.Error(fmt.Errorf("creating document: %w", err).Error())
			continue
		}

		logger.Info("document inserted: " + modelDoc.URL)
	}
}
