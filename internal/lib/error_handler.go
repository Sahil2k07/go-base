package lib

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	app_errors "github.com/Sahil2k07/go-base/internal/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func HandleError(c *echo.Context, err error) error {
	if err == nil {
		return nil
	}

	appErr := toAppError(err)

	req := c.Request()

	attrs := []any{
		slog.String("ip_address", c.RealIP()),
		slog.String("request_path", req.URL.Path),
		slog.String("request_url", c.Scheme()+"://"+req.Host+req.RequestURI),
		slog.String("request_method", req.Method),
		slog.String("error_message", err.Error()),
		slog.Any("stack_trace", appErr.StackTrace()),
	}

	if id := c.Response().Header().Get(echo.HeaderXRequestID); id != "" {
		attrs = append(attrs, slog.Group("details", slog.String("request_id", id)))
	}

	if appErr.Code >= http.StatusInternalServerError {
		slog.ErrorContext(req.Context(), "Request processing failed", attrs...)
	} else {
		slog.WarnContext(req.Context(), "Request processing failed", attrs...)
	}

	return c.JSON(appErr.Code, errorResponse{
		Success: false,
		Message: appErr.Message,
		Code:    appErr.Code,
	})
}

// toAppError maps any error to an AppError. Our own errors win, then known
// library errors, and anything unrecognised becomes a 500.
func toAppError(err error) *app_errors.AppError {
	var appErr *app_errors.AppError
	var httpErr *echo.HTTPError
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &appErr):
		return appErr

	case errors.As(err, &httpErr):
		// raised by echo itself: unknown route, bad bind, etc.
		return app_errors.New(httpErr.Code, fmt.Sprint(httpErr.Message), err)

	case errors.Is(err, gorm.ErrRecordNotFound):
		return app_errors.NotFound("resource not found").Wrap(err)

	case errors.Is(err, gorm.ErrDuplicatedKey):
		return app_errors.Conflict("resource already exists").Wrap(err)

	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case "23505": // unique_violation
			return app_errors.Conflict("resource already exists").Wrap(err)
		case "23503": // foreign_key_violation
			return app_errors.BadRequest("referenced resource does not exist").Wrap(err)
		}
	}

	return app_errors.Internal(err)
}
