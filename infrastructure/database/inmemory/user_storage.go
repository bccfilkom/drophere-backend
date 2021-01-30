package inmemory

import "github.com/bccfilkom/drophere-go/domain"

type userStorageCredentialRepository struct {
	db *DB
}

// NewUserStorageCredentialRepository func
func NewUserStorageCredentialRepository(db *DB) domain.UserStorageCredentialRepository {
	return &userStorageCredentialRepository{db}
}

func isInUintSlice(u uint, slice []uint) bool {
	for _, el := range slice {
		if el == u {
			return true
		}
	}
	return false
}

// Find impl
func (repo *userStorageCredentialRepository) Find(filters domain.UserStorageCredentialFilters, withUserRelation bool) ([]domain.UserStorageCredential, error) {
	creds := make([]domain.UserStorageCredential, 0)
	usersCache := make(map[uint]domain.User)

	// load users first
	if withUserRelation && len(filters.UserIDs) > 0 {
		for _, u := range repo.db.users {
			if isInUintSlice(u.ID, filters.UserIDs) {
				usersCache[u.ID] = u
			}
		}
	}

	for _, usc := range repo.db.userStorageCreds {
		if filters.UserIDs != nil && (len(filters.UserIDs) == 0 ||
			!isInUintSlice(usc.UserID, filters.UserIDs)) {
			continue
		}

		if filters.ProviderIDs != nil && (len(filters.ProviderIDs) == 0 ||
			!isInUintSlice(usc.ProviderID, filters.ProviderIDs)) {
			continue
		}

		if withUserRelation {
			usc.User = usersCache[usc.UserID]
		}

		creds = append(creds, usc)
	}

	return creds, nil
}

// FindByID impl
func (repo *userStorageCredentialRepository) FindByID(id uint, withUserRelation bool) (domain.UserStorageCredential, error) {
	cred := domain.UserStorageCredential{}
	found := false
	for _, usc := range repo.db.userStorageCreds {
		if usc.ID == id {
			cred = usc
			found = true
			break
		}
	}

	if found {
		if withUserRelation {
			for _, u := range repo.db.users {
				if u.ID == cred.UserID {
					cred.User = u
					break
				}
			}
		}

		return cred, nil
	}
	return cred, domain.ErrUserStorageCredentialNotFound
}

// Create impl
func (repo *userStorageCredentialRepository) Create(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {
	repo.db.userStorageCreds = append(repo.db.userStorageCreds, cred)
	return cred, nil
}

// Update impl
func (repo *userStorageCredentialRepository) Update(cred domain.UserStorageCredential) (domain.UserStorageCredential, error) {

	for i := range repo.db.userStorageCreds {
		if repo.db.userStorageCreds[i].ID == cred.ID {
			repo.db.userStorageCreds[i] = cred
			break
		}
	}

	return cred, nil
}

// Delete impl
func (repo *userStorageCredentialRepository) Delete(cred domain.UserStorageCredential) error {

	for i := range repo.db.userStorageCreds {
		if repo.db.userStorageCreds[i].ID == cred.ID {
			repo.db.userStorageCreds = append(repo.db.userStorageCreds[:i], repo.db.userStorageCreds[i+1:]...)
			break
		}
	}

	return nil
}
