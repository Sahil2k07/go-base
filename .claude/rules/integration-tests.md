---
paths:
  - "tests/integration/**"
---

# Integration tests

- Tests: `tests/integration/features/{feature}_test.go` (`package features`). Harness: `config.go`, `harness.go`, `proxy.go`, `client.go` in `tests/integration/` (`package integration`). Config: `configs/test.toml`.
- Start every test with `env := integration.Setup(t)`: a fresh Postgres container (image from `[test] postgres_image`), the real migrations, and the real app through `apps.NewServer`. Subtests share their parent's container; call `Setup` inside `t.Run` for a fresh one.
- Requests: `status, body := env.Call(t, method, path, bearer, payload)`, then `integration.Decode[T](t, body)`. Paths are the `/api/...` routes the handler registers.
- Config loads once per process, so the database address in `test.toml` is fixed and an in-process TCP forwarder (`proxy.go`) points it at the current container. Do not try to reload config or change that mechanism.
- Get tokens by calling signup/signin. Seed roles, permissions and `user_roles` through `env.DB`.
- Per endpoint cover: success, validation failure (400), missing/invalid token (401) and missing role/permission (403) when protected, duplicate (409), and the stored result.
- No `t.Parallel()`. Needs Docker running; skipped by `-short`. Run: `go test -count=1 -v ./tests/integration/...`
- Extend the harness in `tests/integration/*.go`; do not edit `internal/` to make a test pass.
