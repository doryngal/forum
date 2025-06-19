package domain

import "github.com/google/uuid"

type Post struct {
	ID uuid.UUID `json:"id"`
}

type User struct {
	ID uuid.UUID `json:"id"`
}

type Comment struct {
	ID uuid.UUID `json:"id"`
}

type Category struct {
	ID uuid.UUID `json:"id"`
}
