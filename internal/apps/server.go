package apps

import (
	"github.com/Sahil2k07/go-base/internal/features/auth"
	"github.com/Sahil2k07/go-base/internal/features/user"
	"github.com/Sahil2k07/go-base/internal/store"
	"github.com/labstack/echo/v5"
)

func BootstrapServer(publicGroup, privateGroup *echo.Group) {
	// Repositories
	authRepo := store.NewAuthRepository()
	userRepo := store.NewAuthRepository()

	// Services
	authService := auth.NewAuthService(authRepo)
	userService := user.NewUserService(userRepo)

	// Handlers
	auth.NewAuthHandler(publicGroup, authService)
	user.NewUserHandler(privateGroup, userService)
}
