package store

import "github.com/Sahil2k07/go-base/internal/features/user"

type userStore struct{}

func NewUserRepository() user.UserRepository {
	return &userStore{}
}
