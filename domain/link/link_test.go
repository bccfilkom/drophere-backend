package link_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/link"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"
)

var dummyHasher domain.Hasher

func init() {
	dummyHasher = hasher.NewNotAHasher()
}

func newRepo() (domain.LinkRepository, domain.UserRepository) {
	memdb := inmemory.New()
	return inmemory.NewLinkRepository(memdb), inmemory.NewUserRepository(memdb)
}

func str2ptr(s string) *string {
	return &s
}

func time2ptr(t time.Time) *time.Time {
	return &t
}

func TestCreateLink(t *testing.T) {
	type test struct {
		title       string
		slug        string
		description string
		user        *domain.User
		wantLink    *domain.Link
		wantErr     error
	}

	linkRepo, userRepo := newRepo()
	user, _ := userRepo.FindByID(1)

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
			wantErr:     nil,
		},
	}

	linkSvc := link.NewService(linkRepo, dummyHasher)

	for i, tc := range tests {
		_, gotErr := linkSvc.CreateLink(tc.title, tc.slug, tc.description, tc.user)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
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
		wantLink    *domain.Link
		wantErr     error
	}

	linkRepo, userRepo := newRepo()
	user, _ := userRepo.FindByID(1)

	tests := []test{
		{
			linkID:  123,
			title:   "Drop file here",
			slug:    "drop-here",
			wantErr: domain.ErrLinkNotFound,
		},
		{
			linkID:  100,
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
	}

	linkSvc := link.NewService(linkRepo, dummyHasher)

	for i, tc := range tests {
		gotLink, gotErr := linkSvc.UpdateLink(tc.linkID, tc.title, tc.slug, tc.description, tc.deadline, tc.password)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if !reflect.DeepEqual(gotLink, tc.wantLink) {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantLink, gotLink)
		}
	}

}

func TestDeleteLink(t *testing.T) {
	type test struct {
		linkID  uint
		wantErr error
	}

	linkRepo, _ := newRepo()

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

	linkSvc := link.NewService(linkRepo, dummyHasher)

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

	linkRepo, userRepo := newRepo()

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

	linkSvc := link.NewService(linkRepo, dummyHasher)

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

	linkRepo, userRepo := newRepo()

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

	linkSvc := link.NewService(linkRepo, dummyHasher)

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

	linkRepo, userRepo := newRepo()

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
					ID:          100,
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

	linkSvc := link.NewService(linkRepo, dummyHasher)

	for i, tc := range tests {
		gotLinks, gotErr := linkSvc.ListLinks(tc.userID)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}

		if !reflect.DeepEqual(gotLinks, tc.wantLinks) {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantLinks, gotLinks)
		}
	}

}
