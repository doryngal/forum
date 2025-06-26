package home

import (
	"fmt"
	"forum/internal/domain"
	"forum/internal/service/category"
	"forum/internal/service/post"
	"forum/internal/service/user"
	"html/template"
	"net/http"
)

type HomeHandler struct {
	postService     post.Service
	userService     user.Service
	categoryService category.Service
}

func NewHomeHandler(ps post.Service, us user.Service, cs category.Service) *HomeHandler {
	return &HomeHandler{
		postService:     ps,
		userService:     us,
		categoryService: cs,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get all posts
	posts, err := h.postService.GetAllPosts()
	if err != nil {
		http.Error(w, "Error fetching posts", http.StatusInternalServerError)
		return
	}

	// Get top categories (trending topics)
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		http.Error(w, "Error fetching categories", http.StatusInternalServerError)
		return
	}

	// Get top users (you might need to implement this in user service)
	// topUsers, err := h.userService.GetTopUsers(5)

	// Create template data
	data := struct {
		Posts      []*domain.Post
		Categories []*domain.Category
		// TopUsers   []*domain.User
	}{
		Posts:      posts,
		Categories: categories,
		// TopUsers:   topUsers,
	}

	// Render template
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}
}
