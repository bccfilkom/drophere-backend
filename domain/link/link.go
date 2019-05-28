package link

import (
	"github.com/bccfilkom/drophere-go/domain"
)

type service struct {
	linkRepo domain.LinkRepository
}

// NewService returns new service instance
func NewService(linkRepo domain.LinkRepository) domain.LinkService {
	return &service{linkRepo}
}

// CreateLink creates new Link and store it to repository
func (s *service) CreateLink(title, slug, description string, user *domain.User) (*domain.Link, error) {
	l, err := s.linkRepo.FindBySlug(slug)
	if err != nil && err != domain.ErrLinkNotFound {
		return nil, err
	}

	if l != nil {
		return nil, domain.ErrLinkDuplicatedSlug
	}

	l = &domain.Link{
		UserID:      user.ID,
		Title:       title,
		Slug:        slug,
		Description: description,
	}

	return s.linkRepo.Create(l)
}
