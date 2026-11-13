package user

import "github.com/labstack/echo/v5"

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(g *echo.Group, userService *UserService) *UserHandler {
	h := &UserHandler{userService: userService}

	return h
}
