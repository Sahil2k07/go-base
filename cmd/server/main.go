package main

import (
	"log/slog"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
)

func init() {
	configs.LoadConfigs("configs/server.toml")
	database.ConnectPostgres()
}

func main() {
	slog.Info("Server hit...")
}
