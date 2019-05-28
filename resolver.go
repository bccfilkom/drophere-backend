package drophere_go

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"errors"
	"time"

	"github.com/bccfilkom/drophere-go/domain"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type authenticator interface {
	GetAuthenticatedUser(context.Context) *domain.User
}

// Resolver resolves given query from client
type Resolver struct {
	links         []Link
	lastLinkID    int
	userSvc       domain.UserService
	authenticator authenticator
}

// NewResolver func
func NewResolver(userSvc domain.UserService, authenticator authenticator) *Resolver {
	return &Resolver{userSvc: userSvc, authenticator: authenticator}
}

func (r *Resolver) searchLink(ID int) (idx int, found bool) {
	for i, link := range r.links {
		if link.ID == ID {
			idx = i
			found = true
			return
		}
	}

	return
}

func (r *Resolver) searchLinkSlug(slug string) (idx int, found bool) {
	for i, link := range r.links {
		if link.Slug != nil && *link.Slug == slug {
			idx = i
			found = true
			return
		}
	}

	return
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
	r.lastLinkID++
	newLink := Link{
		ID:          r.lastLinkID,
		Title:       title,
		IsProtected: false,
		Slug:        &slug,
		Description: description,
		Deadline:    deadline,
	}
	r.links = append(r.links, newLink)
	return &newLink, nil
}
func (r *mutationResolver) UpdateLink(ctx context.Context, linkID int, title string, slug string, description *string, deadline *time.Time, password *string) (*Message, error) {
	linkIdx, found := r.searchLink(linkID)
	if !found {
		return &Message{Message: "Link not found"}, nil
	}

	r.links[linkIdx].Title = title
	r.links[linkIdx].Slug = &slug
	r.links[linkIdx].Description = description
	r.links[linkIdx].Deadline = deadline

	return &Message{Message: "Link Updated!"}, nil
}
func (r *mutationResolver) DeleteLink(ctx context.Context, linkID int) (*Message, error) {
	linkIdx, found := r.searchLink(linkID)
	if !found {
		return &Message{Message: "Link not found"}, nil
	}

	r.links = append(r.links[:linkIdx], r.links[linkIdx+1:]...)
	return &Message{Message: "Link Deleted!"}, nil
}
func (r *mutationResolver) CheckLinkPassword(ctx context.Context, linkID int, password string) (*Message, error) {
	return &Message{Message: "Not Implemented!"}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	links := make([]*Link, len(r.links))
	for i := range r.links {
		links[i] = &r.links[i]
	}
	return links, nil
}
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errors.New("access denied")
	}
	return &User{
		ID:    int(user.ID),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
func (r *queryResolver) Link(ctx context.Context, slug string) (*Link, error) {
	linkIdx, found := r.searchLinkSlug(slug)
	if !found {
		return nil, errors.New("link not found")
	}

	return &r.links[linkIdx], nil
}
