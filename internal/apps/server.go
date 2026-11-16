package apps

import (
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/features/auth"
	"github.com/Sahil2k07/go-base/internal/features/user"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/middlewares"
	"github.com/Sahil2k07/go-base/internal/store"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func bootstrapServer(appConfig configs.GlobalConfig, publicGroup, privateGroup *echo.Group) {
	// Repositories
	userRepo := store.NewUserRepository()
	userSessionRepo := store.NewUserSessionRepository()

	// Services
	cryptoService := infrastructure.NewCryptoService(appConfig.Server)
	authService := auth.NewAuthService(userRepo, userSessionRepo, cryptoService)
	userService := user.NewUserService(userRepo)

	// Handlers
	auth.NewAuthHandler(publicGroup, authService)
	user.NewUserHandler(privateGroup, userService)
}

func NewServer(appConfig configs.GlobalConfig) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     appConfig.Server.Origins,
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	publicGroup := e.Group("/api")

	privateGroup := e.Group("/api")
	privateGroup.Use(middlewares.AuthMiddleware())

	bootstrapServer(appConfig, publicGroup, privateGroup)

	return e
}
