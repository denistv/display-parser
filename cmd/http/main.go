package main

import (
	"display_parser/cmd/http/controllers"
	"display_parser/internal/app"
	"display_parser/internal/repository"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg := app.NewConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(255)
	}

	rootCmd := newRootCommand(&cfg)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(fmt.Errorf("executing command: %w", err))
	}

	sqlxConn, err := sqlx.Connect("postgres", cfg.DB.ConnStringSQLX())
	if err != nil {
		logger.Fatal(err.Error())
	}

	goquDB := goqu.New("postgres", sqlxConn)
	modelsRepo := repository.NewModel(goquDB)
	modelsController := controllers.NewModelsController(logger, modelsRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/models", modelsController.ModelsIndex)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}

