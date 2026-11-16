package middlewares

import (
	"net/http"
	"testing"

	appmw "github.com/Sahil2k07/go-base/internal/middlewares"
	"github.com/Sahil2k07/go-base/tests/unit"
)

func TestAuthMiddleware(t *testing.T) {
	unit.LoadConfig(t)
	mw := appmw.AuthMiddleware()

	valid := unit.AccessToken(t, 1, 1, "a@x.com")
	expired := unit.ExpiredAccessToken(t, 1, 1, "a@x.com")
	refresh := unit.RefreshToken(t, 1, "11111111-1111-1111-1111-111111111111")

	tests := []struct {
		name   string
		header string
		want   int
	}{
		{"missing header", "", http.StatusUnauthorized},
		{"wrong scheme", "Basic " + valid, http.StatusUnauthorized},
		{"empty bearer", "Bearer ", http.StatusUnauthorized},
		{"garbage token", "Bearer not-a-token", http.StatusUnauthorized},
		{"tampered token", "Bearer " + valid[:len(valid)-2] + "xx", http.StatusUnauthorized},
		{"expired token", "Bearer " + expired, http.StatusUnauthorized},
		{"refresh token used as access token", "Bearer " + refresh, http.StatusUnauthorized},
		{"scheme is case insensitive", "bearer " + valid, http.StatusOK},
		{"valid token", "Bearer " + valid, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := unit.Serve(t, tt.header, mw)

			unit.AssertStatus(t, rec, tt.want)
		})
	}
}
