package boltdb

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// MarshalObject encodes an object to binary format
func (d *Database) MarshalObject(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return data, err
	}

	if d.isEncrypted {
		return encrypt(data, d.encryptionKey)
	}

	return data, nil
}

// UnmarshalObject decodes an object from binary data
// using the jsoniter library. It is mainly used to accelerate environment(endpoint)
// decoding at the moment.
func (d *Database) UnmarshalObject(data []byte, object interface{}) error {
	var (
		err           error
		decryptedData []byte
	)

	if d.isEncrypted {
		decryptedData, err = decrypt(data, d.encryptionKey)
		if err != nil {
			return fmt.Errorf("decrypt object: %v", err)
		}
	}

	if err = json.Unmarshal(decryptedData, &object); err != nil {
		if s, ok := object.(*string); ok {
			*s = string(decryptedData)

			return nil
		}

		return fmt.Errorf("unmarshal: %v", err)
	}

	return nil
}
