---
paths:
  - "internal/features/**/*.go"
  - "internal/apps/**/*.go"
---

# Features

- A feature is `internal/features/{name}/` with `handler.go`, `service.go`, `views.go`, `repository.go`, plus a store implementation in `internal/store/` and wiring in `bootstrapServer` (`internal/apps/server.go`): repositories, then services, then handlers.
- `repository.go`: interface only, written in `internal/models` types.
- `service.go`: concrete struct with a `NewXService(deps...)` constructor. Add an interface only when a test or a second implementation needs one.
- `handler.go`: `NewXHandler(group *echo.Group, svc)` registers routes on the group it is given (`publicGroup`, or `privateGroup` which applies `AuthMiddleware()`). A handler only binds, calls `ValidateRequest()`, calls the service and returns `lib.HandleError(c, err)` on failure.
- Role/permission checks are attached in `handler.go`.
- `views.go`: request/response structs with camelCase json tags. `ValidateRequest()` has a pointer receiver, normalizes first (trim, lowercase email), then validates and returns `BadRequest` errors.
- Status codes: validation 400; bad credentials 401 with one generic message; missing role/permission 403; duplicate 409; unexpected failures `Internal(err)`. Translate `gorm.ErrRecordNotFound` only when the status must differ from 404.
- A feature may use another feature's repository interface (auth uses `user.UserRepository`). A feature never imports `store`, `middlewares`, gorm or gen: `store` imports features and `middlewares` imports `store`, so either would be an import cycle.
- Write tests first: an integration test per handler, unit tests for every service and validation method.
