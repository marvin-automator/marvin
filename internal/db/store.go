package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/marvin-automator/marvin/internal/config"
	"time"
)

var db *badger.DB

func GetStore(name string) Store {
	if db == nil {
		opts := badger.DefaultOptions
		opts.Dir = config.DataDir
		opts.ValueDir = opts.Dir
		var err error
		db, err = badger.Open(opts)

		if err != nil {
			panic(err)
		}
	}

	return Store{name}
}

type Store struct {
	name string
}

func (s Store) makeKey(k string) []byte {
	return []byte(fmt.Sprintf("%v_%v", s.name, k))
}

func (s Store) dbKeyToStorekey(dbKey []byte) string {
	return string(dbKey[len(s.name)+1:])
}

func (s Store) Get(key string, ptr interface{}) error {
	var v []byte
	err := db.View(func(tx *badger.Txn) error {
		item, err := tx.Get(s.makeKey(key))

		if err != nil {
			if err == badger.ErrKeyNotFound {
				return KeyNotFoundError{key}
			}
			return err
		}

		v, err = item.ValueCopy(v)
		return err
	})

	if err != nil {
		return err
	}

	d := gob.NewDecoder(bytes.NewReader(v))
	return d.Decode(ptr)
}

func (s Store) encodeValue(v interface{}) ([]byte, error) {
	b := bytes.NewBuffer([]byte{})
	e := gob.NewEncoder(b)
	err := e.Encode(v)
	return b.Bytes(), err
}

func (s Store) Set(key string, value interface{}) error {
	data, err := s.encodeValue(value)
	if err != nil {
		return err
	}

	return db.Update(func(tx *badger.Txn) error {
		return tx.Set(s.makeKey(key), data)
	})
}

func (s Store) SetWithExpiration(key string, value interface, expires time.Duration) error {
	data, err := s.encodeValue(value)
	if err != nil {
		return err
	}

	return db.Update(func(tx *badger.Txn) error {
		return tx.SetWithTTL(s.makeKey(key), data, expires)
	})
}

func (s Store) Delete(key string) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete(s.makeKey(key))
	})
}

func (s Store) EachKeyWithPrefix(prefix string, ptr interface{}, f func(key string) error) error {
	return db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		bprefix := s.makeKey(prefix)
		for it.Seek(bprefix); it.ValidForPrefix(bprefix); it.Next() {
			item := it.Item()

			var val []byte
			val, err := item.ValueCopy(val)
			if err != nil {
				return err
			}

			dec := gob.NewDecoder(bytes.NewBuffer(val))
			err = dec.Decode(ptr)
			if err != nil {
				return err
			}

			err = f(s.dbKeyToStorekey(item.Key()))

			if err != nil {
				return err
			}
		}
		return nil
	})
}
