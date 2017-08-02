package storage

import (
	"github.com/marvin-automator/marvin/identityproviders/domain"
	"github.com/marvin-automator/marvin/storage"
)

// IdentityStore is an implementation of the domain.IdentityStore interface
type IdentityStore struct {
	store storage.Store
}

// SaveIdentity saves an identity for the given marvin account and provider.
func (is *IdentityStore) SaveIdentity(account, provider string, id domain.Identity) error {
	b, err := is.store.CreateBucketHierarchy("identities_"+account, provider)
	if err != nil {
		return err
	}

	return b.Put(id.ProviderID, id)
}

// GetIdentity returns the identity for the given marvin account, and provider with the given ID (providerID field in the Identity type)
func (is *IdentityStore) GetIdentity(account, provider, id string) (domain.Identity, error) {
	b, err := is.store.GetBucketFromPath("identities_"+account, provider)
	if err != nil {
		return domain.Identity{}, err
	}

	i := domain.Identity{}
	err = b.Get(id, &i)
	return i, err
}

// GetAccountIdentitiesForProvider returns the identities that are stored for the given account and provider.
func (is *IdentityStore) GetAccountIdentitiesForProvider(account, provider string) ([]domain.Identity, error) {
	b, err := is.store.GetBucketFromPath("identities_"+account, provider)
	if err != nil {
		return []domain.Identity{}, err
	}

	ids := []domain.Identity{}
	err = b.Each(func(id string) error {
		i := domain.Identity{}
		err := b.Get(id, &i)
		if err != nil {
			return err
		}
		ids = append(ids, i)
		return nil
	})

	return ids, err
}

// DeleteIdentity deletes a specific Identity
func (is *IdentityStore) DeleteIdentity(account, provider string, i domain.Identity) error {
	b, err := is.store.GetBucketFromPath("identities_"+account, provider)
	if err != nil {
		return err
	}

	return b.Delete(i.ProviderID)
}

// NewIdentityStore returns a new IdentityStore
func NewIdentityStore(s storage.Store) domain.IdentityStore {
	return &IdentityStore{s}
}
