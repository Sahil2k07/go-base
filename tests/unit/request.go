package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
)

type ErrorBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Serve sends one GET request through the middlewares to a handler that answers 200.
// authHeader is the full Authorization header value, or "" for none.
func Serve(t *testing.T, authHeader string, middlewares ...echo.MiddlewareFunc) *httptest.ResponseRecorder {
	t.Helper()

	e := echo.New()
	e.GET("/", func(c *echo.Context) error { return c.String(http.StatusOK, "ok") }, middlewares...)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}

// AssertStatus checks the status code and, for errors, that the body has the standard JSON shape.
func AssertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("want status %d, got %d: %s", want, rec.Code, rec.Body.String())
	}

	if want < http.StatusBadRequest {
		return
	}

	var body ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("error body is not JSON: %q", rec.Body.String())
	}

	if body.Success || body.Code != want || body.Message == "" {
		t.Fatalf("unexpected error body: %s", rec.Body.String())
	}
}
