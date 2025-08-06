package sqlite

import (
	"database/sql"
	"fmt"
	"forum/internal/domain"
	post_repo "forum/internal/repository/post"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) post_repo.Repository {
	return &repository{db: db}
}

func (r repository) Create(post *domain.Post) error {
	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	_, err := r.db.Exec(
		"INSERT INTO posts (id, user_id, title, content, created_at) VALUES (?, ?, ?, ?, ?)",
		post.ID.String(), post.UserID.String(), post.Title, post.Content, post.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", post_repo.ErrInsertPostFailed, err)
	}
	return nil
}

func (r repository) GetByID(id uuid.UUID) (*domain.Post, error) {
	var p domain.Post
	var idStr, userIDStr string
	err := r.db.QueryRow(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
		       u.username, 
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.id = ?
		GROUP BY p.id`, id.String(),
	).Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
		&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	p.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
	}
	p.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
	}

	// Load categories for the post
	categories, err := r.getPostCategories(p.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrLoadCategoriesFailed, err)
	}
	p.Categories = categories

	return &p, nil
}

func (r repository) getPostCategories(postID uuid.UUID) ([]*domain.Category, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.name 
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?`, postID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var c domain.Category
		var idStr string
		if err := rows.Scan(&idStr, &c.Name); err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrScanFailed, err)
		}
		c.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		categories = append(categories, &c)
	}
	return categories, nil
}

func (r repository) GetAll() ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
		       u.username, 
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		GROUP BY p.id
		ORDER BY p.created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrScanFailed, err)
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		posts = append(posts, &p)
	}

	// Load categories for each post
	for _, post := range posts {
		categories, err := r.getPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
	}

	return posts, nil
}

func (r repository) Update(post *domain.Post) error {
	_, err := r.db.Exec(`
		UPDATE posts SET title = $1, content = $2 WHERE id = $3 AND user_id = $4`,
		post.Title, post.Content, post.ID, post.UserID)
	if err != nil {
		return fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	return nil
}

func (r repository) Delete(postID uuid.UUID, userID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	_, err = tx.Exec(`
		DELETE FROM comment_reactions 
		WHERE comment_id IN (SELECT id FROM comments WHERE post_id = $1)
	`, postID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete comment_reactions: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM comments WHERE post_id = $1`, postID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete comments: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM post_reactions WHERE post_id = $1`, postID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete post_reactions: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM post_categories WHERE post_id = $1`, postID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete post_categories: %w", err)
	}

	res, err := tx.Exec(`DELETE FROM posts WHERE id = $1 AND user_id = $2`, postID, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("no post deleted, might not exist or not owned by user")
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r repository) GetByCategory(categoryID uuid.UUID) ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
		       u.username, 
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		JOIN post_categories pc ON p.id = pc.post_id
		WHERE pc.category_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC`, categoryID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrScanFailed, err)
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		posts = append(posts, &p)
	}

	// Load categories for each post
	for _, post := range posts {
		categories, err := r.getPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
	}

	return posts, nil
}

func (r repository) GetByUserID(userID, sessionID uuid.UUID) ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
		       u.username, 
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		WHERE p.user_id = ?
		GROUP BY p.id
		ORDER BY p.created_at DESC`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrScanFailed, err)
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		posts = append(posts, &p)
	}

	// Load categories for each post
	for _, post := range posts {
		categories, err := r.getPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
		post.IsOwner = post.UserID == sessionID
	}

	return posts, nil
}

func (r repository) GetLikedByUser(userID uuid.UUID) ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at, 
		       u.username, 
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		JOIN post_reactions user_pr ON p.id = user_pr.post_id AND user_pr.user_id = ? AND user_pr.reaction = 1
		GROUP BY p.id
		ORDER BY p.created_at DESC`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %v", post_repo.ErrQueryFailed, err)
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrScanFailed, err)
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", post_repo.ErrUUIDParseFailed, err)
		}
		posts = append(posts, &p)
	}

	// Load categories for each post
	for _, post := range posts {
		categories, err := r.getPostCategories(post.ID)
		if err != nil {
			return nil, err
		}
		post.Categories = categories
	}

	return posts, nil
}

