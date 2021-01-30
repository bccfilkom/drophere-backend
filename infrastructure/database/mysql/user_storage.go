package mysql

import (
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/jinzhu/gorm"
)

type userStorageCredentialRepository struct {
	db *gorm.DB
}

// NewUserStorageCredentialRepository func
func NewUserStorageCredentialRepository(db *gorm.DB) domain.UserStorageCredentialRepository {
	return &userStorageCredentialRepository{db}
}

// Find implementation
func (repo *userStorageCredentialRepository) Find(filters domain.UserStorageCredentialFilters, withUserRelation bool) ([]domain.UserStorageCredential, error) {
	var (
		creds []domain.UserStorageCredential
		dbQuery = repo.db
	)

	if withUserRelation {
		dbQuery = dbQuery.Preload("User")
	}

	if filters.UserIDs != nil && len(filters.UserIDs) > 0 {
		dbQuery = dbQuery.Where("`user_id` IN (?)", filters.UserIDs)
	}

	if filters.ProviderIDs != nil && len(filters.ProviderIDs) > 0 {
		dbQuery = dbQuery.Where("`provider_id` IN (?)", filters.ProviderIDs)
	}

	err := dbQuery.Find(&creds).Error
	if err != nil {
		return nil, err
	}

	return creds, nil
}

// FindByID implementation
func (repo *userStorageCredentialRepository) FindByID(id uint, withUserRelation bool) (domain.UserStorageCredential, error) {
	var (
		cred domain.UserStorageCredential
		dbQuery = repo.db
	)

	if withUserRelation {
		dbQuery = dbQuery.Preload("User")
	}

	if q := dbQuery.Find(&cred, id); q.RecordNotFound() {
		return cred, domain.ErrUserStorageCredentialNotFound
	} else if q.Error != nil {
		return cred, q.Error
	}

	return cred, nil
}

// Create implementation
func (repo *userStorageCredentialRepository) Create(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	err := repo.db.Create(&cred).Error
	if err != nil {
		return domain.UserStorageCredential{}, err
	}

	return cred, nil
}

// Update implementation
func (repo *userStorageCredentialRepository) Update(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	err := repo.db.Save(&cred).Error
	if err != nil {
		return domain.UserStorageCredential{}, err
	}

	return cred, nil
}

// Delete implementation
func (repo *userStorageCredentialRepository) Delete(cred domain.UserStorageCredential) error {
	return repo.db.Delete(&cred).Error
}
