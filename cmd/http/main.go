package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	sqlxConn, err := sqlx.Connect("postgres", cfg.DB.ConnStringSQLX())
	if err != nil {
		logger.Fatal(err.Error())
	}

	modelsRepo := repository.NewModel(sqlxConn)
	modelsController := controllers.NewModelsController(logger, modelsRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/models", modelsController.ModelsIndex)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		panic(err)
	}
}
