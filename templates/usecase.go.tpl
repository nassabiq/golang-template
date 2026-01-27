package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/nassabiq/golang-template/internal/modules/{{MODULE}}/domain"
	"github.com/nassabiq/golang-template/internal/modules/{{MODULE}}/dto"
)

type {{MODULE}}Usecase struct {
	repository domain.{{MODULE}}Repository
}

func New{{MODULE}}Usecase(repository domain.{{MODULE}}Repository) *{{MODULE}}Usecase {
	return &{{MODULE}}Usecase{
		repository: repository,
	}
}

func (usecase *{{MODULE}}Usecase) List(ctx context.Context, limit int, offset int) ([]domain.{{MODULE}}, error) {
	if limit <= 0 {
		limit = 10
	}

	return usecase.repository.List(ctx, limit, offset)
}

func (usecase *{{MODULE}}Usecase) GetByID(ctx context.Context, id string) (*domain.{{MODULE}}, error) {
	return usecase.repository.FindByID(ctx, id)
}

func (usecase *{{MODULE}}Usecase) Create(ctx context.Context, request *dto.Create{{MODULE}}Dto) (*domain.{{MODULE}}, error) {
	params := &domain.{{MODULE}}Create{
		ID:   uuid.New().String(),
		Name: request.Name,
	}

	return usecase.repository.Create(ctx, params)
}

func (usecase *{{MODULE}}Usecase) Update(ctx context.Context, request *dto.Update{{MODULE}}Dto) (*domain.{{MODULE}}, error) {
	params := &domain.{{MODULE}}Update{
		ID:   request.ID,
		Name: request.Name,
	}

	return usecase.repository.Update(ctx, params)
}

func (usecase *{{MODULE}}Usecase) Delete(ctx context.Context, {{MODULE|lower}} *domain.{{MODULE}}) error {
	return usecase.repository.Delete(ctx, {{MODULE|lower}})
}