package store

import (
	"context"
	"time"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/Sahil2k07/go-base/internal/database"
	"github.com/Sahil2k07/go-base/internal/features/user"
	"github.com/Sahil2k07/go-base/internal/gen"
	"github.com/Sahil2k07/go-base/internal/models"
	"github.com/google/uuid"
)

type userSessionStore struct {
	q             *gen.Query
	refreshExpiry time.Duration
}

func NewUserSessionRepository() user.UserSessionRepository {
	return &userSessionStore{
		q:             database.GQ,
		refreshExpiry: configs.GetServerConfig().RefreshTokenExpiry,
	}
}

func (s *userSessionStore) CreateNewUserSession(userID uint, sessionID uuid.UUID, refreshTokenHash string) (*models.UserSession, error) {
	session := &models.UserSession{
		UserID:            userID,
		SessionIdentifier: sessionID,
		RefreshTokenHash:  refreshTokenHash,
		ExpiresAt:         time.Now().Add(s.refreshExpiry),
	}

	if err := s.q.UserSession.WithContext(context.Background()).Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userSessionStore) GetUserSessionsByUserID(userID uint) ([]models.UserSession, error) {
	us := s.q.UserSession

	sessions, err := us.WithContext(context.Background()).
		Where(us.UserID.Eq(userID), us.ExpiresAt.Gt(time.Now())).
		Find()

	if err != nil {
		return nil, err
	}

	return dereference(sessions), nil
}

func (s *userSessionStore) GetUserSessionBySessionID(sessionID uuid.UUID) (*models.UserSession, error) {
	us := s.q.UserSession

	session, err := us.WithContext(context.Background()).
		Where(us.SessionIdentifier.Eq(sessionID), us.ExpiresAt.Gt(time.Now())).
		First()

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *userSessionStore) DeleteUserSession(sessionID uuid.UUID) error {
	us := s.q.UserSession

	_, err := us.WithContext(context.Background()).
		Unscoped().
		Where(us.SessionIdentifier.Eq(sessionID)).
		Delete()

	return err
}

func (s *userSessionStore) RotateUserSession(userID uint, sessionID uuid.UUID, oldHash, newHash string) (bool, error) {
	us := s.q.UserSession
	now := time.Now()

	info, err := us.WithContext(context.Background()).
		Where(
			us.UserID.Eq(userID),
			us.SessionIdentifier.Eq(sessionID),
			us.RefreshTokenHash.Eq(oldHash),
			us.ExpiresAt.Gt(now),
		).
		UpdateSimple(
			us.RefreshTokenHash.Value(newHash),
			us.ExpiresAt.Value(now.Add(s.refreshExpiry)),
		)

	if err != nil {
		return false, err
	}

	return info.RowsAffected > 0, nil
}
