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

// Register resolver
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

// Login resolver
func (r *mutationResolver) Login(ctx context.Context, email string, password string) (*Token, error) {
	userCreds, err := r.userSvc.Auth(email, password)
	if err != nil {
		return nil, err
	}
	return &Token{LoginToken: userCreds.Token}, nil
}

// RequestPasswordRecovery resolver
func (r *mutationResolver) RequestPasswordRecovery(ctx context.Context, email string) (*Message, error) {
	err := r.userSvc.RequestPasswordRecovery(email)
	if err != nil {
		return nil, err
	}

	return &Message{"Recover Password instruction has been sent to your email"}, nil
}

// RecoverPassword resolver
func (r *mutationResolver) RecoverPassword(ctx context.Context, email, recoverToken, newPassword string) (*Token, error) {
	err := r.userSvc.RecoverPassword(email, recoverToken, newPassword)
	if err != nil {
		return nil, err
	}

	userCreds, err := r.userSvc.Auth(email, newPassword)
	if err != nil {
		return nil, err
	}

	return &Token{LoginToken: userCreds.Token}, nil
}

// UpdatePassword resolver
func (r *mutationResolver) UpdatePassword(ctx context.Context, oldPassword string, newPassword string) (*Message, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	_, err := r.userSvc.Update(user.ID, nil, &newPassword, &oldPassword)
	if err != nil {
		return nil, err
	}

	return &Message{Message: "You password successfully updated"}, nil
}

// UpdateProfile resolver
func (r *mutationResolver) UpdateProfile(ctx context.Context, newName string) (*Message, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	_, err := r.userSvc.Update(user.ID, &newName, nil, nil)
	if err != nil {
		return nil, err
	}

	return &Message{Message: "Your profile successfully updated"}, nil
}

// CreateLink resolver
func (r *mutationResolver) CreateLink(ctx context.Context, title string, slug string, description *string, deadline *time.Time, password *string, providerID *int) (*Link, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	desc := ""
	if description != nil {
		desc = *description
	}

	var providerIDUintPtr *uint
	if providerID != nil {
		providerIDUint := uint(*providerID)
		providerIDUintPtr = &providerIDUint
	}

	l, err := r.linkSvc.CreateLink(title, slug, desc, deadline, password, user, providerIDUintPtr)
	if err != nil {
		return nil, err
	}

	return formatLink(*l), nil
}

// UpdateLink resolver
func (r *mutationResolver) UpdateLink(ctx context.Context, linkID int, title string, slug string, description *string, deadline *time.Time, password *string, providerID *int) (*Link, error) {
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

	var providerIDUintPtr *uint
	if providerID != nil {
		providerIDUint := uint(*providerID)
		providerIDUintPtr = &providerIDUint
	}

	l, err = r.linkSvc.UpdateLink(
		uint(linkID),
		title,
		slug,
		description,
		deadline,
		password,
		providerIDUintPtr,
	)

	if err != nil {
		return nil, err
	}

	return formatLink(*l), nil
}

// DeleteLink resolver
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

// CheckLinkPassword resolver
func (r *mutationResolver) CheckLinkPassword(ctx context.Context, linkID int, password string) (*Message, error) {
	// this is for public use, no need to check user auth
	l, err := r.linkSvc.FetchLink(uint(linkID))
	if err != nil {
		return nil, err
	}

	msg := "Invalid Password"
	if r.linkSvc.CheckLinkPassword(l, password) {
		msg = "Valid Password"
	}

	return &Message{Message: msg}, nil
}

// ConnectStorageProvider resolver
func (r *mutationResolver) ConnectStorageProvider(ctx context.Context, providerID int, providerToken string) (*Message, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	err := r.userSvc.ConnectStorageProvider(user.ID, uint(providerID), providerToken)
	if err != nil {
		return nil, err
	}

	return &Message{Message: "Storage Provider successfully connected"}, nil
}

// DisconnectStorageProvider resolver
func (r *mutationResolver) DisconnectStorageProvider(ctx context.Context, providerID int) (*Message, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	err := r.userSvc.DisconnectStorageProvider(user.ID, uint(providerID))
	if err != nil {
		return nil, err
	}

	return &Message{Message: "Storage Provider disconnected"}, nil
}

type queryResolver struct{ *Resolver }

// Links resolver
func (r *queryResolver) Links(ctx context.Context) ([]*Link, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	links, err := r.linkSvc.ListLinks(user.ID)
	if err != nil {
		return nil, err
	}

	return formatLinks(links), nil
}

// Me resolver
func (r *queryResolver) Me(ctx context.Context) (*User, error) {
	user := r.authenticator.GetAuthenticatedUser(ctx)
	if user == nil {
		return nil, errUnauthenticated
	}

	uscs, err := r.userSvc.ListStorageProviders(user.ID)
	if err != nil {
		return nil, err
	}

	// map from domain.UserStorageProviderCredential to StorageProvider
	storageProviders := make([]*StorageProvider, len(uscs))
	for i, usc := range uscs {
		storageProviders[i] = &StorageProvider{
			ID:         int(usc.ID),
			ProviderID: int(usc.ProviderID),
			Email:      usc.Email,
			Photo:      usc.Photo,
		}
	}

	return &User{
		ID:                        int(user.ID),
		Email:                     user.Email,
		Name:                      user.Name,
		ConnectedStorageProviders: storageProviders,
	}, nil
}

// Link resolver
func (r *queryResolver) Link(ctx context.Context, slug string) (*Link, error) {
	// this is for public use, no need to check user auth
	link, err := r.linkSvc.FindLinkBySlug(slug)
	if err != nil {
		return nil, err
	}

	return formatLink(*link), nil
}

func formatLink(link domain.Link) *Link {
	formattedLink := &Link{
		ID:          int(link.ID),
		Title:       link.Title,
		IsProtected: link.IsProtected(),
		Slug:        &link.Slug,
		Description: &link.Description,
		Deadline:    link.Deadline,
	}

	if link.UserStorageCredential != nil {
		formattedLink.StorageProvider = &StorageProvider{
			ID:         int(link.UserStorageCredential.ID),
			ProviderID: int(link.UserStorageCredential.ProviderID),
			Email:      link.UserStorageCredential.Email,
			Photo:      link.UserStorageCredential.Photo,
		}
	}

	return formattedLink
}

func formatLinks(links []domain.Link) []*Link {
	formattedLinks := make([]*Link, len(links))
	for i, link := range links {
		formattedLinks[i] = &Link{
			ID:          int(link.ID),
			Title:       link.Title,
			IsProtected: link.IsProtected(),
			Slug:        &links[i].Slug,
			Description: &links[i].Description,
			Deadline:    link.Deadline,
		}

		if link.UserStorageCredential != nil {
			formattedLinks[i].StorageProvider = &StorageProvider{
				ID:         int(link.UserStorageCredential.ID),
				ProviderID: int(link.UserStorageCredential.ProviderID),
				Email:      link.UserStorageCredential.Email,
				Photo:      link.UserStorageCredential.Photo,
			}
		}
	}
	return formattedLinks
}
