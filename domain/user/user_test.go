package user_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	htmlTemplate "html/template"
	textTemplate "text/template"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"
	"github.com/bccfilkom/drophere-go/infrastructure/mailer"
	"github.com/bccfilkom/drophere-go/infrastructure/storageprovider"
	"github.com/bccfilkom/drophere-go/infrastructure/stringgenerator"
)

var (
	authenticator domain.Authenticator
	dummyHasher   domain.Hasher
	mockMailer    domain.Mailer
	strGen        domain.StringGenerator
	htmlTemplates *htmlTemplate.Template
	textTemplates *textTemplate.Template

	storageProviderPool domain.StorageProviderPool
)

func init() {
	authenticator = auth.NewJWTMock()
	dummyHasher = hasher.NewNotAHasher()
	strGen = stringgenerator.NewMock()
	mockMailer = mailer.NewMockMailer()
	stringgenerator.SetMockResult("this_is_not_a_random_string")
	mockStorageProvider := storageprovider.NewMock()
	storageProviderPool.Register(mockStorageProvider)

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

func newRepo() (domain.UserRepository, domain.UserStorageCredentialRepository) {
	memdb := inmemory.New()
	return inmemory.NewUserRepository(memdb), inmemory.NewUserStorageCredentialRepository(memdb)
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

	userRepo, userStorageCredRepo := newRepo()
	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
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

	userRepo, userStorageCredRepo := newRepo()
	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
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

	userRepo, userStorageCredRepo := newRepo()
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
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
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

	userRepo, userStorageCredRepo := newRepo()
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
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
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

	userRepo, userStorageCredRepo := newRepo()
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
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
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

	userRepo, userStorageCredRepo := newRepo()
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

	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
		htmlTemplates,
		textTemplates,
	)

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

func TestConnectStorageProvider(t *testing.T) {
	type test struct {
		userID             uint
		providerID         uint
		providerCredential string
		accountInfo        domain.StorageProviderAccountInfo
		wantErr            error
	}

	userRepo, userStorageCredRepo := newRepo()
	// user1, _ := userRepo.FindByID(1)

	tests := []test{
		{
			userID:  123,
			wantErr: domain.ErrStorageProviderInvalid,
		},
		{
			userID:     123,
			providerID: 1,
			wantErr:    domain.ErrUserNotFound,
		},
		{
			// update existing token
			userID:             1,
			providerID:         1,
			providerCredential: "dropboxToken+0xbadc0de",
			accountInfo: domain.StorageProviderAccountInfo{
				Email: "user_1_another_email@drophere.link",
				Photo: "https://my.photo/user_1.jpg",
			},
			wantErr: nil,
		},
		{
			// create new record
			userID:             357,
			providerID:         1,
			providerCredential: "mockStorageToken+0xbadc0de",
			accountInfo: domain.StorageProviderAccountInfo{
				Email: "user_357_another_email@drophere.link",
				Photo: "https://my.photo/user_357.jpg",
			},
			wantErr: nil,
		},
	}

	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
		htmlTemplates,
		textTemplates,
	)

	for i, tc := range tests {

		storageprovider.SetSharedAccountInfo(tc.accountInfo)

		gotErr := userSvc.ConnectStorageProvider(tc.userID, tc.providerID, tc.providerCredential)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if gotErr == nil {
			ucs, _ := userStorageCredRepo.Find(domain.UserStorageCredentialFilters{
				UserIDs: []uint{tc.userID},
			}, false)

			if tc.providerCredential != ucs[0].ProviderCredential {
				t.Fatalf("test %d: expected: %v, got: %v", i, tc.providerCredential, ucs[0].ProviderCredential)
			}

			if tc.accountInfo.Email != ucs[0].Email {
				t.Fatalf("test %d: expected: %v, got: %v", i, tc.accountInfo.Email, ucs[0].Email)
			}

			if tc.accountInfo.Photo != ucs[0].Photo {
				t.Fatalf("test %d: expected: %v, got: %v", i, tc.accountInfo.Photo, ucs[0].Photo)
			}

		}

	}
}

func TestDisconnectStorageProvider(t *testing.T) {
	type test struct {
		userID     uint
		providerID uint
		wantErr    error
	}

	userRepo, userStorageCredRepo := newRepo()
	// user1, _ := userRepo.FindByID(1)

	tests := []test{
		{
			userID:  123,
			wantErr: domain.ErrStorageProviderInvalid,
		},
		{
			userID:     123,
			providerID: 1,
			wantErr:    domain.ErrUserNotFound,
		},
		{
			// delete existing token
			userID:     1,
			providerID: 1,
			wantErr:    nil,
		},
		{
			// delete empty record
			userID:     357,
			providerID: 1,
			wantErr:    nil,
		},
	}

	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
		htmlTemplates,
		textTemplates,
	)

	for i, tc := range tests {

		gotErr := userSvc.DisconnectStorageProvider(tc.userID, tc.providerID)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if gotErr == nil {
			ucs, _ := userStorageCredRepo.Find(domain.UserStorageCredentialFilters{
				UserIDs:     []uint{tc.userID},
				ProviderIDs: []uint{tc.providerID},
			}, false)

			if len(ucs) > 0 {
				t.Fatalf("test %d: expected: %v, got: %v", i, nil, ucs)
			}

		}

	}
}

func TestListStorageProviders(t *testing.T) {
	type test struct {
		userID            uint
		expectedProviders []domain.UserStorageCredential
		wantErr           error
	}

	userRepo, userStorageCredRepo := newRepo()
	uscsUser1, _ := userStorageCredRepo.Find(domain.UserStorageCredentialFilters{
		UserIDs: []uint{1},
	}, false)

	tests := []test{
		{
			// expect empty
			userID:            123,
			expectedProviders: []domain.UserStorageCredential{},
			wantErr:           nil,
		},
		{
			// expect non-empty credential
			userID:            1,
			expectedProviders: uscsUser1,
			wantErr:           nil,
		},
	}

	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		mockMailer,
		dummyHasher,
		strGen,
		storageProviderPool,
		htmlTemplates,
		textTemplates,
	)

	for _, tc := range tests {

		uscs, gotErr := userSvc.ListStorageProviders(tc.userID)
		assert.ElementsMatch(t, tc.expectedProviders, uscs)
		assert.Equal(t, tc.wantErr, gotErr)

	}
}
