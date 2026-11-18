package apps

import (
	"log/slog"
	"time"

	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/migrations"
	"gorm.io/gorm"
)

func RunMigrations() error {
	start := time.Now()

	slog.Info("Reading Files")

	histories, err := migrations.GetMigrationHistory()
	if err != nil {
		return err
	}

	tableFiles, err := migrations.GetMigrationFiles("migrations/tables", true, histories)
	if err != nil {
		return err
	}

	indexFiles, err := migrations.GetMigrationFiles("migrations/indexes", true, histories)
	if err != nil {
		return err
	}

	viewFiles, err := migrations.GetMigrationFiles("migrations/views", false, histories)
	if err != nil {
		return err
	}

	functionFiles, err := migrations.GetMigrationFiles("migrations/functions", false, histories)
	if err != nil {
		return err
	}

	procedureFiles, err := migrations.GetMigrationFiles("migrations/procedures", false, histories)
	if err != nil {
		return err
	}

	triggerFiles, err := migrations.GetMigrationFiles("migrations/triggers", false, histories)
	if err != nil {
		return err
	}

	slog.Info("Starting Applying Migrations")

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := migrations.ApplyMigration(tx, "TABLE", tableFiles); err != nil {
			return err
		}

		if err := migrations.ApplyMigration(tx, "INDEX", indexFiles); err != nil {
			return err
		}

		if err := migrations.ApplyMigration(tx, "VIEW", viewFiles); err != nil {
			return err
		}

		if err := migrations.ApplyMigration(tx, "FUNCTION", functionFiles); err != nil {
			return err
		}

		if err := migrations.ApplyMigration(tx, "PROCEDURE", procedureFiles); err != nil {
			return err
		}

		if err := migrations.ApplyMigration(tx, "TRIGGER", triggerFiles); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	migrations.PrintMigrationSummary(
		start,
		len(tableFiles),
		len(indexFiles),
		len(viewFiles),
		len(functionFiles),
		len(procedureFiles),
		len(triggerFiles),
	)

	return nil
}
