package dto

type CreateUserDto struct {
	Name     string `validate:"required,min=3,max=100"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=100"`
	RoleID   string `validate:"required"`
}

type UpdateUserDto struct {
	ID     string  `validate:"required"`
	Name   *string `validate:"omitempty,min=3,max=100"`
	Email  *string `validate:"omitempty,email"`
	RoleID *string `validate:"omitempty"`
}
