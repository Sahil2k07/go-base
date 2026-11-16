package integration

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/Sahil2k07/go-base/internal/configs"
)

const testConfigFile = "configs/test.toml"

// Must match the os.Getenv names in internal/configs, which would otherwise override test.toml.
var envPrefixesToClear = []string{"SERVER_", "DB_", "APP_"}

type testMeta struct {
	Test struct {
		PostgresImage string `toml:"postgres_image"`
	} `toml:"test"`
	Database struct {
		Port     string `toml:"db_port"`
		User     string `toml:"db_user"`
		Password string `toml:"db_password"`
		Name     string `toml:"db_name"`
	} `toml:"database"`
}

var (
	initOnce  sync.Once
	initErr   error
	meta      testMeta
	appConfig configs.GlobalConfig
	dbProxy   *proxy
)

// initConfig mirrors cmd/server: it loads the config once per process and keeps the result
// to pass to apps.NewServer. It also starts the forwarder on the database address that
// config points at.
func initConfig() error {
	initOnce.Do(func() {
		clearConfigEnv()

		// the config loader and the migration runner both resolve paths from the repo root
		if initErr = os.Chdir(projectRoot()); initErr != nil {
			return
		}

		if _, initErr = toml.DecodeFile(testConfigFile, &meta); initErr != nil {
			return
		}

		appConfig = configs.LoadConfigs(testConfigFile)

		// LoadConfigs only logs when it cannot read the file, so verify that it really loaded
		if appConfig.Server.AccessTokenSecret == "" {
			initErr = errors.New(testConfigFile + " was not loaded")
			return
		}

		dbProxy, initErr = startProxy(net.JoinHostPort("127.0.0.1", meta.Database.Port))
	})

	return initErr
}

func clearConfigEnv() {
	for _, kv := range os.Environ() {
		name, _, _ := strings.Cut(kv, "=")

		for _, prefix := range envPrefixesToClear {
			if strings.HasPrefix(name, prefix) {
				_ = os.Unsetenv(name)
			}
		}
	}
}

func projectRoot() string {
	dir, _ := os.Getwd()

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			panic("go.mod not found")
		}

		dir = parent
	}
}
