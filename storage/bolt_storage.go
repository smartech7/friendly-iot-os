package storage

import (
	"time"

	"github.com/boltdb/bolt"
)

var gwk_documents = []byte("gwk_documents")

type boltStorage struct {
	db *bolt.DB
}

func openBoltStorage(dbpath string) (Storage, error) {
	db, err := bolt.Open(dbpath, 0600, &bolt.Options{Timeout: 3600 * time.Second})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(gwk_documents)
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return &boltStorage{db}, nil
}

func (s *boltStorage) WALName() string {
	return s.db.Path()
}

func (s *boltStorage) Set(k []byte, v []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(gwk_documents).Put(k, v)
	})
}

func (s *boltStorage) Get(k []byte) (b []byte, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(gwk_documents).Get(k)
		return nil
	})
	return
}

func (s *boltStorage) Delete(k []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(gwk_documents).Delete(k)
	})
}

func (s *boltStorage) ForEach(fn func(k, v []byte) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(gwk_documents)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *boltStorage) Close() error {
	return s.db.Close()
}
