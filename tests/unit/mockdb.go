package unit

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/gen"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MockDB replaces database.DB and database.GQ with a gorm connection backed by sqlmock
// for the duration of the test, mirroring what ConnectPostgres sets up. Call it before
// building anything that captures them, such as store.NewUserRepository() or
// middlewares.RequireRoles(...). Unmet expectations fail the test.
func MockDB(t *testing.T) sqlmock.Sqlmock {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}

	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		&gorm.Config{DisableAutomaticPing: true, Logger: logger.Discard},
	)
	if err != nil {
		t.Fatalf("open gorm over sqlmock: %v", err)
	}

	previousDB, previousGQ := database.DB, database.GQ
	database.DB = db
	database.GQ = gen.Use(db)

	t.Cleanup(func() {
		database.DB, database.GQ = previousDB, previousGQ

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet database expectations: %v", err)
		}

		_ = sqlDB.Close()
	})

	return mock
}

// ExpectCount expects a SELECT count(*) matching the regex and returns the given count.
func ExpectCount(m sqlmock.Sqlmock, queryRegex string, count int, args ...driver.Value) {
	m.ExpectQuery(queryRegex).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))
}

// ExpectQueryError expects a query matching the regex and makes it fail.
func ExpectQueryError(m sqlmock.Sqlmock, queryRegex string, args ...driver.Value) {
	m.ExpectQuery(queryRegex).
		WithArgs(args...).
		WillReturnError(errors.New("connection refused"))
}
