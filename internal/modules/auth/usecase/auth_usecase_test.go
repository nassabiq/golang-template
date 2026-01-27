package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	"github.com/nassabiq/golang-template/internal/modules/auth/event"
)

// Mock implementations
type mockAuthRepository struct {
	users                  map[string]*domain.User
	refreshTokens          map[string]*domain.RefreshToken
	passwordResets         map[string]*domain.PasswordReset
	findUserByEmail        func(email string) (*domain.User, error)
	findUserByID           func(id string) (*domain.User, error)
	createUser             func(user *domain.User) error
	storeRefreshToken      func(token *domain.RefreshToken) error
	findValidRefreshToken  func(tokenHash string) (*domain.RefreshToken, error)
	revokeRefreshToken     func(tokenHash string) error
	storePasswordReset     func(pr *domain.PasswordReset) error
	findValidPasswordReset func(tokenHash string) (*domain.PasswordReset, error)
	markPasswordResetUsed  func(id string) error
	updateUserPassword     func(userID, newHash string) error
}

func (m *mockAuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if m.createUser != nil {
		return m.createUser(user)
	}
	if m.users == nil {
		m.users = make(map[string]*domain.User)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockAuthRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.findUserByEmail != nil {
		return m.findUserByEmail(email)
	}
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (m *mockAuthRepository) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	if m.findUserByID != nil {
		return m.findUserByID(id)
	}
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockAuthRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	if m.storeRefreshToken != nil {
		return m.storeRefreshToken(token)
	}
	if m.refreshTokens == nil {
		m.refreshTokens = make(map[string]*domain.RefreshToken)
	}
	m.refreshTokens[token.TokenHash] = token
	return nil
}

func (m *mockAuthRepository) FindValidRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	if m.findValidRefreshToken != nil {
		return m.findValidRefreshToken(tokenHash)
	}
	if t, ok := m.refreshTokens[tokenHash]; ok && !t.Revoked && t.ExpiresAt.After(time.Now()) {
		return t, nil
	}
	return nil, errors.New("token not found")
}

func (m *mockAuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	if m.revokeRefreshToken != nil {
		return m.revokeRefreshToken(tokenHash)
	}
	if t, ok := m.refreshTokens[tokenHash]; ok {
		t.Revoked = true
	}
	return nil
}

func (m *mockAuthRepository) RevokeAllRefreshTokens(ctx context.Context, userID string) error {
	for _, t := range m.refreshTokens {
		if t.UserID == userID {
			t.Revoked = true
		}
	}
	return nil
}

func (m *mockAuthRepository) StorePasswordReset(ctx context.Context, pr *domain.PasswordReset) error {
	if m.storePasswordReset != nil {
		return m.storePasswordReset(pr)
	}
	if m.passwordResets == nil {
		m.passwordResets = make(map[string]*domain.PasswordReset)
	}
	m.passwordResets[pr.TokenHash] = pr
	return nil
}

func (m *mockAuthRepository) FindValidPasswordReset(ctx context.Context, tokenHash string) (*domain.PasswordReset, error) {
	if m.findValidPasswordReset != nil {
		return m.findValidPasswordReset(tokenHash)
	}
	if pr, ok := m.passwordResets[tokenHash]; ok && !pr.Used && pr.ExpiresAt.After(time.Now()) {
		return pr, nil
	}
	return nil, errors.New("password reset not found")
}

func (m *mockAuthRepository) MarkPasswordResetUsed(ctx context.Context, id string) error {
	if m.markPasswordResetUsed != nil {
		return m.markPasswordResetUsed(id)
	}
	for _, pr := range m.passwordResets {
		if pr.UserID == id {
			pr.Used = true
		}
	}
	return nil
}

func (m *mockAuthRepository) UpdateUserPassword(ctx context.Context, userID, newHash string) error {
	if m.updateUserPassword != nil {
		return m.updateUserPassword(userID, newHash)
	}
	if u, ok := m.users[userID]; ok {
		u.PasswordHash = newHash
	}
	return nil
}

type mockTokenService struct {
	generateAccessToken  func(userID, role string) (string, error)
	generateRefreshToken func() (plain string, hash string, err error)
}

