package types

import (
	"context"
)

// AutoCertCache provides an implementation of an autocert cache.
type AutoCertCache interface {
	// Get returns a certificate data for the specified key.
	// If there's no such key, Get returns autocert.ErrCacheMiss.
	Get(ctx context.Context, key string) ([]byte, error)

	// Put stores the data in the cache under the specified key.
	// Underlying implementations may use any data storage format,
	// as long as the reverse operation, Get, results in the original data.
	Put(ctx context.Context, key string, data []byte) error

	// Delete removes a certificate data from the cache under the specified key.
	// If there's no such key in the cache, Delete returns nil.
	Delete(ctx context.Context, key string) error
}
