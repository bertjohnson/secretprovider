package localfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetAutoCertCache tests GetAutoCertCache().
func TestGetAutoCertCache(t *testing.T) {
	cache := localFilesClient.GetAutoCertCache(ctx)

	// Create the secret.
	secretBytes := []byte("abcdefg")
	err := cache.Put(ctx, "secret1", secretBytes)
	assert.NoError(t, err)

	// Read the secret.
	var readBytes []byte
	readBytes, err = cache.Get(ctx, "secret1")
	assert.NoError(t, err)
	assert.Equal(t, secretBytes, readBytes)

	// Delete the secret.
	err = cache.Delete(ctx, "secret1")
	assert.NoError(t, err)
}
