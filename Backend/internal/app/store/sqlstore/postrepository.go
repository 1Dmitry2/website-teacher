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
		SELECT id, title, content, images, pages, is_published, created_at, updated_at
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
		var pages []byte
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &images, &pages, &p.IsPublished, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.Images = images
		p.Pages = pages
		posts = append(posts, &p)
	}
	return posts, rows.Err()
}

func (r *PostRepository) FindByID(id string) (*model.Post, error) {
	query := `
		SELECT id, title, content, images, pages, is_published, created_at, updated_at
		FROM posts WHERE id = $1
	`
	var p model.Post
	var images []byte
	var pages []byte
	if err := r.store.db.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &images, &pages, &p.IsPublished, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	p.Images = images
	p.Pages = pages
	return &p, nil
}

func (r *PostRepository) Create(p *model.Post) error {
	query := `
		INSERT INTO posts (title, content, images, pages, is_published)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	images := []byte("[]")
	if len(p.Images) > 0 {
		images = p.Images
	}
	pages := []byte("[]")
	if len(p.Pages) > 0 {
		pages = p.Pages
	}
	return r.store.db.QueryRow(query, p.Title, p.Content, images, pages, p.IsPublished).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PostRepository) Update(p *model.Post) error {
	query := `
		UPDATE posts
		SET title = $1,
		    content = $2,
		    images = $3,
		    pages = $4,
		    is_published = $5,
		    updated_at = now()
		WHERE id = $6
		RETURNING created_at, updated_at
	`
	images := []byte("[]")
	if len(p.Images) > 0 {
		images = p.Images
	}
	pages := []byte("[]")
	if len(p.Pages) > 0 {
		pages = p.Pages
	}
	return r.store.db.QueryRow(query, p.Title, p.Content, images, pages, p.IsPublished, p.ID).
		Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *PostRepository) Delete(id string) error {
	_, err := r.store.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	return err
}
