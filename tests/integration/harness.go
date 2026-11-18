package integration

import (
	"context"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/Sahil2k07/go-base/internal/apps"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/gorm"
)

// Env is a running app backed by its own fresh Postgres container.
type Env struct {
	URL string
	DB  *gorm.DB
}

// Setup starts a new Postgres container, applies the migrations and boots the app
// against it. Everything is torn down when the test finishes.
func Setup(t *testing.T) *Env {
	t.Helper()

	if testing.Short() {
		t.Skip("integration tests skipped in -short mode")
	}

	if err := initConfig(); err != nil {
		t.Fatalf("test config: %v", err)
	}

	ctx := context.Background()

	ctr, err := postgres.Run(ctx, meta.Test.PostgresImage,
		postgres.WithDatabase(meta.Database.Name),
		postgres.WithUsername(meta.Database.User),
		postgres.WithPassword(meta.Database.Password),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		t.Fatalf("start postgres: %v", err)
	}

	host, err := ctr.Host(ctx)
	if err != nil {
		t.Fatalf("container host: %v", err)
	}

	mapped, err := ctr.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("container port: %v", err)
	}

	dbProxy.SetTarget(net.JoinHostPort(host, mapped.Port()))

	database.ConnectPostgres()
	db := database.DB
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		dbProxy.DropConnections()
	})

	if err := apps.RunMigrations(); err != nil {
		t.Fatalf("migrations: %v", err)
	}

	srv := httptest.NewServer(apps.NewServer(appConfig))
	t.Cleanup(srv.Close)

	return &Env{URL: srv.URL, DB: db}
}
