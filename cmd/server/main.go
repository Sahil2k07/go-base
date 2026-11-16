package main

import (
	"github.com/Sahil2k07/go-base/internal/apps"
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
)

var appConfig configs.GlobalConfig

func init() {
	appConfig = configs.LoadConfigs("configs/server.toml")
	infrastructure.ConfigureLogging("server")
	database.ConnectPostgres()
}

func main() {
	e := apps.NewServer(appConfig)

	if err := e.Start(appConfig.Server.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
