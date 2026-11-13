package store

import "github.com/Sahil2k07/go-base/internal/features/auth"

type authStore struct{}

func NewAuthRepository() auth.AuthRepository {
	return &authStore{}
}
