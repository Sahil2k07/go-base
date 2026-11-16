package auth

import (
	"regexp"
	"slices"
	"strings"
	"unicode"

	app_errors "github.com/Sahil2k07/go-base/internal/errors"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9][0-9]{9,14}$`)

	allowedGenders = []string{"Male", "Female", "Trans", "Binary", "Non-Binary", "Refused"}
)

type SignupRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Gender    string `json:"gender"`
	Password  string `json:"password"`
}

func (r *SignupRequest) ValidateRequest() error {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.FirstName = strings.TrimSpace(r.FirstName)
	r.LastName = strings.TrimSpace(r.LastName)
	r.Phone = strings.TrimSpace(r.Phone)

	if err := validateEmail(r.Email); err != nil {
		return err
	}

	if r.FirstName == "" || len(r.FirstName) > 100 {
		return app_errors.BadRequest("First name is required and must be at most 100 characters")
	}

	if r.LastName == "" || len(r.LastName) > 100 {
		return app_errors.BadRequest("Last name is required and must be at most 100 characters")
	}

	if !phoneRegex.MatchString(r.Phone) {
		return app_errors.BadRequest("Phone must be 10 to 15 digits with an optional leading +")
	}

	if !slices.Contains(allowedGenders, r.Gender) {
		return app_errors.BadRequest("Gender must be one of: " + strings.Join(allowedGenders, ", "))
	}

	return validatePassword(r.Password)
}

type SigninRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *SigninRequest) ValidateRequest() error {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))

	if err := validateEmail(r.Email); err != nil {
		return err
	}

	if r.Password == "" {
		return app_errors.BadRequest("Password is required")
	}

	return nil
}

func validateEmail(email string) error {
	if len(email) > 150 || !emailRegex.MatchString(email) {
		return app_errors.BadRequest("A valid email address is required")
	}

	return nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return app_errors.BadRequest("Password must be at least 8 characters long")
	}

	if len(password) > 72 {
		return app_errors.BadRequest("Password must be at most 72 characters long")
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSymbol {
		return app_errors.BadRequest("Password must contain an uppercase letter, a lowercase letter, a number and a symbol")
	}

	return nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
