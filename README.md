# go-base

A reusable Go backend project template designed to serve as a starting point for building new Go applications. It provides a structured architecture, database layer, authentication, migrations, testing infrastructure, and AI-assisted development guidelines out of the box.

This repository is a **major update to the previous version**, with a focus on maintainability, safer database operations, modular architecture, and a consistent development workflow.

## Highlights

- **GORM + GORM Gen** — Type-safe, generated database queries to reduce manual query errors and improve query safety.
- **Feature-based architecture** — Code is organized by application features instead of separating everything into global handler, service, and repository folders.
- **Testing setup** — Unit and integration test infrastructure, with PostgreSQL provisioned through Docker using Testcontainers.
- **Custom SQL migrations** — File-based SQL migrations with a custom migration runner instead of relying on GORM's automatic migration system.
- **Claude Code support** — Includes `CLAUDE.md`, rules, and skills to guide AI-assisted development and keep agents aligned with the project's architecture and conventions.

## Tech Stack

- Go
- Echo
- PostgreSQL
- GORM
- GORM Gen
- Testcontainers
- Docker
- Claude Code

## Project Structure

```text
.
├── cmd/
│   ├── server/         # Application server entry point
│   ├── migration/      # Migration runner entry point
│   └── gen/            # GORM Gen code generation
├── configs/            # Application configuration
├── internal/
│   ├── apps/           # Application setup and wiring
│   ├── database/       # Database connection
│   ├── errors/         # Application errors
│   ├── features/       # Feature-based modules
│   ├── infrastructure/ # Infrastructure services
│   ├── lib/            # Shared utilities
│   ├── middlewares/    # HTTP middleware
│   ├── models/         # Database models
│   ├── store/          # Data access
│   └── utils/          # Utilities
├── migrations/
│   ├── indexes/
│   └── tables/
├── tests/
│   ├── integration/
│   └── unit/
├── docker/
├── CLAUDE.md
└── Makefile
```

## Getting Started

### Prerequisites

Make sure you have the following installed:

- Go (version required by `go.mod`)
- PostgreSQL
- Docker (required for integration tests)
- Git

### 1. Clone the repository

```bash
git clone https://github.com/Sahil2k07/go-base.git
cd go-base
```

### 2. Configure the application

Copy the example configuration files:

```bash
cp migration.example.toml migration.toml
cp server.example.toml server.toml
cp test.example.toml test.toml
```

Update the configuration values as required for your local environment, especially database credentials, server settings, and secrets.

### 3. Download dependencies

```bash
go mod download
```

### 4. Run database migrations

Make sure your PostgreSQL database is running and matches the configuration in `migration.toml`.

Run the migration tool:

```bash
go run cmd/migration/main.go
```

### 5. Start the server

Once the migrations have completed successfully:

```bash
go run cmd/server/main.go
```

The server will start using the configuration defined in `server.toml`.

## Running Tests

This project includes unit and integration tests.

Integration tests use **Testcontainers** to provision a PostgreSQL container automatically. Make sure Docker is running before executing tests.

Run all tests without using cached results:

```bash
go test -count=1 -v ./...
```

- `-count=1` disables test caching.
- `-v` enables verbose test output.

## Database Migrations

This project uses its own SQL-based migration system instead of GORM's migration runner.

Migrations are organized by object type under the `migrations/` directory:

```text
migrations/
├── tables/
├── indexes/
├── views/
├── functions/
├── procedures/
└── triggers/
```

Migration files are written in SQL and executed by the custom migration runner.

This approach provides explicit control over database changes, allows migrations to be versioned and reviewed, and avoids relying on automatic schema changes.

## Code Generation

This project uses **GORM Gen** to generate database query code.

After modifying database models, run the code generation command to regenerate the query layer:

```bash
go run cmd/gen/main.go
```

Generated code should be kept in sync with the models and schema.

## Docker

Docker configuration is available under the `docker/` directory.

PostgreSQL can be started locally using Docker Compose:

```bash
docker compose -f docker/docker-compose.yml up -d postgres
```

Stop the container:

```bash
docker compose -f docker/docker-compose.yml down
```

To remove the persisted database volume as well:

```bash
docker compose -f docker/docker-compose.yml down -v
```

> `down -v` permanently removes the PostgreSQL data stored in the Compose volume.

## Make Commands

The project includes a Makefile for common development and build tasks.

| Command          | Description                |
| ---------------- | -------------------------- |
| `make`           | Default build target       |
| `make build`     | Build the server           |
| `make server`    | Build the server binary    |
| `make migration` | Build the migration binary |
| `make tests`     | Run tests                  |

Refer to the Makefile for platform-specific build targets and their output.

## AI-Assisted Development

This repository includes configuration for Claude Code to help maintain architectural consistency and development standards.

- **`CLAUDE.md`** — Main project instructions and architectural overview.
- **Rules** — Additional instructions scoped to specific areas of the codebase.
- **Skills** — Reusable workflows and domain-specific guidance for development tasks.

These instructions help AI agents understand the project structure, follow established conventions, and avoid introducing inconsistent patterns.

## Goals

The goal of this template is to reduce repetitive setup when starting new Go backend projects by providing a foundation that already includes:

- Application architecture and dependency wiring
- Database integration and generated queries
- Authentication infrastructure
- Custom error handling
- SQL-based migration workflow
- Unit and integration testing
- Docker-based development support
- AI coding assistant guidelines

Use this repository as a starting point, adapt it to your application's requirements, and build on top of the existing foundation.

---

**go-base** — A foundation for starting new Go backend applications.
