package apps

import (
	"github.com/Sahil2k07/go-base/internal/features/auth"
	"github.com/Sahil2k07/go-base/internal/store"
	"github.com/labstack/echo/v5"
)

func BootstrapServer(publicGroup, privateGroup *echo.Group) {
	// Repositories
	authRepo := store.NewAuthStore()

	// Services
	authService := auth.NewAuthService(authRepo)

	// Handlers
	auth.NewAuthHandler(publicGroup, authService)
}
