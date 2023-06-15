package localfiles

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"golang.org/x/crypto/acme/autocert"
)

// AutoCertCache implements AutoCertCache using a local directory.
type AutoCertCache string

// GetAutoCertCache returns an autocert-compatible cache.
func (l *LocalFiles) GetAutoCertCache(ctx context.Context) secretprovidertype.AutoCertCache {
	return AutoCertCache(filepath.Join(l.basePath, "/autocert"))
}

// Get reads a certificate data from the specified file name.
func (a AutoCertCache) Get(ctx context.Context, name string) ([]byte, error) {
	// Validate parameters.
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	if strings.Contains(name, "/") {
		return nil, errors.New("name cannot traverse directories")
	}

	name = filepath.Join(string(a), name)
	var (
		data []byte
		err  error
		done = make(chan struct{})
	)
	go func() {
		data, err = ioutil.ReadFile(name) // #nosec G304
		close(done)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
	}
	if os.IsNotExist(err) {
		return nil, autocert.ErrCacheMiss
	}
	return data, err
}

// Put writes the certificate data to the specified file name.
// The file will be created with 0600 permissions.
func (a AutoCertCache) Put(ctx context.Context, name string, data []byte) error {
	// Validate parameters.
	if ctx == nil {
		return errors.New("context is required")
	}
	if strings.Contains(name, "/") {
		return errors.New("name cannot traverse directories")
	}

	if err := os.MkdirAll(string(a), 0700); err != nil {
		return err
	}

	done := make(chan struct{})
	var err error
	go func() {
		defer close(done)
		var tmp string
		if tmp, err = a.writeTempFile(name, data); err != nil {
			return
		}
		select {
		case <-ctx.Done():
			// Don't overwrite the file if the context was canceled.
		default:
			newName := filepath.Join(string(a), name)
			err = os.Rename(tmp, newName)
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	return err
}

// Delete removes the specified file name.
func (a AutoCertCache) Delete(ctx context.Context, name string) error {
	name = filepath.Join(string(a), name)
	var (
		err  error
		done = make(chan struct{})
	)
	go func() {
		err = os.Remove(name)
		close(done)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// writeTempFile writes b to a temporary file, closes the file and returns its path.
func (a AutoCertCache) writeTempFile(prefix string, b []byte) (string, error) {
	// TempFile uses 0600 permissions
	f, err := ioutil.TempFile(string(a), prefix)
	if err != nil {
		return "", err
	}
	if _, err := f.Write(b); err != nil {
		err = f.Close()
		return "", err
	}
	return f.Name(), f.Close()
}
