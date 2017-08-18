package storage

import (
	"github.com/marvin-automator/marvin/actions/domain"
	"time"
	"github.com/marvin-automator/marvin/storage"
)

// TemporaryInstanceStore is an implementation of the domain.TemporaryInstanceStore interface
type TemporaryInstanceStore struct {
	s storage.Store
}

// NewTemporaryInstanceStore creates a new TemporaryInstanceStore.
func NewTemporaryInstanceStore(s storage.Store) domain.TemporaryInstanceStore {
	return &TemporaryInstanceStore{s}
}

// Save a TemporaryActionInstance to the store.
func (ts *TemporaryInstanceStore) Save(tai domain.TemporaryActionInstance) error {
	b, err := ts.getOrCreateBucket()
	if err != nil {
		return err
	}

	return b.Put(tai.ID, tai)
}

// Get a TemporaryActionInstance by its ID.
func (ts *TemporaryInstanceStore) Get(ID string) (domain.TemporaryActionInstance, error) {
	tai := domain.TemporaryActionInstance{}

	b, err := ts.getOrCreateBucket()
	if err != nil {
		return tai, err
	}

	err = b.Get(ID, &tai)
	return tai, err
}

// DeleteInstancesOlderThan the given time.
func (ts *TemporaryInstanceStore) DeleteInstancesOlderThan(t time.Time) (int, error) {
	b, err := ts.getOrCreateBucket()
	if err != nil {
		return 0, err
	}

	n := 0
	err = b.Each(func(ID string) error {
		tai, err := ts.Get(ID)
		if err != nil {
			return err
		}

		if t.After(tai.Created) {
			err = b.Delete(tai.ID)
			if err != nil {
				return err
			}
			n += 1
		}
		return nil
	})

	return n, err
}

func (ts *TemporaryInstanceStore) getOrCreateBucket() (storage.Bucket, error) {
	return ts.s.Bucket("temporary_instances")
}