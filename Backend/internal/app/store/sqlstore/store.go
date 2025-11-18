package sqlstore

import (
	"backed-teacher/internal/app/store"
	"database/sql"
	_ "github.com/lib/pq"
)

type Store struct {
	db              *sql.DB
	userRepository  *UserRepository
	adminRepository *AdminRepository
	blockRepository *BlockRepository
	postRepository  *PostRepository
	galleryRepo     *GalleryRepository
	sliderRepo      *SliderRepository
	commentRepo     *CommentRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

func (s *Store) Admin() store.AdminRepository {
	if s.adminRepository != nil {
		return s.adminRepository
	}
	s.adminRepository = &AdminRepository{
		store: s,
	}
	return s.adminRepository
}

func (s *Store) Block() store.BlockRepository {
	if s.blockRepository != nil {
		return s.blockRepository
	}
	s.blockRepository = &BlockRepository{store: s}
	return s.blockRepository
}

func (s *Store) Post() store.PostRepository {
	if s.postRepository != nil {
		return s.postRepository
	}
	s.postRepository = &PostRepository{store: s}
	return s.postRepository
}

func (s *Store) Gallery() store.GalleryRepository {
	if s.galleryRepo != nil {
		return s.galleryRepo
	}
	s.galleryRepo = &GalleryRepository{store: s}
	return s.galleryRepo
}

func (s *Store) Slider() store.SliderRepository {
	if s.sliderRepo != nil {
		return s.sliderRepo
	}
	s.sliderRepo = &SliderRepository{store: s}
	return s.sliderRepo
}

func (s *Store) Comment() store.CommentRepository {
	if s.commentRepo != nil {
		return s.commentRepo
	}
	s.commentRepo = &CommentRepository{store: s}
	return s.commentRepo
}
