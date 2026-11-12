package configs

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/Sahil2k07/go-base/internal/utils"
)

var (
	IsDevelopment bool = false
	globalConfig  GlobalConfig
	once          sync.Once
)

type GlobalConfig struct {
	App      AppConfig      `toml:"app"`
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
}

func LoadConfigs(filepath string) GlobalConfig {
	once.Do(func() {
		path, exists := utils.ResolvePath(filepath)
		if globalConfig.App.Environment == "" && exists {
			_, err := toml.DecodeFile(path, &globalConfig)
			if err != nil {
				slog.Error("failed to decode app.toml", "error", err)
				panic("Failed to load development configurations")
			}
		}
		globalConfig = GlobalConfig{
			App:      loadAppConfig(globalConfig.App),
			Database: loadDatabaseConfig(globalConfig.Database),
			Server:   loadServerConfig(globalConfig.Server),
		}
		IsDevelopment = strings.EqualFold(
			globalConfig.App.Environment,
			"development",
		)
	})

	return globalConfig
}
