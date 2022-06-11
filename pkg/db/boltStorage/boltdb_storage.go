package boltStorage

import (
	"github.com/boltdb/bolt"
)

type BoltdbStorage struct {
	db   *bolt.DB
	path string
}

func Open(path string) (*BoltdbStorage, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &BoltdbStorage{
		db:   db,
		path: path,
	}, nil
}

// CreateBucket 创建一个桶
func (s *BoltdbStorage) CreateBucket(bucketName []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucketName)
		return err
	})
}

// DeleteBucket 删除一个桶
func (s *BoltdbStorage) DeleteBucket(bucketName []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucketName)
	})
}

// CreateBucketIfNotExist 如果桶不存在则创建
func (s *BoltdbStorage) CreateBucketIfNotExist(bucketName []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

// Get -> if key not exist will return nil.
func (s *BoltdbStorage) Get(key []byte, bucketName []byte) ([]byte, bool) {
	var buffer []byte

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		buffer = b.Get(key)
		return nil
	})

	if err != nil {
		return nil, false
	}
	return buffer, true
}

func (s *BoltdbStorage) Set(key []byte, value []byte, bucketName []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		err := b.Put(key, value)
		return err
	})

	return err
}

// Delete 删除
func (s *BoltdbStorage) Delete(key []byte, bucketName []byte) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Delete(key)
	})

	return err
}

// Close 关闭
func (s *BoltdbStorage) Close() error {
	return s.db.Close()
}
