package domain

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	RoleID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	Revoked   bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PasswordReset struct {
	ID        string
	UserID    string
	TokenHash string
	Used      bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
