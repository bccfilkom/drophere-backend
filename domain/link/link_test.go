package link_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/link"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"

	"github.com/stretchr/testify/assert"
)

var dummyHasher domain.Hasher

func init() {
	dummyHasher = hasher.NewNotAHasher()
}

func newRepo() (domain.LinkRepository, domain.UserRepository, domain.UserStorageCredentialRepository) {
	memdb := inmemory.New()
	return inmemory.NewLinkRepository(memdb), inmemory.NewUserRepository(memdb), inmemory.NewUserStorageCredentialRepository(memdb)
}

func str2ptr(s string) *string {
	return &s
}

func time2ptr(t time.Time) *time.Time {
	return &t
}

func uint2ptr(u uint) *uint {
	return &u
}

func TestCheckLinkPassword(t *testing.T) {
	type test struct {
		link       *domain.Link
		password   string
		wantResult bool
	}

	linkRepo, _, uscRepo := newRepo()
	getLink := func(id uint) *domain.Link {
		l, _ := linkRepo.FindByID(id)
		return l
	}
	tests := []test{
		{
			link:       getLink(1),
			password:   "",
			wantResult: false,
		},
		{
			link:       getLink(1),
			password:   "abcdef",
			wantResult: false,
		},
		{
			link:       getLink(1),
			password:   "123098",
			wantResult: true,
		},
		{
			link:       getLink(2),
			password:   "123098",
			wantResult: true,
		},
		{
			link:       getLink(2),
			password:   "",
			wantResult: true,
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for i, tc := range tests {
		gotResult := linkSvc.CheckLinkPassword(tc.link, tc.password)
		if gotResult != tc.wantResult {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantResult, gotResult)
		}
	}

}

func TestCreateLink(t *testing.T) {
	type test struct {
		title       string
		slug        string
		description string
		deadline    *time.Time
		password    *string
		user        *domain.User
		providerID  *uint
		wantLink    *domain.Link
		wantErr     error
	}

	linkRepo, userRepo, uscRepo := newRepo()
	user, _ := userRepo.FindByID(1)
	uscUser1, _ := uscRepo.FindByID(2000, false)

	linkDeadline := time.Date(2020, time.November, 11, 1, 2, 3, 0, time.UTC)

	tests := []test{
		{
			title:       "Drop file here",
			slug:        "drop-here",
			description: "drop a file here",
			user:        user,
			wantErr:     domain.ErrLinkDuplicatedSlug,
		},
		{
			title:       "Drop CV",
			slug:        "yoursummerintern",
			description: "Drop your CV for summer internship",
			user:        user,
			wantLink: &domain.Link{
				ID:          4,
				UserID:      user.ID,
				Title:       "Drop CV",
				Slug:        "yoursummerintern",
				Description: "Drop your CV for summer internship",
			},
			wantErr: nil,
		},
		{
			title:       "Link with associated storage provider",
			slug:        "linktomockbox",
			description: "hello there, please upload a file",
			user:        user,
			providerID:  uint2ptr(1234),
			wantErr:     domain.ErrUserStorageCredentialNotFound,
		},
		{
			title:       "Link with associated storage provider",
			slug:        "linktomockbox",
			description: "hello there, please upload a file",
			user:        user,
			providerID:  uint2ptr(1),
			wantLink: &domain.Link{
				ID:                      5,
				UserID:                  user.ID,
				Title:                   "Link with associated storage provider",
				Slug:                    "linktomockbox",
				Description:             "hello there, please upload a file",
				Password:                "",
				UserStorageCredentialID: uint2ptr(2000),
				UserStorageCredential:   &uscUser1,
			},
			wantErr: nil,
		},
		{
			title:       "Link with associated storage provider",
			slug:        "guarded-with-pwd-and-deadline",
			description: "hello there, please upload a file",
			deadline:    &linkDeadline,
			password:    str2ptr("abcdef"),
			user:        user,
			providerID:  uint2ptr(1),
			wantLink: &domain.Link{
				ID:                      6,
				UserID:                  user.ID,
				Title:                   "Link with associated storage provider",
				Slug:                    "guarded-with-pwd-and-deadline",
				Description:             "hello there, please upload a file",
				Deadline:                &linkDeadline,
				Password:                "abcdef",
				UserStorageCredentialID: uint2ptr(2000),
				UserStorageCredential:   &uscUser1,
			},
			wantErr: nil,
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for _, tc := range tests {
		gotLink, gotErr := linkSvc.CreateLink(tc.title, tc.slug, tc.description, tc.deadline, tc.password, tc.user, tc.providerID)

		assert.Equal(t, tc.wantErr, gotErr)
		assert.Equal(t, tc.wantLink, gotLink)
	}

}

func TestUpdateLink(t *testing.T) {
	type test struct {
		linkID      uint
		title       string
		slug        string
		description *string
		deadline    *time.Time
		password    *string
		providerID  *uint
		wantLink    *domain.Link
		wantErr     error
	}

	linkRepo, userRepo, uscRepo := newRepo()
	user, _ := userRepo.FindByID(1)
	uscUser1, _ := uscRepo.FindByID(2000, false)

	tests := []test{
		{
			linkID:  123,
			title:   "Drop file here",
			slug:    "drop-here",
			wantErr: domain.ErrLinkNotFound,
		},
		{
			linkID:  2,
			title:   "Drop file here",
			slug:    "drop-here",
			wantErr: domain.ErrLinkDuplicatedSlug,
		},
		{
			linkID:      1,
			title:       "Drop CV 2",
			slug:        "yoursummerintern2",
			description: str2ptr("Drop your CV for summer internship 2019"),
			deadline:    time2ptr(time.Date(2019, 1, 2, 3, 0, 0, 0, time.Local)),
			password:    str2ptr("123098"),
			wantErr:     nil,
			wantLink: &domain.Link{
				ID:          1,
				Title:       "Drop CV 2",
				Slug:        "yoursummerintern2",
				Description: "Drop your CV for summer internship 2019",
				Deadline:    time2ptr(time.Date(2019, 1, 2, 3, 0, 0, 0, time.Local)),
				Password:    "123098",
				UserID:      user.ID,
				User:        user,
			},
		},
		{
			linkID:      1,
			title:       "Drop CV 2",
			slug:        "yoursummerintern2",
			description: str2ptr("Drop your CV for summer internship 2019"),
			deadline:    time2ptr(time.Date(2019, 1, 2, 3, 0, 0, 0, time.Local)),
			password:    str2ptr("123098"),
			providerID:  uint2ptr(1234),
			wantErr:     domain.ErrUserStorageCredentialNotFound,
		},
		{
			linkID:      1,
			title:       "Drop CV 2 With MockBox",
			slug:        "yoursummerintern2mockbox",
			description: str2ptr("Drop your CV for summer internship 2019"),
			deadline:    time2ptr(time.Date(2019, 1, 2, 3, 0, 0, 0, time.Local)),
			password:    str2ptr("123098"),
			providerID:  uint2ptr(1),
			wantErr:     nil,
			wantLink: &domain.Link{
				ID:                      1,
				Title:                   "Drop CV 2 With MockBox",
				Slug:                    "yoursummerintern2mockbox",
				Description:             "Drop your CV for summer internship 2019",
				Deadline:                time2ptr(time.Date(2019, 1, 2, 3, 0, 0, 0, time.Local)),
				Password:                "123098",
				UserID:                  user.ID,
				User:                    user,
				UserStorageCredentialID: uint2ptr(uscUser1.ID),
				UserStorageCredential:   &uscUser1,
			},
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for _, tc := range tests {
		gotLink, gotErr := linkSvc.UpdateLink(tc.linkID, tc.title, tc.slug, tc.description, tc.deadline, tc.password, tc.providerID)

		assert.Equal(t, tc.wantErr, gotErr)
		assert.Equal(t, tc.wantLink, gotLink)
	}

}

func TestDeleteLink(t *testing.T) {
	type test struct {
		linkID  uint
		wantErr error
	}

	linkRepo, _, uscRepo := newRepo()

	tests := []test{
		{
			linkID:  123,
			wantErr: domain.ErrLinkNotFound,
		},
		{
			linkID:  1,
			wantErr: nil,
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for i, tc := range tests {
		gotErr := linkSvc.DeleteLink(tc.linkID)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
	}

}

func TestFetchLink(t *testing.T) {
	type test struct {
		linkID   uint
		wantErr  error
		wantLink *domain.Link
	}

	linkRepo, userRepo, uscRepo := newRepo()

	user, _ := userRepo.FindByID(1)

	tests := []test{
		{
			linkID:  123,
			wantErr: domain.ErrLinkNotFound,
		},
		{
			linkID:  1,
			wantErr: nil,
			wantLink: &domain.Link{
				ID:          1,
				UserID:      user.ID,
				User:        user,
				Title:       "Drop file here",
				Slug:        "drop-here",
				Password:    "123098",
				Description: "drop a file here",
			},
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for i, tc := range tests {
		gotLink, gotErr := linkSvc.FetchLink(tc.linkID)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if !reflect.DeepEqual(gotLink, tc.wantLink) {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantLink, gotLink)
		}
	}

}

func TestFindLinkBySlug(t *testing.T) {
	type test struct {
		slug     string
		wantErr  error
		wantLink *domain.Link
	}

	linkRepo, userRepo, uscRepo := newRepo()

	user, _ := userRepo.FindByID(1)

	tests := []test{
		{
			slug:    "123",
			wantErr: domain.ErrLinkNotFound,
		},
		{
			slug:    "drop-here",
			wantErr: nil,
			wantLink: &domain.Link{
				ID:          1,
				UserID:      user.ID,
				User:        user,
				Title:       "Drop file here",
				Slug:        "drop-here",
				Password:    "123098",
				Description: "drop a file here",
			},
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for i, tc := range tests {
		gotLink, gotErr := linkSvc.FindLinkBySlug(tc.slug)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if !reflect.DeepEqual(gotLink, tc.wantLink) {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantLink, gotLink)
		}
	}

}

func TestListLinks(t *testing.T) {
	type test struct {
		userID    uint
		wantErr   error
		wantLinks []domain.Link
	}

	linkRepo, userRepo, uscRepo := newRepo()

	user, _ := userRepo.FindByID(1)

	tests := []test{
		{
			userID:    123,
			wantErr:   nil,
			wantLinks: []domain.Link{},
		},
		{
			userID:  1,
			wantErr: nil,
			wantLinks: []domain.Link{
				{
					ID:          1,
					UserID:      user.ID,
					User:        user,
					Title:       "Drop file here",
					Slug:        "drop-here",
					Password:    "123098",
					Description: "drop a file here",
				},
				{
					ID:          2,
					UserID:      user.ID,
					User:        user,
					Title:       "Test Link 2",
					Slug:        "test-link-2",
					Password:    "",
					Description: "no description",
				},
			},
		},
	}

	linkSvc := link.NewService(linkRepo, uscRepo, dummyHasher)

	for _, tc := range tests {
		gotLinks, gotErr := linkSvc.ListLinks(tc.userID)

		assert.Equal(t, tc.wantErr, gotErr)
		assert.Equal(t, tc.wantLinks, gotLinks)

	}

}
