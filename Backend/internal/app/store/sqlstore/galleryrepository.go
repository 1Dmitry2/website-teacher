package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"database/sql"
)

type GalleryRepository struct {
	store *Store
}

func (r *GalleryRepository) List() ([]*model.GalleryItem, error) {
	rows, err := r.store.db.Query(`
		SELECT id, image_url, title, description, text, pages, created_at, updated_at
		FROM gallery_items
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.GalleryItem
	for rows.Next() {
		var item model.GalleryItem
		var pages []byte
		if err := rows.Scan(&item.ID, &item.ImageURL, &item.Title, &item.Description, &item.Text, &pages, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.Pages = pages
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *GalleryRepository) FindByID(id string) (*model.GalleryItem, error) {
	var item model.GalleryItem
	var pages []byte
	err := r.store.db.QueryRow(`
		SELECT id, image_url, title, description, text, pages, created_at, updated_at
		FROM gallery_items WHERE id = $1
	`, id).Scan(&item.ID, &item.ImageURL, &item.Title, &item.Description, &item.Text, &pages, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	item.Pages = pages
	return &item, nil
}

func (r *GalleryRepository) Create(item *model.GalleryItem) error {
	pages := []byte("[]")
	if len(item.Pages) > 0 {
		pages = item.Pages
	}
	return r.store.db.QueryRow(`
		INSERT INTO gallery_items (image_url, title, description, text, pages)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, item.ImageURL, item.Title, item.Description, item.Text, pages).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *GalleryRepository) Update(item *model.GalleryItem) error {
	pages := []byte("[]")
	if len(item.Pages) > 0 {
		pages = item.Pages
	}
	return r.store.db.QueryRow(`
		UPDATE gallery_items
		SET image_url = $1,
		    title = $2,
		    description = $3,
		    text = $4,
		    pages = $5,
		    updated_at = now()
		WHERE id = $6
		RETURNING created_at, updated_at
	`, item.ImageURL, item.Title, item.Description, item.Text, pages, item.ID).Scan(&item.CreatedAt, &item.UpdatedAt)
}

func (r *GalleryRepository) Delete(id string) error {
	_, err := r.store.db.Exec(`DELETE FROM gallery_items WHERE id = $1`, id)
	return err
}
