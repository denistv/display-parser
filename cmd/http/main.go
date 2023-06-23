package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"display_parser/cmd/http/controllers"
	"display_parser/internal/app"
	"display_parser/internal/repository"
)

func main() {
	cfg := app.NewConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(app.UnexpectedErrCode)
	}

	rootCmd := newRootCommand(&cfg)
	if err = rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("executing command: %w", err))
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DB.ConnString())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	modelsRepo := repository.NewModel(dbpool)
	modelsController := controllers.NewModelsController(logger, modelsRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/models", modelsController.ModelsIndex)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
