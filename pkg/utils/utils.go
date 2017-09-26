package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Store stores string read from `body` into file at `dstfile`.
func Store(dstfile string, body io.ReadCloser, override bool) error {
	dstfile = filepath.Clean(dstfile)
	tmpfile := filepath.Join(filepath.Dir(dstfile), fmt.Sprintf(".%s.tmp", filepath.Base(dstfile)))
	// Check file exists or not.
	if !override {
		if _, err := os.Stat(tmpfile); !os.IsNotExist(err) {
			return fmt.Errorf("%s does exist", tmpfile)
		}
		if _, err := os.Stat(dstfile); !os.IsNotExist(err) {
			return fmt.Errorf("%s does exist", dstfile)
		}
	}
	// Open.
	f, err := os.OpenFile(tmpfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Copy.
	if _, err := io.Copy(f, body); err != nil {
		os.Remove(tmpfile)
		return err
	}
	// Sync.
	if err := f.Sync(); err != nil {
		os.Remove(tmpfile)
		return err
	}
	// Rename.
	if err := os.Rename(tmpfile, dstfile); err != nil {
		os.Remove(tmpfile)
		return err
	}
	return nil
}
