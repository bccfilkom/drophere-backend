package user_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"
	"github.com/bccfilkom/drophere-go/infrastructure/mailer"
	"github.com/bccfilkom/drophere-go/infrastructure/stringgenerator"
)

var (
	authenticator domain.Authenticator
	dummyHasher   domain.Hasher
	mockMailer    domain.Mailer
	strGen        domain.StringGenerator
	htmlTemplates *htmlTemplate.Template
	textTemplates *textTemplate.Template
)

func init() {
	authenticator = auth.NewJWTMock()
	dummyHasher = hasher.NewNotAHasher()
	strGen = stringgenerator.NewMock()
	mockMailer = mailer.NewMockMailer()
	stringgenerator.SetMockResult("this_is_not_a_random_string")
	var err error

	htmlTemplates, err = htmlTemplate.
		New("request_password_recovery_html").
		Parse("{{.Token}}")
	if err != nil {
		panic(err)
	}

	textTemplates, err = textTemplate.
		New("request_password_recovery_text").
		Parse("{{.Token}}")
	if err != nil {
		panic(err)
	}
}

func newUserRepo() domain.UserRepository {
	memdb := inmemory.New()
	return inmemory.NewUserRepository(memdb)
}

func str2ptr(s string) *string {
	return &s
}

func time2ptr(t time.Time) *time.Time {
	return &t
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

	userSvc := user.NewService(
		newUserRepo(),
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		htmlTemplates,
		textTemplates,
	)

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

	userSvc := user.NewService(
		newUserRepo(),
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		htmlTemplates,
		textTemplates,
	)

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

	userSvc := user.NewService(
		userRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		htmlTemplates,
		textTemplates,
	)

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

	userSvc := user.NewService(
		userRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		htmlTemplates,
		textTemplates,
	)

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

func TestRequestPasswordRecovery(t *testing.T) {
	type test struct {
		email   string
		wantErr error
	}

	userRepo := newUserRepo()
	u, _ := userRepo.FindByEmail("reset+pwd@drophere.link")
	expectedToken := str2ptr("this_is_not_a_random_string")
	emailHTMLTemplate := htmlTemplates.Lookup("request_password_recovery_html")
	emailTextTemplate := textTemplates.Lookup("request_password_recovery_text")
	templateContent := map[string]string{
		"Token": *expectedToken,
	}

	expectedHTMLEmailMessage := &bytes.Buffer{}
	emailHTMLTemplate.Execute(expectedHTMLEmailMessage, templateContent)
	expectedTextEmailMessage := &bytes.Buffer{}
	emailTextTemplate.Execute(expectedTextEmailMessage, templateContent)

	expectedMail := mailer.MockMessage{
		From:         "admin@drophere.link",
		To:           "reset+pwd@drophere.link",
		Title:        "Recover Password",
		MessagePlain: expectedTextEmailMessage.String(),
		MessageHTML:  expectedHTMLEmailMessage.String(),
	}

	tests := []test{
		{email: "", wantErr: domain.ErrUserNotFound},
		{email: "reset+pwd@drophere.link", wantErr: nil},
	}

	userSvc := user.NewService(
		userRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		htmlTemplates,
		textTemplates,
	)

	for i, tc := range tests {
		// reset inbox
		mailer.ClearMessages()

		gotErr := userSvc.RequestPasswordRecovery(tc.email)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if gotErr == nil {
			if !reflect.DeepEqual(u.RecoverPasswordToken, expectedToken) {
				t.Fatalf("test %d: expected: %v, got: %v", i, expectedToken, u.RecoverPasswordToken)
			}
			// TODO: Mock time using https://github.com/bouk/monkey
			if !reflect.DeepEqual(mailer.MockMessages[0], expectedMail) {
				t.Fatalf("test %d: expected: %v, got: %v", i, expectedMail, mailer.MockMessages[0])
			}
		}

	}
}

func TestRecoverPassword(t *testing.T) {
	type test struct {
		email        string
		token        string
		newPassword  string
		expectedUser *domain.User
		wantErr      error
	}

	recoverPasswordToken := "this_is_a_recover_password_token"

	userRepo := newUserRepo()
	u, _ := userRepo.FindByEmail("reset+pwd@drophere.link")
	u.RecoverPasswordToken = str2ptr(recoverPasswordToken)
	u.RecoverPasswordTokenExpiry = time2ptr(time.Now().Add(30 * time.Minute))

	expiredRPTUser, _ := userRepo.FindByEmail("reset+pwd+expired_token@drophere.link")
	expiredRPTUser.RecoverPasswordToken = str2ptr(recoverPasswordToken)

	tests := []test{
		{email: "", wantErr: domain.ErrUserNotFound},
		{email: "reset+pwd@drophere.link", token: "", wantErr: domain.ErrUserNotFound},
		// {email: "reset+pwd@drophere.link", token: recoverPasswordToken, newPassword: "", wantErr: domain.ErrUserNotFound},
		{
			email:       "reset+pwd+expired_token@drophere.link",
			token:       recoverPasswordToken,
			newPassword: "new_password_for_this_user",
			wantErr:     domain.ErrUserPasswordRecoveryTokenExpired,
		},
		{
			email:       "reset+pwd@drophere.link",
			token:       recoverPasswordToken,
			newPassword: "new_password_for_this_user",
			expectedUser: &domain.User{
				ID:                         u.ID,
				Email:                      u.Email,
				Name:                       u.Name,
				Password:                   "new_password_for_this_user",
				RecoverPasswordToken:       nil,
				RecoverPasswordTokenExpiry: nil,
			},
			wantErr: nil,
		},
	}

	userSvc := user.NewService(userRepo, authenticator, mockMailer, dummyHasher, strGen, htmlTemplates, textTemplates)

	for i, tc := range tests {

		gotErr := userSvc.RecoverPassword(tc.email, tc.token, tc.newPassword)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if gotErr == nil {
			if !reflect.DeepEqual(u, tc.expectedUser) {
				t.Fatalf("test %d: expected: %+v, got: %+v", i, tc.expectedUser, u)
			}
		}

	}
}
