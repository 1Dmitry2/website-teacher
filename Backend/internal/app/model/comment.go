package model

import "time"

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id" validate:"required"`
	UserID    int       `json:"user_id" validate:"required"`
	UserEmail *string   `json:"user_email,omitempty"`
	Text      string    `json:"text" validate:"required"`
	ReplyTo   *string   `json:"reply_to,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}
