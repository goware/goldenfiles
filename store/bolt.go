package store

import (
	"github.com/boltdb/bolt"
)

type Bolt struct {
	DB *bolt.DB
}

const bucketName = "mockingbird"

func (b *Bolt) bucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	return tx.CreateBucketIfNotExists([]byte(bucketName))
}

func (b *Bolt) Get(key string) ([]byte, error) {
	var v []byte
	err := b.DB.Update(func(tx *bolt.Tx) error {
		k, err := b.bucket(tx)
		if err != nil {
			return err
		}
		v = k.Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (b *Bolt) Set(key string, val []byte) error {
	err := b.DB.Update(func(tx *bolt.Tx) error {
		k, err := b.bucket(tx)
		if err != nil {
			return err
		}
		return k.Put([]byte(key), val)
	})
	return err
}

func (b *Bolt) Delete(key string) error {
	err := b.DB.Update(func(tx *bolt.Tx) error {
		k, err := b.bucket(tx)
		if err != nil {
			return err
		}
		return k.Delete([]byte(key))
	})
	return err
}
