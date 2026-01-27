package domain

import "time"

type {{MODULE}} struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type {{MODULE}}Create struct {
	ID   string
	Name string
}

type {{MODULE}}Update struct {
	ID   string
	Name *string
}