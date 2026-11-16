package middlewares

import (
	"github.com/Sahil2k07/go-base/internal/constants"
	app_errors "github.com/Sahil2k07/go-base/internal/errors"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/lib"
	"github.com/Sahil2k07/go-base/internal/store"
	"github.com/labstack/echo/v5"
)

func RequireRoles(roles ...string) echo.MiddlewareFunc {
	userRepo := store.NewUserRepository()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claims, ok := infrastructure.GetAuthClaims(c)
			if !ok {
				return lib.HandleError(c, app_errors.Unauthorized("Unauthorized"))
			}

			for _, role := range roles {
				has, err := userRepo.UserHasRole(claims.UserID, role)
				if err != nil {
					return lib.HandleError(c, app_errors.Internal(err))
				}

				if has {
					return next(c)
				}
			}

			return lib.HandleError(c, app_errors.Forbidden("You do not have the required role to access this resource"))
		}
	}
}

func RequirePermissions(permissions ...constants.Permission) echo.MiddlewareFunc {
	userRepo := store.NewUserRepository()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			claims, ok := infrastructure.GetAuthClaims(c)
			if !ok {
				return lib.HandleError(c, app_errors.Unauthorized("Unauthorized"))
			}

			for _, permission := range permissions {
				has, err := userRepo.UserHasPermission(claims.UserID, string(permission))
				if err != nil {
					return lib.HandleError(c, app_errors.Internal(err))
				}

				if !has {
					return lib.HandleError(c, app_errors.Forbidden("You do not have permission to access this resource"))
				}
			}

			return next(c)
		}
	}
}
