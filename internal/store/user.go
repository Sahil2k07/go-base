package store

import (
	"context"

	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/features/user"
	"github.com/Sahil2k07/go-base/internal/gen"
	"github.com/Sahil2k07/go-base/internal/models"
)

type userStore struct {
	q *gen.Query
}

func NewUserRepository() user.UserRepository {
	return &userStore{q: database.GQ}
}

func (s *userStore) CreateNewUser(email, password, firstName, lastName, phone, gender string) (*models.User, error) {
	user := &models.User{
		Email:        email,
		PasswordHash: password,
		Profile: models.UserProfile{
			FirstName: firstName,
			LastName:  lastName,
			Phone:     phone,
			Gender:    gender,
		},
	}

	if err := s.q.User.WithContext(context.Background()).Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userStore) UserHasPermission(userID uint, permission string) (bool, error) {
	p, rp, ur := s.q.Permission, s.q.RolePermission, s.q.UserRole

	count, err := p.WithContext(context.Background()).
		Join(rp, rp.PermissionID.EqCol(p.ID)).
		Join(ur, ur.RoleID.EqCol(rp.RoleID)).
		Where(ur.UserID.Eq(userID), p.Name.Eq(permission)).
		Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *userStore) UserHasRole(userID uint, role string) (bool, error) {
	r, ur := s.q.Role, s.q.UserRole

	count, err := r.WithContext(context.Background()).
		Join(ur, ur.RoleID.EqCol(r.ID)).
		Where(ur.UserID.Eq(userID), r.Name.Eq(role)).
		Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *userStore) GetUserByEmail(email string) (*models.User, error) {
	u := s.q.User

	found, err := u.WithContext(context.Background()).
		Where(u.Email.Eq(email)).
		Preload(u.Profile).
		First()

	if err != nil {
		return nil, err
	}

	return found, nil
}

func (s *userStore) GetUserByUserID(userID uint) (*models.User, error) {
	u := s.q.User

	found, err := u.WithContext(context.Background()).
		Where(u.ID.Eq(userID)).
		Preload(u.Profile).
		First()

	if err != nil {
		return nil, err
	}

	return found, nil
}

func (s *userStore) IsEmailAlreadyRegistered(email string) (bool, error) {
	u := s.q.User

	count, err := u.WithContext(context.Background()).
		Where(u.Email.Eq(email)).
		Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *userStore) GetUserRolesByUserID(userID uint) ([]models.Role, error) {
	r, ur := s.q.Role, s.q.UserRole

	roles, err := r.WithContext(context.Background()).
		Select(r.ALL).
		Join(ur, ur.RoleID.EqCol(r.ID)).
		Where(ur.UserID.Eq(userID)).
		Find()

	if err != nil {
		return nil, err
	}

	return dereference(roles), nil
}

func (s *userStore) GetUserRoleNamesByUserID(userID uint) ([]string, error) {
	r, ur := s.q.Role, s.q.UserRole

	names := make([]string, 0)

	err := r.WithContext(context.Background()).
		Join(ur, ur.RoleID.EqCol(r.ID)).
		Where(ur.UserID.Eq(userID)).
		Pluck(r.Name, &names)

	if err != nil {
		return nil, err
	}

	return names, nil
}

func (s *userStore) GetUserPermissionsByUserID(userID uint) ([]models.Permission, error) {
	p, rp, ur := s.q.Permission, s.q.RolePermission, s.q.UserRole

	permissions, err := p.WithContext(context.Background()).
		Distinct(p.ALL).
		Join(rp, rp.PermissionID.EqCol(p.ID)).
		Join(ur, ur.RoleID.EqCol(rp.RoleID)).
		Where(ur.UserID.Eq(userID)).
		Find()

	if err != nil {
		return nil, err
	}

	return dereference(permissions), nil
}

func (s *userStore) GetUserPermissionNamesByUserID(userID uint) ([]string, error) {
	p, rp, ur := s.q.Permission, s.q.RolePermission, s.q.UserRole

	names := make([]string, 0)

	err := p.WithContext(context.Background()).
		Distinct(p.Name).
		Join(rp, rp.PermissionID.EqCol(p.ID)).
		Join(ur, ur.RoleID.EqCol(rp.RoleID)).
		Where(ur.UserID.Eq(userID)).
		Pluck(p.Name, &names)

	if err != nil {
		return nil, err
	}

	return names, nil
}

func dereference[T any](items []*T) []T {
	out := make([]T, len(items))
	for i, item := range items {
		out[i] = *item
	}

	return out
}
