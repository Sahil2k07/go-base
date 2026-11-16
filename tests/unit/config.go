package unit

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Sahil2k07/go-base/internal/configs"
)

const TestConfigFile = "configs/test.toml"

// Must match the os.Getenv names in internal/configs, which would otherwise override test.toml.
var envPrefixesToClear = []string{"SERVER_", "DB_", "APP_"}

var (
	configOnce sync.Once
	configErr  error
	config     configs.GlobalConfig
)

// LoadConfig loads configs/test.toml once per test process, the way cmd/server loads its
// config, and returns it. Call it before building anything that reads the global config,
// such as middlewares.AuthMiddleware().
func LoadConfig(t *testing.T) configs.GlobalConfig {
	t.Helper()

	configOnce.Do(func() {
		clearConfigEnv()

		// the config loader resolves paths from the repo root
		if configErr = os.Chdir(projectRoot()); configErr != nil {
			return
		}

		config = configs.LoadConfigs(TestConfigFile)

		// LoadConfigs only logs when it cannot read the file, so verify that it really loaded
		if config.Server.AccessTokenSecret == "" {
			configErr = errors.New(TestConfigFile + " was not loaded")
		}
	})

	if configErr != nil {
		t.Fatalf("test config: %v", configErr)
	}

	return config
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
