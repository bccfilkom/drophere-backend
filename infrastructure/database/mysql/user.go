package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository func
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db}
}

// Create implementation
func (repo *userRepository) Create(user *domain.User) (*domain.User, error) {
	if err := repo.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByEmail implementation
func (repo *userRepository) FindByEmail(email string) (*domain.User, error) {
	user := domain.User{}
	if q := repo.db.
		Where("`email` = ? ", email).
		Find(&user); q.RecordNotFound() {
		return nil, domain.ErrUserNotFound
	} else if q.Error != nil {
		return nil, q.Error
	}
	return &user, nil
}
