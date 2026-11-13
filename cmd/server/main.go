package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/Sahil2k07/go-base/internal/apps"
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/middlewares"
)

var appConfig configs.GlobalConfig

func init() {
	appConfig = configs.LoadConfigs("configs/server.toml")
	database.ConnectPostgres()
	infrastructure.ConfigureLogging("server")
}

func main() {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     appConfig.Server.Origins,
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	publicGroup := e.Group("/")

	privateGroup := e.Group("/api")
	privateGroup.Use(middlewares.AuthMiddleware())

	apps.BootstrapServer(publicGroup, privateGroup)

	if err := e.Start(appConfig.Server.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
