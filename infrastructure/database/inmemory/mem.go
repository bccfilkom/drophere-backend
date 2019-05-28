package inmemory

import "github.com/bccfilkom/drophere-go/domain"

// DB struct
type DB struct {
	users []domain.User
}

// New func
func New() *DB {
	db := &DB{}
	db.populate()
	return db
}

func (db *DB) populate() {
	db.users = []domain.User{
		{ID: 1, Email: "user@drophere.link", Name: "User", Password: "123456", DropboxToken: nil, DriveToken: nil},
	}
}

// FindUserByEmail func
func (db *DB) FindUserByEmail(email string) (*domain.User, error) {
	for i, u := range db.users {
		if u.Email == email {
			return &db.users[i], nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// CreateUser func
func (db *DB) CreateUser(u *domain.User) (*domain.User, error) {
	db.users = append(db.users, *u)
	return u, nil
}
