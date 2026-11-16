package auth

import (
	"errors"
	"log/slog"

	app_errors "github.com/Sahil2k07/go-base/internal/errors"
	"github.com/Sahil2k07/go-base/internal/features/user"
	"github.com/Sahil2k07/go-base/internal/infrastructure"
	"github.com/Sahil2k07/go-base/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	ur  user.UserRepository
	usr user.UserSessionRepository
	cs  *infrastructure.CryptoService
}

func NewAuthService(ur user.UserRepository, usr user.UserSessionRepository, cs *infrastructure.CryptoService) *AuthService {
	return &AuthService{ur: ur, usr: usr, cs: cs}
}

func (s *AuthService) SignupUser(req SignupRequest) (TokenResponse, error) {
	alreadyExists, err := s.ur.IsEmailAlreadyRegistered(req.Email)
	if err != nil {
		return TokenResponse{}, err
	}

	if alreadyExists {
		return TokenResponse{}, app_errors.Conflict("Email is already registered.")
	}

	passwordHash, err := s.cs.HashPassword(req.Password)
	if err != nil {
		return TokenResponse{}, app_errors.Internal(err)
	}

	u, err := s.ur.CreateNewUser(req.Email, passwordHash, req.FirstName, req.LastName, req.Phone, req.Gender)
	if err != nil {
		return TokenResponse{}, err
	}

	return s.issueTokens(u)
}

func (s *AuthService) SigninUser(req SigninRequest) (TokenResponse, error) {
	u, err := s.ur.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TokenResponse{}, app_errors.Unauthorized("Invalid email or password").Wrap(err)
		}

		return TokenResponse{}, err
	}

	if err := s.cs.CheckPassword(u.PasswordHash, req.Password); err != nil {
		return TokenResponse{}, app_errors.Unauthorized("Invalid email or password").Wrap(err)
	}

	return s.issueTokens(u)
}

func (s *AuthService) RefreshToken(req RefreshTokenRequest) (TokenResponse, error) {
	claims, err := s.cs.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		return TokenResponse{}, app_errors.Unauthorized("Invalid or expired refresh token").Wrap(err)
	}

	sessionID, err := uuid.Parse(claims.SessionID)
	if err != nil {
		return TokenResponse{}, app_errors.Unauthorized("Invalid or expired refresh token").Wrap(err)
	}

	u, err := s.ur.GetUserByUserID(claims.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return TokenResponse{}, app_errors.Unauthorized("Invalid or expired refresh token").Wrap(err)
		}

		return TokenResponse{}, err
	}

	newRefreshToken, err := s.cs.GenerateRefreshToken(u.ID, sessionID.String())
	if err != nil {
		return TokenResponse{}, app_errors.Internal(err)
	}

	rotated, err := s.usr.RotateUserSession(
		u.ID,
		sessionID,
		s.cs.HashRefreshToken(req.RefreshToken),
		s.cs.HashRefreshToken(newRefreshToken),
	)
	if err != nil {
		return TokenResponse{}, err
	}

	if !rotated {
		if err := s.usr.DeleteUserSession(sessionID); err != nil {
			slog.Error("failed to revoke session after refresh token reuse", "session_id", sessionID, "error", err)
		}

		return TokenResponse{}, app_errors.Unauthorized("Invalid or expired refresh token")
	}

	accessToken, err := s.cs.GenerateAccessToken(u.ID, u.Profile.ID, u.Email)
	if err != nil {
		return TokenResponse{}, app_errors.Internal(err)
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) issueTokens(u *models.User) (TokenResponse, error) {
	sessionID := uuid.New()

	refreshToken, err := s.cs.GenerateRefreshToken(u.ID, sessionID.String())
	if err != nil {
		return TokenResponse{}, app_errors.Internal(err)
	}

	if _, err := s.usr.CreateNewUserSession(u.ID, sessionID, s.cs.HashRefreshToken(refreshToken)); err != nil {
		return TokenResponse{}, err
	}

	accessToken, err := s.cs.GenerateAccessToken(u.ID, u.Profile.ID, u.Email)
	if err != nil {
		return TokenResponse{}, app_errors.Internal(err)
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
