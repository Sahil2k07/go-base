---
paths:
  - "internal/infrastructure/**/*.go"
---

# Infrastructure

Two shapes:

- **One implementation** (crypto): `internal/infrastructure/{name}.go`, a struct built by `NewX(config)`. No interface.
- **Several implementations** (email: local, ses): the interface and shared types in `internal/infrastructure/{name}.go`; each implementation in its own package `internal/infrastructure/{name}/{impl}.go`, with a compile-time check such as `var _ infrastructure.EmailSender = (*SES)(nil)`.

Multi-implementation rules:

- Choose the implementation from a config value (`email_provider = "local" | "ses"`) in a small switch in `internal/apps`, never inside `infrastructure`: implementation packages import `infrastructure`, so a factory there is an import cycle.
- An unknown provider fails at startup. Refuse `local` when `app.environment` is `production`.
- `local` only logs through `slog` and needs no network or credentials.
- Provider SDKs stay inside their implementation package; run `go mod vendor` after adding one.
- Services receive the interface through their constructor.
- Provider settings are new config: follow the configs rule.
- Unit-test consumers with a recording fake of the interface.
