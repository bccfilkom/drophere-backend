package link

import (
	"time"

	"github.com/bccfilkom/drophere-go/domain"
)

type service struct {
	linkRepo       domain.LinkRepository
	passwordHasher domain.Hasher
}

// NewService returns new service instance
func NewService(linkRepo domain.LinkRepository, passwordHasher domain.Hasher) domain.LinkService {
	return &service{linkRepo, passwordHasher}
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

// UpdateLink updates existing Link and save it to repository
func (s *service) UpdateLink(linkID uint, title, slug string, description *string, deadline *time.Time, password *string) (*domain.Link, error) {
	l, err := s.linkRepo.FindByID(linkID)
	if err != nil {
		return nil, err
	}

	// check duplicated slug
	link2, err := s.linkRepo.FindBySlug(slug)
	if err != nil && err != domain.ErrLinkNotFound {
		return nil, err
	}

	if link2 != nil && link2.ID != l.ID {
		return nil, domain.ErrLinkDuplicatedSlug
	}

	l.Title = title
	l.Slug = slug
	l.Deadline = deadline // set null if the user want to remove the deadline
	if description != nil {
		l.Description = *description
	}

	if password != nil {
		l.SetPassword(*password, s.passwordHasher)
	}

	return s.linkRepo.Update(l)
}

// DeleteLink delete existing Link specified by its ID
func (s *service) DeleteLink(id uint) error {
	l, err := s.linkRepo.FindByID(id)
	if err != nil {
		return err
	}

	return s.linkRepo.Delete(l)
}

// FetchLink returns single Link identified by its ID
func (s *service) FetchLink(id uint) (*domain.Link, error) {
	return s.linkRepo.FindByID(id)
}

// FindLinkBySlug returns single Link identified by its slug
func (s *service) FindLinkBySlug(slug string) (*domain.Link, error) {
	return s.linkRepo.FindBySlug(slug)
}

// ListLinks returns list of Link which belongs to a user
func (s *service) ListLinks(userID uint) ([]domain.Link, error) {
	return s.linkRepo.ListByUser(userID)
}
