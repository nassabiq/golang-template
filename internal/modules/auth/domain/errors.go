package domain

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrTokenRevoked         = errors.New("token revoked")
	ErrPasswordResetExpired = errors.New("password reset expired")
	ErrPasswordResetUsed    = errors.New("password reset used")
	ErrPasswordUpdateFailed = errors.New("password update failed")
	ErrPasswordNotMatch     = errors.New("password not match")
	ErrEmailAlreadyUsed     = errors.New("email already registered")
	ErrWeakPassword         = errors.New("weak password")
)
