package user

import (
	"github.com/Sahil2k07/go-base/internal/models"
	"github.com/google/uuid"
)

type UserRepository interface {
	IsEmailAlreadyRegistered(email string) (bool, error)

	GetUserByUserID(userID uint) (*models.User, error)

	GetUserByEmail(email string) (*models.User, error)

	GetUserRolesByUserID(userID uint) ([]models.Role, error)

	GetUserPermissionsByUserID(userID uint) ([]models.Permission, error)

	GetUserRoleNamesByUserID(userID uint) ([]string, error)

	GetUserPermissionNamesByUserID(userID uint) ([]string, error)

	UserHasRole(userID uint, role string) (bool, error)

	UserHasPermission(userID uint, permission string) (bool, error)

	CreateNewUser(email, password, firstName, lastName, phone, gender string) (*models.User, error)
}

type UserSessionRepository interface {
	CreateNewUserSession(userID uint, sessionID uuid.UUID, refreshTokenHash string) (*models.UserSession, error)

	GetUserSessionsByUserID(userID uint) ([]models.UserSession, error)

	GetUserSessionBySessionID(sessionID uuid.UUID) (*models.UserSession, error)

	RotateUserSession(userID uint, sessionID uuid.UUID, oldHash, newHash string) (bool, error)

	DeleteUserSession(sessionID uuid.UUID) error
}
