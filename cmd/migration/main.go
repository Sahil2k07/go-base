package main

import (
	"log/slog"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/migrations"
)

func init() {
	infrastructure.ConfigureLogging("migration")
	configs.LoadConfigs("configs/migration.toml")
	database.ConnectPostgres()
}

func main() {
	err := migrations.RunMigrations()
	if err != nil {
		slog.Error("Migrations failed", "error", err)
		return
	}

	slog.Info("Migrations Applied Successfully")
}
