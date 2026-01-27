package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nassabiq/golang-template/internal/modules/user/domain"
	"github.com/nassabiq/golang-template/internal/modules/user/dto"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
}

type UserUsecase struct {
	repository domain.UserRepository
	hasher     PasswordHasher
}

func NewUserUsecase(repository domain.UserRepository, hasher PasswordHasher) *UserUsecase {
	return &UserUsecase{
		repository: repository,
		hasher:     hasher,
	}
}

func (usecase *UserUsecase) List(ctx context.Context, limit int, offset int) ([]domain.User, int64, error) {
	if limit <= 0 {
		limit = 10
	}

	return usecase.repository.List(ctx, limit, offset)
}

func (usecase *UserUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return usecase.repository.FindByID(ctx, id)
}

func (usecase *UserUsecase) Create(ctx context.Context, request *dto.CreateUserDto) (*domain.User, error) {

	hash, err := usecase.hasher.HashPassword(request.Password)

	if err != nil {
		return nil, err
	}

	params := &domain.UserCreate{
		ID:       uuid.New().String(),
		Name:     request.Name,
		Email:    request.Email,
		Password: hash,
		RoleID:   request.RoleID,
	}

	return usecase.repository.Create(ctx, params)
}

func (usecase *UserUsecase) Update(ctx context.Context, request *dto.UpdateUserDto) (*domain.User, error) {

	params := &domain.UserUpdate{
		ID:     request.ID,
		Name:   request.Name,
		Email:  request.Email,
		RoleID: request.RoleID,
	}

	return usecase.repository.Update(ctx, params)
}

func (usecase *UserUsecase) Delete(ctx context.Context, user *domain.User) error {
	return usecase.repository.Delete(ctx, user)
}
