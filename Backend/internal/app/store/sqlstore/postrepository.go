package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"database/sql"
)

type PostRepository struct {
	store *Store
}

func (r *PostRepository) List(filter store.PostFilter) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, images, videos, pages, is_published, alignment, title_position, content_position, created_at, updated_at
		FROM posts
	`
	args := []interface{}{}
	if !filter.IncludeDrafts {
		query += " WHERE is_published = true"
	}
	query += " ORDER BY created_at DESC"

	rows, err := r.store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		var p model.Post
		var images []byte
		var videos []byte
		var pages []byte
		var alignment sql.NullString
		var titlePosition sql.NullString
		var contentPosition sql.NullString
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &images, &videos, &pages, &p.IsPublished, &alignment, &titlePosition, &contentPosition, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Images = images
		p.Videos = videos
		p.Pages = pages
		if alignment.Valid {
			p.Alignment = alignment.String
		}
		if titlePosition.Valid {
			p.TitlePosition = titlePosition.String
		}
		if contentPosition.Valid {
			p.ContentPosition = contentPosition.String
		}
		posts = append(posts, &p)
	}
	return posts, rows.Err()
}

func (r *PostRepository) FindByID(id string) (*model.Post, error) {
	query := `
		SELECT id, title, content, images, videos, pages, is_published, alignment, title_position, content_position, created_at, updated_at
		FROM posts WHERE id = $1
	`
	var p model.Post
	var images []byte
	var videos []byte
	var pages []byte
	var alignment sql.NullString
	var titlePosition sql.NullString
	var contentPosition sql.NullString
	if err := r.store.db.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &images, &videos, &pages, &p.IsPublished, &alignment, &titlePosition, &contentPosition, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	p.Images = images
	p.Videos = videos
	p.Pages = pages
	if alignment.Valid {
		p.Alignment = alignment.String
	}
	if titlePosition.Valid {
		p.TitlePosition = titlePosition.String
	}
	if contentPosition.Valid {
		p.ContentPosition = contentPosition.String
	}
	return &p, nil
}

func (r *PostRepository) Create(p *model.Post) error {
	query := `
		INSERT INTO posts (title, content, images, videos, pages, is_published, alignment, title_position, content_position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	images := []byte("[]")
	if len(p.Images) > 0 {
		images = p.Images
	}
	videos := []byte("[]")
	if len(p.Videos) > 0 {
		videos = p.Videos
	}
	pages := []byte("[]")
	if len(p.Pages) > 0 {
		pages = p.Pages
	}
	alignment := p.Alignment
	if alignment == "" {
		alignment = "full-width"
	}
	titlePosition := p.TitlePosition
	if titlePosition == "" {
		titlePosition = "top"
	}
	contentPosition := p.ContentPosition
	if contentPosition == "" {
		contentPosition = "bottom"
	}
	return r.store.db.QueryRow(query, p.Title, p.Content, images, videos, pages, p.IsPublished, alignment, titlePosition, contentPosition).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PostRepository) Update(p *model.Post) error {
	query := `
		UPDATE posts
		SET title = $1,
		    content = $2,
		    images = $3,
		    videos = $4,
		    pages = $5,
		    is_published = $6,
		    alignment = $7,
		    title_position = $8,
		    content_position = $9,
		    updated_at = now()
		WHERE id = $10
		RETURNING created_at, updated_at
	`
	images := []byte("[]")
	if len(p.Images) > 0 {
		images = p.Images
	}
	videos := []byte("[]")
	if len(p.Videos) > 0 {
		videos = p.Videos
	}
	pages := []byte("[]")
	if len(p.Pages) > 0 {
		pages = p.Pages
	}
	alignment := p.Alignment
	if alignment == "" {
		alignment = "full-width"
	}
	titlePosition := p.TitlePosition
	if titlePosition == "" {
		titlePosition = "top"
	}
	contentPosition := p.ContentPosition
	if contentPosition == "" {
		contentPosition = "bottom"
	}
	return r.store.db.QueryRow(query, p.Title, p.Content, images, videos, pages, p.IsPublished, alignment, titlePosition, contentPosition, p.ID).
		Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *PostRepository) Delete(id string) error {
	_, err := r.store.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	return err
}
