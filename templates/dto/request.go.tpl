package dto

type Create{{MODULE}}Dto struct {
	Name string `validate:"required,min=3,max=100"`
}

type Update{{MODULE}}Dto struct {
	ID   string  `validate:"required"`
	Name *string `validate:"omitempty,min=3,max=100"`
}
