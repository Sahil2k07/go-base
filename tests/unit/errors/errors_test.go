package errors

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	app_errors "github.com/Sahil2k07/go-base/internal/errors"
	"github.com/Sahil2k07/go-base/internal/features/auth"
)

func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *app_errors.AppError
		code int
	}{
		{"BadRequest", app_errors.BadRequest("m"), http.StatusBadRequest},
		{"Unauthorized", app_errors.Unauthorized("m"), http.StatusUnauthorized},
		{"Forbidden", app_errors.Forbidden("m"), http.StatusForbidden},
		{"NotFound", app_errors.NotFound("m"), http.StatusNotFound},
		{"Conflict", app_errors.Conflict("m"), http.StatusConflict},
		{"PreconditionFailed", app_errors.PreconditionFailed("m"), http.StatusPreconditionFailed},
		{"UnprocessableEntity", app_errors.UnprocessableEntity("m"), http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code || tt.err.Message != "m" || tt.err.Err != nil {
				t.Fatalf("unexpected error: %+v", tt.err)
			}
		})
	}
}

func TestErrorString(t *testing.T) {
	if got := app_errors.NotFound("user not found").Error(); got != "user not found" {
		t.Fatalf("got %q", got)
	}

	if got := app_errors.Internal(errors.New("db down")).Error(); got != "internal server error: db down" {
		t.Fatalf("got %q", got)
	}
}

func TestInternalKeepsTheCauseOutOfTheMessage(t *testing.T) {
	cause := errors.New("connection refused")
	err := app_errors.Internal(cause)

	if err.Code != http.StatusInternalServerError || err.Message != "internal server error" {
		t.Fatalf("unexpected error: %+v", err)
	}

	if !errors.Is(err, cause) {
		t.Fatal("the cause should be reachable through errors.Is")
	}
}

func TestWrapReturnsACopy(t *testing.T) {
	sentinel := app_errors.Unauthorized("invalid credentials")
	cause := errors.New("no rows")

	wrapped := sentinel.Wrap(cause)

	if sentinel.Err != nil {
		t.Fatal("Wrap modified the original error")
	}

	if wrapped == sentinel {
		t.Fatal("Wrap should return a copy")
	}

	if wrapped.Code != http.StatusUnauthorized || wrapped.Message != "invalid credentials" {
		t.Fatalf("copy lost its code or message: %+v", wrapped)
	}

	if !errors.Is(wrapped, cause) {
		t.Fatal("the cause should be reachable through errors.Is")
	}
}

func TestAsFindsAWrappedAppError(t *testing.T) {
	err := fmt.Errorf("login: %w", app_errors.Forbidden("nope"))

	var appErr *app_errors.AppError
	if !errors.As(err, &appErr) || appErr.Code != http.StatusForbidden {
		t.Fatalf("errors.As did not find the AppError in %v", err)
	}
}

func TestNew(t *testing.T) {
	cause := errors.New("x")
	err := app_errors.New(http.StatusTeapot, "short and stout", cause)

	if err.Code != http.StatusTeapot || err.Message != "short and stout" || !errors.Is(err, cause) {
		t.Fatalf("unexpected error: %+v", err)
	}
}

func TestStackTraceOnlyHasProjectFrames(t *testing.T) {
	// created inside internal/features/auth, so its frames belong to the project
	err := (&auth.SigninRequest{Email: "nope"}).ValidateRequest()

	var appErr *app_errors.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("want an AppError, got %v", err)
	}

	frames := appErr.StackTrace()
	if len(frames) == 0 {
		t.Fatal("expected a stack trace")
	}

	foundValidate := false

	for _, f := range frames {
		if !strings.HasPrefix(f.File, "internal/") && !strings.HasPrefix(f.File, "cmd/") {
			t.Errorf("frame outside the project: %+v", f)
		}

		if f.Line <= 0 {
			t.Errorf("frame without a line number: %+v", f)
		}

		if strings.Contains(f.Function, "ValidateRequest") {
			foundValidate = true
		}
	}

	if !foundValidate {
		t.Errorf("no ValidateRequest frame in %+v", frames)
	}
}