func (m *mockTokenService) GenerateAccessToken(userID, role string) (string, error) {
	if m.generateAccessToken != nil {
		return m.generateAccessToken(userID, role)
	}
	return "access-token-" + userID, nil
}

func (m *mockTokenService) GenerateRefreshToken() (plain string, hash string, err error) {
	if m.generateRefreshToken != nil {
		return m.generateRefreshToken()
	}
	return "refresh-token-plain", "refresh-token-hash", nil
}

type mockPasswordHasher struct {
	hashPassword        func(password string) (string, error)
	verifyPassword      func(password, hash string) bool
	generateRandomToken func() (string, error)
	hashToken           func(token string) string
}

func (m *mockPasswordHasher) HashPassword(password string) (string, error) {
	if m.hashPassword != nil {
		return m.hashPassword(password)
	}
	return "hashed-" + password, nil
}

func (m *mockPasswordHasher) VerifyPassword(password, hash string) bool {
	if m.verifyPassword != nil {
		return m.verifyPassword(password, hash)
	}
	return hash == "hashed-"+password
}

func (m *mockPasswordHasher) GenerateRandomToken() (string, error) {
	if m.generateRandomToken != nil {
		return m.generateRandomToken()
	}
	return "random-token", nil
}

func (m *mockPasswordHasher) HashToken(token string) string {
	if m.hashToken != nil {
		return m.hashToken(token)
	}
	return "sha256-" + token
}

type mockUUIDGenerator struct {
	id string
}

func (m *mockUUIDGenerator) GenerateID() string {
	return m.id
}

type mockEventPublisher struct {
	forgotPasswordCalled bool
	forgotPasswordEvent  event.ForgotPasswordEvent
}

func (m *mockEventPublisher) ForgotPassword(e event.ForgotPasswordEvent) error {
	m.forgotPasswordCalled = true
	m.forgotPasswordEvent = e
	return nil
}

func setupTestUsecase() (*AuthUsecase, *mockAuthRepository, *mockTokenService, *mockPasswordHasher, *mockUUIDGenerator, *mockEventPublisher) {
	repo := &mockAuthRepository{
		users:          make(map[string]*domain.User),
		refreshTokens:  make(map[string]*domain.RefreshToken),
		passwordResets: make(map[string]*domain.PasswordReset),
	}
	tokenSvc := &mockTokenService{}
	hasher := &mockPasswordHasher{}
	uuid := &mockUUIDGenerator{id: "test-uuid-123"}
	eventPub := &mockEventPublisher{}

	// Create usecase with all dependencies
	uc := NewAuthUsecase(repo, nil)
	uc.token = tokenSvc
	uc.now = func() time.Time { return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) }
	uc.passwordHasher = hasher
	uc.uuid = uuid
	// eventPub will be set in individual tests that need it

	return uc, repo, tokenSvc, hasher, uuid, eventPub
}

