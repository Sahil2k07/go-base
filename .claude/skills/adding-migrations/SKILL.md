---
name: adding-migrations
description: Adds or changes the database schema: re-run-safe SQL migration files, the matching gorm model, and regenerated gen code. Use when the user asks to add or alter a table, column, index, constraint, view, function, procedure or trigger, or edits anything in internal/models.
---

# Adding migrations

```
- [ ] 1. Folder and file number
- [ ] 2. SQL (copy from sql-patterns.md)
- [ ] 3. Model in internal/models
- [ ] 4. Register in cmd/gen, regenerate
- [ ] 5. Verify
```

1. **File:** folder by object type (`tables`, `indexes`, `views`, `functions`, `procedures`, `triggers`). Number = highest `NNN` in that folder + 1, zero-padded to three digits. File names are lowercase snake_case: `NNN_create_blogs.sql`, `NNN_idx_blogs_user_id.sql`, `NNN_uq_blogs_slug.sql`. The runner sorts by file name and applies in that order (like goose), so the number is the execution order: a new table's number must be higher than every table it references. Never reuse or renumber an existing file.
2. **SQL:** patterns are in [sql-patterns.md](sql-patterns.md). Indexes go in `migrations/indexes/`, not in the table file except for `deleted_at` index. `deleted_at` index must be added only when the table might have a large number of deleted rows and if added it must be in same create_table file itself.  Every statement must be safe to run twice. Never edit an existing file: the runner skips applied files by name, so changes go in a new file. Ask the user before anything that drops data (`DROP TABLE`, `DROP COLUMN`, type changes).
3. **Model:** `internal/models/{name}.go`, types matching the SQL (mapping in sql-patterns.md). Keep gorm tags in sync with the SQL.
4. **Gen:** add a new model to `ApplyBasic` in `cmd/gen/main.go`, run `go run ./cmd/gen` (local Postgres must be up), and `go mod tidy && go mod vendor` only if imports changed. Never edit `internal/gen`.
5. **Verify:** `go build ./...`, then `go test -count=1 ./tests/integration/...`, which applies every migration to a fresh Postgres. For a local database, `go run ./cmd/migration` twice: the second run must skip everything.
