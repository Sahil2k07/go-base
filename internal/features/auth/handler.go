package auth

import "github.com/labstack/echo/v5"

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(g *echo.Group, authService *AuthService) *AuthHandler {
	h := &AuthHandler{authService: authService}

	return h
}
