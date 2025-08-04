package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/domain"
	category2 "forum/internal/repository/category"
	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) category2.Repository {
	return &repository{db: db}
}

func (r repository) Create(name string) error {
	_, err := r.db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", uuid.New().String(), name)
	if err != nil {
		return fmt.Errorf("%w: %v", category2.ErrInsertFailed, err)
	}
	return nil
}

func (r repository) GetAll() ([]*domain.Category, error) {
	rows, err := r.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", category2.ErrQueryFailed, err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		var idStr string
		if err := rows.Scan(&idStr, &c.Name); err != nil {
			return nil, fmt.Errorf("%w: %v", category2.ErrScanFailed, err)
		}
		c.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", category2.ErrUUIDParseFailed, err)
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r repository) AssignToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %v", category2.ErrTransactionBegin, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID.String())
	if err != nil {
		return fmt.Errorf("%w: %v", category2.ErrQueryFailed, err)
	}

	stmt, err := tx.Prepare("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%w: %v", category2.ErrPrepareStmtFailed, err)
	}
	defer stmt.Close()

	for _, categoryID := range categoryIDs {
		_, err = stmt.Exec(postID, categoryID)
		if err != nil {
			return fmt.Errorf("%w: %v", category2.ErrInsertRelation, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", category2.ErrTransactionCommit, err)
	}
	return nil
}

func (r *repository) ExistsByName(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", category2.ErrExistsCheckFailed, err)
	}
	return exists, nil
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", id.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", category2.ErrExistsCheckFailed, err)
	}
	return exists, nil
}
