package main

import (
	"github.com/Sahil2k07/go-base/internal/apps"
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
)

func init() {
	configs.LoadConfigs("configs/migration.toml")
	infrastructure.ConfigureLogging("gen")
	database.ConnectPostgres()
}

func main() {
	apps.GenerateGenModels()
}
