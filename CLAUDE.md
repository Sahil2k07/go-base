# go-base

Go 1.26+, Echo v5 (handlers take `*echo.Context`), GORM + GORM Gen, Postgres. Module `github.com/Sahil2k07/go-base`. Project template: keep it simple, no speculative abstractions.

## Commands

- `go run ./cmd/server` / `./cmd/migration` / `./cmd/gen` (regenerate `internal/gen` after ANY `internal/models` change)
- `go mod tidy && go mod vendor` after dependency changes (`vendor/` is used)
- `go test -short ./...` unit only; `go test -count=1 ./...` unit + integration (needs Docker)
- Run `gofmt` on everything you touch

## Workflow (TDD)

Write the tests first, show them to the user, and wait for approval before implementing. Do not change `internal/` code only to ease testing without asking.

## Layout

- `cmd/{app}/main.go`: load config, call `internal/apps/{app}.go`, handle the error. Nothing else. All bootstrap/wiring lives in `internal/apps/{app}.go`.
- `internal/features/{name}/`: `handler.go` (routes/HTTP), `service.go` (concrete struct, no interface), `views.go` (request/response types, `ValidateRequest`), `repository.go` (interface only).
- `internal/store/`: repository implementations. Query with `internal/gen`, return `internal/models` types. Read `database.GQ` inside methods, never at construction (it is nil until `ConnectPostgres()` runs).
- Features never import `store`, `middlewares`, gorm or gen.
- `internal/models`: DB models. FK lives on the child side. Secrets use `json:"-"`.
- `internal/infrastructure/`: single-implementation service = struct in `{name}.go` (e.g. `crypto.go`). Multi-implementation service = interface in `{name}.go`, implementations in `{name}/{impl}.go` (e.g. `email.go`, `email/local.go`, `email/ses.go`), chosen by a config value.
- Dependencies are injected through constructors; repositories are never package-level vars.

## Config

- TOML files in `configs/` (`server`, `migration`, `test`), env vars override (`SERVER_*`, `DB_*`, `APP_*`). Loaded once (`sync.Once`), never reload.
- New settings: add to the struct, the loader and every toml in `configs/`, same pattern as the existing ones.

## Migrations

- Custom runner (`internal/migrations`). Files: `migrations/{tables,indexes,views,functions,procedures,triggers}/NNN_name.sql`.
- Must be re-run safe (`IF NOT EXISTS`, `CREATE OR REPLACE`). Applied files are skipped, so schema changes go in a new file.
- Soft-deleted tables (`gorm.Model`): unique indexes must be partial (`WHERE deleted_at IS NULL`).

## Errors and logging

- Return `internal/errors` (alias `app_errors`): `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `Conflict`, `UnprocessableEntity`, `Internal(err)`; `.Wrap(err)` keeps the cause.
- Handlers end with `return lib.HandleError(c, err)`. Response is `{success, message, code}`. Never put causes or secrets in messages. Stores may return raw gorm/pg errors; `HandleError` maps them (404/409/400/500).
- Log only with `slog`; never `log` or `fmt`. Never log passwords, tokens or secrets.

## Auth

- Protected routes: `AuthMiddleware()`, then `RequireRoles(...)` (any) / `RequirePermissions(...)` (all). Claims: `infrastructure.GetAuthClaims(c)`.
- Refresh tokens are stored only as SHA-256 hashes and rotated via `RotateUserSession`. Reusing an old token deletes the session.

## Tests

- No `_test.go` files inside `internal/`. Table-driven. Always `-count=1`.
- Integration: `tests/integration/features/{feature}_test.go`, harness/config in `tests/integration/*.go`. Each test calls `integration.Setup(t)` (fresh Postgres container, config from `configs/test.toml`). No `t.Parallel()` (global app state).
- Unit: `tests/unit/{same path as internal}/{file}_test.go` (e.g. `tests/unit/middlewares/auth_test.go`), shared helpers in `tests/unit/*.go` (`MockDB`, `Serve`, tokens, `CaptureLogs`). Every method gets unit tests. External packages: test exported API only.
