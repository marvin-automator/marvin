package storage

import (
	"github.com/marvin-automator/marvin/identityproviders/domain"
	"github.com/marvin-automator/marvin/storage"
)

type IdentityStore struct {
	store storage.Store
}

func (is *IdentityStore) SaveIdentity(account, provider string, id domain.Identity) error {
	b, err := is.store.CreateBucketHierarchy("identities_"+account, provider)
	if err != nil {
		return err
	}

	return b.Put(id.ProviderID, id)
}

func (is *IdentityStore) GetIdentity(account, provider, id string) (domain.Identity, error) {
	b, err := is.store.GetBucketFromPath("identities_"+account, provider)
	if err != nil {
		return domain.Identity{}, err
	}

	i := domain.Identity{}
	err = b.Get(id, &i)
	return i, err
}

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

func (is *IdentityStore) DeleteIdentity(account, provider string, i domain.Identity) error {
	b, err := is.store.GetBucketFromPath("identities_"+account, provider)
	if err != nil {
		return err
	}

	return b.Delete(i.ProviderID)
}

func NewIdentityStore(s storage.Store) domain.IdentityStore {
	return &IdentityStore{s}
}
