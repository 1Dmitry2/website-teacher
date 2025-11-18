package model

import (
	"encoding/json"
	"time"
)

type GalleryItem struct {
	ID          string          `json:"id"`
	ImageURL    string          `json:"image_url" validate:"required"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Text        string          `json:"text,omitempty"`
	Pages       json.RawMessage `json:"pages,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type SliderItem struct {
	ID          string          `json:"id"`
	ImageURL    string          `json:"image_url" validate:"required"`
	Title       string          `json:"title,omitempty"`
	Description string          `json:"description,omitempty"`
	Pages       json.RawMessage `json:"pages,omitempty"`
	Order       int             `json:"order"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
