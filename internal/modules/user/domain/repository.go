package domain

import "context"

type UserRepository interface {
	List(ctx context.Context, limit int, offset int) ([]User, int64, error)
	FindByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, request *UserCreate) (*User, error)
	Update(ctx context.Context, request *UserUpdate) (*User, error)
	Delete(ctx context.Context, user *User) error
}
