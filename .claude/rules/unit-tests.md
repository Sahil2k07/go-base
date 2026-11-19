---
paths:
  - "tests/unit/**"
---

# Unit tests

- Layout: `tests/unit/{path mirroring internal}/{file}_test.go`, package named after the area (`package middlewares`), app packages imported with an alias (`app_mw`, `app_errors`). Shared helpers are non-test files directly in `tests/unit/` (`package unit`). Tests are external, so only exported API is testable.
- Table-driven with `t.Run`. Cover success, every failure branch and the dependency-error branch.
- Helpers in `package unit`:
  - `LoadConfig(t)` loads `configs/test.toml` once. Call it before building anything that reads global config (`middlewares.AuthMiddleware()`).
  - `MockDB(t)` swaps `database.DB` and `database.GQ` for a sqlmock-backed gorm DB and verifies expectations on cleanup. Call it at the start of each subtest. `ExpectCount(m, regex, n, args...)`, `ExpectQueryError(m, regex, args...)`.
  - `Serve(t, authHeader, middlewares...)` sends GET `/` to a handler that answers 200. `AssertStatus` also checks the `{success, message, code}` body on errors.
  - `AccessToken`, `ExpiredAccessToken`, `RefreshToken` mint real tokens from the test config.
  - `CaptureLogs(t)` returns slog records; groups are flattened (`details.request_id`).
- SQL expectations are regexes against gen's SQL (`FROM "roles"`) with args in query order. On a mismatch sqlmock prints the SQL it received; copy from that.
- Stack traces only contain frames from `internal/` and `cmd/`, so create errors through the code under test, not in the test.
- No `t.Parallel()`: tests change globals (`database.*`, slog default, config).
- A helper needed in more than one area goes in `tests/unit/*.go`, not copied.
- Run: `go test -count=1 ./tests/unit/...`
