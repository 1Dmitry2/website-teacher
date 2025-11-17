package store

import "backed-teacher/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindByID(int) (*model.User, error)
}
