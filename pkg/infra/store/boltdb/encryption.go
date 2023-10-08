package boltdb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/kirychukyurii/grafana-reporter-plugin/pkg/apperrors"
)

// don't have a KMS .... aes GCM seems the most likely from
// https://gist.github.com/atoponce/07d8d4c833873be2f68c34f9afc5a78a#symmetric-encryption

func encrypt(plaintext []byte, passphrase []byte) ([]byte, error) {
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return nil, fmt.Errorf("create cypher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(encrypted []byte, passphrase []byte) ([]byte, error) {
	if string(encrypted) == "false" {
		return []byte("false"), nil
	}
	block, err := aes.NewCipher(passphrase)
	if err != nil {
		return nil, fmt.Errorf("create cypher block: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, apperrors.ErrEncryptedStringTooShort
	}

	nonce, ciphertextByteClean := encrypted[:nonceSize], encrypted[nonceSize:]
	plaintextByte, err := gcm.Open(nil, nonce, ciphertextByteClean, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt text: %v", err)
	}

	return plaintextByte, err
}
