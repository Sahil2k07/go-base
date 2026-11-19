---
name: adding-features
description: Adds a new feature or API endpoint to this Go service end to end: migration, model, repository, store, service, handler, wiring in bootstrapServer, and tests. Use when the user asks to add or extend a feature, resource, CRUD, route or endpoint, for example "add blogs" or "add a reviews endpoint".
---

# Adding a feature

Mirror `internal/features/user/` (smallest complete feature), its file in `internal/store/`, and its tests. Read them first.

```
- [ ] 1. Scope: routes, request/response fields, who may call each route (public, any user, role, permission), new tables
- [ ] 2. Tests first (writing-tests skill). Show them, then STOP until the user approves
- [ ] 3. Schema: migration, model, gen (adding-migrations skill)
- [ ] 4. repository.go: interface only
- [ ] 5. internal/store/{name}.go: implementation
- [ ] 6. views.go, then service.go, then handler.go
- [ ] 7. Wire it in bootstrapServer
- [ ] 8. Tests green, gofmt, go build ./...
```

## Names (feature `blog`)

| Piece      | Name                                                       |
| ---------- | ---------------------------------------------------------- |
| model      | `Blog` in `internal/models/blog.go`                        |
| repository | `BlogRepository` in `internal/features/blog/repository.go` |
| store      | `NewBlogRepository()` in `internal/store/blog.go`          |
| service    | `NewBlogService(repo)`                                     |
| handler    | `NewBlogHandler(group, service)`                           |

## Wiring

`internal/apps/server.go` builds everything in `bootstrapServer`. Add one line to each block, in this order:

```go
// Repositories
blogRepo := store.NewBlogRepository()
// Services
blogService := blog.NewBlogService(blogRepo)
// Handlers
blog.NewBlogHandler(privateGroup, blogService)
```

- `publicGroup` has no auth, `privateGroup` applies `AuthMiddleware()`. Both are mounted at `/api`.
- Reuse instances that already exist (a repository another service needs, the crypto service). Never build the same one twice.
- Features cannot import `middlewares` (import cycle through `store`). Attach role/permission checks in `internal/apps` with a sub-group, `admin := privateGroup.Group("/admin", middlewares.RequireRoles("admin"))`, and pass that group to the handler. For per-route differences, pass the middleware into the handler constructor as a parameter.

## Constraints

- Handlers only bind, call `ValidateRequest()`, call the service, and return `lib.HandleError(c, err)`.
- A feature that needs another feature's data depends on that feature's repository interface.
- Do not change an existing interface or response shape without asking: other code and tests depend on them.
