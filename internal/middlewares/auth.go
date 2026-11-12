package middlewares

import (
	"net/http"
	"strings"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/labstack/echo/v5"
)

func AuthMiddleware() echo.MiddlewareFunc {
	cryptoService := infrastructure.NewCryptoService(configs.GetServerConfig())

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			// Check whether the Authorization header exists.
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"code":    http.StatusUnauthorized,
					"message": "Authorization header is required",
				})
			}

			// Extract the Bearer token.
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"code":    http.StatusUnauthorized,
					"message": "Invalid authorization header",
				})
			}

			token := strings.TrimSpace(parts[1])

			// Verify the access token.
			claims, err := cryptoService.VerifyAccessToken(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]any{
					"success": false,
					"code":    http.StatusUnauthorized,
					"message": "Invalid or expired access token",
				})
			}

			// Make verified claims available to downstream handlers.
			c.Set("auth_claims", claims)

			return next(c)
		}
	}
}
