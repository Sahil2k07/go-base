package main

import (
	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/models"
	"gorm.io/gen"
)

func init() {
	configs.LoadConfigs("configs/migration.toml")
	infrastructure.ConfigureLogging("gen")
	database.ConnectPostgres()
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/gen",
	})

	g.ApplyBasic(
		models.User{},
		models.UserProfile{},
		models.UserSession{},
		models.Role{},
		models.Permission{},
		models.UserRole{},
		models.RolePermission{},
	)

	g.Execute()
}
