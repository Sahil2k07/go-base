package configs

import (
	"os"
	"strconv"
	"strings"

	"github.com/Sahil2k07/go-base/internal/utils"
)

type AppConfig struct {
	Environment string `toml:"environment"`
	JSONLogs    bool   `toml:"json_logs"`
	FileLogging bool   `toml:"file_logging"`
}

func GetAppConfig() AppConfig {
	return AppConfig{
		Environment: globalConfig.App.Environment,
		JSONLogs:    globalConfig.App.JSONLogs,
		FileLogging: globalConfig.App.FileLogging,
	}
}

func loadAppConfig(config AppConfig) AppConfig {
	return AppConfig{
		Environment: envOrDefaultString(
			os.Getenv("APP_ENV"),
			config.Environment,
		),
		JSONLogs:    utils.StringToBool(os.Getenv("JSON_LOGS"), globalConfig.App.JSONLogs),
		FileLogging: utils.StringToBool(os.Getenv("FILE_LOGGING"), globalConfig.App.FileLogging),
	}
}

func envOrDefaultSlice(value string, current []string) []string {
	if strings.TrimSpace(value) == "" {
		return current
	}

	parts := strings.Split(value, ",")

	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	return parts
}

func envOrDefaultString(value, current string) string {
	if value != "" {
		return value
	}

	return current
}

func envOrDefaultInt(value string, current int) int {
	if value == "" {
		return current
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return current
	}

	return parsed
}
