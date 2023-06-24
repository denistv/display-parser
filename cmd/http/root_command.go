package main

import (
	"display_parser/internal/app"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nolint
func newRootCommand(cfg *app.Config) *cobra.Command {
	rootCmd := cobra.Command{}

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
