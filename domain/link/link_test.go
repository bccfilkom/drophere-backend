package link_test

import (
	"testing"

	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/link"
	"github.com/bccfilkom/drophere-go/infrastructure/database/inmemory"
)

func newRepo() (domain.LinkRepository, domain.UserRepository) {
	memdb := inmemory.New()
	return inmemory.NewLinkRepository(memdb), inmemory.NewUserRepository(memdb)
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

	linkSvc := link.NewService(linkRepo)

	for i, tc := range tests {
		_, gotErr := linkSvc.CreateLink(tc.title, tc.slug, tc.description, tc.user)
		if gotErr != tc.wantErr {
			t.Fatalf("test %d: expected: %v, got: %v", i, tc.wantErr, gotErr)
		}
	}

}
