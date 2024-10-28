package model

import "time"

type Posts struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	PostText  string    `json:"post_text"`
	Author    string    `json:"author"`
	AuthorID  int       `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostsList []Posts

type User struct {
	ID       int    `json:"id"`
	Username string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
