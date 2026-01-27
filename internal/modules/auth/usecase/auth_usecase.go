package usecase

import (
	"context"
	"time"

	"github.com/nassabiq/golang-template/internal/modules/auth/domain"
	"github.com/nassabiq/golang-template/internal/modules/auth/event"
	"github.com/nassabiq/golang-template/internal/shared/helper"
)

type TokenService interface {
	GenerateAccessToken(userID, role string) (string, error)
	GenerateRefreshToken() (plain string, hash string, err error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) bool
	GenerateRandomToken() (string, error)
	HashToken(token string) string
}

type AuthUsecase struct {
	repository     domain.AuthRepository
	token          TokenService
	now            func() time.Time
	passwordHasher PasswordHasher
	uuid           helper.UUIDGeneratorInterface
	eventPub       *event.Publisher
}

func NewAuthUsecase(repository domain.AuthRepository, pub *event.Publisher) *AuthUsecase {
	return &AuthUsecase{
		repository: repository,
		eventPub:   pub,
	}
}

func (usecase *AuthUsecase) SetPasswordHasher(hasher PasswordHasher) {
	usecase.passwordHasher = hasher
}

func (usecase *AuthUsecase) SetUUIDGenerator(uuid helper.UUIDGeneratorInterface) {
	usecase.uuid = uuid
}

func (usecase *AuthUsecase) SetTokenService(token TokenService) {
	usecase.token = token
}

func (usecase *AuthUsecase) SetNowFunc(now func() time.Time) {
	usecase.now = now
}

func (usecase *AuthUsecase) Register(ctx context.Context, req domain.RegisterInput) error {

	isExists, _ := usecase.repository.FindUserByEmail(ctx, req.Email)

	if isExists != nil {
		return domain.ErrEmailAlreadyUsed
	}

	if req.Password != req.PasswordConfirmation {
		return domain.ErrPasswordNotMatch
	}

	hash, err := usecase.passwordHasher.HashPassword(req.Password)

	if err != nil {
		return err
	}

	user := &domain.User{
		ID:           usecase.uuid.GenerateID(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
		RoleID:       string(domain.RoleIDUser), // default: user role
		CreatedAt:    usecase.now(),
		UpdatedAt:    usecase.now(),
	}

	return usecase.repository.CreateUser(ctx, user)
}

func (usecase *AuthUsecase) Login(ctx context.Context, req domain.LoginInput) (*domain.AuthOutput, error) {
	user, err := usecase.repository.FindUserByEmail(ctx, req.Email)

	if err != nil || user == nil || !usecase.passwordHasher.VerifyPassword(req.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := usecase.token.GenerateAccessToken(user.ID, user.RoleID)

	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenHash, err := usecase.token.GenerateRefreshToken()

	if err != nil {
		return nil, err
	}

	if err := usecase.repository.StoreRefreshToken(ctx, &domain.RefreshToken{
		ID:        usecase.uuid.GenerateID(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: usecase.now().Add(time.Hour * 24 * 7),
		CreatedAt: usecase.now(),
		UpdatedAt: usecase.now(),
	}); err != nil {
		return nil, err
	}

	return &domain.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (usecase *AuthUsecase) RefreshToken(ctx context.Context, token string) (*domain.AuthOutput, error) {
	hashedToken := usecase.passwordHasher.HashToken(token)
	refreshToken, err := usecase.repository.FindValidRefreshToken(ctx, hashedToken)

	if err != nil || refreshToken == nil {
		return nil, domain.ErrInvalidToken
	}

	usecase.repository.RevokeRefreshToken(ctx, refreshToken.TokenHash)

	accessToken, err := usecase.token.GenerateAccessToken(refreshToken.UserID, "")

	if err != nil {
		return nil, err
	}

	plain, hashNew, err := usecase.token.GenerateRefreshToken()

	if err != nil {
		return nil, err
	}

	usecase.repository.StoreRefreshToken(ctx, &domain.RefreshToken{
		ID:        usecase.uuid.GenerateID(),
		UserID:    refreshToken.UserID,
		TokenHash: hashNew,
		ExpiresAt: usecase.now().Add(time.Hour * 24 * 7),
		CreatedAt: usecase.now(),
		UpdatedAt: usecase.now(),
	})

	return &domain.AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: plain,
	}, nil
}

func (usecase *AuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	user, err := usecase.repository.FindUserByEmail(ctx, email)

	if err != nil || user == nil {
		return domain.ErrUserNotFound
	}

	token, err := usecase.passwordHasher.GenerateRandomToken()

	if err != nil {
		return err
	}

	tokenHash := usecase.passwordHasher.HashToken(token)

	err = usecase.repository.StorePasswordReset(ctx, &domain.PasswordReset{
		ID:        usecase.uuid.GenerateID(),
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: usecase.now().Add(15 * time.Minute),
		CreatedAt: usecase.now(),
		UpdatedAt: usecase.now(),
	})

	if err != nil {
		return err
	}

	_ = usecase.eventPub.ForgotPassword(event.ForgotPasswordEvent{
		Email:     user.Email,
		Token:     token,
		ExpiredAt: usecase.now().Add(15 * time.Minute),
	})

	return nil
}

func (usecase *AuthUsecase) ResetPassword(ctx context.Context, req domain.ResetPasswordInput) error {
	// Hash the incoming token to match against stored hash
	tokenHash := usecase.passwordHasher.HashToken(req.Token)
	passwordReset, err := usecase.repository.FindValidPasswordReset(ctx, tokenHash)

	if err != nil || passwordReset == nil {
		return domain.ErrInvalidToken
	}

	if passwordReset.ExpiresAt.Before(usecase.now()) {
		return domain.ErrPasswordResetExpired
	}

	if passwordReset.Used {
		return domain.ErrPasswordResetUsed
	}

	user, err := usecase.repository.FindUserByID(ctx, passwordReset.UserID)

	if err != nil {
		return domain.ErrUserNotFound
	}

	hashedPassword, err := usecase.passwordHasher.HashPassword(req.NewPassword)

	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	usecase.repository.UpdateUserPassword(ctx, user.ID, hashedPassword)
	usecase.repository.MarkPasswordResetUsed(ctx, user.ID)
	usecase.repository.RevokeRefreshToken(ctx, passwordReset.TokenHash)

	return nil
}

func (usecase *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	hashedToken := usecase.passwordHasher.HashToken(refreshToken)
	token, err := usecase.repository.FindValidRefreshToken(ctx, hashedToken)
	if err != nil || token == nil {
		return domain.ErrInvalidToken
	}

	return usecase.repository.RevokeRefreshToken(ctx, hashedToken)
}
