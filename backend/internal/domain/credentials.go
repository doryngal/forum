package domain

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserStats struct {
	PostCount    int
	LikeCount    int
	DislikeCount int
	CommentCount int
}

type ProfileData struct {
	User  *User
	Stats *UserStats
	Posts []*Post
}
