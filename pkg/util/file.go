package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Create one file
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

func OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func CopyReader(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
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

func TemporaryDir(dir string) (string, error) {
	tmpDir := filepath.Join(dir, NewUUID().String())
	if err := EnsureDirRW(tmpDir); err != nil {
		return "", fmt.Errorf("ensure dir RW: %v", err)
	}

	return tmpDir, nil
}

func ReadDir(dir string) ([]string, error) {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		files = append(files, filepath.Join(dir, e.Name()))
	}

	return files, nil
}
