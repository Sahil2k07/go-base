# go-base

Go 1.26+, Echo v5 (handlers take `*echo.Context`), GORM + GORM Gen, Postgres. Module `github.com/Sahil2k07/go-base`. Project template: keep it simple, no speculative abstractions.

## Commands

- `go run ./cmd/server` / `./cmd/migration` / `./cmd/gen` (regenerate `internal/gen` after ANY `internal/models` change)
- `go mod tidy && go mod vendor` after dependency changes (`vendor/` is used)
- `go test -short ./...` unit only; `go test -count=1 ./...` unit + integration (needs Docker)
- Run `gofmt` on everything you touch

## Workflow (TDD)

Write the tests first, show them to the user, and wait for approval before implementing. Do not change `internal/` code only to ease testing without asking. Skills: `adding-features`, `adding-migrations`, `adding-infrastructure-services`, `writing-tests`.

## Layout

- `cmd/{app}/main.go`: load config, call `internal/apps/{app}.go`, handle the error. Nothing else. All bootstrap/wiring lives in `internal/apps/{app}.go`.
- `internal/features/{name}/`: `handler.go`, `service.go` (concrete struct), `views.go` (request/response types), `repository.go` (the only interface). Implementations live in `internal/store/`.
- `internal/store/`: query with `internal/gen`, return `internal/models` types.
- `internal/models`: DB models. `internal/gen` is generated, never edit it.
- `internal/infrastructure/`: one-implementation service = struct in `{name}.go` (e.g. `crypto.go`). Multi-implementation service = interface in `{name}.go`, implementations in `{name}/{impl}.go` (e.g. `email/local.go`, `email/ses.go`).
- Dependencies are injected through constructors. Repositories are never package-level vars.

## Wiring (`internal/apps/server.go`)

- `NewServer(appConfig)`: Echo with `RequestLogger`, `Recover`, CORS; then `publicGroup` (no auth) and `privateGroup` (`AuthMiddleware()`), both at `/api`; then `bootstrapServer`.
- `bootstrapServer` builds in this order: repositories (`store.NewXRepository()`), services (infrastructure first, e.g. `infrastructure.NewCryptoService(appConfig.Server)`, then `NewXService(repos, infra)`), handlers (`x.NewXHandler(group, service)`).
- A new piece is one line in its block. Reuse existing instances; never build the same repository or service twice.
- Features cannot import `middlewares` (import cycle via `store`). `RequireRoles`/`RequirePermissions` are attached in `internal/apps` through sub-groups or constructor arguments.

## Config

TOML files in `configs/`, env vars override. Loaded once, never reload. New settings follow the existing pattern in `internal/configs`.

## Migrations

Custom runner. Files go in `migrations/{tables,indexes,views,functions,procedures,triggers}/NNN_name.sql` and must be re-run safe. Applied files are skipped, so schema changes go in a new file.
File names are lowercase snake*case starting with a zero-padded number (`001*`, `002\_`, continuing). The runner sorts by file name and applies in that order, like goose, so the number is the execution order.

## Errors and logging

- Return `internal/errors` (alias `app_errors`): `BadRequest`, `Unauthorized`, `Forbidden`, `NotFound`, `Conflict`, `UnprocessableEntity`, `Internal(err)`; `.Wrap(err)` keeps the cause.
- Handlers end with `return lib.HandleError(c, err)`. Response is `{success, message, code}`. Never put causes or secrets in messages.
- Log only with `slog`, never `log` or `fmt`. Never log passwords, tokens or secrets.
- Every config file has `[app]` with `environment`, `json_logs` (JSON slog output) and `file_logging`. Each `cmd/{app}/main.go` `init()` loads its toml, calls `infrastructure.ConfigureLogging("{app}")`, then `database.ConnectPostgres()`.
- With `file_logging = true`, logs are also written to `logs/{app}` (e.g. `logs/server`, `logs/migration`). Read them when debugging a failure. `configs/test.toml` sets it to `false`.

## Auth

Private routes run `AuthMiddleware()`; roles and permissions are checked per request by `RequireRoles` (any) / `RequirePermissions` (all). Claims: `infrastructure.GetAuthClaims(c)`.

## Tests

- No `_test.go` files inside `internal/`. Table-driven. Always `-count=1`. No `t.Parallel()`.
- Integration: `tests/integration/features/{feature}_test.go`, harness in `tests/integration/*.go`.
- Unit: `tests/unit/{same path as internal}/{file}_test.go`, helpers in `tests/unit/*.go`. Every method gets unit tests.
