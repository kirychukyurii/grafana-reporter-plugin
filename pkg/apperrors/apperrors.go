package apperrors

import (
	"errors"
)

var (
	ErrBadObject = errors.New("")

	ErrObjectNotFound = errors.New("object not found")

	ErrHaveEncryptedAndUnencrypted = errors.New("detected both an encrypted and un-encrypted database and cannot start: only one database should exist")
	ErrHaveEncryptedWithNoKey      = errors.New("database is encrypted, but no secret was loaded")

	ErrEncryptedStringTooShort = errors.New("encrypted string too short")

	// ErrStop signals the stop of computation when filtering results
	ErrStop = errors.New("stop")
)
