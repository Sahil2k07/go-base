package apps

import (
	"github.com/Sahil2k07/go-base/internal/models"
	"gorm.io/gen"
)

func GenerateGenModels() {
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
