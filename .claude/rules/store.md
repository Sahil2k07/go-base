---
paths:
  - "internal/store/**/*.go"
---

# Store

- One file per repository, implementing the interface in `internal/features/{name}/repository.go`.
- Read the generated queries at call time, never in the constructor (`database.GQ` is nil until `ConnectPostgres()` runs):

```go
type userStore struct {
	q *gen.Query
}

func NewUserRepository() user.UserRepository {
	return &userStore{q: database.GQ}
}
```

- Query with gen and return `internal/models` types. No raw SQL unless gen cannot express it. No business logic, no HTTP types.
- Return errors untouched (`gorm.ErrRecordNotFound`, pg errors); `lib.HandleError` maps them. Never return `nil, nil` for "not found".
- Existence check: `Count() > 0`. Name-only lists: one query with `Pluck` (add `Distinct` across joins), never load structs and loop. Joins: `Join(ur, ur.RoleID.EqCol(r.ID))`.
- Multi-step writes: `s.q.Transaction(func(tx *gen.Query) error {...})`. Create a parent with its has-one/has-many children in a single `Create`.
- Compare-and-swap: one `UpdateSimple` with the old value in `Where`, then check `RowsAffected` (see `RotateUserSession`).
- Soft delete is the default for `gorm.Model` tables. Rows that must really disappear (sessions) use `Unscoped().Delete()`.
- Password and token hashes arrive already hashed. Never hash in the store.
