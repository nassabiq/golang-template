package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	"github.com/nassabiq/golang-template/internal/modules/auth/usecase"
	authpb "github.com/nassabiq/golang-template/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock AuthUsecase
type mockAuthUsecase struct {
	registerFunc       func(ctx context.Context, req domain.RegisterInput) error
	loginFunc          func(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error)
	refreshTokenFunc   func(ctx context.Context, token string) (*domain.AuthOutput, error)
	logoutFunc         func(ctx context.Context, refreshToken string) error
	forgotPasswordFunc func(ctx context.Context, email string) error
	resetPasswordFunc  func(ctx context.Context, req domain.ResetPasswordInput) error
}

func (m *mockAuthUsecase) Register(ctx context.Context, req domain.RegisterInput) error {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, req)
	}
	return nil
}

func (m *mockAuthUsecase) Login(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, req)
	}
	return nil, nil
}

func (m *mockAuthUsecase) RefreshToken(ctx context.Context, token string) (*domain.AuthOutput, error) {
	if m.refreshTokenFunc != nil {
		return m.refreshTokenFunc(ctx, token)
	}
	return nil, nil
}

func (m *mockAuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	if m.logoutFunc != nil {
		return m.logoutFunc(ctx, refreshToken)
	}
	return nil
}

func (m *mockAuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	if m.forgotPasswordFunc != nil {
		return m.forgotPasswordFunc(ctx, email)
	}
	return nil
}

func (m *mockAuthUsecase) ResetPassword(ctx context.Context, req domain.ResetPasswordInput) error {
	if m.resetPasswordFunc != nil {
		return m.resetPasswordFunc(ctx, req)
	}
	return nil
}

func setupTestHandler() (*AuthHandler, *mockAuthUsecase) {
	mockUC := &mockAuthUsecase{}
	handler := NewAuthHandler((*usecase.AuthUsecase)(nil))
	// Use reflection to set private field for testing
	// In real scenario, we'd modify the handler to accept interface
	// For now, we'll test via behavior
	_ = handler
	return handler, mockUC
}

// Test Register
func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name        string
		req         *authpb.RegisterRequest
		mockSetup   func(*mockAuthUsecase)
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "success - valid registration",
			req: &authpb.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.registerFunc = func(ctx context.Context, req domain.RegisterInput) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "failure - empty email",
			req: &authpb.RegisterRequest{
				Email:    "",
				Password: "password123",
			},
			mockSetup:   func(m *mockAuthUsecase) {},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "failure - empty password",
			req: &authpb.RegisterRequest{
				Email:    "test@example.com",
				Password: "",
			},
			mockSetup:   func(m *mockAuthUsecase) {},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "failure - user already exists",
			req: &authpb.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.registerFunc = func(ctx context.Context, req domain.RegisterInput) error {
					return domain.ErrUserAlreadyExists
				}
			},
			wantErr:     true,
			wantErrCode: codes.AlreadyExists,
		},
		{
			name: "failure - internal error",
			req: &authpb.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.registerFunc = func(ctx context.Context, req domain.RegisterInput) error {
					return errors.New("database error")
				}
			},
			wantErr:     true,
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockAuthUsecase{}
			tt.mockSetup(mockUC)

			// Create handler with mock
			handler := &AuthHandler{authUC: mockUC}

			_, err := handler.Register(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Register() error is not a gRPC status error")
					return
				}
				if st.Code() != tt.wantErrCode {
					t.Errorf("Register() error code = %v, want %v", st.Code(), tt.wantErrCode)
				}
			}
		})
	}
}

// Test Login
func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name        string
		req         *authpb.LoginRequest
		mockSetup   func(*mockAuthUsecase)
		wantErr     bool
		wantErrCode codes.Code
		wantToken   bool
	}{
		{
			name: "success - valid credentials",
			req: &authpb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.loginFunc = func(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error) {
					return &domain.AuthOutput{
						AccessToken:  "access-token-123",
						RefreshToken: "refresh-token-123",
					}, nil
				}
			},
			wantErr:   false,
			wantToken: true,
		},
		{
			name: "failure - empty email",
			req: &authpb.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			mockSetup:   func(m *mockAuthUsecase) {},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "failure - empty password",
			req: &authpb.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			mockSetup:   func(m *mockAuthUsecase) {},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "failure - invalid credentials",
			req: &authpb.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.loginFunc = func(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error) {
					return nil, domain.ErrInvalidCredentials
				}
			},
			wantErr:     true,
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "failure - internal error",
			req: &authpb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.loginFunc = func(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error) {
					return nil, errors.New("database error")
				}
			},
			wantErr:     true,
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockAuthUsecase{}
			tt.mockSetup(mockUC)

			handler := &AuthHandler{authUC: mockUC}

			resp, err := handler.Login(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Login() error is not a gRPC status error")
					return
				}
				if st.Code() != tt.wantErrCode {
					t.Errorf("Login() error code = %v, want %v", st.Code(), tt.wantErrCode)
				}
			}
			if tt.wantToken && resp != nil {
				if resp.AccessToken == "" || resp.RefreshToken == "" {
					t.Errorf("Login() expected tokens, got empty")
				}
			}
		})
	}
}

