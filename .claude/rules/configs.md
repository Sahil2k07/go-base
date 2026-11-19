---
paths:
  - "configs/**"
  - "internal/configs/**"
---

# Configs

- Files: `configs/server.toml`, `configs/migration.toml`, `configs/test.toml`. `internal/configs` holds `GlobalConfig{App, Database, Server}`. Config loads once (`sync.Once`) and is read through the `GetXConfig()` getters.
- Adding a setting means all four:
  1. a field with a `toml:"snake_case"` tag on the section struct;
  2. an env override in that section's `loadXConfig`: `envOrDefaultString/Int/Slice(os.Getenv("SECTION_NAME"), file value)` (`SERVER_*`, `DB_*`, `APP_*`);
  3. the matching line in its `GetXConfig()` getter (check it reads its own section);
  4. the key in every toml that loads that section, including `configs/test.toml`.
- New section: add it to `GlobalConfig`, merge it where `LoadConfigs` merges the others, add a getter.
- Durations are strings in toml (`"15m"`, `"72h"`). Env overrides use plain numbers with a unit in the name (`_MINUTES`, `_DAYS`).
- Every toml starts with the same `[app]` section: `environment` (`development`, `test`, `production`), `json_logs`, `file_logging`. A new app's toml and `cmd/{app}/main.go` `init()` follow `cmd/server`: `LoadConfigs("configs/{app}.toml")`, `infrastructure.ConfigureLogging("{app}")`, then `database.ConnectPostgres()`. `file_logging = true` writes logs to `logs/{app}`.
- Committed tomls hold local, fake secrets only. Real secrets come from env vars.
- `LoadConfigs` logs but does not fail when a file is unreadable, so verify a loaded value when a config problem is suspected.
- `[test]` in `test.toml` is read only by the integration harness.
