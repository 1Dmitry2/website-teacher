package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"database/sql"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}
	return r.store.db.QueryRow(
		`INSERT INTO users(email, encrypted_password)
			 VALUES ($1,$2)
			 RETURNING id, is_admin, created_at, updated_at`,
		u.Email,
		u.EncryptedPassword,
	).Scan(&u.ID, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, is_admin, created_at, updated_at FROM users WHERE email = $1", email,
	).Scan(
		&u.ID, &u.Email, &u.EncryptedPassword, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, email, encrypted_password, is_admin, created_at, updated_at FROM users WHERE id = $1", id,
	).Scan(
		&u.ID, &u.Email, &u.EncryptedPassword, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) List() ([]*model.User, error) {
	rows, err := r.store.db.Query(`SELECT id, email, is_admin, created_at, updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, rows.Err()
}
