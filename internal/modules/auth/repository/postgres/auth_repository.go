package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	"github.com/nassabiq/golang-template/internal/shared/helper"
)

type AuthRepository struct {
	db      *sql.DB
	queries map[string]string
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db:      db,
		queries: helper.LoadQuery(rawQuery),
	}
}

func (repository *AuthRepository) query(name string) string {
	query, ok := repository.queries[name]

	if !ok {
		panic("query not found: " + name)
	}

	return query
}

func (repository *AuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	// Check if email already exists
	existing, _ := repository.FindUserByEmail(ctx, user.Email)
	if existing != nil {
		return domain.ErrEmailAlreadyUsed
	}

	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("CreateUser"),
		user.ID, user.Name, user.Email, user.PasswordHash, user.RoleID, user.CreatedAt, user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// RUN QUERY
	row := repository.db.QueryRowContext(ctx, repository.query("FindUserByEmail"), email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (repository *AuthRepository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	// RUN QUERY
	row := repository.db.QueryRowContext(ctx, repository.query("FindUserByID"), id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *AuthRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("StoreRefreshToken"),
		token.ID, token.UserID, token.TokenHash, token.Revoked, token.ExpiresAt, token.CreatedAt, token.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repository *AuthRepository) FindValidRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	// RUN QUERY
	row := repository.db.QueryRowContext(ctx, repository.query("FindValidRefreshToken"), tokenHash)

	var token domain.RefreshToken
	if err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.Revoked,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &token, nil
}

func (repository *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("RevokeRefreshToken"), tokenHash)
	return err
}

func (repository *AuthRepository) RevokeAllRefreshTokens(ctx context.Context, userID string) error {
	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("RevokeAllRefreshTokens"), userID)
	return err
}

func (repository *AuthRepository) StorePasswordReset(ctx context.Context, passwordReset *domain.PasswordReset) error {
	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("StorePasswordReset"),
		passwordReset.ID,
		passwordReset.UserID,
		passwordReset.TokenHash,
		passwordReset.ExpiresAt,
		passwordReset.CreatedAt,
		passwordReset.UpdatedAt,
	)

	return err
}

func (repository *AuthRepository) FindValidPasswordReset(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	// RUN QUERY
	row := repository.db.QueryRowContext(ctx, repository.query("FindValidPasswordReset"), tokenHash)

	var passwordReset domain.PasswordReset
	if err := row.Scan(
		&passwordReset.ID,
		&passwordReset.UserID,
		&passwordReset.TokenHash,
		&passwordReset.ExpiresAt,
		&passwordReset.Used,
		&passwordReset.CreatedAt,
		&passwordReset.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &passwordReset, nil
}

func (repository *AuthRepository) MarkPasswordResetUsed(ctx context.Context, id string) error {
	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("MarkPasswordResetUsed"), id)
	return err
}

func (repository *AuthRepository) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
	// RUN QUERY
	_, err := repository.db.ExecContext(ctx, repository.query("UpdateUserPassword"), newHash, userID)
	return err
}
