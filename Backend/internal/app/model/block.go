package model

import (
	"encoding/json"
	"time"
)

type BlockType string

const (
	BlockTypeText    BlockType = "text"
	BlockTypeSlider  BlockType = "slider"
	BlockTypeGallery BlockType = "gallery"
	BlockTypeVideo   BlockType = "video"
)

type Block struct {
	ID        string          `json:"id"`
	Page      string          `json:"page" validate:"required"`
	Pages     json.RawMessage `json:"-"`
	Type      BlockType       `json:"type" validate:"required,oneof=text slider gallery video"`
	Content   json.RawMessage `json:"content"`
	Order     int             `json:"order"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	
	// Для JSON сериализации
	PagesArray []string `json:"pages,omitempty"`
}

// MarshalJSON кастомная сериализация для правильной обработки массивов
func (b *Block) MarshalJSON() ([]byte, error) {
	type Alias Block
	aux := &struct {
		*Alias
		DisplayOrder int `json:"display_order"`
	}{
		Alias:        (*Alias)(b),
		DisplayOrder: b.Order,
	}
	
	// Парсим Pages из json.RawMessage
	if len(b.Pages) > 0 {
		if err := json.Unmarshal(b.Pages, &aux.PagesArray); err != nil {
			aux.PagesArray = []string{}
		}
	} else {
		aux.PagesArray = []string{}
	}
	
	return json.Marshal(aux)
}