func (r repository) Like(postID, userID uuid.UUID) error {
	return r.setReaction(postID, userID, 1)
}

func (r repository) Dislike(postID, userID uuid.UUID) error {
	return r.setReaction(postID, userID, -1)
}

func (r repository) setReaction(postID, userID uuid.UUID, reaction int) error {
	var existingReaction int
	err := r.db.QueryRow(`
		SELECT reaction FROM post_reactions 
		WHERE user_id = ? AND post_id = ?`,
		userID.String(), postID.String(),
	).Scan(&existingReaction)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("%w: %v", post_repo.ErrReactionUpdateFailed, err)
	}

	if err == nil {
		// реакция уже есть
		if existingReaction == reaction {
			// если такая же — удалить
			_, err := r.db.Exec(`
				DELETE FROM post_reactions 
				WHERE user_id = ? AND post_id = ?`,
				userID.String(), postID.String())
			if err != nil {
				return fmt.Errorf("%w: %v", post_repo.ErrReactionUpdateFailed, err)
			}
			return nil
		}
		// если другая — обновить
		_, err := r.db.Exec(`
			UPDATE post_reactions 
			SET reaction = ? 
			WHERE user_id = ? AND post_id = ?`,
			reaction, userID.String(), postID.String())
		if err != nil {
			return fmt.Errorf("%w: %v", post_repo.ErrReactionUpdateFailed, err)
		}
		return nil
	}

	// реакции нет — вставляем новую
	_, err = r.db.Exec(`
		INSERT INTO post_reactions (user_id, post_id, reaction) 
		VALUES (?, ?, ?)`,
		userID.String(), postID.String(), reaction)
	if err != nil {
		return fmt.Errorf("%w: %v", post_repo.ErrReactionUpdateFailed, err)
	}
	return nil
}

func (r repository) ExistsByID(postID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)",
		postID.String()).Scan(&exists)
	return exists, err
}

func (r repository) GetReaction(postID, userID uuid.UUID) (int, error) {
	var reaction int
	err := r.db.QueryRow(
		"SELECT reaction FROM post_reactions WHERE post_id = ? AND user_id = ?",
		postID.String(), userID.String(),
	).Scan(&reaction)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, post_repo.ErrReactionNotFound
		}
		return 0, fmt.Errorf("%w: %v", post_repo.ErrGetReactionFailed, err)
	}

	return reaction, nil
}

func (r repository) GetLikedPostsByUserID(userID uuid.UUID) ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at,
		       u.username,
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN post_reactions r2 ON p.id = r2.post_id AND r2.user_id = ? AND r2.reaction = 1
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		GROUP BY p.id
		ORDER BY p.created_at DESC`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("GetLikedPostsByUserID query failed: %v", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r repository) GetDislikedPostsByUserID(userID uuid.UUID) ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.created_at,
		       u.username,
		       COALESCE(SUM(CASE WHEN pr.reaction = 1 THEN 1 ELSE 0 END), 0) as likes,
		       COALESCE(SUM(CASE WHEN pr.reaction = -1 THEN 1 ELSE 0 END), 0) as dislikes,
		       (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) as comments_count
		FROM posts p
		JOIN post_reactions r2 ON p.id = r2.post_id AND r2.user_id = ? AND r2.reaction = -1
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_reactions pr ON p.id = pr.post_id
		GROUP BY p.id
		ORDER BY p.created_at DESC`, userID.String())
	if err != nil {
		return nil, fmt.Errorf("GetDislikedPostsByUserID query failed: %v", err)
	}
	defer rows.Close()

	return r.scanPosts(rows)
}

func (r repository) scanPosts(rows *sql.Rows) ([]*domain.Post, error) {
	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, fmt.Errorf("scanPosts: %v", err)
		}
		var err error
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		categories, err := r.getPostCategories(p.ID)
		if err != nil {
			return nil, err
		}
		p.Categories = categories
		posts = append(posts, &p)
	}
	return posts, nil
}
