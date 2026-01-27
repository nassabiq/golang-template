package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service provides JWT token generation
type Service struct {
	secret          []byte
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewService creates a new token service
func NewService(secret string) *Service {
	return &Service{
		secret:          []byte(secret),
		accessTokenTTL:  15 * time.Minute,
		refreshTokenTTL: 7 * 24 * time.Hour,
	}
}

// GenerateAccessToken creates a new JWT access token
func (s *Service) GenerateAccessToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(s.accessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken creates a new refresh token
// Returns: plain token (for client), hashed token (for storage), error
func (s *Service) GenerateRefreshToken() (plain string, hash string, err error) {
	// Generate random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	plain = base64.URLEncoding.EncodeToString(bytes)

	// Hash for storage (SHA256 for deterministic key)
	hashBytes := sha256.Sum256([]byte(plain))

	return plain, hex.EncodeToString(hashBytes[:]), nil
}
