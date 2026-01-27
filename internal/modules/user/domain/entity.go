package domain

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	RoleID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreate struct {
	ID       string
	Name     string
	Email    string
	Password string
	RoleID   string
}

type UserUpdate struct {
	ID     string
	Name   *string
	Email  *string
	RoleID *string
}
