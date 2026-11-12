package infrastructure

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/Sahil2k07/go-base/internal/configs"
	"github.com/go-jose/go-jose/v4"
	josejwt "github.com/go-jose/go-jose/v4/jwt"
	"golang.org/x/crypto/bcrypt"
)

type cryptoService struct {
	config configs.ServerConfig
}

type AccessTokenClaims struct {
	UserID    uint   `json:"user_id"`
	ProfileID uint   `json:"profile_id"`
	Email     string `json:"email"`
}

type RefreshTokenClaims struct {
	UserID    uint   `json:"user_id"`
	SessionID string `json:"session_id"`
}

// Internal payloads: registered claims (exp, iat, jti) plus our own claims.
// Callers only ever see the clean structs above.
type accessPayload struct {
	josejwt.Claims
	AccessTokenClaims
}

type refreshPayload struct {
	josejwt.Claims
	RefreshTokenClaims
}

func NewCryptoService(config configs.ServerConfig) *cryptoService {
	return &cryptoService{
		config: config,
	}
}

// HashPassword hashes a password using bcrypt.
func (c *cryptoService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		c.config.BcryptCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// CheckPassword verifies a password against its bcrypt hash.
func (c *cryptoService) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}

// GenerateAccessToken generates a signed and encrypted JWT.
func (c *cryptoService) GenerateAccessToken(userID, profileID uint, email string) (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       []byte(c.config.AccessTokenSecret),
		},
		nil,
	)
	if err != nil {
		return "", err
	}

	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.DIRECT,
			Key:       []byte(c.config.EncryptionSecret),
		},
		(&jose.EncrypterOptions{}).WithContentType("JWT"),
	)
	if err != nil {
		return "", err
	}

	now := time.Now()

	return josejwt.SignedAndEncrypted(signer, encrypter).
		Claims(accessPayload{
			Claims: josejwt.Claims{
				IssuedAt: josejwt.NewNumericDate(now),
				Expiry:   josejwt.NewNumericDate(now.Add(c.config.AccessTokenExpiry)),
			},
			AccessTokenClaims: AccessTokenClaims{
				UserID:    userID,
				ProfileID: profileID,
				Email:     email,
			},
		}).
		Serialize()
}

// VerifyAccessToken decrypts, verifies and validates an access token.
func (c *cryptoService) VerifyAccessToken(token string) (*AccessTokenClaims, error) {
	parsedToken, err := josejwt.ParseSignedAndEncrypted(
		token,
		[]jose.KeyAlgorithm{jose.DIRECT},
		[]jose.ContentEncryption{jose.A256GCM},
		[]jose.SignatureAlgorithm{jose.HS256},
	)
	if err != nil {
		return nil, err
	}

	innerToken, err := parsedToken.Decrypt([]byte(c.config.EncryptionSecret))
	if err != nil {
		return nil, err
	}

	payload := new(accessPayload)

	if err := innerToken.Claims([]byte(c.config.AccessTokenSecret), payload); err != nil {
		return nil, err
	}

	if err := payload.Claims.Validate(josejwt.Expected{Time: time.Now()}); err != nil {
		return nil, err
	}

	return &payload.AccessTokenClaims, nil
}

// GenerateRefreshToken generates a signed JWT carrying the user and session IDs.
func (c *cryptoService) GenerateRefreshToken(userID uint, sessionID string) (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       []byte(c.config.RefreshTokenSecret),
		},
		nil,
	)
	if err != nil {
		return "", err
	}

	// Random jti so two tokens issued in the same second are never identical.
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}

	now := time.Now()

	return josejwt.Signed(signer).
		Claims(refreshPayload{
			Claims: josejwt.Claims{
				ID:       base64.RawURLEncoding.EncodeToString(idBytes),
				IssuedAt: josejwt.NewNumericDate(now),
				Expiry:   josejwt.NewNumericDate(now.Add(c.config.RefreshTokenExpiry)),
			},
			RefreshTokenClaims: RefreshTokenClaims{
				UserID:    userID,
				SessionID: sessionID,
			},
		}).
		Serialize()
}

// VerifyRefreshToken verifies the signature and expiry of a refresh token.
func (c *cryptoService) VerifyRefreshToken(token string) (*RefreshTokenClaims, error) {
	parsedToken, err := josejwt.ParseSigned(
		token,
		[]jose.SignatureAlgorithm{jose.HS256},
	)
	if err != nil {
		return nil, err
	}

	payload := new(refreshPayload)

	if err := parsedToken.Claims([]byte(c.config.RefreshTokenSecret), payload); err != nil {
		return nil, err
	}

	if err := payload.Claims.Validate(josejwt.Expected{Time: time.Now()}); err != nil {
		return nil, err
	}

	return &payload.RefreshTokenClaims, nil
}

// HashRefreshToken hashes a refresh token before database storage.
func (c *cryptoService) HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))

	return base64.RawURLEncoding.EncodeToString(hash[:])
}
