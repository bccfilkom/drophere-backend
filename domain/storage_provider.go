package domain

import (
	"errors"
	"io"
)

var (
	// ErrStorageProviderInvalid error
	ErrStorageProviderInvalid = errors.New("Invalid Storage Provider ID")
)

// StorageProvider domain model
// type StorageProvider struct {
// 	ID   uint
// 	Name string
// }

// StorageProviderCredential stores data needed to access
// storage provider API
type StorageProviderCredential struct {
	UserAccessToken string
}

// StorageProviderAccountInfo domain model
type StorageProviderAccountInfo struct {
	Email string
	Photo string
}

// StorageProviderService abstraction
type StorageProviderService interface {
	ID() uint
	AccountInfo(creds StorageProviderCredential) (StorageProviderAccountInfo, error)
	Upload(creds StorageProviderCredential, file io.Reader, fileName, slug string) error
}

// StorageProviderPool stores a collection of Storage Provider Service
// with provider ID as the key
type StorageProviderPool struct {
	pool map[uint]StorageProviderService
}

// Get returns a StorageProviderService instance identified by its provider ID
func (p *StorageProviderPool) Get(providerID uint) (StorageProviderService, error) {
	sps, ok := p.pool[providerID]
	if !ok {
		return nil, ErrStorageProviderInvalid
	}
	return sps, nil
}

// Register stores a StorageProviderService instance to the pool
func (p *StorageProviderPool) Register(sps StorageProviderService) {
	if p.pool == nil {
		p.pool = make(map[uint]StorageProviderService)
	}

	if sps != nil {
		p.pool[sps.ID()] = sps
	}
}
