package sqlstore

import (
	"backed-teacher/internal/app/model"
	"backed-teacher/internal/app/store"
	"database/sql"
	"fmt"
	"time"
)

type AdminRepository struct {
	store *Store
}

func (r *AdminRepository) Create(a *model.Admin) error {
	if err := a.Validate(); err != nil {
		return err
	}

	if err := a.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		`INSERT INTO admins (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`,
		a.Email,
		a.PasswordHash,
	).Scan(&a.ID, &a.CreatedAt)
}

func (r *AdminRepository) FindByEmail(email string) (*model.Admin, error) {
	return r.findOneBy("email", email)
}

func (r *AdminRepository) FindByID(id int) (*model.Admin, error) {
	return r.findOneBy("id", id)
}

func (r *AdminRepository) FindByResetToken(token string) (*model.Admin, error) {
	return r.findOneBy("reset_token", token)
}

func (r *AdminRepository) SaveResetToken(adminID int, token string, expires time.Time) error {
	_, err := r.store.db.Exec(
		`UPDATE admins SET reset_token = $1, reset_token_expires = $2 WHERE id = $3`,
		token,
		expires,
		adminID,
	)
	return err
}

func (r *AdminRepository) UpdatePassword(adminID int, passwordHash string) error {
	_, err := r.store.db.Exec(
		`UPDATE admins SET password_hash = $1 WHERE id = $2`,
		passwordHash,
		adminID,
	)
	return err
}

func (r *AdminRepository) ClearResetToken(adminID int) error {
	_, err := r.store.db.Exec(
		`UPDATE admins SET reset_token = NULL, reset_token_expires = NULL WHERE id = $1`,
		adminID,
	)
	return err
}

func (r *AdminRepository) findOneBy(field string, value interface{}) (*model.Admin, error) {
	column := ""
	switch field {
	case "id", "email", "reset_token":
		column = field
	default:
		return nil, fmt.Errorf("unsupported field: %s", field)
	}

	query := fmt.Sprintf(`
		SELECT id, email, password_hash, reset_token, reset_token_expires, created_at
		FROM admins
		WHERE %s = $1
	`, column)

	row := r.store.db.QueryRow(query, value)

	admin := &model.Admin{}
	var resetToken sql.NullString
	var resetExpires sql.NullTime

	err := row.Scan(
		&admin.ID,
		&admin.Email,
		&admin.PasswordHash,
		&resetToken,
		&resetExpires,
		&admin.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	if resetToken.Valid {
		token := resetToken.String
		admin.ResetToken = &token
	}

	if resetExpires.Valid {
		expires := resetExpires.Time
		admin.ResetTokenExpires = &expires
	}

	return admin, nil
}
