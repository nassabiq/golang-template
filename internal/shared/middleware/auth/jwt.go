package auth

import (
	"errors"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWTVerifier struct {
	secret []byte
}

func NewJWTVerifier(secret string) *JWTVerifier {
	return &JWTVerifier{secret: []byte(secret)}
}

func (j *JWTVerifier) Verify(tokenString string) (userID, role string, err error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.secret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	userID = claims["sub"].(string)
	role = claims["role"].(string)

	if userID == "" || role == "" {
		return "", "", errors.New("missing claims")
	}

	return userID, role, nil
}
