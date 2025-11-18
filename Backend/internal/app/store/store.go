package store

type Store interface {
	User() UserRepository
	Admin() AdminRepository
	Block() BlockRepository
	Post() PostRepository
	Gallery() GalleryRepository
	Slider() SliderRepository
	Comment() CommentRepository
}
