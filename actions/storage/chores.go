package storage

import (
	"github.com/marvin-automator/marvin/actions/domain"
	"github.com/marvin-automator/marvin/storage"
)

// ChoreStore is an implementation of the domain.ChoreStore interface, using the storage package
type ChoreStore struct {
	store storage.Store
}

// NewChoreStore returns a new ChoreStore.
func NewChoreStore(s storage.Store) ChoreStore {
	return ChoreStore{s}
}

// SaveChore saves a Chore to the database, under the account with the ID aid.
func (c ChoreStore) SaveChore(aid string, ch domain.Chore) error {
	b, err := c.accountChoresBucket(aid)
	if err != nil {
		return err
	}

	err = b.Put(ch.ID, ch)
	return err
}

// GetChore returns a chore with id cid, owned by account with ID aid.
func (c ChoreStore) GetChore(aid, cid string) (domain.Chore, error) {
	b, err := c.accountChoresBucket(aid)

	ch := domain.Chore{}
	err = b.Get(cid, &ch)
	if err == storage.NotFoundError {
		return ch, domain.ErrChoreNotFound
	}

	return ch, err
}

// GetAccountChores returns all the chores from account with ID aid
func (c ChoreStore) GetAccountChores(aid string) ([]domain.Chore, error) {
	b, err := c.accountChoresBucket(aid)

	if err != nil{
		return []domain.Chore{}, err
	}

	cs := make([]domain.Chore, 0)
	err = b.Each(func(k string) error {
		ch := domain.Chore{}
		err := b.Get(k, &ch)
		if err == nil {
			cs = append(cs, ch)
		}
		return err
	})

	return cs, err
}

// DeleteChore deletes a Chore owned by account aid, with ID cid
func (c ChoreStore) DeleteChore(aid, cid string) (err error) {
	b, err := c.accountChoresBucket(aid)
	if err == nil {
		err = b.Delete(cid)
	}
	return
}

// DeleteAccountChores deletes all chores from account ID aid
func (c ChoreStore) DeleteAccountChores(aid string) error {
	return c.store.DeleteBucket("chores_" + aid)
}

func (c ChoreStore) accountChoresBucket(aid string) (storage.Bucket, error) {
	return c.store.Bucket("chores_" + aid)
}
