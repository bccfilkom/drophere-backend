package user_test

import (
	"reflect"
	"testing"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"
)

var (
	authenticator domain.Authenticator
	dummyHasher   domain.Hasher
)

func init() {
	authenticator = auth.NewJWTMock()
	dummyHasher = hasher.NewNotAHasher()
}

func newUserRepo() domain.UserRepository {
	memdb := inmemory.New()
	return inmemory.NewUserRepository(memdb)
}

func str2ptr(s string) *string {
	return &s
}

func TestRegister(t *testing.T) {
	type test struct {
		email    string
		name     string
		password string
		wantUser *domain.User
		wantErr  error
	}

	tests := []test{
		{email: "user@drophere.link", name: "User", password: "123456", wantErr: domain.ErrUserDuplicated},
		{email: "new_user@drophere.link", name: "New User", password: "123456", wantErr: nil},
	}

	userSvc := user.NewService(newUserRepo(), authenticator, dummyHasher)

	for i, tc := range tests {
		_, gotErr := userSvc.Register(tc.email, tc.name, tc.password)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
	}
}

func TestAuth(t *testing.T) {
	type test struct {
		email     string
		password  string
		wantCreds *domain.UserCredentials
		wantErr   error
	}

	tests := []test{
		{email: "", password: "", wantErr: domain.ErrUserNotFound},
		{email: "user@drophere.link", password: "", wantErr: domain.ErrUserInvalidPassword},
		{email: "user@drophere.link", password: "123456", wantCreds: &domain.UserCredentials{Token: "user_token_1"}},
	}

	userSvc := user.NewService(newUserRepo(), authenticator, dummyHasher)

	for i, tc := range tests {
		gotCreds, gotErr := userSvc.Auth(tc.email, tc.password)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
		if gotCreds != nil && gotCreds.Token != tc.wantCreds.Token {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantCreds.Token, gotCreds.Token)
		}
	}
}

func TestUpdateStorageToken(t *testing.T) {
	type test struct {
		userID       uint
		dropboxToken *string
		wantUser     *domain.User
		wantErr      error
	}

	userRepo := newUserRepo()
	u, _ := userRepo.FindByID(1)

	tests := []test{
		{userID: 123, wantErr: domain.ErrUserNotFound, wantUser: nil},
		{
			userID:       1,
			dropboxToken: str2ptr("my_dropbox_token_here"),
			wantUser: &domain.User{
				ID:           u.ID,
				Email:        u.Email,
				Name:         u.Name,
				Password:     u.Password,
				DropboxToken: str2ptr("my_dropbox_token_here"),
			},
		},
	}

	userSvc := user.NewService(userRepo, authenticator, dummyHasher)

	for i, tc := range tests {
		gotUser, gotErr := userSvc.UpdateStorageToken(tc.userID, tc.dropboxToken)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
		if !reflect.DeepEqual(gotUser, tc.wantUser) {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantUser, gotUser)
		}
	}
}

func TestUpdate(t *testing.T) {
	type test struct {
		userID      uint
		name        *string
		password    *string
		oldPassword *string
		wantUser    *domain.User
		wantErr     error
	}

	userRepo := newUserRepo()
	u, _ := userRepo.FindByID(1)

	tests := []test{
		{userID: 123, wantErr: domain.ErrUserNotFound, wantUser: nil},
		{userID: 1, password: str2ptr("new_password123"), oldPassword: nil, wantErr: domain.ErrUserInvalidPassword},
		{userID: 1, password: str2ptr("new_password123"), oldPassword: str2ptr(""), wantErr: domain.ErrUserInvalidPassword},
		{
			userID:      1,
			name:        str2ptr("new name 123"),
			password:    str2ptr("new_password123"),
			oldPassword: str2ptr("123456"),
			wantUser: &domain.User{
				ID:       u.ID,
				Email:    u.Email,
				Name:     "new name 123",
				Password: "new_password123",
			},
		},
	}

	userSvc := user.NewService(userRepo, authenticator, dummyHasher)

	for i, tc := range tests {
		gotUser, gotErr := userSvc.Update(tc.userID, tc.name, tc.password, tc.oldPassword)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
		if !reflect.DeepEqual(gotUser, tc.wantUser) {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantUser, gotUser)
		}
	}
}
