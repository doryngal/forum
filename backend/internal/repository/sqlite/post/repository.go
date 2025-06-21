package post

import (
	"database/sql"
	"forum/internal/domain"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r repository) Create(post *domain.Post) error {
	post.ID = uuid.New()
	post.CreatedAt = time.Now()
	_, err := r.db.Exec(
		"INSERT INTO posts (id, user_id, title, content, created_at) VALUES (?, ?, ?, ?, ?)",
		post.ID.String(), post.UserID.String(), post.Title, post.Content, post.CreatedAt,
	)
	return err
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
		return nil, err
	}
	p.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	p.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	// Load categories for the post
	categories, err := r.getPostCategories(p.ID)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, err
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, err
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
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

func (r repository) GetByUserID(userID uuid.UUID) ([]*domain.Post, error) {
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
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, err
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		var p domain.Post
		var idStr, userIDStr string
		if err := rows.Scan(&idStr, &userIDStr, &p.Title, &p.Content, &p.CreatedAt,
			&p.AuthorUsername, &p.Likes, &p.Dislikes, &p.CommentsCount); err != nil {
			return nil, err
		}
		p.ID, err = uuid.Parse(idStr)
		if err != nil {
			return nil, err
		}
		p.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
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
	_, err := r.db.Exec(`
		INSERT INTO post_reactions (user_id, post_id, reaction) 
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, post_id) DO UPDATE SET reaction = excluded.reaction`,
		userID.String(), postID.String(), reaction)
	return err
}
