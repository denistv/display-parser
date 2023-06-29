package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"display_parser/internal/config"
)

// nolint
func newRootCommand(cfg *config.HTTPConfig) *cobra.Command {
	rootCmd := cobra.Command{}

	// Database
	rootCmd.PersistentFlags().StringVar(&cfg.DB.DBName, "db-name", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.User, "db-user", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.Password, "db-password", "", "")
	rootCmd.PersistentFlags().StringVar(&cfg.DB.Hostname, "db-hostname", "localhost", "")
	rootCmd.PersistentFlags().IntVar(&cfg.DB.Port, "db-port", 5432, "")
	rootCmd.PersistentFlags().IntVar(&cfg.DB.PoolMaxConns, "db-pool-max-conns", 3, "")

	rootCmd.PersistentFlags().IntVar(&cfg.ListenPort, "listen-port", 3000, "")
	rootCmd.PersistentFlags().StringVar(&cfg.CORSAllowedOrigin, "cors-allowed-origin", "", "")

	viper.AutomaticEnv()

	return &rootCmd
}
