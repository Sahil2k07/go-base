package main

import (
	"log/slog"

	"github.com/Sahil2k07/go-base/internal/apps"
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
)

func init() {
	configs.LoadConfigs("configs/migration.toml")
	infrastructure.ConfigureLogging("migration")
	database.ConnectPostgres()
}

func main() {
	err := apps.RunMigrations()
	if err != nil {
		slog.Error("Migrations failed", "error", err)
		return
	}

	slog.Info("Migrations Applied Successfully")
}
