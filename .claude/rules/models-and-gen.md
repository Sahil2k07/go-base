---
paths:
  - "internal/models/**/*.go"
  - "internal/gen/**"
  - "cmd/gen/**"
---

# Models and gen

- `internal/gen` is generated: never edit it. After a model change: add new models to `ApplyBasic` in `internal/apps/gen.go`, run `go run ./cmd/gen` (it connects to the database in `configs/migration.toml`, so local Postgres must be up), then `go mod tidy && go mod vendor` if imports changed.
- Ship the matching migration in the same change. Gorm tags describe the schema; the SQL enforces it.
- The foreign key lives on the child. Parent side:

```go
  Profile  UserProfile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile"`
  Sessions []UserSession `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
```

- Belongs-to backlinks are pointers (`User *User`). A by-value backlink is an "invalid recursive type" compile error.
- Many-to-many: `many2many:{join_table}` on both sides plus an explicit join model with a composite primary key (`primaryKey` on both ids) and no `gorm.Model`.
- UUID columns: `uuid.UUID` with `type:uuid`. Generate the id in Go before insert when it is needed earlier (the session id goes into the refresh token).
- `json:"-"` on password and token hashes and on backlinks that must not be serialized.
- Models stay database-shaped. API shapes belong in the feature's `views.go`.
