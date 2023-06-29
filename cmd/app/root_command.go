package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"display_parser/internal/config"
)

// В этом проекте не задействуется вся функциональность cobra.
// На данном этапе мне достаточно возможности удобной работы с параметрами (отображение справки, парсинг флагов из CLI и ENV)
// nolint
func newRootCommand(cfg *config.AppConfig) *cobra.Command {
	rootCmd := cobra.Command{}

	// Common flags
	rootCmd.PersistentFlags().DurationVar(&cfg.HTTP.DelayPerRequest, "http-delay-per-request", 2000*time.Millisecond, "use golang time.Duration string format. Example: 1m30s500ms")
	rootCmd.PersistentFlags().DurationVar(&cfg.HTTP.Timeout, "http-timeout", 10*time.Second, "use golang time.Duration string format. Example: 1m30s500ms")

	// Pipeline
	rootCmd.PersistentFlags().IntVar(&cfg.Pipeline.ModelParserCount, "pipeline-model-parser-count", 1, "")
	rootCmd.PersistentFlags().BoolVar(&cfg.Pipeline.PageCollector.UseStoredPagesOnly, "pipeline-use-stored-pages-only", false, "use for rebuild database models only. If this flag enabled, parser will not going to site and using db-cache.")
	rootCmd.PersistentFlags().IntVar(&cfg.Pipeline.PageCollector.Count, "pipeline-page-collector-count", 1, "")

	// Database
	rootCmd.PersistentFlags().StringVar(&cfg.DB.DBName, "db-name", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.User, "db-user", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.Password, "db-password", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.Hostname, "db-hostname", "localhost", "")
	rootCmd.PersistentFlags().IntVar(&cfg.DB.Port, "db-port", 5432, "")
	rootCmd.PersistentFlags().IntVar(&cfg.DB.PoolMaxConns, "db-pool-max-conns", 3, "")

	viper.AutomaticEnv()

	return &rootCmd
}
