package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"database/sql"
)

type CommentRepository struct {
	store *Store
}

func (r *CommentRepository) scanRow(scanner interface {
	Scan(dest ...any) error
}) (*model.Comment, error) {
	var c model.Comment
	var replyTo *string
	var userEmail *string
	if err := scanner.Scan(&c.ID, &c.PostID, &c.UserID, &replyTo, &c.Text, &c.IsAdmin, &c.CreatedAt, &userEmail); err != nil {
		return nil, err
	}
	c.ReplyTo = replyTo
	c.UserEmail = userEmail
	return &c, nil
}

func (r *CommentRepository) ListAll() ([]*model.Comment, error) {
	rows, err := r.store.db.Query(`
		SELECT id, post_id, user_id, reply_to, text, is_admin, created_at
		FROM comments
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *CommentRepository) ListByPost(postID string) ([]*model.Comment, error) {
	rows, err := r.store.db.Query(`
		SELECT c.id, c.post_id, c.user_id, c.reply_to, c.text, c.is_admin, c.created_at, u.email
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		c, err := r.scanRow(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *CommentRepository) FindByID(id string) (*model.Comment, error) {
	comment, err := r.scanRow(r.store.db.QueryRow(`
		SELECT c.id, c.post_id, c.user_id, c.reply_to, c.text, c.is_admin, c.created_at, u.email
		FROM comments c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = $1
	`, id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return comment, nil
}

func (r *CommentRepository) Create(cn *model.Comment) error {
	return r.store.db.QueryRow(`
		INSERT INTO comments (post_id, user_id, reply_to, text, is_admin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, cn.PostID, cn.UserID, cn.ReplyTo, cn.Text, cn.IsAdmin).Scan(&cn.ID, &cn.CreatedAt)
}

func (r *CommentRepository) Delete(id string) error {
	_, err := r.store.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
}
