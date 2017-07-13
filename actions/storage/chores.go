package storage

import (
	"github.com/bigblind/marvin/actions/domain"
	"github.com/bigblind/marvin/storage"
	"github.com/boltdb/bolt"
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
	b, err := c.getOrCreateAccountChoresBucket(aid)
	if err != nil {
		return err
	}

	cb, err := c.store.EncodeBytes(ch)
	if err != nil {
		return err
	}

	err = b.Put([]byte(ch.ID), cb)
	return err
}

// GetChore returns a chore with id cid, owned by account with ID aid.
func (c ChoreStore) GetChore(aid, cid string) (ch domain.Chore, err error) {
	b, err := c.getOrCreateAccountChoresBucket(aid)
	if err == nil {
		cb := b.Get([]byte(cid))
		if cb == nil {
			err = domain.ErrChoreNotFound
			return
		}
		err = c.store.DecodeBytes(&ch, cb)
	}
	return
}

// GetAccountChores returns all the chores from account with ID aid
func (c ChoreStore) GetAccountChores(aid string) (cs []domain.Chore, err error) {
	b, err := c.getOrCreateAccountChoresBucket(aid)
	if err == nil {
		err = b.ForEach(func(k, v []byte) error {
			ch := domain.Chore{}
			err = c.store.DecodeBytes(&ch, v)
			if err == nil {
				cs = append(cs, ch)
			}
			return err
		})
	}
	return
}

// DeleteChore deletes a Chore owned by account aid, with ID cid
func (c ChoreStore) DeleteChore(aid, cid string) (err error) {
	b, err := c.getOrCreateAccountChoresBucket(aid)
	if err == nil {
		err = b.Delete([]byte(cid))
	}
	return
}

// DeleteAccountChores deletes all chores from account ID aid
func (c ChoreStore) DeleteAccountChores(aid string) error {
	return c.store.Tx.DeleteBucket([]byte("chores_" + aid))
}

func (c ChoreStore) getOrCreateAccountChoresBucket(aid string) (*bolt.Bucket, error) {
	return c.store.CreateBucketIfNotExists("chores_" + aid)
}
