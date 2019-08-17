package user

import "github.com/bccfilkom/drophere-go/domain"

// ConnectStorageProvider implementation
func (s *service) ConnectStorageProvider(userID, providerID uint, providerCredential string) error {
	storageProvider, err := s.storageProviderPool.Get(providerID)
	if err != nil {
		return err
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	storageProviderAccount, err := storageProvider.AccountInfo(
		domain.StorageProviderCredential{
			UserAccessToken: providerCredential,
		},
	)
	if err != nil {
		return err
	}

	var cred domain.UserStorageCredential

	creds, err := s.userStorageCredRepo.Find(domain.UserStorageCredentialFilters{
		UserIDs:     []uint{u.ID},
		ProviderIDs: []uint{providerID},
	}, false)
	if err != nil {
		return err
	}

	if len(creds) > 0 {
		cred = creds[0]
		cred.ProviderCredential = providerCredential
		cred.Email = storageProviderAccount.Email
		cred.Photo = storageProviderAccount.Photo
		cred, err = s.userStorageCredRepo.Update(cred)
	} else {
		cred, err = s.userStorageCredRepo.Create(domain.UserStorageCredential{
			UserID:             u.ID,
			ProviderID:         providerID,
			ProviderCredential: providerCredential,
			Email:              storageProviderAccount.Email,
			Photo:              storageProviderAccount.Photo,
		})
	}

	return err

}

// DisconnectStorageProvider implementation
func (s *service) DisconnectStorageProvider(userID, providerID uint) error {
	storageProvider, err := s.storageProviderPool.Get(providerID)
	if err != nil {
		return err
	}

	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	creds, err := s.userStorageCredRepo.Find(domain.UserStorageCredentialFilters{
		UserIDs:     []uint{u.ID},
		ProviderIDs: []uint{storageProvider.ID()},
	}, false)
	if err != nil {
		return err
	}

	if len(creds) > 0 {
		err = s.userStorageCredRepo.Delete(creds[0])
		if err != nil {
			return err
		}
	}

	return nil

}

// ListStorageProviders implementation
func (s *service) ListStorageProviders(userID uint) ([]domain.UserStorageCredential, error) {
	return s.userStorageCredRepo.Find(domain.UserStorageCredentialFilters{
		UserIDs: []uint{userID},
	}, false)
}