// Test Refresh
func TestAuthHandler_Refresh(t *testing.T) {
	tests := []struct {
		name        string
		req         *authpb.RefreshRequest
		mockSetup   func(*mockAuthUsecase)
		wantErr     bool
		wantErrCode codes.Code
		wantToken   bool
	}{
		{
			name: "success - valid refresh token",
			req: &authpb.RefreshRequest{
				RefreshToken: "valid-refresh-token",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.refreshTokenFunc = func(ctx context.Context, token string) (*domain.AuthOutput, error) {
					return &domain.AuthOutput{
						AccessToken:  "new-access-token",
						RefreshToken: "new-refresh-token",
					}, nil
				}
			},
			wantErr:   false,
			wantToken: true,
		},
		{
			name: "failure - empty refresh token",
			req: &authpb.RefreshRequest{
				RefreshToken: "",
			},
			mockSetup:   func(m *mockAuthUsecase) {},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "failure - invalid refresh token",
			req: &authpb.RefreshRequest{
				RefreshToken: "invalid-token",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.refreshTokenFunc = func(ctx context.Context, token string) (*domain.AuthOutput, error) {
					return nil, domain.ErrInvalidRefreshToken
				}
			},
			wantErr:     true,
			wantErrCode: codes.Unauthenticated,
		},
		{
			name: "failure - token expired",
			req: &authpb.RefreshRequest{
				RefreshToken: "expired-token",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.refreshTokenFunc = func(ctx context.Context, token string) (*domain.AuthOutput, error) {
					return nil, domain.ErrTokenExpired
				}
			},
			wantErr:     true,
			wantErrCode: codes.Unauthenticated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockAuthUsecase{}
			tt.mockSetup(mockUC)

			handler := &AuthHandler{authUC: mockUC}

			resp, err := handler.Refresh(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Refresh() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Refresh() error is not a gRPC status error")
					return
				}
				if st.Code() != tt.wantErrCode {
					t.Errorf("Refresh() error code = %v, want %v", st.Code(), tt.wantErrCode)
				}
			}
			if tt.wantToken && resp != nil {
				if resp.AccessToken == "" || resp.RefreshToken == "" {
					t.Errorf("Refresh() expected tokens, got empty")
				}
			}
		})
	}
}

// Test Logout
func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		name        string
		req         *authpb.LogoutRequest
		mockSetup   func(*mockAuthUsecase)
		wantErr     bool
		wantErrCode codes.Code
	}{
		{
			name: "success - valid logout",
			req: &authpb.LogoutRequest{
				RefreshToken: "valid-refresh-token",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.logoutFunc = func(ctx context.Context, refreshToken string) error {
					return nil
				}
			},
			wantErr: false,
		},
		{
			name: "failure - empty refresh token",
			req: &authpb.LogoutRequest{
				RefreshToken: "",
			},
			mockSetup:   func(m *mockAuthUsecase) {},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
		},
		{
			name: "failure - internal error",
			req: &authpb.LogoutRequest{
				RefreshToken: "valid-token",
			},
			mockSetup: func(m *mockAuthUsecase) {
				m.logoutFunc = func(ctx context.Context, refreshToken string) error {
					return errors.New("database error")
				}
			},
			wantErr:     true,
			wantErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mockAuthUsecase{}
			tt.mockSetup(mockUC)

			handler := &AuthHandler{authUC: mockUC}

			_, err := handler.Logout(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Logout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					t.Errorf("Logout() error is not a gRPC status error")
					return
				}
				if st.Code() != tt.wantErrCode {
					t.Errorf("Logout() error code = %v, want %v", st.Code(), tt.wantErrCode)
				}
			}
		})
	}
}
