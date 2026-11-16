package features

import (
	"net/http"
	"strings"
	"testing"

	"github.com/Sahil2k07/go-base/internal/features/auth"
	"github.com/Sahil2k07/go-base/tests/integration"
)

const (
	signupURL       = "/api/v1/auth/signup"
	signinURL       = "/api/v1/auth/signin"
	refreshTokenURL = "/api/v1/auth/refresh-token"
)

func TestSignup(t *testing.T) {
	env := integration.Setup(t)

	validRequest := auth.SignupRequest{
		Email:     "signup@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+919876543210",
		Gender:    "Male",
		Password:  "Password123!",
	}

	tests := []struct {
		name       string
		request    auth.SignupRequest
		setup      func()
		wantStatus int
		wantTokens bool
	}{
		{
			name:       "valid signup",
			request:    validRequest,
			wantStatus: http.StatusAccepted,
			wantTokens: true,
		},
		{
			name: "duplicate email",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Email = "duplicate@example.com"
				return r
			}(),
			setup: func() {
				status, body := env.Call(t, http.MethodPost, signupURL, "", auth.SignupRequest{
					Email:     "duplicate@example.com",
					FirstName: "John",
					LastName:  "Doe",
					Phone:     "+919876543210",
					Gender:    "Male",
					Password:  "Password123!",
				})
				if status != http.StatusAccepted {
					t.Fatalf("setup signup failed: status %d: %s", status, body)
				}
			},
			wantStatus: http.StatusConflict,
		},

		// Email validation
		{
			name: "missing email",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Email = ""
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Email = "invalid-email"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "email too long",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Email = strings.Repeat("a", 145) + "@example.com"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},

		// First name validation
		{
			name: "missing first name",
			request: func() auth.SignupRequest {
				r := validRequest
				r.FirstName = ""
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "first name too long",
			request: func() auth.SignupRequest {
				r := validRequest
				r.FirstName = strings.Repeat("a", 101)
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},

		// Last name validation
		{
			name: "missing last name",
			request: func() auth.SignupRequest {
				r := validRequest
				r.LastName = ""
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "last name too long",
			request: func() auth.SignupRequest {
				r := validRequest
				r.LastName = strings.Repeat("a", 101)
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},

		// Phone validation
		{
			name: "missing phone",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Phone = ""
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid phone",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Phone = "12345"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},

		// Gender validation
		{
			name: "missing gender",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Gender = ""
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "unsupported gender",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Gender = "Unknown"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},

		// Password validation
		{
			name: "missing password",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = ""
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password too short",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = "Ab1!"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password missing uppercase",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = "password123!"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password missing lowercase",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = "PASSWORD123!"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password missing digit",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = "PasswordABC!"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password missing symbol",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = "Password123"
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "password too long",
			request: func() auth.SignupRequest {
				r := validRequest
				r.Password = "Aa1!" + strings.Repeat("x", 69)
				return r
			}(),
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup()
			}

			status, body := env.Call(
				t,
				http.MethodPost,
				signupURL,
				"",
				tc.request,
			)

			if status != tc.wantStatus {
				t.Fatalf("expected status %d, got %d: %s",
					tc.wantStatus, status, body)
			}

			if tc.wantTokens {
				tokens := integration.Decode[auth.TokenResponse](t, body)

				if tokens.AccessToken == "" {
					t.Error("access token is empty")
				}

				if tokens.RefreshToken == "" {
					t.Error("refresh token is empty")
				}
			}
		})
	}
}

func TestSignin(t *testing.T) {
	env := integration.Setup(t)

	email := "signin@example.com"
	password := "Password123!"

	// Create a user for signin test cases.
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
		t.Fatalf("failed to create signin user: status %d: %s", status, body)
	}

	tests := []struct {
		name       string
		request    auth.SigninRequest
		wantStatus int
		wantTokens bool
	}{
		{
			name: "valid credentials",
			request: auth.SigninRequest{
				Email:    email,
				Password: password,
			},
			wantStatus: http.StatusOK,
			wantTokens: true,
		},
		{
			name: "unknown email",
			request: auth.SigninRequest{
				Email:    "unknown@example.com",
				Password: password,
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "incorrect password",
			request: auth.SigninRequest{
				Email:    email,
				Password: "WrongPassword123!",
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "missing email",
			request: auth.SigninRequest{
				Email:    "",
				Password: password,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			request: auth.SigninRequest{
				Email:    "invalid-email",
				Password: password,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing password",
			request: auth.SigninRequest{
				Email:    email,
				Password: "",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			status, body := env.Call(
				t,
				http.MethodPost,
				signinURL,
				"",
				tc.request,
			)

			if status != tc.wantStatus {
				t.Fatalf("expected status %d, got %d: %s",
					tc.wantStatus, status, body)
			}

			if tc.wantTokens {
				tokens := integration.Decode[auth.TokenResponse](t, body)

				if tokens.AccessToken == "" {
					t.Error("access token is empty")
				}

				if tokens.RefreshToken == "" {
					t.Error("refresh token is empty")
				}
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	env := integration.Setup(t)

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "valid refresh token",
			run: func(t *testing.T) {
				// Create a fresh user and session.
				status, body := env.Call(
					t,
					http.MethodPost,
					signupURL,
					"",
					auth.SignupRequest{
						Email:     "refresh-valid@example.com",
						FirstName: "John",
						LastName:  "Doe",
						Phone:     "+919876543210",
						Gender:    "Male",
						Password:  "Password123!",
					},
				)

				if status != http.StatusAccepted {
					t.Fatalf("signup failed: %d: %s", status, body)
				}

				tokens := integration.Decode[auth.TokenResponse](t, body)

				status, body = env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: tokens.RefreshToken,
					},
				)

				if status != http.StatusOK {
					t.Fatalf("expected %d, got %d: %s",
						http.StatusOK, status, body)
				}

				newTokens := integration.Decode[auth.TokenResponse](t, body)

				if newTokens.AccessToken == "" || newTokens.RefreshToken == "" {
					t.Fatal("expected both tokens to be returned")
				}

				if newTokens.RefreshToken == tokens.RefreshToken {
					t.Error("refresh token was not rotated")
				}

				if newTokens.AccessToken == tokens.AccessToken {
					t.Error("access token was not regenerated")
				}
			},
		},
		{
			name: "empty refresh token",
			run: func(t *testing.T) {
				status, body := env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{},
				)

				if status != http.StatusUnauthorized {
					t.Fatalf("expected %d, got %d: %s",
						http.StatusUnauthorized, status, body)
				}
			},
		},
		{
			name: "malformed refresh token",
			run: func(t *testing.T) {
				status, body := env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: "not-a-valid-token",
					},
				)

				if status != http.StatusUnauthorized {
					t.Fatalf("expected %d, got %d: %s",
						http.StatusUnauthorized, status, body)
				}
			},
		},
		{
			name: "reused refresh token revokes session",
			run: func(t *testing.T) {
				// Create a fresh user and session.
				status, body := env.Call(
					t,
					http.MethodPost,
					signupURL,
					"",
					auth.SignupRequest{
						Email:     "refresh-reuse@example.com",
						FirstName: "John",
						LastName:  "Doe",
						Phone:     "+919876543211",
						Gender:    "Male",
						Password:  "Password123!",
					},
				)

				if status != http.StatusAccepted {
					t.Fatalf("signup failed: %d: %s", status, body)
				}

				originalTokens := integration.Decode[auth.TokenResponse](t, body)

				// First use: should rotate successfully.
				status, body = env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: originalTokens.RefreshToken,
					},
				)

				if status != http.StatusOK {
					t.Fatalf("first refresh failed: %d: %s", status, body)
				}

				rotatedTokens := integration.Decode[auth.TokenResponse](t, body)

				// Reuse the old token: should revoke the session.
				status, body = env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: originalTokens.RefreshToken,
					},
				)

				if status != http.StatusUnauthorized {
					t.Fatalf("expected reused token to return %d, got %d: %s",
						http.StatusUnauthorized, status, body)
				}

				// The rotated token should now also fail because the session was revoked.
				status, body = env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: rotatedTokens.RefreshToken,
					},
				)

				if status != http.StatusUnauthorized {
					t.Fatalf("expected revoked session to return %d, got %d: %s",
						http.StatusUnauthorized, status, body)
				}
			},
		},
		{
			name: "new refresh token works after rotation",
			run: func(t *testing.T) {
				status, body := env.Call(
					t,
					http.MethodPost,
					signupURL,
					"",
					auth.SignupRequest{
						Email:     "refresh-rotation@example.com",
						FirstName: "John",
						LastName:  "Doe",
						Phone:     "+919876543212",
						Gender:    "Male",
						Password:  "Password123!",
					},
				)

				if status != http.StatusAccepted {
					t.Fatalf("signup failed: %d: %s", status, body)
				}

				tokens := integration.Decode[auth.TokenResponse](t, body)

				// Rotate once.
				status, body = env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: tokens.RefreshToken,
					},
				)

				if status != http.StatusOK {
					t.Fatalf("first refresh failed: %d: %s", status, body)
				}

				newTokens := integration.Decode[auth.TokenResponse](t, body)

				// Use the new token.
				status, body = env.Call(
					t,
					http.MethodPost,
					refreshTokenURL,
					"",
					auth.RefreshTokenRequest{
						RefreshToken: newTokens.RefreshToken,
					},
				)

				if status != http.StatusOK {
					t.Fatalf("second refresh failed: %d: %s", status, body)
				}

				rotatedTokens := integration.Decode[auth.TokenResponse](t, body)

				if rotatedTokens.RefreshToken == newTokens.RefreshToken {
					t.Error("expected the second refresh token to be rotated")
				}

				if rotatedTokens.AccessToken == "" || rotatedTokens.RefreshToken == "" {
					t.Error("expected both tokens to be returned")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.run)
	}
}
