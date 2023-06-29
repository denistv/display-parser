package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"display_parser/cmd/http/controllers"
	"display_parser/internal/config"
	"display_parser/internal/repository"
)

func main() {
	cfg := config.NewHTTPConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(config.UNIXDefaultErrCode)
	}

	rootCmd := newRootCommand(&cfg)
	if err = rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("executing command: %w", err))
	}

	if err := cfg.Validate(); err != nil {
		logger.Fatal(err.Error())
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DB.ConnString())
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to create connection pool: %w", err).Error())
	}
	defer dbpool.Close()

	modelsRepo := repository.NewModel(dbpool)
	modelsController := controllers.NewModelsController(logger, modelsRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	if cfg.CORSAllowedOrigin != "" {
		// Для корректной работы swagger-ui, который работает в браузере на другом порту
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{cfg.CORSAllowedOrigin},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders: []string{"Link"},
			MaxAge:         300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Get("/models", modelsController.ModelsIndex)

	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.ListenPort), r)
	if err != nil {
		panic(err)
	}
}
