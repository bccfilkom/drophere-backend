package user_test

import (
	"testing"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
)

var authenticator domain.Authenticator

func init() {
	authenticator = auth.NewJWTMock()
}

func newUserRepo() domain.UserRepository {
	memdb := inmemory.New()
	return inmemory.NewUserRepository(memdb)
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
		{email: "user@drophere.link", name: "User", password: "123456", wantErr: nil},
	}

	userSvc := user.NewService(newUserRepo(), authenticator)

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

	userSvc := user.NewService(newUserRepo(), authenticator)

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
