package bolt

import (
	"time"

	"github.com/boltdb/bolt"
	"github.com/goware/mockingbird/store"
)

const defaultBucketName = "mockingbird"

type Bolt struct {
	DB         *bolt.DB
	BucketName string
}

func NewStore(name string) *Bolt {
	db, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		panic(err.Error())
	}

	return &Bolt{
		DB: db,
	}
}

func (b *Bolt) bucket(tx *bolt.Tx) (*bolt.Bucket, error) {
	bucketName := b.BucketName
	if bucketName == "" {
		bucketName = defaultBucketName
	}
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

var _ = store.Store(&Bolt{})
