package storageprovider

import (
	"io"

	"github.com/bccfilkom/drophere-go/domain"
)

var sharedAccountInfo domain.StorageProviderAccountInfo

type mock struct{}

// SetSharedAccountInfo set the sharedAccountInfo object
func SetSharedAccountInfo(accountInfo domain.StorageProviderAccountInfo) {
	sharedAccountInfo = accountInfo
}

// NewMock returns new mock
func NewMock() domain.StorageProviderService {
	return &mock{}
}

// ID returns provider ID
func (m *mock) ID() uint {
	return 1
}

// AccountInfo mock
func (m *mock) AccountInfo(cred domain.StorageProviderCredential) (domain.StorageProviderAccountInfo, error) {
	return sharedAccountInfo, nil
}

// Upload mock
func (m *mock) Upload(cred domain.StorageProviderCredential, file io.Reader, fileName, slug string) error {
	return nil
}