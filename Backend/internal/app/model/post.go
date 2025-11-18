package model

import (
	"encoding/json"
	"time"
)

type Post struct {
	ID          string          `json:"id"`
	Title       string          `json:"title" validate:"required"`
	Content     string          `json:"content" validate:"required"`
	Images      json.RawMessage `json:"-"`
	Pages       json.RawMessage `json:"-"`
	IsPublished bool            `json:"is_published"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	
	// Для JSON сериализации
	ImagesArray []string `json:"images,omitempty"`
	PagesArray  []string `json:"pages,omitempty"`
}

// MarshalJSON кастомная сериализация для правильной обработки массивов
func (p *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	
	// Парсим Images из json.RawMessage
	if len(p.Images) > 0 {
		if err := json.Unmarshal(p.Images, &aux.ImagesArray); err != nil {
			aux.ImagesArray = []string{}
		}
	} else {
		aux.ImagesArray = []string{}
	}
	
	// Парсим Pages из json.RawMessage
	if len(p.Pages) > 0 {
		if err := json.Unmarshal(p.Pages, &aux.PagesArray); err != nil {
			aux.PagesArray = []string{}
		}
	} else {
		aux.PagesArray = []string{}
	}
	
	return json.Marshal(aux)
}
