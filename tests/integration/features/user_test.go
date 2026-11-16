package features

import (
	"net/http"
	"testing"

	"github.com/Sahil2k07/go-base/internal/features/auth"
	"github.com/Sahil2k07/go-base/tests/integration"
)

const userMeURL = "/api/v1/user/me"

func TestMe(t *testing.T) {
	env := integration.Setup(t)

	// Create a user and obtain an access token.
	email := "user-me@example.com"
	password := "Password123!"

	status, body := env.Call(
		t,
		http.MethodPost,
		signupURL,
		"",
		auth.SignupRequest{
			Email:     email,
			FirstName: "John",
			LastName:  "Doe",
			Phone:     "+919876543210",
			Gender:    "Male",
			Password:  password,
		},
	)

	if status != http.StatusAccepted {
		t.Fatalf("failed to create user: status %d: %s", status, body)
	}

	tokens := integration.Decode[auth.TokenResponse](t, body)

	if tokens.AccessToken == "" {
		t.Fatal("signup returned an empty access token")
	}

	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "valid access token",
			token:      tokens.AccessToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing access token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid access token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status, body := env.Call(
				t,
				http.MethodGet,
				userMeURL,
				tc.token,
				nil,
			)

			if status != tc.wantStatus {
				t.Fatalf(
					"expected status %d, got %d: %s",
					tc.wantStatus,
					status,
					body,
				)
			}

			if status != http.StatusOK {
				return
			}

			response := integration.Decode[struct {
				ID          uint     `json:"id"`
				Email       string   `json:"email"`
				FirstName   string   `json:"firstName"`
				LastName    string   `json:"lastName"`
				Phone       string   `json:"phone"`
				Gender      string   `json:"gender"`
				Roles       []string `json:"roles"`
				Permissions []string `json:"permissions"`
			}](t, body)

			if response.ID == 0 {
				t.Error("expected a valid user ID")
			}

			if response.Email != email {
				t.Errorf("expected email %q, got %q", email, response.Email)
			}

			if response.FirstName != "John" {
				t.Errorf("expected first name John, got %q", response.FirstName)
			}

			if response.LastName != "Doe" {
				t.Errorf("expected last name Doe, got %q", response.LastName)
			}

			if response.Phone != "+919876543210" {
				t.Errorf("expected phone +919876543210, got %q", response.Phone)
			}

			if response.Gender != "Male" {
				t.Errorf("expected gender Male, got %q", response.Gender)
			}
		})
	}
}
