package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	app_errors "github.com/Sahil2k07/go-base/internal/errors"
	"github.com/Sahil2k07/go-base/internal/features/auth"
	applib "github.com/Sahil2k07/go-base/internal/lib"
	"github.com/Sahil2k07/go-base/tests/unit"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func serve(t *testing.T, handler echo.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	e.GET("/", handler)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	return rec
}

func failWith(err error) echo.HandlerFunc {
	return func(c *echo.Context) error { return applib.HandleError(c, err) }
}

func decodeBody(t *testing.T, rec *httptest.ResponseRecorder) unit.ErrorBody {
	t.Helper()

	var body unit.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("body is not JSON: %q", rec.Body.String())
	}

	return body
}

func onlyRecord(t *testing.T, logs *unit.LogCapture) unit.LogRecord {
	t.Helper()

	records := logs.Records()
	if len(records) != 1 {
		t.Fatalf("want 1 log record, got %d: %+v", len(records), records)
	}

	return records[0]
}

func TestHandleErrorResponses(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantMessage string
	}{
		{"bad request", app_errors.BadRequest("bad input"), http.StatusBadRequest, "bad input"},
		{"unauthorized", app_errors.Unauthorized("sign in first"), http.StatusUnauthorized, "sign in first"},
		{"forbidden", app_errors.Forbidden("not yours"), http.StatusForbidden, "not yours"},
		{"not found", app_errors.NotFound("user not found"), http.StatusNotFound, "user not found"},
		{"conflict", app_errors.Conflict("already taken"), http.StatusConflict, "already taken"},
		{"precondition failed", app_errors.PreconditionFailed("stale"), http.StatusPreconditionFailed, "stale"},
		{"unprocessable entity", app_errors.UnprocessableEntity("invalid"), http.StatusUnprocessableEntity, "invalid"},
		{"internal hides the cause", app_errors.Internal(errors.New("connection refused")), http.StatusInternalServerError, "internal server error"},
		{"wrapped app error keeps its status", fmt.Errorf("login: %w", app_errors.Unauthorized("invalid credentials")), http.StatusUnauthorized, "invalid credentials"},
		{"gorm record not found", gorm.ErrRecordNotFound, http.StatusNotFound, "resource not found"},
		{"wrapped gorm record not found", fmt.Errorf("find user: %w", gorm.ErrRecordNotFound), http.StatusNotFound, "resource not found"},
		{"gorm duplicated key", gorm.ErrDuplicatedKey, http.StatusConflict, "resource already exists"},
		{"postgres unique violation", &pgconn.PgError{Code: "23505"}, http.StatusConflict, "resource already exists"},
		{"postgres foreign key violation", &pgconn.PgError{Code: "23503"}, http.StatusBadRequest, "referenced resource does not exist"},
		{"other postgres error", &pgconn.PgError{Code: "42P01", Message: "relation does not exist"}, http.StatusInternalServerError, "internal server error"},
		{"echo http error", echo.NewHTTPError(http.StatusTeapot, "short and stout"), http.StatusTeapot, "short and stout"},
		{"unknown error", errors.New("something unexpected"), http.StatusInternalServerError, "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(t, failWith(tt.err))

			unit.AssertStatus(t, rec, tt.wantStatus)

			if got := decodeBody(t, rec).Message; got != tt.wantMessage {
				t.Fatalf("want message %q, got %q", tt.wantMessage, got)
			}
		})
	}
}

func TestHandleErrorDoesNotLeakInternalDetails(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"internal error", app_errors.Internal(errors.New("password=hunter2 rejected by db"))},
		{"unknown error", errors.New("password=hunter2 rejected by db")},
		{"postgres error", &pgconn.PgError{Code: "42P01", Message: "password=hunter2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := serve(t, failWith(tt.err))

			if strings.Contains(rec.Body.String(), "hunter2") {
				t.Fatalf("response leaked the cause: %s", rec.Body.String())
			}
		})
	}
}

func TestHandleErrorNil(t *testing.T) {
	logs := unit.CaptureLogs(t)

	rec := serve(t, failWith(nil))

	if rec.Body.Len() != 0 {
		t.Fatalf("expected no body, got %q", rec.Body.String())
	}

	if n := len(logs.Records()); n != 0 {
		t.Fatalf("expected nothing logged, got %d records", n)
	}
}

func TestHandleErrorLogging(t *testing.T) {
	t.Run("client errors log a warning with the request details", func(t *testing.T) {
		logs := unit.CaptureLogs(t)

		serve(t, failWith(app_errors.NotFound("user not found")))

		record := onlyRecord(t, logs)

		if record.Level != slog.LevelWarn || record.Message != "Request processing failed" {
			t.Fatalf("got %v %q", record.Level, record.Message)
		}

		want := map[string]any{
			"ip_address":     "192.0.2.1",
			"request_path":   "/",
			"request_method": http.MethodGet,
			"request_url":    "http://example.com/",
			"error_message":  "user not found",
		}

		for key, value := range want {
			if record.Attrs[key] != value {
				t.Errorf("%s: want %v, got %v", key, value, record.Attrs[key])
			}
		}
	})

	t.Run("server errors log at error level and keep the cause", func(t *testing.T) {
		logs := unit.CaptureLogs(t)

		serve(t, failWith(app_errors.Internal(errors.New("connection refused"))))

		record := onlyRecord(t, logs)

		if record.Level != slog.LevelError {
			t.Fatalf("want ERROR, got %v", record.Level)
		}

		if record.Attrs["error_message"] != "internal server error: connection refused" {
			t.Fatalf("error_message: %v", record.Attrs["error_message"])
		}
	})

	t.Run("request id is logged when the response has one", func(t *testing.T) {
		logs := unit.CaptureLogs(t)

		serve(t, func(c *echo.Context) error {
			c.Response().Header().Set(echo.HeaderXRequestID, "req-123")
			return applib.HandleError(c, app_errors.BadRequest("bad"))
		})

		if got := onlyRecord(t, logs).Attrs["details.request_id"]; got != "req-123" {
			t.Fatalf("details.request_id: %v", got)
		}
	})

	t.Run("no details without a request id", func(t *testing.T) {
		logs := unit.CaptureLogs(t)

		serve(t, failWith(app_errors.BadRequest("bad")))

		if _, ok := onlyRecord(t, logs).Attrs["details.request_id"]; ok {
			t.Fatal("details.request_id should be absent")
		}
	})

	t.Run("stack trace points at the project code that created the error", func(t *testing.T) {
		logs := unit.CaptureLogs(t)

		serve(t, func(c *echo.Context) error {
			return applib.HandleError(c, (&auth.SigninRequest{Email: "nope"}).ValidateRequest())
		})

		frames, ok := onlyRecord(t, logs).Attrs["stack_trace"].([]app_errors.Frame)
		if !ok || len(frames) == 0 {
			t.Fatal("expected a non-empty stack_trace")
		}

		if !strings.HasPrefix(frames[0].File, "internal/features/auth/") {
			t.Errorf("first frame is %+v, want a file in internal/features/auth", frames[0])
		}
	})
}
