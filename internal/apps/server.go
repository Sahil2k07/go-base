package apps

import (
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/features/auth"
	"github.com/Sahil2k07/go-base/internal/features/user"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/store"
	"github.com/labstack/echo/v5"
)

func BootstrapServer(appConfig configs.GlobalConfig, publicGroup, privateGroup *echo.Group) {
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
