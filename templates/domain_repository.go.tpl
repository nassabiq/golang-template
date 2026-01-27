package domain

import "context"

type {{MODULE}}Repository interface {
	List(ctx context.Context, limit int, offset int) ([]{{MODULE}}, error)
	FindByID(ctx context.Context, id string) (*{{MODULE}}, error)
	Create(ctx context.Context, request *{{MODULE}}Create) (*{{MODULE}}, error)
	Update(ctx context.Context, request *{{MODULE}}Update) (*{{MODULE}}, error)
	Delete(ctx context.Context, {{MODULE|lower}} *{{MODULE}}) error
}