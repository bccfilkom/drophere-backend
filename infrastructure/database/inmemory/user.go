package inmemory

import "github.com/bccfilkom/drophere-go/domain"

type userRepository struct {
	db *DB
}

// NewUserRepository func
func NewUserRepository(db *DB) domain.UserRepository {
	return &userRepository{db}
}

// Create implementation
func (repo *userRepository) Create(user *domain.User) (*domain.User, error) {
	return repo.db.CreateUser(user)
}

// FindByEmail implementation
func (repo *userRepository) FindByEmail(email string) (*domain.User, error) {
	return repo.db.FindUserByEmail(email)
}
