package awssecretsmanager

import (
	"log"
	"testing"

	utilstrings "github.com/bertjohnson/util/strings"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestGetAutoCertCache tests GetAutoCertCache().
func TestGetAutoCertCache(t *testing.T) {
	cache := awsSecretsManager.GetAutoCertCache(ctx)

	// Create the secret.
	secretBytes := []byte("abcdefg")
	secretID, err := utilstrings.HexToBase58(uuid.New().String())
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = cache.Put(ctx, secretID, secretBytes)
	assert.NoError(t, err)

	// Read the secret.
	var readBytes []byte
	readBytes, err = cache.Get(ctx, secretID)
	assert.NoError(t, err)
	assert.Equal(t, secretBytes, readBytes)

	// Delete the secret.
	err = cache.Delete(ctx, secretID)
	assert.NoError(t, err)
}