// Test Register
func TestAuthUsecase_Register(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.RegisterInput
		setup   func(*mockAuthRepository, *mockPasswordHasher)
		wantErr error
	}{
		{
			name: "success - valid registration",
			input: domain.RegisterInput{
				Name:                 "Test User",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: nil,
		},
		{
			name: "failure - email already registered",
			input: domain.RegisterInput{
				Name:                 "Test User",
				Email:                "existing@example.com",
				Password:             "password123",
				PasswordConfirmation: "password123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return &domain.User{Email: email}, nil
				}
			},
			wantErr: domain.ErrEmailAlreadyUsed,
		},
		{
			name: "failure - password confirmation mismatch",
			input: domain.RegisterInput{
				Name:                 "Test User",
				Email:                "test@example.com",
				Password:             "password123",
				PasswordConfirmation: "different",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return nil, errors.New("not found")
				}
			},
			wantErr: domain.ErrPasswordNotMatch,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, _, hasher, _, _ := setupTestUsecase()
			tt.setup(repo, hasher)

			err := uc.Register(context.Background(), tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Login
func TestAuthUsecase_Login(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.LoginInput
		setup   func(*mockAuthRepository, *mockPasswordHasher, *mockTokenService)
		wantErr error
	}{
		{
			name: "success - valid credentials",
			input: domain.LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, tokenSvc *mockTokenService) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return &domain.User{
						ID:           "user-123",
						Email:        email,
						PasswordHash: "hashed-password123",
						RoleID:       "user",
					}, nil
				}
				hasher.verifyPassword = func(password, hash string) bool {
					return password == "password123" && hash == "hashed-password123"
				}
			},
			wantErr: nil,
		},
		{
			name: "failure - invalid credentials",
			input: domain.LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, tokenSvc *mockTokenService) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return &domain.User{
						ID:           "user-123",
						Email:        email,
						PasswordHash: "hashed-password123",
					}, nil
				}
				hasher.verifyPassword = func(password, hash string) bool {
					return false
				}
			},
			wantErr: domain.ErrInvalidCredentials,
		},
		{
			name: "failure - user not found",
			input: domain.LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, tokenSvc *mockTokenService) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr: domain.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, tokenSvc, hasher, _, _ := setupTestUsecase()
			tt.setup(repo, hasher, tokenSvc)

			_, err := uc.Login(context.Background(), tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Mock Event Publisher
type mockEventBus struct {
	publishedSubject string
	publishedData    []byte
}

func (m *mockEventBus) Publish(subject string, payload []byte) error {
	m.publishedSubject = subject
	m.publishedData = payload
	return nil
}

// Test ForgotPassword
func TestAuthUsecase_ForgotPassword(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		setup   func(*mockAuthRepository, *mockPasswordHasher, *AuthUsecase)
		wantErr error
	}{
		{
			name:  "success - user exists",
			email: "test@example.com",
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, uc *AuthUsecase) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return &domain.User{
						ID:    "user-123",
						Email: email,
					}, nil
				}
				// Set up event publisher with mock bus
				mockBus := &mockEventBus{}
				uc.eventPub = event.NewAuthPublisher(mockBus)
			},
			wantErr: nil,
		},
		{
			name:  "failure - user not found",
			email: "nonexistent@example.com",
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, uc *AuthUsecase) {
				repo.findUserByEmail = func(email string) (*domain.User, error) {
					return nil, errors.New("user not found")
				}
			},
			wantErr: domain.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, _, hasher, _, _ := setupTestUsecase()
			tt.setup(repo, hasher, uc)

			err := uc.ForgotPassword(context.Background(), tt.email)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ForgotPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Logout
func TestAuthUsecase_Logout(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken string
		setup        func(*mockAuthRepository, *mockPasswordHasher)
		wantErr      error
	}{
		{
			name:         "success - valid token",
			refreshToken: "valid-refresh-token",
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				hasher.hashPassword = func(password string) (string, error) {
					return "hashed-" + password, nil
				}
				repo.findValidRefreshToken = func(tokenHash string) (*domain.RefreshToken, error) {
					if tokenHash == "hashed-valid-refresh-token" {
						return &domain.RefreshToken{
							TokenHash: tokenHash,
							Revoked:   false,
						}, nil
					}
					return nil, errors.New("token not found")
				}
			},
			wantErr: nil,
		},
		{
			name:         "failure - invalid token",
			refreshToken: "invalid-token",
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				hasher.hashPassword = func(password string) (string, error) {
					return "hashed-" + password, nil
				}
				repo.findValidRefreshToken = func(tokenHash string) (*domain.RefreshToken, error) {
					return nil, errors.New("token not found")
				}
			},
			wantErr: domain.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, _, hasher, _, _ := setupTestUsecase()
			tt.setup(repo, hasher)

			err := uc.Logout(context.Background(), tt.refreshToken)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test ResetPassword
