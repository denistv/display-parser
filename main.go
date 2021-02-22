package main

import (
	"context"
	 "displayCrawler/internal/storage"
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

	conn, err := pgx.Connect(ctx, "storage://displayspec:displayspec@localhost:5432/displayspec")
	if err != nil {
		logger.Fatal(fmt.Errorf("creating pq connector: %w", err).Error())
	}

	defer conn.Close(ctx)

	stor := storage.NewStorage(conn)

	brandsCollector := pipeline.NewBrandsCollector(logger)
	devicesCollector := pipeline.NewDevicesCollector(logger)
	saver := pipeline.NewSaver(logger, stor)

	brandsChan, err := brandsCollector.Run()
	if err != nil {
		logger.Fatal(err.Error())
	}

	devicesChan := devicesCollector.ItemsIndex(brandsChan)
	saver.Persist(ctx, devicesChan)
}
