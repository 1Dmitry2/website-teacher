package model

import (
	"encoding/json"
	"time"
)

type Post struct {
	ID             string          `json:"id"`
	Title          string          `json:"title" validate:"required"`
	Content        string          `json:"content" validate:"required"`
	Images         json.RawMessage `json:"-"`
	Videos         json.RawMessage `json:"-"`
	Pages          json.RawMessage `json:"-"`
	IsPublished   bool            `json:"is_published"`
	Alignment     string          `json:"alignment,omitempty"`
	TitlePosition  string          `json:"title_position,omitempty"`
	ContentPosition string         `json:"content_position,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	
	ImagesArray []string `json:"images,omitempty"`
	VideosArray []string `json:"videos,omitempty"`
	PagesArray  []string `json:"pages,omitempty"`
}

func (p *Post) MarshalJSON() ([]byte, error) {
	type Alias Post
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	
	if len(p.Images) > 0 {
		if err := json.Unmarshal(p.Images, &aux.ImagesArray); err != nil {
			aux.ImagesArray = []string{}
		}
	} else {
		aux.ImagesArray = []string{}
	}
	
	if len(p.Videos) > 0 {
		if err := json.Unmarshal(p.Videos, &aux.VideosArray); err != nil {
			aux.VideosArray = []string{}
		}
	} else {
		aux.VideosArray = []string{}
	}
	
	if len(p.Pages) > 0 {
		if err := json.Unmarshal(p.Pages, &aux.PagesArray); err != nil {
			aux.PagesArray = []string{}
		}
	} else {
		aux.PagesArray = []string{}
	}
	
	return json.Marshal(aux)
}
