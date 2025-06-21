package category

import "errors"

var (
	ErrCategoryExists    = errors.New("category: already exists")
	ErrInvalidCategoryID = errors.New("category: invalid category id")
	ErrInvalidPostID     = errors.New("category: invalid post id")
	ErrNoCategories      = errors.New("category: no categories provided")
	ErrNotFound          = errors.New("category: not found")
)
