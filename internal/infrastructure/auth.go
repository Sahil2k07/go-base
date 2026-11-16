package infrastructure

import (
	"github.com/labstack/echo/v5"
)

const AuthClaimsKey = "auth_claims"

func GetAuthClaims(c *echo.Context) (*AccessTokenClaims, bool) {
	claims, ok := c.Get(AuthClaimsKey).(*AccessTokenClaims)
	return claims, ok && claims != nil
}
