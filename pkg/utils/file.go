package utils

import (
	"fmt"
	"os"
	"time"
)

// Create one file
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

// EnsureDir mkdir dir if not exist
func EnsureDir(fp string) error {
	return os.MkdirAll(fp, os.ModePerm)
}

// Close fd
func Close(fd *os.File) error {
	return fd.Close()
}

// Remove one file
func Remove(name string) error {
	return os.Remove(name)
}

// EnsureDirRW ensure the datadir and make sure it's rw-able
func EnsureDirRW(dataDir string) error {
	if err := EnsureDir(dataDir); err != nil {
		return err
	}

	checkFile := fmt.Sprintf("%s/rw.%d", dataDir, time.Now().UnixNano())
	fd, err := Create(checkFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("open %s: rw permission denied", dataDir)
		}

		return err
	}

	if err = Close(fd); err != nil {
		return fmt.Errorf("close error: %s", err)
	}

	if err = Remove(checkFile); err != nil {
		return fmt.Errorf("remove error: %s", err)
	}

	return nil
}
