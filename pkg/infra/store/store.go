package store

import (
	"fmt"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/log"
	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/infra/store/boltdb"
)

type storeType int

const (
	BoltDB storeType = iota
)

type DatabaseManager interface {
	Migrate() error

	// SetServiceName is a generic function used to create a bucket inside a database.
	SetServiceName(bucketName string) error

	// GetObject is a generic function used to retrieve an unmarshalled object from a database.
	GetObject(bucketName string, key []byte, object interface{}) error

	// UpdateObject is a generic function used to update an object inside a database.
	UpdateObject(bucketName string, key []byte, object interface{}) error

	// DeleteObject is a generic function used to delete an object inside a database.
	DeleteObject(bucketName string, key []byte) error

	// DeleteAllObjects delete all objects where matching() returns (id, ok).
	// TODO: think about how to return the error inside (maybe change ok to type err, and use "notfound"?
	DeleteAllObjects(bucketName string, obj interface{}, matchingFn func(o interface{}) (id int, ok bool)) error

	// CreateObject creates a new object in the bucket, using the next bucket sequence id
	CreateObject(bucketName string, objFn func(uint64) (int, interface{})) error

	// CreateObjectWithId creates a new object in the bucket, using the specified id
	CreateObjectWithId(bucketName string, id int, obj interface{}) error

	// CreateObjectWithStringId creates a new object in the bucket, using the specified id
	CreateObjectWithStringId(bucketName string, id []byte, obj interface{}) error

	// GetAll gets all objects from the bucket, using the jsoniter library to unmarshal them.
	GetAll(bucketName string, obj interface{}, appendFn func(o interface{}) (interface{}, error)) error

	// GetAllWithKeyPrefix gets all objects from the bucket, starting with the given key prefix.
	GetAllWithKeyPrefix(bucketName string, keyPrefix []byte, obj interface{}, appendFn func(o interface{}) (interface{}, error)) error

	// ConvertToKey returns an 8-byte big endian representation of v.
	// This function is typically used for encoding integer IDs to byte slices
	// so that they can be used as BoltDB keys.
	ConvertToKey(v int) []byte
}

type Options struct {
	Type       storeType
	BoltDBOpts *boltdb.Options
}

func New(logger *log.Logger, options *Options) (db DatabaseManager, err error) {
	switch options.Type {
	case BoltDB:
		if options.BoltDBOpts == nil {
			return nil, fmt.Errorf("boltdb storage options can not be empty")
		}

		db, err = boltdb.New(logger, options.BoltDBOpts)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", options.Type, err)
		}

	default:
		return nil, fmt.Errorf("unknown storage type: %s", options.Type)
	}

	return db, nil
}

func (s storeType) String() string {
	return [...]string{"boltdb"}[s]
}
