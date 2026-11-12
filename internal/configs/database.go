package configs

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string `toml:"db_host"`
	Port     string `toml:"db_port"`
	User     string `toml:"db_user"`
	Password string `toml:"db_password"`
	Name     string `toml:"db_name"`
}

func loadDatabaseConfig(config DatabaseConfig) DatabaseConfig {
	return DatabaseConfig{
		Host:     envOrDefaultString(os.Getenv("DB_HOST"), config.Host),
		Port:     envOrDefaultString(os.Getenv("DB_PORT"), config.Port),
		User:     envOrDefaultString(os.Getenv("DB_USER"), config.User),
		Password: envOrDefaultString(os.Getenv("DB_PASSWORD"), config.Password),
		Name:     envOrDefaultString(os.Getenv("DB_NAME"), config.Name),
	}
}

func GetDbConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:     globalConfig.Database.Host,
		Port:     globalConfig.Database.Port,
		User:     globalConfig.Database.User,
		Password: globalConfig.Database.Password,
		Name:     globalConfig.Database.Name,
	}
}

func GetDbString() string {
	conf := GetDbConfig()

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", conf.Host, conf.User, conf.Password, conf.Name, conf.Port)
}
