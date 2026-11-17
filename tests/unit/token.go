package unit

import (
	"testing"
	"time"

	"github.com/Sahil2k07/go-base/internal/infrastructure"
)

func AccessToken(t *testing.T, userID, profileID uint, email string) string {
	t.Helper()

	cs := infrastructure.NewCryptoService(LoadConfig(t).Server)

	token, err := cs.GenerateAccessToken(userID, profileID, email)
	if err != nil {
		t.Fatalf("generate access token: %v", err)
	}

	return token
}

// ExpiredAccessToken is signed with the real secrets but expired an hour ago.
func ExpiredAccessToken(t *testing.T, userID, profileID uint, email string) string {
	t.Helper()

	server := LoadConfig(t).Server
	server.AccessTokenExpiry = -time.Hour

	token, err := infrastructure.NewCryptoService(server).GenerateAccessToken(userID, profileID, email)
	if err != nil {
		t.Fatalf("generate expired access token: %v", err)
	}

	return token
}

func RefreshToken(t *testing.T, userID uint, sessionID string) string {
	t.Helper()

	cs := infrastructure.NewCryptoService(LoadConfig(t).Server)

	token, err := cs.GenerateRefreshToken(userID, sessionID)
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}

	return token
}
