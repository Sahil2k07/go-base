package features

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Sahil2k07/go-base/internal/apps"
	"github.com/Sahil2k07/go-base/tests/integration"
)

func TestSetup(t *testing.T) {
	env := integration.Setup(t)

	t.Run("migrations created the tables", func(t *testing.T) {
		var count int64

		err := env.DB.Raw(`
			SELECT count(*) FROM information_schema.tables
			WHERE table_schema = 'public'
			AND table_name IN ('users','roles','permissions','user_profiles','user_sessions','user_roles','role_permissions')
		`).Scan(&count).Error
		if err != nil {
			t.Fatal(err)
		}

		if count != 7 {
			t.Fatalf("want 7 tables, got %d", count)
		}
	})

	t.Run("runner skips applied migrations", func(t *testing.T) {
		if err := apps.RunMigrations(); err != nil {
			t.Fatalf("second run failed: %v", err)
		}
	})

	t.Run("sql files are rerun safe", func(t *testing.T) {
		dirs, err := os.ReadDir("migrations")
		if err != nil {
			t.Fatal(err)
		}

		for _, dir := range dirs {
			if !dir.IsDir() {
				continue
			}

			files, err := filepath.Glob(filepath.Join("migrations", dir.Name(), "*.sql"))
			if err != nil {
				t.Fatal(err)
			}

			for _, f := range files {
				sql, err := os.ReadFile(f)
				if err != nil {
					t.Fatal(err)
				}

				if err := env.DB.Exec(string(sql)).Error; err != nil {
					t.Errorf("%s is not rerun safe: %v", f, err)
				}
			}
		}
	})
}
