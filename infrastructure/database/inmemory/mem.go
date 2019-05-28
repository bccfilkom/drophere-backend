package inmemory

import "github.com/bccfilkom/drophere-go/domain"

// DB struct
type DB struct {
	users []domain.User
	links []domain.Link
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
		{ID: 357, Email: "user_357@drophere.link", Name: "User 357", Password: "123456", DropboxToken: nil, DriveToken: nil},
	}

	db.links = []domain.Link{
		{ID: 1, UserID: 1, User: &db.users[0], Title: "Drop file here", Slug: "drop-here", Password: "123098", Description: "drop a file here"},
		{ID: 100, UserID: 1, User: &db.users[0], Title: "Test Link 2", Slug: "test-link-2", Password: "", Description: "no description"},
		{ID: 101, UserID: 357, User: &db.users[1], Title: "Another link", Slug: "another-link", Password: "999", Description: "nil here"},
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

// FindUserByID func
func (db *DB) FindUserByID(id uint) (*domain.User, error) {
	for i, u := range db.users {
		if u.ID == id {
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
