package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"context"
	"database/sql"
	"fmt"
)

type SliderRepository struct {
	store *Store
}

func (r *SliderRepository) List() ([]*model.SliderItem, error) {
	rows, err := r.store.db.Query(`
		SELECT id, image_url, title, description, pages, display_order, created_at, updated_at
		FROM slider_items
		ORDER BY display_order ASC, created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.SliderItem
	for rows.Next() {
		var item model.SliderItem
		var pages []byte
		if err := rows.Scan(&item.ID, &item.ImageURL, &item.Title, &item.Description, &pages, &item.Order, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.Pages = pages
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *SliderRepository) FindByID(id string) (*model.SliderItem, error) {
	var item model.SliderItem
	var pages []byte
	err := r.store.db.QueryRow(`
		SELECT id, image_url, title, description, pages, display_order, created_at, updated_at
		FROM slider_items WHERE id = $1
	`, id).Scan(&item.ID, &item.ImageURL, &item.Title, &item.Description, &pages, &item.Order, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	item.Pages = pages
	return &item, nil
}

func (r *SliderRepository) Create(item *model.SliderItem) error {
	pages := []byte("[]")
	if len(item.Pages) > 0 {
		pages = item.Pages
	}
	return r.store.db.QueryRow(`
		INSERT INTO slider_items (image_url, title, description, pages, display_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`, item.ImageURL, item.Title, item.Description, pages, item.Order).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *SliderRepository) Update(item *model.SliderItem) error {
	pages := []byte("[]")
	if len(item.Pages) > 0 {
		pages = item.Pages
	}
	return r.store.db.QueryRow(`
		UPDATE slider_items
		SET image_url = $1,
		    title = $2,
		    description = $3,
		    pages = $4,
		    display_order = $5,
		    updated_at = now()
		WHERE id = $6
		RETURNING created_at, updated_at
	`, item.ImageURL, item.Title, item.Description, pages, item.Order, item.ID).Scan(&item.CreatedAt, &item.UpdatedAt)
}

func (r *SliderRepository) Delete(id string) error {
	_, err := r.store.db.Exec(`DELETE FROM slider_items WHERE id = $1`, id)
	return err
}

func (r *SliderRepository) UpdateOrders(updates []store.SliderOrderUpdate) error {
	tx, err := r.store.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE slider_items SET display_order = $1, updated_at = now() WHERE id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, upd := range updates {
		if _, err := stmt.Exec(upd.Order, upd.ID); err != nil {
			return fmt.Errorf("update slider order for %s: %w", upd.ID, err)
		}
	}

	return tx.Commit()
}
