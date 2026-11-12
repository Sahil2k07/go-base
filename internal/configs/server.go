package configs

import (
	"os"
	"time"
)

type ServerConfig struct {
	Port               string        `toml:"port"`
	Origins            []string      `toml:"origins"`
	AccessTokenSecret  string        `toml:"access_token_secret"`
	AccessTokenExpiry  time.Duration `toml:"access_token_expiry"`
	RefreshTokenSecret string        `toml:"refresh_token_secret"`
	RefreshTokenExpiry time.Duration `toml:"refresh_token_expiry"`
	EncryptionSecret   string        `toml:"encryption_secret"`
	BcryptCost         int           `toml:"bcrypt_cost"`
}

func loadServerConfig(config ServerConfig) ServerConfig {
	return ServerConfig{
		Port:               envOrDefaultString(os.Getenv("SERVER_PORT"), config.Port),
		Origins:            envOrDefaultSlice(os.Getenv("SERVER_ORIGINS"), config.Origins),
		EncryptionSecret:   envOrDefaultString(os.Getenv("SERVER_ENCRYPTION_SECRET"), config.EncryptionSecret),
		AccessTokenSecret:  envOrDefaultString(os.Getenv("SERVER_ACCESS_TOKEN_SECRET"), config.AccessTokenSecret),
		AccessTokenExpiry:  time.Duration(envOrDefaultInt(os.Getenv("SERVER_ACCESS_TOKEN_EXPIRY_MINUTES"), int(config.AccessTokenExpiry.Minutes()))) * time.Minute,
		RefreshTokenSecret: envOrDefaultString(os.Getenv("SERVER_REFRESH_TOKEN_SECRET"), config.RefreshTokenSecret),
		RefreshTokenExpiry: time.Duration(envOrDefaultInt(os.Getenv("SERVER_REFRESH_TOKEN_EXPIRY_DAYS"), int(config.RefreshTokenExpiry.Hours()/24))) * 24 * time.Hour,
		BcryptCost:         envOrDefaultInt(os.Getenv("SERVER_BCRYPT_COST"), config.BcryptCost),
	}
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		Port:    globalConfig.Database.Port,
		Origins: globalConfig.Server.Origins,
	}
}
