package inmemory

import "github.com/bccfilkom/drophere-go/domain"

type linkRepository struct {
	db *DB
}

// NewLinkRepository func
func NewLinkRepository(db *DB) domain.LinkRepository {
	return &linkRepository{db}
}

// Create implementation
func (repo *linkRepository) Create(l *domain.Link) (*domain.Link, error) {
	repo.db.links = append(repo.db.links, *l)
	return l, nil
}

// FindBySlug implementation
func (repo *linkRepository) FindBySlug(slug string) (*domain.Link, error) {
	for i := range repo.db.links {
		if repo.db.links[i].Slug == slug {
			return &repo.db.links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
}
