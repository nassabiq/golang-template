package handler

import (
	"context"
	"log"

	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	authpb "github.com/nassabiq/golang-template/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthUsecaseInterface interface {
	Register(ctx context.Context, req domain.RegisterInput) error
	Login(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error)
	RefreshToken(ctx context.Context, token string) (*domain.AuthOutput, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req domain.ResetPasswordInput) error
}

type AuthHandler struct {
	authUC AuthUsecaseInterface
	authpb.UnimplementedAuthServiceServer
}

func NewAuthHandler(authUC AuthUsecaseInterface) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.MessageResponse, error) {

	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	if err := h.authUC.Register(ctx, domain.RegisterInput{
		Name:                 req.Name,
		Email:                req.Email,
		Password:             req.Password,
		PasswordConfirmation: req.PasswordConfirmation,
	}); err != nil {
		switch err {
		case domain.ErrUserAlreadyExists, domain.ErrEmailAlreadyUsed:
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case domain.ErrPasswordNotMatch:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Printf("[Auth] Register error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authpb.MessageResponse{Message: "Registration successful"}, nil
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.AuthResponse, error) {
	if req.GetEmail() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password required")
	}

	result, err := h.authUC.Login(ctx, domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			log.Printf("[Auth] Login error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authpb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (h *AuthHandler) Refresh(
	ctx context.Context,
	req *authpb.RefreshRequest,
) (*authpb.AuthResponse, error) {

	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token required")
	}

	result, err := h.authUC.RefreshToken(ctx, req.RefreshToken)

	if err != nil {
		switch err {
		case domain.ErrInvalidRefreshToken, domain.ErrTokenExpired:
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			log.Printf("[Auth] Refresh error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authpb.AuthResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (h *AuthHandler) Logout(
	ctx context.Context,
	req *authpb.LogoutRequest,
) (*authpb.MessageResponse, error) {

	if req.GetRefreshToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token required")
	}

	if err := h.authUC.Logout(ctx, req.RefreshToken); err != nil {
		log.Printf("[Auth] Logout error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authpb.MessageResponse{Message: "Logout successful"}, nil
}

func (h *AuthHandler) ForgotPassword(
	ctx context.Context,
	req *authpb.ForgotPasswordRequest,
) (*authpb.MessageResponse, error) {

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if err := h.authUC.ForgotPassword(ctx, req.Email); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			// Return success anyway to prevent email enumeration
			return &authpb.MessageResponse{Message: "If your email is registered, you will receive a password reset link"}, nil
		default:
			log.Printf("[Auth] ForgotPassword error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authpb.MessageResponse{Message: "If your email is registered, you will receive a password reset link"}, nil
}

func (h *AuthHandler) ResetPassword(
	ctx context.Context,
	req *authpb.ResetPasswordRequest,
) (*authpb.MessageResponse, error) {

	if req.GetToken() == "" || req.GetNewPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "token and new_password are required")
	}

	if err := h.authUC.ResetPassword(ctx, domain.ResetPasswordInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}); err != nil {
		switch err {
		case domain.ErrInvalidToken, domain.ErrPasswordResetExpired, domain.ErrPasswordResetUsed:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Printf("[Auth] ResetPassword error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &authpb.MessageResponse{Message: "Password reset successful"}, nil
}