func TestAuthUsecase_ResetPassword(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   domain.ResetPasswordInput
		setup   func(*mockAuthRepository, *mockPasswordHasher)
		wantErr error
	}{
		{
			name: "success - valid reset",
			input: domain.ResetPasswordInput{
				Token:       "valid-token",
				NewPassword: "newpassword123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				repo.findValidPasswordReset = func(tokenHash string) (*domain.PasswordReset, error) {
					return &domain.PasswordReset{
						ID:        "reset-123",
						UserID:    "user-123",
						TokenHash: tokenHash,
						Used:      false,
						ExpiresAt: fixedTime.Add(1 * time.Hour),
					}, nil
				}
				repo.findUserByID = func(id string) (*domain.User, error) {
					return &domain.User{
						ID:    id,
						Email: "test@example.com",
					}, nil
				}
			},
			wantErr: nil,
		},
		{
			name: "failure - token expired",
			input: domain.ResetPasswordInput{
				Token:       "expired-token",
				NewPassword: "newpassword123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				repo.findValidPasswordReset = func(tokenHash string) (*domain.PasswordReset, error) {
					return &domain.PasswordReset{
						ID:        "reset-123",
						UserID:    "user-123",
						TokenHash: tokenHash,
						Used:      false,
						ExpiresAt: fixedTime.Add(-1 * time.Hour), // Expired
					}, nil
				}
			},
			wantErr: domain.ErrPasswordResetExpired,
		},
		{
			name: "failure - token already used",
			input: domain.ResetPasswordInput{
				Token:       "used-token",
				NewPassword: "newpassword123",
			},
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher) {
				repo.findValidPasswordReset = func(tokenHash string) (*domain.PasswordReset, error) {
					return &domain.PasswordReset{
						ID:        "reset-123",
						UserID:    "user-123",
						TokenHash: tokenHash,
						Used:      true, // Already used
						ExpiresAt: fixedTime.Add(1 * time.Hour),
					}, nil
				}
			},
			wantErr: domain.ErrPasswordResetUsed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, _, hasher, _, _ := setupTestUsecase()
			tt.setup(repo, hasher)

			err := uc.ResetPassword(context.Background(), tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test RefreshToken
func TestAuthUsecase_RefreshToken(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		setup   func(*mockAuthRepository, *mockPasswordHasher, *mockTokenService)
		wantErr error
	}{
		{
			name:  "success - valid refresh token",
			token: "valid-refresh-token",
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, tokenSvc *mockTokenService) {
				hasher.hashPassword = func(password string) (string, error) {
					return "hashed-" + password, nil
				}
				repo.findValidRefreshToken = func(tokenHash string) (*domain.RefreshToken, error) {
					if tokenHash == "hashed-valid-refresh-token" {
						return &domain.RefreshToken{
							ID:        "token-123",
							UserID:    "user-123",
							TokenHash: tokenHash,
							Revoked:   false,
							ExpiresAt: time.Now().Add(24 * time.Hour),
						}, nil
					}
					return nil, errors.New("token not found")
				}
			},
			wantErr: nil,
		},
		{
			name:  "failure - invalid refresh token",
			token: "invalid-token",
			setup: func(repo *mockAuthRepository, hasher *mockPasswordHasher, tokenSvc *mockTokenService) {
				hasher.hashPassword = func(password string) (string, error) {
					return "hashed-" + password, nil
				}
				repo.findValidRefreshToken = func(tokenHash string) (*domain.RefreshToken, error) {
					return nil, errors.New("token not found")
				}
			},
			wantErr: domain.ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc, repo, tokenSvc, hasher, _, _ := setupTestUsecase()
			tt.setup(repo, hasher, tokenSvc)

			_, err := uc.RefreshToken(context.Background(), tt.token)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Benchmark Register
func BenchmarkAuthUsecase_Register(b *testing.B) {
	uc, repo, _, _, _, _ := setupTestUsecase()
	repo.findUserByEmail = func(email string) (*domain.User, error) {
		return nil, errors.New("not found")
	}

	input := domain.RegisterInput{
		Name:                 "Test User",
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = uc.Register(context.Background(), input)
	}
}

// Benchmark Login
func BenchmarkAuthUsecase_Login(b *testing.B) {
	uc, repo, _, _, _, _ := setupTestUsecase()
	repo.findUserByEmail = func(email string) (*domain.User, error) {
		return &domain.User{
			ID:           "user-123",
			Email:        email,
			PasswordHash: "hashed-password123",
		}, nil
	}

	input := domain.LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = uc.Login(context.Background(), input)
	}
}
