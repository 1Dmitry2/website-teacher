package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"context"
	"database/sql"
	"fmt"
)

type BlockRepository struct {
	store *Store
}

func (r *BlockRepository) List(filter store.BlockFilter) ([]*model.Block, error) {
	query := `
		SELECT id, page, pages, type, content, display_order, created_at, updated_at
		FROM blocks
	`
	args := []interface{}{}
	argIndex := 1
	
	if filter.Page != "" {
		query += ` WHERE (
			EXISTS (
				SELECT 1 
				FROM jsonb_array_elements_text(COALESCE(pages::jsonb, '[]'::jsonb)) AS page_value
				WHERE page_value = $` + fmt.Sprintf("%d", argIndex) + `
			) OR page = $` + fmt.Sprintf("%d", argIndex+1) + `
		)`
		args = append(args, filter.Page, filter.Page)
		argIndex += 2
	}
	query += " ORDER BY display_order ASC, created_at ASC"

	rows, err := r.store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Block
	for rows.Next() {
		var b model.Block
		var content []byte
		var pages []byte
		if err := rows.Scan(&b.ID, &b.Page, &pages, &b.Type, &content, &b.Order, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.Content = content
		b.Pages = pages
		result = append(result, &b)
	}
	return result, rows.Err()
}

func (r *BlockRepository) Create(b *model.Block) error {
	query := `
		INSERT INTO blocks (page, pages, type, content, display_order)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	content := []byte("{}")
	if len(b.Content) > 0 {
		content = b.Content
	}
	pages := []byte("[]")
	if len(b.Pages) > 0 {
		pages = b.Pages
	}
	return r.store.db.QueryRow(query, b.Page, pages, string(b.Type), content, b.Order).
		Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *BlockRepository) FindByID(id string) (*model.Block, error) {
	query := `
		SELECT id, page, pages, type, content, display_order, created_at, updated_at
		FROM blocks WHERE id = $1
	`
	var b model.Block
	var content []byte
	var pages []byte
	if err := r.store.db.QueryRow(query, id).Scan(&b.ID, &b.Page, &pages, &b.Type, &content, &b.Order, &b.CreatedAt, &b.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	b.Content = content
	b.Pages = pages
	return &b, nil
}

func (r *BlockRepository) Update(b *model.Block) error {
	query := `
		UPDATE blocks
		SET page = $1,
		    pages = $2,
		    type = $3,
		    content = $4,
		    display_order = $5,
		    updated_at = now()
		WHERE id = $6
		RETURNING created_at, updated_at
	`
	content := []byte("{}")
	if len(b.Content) > 0 {
		content = b.Content
	}
	pages := []byte("[]")
	if len(b.Pages) > 0 {
		pages = b.Pages
	}
	return r.store.db.QueryRow(query, b.Page, pages, string(b.Type), content, b.Order, b.ID).
		Scan(&b.CreatedAt, &b.UpdatedAt)
}

func (r *BlockRepository) Delete(id string) error {
	_, err := r.store.db.Exec("DELETE FROM blocks WHERE id = $1", id)
	return err
}

func (r *BlockRepository) UpdateOrders(updates []store.BlockOrderUpdate) error {
	tx, err := r.store.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE blocks SET display_order = $1, updated_at = now() WHERE id = $2`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, upd := range updates {
		if _, err := stmt.Exec(upd.Order, upd.ID); err != nil {
			return fmt.Errorf("update order for block %s: %w", upd.ID, err)
		}
	}

	return tx.Commit()
}
