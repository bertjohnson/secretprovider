package secretprovider

import (
	"context"
	"testing"

	"github.com/bertjohnson/logger"
	secretprovidertype "github.com/bertjohnson/secretprovider/types"
	"github.com/bertjohnson/startup"
	"github.com/stretchr/testify/assert"
)

// TestGet tests Get().
func TestGet(t *testing.T) {
	// Declare that the configuration is ready.
	ctx := context.Background()
	err := startup.Ready()
	if err != nil {
		logger.Fatal(ctx, "Error loading configuration values: "+err.Error())
	}

	// Wait for logger.
	logger.Wait(ctx)

	// Get invalid Vault client.
	_, err = Get(ctx, &secretprovidertype.SecretProvider{
		Type: "Vault",
	})
	assert.Error(t, err)

	// Get invalid local files client.
	_, err = Get(ctx, &secretprovidertype.SecretProvider{
		Type: "LocalFiles",
	})
	assert.Error(t, err)

	// Get valid local files client.
	_, err = Get(ctx, &secretprovidertype.SecretProvider{
		Type: "LocalFiles",
		URI:  "test",
	})
	assert.NoError(t, err)
}
