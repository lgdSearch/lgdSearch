package badgerStorage

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"sync"
)

type BadgerStorage struct {
	db   *badger.DB
	path string
	sync.Mutex
}

// OpenWithDefault open a badgerDB with default setting
func OpenWithDefault(path string) *BadgerStorage {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		log.Fatal(err)
	}
	return &BadgerStorage{
		db:   db,
		path: path,
	}
}

func Open(option badger.Options) *BadgerStorage {
	db, err := badger.Open(option)
	if err != nil {
		log.Fatal(err)
	}
	return &BadgerStorage{
		db:   db,
		path: option.Dir,
	}
}

// Get if not find bool will return false
func (b *BadgerStorage) Get(key []byte) ([]byte, bool) {
	result := make([]byte, 0)
	err := b.db.View(func(txn *badger.Txn) error {
		value, err := txn.Get(key)
		if err != nil {
			return err
		}
		return value.Value(func(val []byte) error {
			result = append(result, val...)
			return nil
		})
	})
	if err != nil {
		return nil, false
	}
	return result, true
}

func (b *BadgerStorage) Set(key []byte, value []byte) error {
	b.Lock()
	defer b.Unlock()
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

func (b *BadgerStorage) Delete(key []byte) error {
	b.Lock()
	defer b.Unlock()
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

func (b *BadgerStorage) Close() error {
	return b.db.Close()
}
