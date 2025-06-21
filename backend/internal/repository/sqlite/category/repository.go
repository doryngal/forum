package category

import (
	"database/sql"
	"forum/internal/domain"
	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r repository) Create(name string) error {
	_, err := r.db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)",
		uuid.New().String(), name)
	return err
}

func (r repository) GetAll() ([]*domain.Category, error) {
	rows, err := r.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		var idStr string
		if err := rows.Scan(&idStr, &c.Name); err != nil {
			return nil, err
		}
		c.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r repository) AssignToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID.String())
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, categoryID := range categoryIDs {
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *repository) ExistsByName(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", name).Scan(&exists)
	return exists, err
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", id.String()).Scan(&exists)
	return exists, err
}
