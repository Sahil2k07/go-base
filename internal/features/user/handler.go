package user

import (
	"net/http"

	app_errors "github.com/Sahil2k07/go-base/internal/errors"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/lib"
	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(g *echo.Group, userService *UserService) *UserHandler {
	h := &UserHandler{userService: userService}

	g.GET("/v1/user/me", h.Me)

	return h
}

func (h *UserHandler) Me(c *echo.Context) error {
	claims, ok := infrastructure.GetAuthClaims(c)
	if !ok {
		return lib.HandleError(c, app_errors.Unauthorized("User is not Authorized"))
	}

	userData, err := h.userService.GetCurrentUserInformation(claims.UserID)
	if err != nil {
		return lib.HandleError(c, err)
	}

	return c.JSON(http.StatusOK, userData)
}
