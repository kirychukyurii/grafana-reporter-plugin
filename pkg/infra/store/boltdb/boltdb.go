package boltdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/apperrors"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"math"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	DatabaseFileName          = "reporter.db"
	EncryptedDatabaseFileName = "reporter.edb"
)

type Options struct {
	DataDirectory   string
	EncryptionKey   []byte
	Timeout         int
	InitialMmapSize int
	MaxBatchSize    int
	MaxBatchDelay   int
}

type Database struct {
	*bolt.DB

	isEncrypted   bool
	encryptionKey []byte
}

// New opens and initializes the BoltDB database.
func New(logger *log.Logger, options *Options) (*Database, error) {
	var db Database

	if options.EncryptionKey != nil {
		db.isEncrypted = true
		db.encryptionKey = options.EncryptionKey
	}

	databasePath := filepath.Join(options.DataDirectory, db.DatabaseFileName())
	database, err := bolt.Open(databasePath, 0600, &bolt.Options{
		Timeout:         time.Duration(options.Timeout) * time.Second,
		InitialMmapSize: options.InitialMmapSize,
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %v", err)
	}

	database.MaxBatchSize = options.MaxBatchSize
	database.MaxBatchDelay = time.Duration(options.MaxBatchDelay) * time.Second
	db.DB = database

	return &db, nil
}

// Close closes the BoltDB database.
// Safe to being called multiple times.
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}

	return nil
}

// DatabaseFileName get the database filename
func (d *Database) DatabaseFileName() string {
	if d.isEncrypted {
		return EncryptedDatabaseFileName
	}

	return DatabaseFileName
}

// SetServiceName is a generic function used to create a bucket inside a database.
func (d *Database) SetServiceName(bucketName string) error {
	fn := func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		return err
	}

	return d.Update(fn)
}

// GetObject is a generic function used to retrieve an unmarshalled object from a database.
func (d *Database) GetObject(bucketName string, key []byte, object interface{}) error {
	fn := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		value := bucket.Get(key)
		if value == nil {
			return fmt.Errorf("%w (bucket=%s, key=%s)", apperrors.ErrObjectNotFound, bucketName, keyToString(key))
		}

		return d.UnmarshalObject(value, object)
	}

	return d.View(fn)
}

// UpdateObject is a generic function used to update an object inside a database.
func (d *Database) UpdateObject(bucketName string, key []byte, object interface{}) error {
	fn := func(tx *bolt.Tx) error {
		data, err := d.MarshalObject(object)
		if err != nil {
			return err
		}

		bucket := tx.Bucket([]byte(bucketName))

		return bucket.Put(key, data)
	}

	return d.Update(fn)
}

// DeleteObject is a generic function used to delete an object inside a database.
func (d *Database) DeleteObject(bucketName string, key []byte) error {
	fn := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		return bucket.Delete(key)
	}

	return d.Update(fn)
}

// DeleteAllObjects delete all objects where matching() returns (id, ok).
// TODO: think about how to return the error inside (maybe change ok to type err, and use "notfound"?
func (d *Database) DeleteAllObjects(bucketName string, obj interface{}, matchingFn func(o interface{}) (id int, ok bool)) error {
	fn := func(tx *bolt.Tx) error {
		var ids []int

		bucket := tx.Bucket([]byte(bucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if err := d.UnmarshalObject(v, &obj); err != nil {
				return err
			}

			if id, ok := matchingFn(obj); ok {
				ids = append(ids, id)
			}
		}

		for _, id := range ids {
			if err := bucket.Delete(d.ConvertToKey(id)); err != nil {
				return err
			}
		}

		return nil
	}

	return d.Update(fn)
}

// CreateObject creates a new object in the bucket, using the next bucket sequence id
func (d *Database) CreateObject(bucketName string, objFn func(uint64) (int, interface{})) error {
	fn := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		seqId, _ := bucket.NextSequence()
		id, obj := objFn(seqId)

		data, err := d.MarshalObject(obj)
		if err != nil {
			return err
		}

		return bucket.Put(d.ConvertToKey(id), data)
	}

	return d.Update(fn)
}

// CreateObjectWithId creates a new object in the bucket, using the specified id
func (d *Database) CreateObjectWithId(bucketName string, id int, obj interface{}) error {
	fn := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		data, err := d.MarshalObject(obj)
		if err != nil {
			return err
		}

		return bucket.Put(d.ConvertToKey(id), data)
	}

	return d.Update(fn)
}

// CreateObjectWithStringId creates a new object in the bucket, using the specified id
func (d *Database) CreateObjectWithStringId(bucketName string, id []byte, obj interface{}) error {
	fn := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		data, err := d.MarshalObject(obj)

		if err != nil {
			return err
		}

		return bucket.Put(id, data)
	}

	return d.Update(fn)
}

func (d *Database) GetAll(bucketName string, obj interface{}, appendFn func(o interface{}) (interface{}, error)) error {
	fn := func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))

		return bucket.ForEach(func(k []byte, v []byte) error {
			err := d.UnmarshalObject(v, obj)
			if err == nil {
				obj, err = appendFn(obj)
			}

			return err
		})
	}

	return d.View(fn)
}

func (d *Database) GetAllWithKeyPrefix(bucketName string, keyPrefix []byte, obj interface{}, appendFn func(o interface{}) (interface{}, error)) error {
	fn := func(tx *bolt.Tx) error {
		cursor := tx.Bucket([]byte(bucketName)).Cursor()

		for k, v := cursor.Seek(keyPrefix); k != nil && bytes.HasPrefix(k, keyPrefix); k, v = cursor.Next() {
			err := d.UnmarshalObject(v, obj)
			if err != nil {
				return err
			}

			obj, err = appendFn(obj)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return d.View(fn)
}

// ConvertToKey returns an 8-byte big endian representation of v.
// This function is typically used for encoding integer IDs to byte slices
// so that they can be used as BoltDB keys.
func (d *Database) ConvertToKey(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))

	return b
}

// keyToString Converts a key to a string value suitable for logging
func keyToString(b []byte) string {
	if len(b) != 8 {
		return string(b)
	}

	v := binary.BigEndian.Uint64(b)
	if v <= math.MaxInt32 {
		return fmt.Sprintf("%d", v)
	}

	return string(b)
}
