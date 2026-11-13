package errors

import (
	"net/http"
	"runtime"
	"strings"
)

type AppError struct {
	Code    int
	Message string
	Err     error
	stack   []uintptr
}

type Frame struct {
	File     string `json:"file"`
	Function string `json:"function"`
	Line     int    `json:"line"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// Wrap attaches a cause and re-captures the stack, so package-level sentinel
// errors report where they were used instead of where they were declared.
func (e *AppError) Wrap(err error) *AppError {
	cp := *e
	cp.Err = err
	cp.stack = callers(3)
	return &cp
}

// StackTrace returns this project's frames only (internal/ and cmd/).
func (e *AppError) StackTrace() []Frame {
	frames := []Frame{}
	it := runtime.CallersFrames(e.stack)

	for {
		f, more := it.Next()

		if file, ok := projectFile(f.File); ok {
			frames = append(frames, Frame{
				File:     file,
				Function: f.Function[strings.LastIndex(f.Function, ".")+1:],
				Line:     f.Line,
			})
		}

		if !more {
			break
		}
	}

	return frames
}

func BadRequest(msg string) *AppError { return newErr(http.StatusBadRequest, msg, nil) }

func Unauthorized(msg string) *AppError { return newErr(http.StatusUnauthorized, msg, nil) }

func Forbidden(msg string) *AppError { return newErr(http.StatusForbidden, msg, nil) }

func NotFound(msg string) *AppError { return newErr(http.StatusNotFound, msg, nil) }

func Conflict(msg string) *AppError { return newErr(http.StatusConflict, msg, nil) }

func PreconditionFailed(msg string) *AppError {
	return newErr(http.StatusPreconditionFailed, msg, nil)
}

func UnprocessableEntity(msg string) *AppError {
	return newErr(http.StatusUnprocessableEntity, msg, nil)
}

// Internal hides the cause from the client and keeps it for the logs.
func Internal(err error) *AppError {
	return newErr(http.StatusInternalServerError, "internal server error", err)
}

// Exported so lib.HandleError can build one for errors it did not receive as an AppError.
func New(code int, msg string, err error) *AppError { return newErr(code, msg, err) }

func newErr(code int, msg string, err error) *AppError {
	// skip 4: runtime.Callers, callers, newErr, the constructor, so the trace starts at the caller
	return &AppError{Code: code, Message: msg, Err: err, stack: callers(4)}
}

func callers(skip int) []uintptr {
	pcs := make([]uintptr, 16)
	n := runtime.Callers(skip, pcs)
	return pcs[:n]
}

func projectFile(file string) (string, bool) {
	file = strings.ReplaceAll(file, "\\", "/")

	if strings.Contains(file, "/vendor/") {
		return "", false
	}

	for _, root := range []string{"/internal/", "/cmd/"} {
		if i := strings.Index(file, root); i >= 0 {
			return file[i+1:], true
		}
	}

	return "", false
}
