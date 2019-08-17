package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type linkRepository struct {
	db *gorm.DB
}

// NewLinkRepository func
func NewLinkRepository(db *gorm.DB) domain.LinkRepository {
	return &linkRepository{db}
}

// Create implementation
func (repo *linkRepository) Create(l *domain.Link) (*domain.Link, error) {
	if err := repo.db.Create(l).Error; err != nil {
		return nil, err
	}
	return l, nil
}

// Delete implementation
func (repo *linkRepository) Delete(l *domain.Link) error {
	return repo.db.Delete(l).Error
}

// FindByID implementation
func (repo *linkRepository) FindByID(id uint) (*domain.Link, error) {
	l := domain.Link{}
	if q := repo.db.
		Preload("User").
		Preload("UserStorageCredential").
		Find(&l, id); q.RecordNotFound() {
		return nil, domain.ErrLinkNotFound
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &l, nil
}

// FindBySlug implementation
func (repo *linkRepository) FindBySlug(slug string) (*domain.Link, error) {
	l := domain.Link{}
	if q := repo.db.
		Where("`slug` = ? ", slug).
		Preload("User").
		Preload("UserStorageCredential").
		Find(&l); q.RecordNotFound() {
		return nil, domain.ErrLinkNotFound
	} else if q.Error != nil {
		return nil, q.Error
	}

	return &l, nil
}

// ListByUser implementation
func (repo *linkRepository) ListByUser(userID uint) ([]domain.Link, error) {
	var links []domain.Link
	if err := repo.db.
		Where("`user_id` = ? ", userID).
		Preload("User").
		Preload("UserStorageCredential").
		Find(&links).
		Error; err != nil {
		return nil, err
	}

	return links, nil
}

// Update implementation
func (repo *linkRepository) Update(l *domain.Link) (link *domain.Link, err error) {
	if err := repo.db.Save(l).Error; err != nil {
		return nil, err
	}

	return l, nil
}
