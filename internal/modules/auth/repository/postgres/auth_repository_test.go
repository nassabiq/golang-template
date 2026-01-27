package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	cleanup := func() {
		db.Close()
	}

	return db, mock, cleanup
}

// Test CreateUser
func TestAuthRepository_CreateUser(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)

	tests := []struct {
		name    string
		user    *domain.User
		mock    func()
		wantErr bool
	}{
		{
			name: "success - create user",
			user: &domain.User{
				ID:           "user-123",
				Name:         "Test User",
				Email:        "test@example.com",
				PasswordHash: "hashed-password",
				RoleID:       "user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("user-123", "Test User", "test@example.com", "hashed-password", "user", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "failure - duplicate email",
			user: &domain.User{
				ID:           "user-123",
				Name:         "Test User",
				Email:        "existing@example.com",
				PasswordHash: "hashed-password",
				RoleID:       "user",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("user-123", "Test User", "existing@example.com", "hashed-password", "user", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.CreateUser(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test FindUserByEmail
func TestAuthRepository_FindUserByEmail(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		email   string
		mock    func()
		want    *domain.User
		wantErr bool
	}{
		{
			name:  "success - user found",
			email: "test@example.com",
			mock: func() {
				// Query returns: id, email, password, role_id, role_name, created_at (6 columns)
				// But repo.Scan expects 7 fields including Name and UpdatedAt
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role_id", "created_at", "updated_at"}).
					AddRow("user-123", "Test User", "test@example.com", "hashed-password", "user", fixedTime, fixedTime)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			want: &domain.User{
				ID:           "user-123",
				Name:         "Test User",
				Email:        "test@example.com",
				PasswordHash: "hashed-password",
				RoleID:       "user",
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			wantErr: false,
		},
		{
			name:  "failure - user not found",
			email: "nonexistent@example.com",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindUserByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID || got.Email != tt.want.Email {
					t.Errorf("FindUserByEmail() = %v, want %v", got, tt.want)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test FindUserByID
func TestAuthRepository_FindUserByID(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		id      string
		mock    func()
		want    *domain.User
		wantErr bool
	}{
		{
			name: "success - user found",
			id:   "user-123",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role_id", "created_at", "updated_at"}).
					AddRow("user-123", "Test User", "test@example.com", "hashed-password", "user", fixedTime, fixedTime)
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("user-123").
					WillReturnRows(rows)
			},
			want: &domain.User{
				ID:           "user-123",
				Name:         "Test User",
				Email:        "test@example.com",
				PasswordHash: "hashed-password",
				RoleID:       "user",
				CreatedAt:    fixedTime,
				UpdatedAt:    fixedTime,
			},
			wantErr: false,
		},
		{
			name: "failure - user not found",
			id:   "nonexistent-id",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM users").
					WithArgs("nonexistent-id").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindUserByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("FindUserByID() = %v, want %v", got, tt.want)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test StoreRefreshToken
func TestAuthRepository_StoreRefreshToken(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		token   *domain.RefreshToken
		mock    func()
		wantErr bool
	}{
		{
			name: "success - store token",
			token: &domain.RefreshToken{
				ID:        "token-123",
				UserID:    "user-123",
				TokenHash: "hash-123",
				Revoked:   false,
				ExpiresAt: fixedTime.Add(24 * time.Hour),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mock: func() {
				// Query has 7 params: id, user_id, token_hash, revoked, expires_at, created_at, updated_at
				mock.ExpectExec("INSERT INTO refresh_tokens").
					WithArgs("token-123", "user-123", "hash-123", false, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.StoreRefreshToken(context.Background(), tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("StoreRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test FindValidRefreshToken
func TestAuthRepository_FindValidRefreshToken(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		tokenHash string
		mock      func()
		want      *domain.RefreshToken
		wantErr   bool
	}{
		{
			name:      "success - token found",
			tokenHash: "hash-123",
			mock: func() {
				// Scan expects 7 fields: ID, UserID, TokenHash, Revoked, ExpiresAt, CreatedAt, UpdatedAt
				rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "revoked", "expires_at", "created_at", "updated_at"}).
					AddRow("token-123", "user-123", "hash-123", false, fixedTime.Add(24*time.Hour), fixedTime, fixedTime)
				mock.ExpectQuery("SELECT (.+) FROM refresh_tokens").
					WithArgs("hash-123").
					WillReturnRows(rows)
			},
			want: &domain.RefreshToken{
				ID:        "token-123",
				UserID:    "user-123",
				TokenHash: "hash-123",
				Revoked:   false,
				ExpiresAt: fixedTime.Add(24 * time.Hour),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			wantErr: false,
		},
		{
			name:      "success - token not found (no error)",
			tokenHash: "nonexistent-hash",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM refresh_tokens").
					WithArgs("nonexistent-hash").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindValidRefreshToken(context.Background(), tt.tokenHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindValidRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("FindValidRefreshToken() = %v, want %v", got, tt.want)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test RevokeRefreshToken
func TestAuthRepository_RevokeRefreshToken(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)

	tests := []struct {
		name      string
		tokenHash string
		mock      func()
		wantErr   bool
	}{
		{
			name:      "success - revoke token",
			tokenHash: "hash-123",
			mock: func() {
				mock.ExpectExec("UPDATE refresh_tokens SET revoked").
					WithArgs("hash-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.RevokeRefreshToken(context.Background(), tt.tokenHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("RevokeRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test StorePasswordReset
func TestAuthRepository_StorePasswordReset(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		passwordReset *domain.PasswordReset
		mock          func()
		wantErr       bool
	}{
		{
			name: "success - store password reset",
			passwordReset: &domain.PasswordReset{
				ID:        "reset-123",
				UserID:    "user-123",
				TokenHash: "hash-123",
				ExpiresAt: fixedTime.Add(15 * time.Minute),
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			mock: func() {
				mock.ExpectExec("INSERT INTO password_resets").
					WithArgs("reset-123", "user-123", "hash-123", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.StorePasswordReset(context.Background(), tt.passwordReset)
			if (err != nil) != tt.wantErr {
				t.Errorf("StorePasswordReset() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test FindValidPasswordReset
func TestAuthRepository_FindValidPasswordReset(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		tokenHash string
		mock      func()
		want      *domain.PasswordReset
		wantErr   bool
	}{
		{
			name:      "success - password reset found",
			tokenHash: "hash-123",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "token_hash", "expires_at", "used", "created_at", "updated_at"}).
					AddRow("reset-123", "user-123", "hash-123", fixedTime.Add(15*time.Minute), false, fixedTime, fixedTime)
				mock.ExpectQuery("SELECT (.+) FROM password_resets").
					WithArgs("hash-123").
					WillReturnRows(rows)
			},
			want: &domain.PasswordReset{
				ID:        "reset-123",
				UserID:    "user-123",
				TokenHash: "hash-123",
				ExpiresAt: fixedTime.Add(15 * time.Minute),
				Used:      false,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			wantErr: false,
		},
		{
			name:      "success - not found (no error)",
			tokenHash: "nonexistent-hash",
			mock: func() {
				mock.ExpectQuery("SELECT (.+) FROM password_resets").
					WithArgs("nonexistent-hash").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindValidPasswordReset(context.Background(), tt.tokenHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindValidPasswordReset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil && got != nil {
				if got.ID != tt.want.ID {
					t.Errorf("FindValidPasswordReset() = %v, want %v", got, tt.want)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

// Test UpdateUserPassword
func TestAuthRepository_UpdateUserPassword(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)

	tests := []struct {
		name    string
		userID  string
		newHash string
		mock    func()
		wantErr bool
	}{
		{
			name:    "success - update password",
			userID:  "user-123",
			newHash: "new-hashed-password",
			mock: func() {
				mock.ExpectExec("UPDATE users SET password").
					WithArgs("new-hashed-password", "user-123").
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.UpdateUserPassword(context.Background(), tt.userID, tt.newHash)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}
