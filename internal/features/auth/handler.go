package auth

import (
	"net/http"

	"github.com/Sahil2k07/go-base/internal/lib"
	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authService *AuthService
}

func NewAuthHandler(g *echo.Group, authService *AuthService) *AuthHandler {
	h := &AuthHandler{authService: authService}

	g.POST("/v1/auth/signup", h.SignupUser)
	g.POST("/v1/auth/signin", h.SigninUser)
	g.POST("/v1/auth/refresh-token", h.RefreshToken)

	return h
}

func (h *AuthHandler) SignupUser(c *echo.Context) error {
	var req SignupRequest

	if err := c.Bind(&req); err != nil {
		return lib.HandleError(c, err)
	}

	if err := req.ValidateRequest(); err != nil {
		return lib.HandleError(c, err)
	}

	resp, err := h.authService.SignupUser(req)
	if err != nil {
		return lib.HandleError(c, err)
	}

	return c.JSON(http.StatusAccepted, resp)
}

func (h *AuthHandler) SigninUser(c *echo.Context) error {
	var req SigninRequest

	if err := c.Bind(&req); err != nil {
		return lib.HandleError(c, err)
	}

	if err := req.ValidateRequest(); err != nil {
		return lib.HandleError(c, err)
	}

	resp, err := h.authService.SigninUser(req)
	if err != nil {
		return lib.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(c *echo.Context) error {
	var req RefreshTokenRequest

	if err := c.Bind(&req); err != nil {
		return lib.HandleError(c, err)
	}

	resp, err := h.authService.RefreshToken(req)
	if err != nil {
		return lib.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}
