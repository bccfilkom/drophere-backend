package drophere_go

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"errors"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var (
	errUnauthenticated = errors.New("Access denied")
	errUnauthorized    = errors.New("You are not allowed to do this operation")
)

type authenticator interface {
	GetAuthenticatedUser(context.Context) *domain.User
}

// Resolver resolves given query from client
type Resolver struct {
	linkSvc       domain.LinkService
	userSvc       domain.UserService
	authenticator authenticator
}

// NewResolver func
func NewResolver(
	userSvc domain.UserService,
	authenticator authenticator,
	linkSvc domain.LinkService,
) *Resolver {
	return &Resolver{
		linkSvc:       linkSvc,
		userSvc:       userSvc,
		authenticator: authenticator,
	}
}

// Mutation returns a group of resolvers for mutation query
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns a group of resolvers for query
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Register(ctx context.Context, email string, password string, name string) (*Token, error) {
	user, err := r.userSvc.Register(email, name, password)
	if err != nil {
		return nil, err
	}

	userCreds, err := r.userSvc.Auth(user.Email, password)
	if err != nil {
		return nil, err
	}
	return &Token{LoginToken: userCreds.Token}, nil
}
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*Token, error) {
	userCreds, err := r.userSvc.Auth(email, password)
	if err != nil {
		return nil, err
	}
	return &Token{LoginToken: userCreds.Token}, nil
}
func (r *mutationResolver) UpdatePassword(ctx context.Context, oldPassword string, newPassword string) (*Message, error) {
	return &Message{Message: "update password: OK"}, nil
}
func (r *mutationResolver) UpdateProfile(ctx context.Context, newName string) (*Message, error) {
	return &Message{Message: "update profile: OK"}, nil
}
func (r *mutationResolver) CreateLink(ctx context.Context, title string, slug string, description *string, deadline *time.Time, password *string) (*Link, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	desc := ""
	if description != nil {
		desc = *description
	}

	l, err := r.linkSvc.CreateLink(title, slug, desc, user)
	if err != nil {
		return nil, err
	}

	return &Link{
		ID:          int(l.ID),
		Title:       l.Title,
		IsProtected: l.IsProtected(),
		Slug:        &l.Slug,
		Description: &l.Description,
		Deadline:    l.Deadline,
	}, nil
}
func (r *mutationResolver) UpdateLink(ctx context.Context, linkID int, title string, slug string, description *string, deadline *time.Time, password *string) (*Message, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	l, err := r.linkSvc.FetchLink(uint(linkID))
	if err != nil {
		return nil, err
	}

	if l.UserID != user.ID {
		return nil, errUnauthorized
	}

	_, err = r.linkSvc.UpdateLink(
		uint(linkID),
		title,
		slug,
		description,
		deadline,
		password,
	)

	if err != nil {
		return nil, err
	}

	return &Message{Message: "Link Updated!"}, nil
}
func (r *mutationResolver) DeleteLink(ctx context.Context, linkID int) (*Message, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	l, err := r.linkSvc.FetchLink(uint(linkID))
	if err != nil {
		return nil, err
	}

	if l.UserID != user.ID {
		return nil, errUnauthorized
	}

	err = r.linkSvc.DeleteLink(uint(linkID))
	if err != nil {
		return nil, err
	}

	return &Message{Message: "Link Deleted!"}, nil
}
func (r *mutationResolver) CheckLinkPassword(ctx context.Context, linkID int, password string) (*Message, error) {
	return &Message{Message: "Not Implemented!"}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	links, err := r.linkSvc.ListLinks(user.ID)
	if err != nil {
		return nil, err
	}

	formattedLinks := make([]*Link, len(links))
	for i := range links {
		formattedLinks[i] = &Link{
			ID:          int(links[i].ID),
			Title:       links[i].Title,
			IsProtected: links[i].IsProtected(),
			Slug:        &links[i].Slug,
			Description: &links[i].Description,
			Deadline:    links[i].Deadline,
		}
	}
	return formattedLinks, nil
}
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}
	return &User{
		ID:    int(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
func (r *queryResolver) Link(ctx context.Context, slug string) (*Link, error) {
	// this is for public use, no need to check user auth
	link, err := r.linkSvc.FindLinkBySlug(slug)
	if err != nil {
		return nil, err
	}

	return &Link{
		ID:          int(link.ID),
		Title:       link.Title,
		IsProtected: link.IsProtected(),
		Slug:        &link.Slug,
		Description: &link.Description,
		Deadline:    link.Deadline,
	}, nil
}
