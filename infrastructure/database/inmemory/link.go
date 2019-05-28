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

// Delete implementation
func (repo *linkRepository) Delete(l *domain.Link) error {

	for i := range repo.db.links {
		if repo.db.links[i].ID == l.ID {
			repo.db.links = append(repo.db.links[:i], repo.db.links[i+1:]...)
			break
		}
	}

	return nil
}

// FindByID implementation
func (repo *linkRepository) FindByID(id uint) (*domain.Link, error) {
	for i := range repo.db.links {
		if repo.db.links[i].ID == id {
			return &repo.db.links[i], nil
		}
	}

	return nil, domain.ErrLinkNotFound
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

// ListByUser implementation
func (repo *linkRepository) ListByUser(userID uint) ([]domain.Link, error) {
	links := make([]domain.Link, 0, len(repo.db.links))
	for _, link := range repo.db.links {
		if link.UserID == userID {
			links = append(links, link)
		}
	}

	return links, nil
}

// Update implementation
func (repo *linkRepository) Update(l *domain.Link) (link *domain.Link, err error) {
	link = l
	for i := range repo.db.links {
		if repo.db.links[i].ID == l.ID {
			repo.db.links[i] = *l
			return
		}
	}
	repo.db.links = append(repo.db.links, *l)
	return
}
