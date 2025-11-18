package store

import (
	"backed-teacher/internal/app/model"
	"time"
)

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindByID(int) (*model.User, error)
	FindByVerificationToken(string) (*model.User, error)
	List() ([]*model.User, error)
	UpdateBanned(id int, banned bool) error
	SaveVerificationToken(userID int, token string, expires time.Time) error
	VerifyEmail(userID int) error
	GetCommentsCount(userID int) (int, error)
}

type AdminRepository interface {
	Create(*model.Admin) error
	FindByEmail(string) (*model.Admin, error)
	FindByID(int) (*model.Admin, error)
	FindByResetToken(string) (*model.Admin, error)
	SaveResetToken(adminID int, token string, expires time.Time) error
	UpdatePassword(adminID int, passwordHash string) error
	ClearResetToken(adminID int) error
}

type BlockFilter struct {
	Page string
}

type BlockOrderUpdate struct {
	ID    string
	Order int
}

type BlockRepository interface {
	List(BlockFilter) ([]*model.Block, error)
	FindByID(id string) (*model.Block, error)
	Create(*model.Block) error
	Update(*model.Block) error
	Delete(id string) error
	UpdateOrders(updates []BlockOrderUpdate) error
}

type PostFilter struct {
	IncludeDrafts bool
}

type PostRepository interface {
	List(PostFilter) ([]*model.Post, error)
	FindByID(id string) (*model.Post, error)
	Create(*model.Post) error
	Update(*model.Post) error
	Delete(id string) error
}

type GalleryRepository interface {
	List() ([]*model.GalleryItem, error)
	FindByID(id string) (*model.GalleryItem, error)
	Create(*model.GalleryItem) error
	Update(*model.GalleryItem) error
	Delete(id string) error
}

type SliderOrderUpdate struct {
	ID    string
	Order int
}

type SliderRepository interface {
	List() ([]*model.SliderItem, error)
	FindByID(id string) (*model.SliderItem, error)
	Create(*model.SliderItem) error
	Update(*model.SliderItem) error
	Delete(id string) error
	UpdateOrders(updates []SliderOrderUpdate) error
}

type CommentRepository interface {
	ListAll() ([]*model.Comment, error)
	ListByPost(postID string) ([]*model.Comment, error)
	FindByID(id string) (*model.Comment, error)
	Create(*model.Comment) error
	Delete(id string) error
}
