---
name: adding-infrastructure-services
description: Adds a cross-cutting service under internal/infrastructure, either a single implementation (like crypto) or swappable implementations chosen by config (like email local/SES). Use when the user asks for email, notifications, SMS, file storage, payments, caching or any external provider integration.
---

# Adding an infrastructure service

```
- [ ] 1. Pick the shape
- [ ] 2. Tests first (writing-tests skill), wait for approval
- [ ] 3. Implement
- [ ] 4. Config
- [ ] 5. Wire in bootstrapServer
- [ ] 6. go mod vendor if an SDK was added; tests green
```

## 1. Shape

- **One implementation** (nothing to swap): struct in `internal/infrastructure/{name}.go`, constructor `NewX(config)` taking the config value, like `NewCryptoService(appConfig.Server)`. No interface.
- **Several implementations** (local, ses, smtp, ...): interface and shared types in `internal/infrastructure/{name}.go`; each implementation in its own package, `internal/infrastructure/{name}/{impl}.go`, with `var _ infrastructure.X = (*Impl)(nil)`. Full example: [email-example.md](email-example.md).

## 2. Choosing an implementation

- Selection lives in `internal/apps` as an exported function (`NewEmailSender(cfg, env)`), never in `infrastructure`: implementation packages import it, so a factory there is an import cycle. Export it so unit tests can reach it.
- Driven by a config value (`email_provider = "local"`). Unknown value returns an error at startup. `local` is refused when `app.environment` is `production`.
- `local` only logs through `slog`: no network, no credentials.

## 3. Config

Add the provider name and each provider's settings as new config (see the configs rule): struct field, env override, getter, every toml including `configs/test.toml` (set it to `local`).

## 4. Wiring

In `bootstrapServer`, build infrastructure services before the feature services that use them, and pass the interface to constructors. Services depend on the interface, never on a provider package.

## 5. Tests

Unit-test each implementation and the selection function. Consumers use a recording fake of the interface. Never call real providers in tests.
