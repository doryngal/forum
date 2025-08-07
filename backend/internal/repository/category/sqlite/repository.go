package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/domain"
	category_repo "forum/internal/repository/category"
	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) category_repo.Repository {
	return &repository{db: db}
}

func (r *repository) Create(name string) error {
	_, err := r.db.Exec("INSERT INTO categories (id, name) VALUES (?, ?)", uuid.New().String(), name)
	if err != nil {
		return fmt.Errorf("%w: %v", category_repo.ErrInsertFailed, err)
	}
	return nil
}

func (r *repository) GetAll() ([]*domain.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", category_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		var idStr string
		if err := rows.Scan(&idStr, &c.Name); err != nil {
			return nil, fmt.Errorf("%w: %v", category_repo.ErrScanFailed, err)
		}
		c.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", category_repo.ErrUUIDParseFailed, err)
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r *repository) AssignToPost(postID uuid.UUID, categoryIDs []uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %v", category_repo.ErrTransactionBegin, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM post_categories WHERE post_id = ?", postID.String())
	if err != nil {
		return fmt.Errorf("%w: %v", category_repo.ErrQueryFailed, err)
	}

	stmt, err := tx.Prepare("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("%w: %v", category_repo.ErrPrepareStmtFailed, err)
	}
	defer stmt.Close()

	for _, categoryID := range categoryIDs {
		_, err = stmt.Exec(postID.String(), categoryID.String())
		if err != nil {
			return fmt.Errorf("%w: %v", category_repo.ErrInsertRelation, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", category_repo.ErrTransactionCommit, err)
	}
	return nil
}

func (r *repository) ExistsByName(name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", category_repo.ErrExistsCheckFailed, err)
	}
	return exists, nil
}

func (r *repository) ExistsByID(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", id.String()).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", category_repo.ErrExistsCheckFailed, err)
	}
	return exists, nil
}
