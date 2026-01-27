package domain

import "context"

type AuthRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)

	// ===== REFRESH TOKEN =====
	StoreRefreshToken(ctx context.Context, token *RefreshToken) error
	FindValidRefreshToken(ctx context.Context, tokenHash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllRefreshTokens(ctx context.Context, userID string) error

	// ===== PASSWORD RESET =====
	StorePasswordReset(ctx context.Context, pr *PasswordReset) error
	FindValidPasswordReset(ctx context.Context, tokenHash string) (*PasswordReset, error)
	MarkPasswordResetUsed(ctx context.Context, id string) error

	// ===== PASSWORD UPDATE =====
	UpdateUserPassword(ctx context.Context, userID, newHash string) error
}
